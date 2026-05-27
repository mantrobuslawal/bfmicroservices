# Protobuf Style Guide

## 1. Purpose

This document defines Protobuf standards for bfstore.

bfstore uses Protocol Buffers for:

```text
gRPC service APIs
Kafka event payloads
shared common messages
event metadata
```

---

## 2. Package Naming

Service APIs:

```proto
package acme.order.v1;
```

Domain events:

```proto
package acme.order.events.v1;
```

Shared event metadata:

```proto
package acme.events.v1;
```

---

## 3. Event Message Pattern

Kafka event payloads should be typed Protobuf messages.

Recommended:

```proto
message OrderCreatedEvent {
  acme.events.v1.EventMetadata metadata = 1;
  OrderCreated payload = 2;
}
```

Avoid:

```proto
message OrderCreated {
  string payload_json = 1;
}
```

---

## 4. Event Metadata

Shared metadata should live in:

```text
proto/acme/events/v1/event_metadata.proto
```

Recommended fields:

```text
event_id
event_type
event_version
occurred_at
producer
subject
correlation_id
causation_id
trace_id
idempotency_key
```

---

## 5. Field Rules

```text
use lower snake case
never reuse field numbers
reserve removed field numbers
reserve removed field names
prefer additive changes
avoid changing field meaning
```

---

## 6. Enums

Enums must include an unspecified zero value.

```proto
enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING = 1;
  ORDER_STATUS_CONFIRMED = 2;
}
```

---

## 7. Money and Timestamps

Use shared `Money` for monetary values.

Use `google.protobuf.Timestamp` for time values.

Event time should use `occurred_at`.

---

## 8. Kafka Event Compatibility

Event messages must be especially stable because old events may remain in Kafka or DLQs.

Rules:

```text
new fields must be optional from old consumer perspective
old consumers should tolerate unknown fields
removed fields must be reserved
breaking changes require new event version
major breaking changes may require new topic version
```

---

## 9. Buf Standards

CI should run:

```sh
buf lint
buf breaking
buf generate
```

Buf checks should cover both API and event proto packages.

---

## 10. Security and Privacy

Protobuf contracts must not include:

```text
raw card data
CVV
passwords
access tokens
provider secrets
unnecessary personal data
```

Events may be retained, replayed, logged, or copied to DLQs, so event payloads should be especially careful.

---

## 11. Summary

bfstore uses Protobuf for both gRPC APIs and Kafka events. This gives the project one contract-first language across synchronous and asynchronous communication.


---

## Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/events/event-envelope.md
docs/events/event-catalog.md
docs/events/kafka-topic-design.md
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
```
