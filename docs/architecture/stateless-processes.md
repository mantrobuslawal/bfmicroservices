# Stateless Processes

This document defines how **bfstore** runs services as stateless, share-nothing processes.

---

## Purpose

This document explains:

```text
process model
stateless/share-nothing rule
memory vs persistent state
filesystem rules
sticky session policy
service-specific examples
```

---

## Core Rule

```text
Processes do work.
Backing services keep state.
```

A bfstore process should be replaceable without losing business data.

---

## Stateless Rule

A process must not rely on its own memory or filesystem for anything that must survive beyond a request, job, restart, or rollout.

Durable state belongs in backing services:

```text
MySQL
Kafka
Redis later
object storage later
```

---

## Share-nothing Rule

Replicas should not rely on another replica’s memory, files, or local cache.

Good:

```text
basket-service replica 1 -> basket MySQL
basket-service replica 2 -> basket MySQL
```

Bad:

```text
basket-service replica 1 stores basket in memory
basket-service replica 2 cannot see it
```

---

## Memory Policy

Allowed:

```text
temporary request data
parsed config
DB connection pools
compiled regex
short-lived cache that is safe to lose
```

Not allowed as source of truth:

```text
basket state
order state
payment state
notification sent status
product truth
checkout progress
```

Rule:

```text
If losing memory changes correctness, the data belongs elsewhere.
```

---

## Filesystem Policy

Container filesystem is temporary.

Do not store durable business data on local disk.

Bad:

```text
product images in /uploads
orders in local JSON files
notification history in local files
```

Good:

```text
product image metadata in MySQL
product images in object storage later
orders in order MySQL
notification status in DB if required
```

Temporary files are acceptable for single-operation processing.

---

## Sticky Session Policy

Do not rely on sticky sessions for correctness.

Bad:

```text
api-gateway stores sessions in memory
load balancer routes same user to same pod
```

Better:

```text
JWT/session tokens validated by any replica
or session state in Redis with TTL later
```

Rule:

```text
Any healthy replica should be able to handle the next request.
```

---

## Service-specific Guidance

### catalog-service

Source of truth:

```text
catalog MySQL
object storage later for images
```

Avoid storing product truth in memory or product images on local disk.

### basket-service

Source of truth:

```text
basket MySQL
Redis later only if explicitly chosen
```

Do not store basket state in a process-level Go map.

### order-service

Source of truth:

```text
order MySQL
outbox table
Kafka event stream
```

Persist important checkout state transitions.

### payment-service

Source of truth:

```text
payment MySQL
idempotency records
provider response records
```

Payment state must survive process restarts.

### notification-service

Source of truth/progress:

```text
Kafka offsets
notification delivery records if needed
idempotency keys
```

Consumers must be safe to restart and replay.

---

## Kubernetes Implications

Pods may be:

```text
deleted
rescheduled
evicted
scaled
rolled out
rolled back
```

Deleting a pod must not delete business data.

Rule:

```text
If deleting a pod deletes business data, the pod was doing the wrong job.
```

---

## Practical Rules

```text
Treat containers as disposable.
Treat memory as temporary.
Treat local disk as temporary.
Store durable state in backing services.
Avoid sticky sessions.
Make replicas interchangeable.
Make workers idempotent.
Use Kafka and databases for durable progress.
```

---

## Final Rule

```text
A bfstore process can disappear; bfstore business state must not.
```
