# Local Observability

This document explains how to run and verify the local bfstore observability stack.

The current local flow is:

```text
catalog-service
  -> OpenTelemetry SDK
  -> otelgrpc gRPC server instrumentation
  -> otelsql database/sql instrumentation
  -> OTLP exporter
  -> OpenTelemetry Collector
  -> Jaeger
```

The goal is to prove that a local gRPC request to `catalog-service` produces a trace that can be inspected in Jaeger, including database spans for catalog repository work.

## Components

### catalog-service

The catalog service emits telemetry when enabled.

It uses:

- `pkg/platform/telemetry` for OpenTelemetry bootstrap;
- `otelgrpc.NewServerHandler()` for gRPC request instrumentation;
- `otelsql` for database/sql instrumentation;
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

## Sending a trace-producing request

Send a catalog request with an explicit correlation ID:

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-db-otel-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Expected behaviour:

```text
request succeeds or fails normally
catalog-service logs include correlation_id=local-dev-db-otel-123
Collector logs show received telemetry
Jaeger shows a trace for catalog-service
Jaeger trace includes database spans underneath the gRPC span
```

## Expected trace shape

For `ListProducts`, the trace should look roughly like:

```text
/bfstore.catalog.v1.CatalogService/ListProducts
  -> database/sql span
  -> database/sql span
```

The exact span names may vary.

The important result is that database work appears underneath the request that caused it.

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
		-H 'x-correlation-id: local-dev-db-otel-123' \
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

### Database spans appear as separate traces

This usually means the request context is not flowing into the database call.

Check the call chain:

```text
handler ctx
-> service ctx
-> repository ctx
-> QueryContext(ctx, ...)
```

## Current milestone

This local observability setup proves:

```text
catalog-service emits OpenTelemetry data
otelgrpc creates gRPC request spans
otelsql creates database spans
context flows from gRPC handler to repository
Collector receives OTLP telemetry
Collector exports traces to Jaeger
Jaeger displays request traces with database child spans
```

Keep it boring where production matters.
