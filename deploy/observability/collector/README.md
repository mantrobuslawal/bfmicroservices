# OpenTelemetry Collector

This directory contains bfstore OpenTelemetry Collector documentation and future configuration.

Recommended path:

```text
deploy/observability/collector/
```

Future files:

```text
deploy/observability/collector/
├── README.md
├── otel-collector.local.yaml
└── otel-collector.k8s.yaml
```

---

## Purpose

The OpenTelemetry Collector is the central telemetry pipeline for bfstore.

It receives telemetry from services, processes it, and exports it to observability backends.

Recommended service flow:

```text
bfstore service
        |
        | OTLP
        v
OpenTelemetry Collector
        |
        +--> traces backend
        +--> metrics backend
        +--> logs backend
```

---

## Local Development Goal

Local development should start simple.

Initial Collector goals:

```text
receive OTLP over gRPC and HTTP
batch telemetry
export debug output
optionally expose metrics for Prometheus
optionally forward traces to Jaeger or Tempo
```

This helps verify that services are emitting telemetry correctly before adding a full backend stack.

---

## Basic Collector Shape

Example conceptual config:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:

exporters:
  debug:

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]

    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]

    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]
```

---

## Service Configuration

bfstore services should export to the Collector using OTLP.

Example environment variables:

```text
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
OTEL_SERVICE_NAME=bfstore-order-service
OTEL_RESOURCE_ATTRIBUTES=service.namespace=bfstore,deployment.environment=local
```

---

## Backend Routing

Later, the Collector may route telemetry to:

```text
Tempo or Jaeger for traces
Prometheus-compatible backend for metrics
Loki or log backend for logs
```

Example conceptual routing:

```text
traces  -> Tempo/Jaeger
metrics -> Prometheus
logs    -> Loki/stdout
```

---

## Processors

Start with:

```text
batch
```

Later consider:

```text
memory_limiter
attributes
resource
filter
tail_sampling
```

Do not add complex processing before the basic pipeline works.

---

## Security and Data Hygiene

The Collector can help filter or drop unsafe attributes, but services should not emit unsafe telemetry in the first place.

Do not emit:

```text
raw JWTs
passwords
API keys
card numbers
CVV
full shipping addresses
customer email
full basket JSON
raw Kafka payloads
```

Practical rule:

```text
Fix unsafe telemetry at the source.
Use Collector filtering as a safety net, not the main defence.
```

---

## Kubernetes Later

In Kubernetes, the Collector may run as:

```text
Deployment
DaemonSet
sidecar, less likely for bfstore initially
```

The OpenTelemetry Operator may later manage Collector resources and auto-instrumentation.

This is a later-stage improvement, not a day-one requirement.

---

## Testing Guidance

Recommended checks:

```text
Collector starts successfully
OTLP gRPC receiver is reachable
OTLP HTTP receiver is reachable
debug exporter shows traces from a local service
service.name appears correctly
deployment.environment appears correctly
trace context connects across services
```

---

## Practical Rules

```text
Services export OTLP to the Collector.
The Collector routes telemetry to backends.
Start with debug output locally.
Add backends gradually.
Keep Collector config version-controlled.
Do not rely on Collector filtering to fix unsafe instrumentation.
Keep local config simple before Kubernetes config.
```

---

## Final Rule

```text
The Collector is bfstore’s observability sorting office.
```
