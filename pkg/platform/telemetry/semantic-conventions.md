# Telemetry Semantic Conventions

This document describes how `pkg/platform/telemetry` should support semantic conventions for **bfstore**.

It complements:

```text
docs/architecture/semantic-conventions.md
docs/architecture/observability.md
docs/architecture/opentelemetry.md
docs/architecture/opentelemetry-components.md
docs/architecture/instrumentation.md
docs/architecture/telemetry-signals.md
docs/architecture/context-propagation.md
```

---

## Purpose

`pkg/platform/telemetry` may provide helpers or constants for bfstore-specific telemetry attributes.

It should encourage standard OpenTelemetry naming without wrapping everything unnecessarily.

---

## General Rule

Use OpenTelemetry semantic conventions for technical facts.

Use bfstore-specific attributes for domain facts.

```text
Technical:
  rpc.method
  db.system
  messaging.system
  service.name

Business:
  bfstore.checkout.stage
  bfstore.payment.provider
  bfstore.stock.reservation_result
```

---

## bfstore-specific Attribute Constants

A future Go package may define constants for bfstore-specific attributes.

Possible location:

```text
pkg/platform/telemetry/attrs
```

Example:

```go
package attrs

const (
    BFStoreCheckoutStage          = "bfstore.checkout.stage"
    BFStorePaymentProvider        = "bfstore.payment.provider"
    BFStoreStockReservationResult = "bfstore.stock.reservation_result"
    BFStoreNotificationChannel    = "bfstore.notification.channel"
    BFStoreCatalogAttributeCount  = "bfstore.catalog.attribute_count"
)
```

Avoid creating duplicate constants for standard OpenTelemetry attributes when official packages already provide them.

---

## Resource Helpers

The telemetry package may provide helpers for common resource attributes:

```text
service.name
service.namespace
service.version
deployment.environment
```

Example config:

```go
type Config struct {
    ServiceName string
    Namespace string
    Version string
    Environment string
}
```

---

## gRPC Helpers

gRPC instrumentation should prefer OpenTelemetry instrumentation libraries.

Expected attributes:

```text
rpc.system = grpc
rpc.service
rpc.method
grpc.status_code
```

bfstore-specific additions should be deliberate:

```text
bfstore.checkout.stage
bfstore.payment.provider
```

---

## Kafka Helpers

Kafka helper code may standardise these operational attributes:

```text
messaging.system = kafka
messaging.destination.name
messaging.operation.name
event.type
```

Kafka header helpers may also standardise:

```text
traceparent
correlation_id
event_type
event_version
content_type
```

Keep payloads and telemetry separate.

---

## Database Helpers

Database helper code may standardise safe database attributes:

```text
db.system = mysql
db.name
db.operation
db.collection.name
```

Do not include full SQL with sensitive values.

---

## Logging Helpers

Log enrichment helpers may add:

```text
trace_id
span_id
correlation_id
service.name
deployment.environment
rpc.service
rpc.method
grpc.status_code
event.type
error.type
```

Logging helpers must not add secrets or personal data.

---

## Cardinality Rules

Helpers should avoid encouraging high-cardinality labels.

Do not create metric helpers that label by:

```text
order_id
customer_id
customer_email
basket_id
raw_product_name
full_error_message
raw URL path with IDs
```

Safe-style metric labels:

```text
service.name
rpc.method
grpc.status_code
deployment.environment
event.type
bfstore.checkout.stage
```

---

## Sensitive Data Rules

The telemetry package must not provide helpers that record:

```text
raw JWTs
passwords
API keys
payment card numbers
CVV
full shipping address
customer email
full basket JSON
raw Kafka payloads
full SQL with sensitive values
```

---

## Testing Guidance

Recommended tests:

```text
bfstore-specific attribute constants match documented names
resource helpers set service.name consistently
log enrichment adds trace_id/span_id when available
metric helpers avoid high-cardinality labels
Kafka helpers use documented header names
no helper records unsafe sensitive fields
```

---

## What This Package Should Not Do

Do not create a wrapper universe around OpenTelemetry.

Bad:

```text
pkg/platform/telemetry redefines every standard OTel attribute
services cannot see normal OTel concepts
helpers hide business meaning
metric helpers accept arbitrary high-cardinality labels
```

Good:

```text
use OTel semantic conventions directly where practical
provide bfstore-specific constants for business attributes
make safe naming easy
keep instrumentation explicit
```

---

## Practical Rules

```text
Prefer OpenTelemetry names for standard telemetry.
Use bfstore.* only for business meaning.
Do not duplicate official constants unnecessarily.
Keep helper APIs small.
Protect metric cardinality.
Protect sensitive data.
Make trace/log correlation easy.
Keep service instrumentation readable.
```

---

## Final Rule

```text
pkg/platform/telemetry should help bfstore speak clean telemetry grammar, not invent a new language.
```
