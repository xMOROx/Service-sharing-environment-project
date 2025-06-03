package internal

import (
	"context"
	"io"
	"sync"
	"time"

	invpb "Service-sharing-environment-project/proto/inventory"
	orderpb "Service-sharing-environment-project/proto/order"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type OrderServer struct {
    orderpb.UnimplementedOrderServiceServer
    inventory      invpb.InventoryServiceClient
    sessions       map[string][]*invpb.OrderItemRequest
    mu             sync.Mutex
    requestCounter metric.Int64Counter
    latencyHist    metric.Float64Histogram
}

func NewOrderServer(invClient invpb.InventoryServiceClient, m metric.Meter) *OrderServer {
    ctr, err := m.Int64Counter("order_requests_total")
    if err != nil {
        panic(err)
    }
    hist, err := m.Float64Histogram("order_request_latency_ms")
    if err != nil {
        panic(err)
    }

    return &OrderServer{
        inventory:      invClient,
        sessions:       make(map[string][]*invpb.OrderItemRequest),
        requestCounter: ctr,
        latencyHist:    hist,
    }
}

// instrument otwiera span i po zakończeniu rejestruje liczniki i histogram
func (s *OrderServer) instrument(ctx context.Context, name string) (context.Context, func()) {
    start := time.Now()
    ctx, span := otel.Tracer("order-service").Start(ctx, name)
    return ctx, func() {
        span.End()
        ms := float64(time.Since(start).Milliseconds())
        s.requestCounter.Add(ctx, 1,
            metric.WithAttributes(attribute.String("rpc", name)),
        )
        s.latencyHist.Record(ctx, ms,
            metric.WithAttributes(attribute.String("rpc", name)),
        )
    }
}

func (s *OrderServer) CheckItemAvailability(ctx context.Context, req *invpb.ProductId) (*invpb.ProductInfo, error) {
    ctx, end := s.instrument(ctx, "CheckItemAvailability")
    defer end()
    return s.inventory.GetProductInfo(ctx, req)
}

func (s *OrderServer) BuildOrder(stream orderpb.OrderService_BuildOrderServer) error {
    tracer := otel.Tracer("order-service")

    for {
        // 1) Span wokół Recv
        ctxRecv, spanRecv := tracer.Start(stream.Context(), "ReceiveOrderRequest")
        req, err := stream.Recv()
        spanRecv.End()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        // 2) Span & metryki wokół logiki
        ctx, end := s.instrument(ctxRecv, "BuildOrder")
        s.mu.Lock()
        s.sessions[req.SessionId] = append(s.sessions[req.SessionId], req)
        s.mu.Unlock()

        // 3) Sprawdzenie stanu magazynowego
        invResp, err := s.inventory.GetProductInfo(ctx, &invpb.ProductId{ProductId: req.ProductId})
        if err != nil {
            end()
            return err
        }
        available := invResp.AvailableQuantity >= req.RequestedQuantity

        // 4) Zwracamy odpowiedni typ z inventory (pole AvailableQuantity)
        resp := &invpb.OrderItemResponse{
            ProductId:         req.ProductId,
            Available:         available,
            AvailableQuantity: invResp.AvailableQuantity,
        }
        if err := stream.Send(resp); err != nil {
            end()
            return err
        }
        end()
    }
}

func (s *OrderServer) FinalizeOrder(ctx context.Context, req *orderpb.FinalizeOrderRequest) (*orderpb.FinalizeOrderResponse, error) {
    ctx, end := s.instrument(ctx, "FinalizeOrder")
    defer end()

    var results []*orderpb.ItemResult
    okAll := true

    for _, item := range req.GetItems() {
        prod, err := s.inventory.GetProductInfo(ctx, &invpb.ProductId{ProductId: item.ProductId})
        res := &orderpb.ItemResult{ProductId: item.ProductId}

        if err != nil || prod.GetAvailableQuantity() < item.GetQuantity() {
            res.Reserved = false
            res.Message = "Insufficient stock"
            okAll = false
        } else {
            _, err := s.inventory.AdjustStock(ctx, &invpb.StockAdjustment{
                ProductId:      item.ProductId,
                QuantityChange: -item.Quantity,
            })
            if err != nil {
                res.Reserved = false
                res.Message = "Reservation failed"
                okAll = false
            } else {
                res.Reserved = true
                res.Message = "Reserved"
            }
        }
        results = append(results, res)
    }

    s.mu.Lock()
    delete(s.sessions, req.SessionId)
    s.mu.Unlock()

    msg := "Order finalized"
    if !okAll {
        msg = "One or more items failed"
    }
    return &orderpb.FinalizeOrderResponse{
        Success:     okAll,
        Message:     msg,
        ItemResults: results,
    }, nil
}

func (s *OrderServer) ConfirmOrderStock(ctx context.Context, req *orderpb.FinalizeOrderRequest) (*invpb.OperationStatus, error) {
    ctx, end := s.instrument(ctx, "ConfirmOrderStock")
    defer end()

    for _, item := range req.GetItems() {
        _, err := s.inventory.AdjustStock(ctx, &invpb.StockAdjustment{
            ProductId:      item.ProductId,
            QuantityChange: -item.Quantity,
        })
        if err != nil {
            return &invpb.OperationStatus{Success: false, Message: "Stock error"}, nil
        }
    }
    return &invpb.OperationStatus{Success: true, Message: "Stock confirmed"}, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
    _, end := s.instrument(ctx, "CancelOrder")
    defer end()

    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.sessions[req.SessionId]; !ok {
        return &orderpb.CancelOrderResponse{Released: false, Message: "Session not found"}, nil
    }
    delete(s.sessions, req.SessionId)
    return &orderpb.CancelOrderResponse{Released: true, Message: "Order cancelled and reservations released"}, nil
}
