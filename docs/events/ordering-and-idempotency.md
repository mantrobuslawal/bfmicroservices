# Ordering and Idempotency

## 1. Purpose

This document defines ordering and idempotency rules for bfstore event-driven workflows.

bfstore uses Kafka for events and Protocol Buffers for event payloads.

---

## 2. Core Principles

```text
assume events may be delivered more than once
assume consumers may restart at any time
design consumers to be idempotent
choose Kafka keys based on ordering requirements
use event_id for deduplication
use business IDs for partitioning where ordering matters
```

---

## 3. Protobuf Event Metadata

Every event should include metadata with:

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

`event_id` is the primary deduplication identifier.

---

## 4. Kafka Key Strategy

```text
Order lifecycle events       -> order_id
Product catalogue events     -> product_id
Payment events               -> order_id or payment_id
Shipment events              -> order_id or shipment_id
Inventory reservation events -> reservation_id or product/variant key
Notification events          -> notification_id
```

Ordering is only guaranteed within a Kafka partition.

---

## 5. Idempotent Producers

Recommended:

```text
use transactional outbox for critical events
persist event_id
persist idempotency_key for command-triggered events
publish from durable outbox state
do not generate a new event_id for the same committed fact
```

`OrderCreated` is a strong outbox candidate.

---

## 6. Idempotent Consumers

Consumers should record processed events where duplicate side effects are harmful.

Examples:

```text
notification-service stores processed OrderCreated event IDs
search-service tracks catalogue projection offsets and event IDs
recommendation-service deduplicates order/review/basket signals
```

---

## 7. Protobuf Compatibility and Idempotency

Consumers should tolerate compatible Protobuf changes.

Rules:

```text
new fields should not break old consumers
unknown fields should be safely ignored where supported
removed fields must be reserved
field numbers must never be reused
unsupported major event versions must be handled safely
```

---

## 8. Checkout Idempotency

Critical commands:

```text
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
```

Each should use:

```text
idempotency_key
request_hash
stored operation result
safe retry behaviour
```

---

## 9. Testing Expectations

```text
duplicate event processing
same event_id processed twice
same idempotency_key command retried
same idempotency_key with different request rejected
consumer handles compatible Protobuf additions
consumer handles unsupported event version safely
event replay does not duplicate notifications
```

---

## 10. Summary

Protobuf event metadata, Kafka keys, idempotency keys, processed event records, and outbox patterns work together to make bfstore event-driven workflows safe and predictable.


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
