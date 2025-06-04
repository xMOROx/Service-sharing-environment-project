package internal

import (
	"context"
	"errors"
	"io"
	"log"
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
	log.Printf("[Inventory][GetProductInfo] called with product_id=%s", req.ProductId)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "GetProductInfo")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "GetProductInfo")),
		)
		log.Printf("[Inventory][GetProductInfo] latency=%.2fms", elapsedMs)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		log.Printf("[Inventory][GetProductInfo] product not found: %s", req.ProductId)
		return nil, errors.New("product not found")
	}
	log.Printf("[Inventory][GetProductInfo] found product: %s, quantity=%d", product.ProductId, product.AvailableQuantity)
	return product, nil
}

// AddProduct dodaje nowy produkt do mapy
func (s *InventoryServer) AddProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	log.Printf("[Inventory][AddProduct] called with product_id=%s name=%s", req.ProductId, req.Name)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "AddProduct")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "AddProduct")),
		)
		log.Printf("[Inventory][AddProduct] latency=%.2fms", elapsedMs)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; exists {
		log.Printf("[Inventory][AddProduct] product already exists: %s", req.ProductId)
		return &pb.OperationStatus{Success: false, Message: "Product already exists"}, nil
	}
	s.products[req.ProductId] = req
	log.Printf("[Inventory][AddProduct] product added: %s", req.ProductId)
	return &pb.OperationStatus{Success: true, Message: "Product added"}, nil
}

// UpdateProduct aktualizuje istniejący produkt
func (s *InventoryServer) UpdateProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	log.Printf("[Inventory][UpdateProduct] called with product_id=%s", req.ProductId)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "UpdateProduct")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "UpdateProduct")),
		)
		log.Printf("[Inventory][UpdateProduct] latency=%.2fms", elapsedMs)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; !exists {
		log.Printf("[Inventory][UpdateProduct] product not found: %s", req.ProductId)
		return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
	}
	s.products[req.ProductId] = req
	log.Printf("[Inventory][UpdateProduct] product updated: %s", req.ProductId)
	return &pb.OperationStatus{Success: true, Message: "Product updated"}, nil
}

// RemoveProduct oznacza produkt jako wycofany (discontinued)
func (s *InventoryServer) RemoveProduct(ctx context.Context, req *pb.ProductId) (*pb.OperationStatus, error) {
	log.Printf("[Inventory][RemoveProduct] called with product_id=%s", req.ProductId)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "RemoveProduct")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "RemoveProduct")),
		)
		log.Printf("[Inventory][RemoveProduct] latency=%.2fms", elapsedMs)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	if product, exists := s.products[req.ProductId]; exists {
		product.Discontinued = true
		log.Printf("[Inventory][RemoveProduct] marked discontinued: %s", req.ProductId)
		return &pb.OperationStatus{Success: true, Message: "Product discontinued"}, nil
	}
	log.Printf("[Inventory][RemoveProduct] product not found: %s", req.ProductId)
	return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
}

// AdjustStock modyfikuje AvailableQuantity o QuantityChange
func (s *InventoryServer) AdjustStock(ctx context.Context, req *pb.StockAdjustment) (*pb.OperationStatus, error) {
	log.Printf("[Inventory][AdjustStock] called with product_id=%s quantity_change=%d", req.ProductId, req.QuantityChange)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "AdjustStock")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "AdjustStock")),
		)
		log.Printf("[Inventory][AdjustStock] latency=%.2fms", elapsedMs)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		log.Printf("[Inventory][AdjustStock] product not found: %s", req.ProductId)
		return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
	}

	product.AvailableQuantity += req.QuantityChange
	product.IsAvailable = product.AvailableQuantity > 0
	log.Printf(
		"[Inventory][AdjustStock] new quantity for %s = %d",
		req.ProductId, product.AvailableQuantity,
	)
	return &pb.OperationStatus{Success: true, Message: "Stock adjusted"}, nil
}

// BulkStockUpdate to RPC typu client‐streaming
func (s *InventoryServer) BulkStockUpdate(stream pb.InventoryService_BulkStockUpdateServer) error {
	log.Printf("[Inventory][BulkStockUpdate] stream started")
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "BulkStockUpdate")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "BulkStockUpdate")),
		)
		log.Printf("[Inventory][BulkStockUpdate] latency=%.2fms", elapsedMs)
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[Inventory][BulkStockUpdate] stream EOF")
			return stream.SendAndClose(&pb.OperationStatus{Success: true, Message: "Bulk update complete"})
		}
		if err != nil {
			log.Printf("[Inventory][BulkStockUpdate] Recv error: %v", err)
			return err
		}
		log.Printf(
			"[Inventory][BulkStockUpdate] adjusting product_id=%s quantity_change=%d",
			req.ProductId, req.QuantityChange,
		)
		_, err = s.AdjustStock(context.Background(), req)
		if err != nil {
			log.Printf("[Inventory][BulkStockUpdate] AdjustStock error: %v", err)
			return err
		}
	}
}

// GetStockLevel zwraca stan magazynu (tożsamy z GetProductInfo)
func (s *InventoryServer) GetStockLevel(ctx context.Context, req *pb.ProductId) (*pb.ProductInfo, error) {
	log.Printf("[Inventory][GetStockLevel] called with product_id=%s", req.ProductId)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(ctx, 1,
			metric.WithAttributes(attribute.String("method", "GetStockLevel")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(ctx, elapsedMs,
			metric.WithAttributes(attribute.String("method", "GetStockLevel")),
		)
		log.Printf("[Inventory][GetStockLevel] latency=%.2fms", elapsedMs)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		log.Printf("[Inventory][GetStockLevel] product not found: %s", req.ProductId)
		return nil, errors.New("product not found")
	}
	log.Printf("[Inventory][GetStockLevel] found product: %s, quantity=%d", product.ProductId, product.AvailableQuantity)
	return product, nil
}

// ListProducts strumieniowo zwraca wszystkie produkty (opcjonalne filtrowanie)
func (s *InventoryServer) ListProducts(req *pb.ProductFilter, stream pb.InventoryService_ListProductsServer) error {
	log.Printf(
		"[Inventory][ListProducts] called with category=%q include_discontinued=%t",
		req.Category, req.IncludeDiscontinued,
	)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "ListProducts")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "ListProducts")),
		)
		log.Printf("[Inventory][ListProducts] latency=%.2fms", elapsedMs)
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
		log.Printf(
			"[Inventory][ListProducts] sending product_id=%s quantity=%d",
			p.ProductId, p.AvailableQuantity,
		)
		if err := stream.Send(p); err != nil {
			log.Printf("[Inventory][ListProducts] Send error: %v", err)
			return err
		}
	}
	log.Printf("[Inventory][ListProducts] stream completed")
	return nil
}

// SubscribeLowStockAlerts wysyła alerty co pewien czas, gdy AvailableQuantity <= threshold
func (s *InventoryServer) SubscribeLowStockAlerts(req *pb.LowStockSubscription, stream pb.InventoryService_SubscribeLowStockAlertsServer) error {
	log.Printf(
		"[Inventory][SubscribeLowStockAlerts] called with threshold=%d product_ids=%v",
		req.Threshold, req.ProductIds,
	)
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "SubscribeLowStockAlerts")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "SubscribeLowStockAlerts")),
		)
		log.Printf("[Inventory][SubscribeLowStockAlerts] latency=%.2fms", elapsedMs)
	}()

	threshold := req.Threshold
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("[Inventory][SubscribeLowStockAlerts] client canceled")
			return nil
		case <-ticker.C:
			s.mu.Lock()
			for _, p := range s.products {
				if p.AvailableQuantity <= threshold {
					log.Printf(
						"[Inventory][SubscribeLowStockAlerts] alert for product_id=%s quantity=%d",
						p.ProductId, p.AvailableQuantity,
					)
					if err := stream.Send(&pb.LowStockAlert{
						ProductId:       p.ProductId,
						CurrentQuantity: p.AvailableQuantity,
						Message:         "Low stock",
					}); err != nil {
						s.mu.Unlock()
						log.Printf("[Inventory][SubscribeLowStockAlerts] Send error: %v", err)
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
	log.Printf("[Inventory][InteractiveOrderStock] stream started")
	start := time.Now()
	defer func() {
		s.requestCounter.Add(stream.Context(), 1,
			metric.WithAttributes(attribute.String("method", "InteractiveOrderStock")),
		)
		elapsedMs := float64(time.Since(start).Milliseconds())
		s.latencyHist.Record(stream.Context(), elapsedMs,
			metric.WithAttributes(attribute.String("method", "InteractiveOrderStock")),
		)
		log.Printf("[Inventory][InteractiveOrderStock] latency=%.2fms", elapsedMs)
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[Inventory][InteractiveOrderStock] stream EOF")
			return nil
		}
		if err != nil {
			log.Printf("[Inventory][InteractiveOrderStock] Recv error: %v", err)
			return err
		}
		log.Printf(
			"[Inventory][InteractiveOrderStock] received request product_id=%s requested_quantity=%d",
			req.ProductId, req.RequestedQuantity,
		)

		s.mu.Lock()
		p, ok := s.products[req.ProductId]
		var resp pb.OrderItemResponse

		if !ok || p.AvailableQuantity < req.RequestedQuantity {
			log.Printf(
				"[Inventory][InteractiveOrderStock] insufficient stock for product_id=%s current=%d requested=%d",
				req.ProductId, p.AvailableQuantity, req.RequestedQuantity,
			)
			resp = pb.OrderItemResponse{
				ProductId:         req.ProductId,
				Available:         false,
				AvailableQuantity: p.AvailableQuantity,
				Message:           "Insufficient stock",
			}
		} else {
			p.AvailableQuantity -= req.RequestedQuantity
			p.IsAvailable = p.AvailableQuantity > 0
			log.Printf(
				"[Inventory][InteractiveOrderStock] reserved product_id=%s new_quantity=%d",
				req.ProductId, p.AvailableQuantity,
			)
			resp = pb.OrderItemResponse{
				ProductId:         req.ProductId,
				Available:         true,
				AvailableQuantity: p.AvailableQuantity,
				Message:           "Reserved",
			}
		}
		s.mu.Unlock()

		if err := stream.Send(&resp); err != nil {
			log.Printf("[Inventory][InteractiveOrderStock] Send error: %v", err)
			return err
		}
		log.Printf(
			"[Inventory][InteractiveOrderStock] sent response for product_id=%s available=%v remaining=%d",
			resp.ProductId, resp.Available, resp.AvailableQuantity,
		)
	}
}
