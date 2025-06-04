package internal

import (
	"context"
	"io"
	"log"
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
    
	log.Printf("[Order][CheckItemAvailability] product_id=%s", req.ProductId)
	resp, err := s.inventory.GetProductInfo(ctx, req)
	if err != nil {
		log.Printf("[Order][CheckItemAvailability] inventory error: %v", err)
		return nil, err
	}
	log.Printf("[Order][CheckItemAvailability] inventory success: available_quantity=%d", resp.AvailableQuantity)
	return resp, nil
}

func (s *OrderServer) BuildOrder(stream orderpb.OrderService_BuildOrderServer) error {
    tracer := otel.Tracer("order-service")

    for {
        // 1) Span wokół Recv
        ctxRecv, spanRecv := tracer.Start(stream.Context(), "ReceiveOrderRequest")
        req, err := stream.Recv()
        spanRecv.End()
        if err == io.EOF {
            log.Printf("[Order][BuildOrder] stream EOF")
            return nil
        }
        if err != nil {
			log.Printf("[Order][BuildOrder] Recv error: %v", err)
            return err
        }

        log.Printf(
			"[Order][BuildOrder] received: session_id=%s product_id=%s requested_quantity=%d",
			req.SessionId, req.ProductId, req.RequestedQuantity,
		)

        // 2) Span & metryki wokół logiki
        ctx, end := s.instrument(ctxRecv, "BuildOrder")
        s.mu.Lock()
        s.sessions[req.SessionId] = append(s.sessions[req.SessionId], req)
        s.mu.Unlock()
		log.Printf("[Order][BuildOrder] appended to session %s", req.SessionId)

        // 3) Sprawdzenie stanu magazynowego
		log.Printf("[Order][BuildOrder] calling Inventory.GetProductInfo for product_id=%s", req.ProductId)
        invResp, err := s.inventory.GetProductInfo(ctx, &invpb.ProductId{ProductId: req.ProductId})
        if err != nil {
            log.Printf("[Order][BuildOrder] Inventory.GetProductInfo error: %v", err)
            end()
            return err
        }
        log.Printf("[Order][BuildOrder] Inventory.GetProductInfo returned available_quantity=%d", invResp.AvailableQuantity)
        available := invResp.AvailableQuantity >= req.RequestedQuantity

        // 4) Zwracamy odpowiedni typ z inventory (pole AvailableQuantity)
        resp := &invpb.OrderItemResponse{
            ProductId:         req.ProductId,
            Available:         available,
            AvailableQuantity: invResp.AvailableQuantity,
        }
        log.Printf(
			"[Order][BuildOrder] sending response: product_id=%s available=%v remaining=%d",
			resp.ProductId, resp.Available, resp.AvailableQuantity,
		)
        if err := stream.Send(resp); err != nil {
            log.Printf("[Order][BuildOrder] Send error: %v", err)
            end()
            return err
        }
		log.Printf("[Order][BuildOrder] stream.Send successful")
        end()
    }
}

func (s *OrderServer) FinalizeOrder(ctx context.Context, req *orderpb.FinalizeOrderRequest) (*orderpb.FinalizeOrderResponse, error) {
    ctx, end := s.instrument(ctx, "FinalizeOrder")
    defer end()

    log.Printf("[Order][FinalizeOrder] session_id=%s items_count=%d", req.SessionId, len(req.Items))
        
    var results []*orderpb.ItemResult
    okAll := true

    for _, item := range req.GetItems() {
        log.Printf("[Order][FinalizeOrder] checking product_id=%s quantity=%d", item.ProductId, item.Quantity)
        prod, err := s.inventory.GetProductInfo(ctx, &invpb.ProductId{ProductId: item.ProductId})
        res := &orderpb.ItemResult{ProductId: item.ProductId}

        if err != nil || prod.GetAvailableQuantity() < item.GetQuantity() {
			availableQty := int32(0)
			if prod != nil {
				availableQty = prod.GetAvailableQuantity()
			}
			log.Printf(
				"[Order][FinalizeOrder] insufficient stock for product_id=%s available=%d requested=%d",
				item.ProductId, availableQty, item.Quantity,
			)
            res.Reserved = false
            res.Message = "Insufficient stock"
            okAll = false
        } else {
			log.Printf(
				"[Order][FinalizeOrder] reserving stock for product_id=%s quantity=%d",
				item.ProductId, item.Quantity,
			)
            _, err := s.inventory.AdjustStock(ctx, &invpb.StockAdjustment{
                ProductId:      item.ProductId,
                QuantityChange: -item.Quantity,
            })
            if err != nil {
				log.Printf("[Order][FinalizeOrder] AdjustStock error: %v", err)
                res.Reserved = false
                res.Message = "Reservation failed"
                okAll = false
            } else {
				log.Printf("[Order][FinalizeOrder] Reserved product_id=%s", item.ProductId)
                res.Reserved = true
                res.Message = "Reserved"
            }
        }
        results = append(results, res)
    }

    log.Printf("[Order][FinalizeOrder] clearing session_id=%s", req.SessionId)
    s.mu.Lock()
    delete(s.sessions, req.SessionId)
    s.mu.Unlock()

    msg := "Order finalized"
    if !okAll {
		log.Printf("[Order][FinalizeOrder] partial failure")
        msg = "One or more items failed"
	} else {
		log.Printf("[Order][FinalizeOrder] finalize success")
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

	log.Printf("[Order][ConfirmOrderStock] session_id=%s", req.SessionId)
	for _, item := range req.GetItems() {
		log.Printf(
			"[Order][ConfirmOrderStock] adjusting stock for product_id=%s quantity=%d",
			item.ProductId, item.Quantity,
		)
        _, err := s.inventory.AdjustStock(ctx, &invpb.StockAdjustment{
            ProductId:      item.ProductId,
            QuantityChange: -item.Quantity,
        })
        if err != nil {
			log.Printf("[Order][ConfirmOrderStock] AdjustStock error: %v", err)
            return &invpb.OperationStatus{Success: false, Message: "Stock error"}, nil
        }
        log.Printf("[Order][ConfirmOrderStock] Stock adjusted for product_id=%s", item.ProductId)
    }

	log.Printf("[Order][ConfirmOrderStock] confirm success")
    return &invpb.OperationStatus{Success: true, Message: "Stock confirmed"}, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
    _, end := s.instrument(ctx, "CancelOrder")
    defer end()

    log.Printf("[Order][CancelOrder] session_id=%s", req.SessionId)
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.sessions[req.SessionId]; !ok {
        log.Printf("[Order][CancelOrder] session_id=%s not found", req.SessionId)
        return &orderpb.CancelOrderResponse{Released: false, Message: "Session not found"}, nil
    }
    delete(s.sessions, req.SessionId)
    log.Printf("[Order][CancelOrder] cancel success for session_id=%s", req.SessionId)
    return &orderpb.CancelOrderResponse{Released: true, Message: "Order cancelled and reservations released"}, nil
}
