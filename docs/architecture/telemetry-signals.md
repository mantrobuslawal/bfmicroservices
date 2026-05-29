# Telemetry Signals

This document defines the telemetry signal policy for **bfstore**.

OpenTelemetry signals are the different types of observability data emitted by services and infrastructure.

bfstore uses signals to make the platform:

```text
debuggable
measurable
operable
safe to troubleshoot
safe to evolve
```

---

## Purpose

bfstore is a distributed ecommerce platform. A single business flow may cross:

```text
api-gateway
basket-service
order-service
inventory-service
payment-service
shipping-service
Kafka
notification-service
```

No single signal is enough to understand the whole system.

bfstore uses:

```text
traces
metrics
logs
baggage, only when clearly justified
profiles, later for code-level performance analysis
```

---

## Signal Summary

```text
Traces:
  request journeys across services

Metrics:
  numeric measurements over time

Logs:
  timestamped event/error details

Baggage:
  propagated contextual key-value data

Profiles:
  code-level resource usage
```

Practical model:

```text
Metrics alert.
Traces locate.
Logs explain.
Profiles optimise.
Baggage enriches, carefully.
```

---

## Traces

Use traces for distributed request and workflow journeys.

Initial trace coverage:

```text
gRPC server calls
gRPC client calls
catalog-service GetProduct
basket-service AddItem
order-service Checkout
inventory-service ReserveStock
payment-service AuthorisePayment
shipping-service CreateShipment
Kafka publish/consume
key MySQL operations
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

Trace context should propagate through gRPC metadata, Kafka headers, and HTTP headers where applicable.

---

## Metrics

Use metrics for trends, dashboards, alerting, and SLO tracking.

Initial metric areas:

```text
gRPC request count
gRPC request duration
gRPC error count
checkout attempts
checkout failures
payment authorisation duration
inventory reservation failures
Kafka publish failures
Kafka consume failures
MySQL query duration
```

Suggested metric examples:

```text
grpc.server.requests
grpc.server.duration
checkout.attempts
checkout.failures
payment.authorisation.duration
inventory.reservation.failures
kafka.messages.published
kafka.publish.failures
mysql.query.duration
```

Use histograms for durations and counters for totals.

---

## Metric Cardinality Rules

Keep metric labels low-cardinality.

Allowed-style labels:

```text
service
method
status_code
environment
event_type
topic
```

Avoid labels like:

```text
order_id
customer_id
customer_email
basket_id
raw_product_name
full_error_message
full_request_path_with_ids
```

Practical rule:

```text
Metrics should group behaviour, not identify individual users or orders.
```

---

## Logs

Logs should be structured and correlated with traces.

Recommended fields:

```text
timestamp
level
service
environment
message
correlation_id
trace_id
span_id
grpc.method
grpc.status_code
event_type
error
```

Add business identifiers only where safe and useful:

```text
order_id
product_id
basket_id
```

Never log:

```text
raw JWTs
passwords
payment card data
CVV
API keys
full shipping address
personal data
full basket JSON
```

---

## Baggage Policy

Default policy:

```text
avoid baggage unless clearly needed
never use baggage for secrets
never use baggage for personal data
strip untrusted baggage at public boundaries
```

Potential future low-risk baggage examples:

```text
tenant_id=bfstore
region=uk
```

Do not use baggage for customer email, shipping address, JWTs, payment tokens, basket JSON, or card details.

---

## Profiles

Profiles are a later-stage signal for code-level performance analysis.

Use profiles when metrics and traces show a problem but the expensive code path is unclear.

Useful cases:

```text
catalog-service CPU hotspots
product attribute mapping inefficiency
payment-service retry/signing overhead
Kafka consumer memory allocation issues
high CPU during Protobuf serialisation
```

Profiles should support optimisation decisions with evidence.

---

## Events

Do not confuse OpenTelemetry telemetry events with Kafka business events.

```text
Kafka business event:
  bfstore.order.events.v1.OrderCreated

Telemetry event:
  payment_retry_started
  stock_reservation_conflict_detected
```

Business events describe domain facts. Telemetry events describe notable moments during system execution.

---

## Sensitive Data Rules

Sensitive data must stay out of every signal:

```text
logs
metrics
traces
span attributes
baggage
profiles where applicable
```

Do not record raw JWTs, API keys, passwords, card numbers, CVV, payment tokens, full addresses, customer email as a telemetry attribute, or full basket JSON.

Practical rule:

```text
Telemetry should explain system behaviour, not leak customer or security data.
```

---

## Example Investigation Workflow

Problem:

```text
Users report checkout is slow.
```

Metrics show:

```text
checkout p95 latency rose from 900ms to 4.5s
payment authorisation p95 rose from 500ms to 4s
```

Traces show:

```text
OrderService/Checkout spends most time waiting on PaymentService/AuthorisePayment
```

Logs show:

```text
payment-service provider simulation timed out after retries
```

Profiles may show:

```text
payment-service CPU is heavily used in retry request signing
```

Outcome:

```text
tune payment retry/timeout policy
optimise expensive code path
add alert for payment provider latency
```

---

## Practical Rules

```text
Use traces for request journeys.
Use metrics for trends, dashboards, alerts, and SLOs.
Use logs for detailed event/error context.
Use baggage sparingly and never for secrets.
Use profiles when code-level performance evidence is needed.
Correlate logs with trace_id and correlation_id.
Keep metric labels low-cardinality.
Keep sensitive data out of every signal.
Use consistent service names.
Start with catalogue and checkout flows.
```

---

## Final Rule

```text
Signals are the different camera angles on bfstore.
```
