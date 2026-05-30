# Concurrency Process Model

This document defines bfstore's process model for concurrency and scale-out.

---

## Purpose

This document explains:

```text
process types
process formation
web/API processes
gRPC service processes
worker processes
one-off jobs
how concurrency differs from statelessness
```

---

## Core Rule

```text
Different workloads should be represented as different process types.
```

---

## Process Types

bfstore process types include:

```text
api-gateway
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-worker
search-indexer later
migration-job
seed-data-job
```

Each process type should have a clear responsibility and lifecycle.

---

## Process Formation

Process formation is the number of running processes of each type.

Example:

```text
api-gateway=2
catalog-service=3
order-service=2
notification-worker=4
```

Formation should be based on workload, reliability, and capacity needs.

---

## API Process

The API gateway handles:

```text
public HTTP traffic
routing to internal gRPC services
authentication later
request/response shaping
```

It should scale independently from background workers.

---

## gRPC Service Processes

Core gRPC service processes include:

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
```

They should remain stateless and scale horizontally through replicas.

---

## Worker Processes

Worker processes handle asynchronous/background work.

Examples:

```text
notification-worker
search-indexer later
catalogue-image-processor later
report-worker later
```

Workers should use Kafka consumer groups, idempotency, retries, and clear observability.

---

## One-off Jobs

Jobs handle finite operational work.

Examples:

```text
database migrations
seed data
backfills
search reindex
stock reconciliation
```

Jobs should be traceable and repeatable where practical.

---

## Concurrency vs Statelessness

Statelessness says:

```text
a process must not own durable state
```

Concurrency says:

```text
run the right number of processes for each kind of work
```

They work together.

---

## Practical Rules

```text
Separate request-serving and background work where appropriate.
Use Deployments for long-running services/workers.
Use Jobs/CronJobs for finite/scheduled work.
Scale each process type independently.
Do not hide workers inside unrelated services.
Document process formation per environment.
```

---

## Final Rule

```text
Concurrency starts with naming the work clearly.
```
