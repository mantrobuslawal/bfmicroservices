# Process Formation

This document defines process formation guidance for bfstore environments.

---

## Purpose

This document explains:

```text
local/staging/prod-style replica counts
manual scaling commands
capacity assumptions
environment-specific formations
```

---

## Core Rule

```text
Formation should be justified by workload, not vibes.
```

---

## Local Development Formation

Recommended initial local formation:

```text
api-gateway=1
catalog-service=1
basket-service=1
inventory-service=1
order-service=1
payment-service=1
shipping-service=1
notification-worker=1
mysql=1
kafka=1
otel-collector=1
```

Purpose:

```text
prove wiring
run integration tests
debug easily
keep local resource use reasonable
```

---

## Staging Formation

Suggested staging formation:

```text
api-gateway=2
catalog-service=2
basket-service=2
inventory-service=2
order-service=2
payment-service=2
shipping-service=2
notification-worker=2
otel-collector=1 or 2
```

Purpose:

```text
prove high availability basics
test rolling deployments
observe service behaviour with more than one replica
```

---

## Production-style Portfolio Formation

Suggested production-style formation:

```text
api-gateway=3
catalog-service=3
basket-service=3
inventory-service=3
order-service=3
payment-service=3
shipping-service=2
notification-worker=3
search-indexer=2 later
otel-collector=2
```

Purpose:

```text
demonstrate platform thinking
prove scaling patterns
test worker concurrency
test observability and resilience
```

---

## Manual Scaling Commands

Kubernetes:

```bash
kubectl scale deployment/catalog-service --replicas=3
kubectl scale deployment/notification-worker --replicas=4
```

Docker Compose:

```bash
docker compose up --scale catalog-service=3
docker compose up --scale notification-worker=3
```

---

## Formation Metadata

Document:

```text
environment
service/process type
replica count
reason for count
resource requests/limits
known bottlenecks
scaling signals
```

---

## Practical Rules

```text
Start small locally.
Use at least two replicas in staging where useful.
Scale workers based on queues/events, not guesswork.
Do not scale services without checking backing-service pressure.
Record why replica counts exist.
```

---

## Final Rule

```text
Process formation is the runtime shape of bfstore.
```
