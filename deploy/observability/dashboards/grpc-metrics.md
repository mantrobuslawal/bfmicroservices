# gRPC Metrics Dashboards

This document defines the initial dashboard expectations for bfstore gRPC OpenTelemetry metrics.

It complements:

```text
docs/architecture/grpc-opentelemetry-metrics.md
pkg/platform/telemetry/grpc-metrics.md
```

---

## Purpose

Dashboards should make bfstore RPC health visible across services.

They should answer:

```text
Which methods are called most?
Which methods are slow?
Which methods are failing?
Which dependencies are slow from the caller view?
Are retries increasing?
Are message sizes growing?
```

---

## Dashboard 1: gRPC Overview

Panels:

```text
request volume by service.name and grpc.method
server latency p50/p95/p99 by grpc.method
error rate by grpc.method and grpc.status
top slow methods
top failing methods
```

Primary metric:

```text
grpc.server.call.duration
```

---

## Dashboard 2: Service Dependency Latency

Panels:

```text
client call duration p95 by service.name and grpc.target
client call duration p95 by grpc.method
dependency error rate by grpc.target
dependency timeout rate
```

Primary metric:

```text
grpc.client.call.duration
```

Useful for:

```text
order-service -> payment-service
order-service -> inventory-service
order-service -> shipping-service
basket-service -> catalog-service
```

---

## Dashboard 3: Checkout RPC Health

Panels:

```text
OrderService/Checkout request volume
OrderService/Checkout p95 latency
OrderService/Checkout error rate
OrderService client latency to PaymentService
OrderService client latency to InventoryService
OrderService client latency to ShippingService
payment retry count later
payment retry delay later
```

---

## Dashboard 4: Catalogue RPC Health

Panels:

```text
CatalogService/GetProduct latency
CatalogService/ListProducts latency
CatalogService/ListProducts response size
CatalogService status-code breakdown
top catalogue errors
```

Watch especially:

```text
grpc.server.call.sent_total_compressed_message_size
```

---

## Dashboard 5: Message Size

Panels:

```text
sent compressed message size by grpc.method
received compressed message size by grpc.method
largest methods by response size
largest methods by request size
```

Use this to catch bloated APIs.

---

## Dashboard 6: Retry Behaviour Later

Add after retries are deliberately enabled.

Panels:

```text
grpc.client.call.retries by grpc.method and grpc.target
grpc.client.call.transparent_retries by grpc.method and grpc.target
grpc.client.call.retry_delay p95
grpc.client.attempt.duration p95
attempt count compared with call count
```

---

## Alert Ideas

Start with dashboards before alerts.

Later alerts:

```text
PaymentService AuthorisePayment error rate > 2% for 5 minutes
OrderService Checkout p95 latency > 3 seconds for 5 minutes
Payment dependency client p95 latency > 2 seconds for 5 minutes
DEADLINE_EXCEEDED count increased sharply
CatalogService ListProducts response size above threshold
retry count increased sharply after deployment
```

---

## Label Grouping

Useful groupings:

```text
service.name
grpc.method
grpc.status
grpc.target
deployment.environment.name
```

Avoid:

```text
order_id
customer_email
basket_id
product_slug
full error message
```

---

## Practical Rules

```text
Start with overview dashboards.
Compare client and server latency.
Graph error rate as a metric.
Track message size.
Add retry dashboards only once retries are enabled.
Keep dashboard labels low-cardinality.
Build dashboards before tuning alerts.
```

---

## Final Rule

```text
Dashboards turn gRPC metrics into operational understanding.
```
