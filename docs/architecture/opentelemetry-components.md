# OpenTelemetry Components

This document defines how **bfstore** should understand and use OpenTelemetry components.

OpenTelemetry is made of several components that work together to produce, process, export, and query telemetry.

---

## Purpose

This document explains the role of each OpenTelemetry component in the bfstore platform.

It exists so the project has clear boundaries between:

```text
application code
OpenTelemetry API
OpenTelemetry SDK
instrumentation libraries
exporters
Collector
observability backends
Kubernetes automation
```

---

## bfstore Telemetry Flow

Recommended flow:

```text
bfstore service code
        |
        | uses OpenTelemetry API
        v
OpenTelemetry SDK
        |
        | OTLP exporter
        v
OpenTelemetry Collector
        |
        v
Observability backends
```

Example backend mapping:

```text
traces  -> Tempo or Jaeger
metrics -> Prometheus/Grafana
logs    -> Loki or structured log backend
```

---

## Specification

The OpenTelemetry specification is the shared rulebook.

It defines concepts such as:

```text
trace
span
metric
log
resource
attribute
propagation
sampler
exporter
OTLP
semantic conventions
```

bfstore should follow OpenTelemetry conventions wherever practical.

---

## API

The API is what bfstore code calls to create telemetry.

Shared packages may use OpenTelemetry API types:

```text
pkg/platform/grpc/interceptors
pkg/platform/telemetry
pkg/platform/kafka
pkg/platform/mysql
```

Shared packages must not hard-code a telemetry backend.

Practical rule:

```text
Shared packages may use the OpenTelemetry API.
Applications configure the SDK, exporters, Collector, and backend.
```

---

## SDK

The SDK is the implementation that processes telemetry.

Service startup/bootstrap code should configure the SDK.

SDK responsibilities include:

```text
tracer provider setup
meter provider setup
resource attributes
span processors
metric readers
sampling
OTLP exporters
shutdown/flush
```

This should be centralised in:

```text
pkg/platform/telemetry
```

---

## Instrumentation Libraries

Use instrumentation libraries for common plumbing:

```text
gRPC server requests
gRPC client calls
HTTP gateway requests
MySQL calls
Kafka producer/consumer operations
Go runtime metrics
```

Use manual instrumentation for business meaning:

```text
checkout stages
stock reservation result
payment provider behaviour
Kafka business event publish
notification delivery result
```

---

## Exporters

Preferred bfstore pattern:

```text
service -> OTLP exporter -> OpenTelemetry Collector
```

Avoid direct service-to-many-backend exports.

Practical rule:

```text
Services export once.
The Collector routes many.
```

---

## Collector

The OpenTelemetry Collector is the central telemetry pipeline.

Collector responsibilities:

```text
receive OTLP telemetry from services
batch telemetry
filter/process telemetry
route traces to trace backend
route metrics to metrics backend
route logs to log backend
support local debugging
```

Local deployment should start with a simple Collector config and debug exporter.

Production-like deployment can add backend routing and processors.

---

## Resource Attributes

All services must set resource attributes.

Required starting attributes:

```text
service.name
service.namespace
service.version
deployment.environment
```

Example:

```text
service.name = bfstore-order-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment = local
```

Kubernetes later may add:

```text
k8s.namespace.name
k8s.pod.name
k8s.deployment.name
k8s.cluster.name
```

---

## Propagators

Propagators carry context across service boundaries.

bfstore propagation policy:

```text
gRPC metadata:
  traceparent
  x-correlation-id

Kafka headers:
  traceparent
  correlation_id
```

Propagation connects distributed traces across:

```text
api-gateway
order-service
payment-service
Kafka
notification-service
```

---

## Samplers

Sampling controls trace volume.

Local development:

```text
sample all traces
```

Production later:

```text
sample high-volume successful traffic
retain error traces where possible
consider Collector tail sampling later
```

Sampling policy should be deliberate and documented before production use.

---

## Zero-code Instrumentation

Zero-code instrumentation may be considered later.

For early bfstore phases, prefer explicit code-based setup because it demonstrates understanding of:

```text
service.name
resource attributes
exporters
propagation
gRPC instrumentation
Collector flow
```

---

## Kubernetes Operator

The OpenTelemetry Operator may later manage:

```text
Collector deployment
auto-instrumentation
Kubernetes telemetry configuration
```

It is not a day-one requirement.

Use it after the basic telemetry pipeline is understood.

---

## FaaS Assets

Function-as-a-Service assets are not currently part of the core bfstore direction.

They may become relevant if bfstore later adds serverless features such as:

```text
image thumbnail generation
scheduled catalogue import
invoice PDF generation
```

---

## Recommended Rollout

### Phase 1: Local foundation

```text
pkg/platform/telemetry initialises SDK
services export OTLP to Collector
Collector debug exports telemetry
resource attributes set correctly
propagation configured
```

### Phase 2: gRPC and catalogue

```text
gRPC instrumentation
CatalogService/GetProduct traces
catalogue MySQL spans
basic request metrics
trace/log correlation
```

### Phase 3: Checkout

```text
OrderService/Checkout spans
inventory/payment/shipping spans
business metrics
Kafka publish OrderCreated span
```

### Phase 4: Async consumers

```text
NotificationService Kafka consume spans
trace context extracted from Kafka headers
notification delivery metrics/logs
```

### Phase 5: Collector pipeline

```text
traces routed to Tempo/Jaeger
metrics exposed for Prometheus
logs routed to Loki/stdout
batch/filter processors added
```

### Phase 6: Kubernetes

```text
Collector deployed in Kubernetes
resource detectors pick up Kubernetes attributes
Operator considered
sampling policy refined
```

---

## Practical Rules

```text
Use the OpenTelemetry API in shared instrumentation code.
Use the Go SDK in service startup/bootstrap code.
Prefer OTLP exporters from services.
Send telemetry to the Collector.
Set resource attributes consistently.
Use instrumentation libraries for plumbing.
Use manual instrumentation for business meaning.
Use propagators for gRPC metadata and Kafka headers.
Sample deliberately.
Keep observability backends replaceable.
Do not start with the Kubernetes Operator before understanding the pipeline.
```

---

## Final Rule

```text
OpenTelemetry components are the observability toolkit.
bfstore should assemble them deliberately, not accidentally.
```
