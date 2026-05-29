# Semantic Conventions

This document defines the semantic convention policy for **bfstore** telemetry.

Semantic conventions are agreed names and meanings for telemetry attributes. They make traces, metrics, logs, resources, and future profiles easier to query, dashboard, alert on, and troubleshoot.

---

## Purpose

bfstore uses semantic conventions to avoid inconsistent telemetry naming.

Bad:

```text
service=orders
service_name=order
app=bfstore-order
svc=order-service
```

Good:

```text
service.name = bfstore-order-service
```

Practical rule:

```text
Use OpenTelemetry semantic conventions for common technical facts.
Use bfstore.* attributes for business/domain facts.
```

---

## Resource Attributes

Every bfstore service must set baseline resource attributes.

Required:

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

Later Kubernetes attributes may include:

```text
k8s.namespace.name
k8s.deployment.name
k8s.pod.name
k8s.container.name
k8s.cluster.name
```

---

## gRPC / RPC Attributes

All gRPC spans, metrics, and correlated logs should use consistent RPC attributes.

Recommended:

```text
rpc.system = grpc
rpc.service
rpc.method
grpc.status_code
```

Examples:

```text
rpc.service = bfstore.order.v1.OrderService
rpc.method = Checkout
```

```text
rpc.service = bfstore.payment.v1.PaymentService
rpc.method = AuthorisePayment
```

---

## Database Attributes

For MySQL operations, use:

```text
db.system = mysql
db.name
db.operation
db.collection.name
```

Examples:

```text
db.system = mysql
db.name = catalog
db.operation = SELECT
db.collection.name = products
```

```text
db.system = mysql
db.name = orders
db.operation = INSERT
db.collection.name = outbox_events
```

Do not record full SQL containing sensitive values.

---

## Kafka / Messaging Attributes

For Kafka operations, use:

```text
messaging.system = kafka
messaging.destination.name
messaging.operation.name
event.type
```

Publishing `OrderCreated`:

```text
messaging.system = kafka
messaging.destination.name = bfstore.order.orders.v1
messaging.operation.name = publish
event.type = bfstore.order.events.v1.OrderCreated
```

Consuming `OrderCreated`:

```text
messaging.system = kafka
messaging.destination.name = bfstore.order.orders.v1
messaging.operation.name = process
event.type = bfstore.order.events.v1.OrderCreated
```

Kafka message layout:

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

---

## HTTP Attributes

At HTTP/API gateway boundaries, use route templates for grouping.

Examples:

```text
http.request.method = POST
url.path = /checkout
http.route = /checkout
http.response.status_code = 201
```

```text
http.request.method = GET
url.path = /products/gopher-desk-lamp
http.route = /products/{slug}
http.response.status_code = 200
```

Do not use raw high-cardinality paths as metric labels.

---

## Log Fields

Structured logs should use consistent fields.

Recommended:

```text
timestamp
level
service.name
deployment.environment
message
trace_id
span_id
correlation_id
rpc.service
rpc.method
grpc.status_code
event.type
error.type
```

---

## bfstore-specific Attributes

Use `bfstore.*` only for business/domain meaning.

Recommended starting attributes:

```text
bfstore.checkout.stage
bfstore.payment.provider
bfstore.stock.reservation_result
bfstore.notification.channel
bfstore.catalog.attribute_count
```

Example values:

```text
bfstore.checkout.stage = payment_authorisation
bfstore.payment.provider = simulated
bfstore.stock.reservation_result = reserved
bfstore.notification.channel = email
bfstore.catalog.attribute_count = 8
```

Do not duplicate standard OpenTelemetry attributes under bfstore names.

Bad:

```text
bfstore.grpc.method = Checkout
```

Good:

```text
rpc.method = Checkout
```

---

## Cardinality Rules

Metric labels must remain low-cardinality.

Good metric labels:

```text
service.name
rpc.method
grpc.status_code
deployment.environment
event.type
bfstore.checkout.stage
```

Avoid metric labels:

```text
order_id
customer_id
customer_email
basket_id
raw_product_name
full_error_message
raw URL paths with IDs
```

Rule:

```text
Metrics are for grouping behaviour.
Traces and logs are for individual investigations.
```

Even in traces and logs, do not include secrets or personal data.

---

## Sensitive Data Rules

Never record:

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

Telemetry must explain behaviour without leaking private or sensitive data.

---

## Starter Standard

```text
Resources:
  service.name
  service.namespace
  service.version
  deployment.environment

gRPC:
  rpc.system
  rpc.service
  rpc.method
  grpc.status_code

MySQL:
  db.system
  db.name
  db.operation
  db.collection.name

Kafka:
  messaging.system
  messaging.destination.name
  messaging.operation.name
  event.type

Logs:
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

bfstore-specific:
  bfstore.checkout.stage
  bfstore.payment.provider
  bfstore.stock.reservation_result
  bfstore.notification.channel
  bfstore.catalog.attribute_count
```

---

## Practical Rules

```text
Use standard OpenTelemetry names first.
Use bfstore.* names only for business meaning.
Set resource attributes in every service.
Keep metric labels low-cardinality.
Keep secrets and personal data out of telemetry.
Use route templates for HTTP metrics.
Use Kafka attributes for messaging operations.
Use database attributes for MySQL spans.
Use consistent log fields for trace/log correlation.
```

---

## Final Rule

```text
Semantic conventions are the grammar of bfstore telemetry.
```
