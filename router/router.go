package router

import (
	"cell-router/config"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Router struct {
	config       *config.Config
	meter        metric.Meter
	tracer       trace.Tracer
	requestCount metric.Int64Counter
	latency      metric.Float64Histogram
	successCount metric.Int64Counter
	errorCount   metric.Int64Counter
}

func NewRouter() *Router {
	cfg := config.LoadConfig()
	log.Println("Configuration loaded successfully")

	meter := otel.GetMeterProvider().Meter("cell-router")
	tracer := otel.GetTracerProvider().Tracer("cell-router")

	requestCount, _ := meter.Int64Counter("requests_total",
		metric.WithDescription("Total number of requests received"))

	latency, _ := meter.Float64Histogram("request_latency",
		metric.WithDescription("Request latency in milliseconds"))

	successCount, _ := meter.Int64Counter("requests_success_total",
		metric.WithDescription("Total number of successful requests"))

	errorCount, _ := meter.Int64Counter("requests_error_total",
		metric.WithDescription("Total number of failed requests"))

	return &Router{
		config:       cfg,
		meter:        meter,
		tracer:       tracer,
		requestCount: requestCount,
		latency:      latency,
		successCount: successCount,
		errorCount:   errorCount,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Add health check endpoint
	if req.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	ctx, span := r.tracer.Start(context.Background(), "ServeHTTP")
	defer span.End()

	startTime := time.Now()
	r.requestCount.Add(ctx, 1)

	clientIDStr := req.URL.Query().Get("client_id")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		log.Printf("Invalid client_id: %s", clientIDStr)
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		r.errorCount.Add(ctx, 1)
		return
	}

	log.Printf("Received request for client_id: %d", clientID)

	for _, cell := range r.config.Cells {
		if clientID >= cell.RangeFrom && clientID <= cell.RangeTo {
			log.Printf("Routing to cell: %s, endpoint: %s", cell.Name, cell.Endpoint)

			proxyURL, err := url.Parse(cell.Endpoint)
			if err != nil {
				log.Printf("Error parsing endpoint URL: %v", err)
				http.Error(w, "Invalid endpoint URL", http.StatusInternalServerError)
				r.errorCount.Add(ctx, 1)
				return
			}

			proxyReq, err := http.NewRequest(req.Method, proxyURL.String(), req.Body)
			if err != nil {
				log.Printf("Error creating proxy request: %v", err)
				http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
				r.errorCount.Add(ctx, 1)
				return
			}
			proxyReq.Header = req.Header

			client := &http.Client{}
			resp, err := client.Do(proxyReq)
			if err != nil {
				log.Printf("Error making request to cell: %v", err)
				http.Error(w, "Error making request to cell", http.StatusInternalServerError)
				r.errorCount.Add(ctx, 1)
				return
			}
			defer resp.Body.Close()

			log.Printf("Received response from cell: %s, status: %d", cell.Name, resp.StatusCode)

			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)

			r.latency.Record(ctx, float64(time.Since(startTime).Milliseconds()))
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				r.successCount.Add(ctx, 1)
			} else {
				r.errorCount.Add(ctx, 1)
			}
			return
		}
	}

	log.Printf("No matching cell found for client_id: %d", clientID)
	http.Error(w, "No matching cell found", http.StatusNotFound)
	r.errorCount.Add(ctx, 1)
}
