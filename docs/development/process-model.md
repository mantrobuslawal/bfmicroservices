# Process Model

This document defines the process model for **bfstore**.

---

## Purpose

This document explains:

```text
service processes
worker processes
one-off jobs
local Docker Compose processes
Kubernetes Deployments and Jobs
```

---

## Process Types

bfstore may use several process types.

### API process

```text
api-gateway
```

### gRPC service processes

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

### Worker processes

```text
notification-worker
search-indexer later
catalogue-image-processor later
data-quality-checker later
```

### One-off processes

```text
database migration job
seed-data job
backfill job
search reindex job later
stock reconciliation job later
```

---

## Local Docker Compose

In local development, Docker Compose runs process containers.

Example:

```text
mysql
kafka
otel-collector
catalog-service
order-service
payment-service
notification-service
```

Application process containers should remain stateless.

Backing-service containers may hold state through volumes.

---

## Kubernetes Deployments

Long-running service processes should usually run as Kubernetes Deployments.

Examples:

```text
catalog-service Deployment
order-service Deployment
payment-service Deployment
notification-service Deployment
```

Deployments manage:

```text
replicas
rollouts
rollbacks
readiness
restart behaviour
```

---

## Kubernetes Jobs

One-off operational work should run as Kubernetes Jobs where appropriate.

Examples:

```text
migration job
catalogue seed job
data backfill job
search reindex job
```

Jobs should be:

```text
traceable
idempotent where practical
configured through env vars/secrets
run from known image tags
logged
```

---

## Worker Processes

Workers must also be stateless.

Good worker state:

```text
Kafka offsets
database job records
idempotency records
object storage artefacts
```

Bad worker state:

```text
local file progress
in-memory deduplication only
manual SSH command history
```

---

## Process Configuration

All processes should receive config through environment variables or mounted secrets.

Examples:

```text
MYSQL_DSN
KAFKA_BROKERS
OTEL_EXPORTER_OTLP_ENDPOINT
SMTP_HOST
PAYMENT_PROVIDER
```

---

## Practical Rules

```text
Define each process type clearly.
Use Deployments for long-running services.
Use Jobs for controlled one-off tasks.
Keep all app processes stateless.
Do not hide durable state in worker memory.
Configure processes through env vars and secrets.
Make jobs traceable and repeatable.
```

---

## Final Rule

```text
Every bfstore process should have a clear purpose, clear config, and no hidden durable state.
```
