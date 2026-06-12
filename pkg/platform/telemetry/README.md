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

Suggested package structure:

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

Logs remain handled by `log/slog` for now. The OpenTelemetry Go logs signal is still experimental, so bfstore should keep logging simple until the logs signal is more stable.

## Key concepts

### Resource

A resource describes the entity producing telemetry.

For a service, that means attributes such as:

```text
service.name=catalog-service
service.version=<git-sha-or-build-version>
deployment.environment.name=local
```

Resource attributes help telemetry backends group data by service and environment.

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

This package uses OTLP over gRPC, which is the standard OpenTelemetry Protocol export path.

A typical local flow is:

```text
catalog-service
  -> OTLP gRPC exporter
  -> OpenTelemetry Collector
  -> backend such as Jaeger, Tempo, Prometheus, or Grafana
```

### Propagator

A propagator carries trace context across process boundaries.

For example:

```text
api-gateway
  -> catalog-service
  -> inventory-service
```

The propagator allows those service calls to become part of the same distributed trace.

This package configures:

```text
TraceContext
Baggage
```

## Usage

In a service startup path, create a config and call `Setup`.

Example:

```go
telemetryConfig := telemetry.DefaultConfig("catalog-service")
telemetryConfig.Environment = cfg.Environment
telemetryConfig.OTLPEndpoint = cfg.OTLPEndpoint
telemetryConfig.OTLPInsecure = cfg.OTLPInsecure

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

For Docker Compose, this will usually become something like:

```text
otel-collector:4317
```

depending on whether the service runs on the host or inside the Compose network.

## Current scope

This first telemetry package only initialises providers and exporters.

It does not yet:

- add gRPC server instrumentation;
- add database instrumentation;
- add Kafka instrumentation;
- define custom business metrics;
- configure an OpenTelemetry Collector;
- emit OpenTelemetry logs.

Those should be added in later slices.

## Recommended next steps

After this package is in place:

```text
1. Add telemetry config fields to catalog-service config.
2. Call telemetry.Setup from catalog-service main.go.
3. Add graceful telemetry shutdown.
4. Add gRPC server instrumentation with otelgrpc.
5. Add an OpenTelemetry Collector to docker-compose.
6. Run a local trace smoke test.
```

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

