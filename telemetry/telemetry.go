package telemetry

import (
	"context"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	apiMetric "go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	goprom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

// InitMetrics sets up a Prometheus exporter and MeterProvider.
// Returns the MeterProvider and an http.Handler to serve /metrics.
func InitMetrics(ctx context.Context, serviceName string) (apiMetric.MeterProvider, http.Handler, error) {
	// Create a Prometheus exporter that registers with the default Prometheus
	// registry so that promhttp.Handler() will pick it up.
	exp, err := otelprom.New(
		otelprom.WithRegisterer(goprom.DefaultRegisterer),
	)
	if err != nil {
		return nil, nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", os.Getenv("ENVIRONMENT")),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	// Use sdkMetric.NewMeterProvider to create the MeterProvider, passing
	// the Prometheus exporter as a Reader, then set it as the global.
	mp := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(exp),
		sdkMetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	// promhttp.Handler() will serve all metrics registered in goprom.DefaultRegisterer.
	handler := promhttp.Handler()
	return mp, handler, nil
}
