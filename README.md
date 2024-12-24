# Cell Router

## Overview

The Cell Router is a service that routes requests to different cells based on the `client_id` provided in the query parameters. It uses OpenTelemetry for metrics and tracing, Prometheus for metrics collection, and Grafana for visualization.

## Libraries and Frameworks Used

- **OpenTelemetry**: For collecting metrics and tracing.
- **Prometheus**: For scraping and storing metrics.
- **Grafana**: For visualizing metrics.

## Build the Project

To build the project, run the following command:

```sh
docker-compose build
```

## Run Locally

To run the project locally, use the following command:

```sh
docker-compose up
```

This will start the following services:
- `cell-router` on port 8080
- `cell1` mock cell on port 8081
- `cell2` mock cell on port 8082
- `prometheus` on port 9090
- `grafana` on port 3000

## Test Locally

### Test Routing

To test the routing, you can use `curl` or any HTTP client to send requests to the `cell-router` service:

```sh
# Test routing to cell1
curl "http://localhost:8080?client_id=50"

# Test routing to cell2
curl "http://localhost:8080?client_id=150"
```

You should receive responses from the respective mock cells based on the `client_id` provided.

### Access Prometheus

Prometheus can be accessed at [http://localhost:9090](http://localhost:9090). You can use Prometheus to query metrics collected from the `cell-router`.

### Access Grafana

Grafana can be accessed at [http://localhost:3000](http://localhost:3000). The default login credentials are:
- **Username**: admin
- **Password**: admin

After logging in, you can import the dashboard configuration from `grafana/dashboards/cell-router-dashboard.json` to visualize the metrics.

## Metrics Collected

The following metrics are collected and available for visualization:

- **Number of requests per second**: `rate(requests_total[1m])`
- **Latency of requests**: `request_latency`
  - Percentiles: p50, p90, p99, p100
- **Availability of the cell router**: Derived from the success and error rates.
- **Number of successful requests**: `rate(requests_success_total[1m])`
- **Number of failed requests**: `rate(requests_error_total[1m])`
  - Includes 3xx, 4xx, and 5xx status codes

## Distributed Tracing

Distributed tracing is enabled using OpenTelemetry. Traces can be collected and visualized using compatible tracing backends. Ensure that the tracing backend is properly configured to collect traces from the `cell-router`.