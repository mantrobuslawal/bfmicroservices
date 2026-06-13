# Local Observability

This document explains how to run and verify the local bfstore observability stack.

The current local flow is:

```text
catalog-service
  -> OpenTelemetry SDK
  -> otelgrpc gRPC server instrumentation
  -> otelsql database/sql instrumentation
  -> dbmetrics database pool metrics
  -> OTLP exporter
  -> OpenTelemetry Collector
  -> Jaeger for traces
  -> Collector debug logs for metrics
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
Collector debug logs show metrics
Jaeger displays request traces with database child spans
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

In local development, it receives OTLP gRPC traffic on:

```text
localhost:4317
```

It also exposes OTLP HTTP on:

```text
localhost:4318
```

The current Collector config exports:

```text
traces -> Jaeger and debug logs
metrics -> debug logs
```

### Jaeger

Jaeger provides a local trace UI.

Open the UI at:

```text
http://localhost:16686
```

Search for:

```text
catalog-service
```

## Collector config

Recommended local config:

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

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug, otlp/jaeger]

    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]
```

Recommended path:

```text
deployments/local/otel-collector/config.yaml
```

## Running catalog-service with telemetry

If running the service from the host with `go run`:

```bash
cd services/catalog-service

TELEMETRY_ENABLED=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_EXPORTER_OTLP_INSECURE=true \
GRPC_REFLECTION_ENABLED=true \
go run ./cmd/catalog-service
```

If running the service inside Docker Compose, the endpoint should normally be:

```text
otel-collector:4317
```

## Sending a trace-producing and metric-producing request

Send a catalog request with an explicit correlation ID:

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-dbmetrics-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Expected behaviour:

```text
request succeeds or fails normally
catalog-service logs include correlation_id=local-dev-dbmetrics-123
Collector logs show received telemetry
Jaeger shows a trace for catalog-service
Jaeger trace includes database spans underneath the gRPC span
Collector logs include database pool metrics
```

## Viewing traces in Jaeger

Open:

```text
http://localhost:16686
```

Select service:

```text
catalog-service
```

Click **Find Traces**.

Expected trace shape:

```text
/bfstore.catalog.v1.CatalogService/ListProducts
  -> database/sql span
  -> database/sql span
```

The exact span names may vary.

## Viewing metrics in Collector logs

Watch Collector logs:

```bash
docker compose logs -f otel-collector
```

Look for metric names such as:

```text
db.client.connections.max
db.client.connections.open
db.client.connections.in_use
db.client.connections.idle
db.client.connections.wait_count
db.client.connections.wait_duration
db.client.connections.max_idle_closed
db.client.connections.max_idle_time_closed
db.client.connections.max_lifetime_closed
```

## Suggested Make targets

```makefile
.PHONY: observability-up
observability-up:
	docker compose -f $(COMPOSE_FILE) up -d otel-collector jaeger

.PHONY: observability-logs
observability-logs:
	docker compose -f $(COMPOSE_FILE) logs -f otel-collector jaeger

.PHONY: catalog-run-telemetry
catalog-run-telemetry:
	cd services/catalog-service && \
		TELEMETRY_ENABLED=true \
		OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
		OTEL_EXPORTER_OTLP_INSECURE=true \
		GRPC_REFLECTION_ENABLED=true \
		go run ./cmd/catalog-service

.PHONY: catalog-list-products-with-correlation
catalog-list-products-with-correlation:
	grpcurl -plaintext \
		-H 'x-correlation-id: local-dev-dbmetrics-123' \
		-d '{"page":{"page_size":5}}' \
		localhost:50051 \
		bfstore.catalog.v1.CatalogService/ListProducts
```

Remember: Makefile recipe lines must use tabs, not spaces.

## Troubleshooting

### Jaeger has no traces

Check that:

```text
otel-collector is running
jaeger is running
catalog-service started with TELEMETRY_ENABLED=true
OTEL_EXPORTER_OTLP_ENDPOINT is correct for host vs container
a fresh grpcurl request was sent after startup
```

Useful commands:

```bash
docker compose ps
docker compose logs -f otel-collector
docker compose logs -f jaeger
```

### Jaeger shows gRPC traces but no database spans

Check that:

```text
catalog-service was restarted after database instrumentation was added
repository methods use QueryContext / QueryRowContext / ExecContext
ctx is passed from handler to service to repository
```

### Collector logs show traces but no DB metrics

Check that:

```text
MetricsEnabled is true
pkg/platform/dbmetrics is wired in catalog-service startup
dbmetrics.Register is called after database.Open
Collector metrics pipeline includes the debug exporter
```

### Metrics appear but values stay low

That may be normal in local development.

Send repeated requests or run a small load test later to create more visible pool activity.

## Next step

The next observability backend is Prometheus.

The future local flow should become:

```text
catalog-service
  -> OpenTelemetry Collector
  -> Jaeger for traces
  -> Prometheus for metrics
```

After Prometheus, Grafana can be added for dashboards.

Keep it boring where production matters.
