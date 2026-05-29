# Instrumentation

This document defines the instrumentation strategy for **bfstore**.

Instrumentation is the code and configuration that makes services emit useful telemetry, including traces, metrics, and logs.

---

## Purpose

bfstore uses instrumentation to make services:

```text
observable
debuggable
measurable
safe to troubleshoot
safe to operate
```

Instrumentation ensures engineers can understand what services are doing during normal operation, deployments, failures, and performance investigations.

---

## Instrumentation Approaches

bfstore uses two complementary approaches:

```text
automatic instrumentation
code-based instrumentation
```

### Automatic Instrumentation

Use automatic instrumentation for standard technical plumbing:

```text
gRPC server requests
gRPC client calls
HTTP requests
database client calls
runtime metrics
```

### Code-based Instrumentation

Use code-based instrumentation for business meaning:

```text
checkout stages
stock reservation results
payment provider behaviour
Kafka business event publishing
order creation outcomes
notification delivery outcomes
```

Practical rule:

```text
Let automatic instrumentation cover the plumbing.
Use code-based instrumentation for business meaning.
```

---

## gRPC Instrumentation

All bfstore gRPC services should instrument both server and client calls.

Capture:

```text
rpc.system
rpc.service
rpc.method
grpc.status_code
duration
error
trace_id
span_id
correlation_id
```

Server-side shape:

```go
grpcServer := grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()),
)
```

Client-side shape:

```go
conn, err := grpc.NewClient(
    target,
    grpc.WithTransportCredentials(creds),
    grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
)
```

bfstore interceptors should continue to handle:

```text
correlation ID policy
logging
auth
recovery
custom metrics
```

---

## Kafka Instrumentation

Kafka producers and consumers must be instrumented.

Kafka message separation:

```text
key:
  stable business identifier, such as order_id

value:
  protobuf-encoded business event

headers:
  traceparent
  correlation_id
  event_type
  event_version
  content_type
```

Instrumentation must cover both publishing and consuming so asynchronous flows remain observable.

---

## MySQL Instrumentation

Instrument important database operations.

Example span names:

```text
mysql.catalog.products.select
mysql.catalog.product_attributes.select
mysql.orders.insert
mysql.outbox_events.insert
```

Useful attributes:

```text
db.system = mysql
db.name
db.operation
db.collection.name
```

Avoid raw SQL containing sensitive values.

---

## Business Metrics

bfstore should define business-flow metrics in addition to automatic technical metrics.

Initial business metrics:

```text
checkout.attempts_total
checkout.completed_total
checkout.failed_total
checkout.duration
payment.authorisation.attempts_total
payment.authorisation.failed_total
inventory.reservation.failed_total
kafka.publish.failed_total
notification.delivery.failed_total
```

Metric labels must remain low-cardinality.

Allowed-style labels:

```text
service
method
status
environment
failure_reason_class
event_type
```

Avoid:

```text
order_id
customer_email
basket_id
full_error_message
raw_product_name
```

---

## Logs

Logs should be structured and correlated.

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
event_type
error
```

Do not log secrets or sensitive personal data.

---

## Semantic Conventions

Use OpenTelemetry semantic conventions for common technical facts.

Prefer:

```text
service.name
deployment.environment
rpc.system
rpc.service
rpc.method
db.system
messaging.system
messaging.destination.name
```

Use bfstore-specific attributes for business facts:

```text
bfstore.checkout.stage
bfstore.payment.provider
bfstore.stock.reservation_result
```

---

## Resources

Every service should set resource attributes:

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

---

## Instrumentation Scopes

Instrumentation scopes identify which library or package emitted telemetry.

Examples:

```text
github.com/mantrobuslawal/bfstore/pkg/platform/telemetry
github.com/mantrobuslawal/bfstore/pkg/platform/grpc/interceptors
github.com/mantrobuslawal/bfstore/services/order-service
```

Resources identify the actor.

Instrumentation scopes identify the narrator.

---

## Rollout Plan

### Phase 1: Foundation

```text
create pkg/platform/telemetry
set service.name, environment, version
configure OTLP exporter
configure propagators
flush telemetry on shutdown
```

### Phase 2: gRPC

```text
instrument gRPC server calls
instrument gRPC client calls
include trace/log correlation
record method/status/duration
```

### Phase 3: Catalog Path

```text
instrument CatalogService/GetProduct
instrument MySQL product queries
add catalog latency/error metrics
```

### Phase 4: Checkout Path

```text
instrument OrderService/Checkout
instrument inventory/payment/shipping calls
instrument Kafka OrderCreated publish
add checkout success/failure metrics
```

### Phase 5: Async Consumers

```text
instrument notification-service Kafka consume
link OrderCreated consume to trace context
add notification delivery metrics
```

### Phase 6: Profiling

```text
add profiles later
investigate CPU/memory hotspots
optimise based on evidence
```

---

## Sensitive Data Rules

Never record:

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
high-cardinality metric labels
```

Practical rule:

```text
Telemetry should explain behaviour, not leak data.
```

---

## Practical Rules

```text
Use automatic instrumentation for standard plumbing.
Use code-based instrumentation for business meaning.
Instrument both gRPC servers and clients.
Instrument Kafka producers and consumers.
Instrument key MySQL operations.
Correlate logs with trace_id and correlation_id.
Use semantic conventions for common technical attributes.
Use bfstore-specific attributes for business facts.
Keep metric labels low-cardinality.
Never record secrets or personal data.
Start with catalog and checkout flows.
Flush telemetry during graceful shutdown.
```

---

## Final Rule

```text
Instrumentation is how bfstore services explain themselves.
```
