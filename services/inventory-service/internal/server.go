package internal

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	pb "Service-sharing-environment-project/proto/inventory"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type InventoryServer struct {
    pb.UnimplementedInventoryServiceServer
    mu             sync.Mutex
    products       map[string]*pb.ProductInfo
    requestCounter instrument.Int64Counter
    latencyHist    instrument.Float64Histogram
}

// NewInventoryServer inicjalizuje mapę produktów i metryki
func NewInventoryServer(m metric.Meter) *InventoryServer {
    counter := metric.Must(m).NewInt64Counter(
        "inventory_requests_total",
        instrument.WithDescription("Total inventory service requests"),
    )
    hist := metric.Must(m).NewFloat64Histogram(
        "inventory_request_latency_ms",
        instrument.WithDescription("Latency of inventory requests in ms"),
    )

    return &InventoryServer{
        products: map[string]*pb.ProductInfo{
            "P001": {ProductId: "P001", Name: "Wireless Mouse", Description: "Ergonomic wireless mouse with USB receiver", Category: "Electronics", Discontinued: false, AvailableQuantity: 120, IsAvailable: true},
            "P002": {ProductId: "P002", Name: "Mechanical Keyboard", Description: "RGB backlit mechanical keyboard with blue switches", Category: "Electronics", Discontinued: false, AvailableQuantity: 75, IsAvailable: true},
            "P003": {ProductId: "P003", Name: "Water Bottle", Description: "Stainless steel water bottle, 750ml", Category: "Home & Kitchen", Discontinued: false, AvailableQuantity: 200, IsAvailable: true},
            "P004": {ProductId: "P004", Name: "Notebook", Description: "A5 size ruled notebook, 200 pages", Category: "Office Supplies", Discontinued: false, AvailableQuantity: 0, IsAvailable: false},
            "P005": {ProductId: "P005", Name: "LED Desk Lamp", Description: "Adjustable LED desk lamp with USB charging port", Category: "Home & Kitchen", Discontinued: false, AvailableQuantity: 45, IsAvailable: true},
        },
        requestCounter: counter,
        latencyHist:    hist,
    }
}

func (s *InventoryServer) GetProductInfo(ctx context.Context, req *pb.ProductId) (*pb.ProductInfo, error) {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(ctx, 1, attribute.String("method", "GetProductInfo"))
        s.latencyHist.Record(ctx, float64(time.Since(start).Milliseconds()), attribute.String("method", "GetProductInfo"))
    }()
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[req.ProductId]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *InventoryServer) AddProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(ctx, 1, attribute.String("method", "AddProduct"))
        s.latencyHist.Record(ctx, float64(time.Since(start).Milliseconds()), attribute.String("method", "AddProduct"))
    }()
	
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; exists {
		return &pb.OperationStatus{Success: false, Message: "Product already exists"}, nil
	}

	s.products[req.ProductId] = req
	return &pb.OperationStatus{Success: true, Message: "Product added"}, nil
}

func (s *InventoryServer) UpdateProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(ctx, 1, attribute.String("method", "UpdateProduct"))
        s.latencyHist.Record(ctx, float64(time.Since(start).Milliseconds()), attribute.String("method", "UpdateProduct"))
    }()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; !exists {
		return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
	}

	s.products[req.ProductId] = req
	return &pb.OperationStatus{Success: true, Message: "Product updated"}, nil
}

func (s *InventoryServer) RemoveProduct(ctx context.Context, req *pb.ProductId) (*pb.OperationStatus, error) {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(ctx, 1, attribute.String("method", "RemoveProduct"))
        s.latencyHist.Record(ctx, float64(time.Since(start).Milliseconds()), attribute.String("method", "RemoveProduct"))
    }()
	
	s.mu.Lock()
	defer s.mu.Unlock()

	if product, exists := s.products[req.ProductId]; exists {
		product.Discontinued = true
		return &pb.OperationStatus{Success: true, Message: "Product discontinued"}, nil
	}
	return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
}

func (s *InventoryServer) AdjustStock(ctx context.Context, req *pb.StockAdjustment) (*pb.OperationStatus, error) {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(ctx, 1, attribute.String("method", "AdjustStock"))
        s.latencyHist.Record(ctx, float64(time.Since(start).Milliseconds()), attribute.String("method", "AdjustStock"))
    }()
	
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
	}

	product.AvailableQuantity += req.QuantityChange
	product.IsAvailable = product.AvailableQuantity > 0

	return &pb.OperationStatus{Success: true, Message: "Stock adjusted"}, nil
}

func (s *InventoryServer) BulkStockUpdate(stream pb.InventoryService_BulkStockUpdateServer) error {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(stream.Context(), 1, attribute.String("method", "GetProBulkStockUpdateductInfo"))
        s.latencyHist.Record(stream.Context(), float64(time.Since(start).Milliseconds()), attribute.String("method", "BulkStockUpdate"))
    }()
	
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.OperationStatus{Success: true, Message: "Bulk update complete"})
		}
		if err != nil {
			return err
		}
		_, err = s.AdjustStock(context.Background(), req)
		if err != nil {
			return err
		}
	}
}

func (s *InventoryServer) GetStockLevel(ctx context.Context, req *pb.ProductId) (*pb.ProductInfo, error) {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(ctx, 1, attribute.String("method", "GetStockLevel"))
        s.latencyHist.Record(ctx, float64(time.Since(start).Milliseconds()), attribute.String("method", "GetStockLevel"))
    }()
	
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *InventoryServer) ListProducts(req *pb.ProductFilter, stream pb.InventoryService_ListProductsServer) error {
	start := time.Now()
    defer func() {
        s.requestCounter.Add(stream.Context(), 1, attribute.String("method", "ListProducts"))
        s.latencyHist.Record(stream.Context(), float64(time.Since(start).Milliseconds()), attribute.String("method", "ListProducts"))
    }()
	
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.products {
		if !req.IncludeDiscontinued && p.Discontinued {
			continue
		}
		if req.Category != "" && p.Category != req.Category {
			continue
		}
		if err := stream.Send(p); err != nil {
			return err
		}
	}
	return nil
}

func (s *InventoryServer) SubscribeLowStockAlerts(req *pb.LowStockSubscription, stream pb.InventoryService_SubscribeLowStockAlertsServer) error {
    start := time.Now()
    defer func() {
        s.requestCounter.Add(stream.Context(), 1, attribute.String("method", "SubscribeLowStockAlerts"))
        s.latencyHist.Record(stream.Context(), float64(time.Since(start).Milliseconds()), attribute.String("method", "SubscribeLowStockAlerts"))
    }()

    threshold := req.LowStockThreshold
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-stream.Context().Done():
            return nil
        case <-ticker.C:
            s.mu.Lock()
            for _, p := range s.products {
                if p.AvailableQuantity <= threshold {
                    if err := stream.Send(&pb.LowStockAlert{
                        ProductId:         p.ProductId,
                        AvailableQuantity: p.AvailableQuantity,
                    }); err != nil {
                        s.mu.Unlock()
                        return err
                    }
                }
            }
            s.mu.Unlock()
        }
    }
}

func (s *InventoryServer) InteractiveOrderStock(stream pb.InventoryService_InteractiveOrderStockServer) error {
    start := time.Now()
    defer func() {
        s.requestCounter.Add(stream.Context(), 1, attribute.String("method", "InteractiveOrderStock"))
        s.latencyHist.Record(stream.Context(), float64(time.Since(start).Milliseconds()), attribute.String("method", "InteractiveOrderStock"))
    }()

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        s.mu.Lock()
        p, ok := s.products[req.ProductId]
        var resp pb.OrderItemResponse
        if !ok || p.AvailableQuantity < req.Quantity {
            resp = pb.OrderItemResponse{
                ProductId:         req.ProductId,
                Available:         false,
                RemainingQuantity: p.AvailableQuantity,
            }
        } else {
            p.AvailableQuantity -= req.Quantity
            p.IsAvailable = p.AvailableQuantity > 0
            resp = pb.OrderItemResponse{
                ProductId:         req.ProductId,
                Available:         true,
                RemainingQuantity: p.AvailableQuantity,
            }
        }
        s.mu.Unlock()

        if err := stream.Send(&resp); err != nil {
            return err
        }
    }
}