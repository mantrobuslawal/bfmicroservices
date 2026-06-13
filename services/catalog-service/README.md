# Catalog Service

The `catalog-service` owns product catalog data for bfstore.

It provides gRPC APIs for reading catalog information such as products, categories, product variants, product images, and product attribute definitions.

The service is implemented in Go and uses:

- gRPC for service APIs;
- Protobuf-generated contracts from the root `proto` directory;
- MySQL for catalog persistence;
- `database/sql` for database access;
- `otelsql` for database instrumentation;
- `pkg/platform/dbmetrics` for database connection pool metrics;
- shared platform gRPC interceptors;
- standard gRPC health checks;
- structured logging with `log/slog`;
- OpenTelemetry tracing and metrics bootstrap;
- Jaeger trace visualisation through the local OpenTelemetry Collector;
- graceful shutdown for local and container runtime behaviour.

## Runtime architecture

At runtime, the catalog service starts like this:

```text
load configuration
-> create logger
-> initialise telemetry when enabled
-> open instrumented MySQL connection pool
-> register database pool metrics
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
-> close database connection pool
-> flush and shutdown telemetry providers
```

## Observability

The catalog service currently emits:

```text
gRPC request spans
database/sql spans
database connection pool metrics
structured request logs
correlation IDs
```

Current telemetry flow:

```text
catalog-service
  -> OpenTelemetry Collector
  -> Jaeger for traces
  -> Collector debug logs for metrics
```

## Database spans

Database spans are created through `otelsql`.

The instrumentation is configured in:

```text
services/catalog-service/internal/database/mysql.go
```

Expected trace shape in Jaeger:

```text
/bfstore.catalog.v1.CatalogService/ListProducts
  -> database/sql span
  -> database/sql span
```

Repository methods must use context-aware SQL calls:

```go
QueryContext(ctx, ...)
QueryRowContext(ctx, ...)
ExecContext(ctx, ...)
```

## Database metrics

Database connection pool metrics are registered through:

```text
pkg/platform/dbmetrics
```

The catalog service wires these metrics after opening the database connection pool.

Current metric names:

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

These metrics are sourced from:

```go
db.Stats()
```

## Running with telemetry

Start observability services:

```bash
make observability-up
```

Run catalog-service with telemetry enabled:

```bash
cd services/catalog-service

TELEMETRY_ENABLED=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_EXPORTER_OTLP_INSECURE=true \
GRPC_REFLECTION_ENABLED=true \
go run ./cmd/catalog-service
```

Send a request:

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-dbmetrics-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

## Viewing traces

Open Jaeger:

```text
http://localhost:16686
```

Search for:

```text
catalog-service
```

## Viewing metrics

Watch Collector logs:

```bash
docker compose logs -f otel-collector
```

Look for:

```text
db.client.connections.open
db.client.connections.in_use
db.client.connections.idle
db.client.connections.wait_count
db.client.connections.wait_duration
```

## Running tests

Run catalog tests:

```bash
make catalog-test
```

Run database package tests:

```bash
go test ./services/catalog-service/internal/database -v
```

Run DB metrics package tests:

```bash
go test ./pkg/platform/dbmetrics -v
```

Run all tests:

```bash
go test ./...
```

## Troubleshooting

### Jaeger shows gRPC spans but no database spans

Check that repository methods use:

```go
QueryContext
QueryRowContext
ExecContext
```

Also confirm the service was restarted after database instrumentation was added.

### Collector logs show no DB metrics

Check that:

```text
TELEMETRY_ENABLED=true
MetricsEnabled is true
dbmetrics.Register is called after database.Open
Collector metrics pipeline includes debug exporter
```

### Metrics do not change much locally

That may be normal with low local traffic.

Send repeated requests or add a small load test later to create more visible activity.

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
database/sql tracing through otelsql
database pool metrics through dbmetrics
OpenTelemetry Collector integration
Jaeger trace visualisation
Collector debug metric verification
Makefile-driven local smoke tests
```

## Practical rule

```text
Traces explain request paths.
Metrics explain resource health over time.
```

Keep it boring where production matters.
