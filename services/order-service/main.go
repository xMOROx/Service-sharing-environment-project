package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	invpb "Service-sharing-environment-project/proto/inventory"
	orderpb "Service-sharing-environment-project/proto/order"
	internal "Service-sharing-environment-project/services/order-service/internal"
	telemetry "Service-sharing-environment-project/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
    orderServiceListenPort            = ":50052"
    defaultInventoryServiceTargetPort = "50051"
)

func getEnv(key, fallback string) string {
    if v, ok := os.LookupEnv(key); ok {
        return v
    }
    log.Printf("Warning: %s not set, using fallback %q", key, fallback)
    return fallback
}

func Test(client invpb.InventoryServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 1) GetProductInfo
    pid := &invpb.ProductId{ProductId: "P001"}
    info, err := client.GetProductInfo(ctx, pid)
    if err != nil {
        log.Printf("GetProductInfo error: %v", err)
    } else {
        fmt.Printf("Product Info: %+v\n", info)
    }

    // 2) AdjustStock
    adj := &invpb.StockAdjustment{ProductId: "P001", QuantityChange: 5}
    resp2, err := client.AdjustStock(ctx, adj)
    if err != nil {
        log.Printf("AdjustStock error: %v", err)
    } else {
        fmt.Printf("AdjustStock response: %+v\n", resp2)
    }

    // 3) ListProducts
    stream, err := client.ListProducts(ctx, &invpb.ProductFilter{IncludeDiscontinued: true})
    if err != nil {
        log.Printf("ListProducts error: %v", err)
        return
    }
    fmt.Println("Products:")
    for {
        p, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Stream error: %v", err)
            break
        }
        fmt.Printf("- %s (%s): %d in stock\n", p.Name, p.ProductId, p.AvailableQuantity)
    }
}

func main() {
    ctx := context.Background()

    // ── Tracing setup ─────────────────────────────────────────────────────────
    tp, err := telemetry.InitTracer(ctx, "order-service")
    if err != nil {
        log.Fatalf("[Order] tracer init error: %v", err)
    }
    defer func() {
        if err := tp.Shutdown(ctx); err != nil {
            log.Printf("[Order] error shutting down tracer: %v", err)
        }
    }()

    // ── Metrics setup (OTLP/gRPC) ──────────────────────────────────────────────
    mp, err := telemetry.InitMetrics(ctx, "order-service")
    if err != nil {
        log.Fatalf("[Order] metrics init error: %v", err)
    }

    // ── Connect to Inventory Service ─────────────────────────────────────────
    invTarget := getEnv("INVENTORY_SERVICE_ENDPOINT", "localhost:"+defaultInventoryServiceTargetPort)
    log.Printf("[Order] connecting to Inventory at %s", invTarget)

    conn, err := grpc.DialContext(ctx, invTarget,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
    )
    if err != nil {
        log.Fatalf("[Order] failed to dial inventory: %v", err)
    }
    defer conn.Close()

    invClient := invpb.NewInventoryServiceClient(conn)
    log.Println("[Order] connected to Inventory.")

    // quick background smoke‐test
    go Test(invClient)

    // ── Start gRPC Server ────────────────────────────────────────────────────
    lis, err := net.Listen("tcp", orderServiceListenPort)
    if err != nil {
        log.Fatalf("[Order] listen error: %v", err)
    }

    grpcServer := grpc.NewServer(
        grpc.StatsHandler(otelgrpc.NewServerHandler()), // server‐side StatsHandler :contentReference[oaicite:3]{index=3}
    )

    // Przekazujemy meter do konstruktora serwera
    orderSrv := internal.NewOrderServer(invClient, mp.Meter("order-service"))
    orderpb.RegisterOrderServiceServer(grpcServer, orderSrv)

    log.Printf("[Order] gRPC listening on %s", orderServiceListenPort)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("[Order] Serve() error: %v", err)
    }
}
