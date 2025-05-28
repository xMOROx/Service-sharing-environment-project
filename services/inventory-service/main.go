package main

import (
	"context"
	"log"
	"net"
	"net/http"

	invpb "Service-sharing-environment-project/proto/inventory"
	"Service-sharing-environment-project/services/inventory-service/internal"
	"Service-sharing-environment-project/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

const port = ":50051"

func main() {
	    ctx := context.Background()
    // Tracing
    tp, err := telemetry.InitTracer(ctx, "inventory-service")
    if err != nil {
        log.Fatalf("Inventory: tracer init error: %v", err)
    }
    defer tp.Shutdown(ctx)

    // Metrics
    mp, promHandler, err := telemetry.InitMetrics(ctx, "inventory-service")
    if err != nil {
        log.Fatalf("Inventory: metrics init error: %v", err)
    }
    // Serve Prometheus metrics
    go func() {
        http.Handle("/metrics", promHandler)
        log.Println("Inventory: metrics endpoint -> :2222/metrics")
        log.Fatal(http.ListenAndServe(":2222", nil))
    }()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Inventory Service: failed to listen: %v", err)
	}

    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
        grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
    )
    invSrv := internal.NewInventoryServer(mp.Meter("inventory-service"))
    invpb.RegisterInventoryServiceServer(grpcServer, invSrv)

	log.Printf("Inventory Service: Starting gRPC server, listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Inventory Service: failed to serve gRPC: %v", err)
	}
}
