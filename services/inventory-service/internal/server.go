package internal

import (
	pb "Service-sharing-environment-project/proto/inventory"
	"context"
	"errors"
	"io"
	"sync"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	mu       sync.Mutex
	products map[string]*pb.ProductInfo
}

func NewInventoryServer() *InventoryServer {
	return &InventoryServer{
		products: map[string]*pb.ProductInfo{
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
				AvailableQuantity: 0, // Out of stock
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
		},
	}
}

func (s *InventoryServer) GetProductInfo(ctx context.Context, req *pb.ProductId) (*pb.ProductInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[req.ProductId]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *InventoryServer) AddProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; exists {
		return &pb.OperationStatus{Success: false, Message: "Product already exists"}, nil
	}

	s.products[req.ProductId] = req
	return &pb.OperationStatus{Success: true, Message: "Product added"}, nil
}

func (s *InventoryServer) UpdateProduct(ctx context.Context, req *pb.ProductInfo) (*pb.OperationStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[req.ProductId]; !exists {
		return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
	}

	s.products[req.ProductId] = req
	return &pb.OperationStatus{Success: true, Message: "Product updated"}, nil
}

func (s *InventoryServer) RemoveProduct(ctx context.Context, req *pb.ProductId) (*pb.OperationStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if product, exists := s.products[req.ProductId]; exists {
		product.Discontinued = true
		return &pb.OperationStatus{Success: true, Message: "Product discontinued"}, nil
	}
	return &pb.OperationStatus{Success: false, Message: "Product not found"}, nil
}

func (s *InventoryServer) AdjustStock(ctx context.Context, req *pb.StockAdjustment) (*pb.OperationStatus, error) {
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
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ProductId]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *InventoryServer) ListProducts(req *pb.ProductFilter, stream pb.InventoryService_ListProductsServer) error {
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
	return errors.New("SubscribeLowStockAlerts not implemented")
}

func (s *InventoryServer) InteractiveOrderStock(stream pb.InventoryService_InteractiveOrderStockServer) error {
	return errors.New("InteractiveOrderStock not implemented")
}
