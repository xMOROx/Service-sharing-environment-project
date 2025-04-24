package main

import (
	"Service-sharing-environment-project/services/inventory-service/internal"
	pb "Service-sharing-environment-project/services/inventory-service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":50051"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, &internal.InventoryServer{})

	log.Printf("Inventory Service listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
