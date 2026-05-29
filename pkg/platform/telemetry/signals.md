# Telemetry Signals Package Guidance

This document describes how `pkg/platform/telemetry` should support OpenTelemetry signals for **bfstore**.

It complements:

```text
docs/architecture/telemetry-signals.md
docs/architecture/observability.md
docs/architecture/opentelemetry.md
docs/architecture/context-propagation.md
pkg/platform/telemetry/README.md
pkg/platform/telemetry/propagation.md
```

---

## Purpose

The telemetry package should provide shared helpers and conventions for:

```text
traces
metrics
structured log correlation
baggage policy, if ever enabled
profiles, later
```

It should keep common telemetry setup consistent across bfstore services.

---

## Traces

The package may provide helpers for:

```text
tracer provider setup
service tracer creation
span naming helpers
common span attributes
gRPC instrumentation wiring
Kafka publish/consume tracing helpers
```

Recommended span names:

```text
CatalogService/GetProduct
BasketService/AddItem
OrderService/Checkout
InventoryService/ReserveStock
PaymentService/AuthorisePayment
ShippingService/CreateShipment
Kafka publish OrderCreated
Kafka consume OrderCreated
```

Do not hide important business instrumentation behind magical helpers. Services should still choose meaningful business spans deliberately.

---

## Metrics

The package may provide helpers for:

```text
meter provider setup
request counters
duration histograms
error counters
Kafka publish/consume metrics
dependency latency metrics
```

Recommended metric categories:

```text
request count
request duration
error count
checkout success/failure
Kafka publish failures
database latency
resource usage
```

Metric labels must remain low-cardinality.

Allowed-style labels:

```text
service
method
status_code
environment
event_type
topic
```

Avoid:

```text
order_id
customer_email
basket_id
raw_product_name
full_error_message
```

---

## Logs

The package may provide helpers to enrich logs with:

```text
trace_id
span_id
correlation_id
service.name
environment
grpc.method
grpc.status_code
event_type
```

Start with structured `slog` and only add helpers where they reduce repeated code.

Do not create a logging framework maze.

Keep it boring where production matters.

---

## Baggage

Default package behaviour:

```text
do not enable broad baggage propagation
do not create baggage helpers until a clear use case exists
strip or ignore untrusted baggage at public boundaries
```

Never allow baggage to contain raw JWTs, API keys, passwords, payment tokens, card details, CVV, customer email, shipping address, basket JSON, or personal data.

---

## Profiles

Profiles are a later-stage concern.

The package may eventually document how to enable profiling for:

```text
local development
performance investigations
load testing
staging diagnostics
```

Profile usage should be intentional and evidence-driven.

Useful cases:

```text
CPU hotspots
memory allocation issues
slow catalogue mapping
expensive checkout code
Kafka consumer performance
```

---

## Sensitive Data Rules

The package should make the safe path easy.

Do not put sensitive data into:

```text
span attributes
metric labels
log fields
baggage
profile labels
```

Avoid raw JWTs, passwords, card data, CVV, full shipping address, customer email, full request payloads, and full Kafka payloads.

---

## Testing Guidance

Recommended tests:

```text
common resource attributes are set correctly
metric labels remain bounded
log enrichment includes trace_id/span_id when available
correlation ID is preserved
baggage is disabled or filtered by default
helpers do not record sensitive values
```

---

## What This Package Should Not Do

Do not put business logic here.

Bad:

```text
telemetry package decides checkout result
telemetry package knows catalogue schema details
telemetry package embeds payment retry rules
telemetry package records every order ID as a metric label
```

Good:

```text
telemetry package provides common instrumentation plumbing
services add deliberate business spans, metrics, and logs
```

---

## Practical Rules

```text
Keep telemetry helpers small.
Keep service instrumentation explicit.
Keep metric labels low-cardinality.
Keep sensitive data out.
Use consistent service names.
Prefer OpenTelemetry conventions where practical.
Make trace/log correlation easy.
Do not over-abstract before the first services are instrumented.
```

---

## Final Rule

```text
The telemetry package should provide the signal wiring.
Services should decide what meaningful behaviour to observe.
```
