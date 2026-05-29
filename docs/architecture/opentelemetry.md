# OpenTelemetry

This document defines how **bfstore** should use OpenTelemetry.

OpenTelemetry is the standard instrumentation and telemetry pipeline approach for producing traces, metrics, and logs from bfstore services.

---

## Purpose

bfstore uses OpenTelemetry to:

```text
trace requests across services
measure service and dependency behaviour
correlate logs with traces
propagate context through gRPC and Kafka
export telemetry to local and production backends
support dashboards, alerts, and troubleshooting
```

OpenTelemetry is not the dashboard itself. It is the instrumentation and export standard.

Telemetry may later be visualised in:

```text
Grafana
Prometheus
Jaeger
Tempo
Loki
vendor observability platforms
```

---

## Initial Architecture

Recommended local development flow:

```text
bfstore Go service
        |
        v
OpenTelemetry SDK/exporter
        |
        v
OpenTelemetry Collector
        |
        +--> Prometheus-compatible metrics backend
        +--> Jaeger/Tempo tracing backend
        +--> structured logs backend/stdout
```

---

## Resource Attributes

Every service should configure:

```text
service.name
service.namespace
service.version
deployment.environment
```

Example:

```text
service.name = bfstore-catalog-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment = local
```

---

## gRPC Instrumentation

gRPC server and client calls should be instrumented.

Server-side instrumentation should capture:

```text
rpc.system = grpc
rpc.service
rpc.method
grpc.status_code
duration
error status
```

Client-side instrumentation should capture outbound calls.

Use interceptors rather than duplicating instrumentation inside handlers.

---

## Kafka Instrumentation

Kafka producers and consumers should propagate trace context using headers.

Recommended Kafka headers:

```text
traceparent
correlation_id
event_type
event_version
```

Producer example:

```text
order-service publishes OrderCreated
trace context is added to Kafka headers
```

Consumer example:

```text
notification-service consumes OrderCreated
trace context is extracted from Kafka headers
notification work is linked to event consumption
```

---

## MySQL Instrumentation

Database calls should be represented with spans where useful.

Example spans:

```text
mysql.catalog.products.select
mysql.orders.insert
mysql.outbox_events.insert
```

Do not record raw SQL with sensitive values.

Useful attributes:

```text
db.system = mysql
db.name
db.operation
table name where safe/useful
```

Avoid sensitive values and high-cardinality attributes.

---

## Logging Integration

Logs should include trace and correlation context.

Recommended fields:

```text
service
level
message
correlation_id
trace_id
span_id
grpc.method
grpc.status_code
error
```

Do not log secrets.

---

## Propagation

For gRPC:

```text
propagate trace context through gRPC metadata
propagate correlation_id through gRPC metadata
```

For Kafka:

```text
propagate trace context through Kafka headers
propagate correlation_id through Kafka headers
```

Without propagation, traces become disconnected.

---

## Initial Instrumentation Scope

Start with:

```text
gRPC server instrumentation
gRPC client instrumentation
catalog-service GetProduct trace
order-service Checkout trace
payment-service AuthorisePayment trace
Kafka OrderCreated publish/consume trace
structured logs with trace_id and correlation_id
```

---

## OpenTelemetry Collector

Use the Collector as the central local telemetry pipeline once multiple services emit telemetry.

Benefits:

```text
services export to one local endpoint
backends can change without rewriting service code
sampling/export policy can be centralised
local dev stack can resemble production shape
```

Recommended local components later:

```text
OpenTelemetry Collector
Prometheus
Grafana
Jaeger or Tempo
Loki or structured stdout
```

---

## Metrics

Initial metrics should include:

```text
gRPC request count
gRPC request duration
gRPC error count
checkout attempts
checkout failures
payment authorisation latency
inventory reservation failures
Kafka publish failures
```

Keep metric labels low-cardinality.

---

## Tracing

Initial traces should focus on:

```text
catalogue browsing
basket update
checkout
order creation
Kafka notification flow
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

## Shutdown

Telemetry should be flushed during graceful shutdown.

Shutdown order should include:

```text
mark service not serving
stop accepting new traffic
finish in-flight requests where possible
flush telemetry exporters/providers
close telemetry resources
exit
```

---

## Sensitive Data Rules

Do not put secrets or sensitive personal data in telemetry.

Avoid:

```text
raw JWTs
payment card data
CVV
passwords
full shipping address
customer email as attribute
full basket JSON
raw Kafka payloads
```

---

## Practical Rules

```text
Use OpenTelemetry for traces and metrics.
Use structured logs and correlate them with trace_id.
Set service.name consistently.
Propagate context through gRPC metadata.
Propagate context through Kafka headers.
Keep attributes low-cardinality.
Do not record secrets or personal data.
Start with catalog and checkout flows.
Use the Collector once multiple services emit telemetry.
Flush telemetry during shutdown.
```

---

## Final Rule

```text
OpenTelemetry is the wiring that lets bfstore explain itself.
```
