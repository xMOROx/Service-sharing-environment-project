package internal

import (
	invpb "Service-sharing-environment-project/proto/inventory"
	orderpb "Service-sharing-environment-project/proto/order"
	"context"
	"errors"
	"sync"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	inventory invpb.InventoryServiceClient
	sessions  map[string][]*invpb.OrderItemRequest
	mu        sync.Mutex
}

func NewOrderServer(client invpb.InventoryServiceClient) *OrderServer {
	return &OrderServer{
		inventory: client,
		sessions:  make(map[string][]*invpb.OrderItemRequest),
	}
}

func (s *OrderServer) CheckItemAvailability(ctx context.Context, req *invpb.ProductId) (*invpb.ProductInfo, error) {
	return s.inventory.GetProductInfo(ctx, req)
}

func (s *OrderServer) BuildOrder(stream orderpb.OrderService_BuildOrderServer) error {
	return errors.New("Not implemented")
}

func (s *OrderServer) FinalizeOrder(ctx context.Context, req *orderpb.FinalizeOrderRequest) (*orderpb.FinalizeOrderResponse, error) {
	results := []*orderpb.ItemResult{}
	success := true

	for _, item := range req.Items {
		product, err := s.inventory.GetProductInfo(ctx, &invpb.ProductId{ProductId: item.ProductId})
		res := &orderpb.ItemResult{ProductId: item.ProductId}

		if err != nil || product.AvailableQuantity < item.Quantity {
			res.Reserved = false
			res.Message = "Insufficient stock"
			success = false
		} else {
			adj := &invpb.StockAdjustment{
				ProductId:      item.ProductId,
				QuantityChange: -item.Quantity,
			}
			_, err := s.inventory.AdjustStock(ctx, adj)
			if err != nil {
				res.Reserved = false
				res.Message = "Reservation failed"
				success = false
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
	if !success {
		msg = "One or more items failed"
	}
	return &orderpb.FinalizeOrderResponse{
		Success:     success,
		Message:     msg,
		ItemResults: results,
	}, nil
}

func (s *OrderServer) ConfirmOrderStock(ctx context.Context, req *orderpb.FinalizeOrderRequest) (*invpb.OperationStatus, error) {
	for _, item := range req.Items {
		adj := &invpb.StockAdjustment{
			ProductId:      item.ProductId,
			QuantityChange: -item.Quantity,
		}
		_, err := s.inventory.AdjustStock(ctx, adj)
		if err != nil {
			return &invpb.OperationStatus{Success: false, Message: "Stock error"}, nil
		}
	}
	return &invpb.OperationStatus{Success: true, Message: "Stock confirmed"}, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[req.SessionId]; !ok {
		return &orderpb.CancelOrderResponse{Released: false, Message: "Session not found"}, nil
	}
	delete(s.sessions, req.SessionId)
	// Optionally, call Inventory to release soft reservations here
	return &orderpb.CancelOrderResponse{Released: true, Message: "Order cancelled and reservations released"}, nil
}
