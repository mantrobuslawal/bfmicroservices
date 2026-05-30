# Worker Processes

This document defines worker process guidance for bfstore.

---

## Purpose

This document explains:

```text
notification workers
search indexers later
Kafka consumer groups
idempotency
retry handling
dead-letter topics later
```

---

## Core Rule

```text
Background work should have explicit worker processes.
```

Do not hide long-running background work inside unrelated request-serving services.

---

## Notification Worker

Responsibilities:

```text
consume OrderCreated events
send confirmation emails
record delivery attempts where needed
handle retries
avoid duplicate notifications where practical
```

Scaling signal:

```text
Kafka consumer lag
email delivery latency
retry/failure rate
email provider rate limits
```

---

## Search Indexer Later

Responsibilities:

```text
consume ProductChanged events
update search projection/index
rebuild facets
handle reindex jobs
```

Scaling signal:

```text
Kafka consumer lag
index update latency
search backend write latency
CPU/memory usage
```

---

## Kafka Consumer Groups

Workers should use consumer groups for parallelism.

Example:

```text
topic: bfstore.order.orders.v1
consumer group: notification-service
worker replicas: 3
```

Parallelism is limited by topic partitions for that consumer group.

---

## Idempotency

Workers must be idempotent because Kafka events can be retried or replayed.

Examples:

```text
event_id
order_id
notification_type
delivery status record
unique constraints
idempotency keys
```

Rule:

```text
A repeated event should not create incorrect duplicate side effects.
```

---

## Retry Handling

Retries should be deliberate.

Track:

```text
attempt count
last attempted at
failure reason where safe
next retry time where applicable
```

Avoid infinite hot loops.

---

## Dead-letter Topics Later

Later, failed events may go to DLQs.

Possible examples:

```text
bfstore.notification.notifications.dlq.v1
bfstore.search.indexing.dlq.v1
```

DLQs need:

```text
alerting
inspection process
replay process
data safety rules
```

---

## Worker Observability

Track:

```text
consumer lag
processing duration
success/failure count
retry count
DLQ count later
duplicate suppression count
downstream provider latency
```

---

## Practical Rules

```text
Use workers for asynchronous work.
Use Kafka consumer groups.
Make handlers idempotent.
Respect downstream rate limits.
Track retries and failures.
Do not block checkout on non-critical notifications.
Use DLQs later with clear runbooks.
```

---

## Final Rule

```text
Workers turn background work into a first-class, observable part of the platform.
```
