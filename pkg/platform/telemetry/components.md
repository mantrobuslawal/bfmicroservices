# Telemetry Components

This document explains which OpenTelemetry components `pkg/platform/telemetry` may own or configure for **bfstore**.

It complements:

```text
docs/architecture/opentelemetry-components.md
docs/architecture/opentelemetry.md
docs/architecture/observability.md
docs/architecture/context-propagation.md
docs/architecture/instrumentation.md
docs/architecture/telemetry-signals.md
```

---

## Purpose

`pkg/platform/telemetry` should provide shared telemetry setup for bfstore Go services.

Its job is to make common OpenTelemetry configuration consistent without hiding service-specific business instrumentation.

---

## Components Owned by This Package

This package may configure:

```text
resource attributes
tracer provider
meter provider
OTLP exporters
propagators
samplers
shutdown/flush behaviour
```

It may also expose helpers for:

```text
trace/log correlation
common metric instruments
Kafka propagation carriers
service resource setup
```

---

## Components Not Owned by This Package

This package should not own:

```text
business span design
checkout stage decisions
payment retry policy
catalogue domain attributes
backend deployment
Collector deployment lifecycle
Kubernetes Operator lifecycle
```

Services own business meaning.

Deployment code owns infrastructure lifecycle.

---

## API vs SDK Boundary

Shared lower-level packages should depend only on the OpenTelemetry API where possible.

Service startup code may configure the SDK.

Practical model:

```text
pkg/platform/grpc/interceptors:
  uses API/instrumentation hooks

service main.go:
  initialises SDK via pkg/platform/telemetry

deploy/observability/collector:
  configures Collector pipeline
```

---

## Resource Setup

The package should help services set:

```text
service.name
service.namespace
service.version
deployment.environment
```

Example config shape:

```go
type Config struct {
    ServiceName  string
    Namespace    string
    Version      string
    Environment  string
    OTLPEndpoint string
}
```

---

## Exporter Setup

Preferred exporter:

```text
OTLP
```

Preferred path:

```text
service -> OTLP exporter -> OpenTelemetry Collector
```

The package should avoid hard-coding backend-specific exporters unless there is a clear reason.

---

## Propagator Setup

Default propagation should support:

```text
traceparent
tracestate, if used
```

bfstore-specific helpers may also handle:

```text
x-correlation-id for gRPC
correlation_id for Kafka
```

---

## Sampling

Local development default:

```text
sample all traces
```

Later production configuration may allow:

```text
ratio-based sampling
parent-based sampling
Collector tail sampling
error-priority policies
```

Sampling policy should be explicit and configurable.

---

## Shutdown

Telemetry must be flushed during graceful shutdown.

Expected service shutdown order:

```text
mark health NOT_SERVING
stop accepting new work
finish in-flight requests where possible
flush telemetry providers/exporters
stop gRPC server
close dependencies
exit
```

The package should return a shutdown function from telemetry initialisation.

---

## Testing Guidance

Recommended tests:

```text
resource attributes are created correctly
OTLP endpoint configuration is respected
propagator is configured
shutdown function flushes providers
disabled telemetry mode is safe for unit tests
sampling config is applied
no backend vendor is hard-coded unexpectedly
```

---

## Practical Rules

```text
Keep telemetry setup centralised.
Keep business instrumentation explicit in services.
Use OTLP exporters.
Send telemetry to the Collector.
Set resource attributes consistently.
Do not hard-code observability vendors.
Make shutdown safe and bounded.
Make local development easy.
```

---

## Final Rule

```text
pkg/platform/telemetry owns the OpenTelemetry wiring.
Services own the business meaning.
```
