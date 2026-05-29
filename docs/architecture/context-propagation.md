# Context Propagation

This document defines the context propagation policy for **bfstore**.

Context propagation is the process of carrying trace, correlation, and request context across service boundaries so distributed work can be observed as one connected flow.

---

## Purpose

bfstore is a distributed platform. One user journey may cross:

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

Context propagation allows engineers to connect telemetry across these boundaries.

It helps answer:

```text
Which services were involved in this checkout?
Where did the request fail?
Which Kafka event belongs to this order flow?
Which logs belong to this trace?
Was the notification linked to the original order?
```

---

## Core Concepts

```text
Context:
  trace/correlation information attached to current work

Propagation:
  moving that context to the next process or service

Carrier:
  where the context is stored during transport
```

Carriers in bfstore:

```text
Go function calls:
  context.Context

gRPC:
  metadata

HTTP:
  headers

Kafka:
  message headers
```

---

## Trace Context

bfstore uses OpenTelemetry trace context.

Standard keys:

```text
traceparent
tracestate, if used
```

Example `traceparent`:

```text
00-a0892f3577b34da6a3ce929d0e0e4736-f03067aa0ba902b7-01
```

Shape:

```text
<version>-<trace-id>-<parent-id>-<trace-flags>
```

Trace context should be propagated through:

```text
gRPC metadata
Kafka headers
HTTP headers where applicable
```

---

## Correlation ID

bfstore also uses a correlation ID.

Purpose:

```text
human-friendly business/request breadcrumb
```

Recommended keys:

```text
gRPC metadata:
  x-correlation-id

Kafka headers:
  correlation_id

Logs:
  correlation_id
```

The correlation ID should remain stable across one business/request flow.

Trace ID and correlation ID are different but complementary.

```text
trace_id:
  OpenTelemetry distributed trace identifier

correlation_id:
  bfstore/platform request breadcrumb
```

---

## gRPC Propagation

For gRPC, propagate:

```text
traceparent
tracestate, if used
x-correlation-id
x-request-id
authorization, where needed
```

Use OpenTelemetry gRPC instrumentation where possible.

Example server shape:

```go
grpcServer := grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()),
)
```

Example client shape:

```go
conn, err := grpc.NewClient(
    target,
    grpc.WithTransportCredentials(creds),
    grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
)
```

bfstore interceptors should still handle:

```text
correlation ID policy
logging
auth
panic recovery
custom metrics
```

---

## Kafka Propagation

Kafka values carry business facts.

Kafka headers carry operational context.

```text
Kafka value:
  protobuf-encoded event

Kafka headers:
  traceparent
  tracestate, if used
  correlation_id
  event_type
  event_version
  content_type
```

Example:

```text
value:
  bfstore.order.events.v1.OrderCreated

headers:
  traceparent = 00-a0892f3577b34da6a3ce929d0e0e4736-f03067aa0ba902b7-01
  correlation_id = checkout-abc-123
  event_type = bfstore.order.events.v1.OrderCreated
  event_version = 1
  content_type = application/protobuf
```

Do not put trace context in the Protobuf event message unless there is a strong business/audit reason.

---

## Logs

Logs should include propagation context where available.

Recommended fields:

```text
trace_id
span_id
correlation_id
service.name
grpc.method or event_type
status
error
```

---

## Baggage Policy

OpenTelemetry baggage can propagate arbitrary key-value pairs.

bfstore should avoid baggage initially unless there is a clear use case.

Never put these in baggage:

```text
raw JWTs
API keys
passwords
payment tokens
card details
CVV
customer email
shipping address
basket JSON
personal data
```

Use baggage deliberately and sparingly.

---

## Security Boundaries

### Public boundary

At `api-gateway`, incoming public context may be forged.

Rules:

```text
validate/sanitise incoming traceparent
generate correlation ID if missing or invalid
strip or reject baggage by default
do not trust propagation context for security decisions
```

### Internal boundary

Between trusted bfstore services:

```text
propagate traceparent
propagate correlation_id
propagate request_id where useful
```

### Third-party boundary

When calling external providers:

```text
do not blindly send internal baggage
do not leak internal auth/context
consider whether traceparent should be propagated
```

---

## Testing Expectations

Recommended tests:

```text
gRPC client sends trace/correlation context
gRPC server extracts trace/correlation context
Kafka producer adds traceparent and correlation_id headers
Kafka consumer extracts trace context from headers
logs include trace_id, span_id, and correlation_id when available
public boundary strips unsafe baggage
```

Testing principle:

```text
Do not only test that spans exist.
Test that the story connects.
```

---

## Practical Rules

```text
Use context.Context inside Go.
Use OpenTelemetry propagation for trace context.
Use gRPC metadata for gRPC propagation.
Use Kafka headers for Kafka propagation.
Use traceparent as the standard trace context key.
Use x-correlation-id for gRPC correlation.
Use correlation_id for Kafka correlation.
Keep Protobuf payloads for business facts.
Do not put secrets or personal data in baggage.
Filter untrusted public context.
Avoid leaking internal context to third parties.
Include trace_id, span_id, and correlation_id in logs.
```

---

## Final Rule

```text
Context propagation is the thread that connects bfstore telemetry.
```
