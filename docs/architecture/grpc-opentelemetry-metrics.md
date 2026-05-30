# gRPC OpenTelemetry Metrics

This document defines how **bfstore** uses gRPC OpenTelemetry metrics as the baseline for RPC observability.

bfstore services communicate using gRPC. gRPC OpenTelemetry metrics provide standard visibility into RPC latency, request volume, status codes, retries, attempts, and message sizes.

---

## Purpose

This document defines:

```text
why bfstore uses gRPC metrics
which metrics matter first
client-side vs server-side metrics
per-call vs per-attempt metrics
error-rate and latency guidance
message-size guidance
dashboard expectations
relationship to business metrics
```

---

## Core Rule

```text
Use built-in gRPC OpenTelemetry metrics as the default RPC metrics baseline.
```

Do not create custom RPC metrics until the standard gRPC metrics have been reviewed.

---

## Metrics Flow

```text
bfstore service
  -> gRPC library emits metrics
  -> OpenTelemetry MeterProvider
  -> OTLP exporter
  -> OpenTelemetry Collector
  -> metrics backend
  -> Grafana dashboards / alerts
```

`pkg/platform/telemetry` should configure the OpenTelemetry `MeterProvider`.

gRPC server and client setup should wire metrics instrumentation to that provider.

---

## Client-side Metrics

Client-side metrics describe what the caller experienced.

Important metric:

```text
grpc.client.call.duration
```

Use it to answer:

```text
How long did the dependency call take from the caller’s point of view?
Are deadlines being reached?
Are retries or connection behaviour increasing latency?
```

Rule:

```text
Client duration tells you what the caller experienced.
```

---

## Server-side Metrics

Server-side metrics describe what the receiving service did.

Important metric:

```text
grpc.server.call.duration
```

Use it to answer:

```text
How long did the service spend handling the RPC?
Which methods are slow?
Which methods return non-OK status codes?
```

Rule:

```text
Server duration tells you what the receiver did.
```

Compare client-side and server-side metrics when debugging latency.

---

## Per-call vs Per-attempt

A call is the logical RPC from the application perspective.

An attempt is one actual try made by the gRPC client underneath that call.

```text
call:
  one logical RPC

attempt:
  one network/library try
```

One call can have multiple attempts when retries are enabled.

Important attempt metrics:

```text
grpc.client.attempt.started
grpc.client.attempt.duration
grpc.client.attempt.sent_total_compressed_message_size
grpc.client.attempt.rcvd_total_compressed_message_size
```

Rule:

```text
Per-call metrics show the customer-facing story.
Per-attempt metrics show the mechanical truth underneath.
```

---

## Retry Metrics

When retries are enabled, watch:

```text
grpc.client.call.retries
grpc.client.call.transparent_retries
grpc.client.call.hedges
grpc.client.call.retry_delay
```

bfstore should monitor retries for:

```text
OrderService -> PaymentService
OrderService -> InventoryService
OrderService -> ShippingService
BasketService -> CatalogService
```

Rule:

```text
Retries are not free. Measure them.
```

---

## Message Size Metrics

Watch compressed sent and received message sizes.

Important paths:

```text
CatalogService/ListProducts
CatalogService/GetProduct
BasketService/GetBasket
OrderService/Checkout
```

Potential issue:

```text
ListProducts returns too much product detail.
```

Better API shape:

```text
ListProducts returns summary cards.
GetProduct returns detailed product view.
```

Rule:

```text
If gRPC message size grows, revisit API shape before blaming the network.
```

---

## Attributes and Cardinality

Important labels:

```text
grpc.method
grpc.status
grpc.target
```

Good:

```text
grpc.method = bfstore.order.v1.OrderService/Checkout
grpc.status = DEADLINE_EXCEEDED
grpc.target = dns:///payment-service:50051
```

Avoid high-cardinality labels:

```text
order_id
customer_email
basket_id
product_slug
full error message
```

Rule:

```text
Metrics need low-cardinality labels.
Individual stories belong in traces and logs.
```

---

## Throughput

Throughput can be derived from duration histogram counts.

Use:

```text
grpc.server.call.duration count
```

to calculate:

```text
Checkout RPCs per second
GetProduct RPCs per second
AuthorisePayment RPCs per minute
```

Rule:

```text
The duration histogram gives you both latency and request volume.
```

---

## Error Rate

Error rate can be calculated by filtering where:

```text
grpc.status != OK
```

Example:

```text
PaymentService AuthorisePayment error rate =
  count grpc.server.call.duration where grpc.status != OK
  /
  count all grpc.server.call.duration
```

Rule:

```text
Error rate should be a metric, not a log search.
```

---

## Dashboard Expectations

Build dashboards for:

```text
gRPC request volume
gRPC server latency p50/p95/p99
gRPC client dependency latency p95
gRPC error rate
gRPC message size
gRPC retry behaviour later
```

Group by:

```text
service.name
grpc.method
grpc.status
grpc.target
deployment.environment.name
```

---

## Service-specific Guidance

### catalog-service

Watch:

```text
GetProduct latency
ListProducts latency
response message size
grpc.status != OK
```

### basket-service

Watch:

```text
AddItem latency
GetBasket latency
RemoveItem error rate
message size for large baskets
```

### order-service

Watch:

```text
Checkout duration
Checkout status codes
client call duration to inventory/payment/shipping
retry delay
```

### payment-service

Watch:

```text
AuthorisePayment server duration
DEADLINE_EXCEEDED
UNAVAILABLE
client retries from order-service
```

### notification-service

Watch:

```text
SendOrderConfirmation duration
status codes
message sizes
```

---

## gRPC Metrics vs Business Metrics

gRPC metrics explain transport behaviour.

Business metrics explain outcomes.

gRPC metrics:

```text
latency
status
message size
attempts
retries
targets
```

Business metrics:

```text
checkout.completed_total
checkout.failed_total
payment.authorised_total
stock.reserved_total
notification.sent_total
```

Rule:

```text
gRPC metrics explain the transport.
Business metrics explain the outcome.
```

Do not confuse `grpc.status = OK` with business success unless the contract makes that explicit.

---

## Implementation Guidance

Recommended steps:

```text
1. Configure MeterProvider in pkg/platform/telemetry.
2. Wire gRPC server metrics instrumentation.
3. Wire gRPC client metrics instrumentation.
4. Export metrics via OTLP to the Collector.
5. Route metrics to Prometheus-compatible backend.
6. Build Grafana dashboards.
7. Add business metrics separately.
```

---

## Practical Rules

```text
Track both client and server metrics.
Use grpc.server.call.duration for service latency and request volume.
Use grpc.client.call.duration for dependency latency.
Use grpc.status != OK for error rates.
Watch per-attempt metrics when retries are enabled.
Watch message-size histograms.
Keep metric labels low-cardinality.
Add business metrics separately.
Start with per-call metrics before LB/xDS metrics.
Build dashboards before alerts.
```

---

## Final Rule

```text
gRPC metrics are bfstore’s baseline RPC health instruments.
```
