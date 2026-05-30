# gRPC Metrics

This document describes how `pkg/platform/telemetry` should support gRPC OpenTelemetry metrics in **bfstore**.

It complements:

```text
docs/architecture/grpc-opentelemetry-metrics.md
docs/architecture/instrumentation.md
docs/architecture/opentelemetry-components.md
docs/architecture/semantic-conventions.md
```

---

## Purpose

`pkg/platform/telemetry` should help bfstore services configure OpenTelemetry metrics consistently.

This includes:

```text
MeterProvider setup
OTLP metrics exporter setup
gRPC server metrics wiring
gRPC client metrics wiring
metrics views where needed
safe metric label guidance
```

---

## Responsibilities

This package may provide:

```text
NewMeterProvider
metrics exporter setup
resource attachment
shutdown handling
environment-specific config
helper functions for gRPC metrics wiring
```

It should not hide normal OpenTelemetry concepts behind unnecessary abstraction.

---

## Metrics Flow

```text
bfstore service
  -> gRPC metrics instrumentation
  -> MeterProvider
  -> OTLP exporter
  -> Collector
  -> metrics backend
```

---

## Required Metric Categories

bfstore should support:

```text
server call duration
client call duration
client attempt duration
message size metrics
retry metrics later
```

---

## Client Metrics

Client metrics should be enabled for outgoing dependencies such as:

```text
order-service -> inventory-service
order-service -> payment-service
order-service -> shipping-service
basket-service -> catalog-service
api-gateway -> catalog-service
```

Important:

```text
grpc.client.call.duration
grpc.client.attempt.duration
```

Use client metrics to understand dependency latency from the caller’s point of view.

---

## Server Metrics

Server metrics should be enabled for all gRPC services.

Important:

```text
grpc.server.call.duration
grpc.server.call.started
```

Use server metrics to understand service-side latency, request volume, and status codes.

---

## Views and Buckets

Metrics views may be used later to tune histogram buckets or reduce cardinality.

Do not tune prematurely.

Rule:

```text
Use default gRPC metrics first.
Tune views after observing real data.
```

---

## Attribute Safety

Allowed-style attributes:

```text
grpc.method
grpc.status
grpc.target
service.name
deployment.environment.name
```

Do not add high-cardinality custom labels:

```text
order_id
customer_email
basket_id
product_slug
full error message
```

Rule:

```text
Metric labels must remain low-cardinality.
```

---

## Business Metrics Are Separate

Do not use gRPC metrics as a replacement for business metrics.

Add business metrics such as:

```text
checkout.completed_total
checkout.failed_total
payment.authorised_total
stock.reserved_total
notification.sent_total
```

gRPC metrics describe RPC transport. Business metrics describe outcomes.

---

## Testing Guidance

Unit/integration tests should check:

```text
MeterProvider is configured
OTLP metric exporter is configured
gRPC server emits metrics
gRPC client emits metrics
service.name resource is present
grpc.method attribute is present
grpc.status attribute is present
high-cardinality labels are not added
```

---

## Practical Rules

```text
Wire metrics once, consistently.
Use OpenTelemetry MeterProvider.
Export metrics through OTLP to the Collector.
Prefer built-in gRPC metrics over custom RPC counters.
Keep labels low-cardinality.
Add business metrics separately.
Do not tune views before observing real data.
```

---

## Final Rule

```text
pkg/platform/telemetry should make gRPC metrics boring, consistent, and hard to misuse.
```
