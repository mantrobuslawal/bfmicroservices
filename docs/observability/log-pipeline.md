# Log Pipeline

This document defines the bfstore log capture and routing model.

---

## Purpose

This document explains:

```text
local stdout
Docker Compose logs
Kubernetes logs
collector/router
central storage/search
retention and alerting later
```

---

## Core Rule

```text
The app writes logs.
The platform routes logs.
```

---

## Local Logging

Local services write JSON logs to stdout/stderr.

Useful commands:

```bash
docker compose logs -f catalog-service
docker compose logs -f order-service
docker compose logs -f notification-worker
```

---

## Kubernetes Logging

Kubernetes captures container stdout/stderr.

Useful commands:

```bash
kubectl logs deployment/order-service
kubectl logs pod/order-service-abc123
kubectl logs -f deployment/notification-worker
```

---

## Collection Options

Possible log collectors:

```text
OpenTelemetry Collector
Fluent Bit
Fluentd
Vector
cloud-native logging agents
```

Choose later based on the observability stack.

---

## Storage/Search Options

Possible destinations:

```text
Loki
OpenSearch
Elasticsearch
Splunk
cloud logging
object storage for archive
```

bfstore should keep the application independent from the chosen destination.

---

## Staged Adoption

Recommended stages:

```text
Phase 1:
  JSON logs to stdout
  docker compose logs / kubectl logs

Phase 2:
  collect logs with OTel Collector or Fluent Bit

Phase 3:
  route to Loki/OpenSearch/cloud logging
  dashboards and alerts

Phase 4:
  retention, redaction validation, incident workflows
```

---

## Alerting

Prefer metrics for primary alerting and logs for explanation.

Example:

```text
metric alert:
  payment failure rate > threshold

log query:
  payment_authorisation_failed grouped by provider/error_code
```

---

## Enrichment

The platform may enrich logs with:

```text
k8s.namespace.name
k8s.pod.name
k8s.deployment.name
container.name
service.version
deployment.environment.name
```

Application logs should already include service and correlation fields.

---

## Practical Rules

```text
Keep apps vendor-neutral.
Do not hard-code log destinations in services.
Collect stdout/stderr.
Preserve structured JSON.
Enrich logs with platform metadata.
Use retention deliberately.
Protect sensitive fields.
```

---

## Final Rule

```text
The log pipeline turns service event streams into searchable operational memory.
```
