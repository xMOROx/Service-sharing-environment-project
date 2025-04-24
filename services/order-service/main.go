package main

import (
	order "Service-sharing-environment-project/services/order-service/proto"
	"context"
	inventorypb "example.com/inventory-service/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

const (
	orderServicePort     = ":50052" // Port where order service will run (but no server, just client)
	inventoryServicePort = ":50051" // Port where inventory service runs
)

func main() {
	// Set up a connection to the InventoryService
	conn, err := grpc.Dial("localhost"+inventoryServicePort, grpc.WithInsecure()) // Connect to inventory service
	if err != nil {
		log.Fatalf("failed to connect to Inventory Service: %v", err)
	}
	defer conn.Close()

	// Create a new InventoryService client
	client := inventorypb.NewInventoryServiceClient(conn)

	// Simulate placing an order and calling the InventoryService
	orderRequest := &order.OrderRequest{
		ProductId: "1234", // Example product ID
		Quantity:  2,      // Example quantity
	}

	// Place order and call InventoryService SayHello method
	helloReq := &inventorypb.HelloRequest{Name: orderRequest.ProductId} // Use product ID as the name
	helloRes, err := client.SayHello(context.Background(), helloReq)
	if err != nil {
		log.Fatalf("could not call InventoryService: %v", err)
	}

	// Handle the response from the InventoryService
	fmt.Printf("Order for %d of product %s placed. Inventory says: %s\n", orderRequest.Quantity, orderRequest.ProductId, helloRes.Message)
}
