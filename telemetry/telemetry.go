package telemetry

import (
	"context"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	promexp "go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// InitTracer sets up an OTLP exporter with a TracerProvider for the given service.
func InitTracer(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
    endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
    if endpoint == "" {
        endpoint = "localhost:4317"
    }
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithInsecure(),
        otlptracegrpc.WithEndpoint(endpoint),
    )
    if err != nil {
        return nil, err
    }
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String(serviceName),
            attribute.String("environment", os.Getenv("ENVIRONMENT")),
        ),
    )
    if err != nil {
        return nil, err
    }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)
    return tp, nil
}

// InitMetrics sets up a Prometheus metric exporter and MeterProvider.
// Returns the MeterProvider and http.Handler for /metrics.
func InitMetrics(ctx context.Context, serviceName string) (*metric.MeterProvider, http.Handler, error) {
    // Create Prometheus exporter using OTEL metric exporter
    exporter, err := promexp.New(promexp.Config{})
    if err != nil {
        return nil, nil, err
    }
    // Setup resource for metrics
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String(serviceName),
            attribute.String("environment", os.Getenv("ENVIRONMENT")),
        ),
    )
    if err != nil {
        return nil, nil, err
    }
    // Create MeterProvider with the exporter as Reader
    mp := metric.NewMeterProvider(
        metric.WithReader(exporter),
        metric.WithResource(res),
    )
    otel.SetMeterProvider(mp)
    // Exporter implements http.Handler (ServeHTTP)
    handler := exporter
    return mp, handler, nil
}