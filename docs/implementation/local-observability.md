# Local Observability

This document explains how to run and verify the local bfstore observability stack.

Current local flow:

```text
catalog-service
  -> OpenTelemetry SDK
  -> otelgrpc gRPC server instrumentation
  -> otelsql database/sql instrumentation
  -> dbmetrics database pool metrics
  -> OTLP exporter
  -> OpenTelemetry Collector
  -> Jaeger for traces
  -> Prometheus for metrics
```

## Current milestone

The local observability setup currently proves:

```text
catalog-service emits OpenTelemetry data
otelgrpc creates gRPC request spans
otelsql creates database spans
dbmetrics emits database connection pool metrics
Collector receives OTLP telemetry
Collector exports traces to Jaeger
Collector exposes metrics for Prometheus
Prometheus scrapes metrics from the Collector
Jaeger displays request traces with database child spans
Prometheus queries database pool metrics
```

## Components

### catalog-service

The catalog service emits telemetry when enabled.

It uses:

- `pkg/platform/telemetry` for OpenTelemetry bootstrap;
- `otelgrpc.NewServerHandler()` for gRPC request instrumentation;
- `otelsql` for database/sql tracing;
- `pkg/platform/dbmetrics` for database connection pool metrics;
- platform interceptors for recovery, correlation ID propagation, and structured request logging.

### OpenTelemetry Collector

The Collector receives telemetry from services and routes it to one or more backends.

Local endpoints:

```text
OTLP gRPC: 4317
OTLP HTTP: 4318
Prometheus exporter: 9464
```

Current Collector routing:

```text
traces -> Jaeger and debug logs
metrics -> debug logs and Prometheus scrape endpoint
```

### Jaeger

Jaeger provides a local trace UI:

```text
http://localhost:16686
```

Search for:

```text
catalog-service
```

### Prometheus

Prometheus provides local metrics storage and querying:

```text
http://localhost:9090
```

Target health:

```text
http://localhost:9090/targets
```

Expected target:

```text
otel-collector
```

Expected state:

```text
UP
```

## Collector config

Recommended path:

```text
deployments/local/otel-collector/config.yaml
```

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:

exporters:
  debug:
    verbosity: detailed

  otlp/jaeger:
    endpoint: jaeger:4317
    tls:
      insecure: true

  prometheus:
    endpoint: 0.0.0.0:9464

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug, otlp/jaeger]

    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug, prometheus]
```

## Prometheus config

Recommended path:

```text
deployments/local/prometheus/prometheus.yml
```

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: otel-collector
    static_configs:
      - targets: ["otel-collector:9464"]
```

## Running observability locally

```bash
make observability-up
```

Or:

```bash
docker compose up -d otel-collector jaeger prometheus
```

Check status:

```bash
docker compose ps
```

Watch logs:

```bash
docker compose logs -f otel-collector jaeger prometheus
```

## Running catalog-service with telemetry

From the host:

```bash
cd services/catalog-service

TELEMETRY_ENABLED=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_EXPORTER_OTLP_INSECURE=true \
GRPC_REFLECTION_ENABLED=true \
go run ./cmd/catalog-service
```

If running inside Docker Compose, use:

```text
otel-collector:4317
```

## Sending a request

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-prometheus-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Expected behaviour:

```text
request succeeds or fails normally
catalog-service logs include correlation_id=local-dev-prometheus-123
Jaeger shows a trace for catalog-service
Jaeger trace includes database spans under the gRPC span
Prometheus can query database pool metrics
```

## Viewing traces in Jaeger

Open:

```text
http://localhost:16686
```

Expected trace shape:

```text
/bfstore.catalog.v1.CatalogService/ListProducts
  -> database/sql span
  -> database/sql span
```

## Viewing metrics in Prometheus

Open:

```text
http://localhost:9090
```

Prometheus normalises metric names from dots to underscores.

Try:

```promql
db_client_connections_open
```

```promql
db_client_connections_in_use
```

```promql
db_client_connections_idle
```

```promql
db_client_connections_wait_count
```

```promql
db_client_connections_wait_duration
```

Useful rate queries:

```promql
rate(db_client_connections_wait_count[5m])
```

```promql
rate(db_client_connections_wait_duration[5m])
```

## Suggested Make targets

```makefile
.PHONY: observability-up
observability-up:
	docker compose -f $(COMPOSE_FILE) up -d otel-collector jaeger prometheus

.PHONY: observability-logs
observability-logs:
	docker compose -f $(COMPOSE_FILE) logs -f otel-collector jaeger prometheus

.PHONY: metrics-up
metrics-up:
	docker compose -f $(COMPOSE_FILE) up -d otel-collector prometheus

.PHONY: metrics-logs
metrics-logs:
	docker compose -f $(COMPOSE_FILE) logs -f otel-collector prometheus

.PHONY: catalog-run-telemetry
catalog-run-telemetry:
	cd services/catalog-service && \
		TELEMETRY_ENABLED=true \
		OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
		OTEL_EXPORTER_OTLP_INSECURE=true \
		GRPC_REFLECTION_ENABLED=true \
		go run ./cmd/catalog-service
```

Remember: Makefile recipe lines must use tabs.

## Troubleshooting

### Prometheus target is down

Open:

```text
http://localhost:9090/targets
```

Check target:

```text
otel-collector:9464
```

Check logs:

```bash
docker compose logs -f prometheus
docker compose logs -f otel-collector
```

### Prometheus target is up but metrics are missing

Check:

```text
catalog-service is running
TELEMETRY_ENABLED=true
MetricsEnabled is true
dbmetrics.Register is called
a fresh catalog request has been sent
```

### Jaeger has no traces

Check:

```text
otel-collector is running
jaeger is running
catalog-service started with TELEMETRY_ENABLED=true
OTEL_EXPORTER_OTLP_ENDPOINT is correct for host vs container
```

## Next step

The next observability step is Grafana:

```text
catalog-service
  -> OpenTelemetry Collector
  -> Jaeger for traces
  -> Prometheus for metrics
  -> Grafana for dashboards
```

Keep it boring where production matters.
