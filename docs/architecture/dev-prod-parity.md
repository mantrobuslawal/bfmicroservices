# Dev/Prod Parity

This document defines how **bfstore** keeps development, staging, and production-style environments as similar as practical.

---

## Purpose

This document explains:

```text
time, personnel, and tools gaps
same-shape environments
backing service parity
acceptable vs unacceptable differences
```

---

## Core Rule

```text
Parity means same shape, not same size.
```

Local development can be smaller and safer, but it should not be a different architecture.

---

## The Three Gaps

### Time gap

Reduce the time between writing code and deploying it.

```text
small PRs
CI validation
frequent deploys to dev/staging
short feedback loops
```

### Personnel gap

Engineers should understand how their services run.

```text
build
config
deploy
observe
scale
rollback
debug
```

### Tools gap

Development and production-style environments should use the same kinds of tools.

For bfstore:

```text
Go
gRPC
Protobuf
Kafka
MySQL
OpenTelemetry
Docker
Kubernetes later
```

---

## Same-shape Environments

Local should be a smaller model of the real platform.

```text
local:
  Docker Compose
  one replica
  local MySQL
  local Kafka
  local OTel Collector

staging:
  Kubernetes
  multiple replicas where useful
  realistic config/secrets
  realistic observability

production-style:
  Kubernetes
  managed or production-like backing services
  rollout/rollback
  network controls
```

---

## Acceptable Differences

Acceptable local differences:

```text
single Kafka broker instead of cluster
single MySQL container instead of managed MySQL
payment simulator instead of real provider
fake SMTP instead of real email provider
debug telemetry exporter instead of full backend
one replica instead of many
```

---

## Unacceptable Differences

Avoid:

```text
SQLite locally but MySQL in staging
in-memory events locally but Kafka in staging
JSON events locally but Protobuf events in staging
manual config locally but env/secrets in staging
local disk images locally but object storage expected later
```

---

## Practical Rules

```text
Simplify scale locally, not architecture.
Use the same backing service types where behaviour matters.
Use the same protocol contracts locally and in deployed environments.
Keep config shape consistent.
Use realistic integration tests.
Document all intentional differences.
```

---

## Final Rule

```text
Local success should be meaningful evidence that deployed success is likely.
```
