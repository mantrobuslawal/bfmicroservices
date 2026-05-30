# Scaling Strategy

This document defines bfstore's scaling strategy for services and workers.

---

## Purpose

This document explains:

```text
horizontal scaling
vertical scaling limits
Kafka worker scaling
DB connection pressure
gRPC deadlines
metrics required before autoscaling
```

---

## Core Rule

```text
Scale deliberately, measure before and after.
```

---

## Horizontal Scaling

Horizontal scaling means increasing replicas.

Examples:

```text
catalog-service replicas: 1 -> 3
notification-worker replicas: 1 -> 4
api-gateway replicas: 2 -> 4
```

This works best when processes are stateless and durable state lives in backing services.

---

## Vertical Scaling

Vertical scaling means increasing resource limits/requests.

Examples:

```text
more CPU
more memory
larger DB instance
larger Kafka cluster
```

Vertical scaling can help, but one process or node can only grow so far.

---

## Kafka Worker Scaling

Kafka worker scaling depends on:

```text
topic partition count
consumer group behaviour
worker idempotency
downstream rate limits
retry strategy
dead-letter handling later
```

Rule:

```text
More workers than partitions may not increase parallelism for one consumer group.
```

---

## Database Pressure

More service replicas can increase:

```text
DB connections
query load
lock contention
transaction conflicts
CPU/memory pressure
```

Configure:

```text
max open connections
max idle connections
connection lifetime
request deadlines
query indexes
```

---

## gRPC Scaling Controls

Use:

```text
deadlines
cancellation
server-side limits where useful
client-side load balancing later
interceptors
rate limiting later
circuit breaking later
```

Rule:

```text
Every RPC should have a deadline.
```

---

## Scaling Signals

Useful metrics:

```text
request rate
p95/p99 latency
error rate
CPU/memory
Kafka consumer lag
worker processing duration
DB query latency
DB connection pool usage
pod restarts
gRPC deadline exceeded count
```

---

## Autoscaling Later

Potential tools:

```text
Horizontal Pod Autoscaler
KEDA
Cluster Autoscaler
```

Do not add autoscaling before:

```text
manual scaling is understood
resource requests/limits exist
metrics are reliable
workers are idempotent
readiness/liveness are correct
backing-service bottlenecks are visible
```

---

## Practical Rules

```text
Scale stateless services horizontally.
Scale workers with partition count and downstream limits in mind.
Watch DB pressure when adding replicas.
Use deadlines and timeouts.
Measure before and after scaling.
Add autoscaling later, not first.
```

---

## Final Rule

```text
Scaling should solve a measured bottleneck, not hide an unknown one.
```
