# Local Observability

This document explains how to run and verify the local bfstore observability stack.

The current local flow is:

```text
catalog-service
  -> OpenTelemetry SDK
  -> otelgrpc gRPC server instrumentation
  -> OTLP exporter
  -> OpenTelemetry Collector
  -> Jaeger
```

The goal is to prove that a local gRPC request to `catalog-service` produces a trace that can be inspected in Jaeger.

## Components

### catalog-service

The catalog service emits telemetry when enabled.

It uses:

- `pkg/platform/telemetry` for OpenTelemetry bootstrap;
- `otelgrpc.NewServerHandler()` for gRPC request instrumentation;
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

## Docker Compose services

```yaml
otel-collector:
  image: otel/opentelemetry-collector-contrib:latest
  container_name: bfstore-otel-collector
  command: ["--config=/etc/otelcol-contrib/config.yaml"]
  volumes:
    - ./deployments/local/otel-collector/config.yaml:/etc/otelcol-contrib/config.yaml:ro
  ports:
    - "4317:4317"
    - "4318:4318"
  depends_on:
    - jaeger
  networks:
    - bfstore-local

jaeger:
  image: jaegertracing/all-in-one:latest
  container_name: bfstore-jaeger
  environment:
    COLLECTOR_OTLP_ENABLED: "true"
  ports:
    - "16686:16686"
  networks:
    - bfstore-local
```

For repeatable builds, prefer pinned image versions later instead of `latest`.

## Running the stack

Start observability services:

```bash
docker compose up -d otel-collector jaeger
```

Or with Make targets if available:

```bash
make observability-up
```

Watch Collector logs:

```bash
docker compose logs -f otel-collector
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

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-jaeger-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Expected behaviour:

```text
request succeeds or fails normally
catalog-service logs include correlation_id=local-dev-jaeger-123
Collector logs show received telemetry
Jaeger shows a trace for catalog-service
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

Useful things to inspect:

```text
operation name
duration
service name
span attributes
error status if a request failed
timeline
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
		-H 'x-correlation-id: local-dev-jaeger-123' \
		-d '{"page":{"page_size":5}}' \
		localhost:50051 \
		bfstore.catalog.v1.CatalogService/ListProducts
```

Remember: Makefile recipe lines must use tabs, not spaces.

## Troubleshooting

### Collector says `protocols.grpc expected a map or struct got string`

Correct:

```yaml
grpc:
  endpoint: 0.0.0.0:4317
```

Incorrect:

```yaml
grpc: 0.0.0.0:4317
```

### Collector still fails after adding `endpoint`

Check spacing after the colon.

Correct:

```yaml
endpoint: 0.0.0.0:4317
```

Incorrect:

```yaml
endpoint:0.0.0.0:4317
```

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

### Collector receives telemetry but Jaeger is empty

Check that the Collector trace pipeline includes the Jaeger exporter:

```yaml
service:
  pipelines:
    traces:
      exporters: [debug, otlp/jaeger]
```

Also check:

```yaml
otlp/jaeger:
  endpoint: jaeger:4317
  tls:
    insecure: true
```

## Design notes

The application should export telemetry to the Collector, not directly to Jaeger.

Preferred:

```text
catalog-service
  -> OpenTelemetry Collector
  -> Jaeger
```

Avoid:

```text
catalog-service
  -> Jaeger directly
```

The Collector gives the platform one routing point for telemetry and lets the backend change later without changing service code.

## Current milestone

This local observability setup proves:

```text
catalog-service emits OpenTelemetry data
otelgrpc creates gRPC request spans
Collector receives OTLP telemetry
Collector exports traces to Jaeger
Jaeger displays service traces locally
```

Keep it boring where production matters.
