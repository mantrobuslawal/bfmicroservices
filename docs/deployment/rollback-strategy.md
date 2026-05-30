# Rollback Strategy

This document defines rollback expectations for bfstore releases.

It complements:

```text
docs/architecture/build-release-run.md
docs/deployment/release-management.md
```

---

## Purpose

This document explains:

```text
image rollback
config rollback
migration rollback limitations
Kafka/Protobuf compatibility
health checks
observability during rollback
```

---

## Core Rule

```text
Rollback is only safe when releases are identifiable and mostly immutable.
```

---

## Rollback Types

bfstore may need to roll back:

```text
service image
environment config
Kubernetes manifests
database migration
Kafka/Protobuf contract change
feature flag/config value
```

Not all rollback types are equally safe.

---

## Image Rollback

Image rollback is usually the simplest case.

Example:

```text
current:
  order-service:def456

rollback:
  order-service:abc123
```

Requirements:

```text
previous image still exists
previous config is compatible
database schema is compatible
event contracts remain compatible
```

---

## Config Rollback

Config rollback may restore previous values.

Examples:

```text
Kafka broker list
timeout value
OTLP endpoint
payment provider mode
replica count
resource limits
```

Config rollback should be tracked through GitOps, Helm history, Kustomize overlays, or release metadata.

---

## Database Migration Rollback

Database rollbacks are often risky.

Prefer:

```text
backward-compatible migrations
expand/contract pattern
roll-forward fixes where safer
dedicated migration jobs
migration metadata tracking
```

Avoid assuming every migration can be reversed safely.

Rule:

```text
Rollback strategy must consider data, not just code.
```

---

## Kafka and Protobuf Compatibility

Event-driven releases cross service boundaries.

Rollback planning should check:

```text
producer compatibility
consumer compatibility
Protobuf field compatibility
reserved field numbers/names
topic versioning
Buf breaking checks
```

Adding optional fields is usually safer than removing or changing existing fields.

---

## Health Checks

Rollback should be guided by health and readiness signals.

Check:

```text
Kubernetes readiness
gRPC health checks
error rates
latency
Kafka consumer lag
database connectivity
OpenTelemetry telemetry
```

---

## Observability During Rollback

Telemetry should include:

```text
service.name
service.version
deployment.environment.name
release ID if available
```

This helps identify whether rollback improved the issue.

---

## Rollback Runbook Shape

Example rollback flow:

```text
identify failing release
confirm affected service/environment
check migration/contract compatibility
select previous known-good release
apply rollback
monitor health checks
monitor metrics/traces/logs
confirm business flow works
document outcome
```

---

## Practical Rules

```text
Keep previous images available.
Avoid latest tags.
Record release metadata.
Make config history traceable.
Treat database rollback carefully.
Use Buf checks for contract safety.
Monitor rollback with telemetry.
Prefer boring, rehearsed rollback steps.
```

---

## Final Rule

```text
A rollback is not a panic button; it is an engineered path back to a known-good state.
```
