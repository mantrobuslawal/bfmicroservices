# Retry and DLQ Strategy

## 1. Purpose

This document defines retry and dead-letter queue strategy for bfstore Kafka consumers.

bfstore uses Kafka for events and Protocol Buffers for event payloads.

---

## 2. Failure Types

| Failure Type | Example | Handling |
|---|---|---|
| Transient dependency failure | database outage | retry |
| Permanent validation failure | invalid event payload | DLQ |
| Protobuf decode failure | corrupt or wrong message type | DLQ and alert |
| Unsupported event version | consumer cannot process v2 | DLQ or safe rejection |
| Duplicate event | same event_id already processed | ignore or return existing result |
| Poison message | always fails | bounded retry then DLQ |

---

## 3. Protobuf-Specific Failure Handling

Consumers must explicitly handle:

```text
deserialisation failure
unknown event type
unsupported event version
missing required business fields
invalid enum values
contract validation failure
```

Handling:

```text
log event_type and event_version from headers where possible
record decode failure metric
send original binary payload to DLQ where appropriate
include safe diagnostic metadata
alert if decode failures exceed threshold
```

---

## 4. DLQ Payload Requirements

For Protobuf payloads, DLQ records should preserve:

```text
original binary payload
headers
topic
partition
offset
consumer group
failure reason
failure category
decode error where applicable
timestamp
```

Sensitive data rules still apply.

---

## 5. Consumer Flow

```text
read Kafka message
decode Protobuf payload
validate metadata
validate supported event version
check idempotency/processed event record
process event
commit offset after successful processing
```

---

## 6. Metrics

Track:

```text
events consumed
events processed successfully
events retried
events sent to DLQ
Protobuf decode failures
unsupported event versions
contract validation failures
consumer lag
oldest failed event age
```

---

## 7. Replay Rules

Before replay:

```text
identify root cause
confirm consumer compatibility
confirm idempotency behaviour
confirm side-effect safety
decide whether replay is normal or corrective
```

Notification and payment-related side effects require extra care.

---

## 8. Testing Expectations

```text
transient failure retry
permanent failure DLQ
Protobuf decode failure DLQ
unsupported event version handling
duplicate event handling
retry exhaustion
DLQ metadata preservation
safe replay behaviour
```

---

## 9. Summary

bfstore treats Protobuf decode failures, unsupported event versions, and poison messages as first-class operational cases.


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
