package main

import (
	"cell-router/router"
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initOpenTelemetry() func() {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("cell-router"),
		)),
	)
	otel.SetMeterProvider(meterProvider)

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("cell-router"),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		_ = meterProvider.Shutdown(context.Background())
		_ = tracerProvider.Shutdown(context.Background())
	}
}

func main() {
	cleanup := initOpenTelemetry()
	defer cleanup()

	// Create a new mux for metrics
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	// Start metrics server in a goroutine
	go func() {
		log.Println("Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", metricsMux); err != nil {
			log.Fatal(err)
		}
	}()

	// Start application server
	r := router.NewRouter()
	log.Println("Starting application server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
