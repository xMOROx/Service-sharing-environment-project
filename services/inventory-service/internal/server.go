package internal

import (
	pb "Service-sharing-environment-project/services/inventory-service/proto"
	"context"
	"fmt"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
}

func (s *InventoryServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello, %s from Inventory Service!", req.Name)}, nil
}
