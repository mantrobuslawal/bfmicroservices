# Retry and Idempotency

This document defines retry and idempotency guidance for bfstore services and workers.

---

## Purpose

This document explains:

```text
retry-safe operations
idempotency keys
transactions
outbox pattern
Kafka replay
duplicate notification prevention
```

---

## Core Rule

```text
If work may be interrupted and retried, it must be safe to retry.
```

Disposability depends on idempotency because processes may shut down, crash, restart, or receive duplicate requests/events.

---

## Retry Sources

Retries can happen because of:

```text
client timeout
network failure
pod restart
worker crash
Kafka redelivery
database deadlock/retry
payment provider timeout
deployment rollout
```

---

## Idempotency Keys

Use stable keys to recognise duplicate work.

Examples:

```text
checkout_id:
  prevents duplicate order creation

payment_idempotency_key:
  prevents duplicate payment authorisation

reservation_id:
  prevents duplicate stock reservation

event_id:
  prevents duplicate event handling

order_id + notification_type:
  prevents duplicate customer notification
```

---

## Transactions

Use database transactions for critical state changes.

Examples:

```text
create order and order items atomically
reserve inventory atomically
record payment attempt/result atomically
write outbox event with business state change
```

---

## Outbox Pattern

For critical events:

```text
write business state
write outbox event in same DB transaction
publisher later publishes event to Kafka
mark outbox event as published
```

This helps recover if a process dies after DB write but before Kafka publish.

---

## Kafka Replay

Kafka consumers must tolerate replay.

Rules:

```text
commit offset only after successful processing
make handlers idempotent
track processed event IDs where needed
use unique constraints for duplicate prevention
send failed messages to DLQ later where appropriate
```

---

## Duplicate Notification Prevention

For notification-worker:

```text
event_id
order_id
notification_type
delivery_status
sent_at
provider_message_id if available
```

Use uniqueness rules to prevent sending the same notification incorrectly.

---

## Payment Safety

Payment operations need strong idempotency.

Store:

```text
payment attempt ID
idempotency key
provider request/reference
status
authorised amount
failure/timeout result where safe
```

Rule:

```text
Never rely on process memory to remember whether payment happened.
```

---

## Practical Rules

```text
Use idempotency keys for externally visible commands.
Use transactions for critical state changes.
Use outbox for critical event publication.
Make Kafka handlers replay-safe.
Commit offsets after success.
Use unique constraints to prevent duplicates.
Avoid blind retries without limits.
Track attempts and failures.
```

---

## Final Rule

```text
Retries are normal; duplicate business effects are not.
```
