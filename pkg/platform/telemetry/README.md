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

Package structure:

```text
pkg/platform/telemetry/
├── README.md
├── config.go
├── telemetry.go
└── telemetry_test.go
```

## What problem this solves

Without a shared telemetry bootstrap, every service would need to know how to configure OpenTelemetry providers and exporters.

That would create duplicated setup code across services such as:

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

This package gives each service one consistent setup path.

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

## Key concepts

### Resource

A resource describes the entity producing telemetry.

For a service, that means attributes such as:

```text
service.name=catalog-service
service.version=<git-sha-or-build-version>
deployment.environment.name=local
```

### Tracer provider

A tracer provider creates tracers.

Tracers create spans.

A span represents a unit of work, such as:

```text
handling a gRPC request
running a database query
publishing a Kafka event
calling another service
```

### Meter provider

A meter provider creates meters.

Meters create metrics instruments.

Metrics answer questions such as:

```text
How many requests are happening?
How long are requests taking?
How many errors are occurring?
How many database connections are open?
```

### Exporter

An exporter sends telemetry data out of the process.

This package uses OTLP over gRPC.

A typical local flow is:

```text
catalog-service
  -> OTLP gRPC exporter
  -> OpenTelemetry Collector
  -> Jaeger
```

### Propagator

A propagator carries trace context across process boundaries.

This package configures:

```text
TraceContext
Baggage
```

## Usage

```go
telemetryConfig := telemetry.DefaultConfig("catalog-service")
telemetryConfig.Environment = cfg.Environment
telemetryConfig.ServiceVersion = cfg.ServiceVersion
telemetryConfig.OTLPEndpoint = cfg.OTLPEndpoint
telemetryConfig.OTLPInsecure = cfg.OTLPInsecure
telemetryConfig.TracesEnabled = cfg.TracesEnabled
telemetryConfig.MetricsEnabled = cfg.MetricsEnabled
telemetryConfig.MetricExportInterval = cfg.MetricExportInterval

telemetryRuntime, err := telemetry.Setup(ctx, telemetryConfig)
if err != nil {
    logger.Error("failed to setup telemetry", "error", err)
    os.Exit(1)
}

defer func() {
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := telemetryRuntime.Shutdown(shutdownCtx); err != nil {
        logger.Error("failed to shutdown telemetry", "error", err)
    }
}()
```

## Config

```go
type Config struct {
    ServiceName          string
    ServiceVersion       string
    Environment          string
    OTLPEndpoint         string
    OTLPInsecure         bool
    TracesEnabled        bool
    MetricsEnabled       bool
    MetricExportInterval time.Duration
}
```

Use `DefaultConfig` for local defaults:

```go
cfg := telemetry.DefaultConfig("catalog-service")
```

Defaults:

```text
Environment:          local
OTLPEndpoint:         localhost:4317
OTLPInsecure:         true
TracesEnabled:        true
MetricsEnabled:       true
MetricExportInterval: 30 seconds
```

## Local development

The default endpoint assumes an OpenTelemetry Collector listening on:

```text
localhost:4317
```

For Docker Compose, this will usually become:

```text
otel-collector:4317
```

depending on whether the service runs on the host or inside the Compose network.

## Current scope

This package initialises providers and exporters.

The catalog service also uses gRPC server instrumentation with:

```go
grpc.StatsHandler(otelgrpc.NewServerHandler())
```

That instrumentation belongs in the service gRPC server setup, not in the telemetry bootstrap package.

The package does not yet:

- add database instrumentation;
- add Kafka instrumentation;
- define custom business metrics;
- emit OpenTelemetry logs.

Those should be added in later slices.

## Testing

Run package tests:

```bash
go test ./pkg/platform/telemetry -v
```

The tests focus on configuration behaviour and disabled-signal setup. They deliberately avoid requiring a live OpenTelemetry Collector.

## Design rule

```text
This package initialises telemetry plumbing.
Service code decides what work is worth instrumenting.
```

Keep it boring where production matters.
