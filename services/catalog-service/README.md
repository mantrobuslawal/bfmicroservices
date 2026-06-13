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
- OpenTelemetry Collector for telemetry routing;
- Jaeger for trace visualisation;
- Prometheus for metrics querying;
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
  -> Prometheus for metrics
```

## Database spans

Database spans are created through `otelsql`.

Instrumentation location:

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

In Prometheus, these commonly appear as:

```text
db_client_connections_open
db_client_connections_in_use
db_client_connections_idle
db_client_connections_wait_count
db_client_connections_wait_duration
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
  -H 'x-correlation-id: local-dev-prometheus-123' \
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

Open Prometheus:

```text
http://localhost:9090
```

Check target health:

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

Try PromQL queries:

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
rate(db_client_connections_wait_count[5m])
```

```promql
rate(db_client_connections_wait_duration[5m])
```

## Running tests

```bash
make catalog-test
go test ./services/catalog-service/internal/database -v
go test ./pkg/platform/dbmetrics -v
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

### Prometheus target is down

Check:

```bash
docker compose logs -f prometheus
docker compose logs -f otel-collector
```

Confirm Prometheus scrapes:

```text
otel-collector:9464
```

### Prometheus target is up but no DB metrics appear

Check:

```text
TELEMETRY_ENABLED=true
MetricsEnabled is true
dbmetrics.Register is called after database.Open
a fresh catalog request has been sent
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
database/sql tracing through otelsql
database pool metrics through dbmetrics
OpenTelemetry Collector integration
Jaeger trace visualisation
Prometheus metric querying
Makefile-driven local smoke tests
```

## Practical rule

```text
Traces explain request paths.
Metrics explain resource health over time.
Dashboards should come after queries work.
```

Keep it boring where production matters.
