package internal

import (
	pb "Service-sharing-environment-project/services/order-service/proto"
	"context"
	inventorypb "example.com/inventory-service/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
}

func (s *OrderServer) PlaceOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	// Setup the gRPC connection to the InventoryService
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Make sure the port is correct for your Inventory service
	if err != nil {
		log.Fatalf("Could not connect to Inventory Service: %v", err)
		return nil, err
	}
	defer conn.Close()

	// Create an InventoryService client
	inventoryClient := inventorypb.NewInventoryServiceClient(conn)

	helloReq := &inventorypb.HelloRequest{Name: req.ProductId} // Example: Use product ID as name for simplicity
	helloRes, err := inventoryClient.SayHello(ctx, helloReq)
	if err != nil {
		log.Printf("Error calling Inventory Service: %v", err)
		return nil, err
	}

	return &pb.OrderResponse{
		Confirmed: true,
		Message:   fmt.Sprintf("Order for %d of product %s has been placed. Inventory says: %s", req.Quantity, req.ProductId, helloRes.Message),
	}, nil
}
