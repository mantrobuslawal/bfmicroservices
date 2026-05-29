# Telemetry Sampling

This document describes how `pkg/platform/telemetry` should support trace sampling for **bfstore** services.

It complements:

```text
docs/architecture/sampling.md
docs/architecture/instrumentation.md
docs/architecture/opentelemetry-components.md
docs/architecture/opentelemetry-resources.md
docs/architecture/context-propagation.md
```

---

## Purpose

`pkg/platform/telemetry` may provide SDK sampling configuration for bfstore services.

Sampling should be explicit, environment-aware, and easy to reason about.

---

## Responsibilities

This package may configure:

```text
sample-all behaviour for local development
probability-based head sampling
environment-specific sampling defaults
configuration parsing
safe fallback behaviour
```

It should not own complex production tail-sampling policy. That belongs in Collector configuration and architecture docs.

---

## Recommended Defaults

Local:

```text
sample all traces
```

Dev/test initially:

```text
sample all traces
```

Load testing:

```text
allow configurable probability sampling
```

Production-style later:

```text
use SDK sampling only for simple volume control
prefer Collector policy for complex tail sampling
```

---

## Config Shape

Possible Go config:

```go
type SamplingConfig struct {
    Mode        string  // "always_on", "always_off", "trace_id_ratio"
    Ratio       float64 // 0.0 to 1.0
    Environment string
}
```

Example local config:

```go
SamplingConfig{
    Mode: "always_on",
    Ratio: 1.0,
    Environment: "local",
}
```

Example load-test config:

```go
SamplingConfig{
    Mode: "trace_id_ratio",
    Ratio: 0.10,
    Environment: "load-test",
}
```

---

## Environment Variables

Possible environment variables:

```text
BFSTORE_OTEL_SAMPLING_MODE=always_on
BFSTORE_OTEL_SAMPLING_RATIO=1.0
```

For OpenTelemetry-native configuration, document any supported `OTEL_*` variables used by the Go SDK.

---

## Sampler Helper

Possible helper shape:

```go
func NewSampler(cfg SamplingConfig) trace.Sampler {
    switch cfg.Mode {
    case "always_on":
        return trace.AlwaysSample()
    case "always_off":
        return trace.NeverSample()
    case "trace_id_ratio":
        return trace.TraceIDRatioBased(cfg.Ratio)
    default:
        return trace.AlwaysSample()
    }
}
```

For local development, default to sampling all traces.

Do not silently apply aggressive sampling in local or early dev.

---

## What Belongs in the Collector

Complex policies should live in Collector config, not service code.

Collector-owned examples:

```text
keep all error traces
keep all slow checkout traces
tail sampling
latency-based sampling
attribute-based sampling
sampling based on full trace contents
```

Service SDK-owned examples:

```text
always sample local traces
simple head sampling ratio
parent-based sampling
```

---

## Sampling and Propagation

Sampling helpers should work with context propagation.

Expected propagation boundaries:

```text
gRPC metadata
Kafka headers
traceparent
```

If a parent trace is sampled, child services should honour the sampled decision where appropriate.

---

## Metrics Are Not Sampled Away

The telemetry package must not treat trace sampling as a replacement for metrics.

Metrics should continue to record aggregate behaviour:

```text
checkout.failed_total
payment.timeout_total
kafka.publish.failed_total
notification.delivery.failed_total
```

---

## Testing Guidance

Unit tests:

```text
local config returns always-on sampler
ratio config applies expected ratio
invalid config falls back safely
ratio below 0 or above 1 is rejected or clamped deliberately
production config is explicit
```

Integration checks:

```text
local service emits all test traces
ratio sampling reduces high-volume traces
sampled trace context propagates across gRPC
sampled trace context propagates through Kafka headers
```

---

## What This Package Should Not Do

Do not hide complex sampling decisions in service code.

Avoid:

```text
hard-coded production sampling percentages
tail sampling logic inside services
business-flow retention policy hidden in SDK config
silent aggressive sampling defaults
```

Prefer:

```text
simple local/dev defaults in code
explicit config
complex policies in Collector docs/config
architecture docs explaining decisions
```

---

## Practical Rules

```text
Default local development to sample all traces.
Keep SDK sampling simple.
Make sampling config explicit.
Use Collector sampling for complex policies.
Do not use trace sampling as a replacement for metrics.
Propagate sampling context correctly.
Do not silently drop useful traces in early environments.
```

---

## Final Rule

```text
pkg/platform/telemetry should make simple sampling safe, explicit, and boring.
```
