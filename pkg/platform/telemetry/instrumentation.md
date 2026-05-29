# Telemetry Instrumentation Package Guidance

This document describes how `pkg/platform/telemetry` should support instrumentation across **bfstore** services.

It complements:

```text
docs/architecture/instrumentation.md
docs/architecture/observability.md
docs/architecture/opentelemetry.md
docs/architecture/context-propagation.md
docs/architecture/telemetry-signals.md
pkg/platform/telemetry/README.md
pkg/platform/telemetry/propagation.md
pkg/platform/telemetry/signals.md
```

---

## Purpose

The telemetry package should provide shared helpers for:

```text
OpenTelemetry SDK setup
resource attributes
tracer provider setup
meter provider setup
propagator configuration
telemetry shutdown/flush
common instrumentation helpers
```

It should keep telemetry setup consistent without hiding important business instrumentation.

---

## Package Responsibilities

The package may provide:

```text
Config type for telemetry setup
Init function for tracer/meter providers
resource attribute helpers
OTLP exporter setup
propagation setup
shutdown function
trace/log correlation helpers
```

Future file shape:

```text
pkg/platform/telemetry/
├── README.md
├── instrumentation.md
├── tracing.go
├── metrics.go
├── logging.go
├── propagation.go
└── resource.go
```

---

## What Services Still Own

Services should own business instrumentation decisions.

Examples:

```text
order-service decides checkout span stages
payment-service decides payment provider attributes
inventory-service decides stock reservation result attributes
notification-service decides notification delivery metrics
```

The shared package provides the wiring.

Services provide the meaning.

---

## Tracing Helpers

Possible helpers:

```text
create service tracer
start named business span
record error on span
set common bfstore attributes
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

---

## Metrics Helpers

Possible helpers:

```text
create service meter
create request duration histogram
create request/error counters
create Kafka publish/consume counters
create checkout business metrics
```

Metric labels should remain low-cardinality:

```text
service
method
status
environment
failure_reason_class
event_type
```

Do not create helpers that encourage labels like:

```text
order_id
customer_email
basket_id
full_error_message
```

---

## Logging Helpers

Possible helpers:

```text
extract trace_id from context
extract span_id from context
extract correlation_id from context
add trace/correlation attributes to slog records
```

Do not create a logging framework maze.

Start simple with structured `slog`.

---

## Resource Helpers

Every service should set:

```text
service.name
service.namespace
service.version
deployment.environment
```

Example config shape:

```go
type Config struct {
    ServiceName string
    Namespace   string
    Version     string
    Environment string
    OTLPEndpoint string
}
```

---

## Shutdown

Telemetry should be flushed during graceful shutdown.

Expected service shutdown sequence:

```text
mark health NOT_SERVING
stop accepting new work
finish in-flight requests where possible
flush telemetry
stop gRPC server
close dependencies
exit
```

The telemetry package should provide a bounded shutdown function.

---

## Sensitive Data Rules

Helpers must not encourage recording:

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
full SQL with sensitive values
```

Safe helpers should make unsafe telemetry harder to accidentally add.

---

## Testing Guidance

Recommended tests:

```text
resource attributes are set correctly
telemetry init returns shutdown function
shutdown flushes providers
trace/log helpers add trace_id and span_id when available
metric helper labels remain bounded
sensitive values are not added by shared helpers
disabled telemetry mode is safe for unit tests
```

---

## What This Package Should Not Do

Do not put business logic here.

Bad:

```text
telemetry package decides checkout outcome
telemetry package records every order ID as a metric label
telemetry package knows payment retry policy
telemetry package hard-codes a vendor backend
```

Good:

```text
telemetry package configures common OpenTelemetry plumbing
services add deliberate business spans, metrics, and logs
applications decide exporters/backends
```

---

## Practical Rules

```text
Keep helpers small.
Keep instrumentation explicit where business meaning matters.
Use OpenTelemetry semantic conventions.
Do not hard-code a telemetry backend.
Do not leak secrets.
Keep metric labels low-cardinality.
Make trace/log correlation easy.
Flush telemetry on shutdown.
Avoid over-abstraction before services are instrumented.
```

---

## Final Rule

```text
The telemetry package should provide the instrumentation wiring.
Services should decide what meaningful behaviour to observe.
```
