package main

import (
	invpb "Service-sharing-environment-project/proto/inventory"
	orderpb "Service-sharing-environment-project/proto/order"
	"Service-sharing-environment-project/services/order-service/internal"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
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

func main() {
	inventoryServiceTarget := getEnv("INVENTORY_SERVICE_ENDPOINT", "localhost:"+defaultInventoryServiceTargetPort)

	log.Printf("Order Service: Attempting to connect to Inventory Service at %s", inventoryServiceTarget)

	conn, err := grpc.Dial(inventoryServiceTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Order Service: failed to connect to Inventory Service: %v", err)
	}
	defer conn.Close()

	inventoryClient := invpb.NewInventoryServiceClient(conn)
	log.Println("Order Service: Successfully connected to Inventory Service.")

	lis, err := net.Listen("tcp", orderServiceListenPort)
	if err != nil {
		log.Fatalf("Order Service: failed to listen on port %s: %v", orderServiceListenPort, err)
	}

	grpcServer := grpc.NewServer()
	orderSrv := internal.NewOrderServer(inventoryClient)
	orderpb.RegisterOrderServiceServer(grpcServer, orderSrv)

	log.Printf("Order Service: Starting gRPC server, listening on %s", orderServiceListenPort)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Order Service: failed to serve gRPC: %v", err)
		}
	}()
	fmt.Println("Services running on port 50051")
	grpcServer.Serve(lis)
}
