package main

import (
	invpb "Service-sharing-environment-project/proto/inventory"
	orderpb "Service-sharing-environment-project/proto/order"
	"Service-sharing-environment-project/services/order-service/internal"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"os"
	"time"
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

func Test(client invpb.InventoryServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. GetProductInfo
	productID := &invpb.ProductId{ProductId: "P001"}
	info, err := client.GetProductInfo(ctx, productID)
	if err != nil {
		log.Printf("GetProductInfo error: %v", err)
	} else {
		fmt.Printf("Product Info: %+v\n", info)
	}

	// 2. AdjustStock
	adjustment := &invpb.StockAdjustment{
		ProductId:      "P001",
		QuantityChange: 5,
		Reason:         "Manual test",
	}
	adjustResp, err := client.AdjustStock(ctx, adjustment)
	if err != nil {
		log.Printf("AdjustStock error: %v", err)
	} else {
		fmt.Printf("AdjustStock response: %+v\n", adjustResp)
	}

	// 3. ListProducts
	filter := &invpb.ProductFilter{IncludeDiscontinued: true}
	stream, err := client.ListProducts(ctx, filter)
	if err != nil {
		log.Printf("ListProducts error: %v", err)
		return
	}
	fmt.Println("Products list:")
	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Stream error: %v", err)
			break
		}
		fmt.Printf("- %s (%s): %d in stock\n", product.Name, product.ProductId, product.AvailableQuantity)
	}
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
	go Test(inventoryClient)
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
