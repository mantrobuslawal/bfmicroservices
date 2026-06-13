# Platform Telemetry

The `telemetry` package provides shared OpenTelemetry bootstrap code for bfstore services.

It is responsible for initialising service-level telemetry plumbing:

- service identity;
- resource attributes;
- trace provider setup;
- metric provider setup;
- OTLP exporter configuration;
- global trace context propagation;
- graceful telemetry shutdown.

It should stay small and boring. Service code should use this package to initialise telemetry during startup, then add instrumentation deliberately where it is useful.

## Package location

```text
pkg/platform/telemetry
```

## Mental model

Telemetry is runtime evidence about what the service is doing.

For bfstore, think of the signals like this:

```text
logs
  What happened?

traces
  Where did this request go?

metrics
  How much, how often, how slow, how many?
```

This package focuses on OpenTelemetry traces and metrics.

Logs remain handled by `log/slog` for now.

## Instrumentation boundaries

This package initialises telemetry providers and exporters.

It does not instrument every technology directly.

Current instrumentation points:

```text
gRPC server instrumentation
  configured in the service gRPC server setup with otelgrpc

database/sql instrumentation
  configured in the service database package with otelsql
```

This separation keeps the platform telemetry package focused.

## gRPC instrumentation

The catalog service uses:

```go
grpc.StatsHandler(otelgrpc.NewServerHandler())
```

This creates gRPC server spans for incoming catalog requests.

## Database instrumentation

The catalog service instruments MySQL through `database/sql` using:

```text
github.com/XSAM/otelsql
```

That instrumentation is configured in:

```text
services/catalog-service/internal/database/mysql.go
```

The database flow is:

```text
catalog repository
  -> *sql.DB
  -> otelsql instrumented driver
  -> go-sql-driver/mysql
  -> MySQL
```

This produces database spans underneath request traces when repository code uses context-aware SQL calls:

```go
QueryContext(ctx, ...)
QueryRowContext(ctx, ...)
ExecContext(ctx, ...)
```

The platform telemetry package should not need to know the details of each service database driver.

## Local development

When running the service from the host, export telemetry to:

```text
localhost:4317
```

When running inside Docker Compose, export telemetry to:

```text
otel-collector:4317
```

The local observability flow is:

```text
catalog-service
  -> OpenTelemetry Collector
  -> Jaeger
```

## Testing

Run package tests:

```bash
go test ./pkg/platform/telemetry -v
```

Database instrumentation tests live closer to the database package:

```bash
go test ./services/catalog-service/internal/database -v
```

Those unit tests should verify driver registration behaviour without requiring a live MySQL instance.

## Current scope

This package currently supports:

```text
resource creation
trace provider setup
metric provider setup
OTLP exporter setup
global propagator setup
graceful shutdown
```

It does not yet:

```text
define custom business metrics
instrument Kafka
instrument outbound gRPC clients
emit OpenTelemetry logs
provide production Collector configuration
```

Those should be added in later slices.

## Design rule

```text
The telemetry package initialises telemetry plumbing.
Service packages instrument their own stable technology boundaries.
Repository code remains context-led and boring.
```

Keep it boring where production matters.
