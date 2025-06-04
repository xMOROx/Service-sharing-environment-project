package internal

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	// Po wygenerowaniu kodu *.pb.go import powinien wskazywać dokładnie
	// tam, gdzie powstały pliki Go z inventory.proto:
	pb "Service-sharing-environment-project/proto/inventory"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// InventoryServer to domyślna implementacja pb.InventoryServiceServer
type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer

	mu       sync.Mutex
	products map[string]*pb.ProductInfo

	requestCounter metric.Int64Counter
	latencyHist    metric.Float64Histogram
}

// NewInventoryServer tworzy nowy serwer, inicjalizuje mapę produktów i instrumenty metryk
func NewInventoryServer(m metric.Meter) *InventoryServer {
	ctr, err := m.Int64Counter(
		"inventory_requests_total",
		metric.WithDescription("Total inventory service requests"),
	)
	if err != nil {
		panic(err)
	}

	hist, err := m.Float64Histogram(
		"inventory_request_latency_ms",
		metric.WithDescription("Latency of inventory service requests in ms"),
	)
	if err != nil {
		panic(err)
	}

	// Przykładowe wypełnienie danymi
	initialProducts := map[string]*pb.ProductInfo{
		"P001": {
			ProductId:         "P001",
			Name:              "Wireless Mouse",
			Description:       "Ergonomic wireless mouse with USB receiver",
			Category:          "Electronics",
			Discontinued:      false,
			AvailableQuantity: 120,
			IsAvailable:       true,
		},
		"P002": {
			ProductId:         "P002",
			Name:              "Mechanical Keyboard",
			Description:       "RGB backlit mechanical keyboard with blue switches",
			Category:          "Electronics",
			Discontinued:      false,
			AvailableQuantity: 75,
			IsAvailable:       true,
		},
		"P003": {
			ProductId:         "P003",
			Name:              "Water Bottle",
			Description:       "Stainless steel water bottle, 750ml",
			Category:          "Home & Kitchen",
			Discontinued:      false,
			AvailableQuantity: 200,
			IsAvailable:       true,
		},
		"P004": {
			ProductId:         "P004",
			Name:              "Notebook",
			Description:       "A5 size ruled notebook, 200 pages",
			Category:          "Office Supplies",
			Discontinued:      false,
			AvailableQuantity: 0,
			IsAvailable:       false,
		},
		"P005": {
			ProductId:         "P005",
			Name:              "LED Desk Lamp",
			Description:       "Adjustable LED desk lamp with USB charging port",
			Category:          "Home & Kitchen",
			Discontinued:      false,
			AvailableQuantity: 45,
			IsAvailable:       true,
		},
	}

	return &InventoryServer{
		products:       initialProducts,
		requestCounter: ctr,
		latencyHist:    hist,
	}
}

// GetProductInfo zwraca szczegóły produktu dla podanego ProductId
func (s *InventoryServer) GetProductInfo(ctx context.Context, req *pb.ProductId) (*pb.ProductInfo, error) {
	start := time.Now()
	defer func() {
		// Rejestrujemy liczbę wywołań i czas trwania
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "GetProductInfo")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "GetProductInfo")),
		)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// AddProduct dodaje nowy produkt do mapy
func (s *InventoryServer) AddProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "AddProduct")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "AddProduct")),
		)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; exists {
		return &pb.OperationStatus{Success: false, Message: "Product already exists"}, nil
	}
	s.products[req.ProductId] = req
	return &pb.OperationStatus{Success: true, Message: "Product added"}, nil
}

// UpdateProduct aktualizuje istniejący produkt
func (s *InventoryServer) UpdateProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "UpdateProduct")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "UpdateProduct")),
		)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; !exists {
		return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
	}
	s.products[req.ProductId] = req
	return &pb.OperationStatus{Success: true, Message: "Product updated"}, nil
}

// RemoveProduct oznacza produkt jako wycofany (discontinued)
func (s *InventoryServer) RemoveProduct(ctx context.Context, req *pb.ProductId) (*pb.OperationStatus, error) {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "RemoveProduct")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "RemoveProduct")),
		)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	if product, exists := s.products[req.ProductId]; exists {
		product.Discontinued = true
		return &pb.OperationStatus{Success: true, Message: "Product discontinued"}, nil
	}
	return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
}

// AdjustStock modyfikuje AvailableQuantity o QuantityChange
func (s *InventoryServer) AdjustStock(ctx context.Context, req *pb.StockAdjustment) (*pb.OperationStatus, error) {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "AdjustStock")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "AdjustStock")),
		)
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

// BulkStockUpdate to RPC typu client‐streaming
func (s *InventoryServer) BulkStockUpdate(stream pb.InventoryService_BulkStockUpdateServer) error {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "BulkStockUpdate")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "BulkStockUpdate")),
		)
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.OperationStatus{Success: true, Message: "Bulk update complete"})
		}
		if err != nil {
			return err
		}
		// Używamy AdjustStock (bez metryk, bo już tu je zarejestrowaliśmy)
		_, err = s.AdjustStock(context.Background(), req)
		if err != nil {
			return err
		}
	}
}

// GetStockLevel zwraca stan magazynu (tożsamy z GetProductInfo)
func (s *InventoryServer) GetStockLevel(ctx context.Context, req *pb.ProductId) (*pb.ProductInfo, error) {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "GetStockLevel")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "GetStockLevel")),
		)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// ListProducts strumieniowo zwraca wszystkie produkty (opcjonalne filtrowanie)
func (s *InventoryServer) ListProducts(req *pb.ProductFilter, stream pb.InventoryService_ListProductsServer) error {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "ListProducts")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "ListProducts")),
		)
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

// SubscribeLowStockAlerts wysyła alerty co pewien czas, gdy AvailableQuantity <= threshold
func (s *InventoryServer) SubscribeLowStockAlerts(req *pb.LowStockSubscription, stream pb.InventoryService_SubscribeLowStockAlertsServer) error {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "SubscribeLowStockAlerts")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "SubscribeLowStockAlerts")),
		)
	}()

	// Poprawne pole w LowStockSubscription to Threshold, nie LowStockThreshold
	threshold := req.Threshold
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
					// LowStockAlert ma pole CurrentQuantity, nie AvailableQuantity
					if err := stream.Send(&pb.LowStockAlert{
						ProductId:       p.ProductId,
						CurrentQuantity: p.AvailableQuantity,
						Message:         "Low stock", // Opcjonalne: możesz dodać własny komunikat
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

// InteractiveOrderStock to RPC bidirectional streaming (recv/send)
func (s *InventoryServer) InteractiveOrderStock(stream pb.InventoryService_InteractiveOrderStockServer) error {
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "InteractiveOrderStock")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "InteractiveOrderStock")),
		)
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

		// Poprawa pola RequestedQuantity zamiast Quantity
		if !ok || p.AvailableQuantity < req.RequestedQuantity {
			resp = pb.OrderItemResponse{
				ProductId:         req.ProductId,
				Available:         false,
				AvailableQuantity: p.AvailableQuantity, // zamiast RemainingQuantity
				Message:           "Insufficient stock",
			}
		} else {
			p.AvailableQuantity -= req.RequestedQuantity
			p.IsAvailable = p.AvailableQuantity > 0
			resp = pb.OrderItemResponse{
				ProductId:         req.ProductId,
				Available:         true,
				AvailableQuantity: p.AvailableQuantity, // pole AvailableQuantity
				Message:           "Reserved",
			}
		}
		s.mu.Unlock()

		if err := stream.Send(&resp); err != nil {
			return err
		}
	}
}
