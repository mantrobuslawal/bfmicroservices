# Catalog Service

The `catalog-service` owns product catalog data for bfstore.

It provides gRPC APIs for reading catalog information such as products, categories, product variants, product images, and product attribute definitions.

The service is implemented in Go and uses:

- gRPC for service APIs;
- Protobuf-generated contracts from the root `proto` directory;
- MySQL for catalog persistence;
- `database/sql` for database access;
- `otelsql` for database instrumentation;
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

## Request path

A typical catalog request flows through the service like this:

```text
grpc client
  -> catalog gRPC server
  -> otelgrpc server instrumentation
  -> recovery interceptor
  -> correlation ID interceptor
  -> logging interceptor
  -> catalog gRPC handler
  -> catalog application service
  -> catalog repository
  -> instrumented database/sql driver
  -> MySQL
```

When telemetry is enabled, Jaeger should show a trace shaped roughly like:

```text
/bfstore.catalog.v1.CatalogService/ListProducts
  -> database/sql span
  -> database/sql span
```

The exact database span names may vary depending on the SQL operation and the instrumentation library behaviour.

## Database instrumentation

The catalog service instruments MySQL at the database connection boundary:

```text
services/catalog-service/internal/database/mysql.go
```

This is deliberate.

Instrumentation belongs at the connection boundary first, not scattered through every repository method.

The service uses:

```text
github.com/XSAM/otelsql
```

to wrap Go's standard `database/sql` MySQL driver.

The resulting flow is:

```text
catalog repository
  -> *sql.DB
  -> otelsql instrumented driver wrapper
  -> go-sql-driver/mysql
  -> MySQL
```

## Driver registration

The instrumented SQL driver should be registered once per process.

The recommended pattern is:

```go
var (
    registerInstrumentedDriverOnce sync.Once
    registerInstrumentedDriverName string
    registerInstrumentedDriverErr  error
)
```

Then:

```go
func instrumentedMySQLDriver() (string, error) {
    registerInstrumentedDriverOnce.Do(func() {
        registerInstrumentedDriverName, registerInstrumentedDriverErr = otelsql.Register(baseMySQLDriverName)
    })

    if registerInstrumentedDriverErr != nil {
        return "", registerInstrumentedDriverErr
    }

    return registerInstrumentedDriverName, nil
}
```

This avoids duplicate SQL driver registration if `Open` is called more than once in a process.

## Context propagation requirement

Database spans attach correctly to request traces when repository methods use context-aware database calls:

```go
db.QueryContext(ctx, query, args...)
db.QueryRowContext(ctx, query, args...)
db.ExecContext(ctx, query, args...)
```

Avoid non-context calls in request paths:

```go
db.Query(query, args...)
db.QueryRow(query, args...)
db.Exec(query, args...)
```

The `ctx` from the gRPC request should flow through:

```text
gRPC handler
  -> catalog service method
  -> repository method
  -> QueryContext / QueryRowContext / ExecContext
```

That is what makes the DB spans appear underneath the gRPC request span in Jaeger.

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

## Smoke testing database spans

Run the catalog service with telemetry enabled:

```bash
cd services/catalog-service

TELEMETRY_ENABLED=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_EXPORTER_OTLP_INSECURE=true \
GRPC_REFLECTION_ENABLED=true \
go run ./cmd/catalog-service
```

Send a request with an explicit correlation ID:

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-db-otel-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Open Jaeger:

```text
http://localhost:16686
```

Search for:

```text
catalog-service
```

Expected result:

```text
a catalog gRPC request span
one or more database/sql child spans
```

## Troubleshooting

### Jaeger shows gRPC spans but no DB spans

Check that repository methods use context-aware SQL calls:

```go
QueryContext
QueryRowContext
ExecContext
```

Also confirm telemetry was enabled before the request was sent.

### DB spans appear as separate traces

This usually means request context is not flowing into the repository.

Check the call chain:

```text
handler ctx
-> service ctx
-> repository ctx
-> QueryContext(ctx, ...)
```

### SQL driver registration fails

If you see a duplicate registration error, ensure the instrumented driver is registered with `sync.Once`.

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
database/sql instrumentation through otelsql
OpenTelemetry Collector integration
Jaeger trace visualisation
Makefile-driven local smoke tests
```

## Practical rule

```text
Platform packages provide reusable plumbing.
Service packages own service-specific truth.
Instrumentation belongs at stable boundaries first.
```

Keep it boring where production matters.
