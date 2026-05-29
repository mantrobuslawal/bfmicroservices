# Observability

This document defines the observability principles for **bfstore**.

Observability is the ability to understand what the system is doing from the outside by using telemetry emitted by services and infrastructure.

bfstore uses observability to make the platform:

```text
debuggable
measurable
operable
safer to deploy
easier to troubleshoot
better aligned with user-facing reliability
```

---

## Purpose

bfstore is a distributed ecommerce platform. A single user journey may cross several services.

Example checkout flow:

```text
api-gateway
-> basket-service
-> order-service
-> inventory-service
-> payment-service
-> shipping-service
-> Kafka
-> notification-service
```

Observability exists so engineers can answer:

```text
Where did this request go?
Which service failed?
Which dependency was slow?
Did the Kafka event publish?
Did the notification consume the event?
Did a deployment make things worse?
```

---

## Core Telemetry Signals

bfstore uses:

```text
traces
metrics
logs
```

### Traces

Traces show the journey of a request or workflow across services.

Example:

```text
OrderService/Checkout
├── BasketService/GetBasket
├── InventoryService/ReserveStock
├── PaymentService/AuthorisePayment
├── ShippingService/CreateShipment
├── MySQL INSERT order
└── Kafka publish OrderCreated
```

### Metrics

Metrics are numeric measurements over time.

Examples:

```text
request count
request duration
error count
checkout success rate
Kafka publish failures
database query latency
```

### Logs

Logs are timestamped events or messages emitted by services.

bfstore logs should be structured and include enough context to correlate with traces.

Example fields:

```text
service
method
grpc_status
correlation_id
trace_id
span_id
order_id where safe/useful
message
error
```

---

## Metrics, Traces, and Logs Together

```text
Metrics tell you something is wrong.
Traces show where it went wrong.
Logs explain details around what happened.
```

Example:

```text
Metric:
  Checkout error rate increased.

Trace:
  Most failed checkouts timeout inside payment-service.

Log:
  payment provider simulation timed out for order ord_123.
```

---

## Correlation IDs and Trace IDs

bfstore uses both correlation IDs and OpenTelemetry trace IDs.

```text
correlation_id:
  business/request tracking ID chosen by bfstore

trace_id:
  OpenTelemetry distributed trace ID

span_id:
  ID for a single operation inside a trace
```

Logs should include both when available.

Correlation IDs should propagate through:

```text
gRPC metadata
Kafka headers
structured logs
```

Trace context should propagate through:

```text
gRPC metadata
Kafka headers
OpenTelemetry instrumentation
```

---

## Service Naming

Recommended service names:

```text
bfstore-api-gateway
bfstore-catalog-service
bfstore-basket-service
bfstore-inventory-service
bfstore-order-service
bfstore-payment-service
bfstore-shipping-service
bfstore-notification-service
```

These names should be used consistently in:

```text
OpenTelemetry resource attributes
logs
dashboards
alerts
deployment labels
Kubernetes workloads
```

---

## Resource Attributes

Every service should set:

```text
service.name
service.version
deployment.environment
service.namespace
```

Example:

```text
service.name = bfstore-order-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment = local
```

---

## Sensitive Data Rules

Telemetry must not contain secrets or sensitive personal data.

Do not record:

```text
raw JWTs
passwords
card numbers
CVV
full shipping addresses
customer email as high-cardinality attribute
full basket JSON
payment provider secrets
database passwords
```

Practical rule:

```text
Telemetry should explain behaviour, not leak data.
```

---

## SLIs and SLOs

Reliability should be measured from the user or business perspective.

Example SLI:

```text
Checkout success rate
```

Example SLO:

```text
99.5% of checkout attempts complete successfully over 30 days
```

Suggested bfstore SLIs:

| Service | SLI |
|---|---|
| `catalog-service` | `GetProduct` success rate |
| `basket-service` | `AddItem` success rate |
| `order-service` | Checkout success rate |
| `payment-service` | Payment authorisation latency |
| `notification-service` | Order confirmation delivery delay |

---

## Initial Observability Scope

Start with:

```text
gRPC server request count/duration/status
gRPC client request count/duration/status
structured logs with correlation_id and trace_id
catalog-service GetProduct trace
order-service Checkout trace
MySQL query spans for key paths
Kafka publish/consume spans for key events
```

Recommended first paths:

```text
api-gateway -> catalog-service -> MySQL
order-service -> inventory/payment/shipping -> Kafka
```

---

## Operational Questions

bfstore observability should answer:

```text
Which service is slow?
Which dependency is failing?
Are errors user-facing or internal?
Did checkout fail before or after payment?
Did stock reservation happen?
Was OrderCreated published?
Did notification-service consume the event?
Are retries happening?
Are failures isolated to one pod/version?
Did a deployment make latency worse?
```

---

## Practical Rules

```text
Use traces, metrics, and logs together.
Use OpenTelemetry as the standard instrumentation approach.
Use structured logs.
Include correlation_id and trace_id in logs where available.
Propagate context through gRPC metadata.
Propagate context through Kafka headers.
Use consistent service names.
Avoid secrets and sensitive personal data in telemetry.
Start with high-value request paths.
Define SLIs/SLOs around user-facing reliability.
```

---

## Final Rule

```text
Observability is how bfstore becomes operable, not just runnable.
```
