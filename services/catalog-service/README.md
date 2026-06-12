# Catalog Service

The `catalog-service` owns product catalog data for bfstore.

It provides gRPC APIs for reading catalog information such as products, categories, product variants, product images, and product attribute definitions.

The service is implemented in Go and uses:

- gRPC for service APIs;
- Protobuf-generated contracts from the root `proto` directory;
- MySQL for catalog persistence;
- shared platform gRPC interceptors;
- standard gRPC health checks;
- structured logging with `log/slog`;
- OpenTelemetry tracing and metrics bootstrap;
- Jaeger trace visualisation through the local OpenTelemetry Collector;
- graceful shutdown for local and container runtime behaviour.

## Responsibilities

The catalog service is responsible for:

- serving product catalog read APIs;
- retrieving product and category data from MySQL;
- mapping catalog domain models to Protobuf responses;
- exposing standard gRPC health status;
- participating in shared platform runtime behaviour such as logging, recovery, correlation ID propagation, and telemetry.

The catalog service is not responsible for:

- basket management;
- inventory reservation;
- order orchestration;
- payment processing;
- shipping;
- notification delivery;
- search indexing;
- recommendation logic.

Those responsibilities belong to separate bfstore services.

## Runtime architecture

At runtime, the catalog service starts like this:

```text
load configuration
-> create logger
-> initialise telemetry when enabled
-> open MySQL connection
-> run catalog readiness check
-> create catalog repository
-> create catalog service
-> create gRPC server
-> register platform interceptors
-> register OpenTelemetry gRPC instrumentation
-> register catalog gRPC handler
-> register gRPC health service
-> optionally register gRPC reflection
-> start serving requests
```

During shutdown:

```text
receive SIGINT or SIGTERM
-> mark service NOT_SERVING
-> stop accepting new gRPC traffic
-> allow in-flight requests to finish
-> force stop if graceful shutdown times out
-> close database connection
-> flush and shutdown telemetry providers
```

## gRPC server behaviour

The catalog service uses shared platform interceptors from:

```text
pkg/platform/grpc/interceptors
```

Recommended chain:

```go
grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()),
    grpc.ChainUnaryInterceptor(
        platforminterceptors.UnaryRecoveryInterceptor(logger),
        platforminterceptors.UnaryCorrelationIDInterceptor(),
        platforminterceptors.UnaryLoggingInterceptor(logger),
    ),
)
```

The order matters:

```text
OpenTelemetry gRPC stats handler
-> Recovery
-> Correlation ID
-> Logging
-> Handler
```

### OpenTelemetry gRPC instrumentation

The service uses `otelgrpc.NewServerHandler()` through `grpc.StatsHandler`.

This adds automatic gRPC telemetry for incoming requests so they can appear as spans in a tracing backend such as Jaeger.

The application still exports telemetry to the OpenTelemetry Collector rather than directly to Jaeger.

```text
catalog-service
  -> OpenTelemetry SDK
  -> OpenTelemetry Collector
  -> Jaeger
```

### Recovery interceptor

The recovery interceptor catches panics from later interceptors and handlers.

It:

- prevents one bad request from crashing the whole process;
- logs panic details server-side;
- returns a safe `codes.Internal` response to the client;
- avoids exposing stack traces or implementation details to callers.

### Correlation ID interceptor

The correlation ID interceptor ensures every request has a correlation ID.

It:

- reads `x-correlation-id` from incoming gRPC metadata;
- reuses the incoming value when present;
- generates a new correlation ID when missing;
- stores the correlation ID in request context;
- returns the correlation ID in response metadata.

Metadata key:

```text
x-correlation-id
```

### Logging interceptor

The logging interceptor writes structured logs for each unary gRPC request.

Current fields include:

```text
grpc.method
grpc.code
duration_ms
correlation_id
error
```

Example successful request log fields:

```text
grpc.method=/bfstore.catalog.v1.CatalogService/ListProducts
grpc.code=OK
duration_ms=12
correlation_id=local-dev-123
```

## Health checks

The catalog service exposes the standard gRPC health API:

```text
grpc.health.v1.Health
```

Health status is managed by the shared platform health manager:

```text
pkg/platform/healthcheck
```

The catalog service owns its own readiness truth through:

```text
services/catalog-service/internal/health
```

## gRPC reflection

gRPC reflection can be enabled for local development and testing.

Enable reflection with:

```bash
GRPC_REFLECTION_ENABLED=true go run ./cmd/catalog-service
```

Reflection should be used for local development and testing. It should not be enabled by default in production.

## Telemetry

Telemetry can be enabled locally with:

```bash
TELEMETRY_ENABLED=true
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true
```

When running the service from the host with `go run`, use:

```text
localhost:4317
```

When running the service from inside Docker Compose, use:

```text
otel-collector:4317
```

Useful local telemetry flow:

```text
catalog-service
  -> OTLP gRPC localhost:4317
  -> otel-collector
  -> jaeger
```

## Running locally

From the repository root, start dependencies:

```bash
make up
```

Start observability services if needed:

```bash
make observability-up
```

Then start the catalog service:

```bash
make catalog-run
```

Run with telemetry enabled:

```bash
cd services/catalog-service

TELEMETRY_ENABLED=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_EXPORTER_OTLP_INSECURE=true \
GRPC_REFLECTION_ENABLED=true \
go run ./cmd/catalog-service
```

## Running tests

```bash
make catalog-test
make catalog-integration-test
make test
go test ./pkg/platform/grpc/interceptors -v
go test ./pkg/platform/healthcheck -v
go test ./pkg/platform/telemetry -v
```

## Smoke testing with grpcurl

List available gRPC services:

```bash
make catalog-grpc-list
```

Or directly:

```bash
grpcurl -plaintext localhost:50051 list
```

Expected services include:

```text
bfstore.catalog.v1.CatalogService
grpc.health.v1.Health
grpc.reflection.v1.ServerReflection
```

## Health check smoke test

```bash
make catalog-health
```

Or directly:

```bash
grpcurl -plaintext \
  -d '{}' \
  localhost:50051 \
  grpc.health.v1.Health/Check
```

Expected response:

```json
{
  "status": "SERVING"
}
```

## Catalog API smoke tests

List products:

```bash
grpcurl -plaintext \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

List categories:

```bash
grpcurl -plaintext \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListCategories
```

Send a request with an explicit correlation ID:

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
logs include correlation_id=local-dev-jaeger-123
Jaeger shows a catalog-service trace when telemetry is enabled
```

## Jaeger

When observability is running, open:

```text
http://localhost:16686
```

Search for:

```text
catalog-service
```

## Troubleshooting

### Jaeger does not show traces

Check:

```bash
docker compose logs -f otel-collector
docker compose logs -f jaeger
```

Confirm the catalog service was started with:

```bash
TELEMETRY_ENABLED=true
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true
```

Then send a fresh `grpcurl` request.

### Collector config says protocol expected map or struct

Correct:

```yaml
grpc:
  endpoint: 0.0.0.0:4317
http:
  endpoint: 0.0.0.0:4318
```

Incorrect:

```yaml
grpc: 0.0.0.0:4317
http: 0.0.0.0:4318
```

Also make sure there is a space after `endpoint:`.

Correct:

```yaml
endpoint: 0.0.0.0:4317
```

Incorrect:

```yaml
endpoint:0.0.0.0:4317
```

## Current runtime foundation

The service currently demonstrates:

```text
gRPC API serving
standard gRPC health checks
gRPC reflection for local development
structured request logging
correlation ID propagation
panic recovery
graceful shutdown
database readiness checks
OpenTelemetry bootstrap
gRPC server tracing instrumentation
OpenTelemetry Collector integration
Jaeger trace visualisation
Makefile-driven local smoke tests
```

## Practical rule

```text
Platform packages provide reusable plumbing.
Service packages own service-specific truth.
```

Keep it boring where production matters.

