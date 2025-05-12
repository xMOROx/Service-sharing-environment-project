package main

import (
	"Service-sharing-environment-project/proto/inventory"
	"Service-sharing-environment-project/proto/order"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	orderServiceListenPort            = ":50052"
	defaultInventoryServiceTargetPort = "50051"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Warning: Environment variable %s not set, using fallback: %s", key, fallback)
	return fallback
}

type orderServer struct {
	order.UnimplementedOrderServiceServer
	inventoryClient inventory.InventoryServiceClient
}

func (s *orderServer) PlaceOrder(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	log.Printf("OrderService: Received PlaceOrder request for ProductID: %s, Quantity: %d", req.ProductId, req.Quantity)

	stockReq := &inventory.HelloRequest{Name: req.ProductId}
	stockRes, err := s.inventoryClient.SayHello(ctx, stockReq)
	if err != nil {
		log.Printf("OrderService: could not call InventoryService during PlaceOrder: %v", err)
		return nil, fmt.Errorf("failed to check stock: %w", err)
	}

	log.Printf("OrderService: Inventory response for ProductID %s: %s", req.ProductId, stockRes.Message)

	return &order.OrderResponse{Message: "Order placed successfully, stock says: " + stockRes.Message}, nil
}

func main() {
	inventoryServiceTarget := getEnv("INVENTORY_SERVICE_ENDPOINT", "localhost:"+defaultInventoryServiceTargetPort)

	log.Printf("Order Service: Attempting to connect to Inventory Service at %s", inventoryServiceTarget)

	conn, err := grpc.Dial(inventoryServiceTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Order Service: failed to connect to Inventory Service: %v", err)
	}
	defer conn.Close()

	inventoryClient := inventory.NewInventoryServiceClient(conn)
	log.Println("Order Service: Successfully connected to Inventory Service.")

	lis, err := net.Listen("tcp", orderServiceListenPort)
	if err != nil {
		log.Fatalf("Order Service: failed to listen on port %s: %v", orderServiceListenPort, err)
	}

	grpcServer := grpc.NewServer()

	order.RegisterOrderServiceServer(grpcServer, &orderServer{inventoryClient: inventoryClient})

	log.Printf("Order Service: Starting gRPC server, listening on %s", orderServiceListenPort)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Order Service: failed to serve gRPC: %v", err)
		}
	}()

	log.Println("Order Service: Starting client loop to call Inventory Service...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		testProductId := "product-abc"
		log.Printf("Order Service: Calling InventoryService.SayHello for ProductID: %s", testProductId)

		helloReq := &inventory.HelloRequest{Name: testProductId}
		helloRes, err := inventoryClient.SayHello(context.Background(), helloReq)
		if err != nil {
			log.Printf("Order Service: ERROR could not call InventoryService.SayHello: %v", err)
			continue
		}
		log.Printf("Order Service: Response from InventoryService.SayHello for %s: %s", testProductId, helloRes.Message)
	}
}
