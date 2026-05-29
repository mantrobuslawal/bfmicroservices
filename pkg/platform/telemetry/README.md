# Telemetry Package

This package will contain shared OpenTelemetry setup helpers for **bfstore** Go services.

Recommended location:

```text
pkg/platform/telemetry/
```

Suggested future files:

```text
pkg/platform/telemetry/
├── README.md
├── tracing.go
├── metrics.go
├── logging.go
└── propagation.go
```

---

## Purpose

This package should provide reusable helpers for:

```text
initialising OpenTelemetry resources
setting service.name and environment attributes
creating tracer providers
creating meter providers
configuring OTLP exporters
configuring context propagation
supporting graceful telemetry shutdown
linking logs with trace/correlation context
```

It should keep telemetry setup consistent across bfstore services.

---

## Why This Lives Under `pkg/platform`

Telemetry is a shared platform concern.

It should be consistent across:

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

A shared package avoids each service copying slightly different OpenTelemetry setup.

---

## Expected Service Usage

Future service startup shape:

```go
telemetryShutdown, err := telemetry.Init(ctx, telemetry.Config{
    ServiceName: "bfstore-catalog-service",
    Environment: "local",
    Version: version,
    OTLPEndpoint: cfg.OTLPEndpoint,
})
if err != nil {
    return err
}
defer telemetryShutdown(context.Background())
```

This is a future shape, not a requirement for the first implementation.

---

## Configuration

Suggested config values:

```text
service_name
service_version
deployment_environment
otlp_endpoint
metrics_enabled
traces_enabled
logs_enabled
sampling_ratio
```

Environment variables may include:

```text
OTEL_SERVICE_NAME
OTEL_EXPORTER_OTLP_ENDPOINT
OTEL_RESOURCE_ATTRIBUTES
OTEL_TRACES_SAMPLER
```

---

## Tracing

Tracing helpers should support:

```text
service tracer creation
gRPC server/client instrumentation
manual spans for important business operations
Kafka publish/consume spans
MySQL spans where appropriate
```

---

## Metrics

Metrics helpers should support:

```text
request counts
request duration
error counts
checkout attempts
checkout failures
Kafka publish failures
dependency latency
```

Keep metric labels low-cardinality.

---

## Logging

Logging helpers should make it easier to include:

```text
correlation_id
trace_id
span_id
service.name
grpc.method
grpc.status_code
```

Start with structured `slog` and add helpers only where they reduce repeated code.

---

## Propagation

Propagation helpers should support:

```text
gRPC metadata propagation
Kafka header propagation
correlation_id propagation
traceparent propagation
```

---

## Shutdown

Telemetry providers/exporters must be flushed during graceful shutdown.

Recommended sequence:

```text
mark health NOT_SERVING
stop accepting new work
finish in-flight requests where possible
flush telemetry
stop gRPC server
close dependencies
exit
```

---

## What This Package Should Not Do

Do not put business observability decisions directly into this package.

Bad:

```text
telemetry package knows checkout business rules
telemetry package knows catalogue schema details
telemetry package records payment-specific attributes automatically everywhere
```

Good:

```text
telemetry package configures common OpenTelemetry plumbing
services add business-specific spans/attributes deliberately
```

---

## Testing Guidance

Recommended tests:

```text
config builds expected resource attributes
disabled telemetry path is safe in local tests
shutdown function flushes providers
trace context can be injected/extracted
correlation ID is preserved where helper supports it
metrics avoid high-cardinality labels
```

---

## Practical Rules

```text
Keep telemetry setup consistent.
Set service.name correctly.
Use OpenTelemetry standards where possible.
Use structured logs.
Propagate trace context through gRPC and Kafka.
Keep attributes safe and low-cardinality.
Do not log secrets.
Flush telemetry during shutdown.
Keep package helpers boring and small.
```

---

## Final Rule

```text
The telemetry package should provide the wiring.
Services should decide what business events are worth observing.
```
