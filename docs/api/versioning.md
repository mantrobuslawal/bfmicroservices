# API and Event Versioning

## 1. Purpose

This document defines versioning rules for bfstore APIs and events.

bfstore uses Protocol Buffers for:

```text
gRPC APIs
Kafka event payloads
```

---

## 2. Versioning Principles

```text
avoid breaking consumers
prefer additive changes
reserve removed Protobuf fields
never reuse field numbers
version packages clearly
version Kafka topics for major breaking changes
test compatibility in CI
```

---

## 3. Package Versioning

Service APIs:

```proto
package acme.order.v1;
```

Event payloads:

```proto
package acme.order.events.v1;
```

Shared metadata:

```proto
package acme.events.v1;
```

---

## 4. Event Versioning

Kafka events have related versions:

```text
Protobuf package version
event_version metadata
Kafka topic major version
```

Example:

```text
package acme.order.events.v1
event_version = v1
topic = bfstore.order.orders.v1
```

---

## 5. Compatible Changes

Generally compatible:

```text
adding new fields
adding new messages
adding new event types
adding new enum values when consumers handle unknowns safely
adding comments
```

---

## 6. Breaking Changes

Breaking:

```text
renaming fields
renumbering fields
removing fields without reserving them
changing field type incompatibly
changing field meaning
changing enum number meaning
removing enum values without reserving them
```

---

## 7. Kafka Topic Versioning

Use a new topic version when old consumers cannot safely process the new event stream.

```text
bfstore.order.orders.v1
bfstore.order.orders.v2
```

Do not create a new topic for every additive Protobuf field.

---

## 8. Event Migration Strategy

```text
create v2 Protobuf event message
create v2 topic if needed
dual publish v1 and v2 during migration
upgrade consumers
monitor v1 consumer lag
deprecate v1
remove v1 after agreed retention period
```

---

## 9. Summary

bfstore treats gRPC APIs and Kafka events as first-class Protobuf contracts. Compatibility discipline, Buf checks, and explicit event/topic versioning are mandatory.


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
