# Telemetry Propagation

This document describes the planned propagation helpers for `pkg/platform/telemetry`.

The implementation may later include helpers for:

```text
OpenTelemetry propagator setup
gRPC propagation support
Kafka header carriers
correlation ID helpers
log enrichment helpers
```

---

## Purpose

Propagation helpers should make it easy for bfstore services to carry telemetry context across:

```text
Go function calls
gRPC calls
Kafka messages
logs
```

The goal is to ensure traces, logs, and events connect across service boundaries.

---

## Core Responsibilities

This package area may provide:

```text
default OpenTelemetry propagator configuration
Kafka header carrier implementation
correlation ID context helpers
trace/log enrichment helpers
safe baggage policy helpers, if needed
```

It should not contain business workflow logic.

---

## Default Propagator

bfstore should use W3C Trace Context propagation.

Expected keys:

```text
traceparent
tracestate, if used
```

Baggage should not be enabled broadly until there is a clear use case and policy.

---

## Kafka Header Carrier

Kafka uses message headers rather than HTTP/gRPC metadata.

A small carrier adapter can allow OpenTelemetry to inject/extract context from Kafka headers.

Conceptual shape:

```go
type KafkaHeaderCarrier struct {
    Headers []kafka.Header
}

func (c KafkaHeaderCarrier) Get(key string) string {
    // Return header value for key.
}

func (c KafkaHeaderCarrier) Set(key, value string) {
    // Set header key/value.
}

func (c KafkaHeaderCarrier) Keys() []string {
    // Return all header keys.
}
```

The exact type depends on the Kafka client library chosen.

---

## Kafka Headers

Recommended bfstore Kafka headers:

```text
traceparent
tracestate, if used
correlation_id
event_type
event_version
content_type
```

Example:

```text
traceparent = 00-a0892f3577b34da6a3ce929d0e0e4736-f03067aa0ba902b7-01
correlation_id = checkout-abc-123
event_type = bfstore.order.events.v1.OrderCreated
event_version = 1
content_type = application/protobuf
```

---

## Correlation ID Helpers

The package may provide helpers to:

```text
read correlation ID from context
write correlation ID to context
inject correlation ID into gRPC metadata
inject correlation ID into Kafka headers
extract correlation ID from Kafka headers
```

This should align with:

```text
gRPC metadata key: x-correlation-id
Kafka header key: correlation_id
```

---

## Log Enrichment

Helpers may support adding these fields to logs:

```text
trace_id
span_id
correlation_id
```

Example target log:

```text
level=error
service=bfstore-payment-service
trace_id=a0892f3577b34da6a3ce929d0e0e4736
span_id=f03067aa0ba902b7
correlation_id=checkout-abc-123
message="payment authorisation timed out"
```

---

## Baggage Policy

Do not propagate baggage broadly by default.

Never allow baggage to contain:

```text
raw JWTs
API keys
passwords
card details
CVV
payment tokens
customer email
shipping address
basket JSON
personal data
```

Baggage support should be added only after a clear use case exists.

---

## Public Boundary Behaviour

At public ingress points, such as `api-gateway`, helpers may support:

```text
validate traceparent
generate correlation ID if missing or invalid
strip baggage
avoid trusting caller-supplied propagation context for security decisions
```

Security decisions must come from authentication and authorisation systems, not trace context.

---

## Testing Guidance

Recommended tests:

```text
Kafka carrier can set and get traceparent
Kafka carrier returns expected keys
correlation ID is injected into Kafka headers
correlation ID is extracted from Kafka headers
log enrichment adds trace_id and span_id when span context exists
unsafe baggage is ignored or rejected, if baggage support exists
```

---

## What This Package Should Not Do

Do not put business logic here.

Bad:

```text
propagation package knows checkout rules
propagation package modifies OrderCreated event payloads
propagation package decides payment behaviour
```

Good:

```text
propagation package carries telemetry context
services decide business behaviour
```

---

## Practical Rules

```text
Use context.Context inside Go.
Use gRPC metadata for gRPC boundaries.
Use Kafka headers for Kafka boundaries.
Keep Protobuf messages focused on business facts.
Propagate traceparent.
Propagate correlation_id.
Avoid baggage by default.
Never propagate secrets or personal data.
Keep helpers small and testable.
```

---

## Final Rule

```text
Propagation helpers should carry the breadcrumbs, not the business parcel.
```
