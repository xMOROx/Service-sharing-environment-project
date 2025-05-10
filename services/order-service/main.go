package main

import (
	"Service-sharing-environment-project/proto/inventory"
	"Service-sharing-environment-project/proto/order"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

const (
	orderServicePort     = ":50052"
	inventoryServicePort = ":50051"
)

func main() {
	conn, err := grpc.Dial("localhost"+inventoryServicePort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Inventory Service: %v", err)
	}
	defer conn.Close()

	client := inventory.NewInventoryServiceClient(conn)

	orderRequest := &order.OrderRequest{
		ProductId: "1234",
		Quantity:  2,
	}

	helloReq := &inventory.HelloRequest{Name: orderRequest.ProductId}
	helloRes, err := client.SayHello(context.Background(), helloReq)
	if err != nil {
		log.Fatalf("could not call InventoryService: %v", err)
	}

	fmt.Printf("Order for %d of product %s placed. Inventory says: %s\n", orderRequest.Quantity, orderRequest.ProductId, helloRes.Message)
}
