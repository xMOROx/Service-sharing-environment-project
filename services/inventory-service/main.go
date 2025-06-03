package main

import (
	"context"
	"log"
	"net"

	invpb "Service-sharing-environment-project/proto/inventory"
	internal "Service-sharing-environment-project/services/inventory-service/internal"
	telemetry "Service-sharing-environment-project/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

const port = ":50051"

func main() {
    ctx := context.Background()

    // ── Tracing setup ─────────────────────────────────────────────────────────
    tp, err := telemetry.InitTracer(ctx, "inventory-service")
    if err != nil {
        log.Fatalf("Inventory: tracer init error: %v", err)
    }
    defer func() {
        if err := tp.Shutdown(ctx); err != nil {
            log.Printf("Inventory: error shutting down tracer: %v", err)
        }
    }()

    // ── Metrics setup (OTLP/gRPC) ──────────────────────────────────────────────
    mp, err := telemetry.InitMetrics(ctx, "inventory-service")
    if err != nil {
        log.Fatalf("Inventory: metrics init error: %v", err)
    }
    // UWAGA: nie wystawiamy Prometheus HTTP /metrics – metryki pushują się do OTLP Collector

    // ── Start gRPC server ──────────────────────────────────────────────────────
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Inventory Service: failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer(
        // ❶ Usuwamy grpc.WithTransportCredentials(insecure.NewCredentials())
        // ❷ Rejestrujemy StatsHandler, żeby OTel/Instruments automatycznie łapało spany i metryki
        grpc.StatsHandler(otelgrpc.NewServerHandler()),
    )

    invSrv := internal.NewInventoryServer(mp.Meter("inventory-service"))
    invpb.RegisterInventoryServiceServer(grpcServer, invSrv)

    log.Printf("Inventory Service: Starting gRPC server, listening on %s", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Inventory Service: failed to serve gRPC: %v", err)
    }
}
