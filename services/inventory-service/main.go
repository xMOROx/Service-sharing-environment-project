package main

import (
	"Service-sharing-environment-project/proto/inventory"
	"Service-sharing-environment-project/services/inventory-service/internal"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":50051"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Inventory Service: failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	inventory.RegisterInventoryServiceServer(grpcServer, internal.NewInventoryServer())

	log.Printf("Inventory Service: Starting gRPC server, listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Inventory Service: failed to serve gRPC: %v", err)
	}
}
