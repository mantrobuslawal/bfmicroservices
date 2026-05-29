# Collector Hosting Security

This document defines the hosting security policy for the OpenTelemetry Collector in **bfstore**.

The OpenTelemetry Collector receives, processes, and exports telemetry. Because it sits in the middle of the observability pipeline, it must be treated as production infrastructure.

---

## Purpose

This document defines bfstore security expectations for:

```text
Collector configuration
Collector secrets
receiver exposure
least privilege
Kubernetes RBAC
network access
resource utilisation
Collector monitoring
```

---

## Core Rule

```text
The Collector is infrastructure, not just config.
```

It should be secured, monitored, versioned, and operated like any other platform service.

---

## Collector Data and Configuration

Collector configuration may be stored in Git when it does not contain secrets.

Allowed in Git:

```text
receivers
processors
exporters without credentials
extensions
pipeline structure
safe local examples
```

Not allowed in Git:

```text
real API keys
private TLS keys
production bearer tokens
vendor credentials
cloud credentials
```

Rule:

```text
Collector config can be version-controlled.
Collector secrets must not be hard-coded.
```

---

## Secrets Policy

Use secure secret storage for:

```text
exporter API tokens
TLS certificates
TLS private keys
CA certificates
bearer tokens
cloud credentials
```

Kubernetes examples:

```text
ConfigMap:
  non-secret Collector config

Secret:
  exporter token
  TLS material
```

Secrets should be mounted read-only where possible.

---

## Least Privilege

The Collector should only receive the permissions it needs.

For a gateway-style Collector receiving OTLP from services, it usually only needs to:

```text
listen on internal OTLP ports
read config
read mounted secrets
export telemetry to configured backend
```

It should not need:

```text
privileged mode
host filesystem access
cluster-admin permissions
```

Only add node-level privileges when intentionally collecting node-level data.

---

## Receiver Exposure Policy

Collector receivers are ingestion doors.

Default bfstore policy:

```text
OTLP receivers are internal only.
Do not expose OTLP receivers publicly.
Restrict access with network policy or equivalent controls.
Use auth/TLS where appropriate.
```

Allowed:

```text
bfstore services -> otel-collector:4317
```

Not allowed:

```text
public internet -> otel-collector
untrusted pod -> otel-collector
unknown namespace -> otel-collector
```

---

## Local Docker Compose Guidance

Local development may expose Collector ports for debugging, but should not use real secrets.

Recommended local behaviour:

```text
debug exporter enabled
no real vendor credentials
upstream Collector image
OTLP receiver inside Docker network
config mounted read-only
```

Example:

```yaml
otel-collector:
  image: otel/opentelemetry-collector-contrib:0.x.x
  command: ["--config=/etc/otelcol/config.yaml"]
  volumes:
    - ./deploy/observability/collector/otel-collector.local.yaml:/etc/otelcol/config.yaml:ro
```

---

## Kubernetes Guidance

Recommended Kubernetes shape:

```text
Deployment:
  Collector pod

Service:
  ClusterIP only

ConfigMap:
  non-secret Collector config

Secret:
  credentials and TLS material

ServiceAccount:
  minimal permissions

NetworkPolicy:
  restrict ingress to bfstore services

Prometheus/ServiceMonitor:
  scrape Collector internal metrics
```

---

## Resource Utilisation

The Collector must have resource requests and limits.

Example:

```yaml
resources:
  requests:
    cpu: "100m"
    memory: "256Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"
```

Mature environments should consider:

```text
memory_limiter processor
batch processor
horizontal scaling
queue monitoring
exporter failure alerts
```

---

## Collector Internal Telemetry

Monitor the Collector itself.

Important signals:

```text
CPU usage
memory usage
received spans/sec
exported spans/sec
dropped spans
queue length
exporter failures
processor latency
refused connections
```

Rule:

```text
The Collector needs observability too.
```

---

## Production-style Hardening

Consider:

```text
mTLS or authenticated ingestion
memory_limiter processor
horizontal scaling
pod disruption budget
securityContext
read-only root filesystem
non-root container
image pinning
image scanning
SBOM/signing for custom images
```

Example security intent:

```yaml
securityContext:
  runAsNonRoot: true
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
```

---

## What Not To Do

Avoid:

```text
committing Collector secrets to Git
exposing OTLP receiver publicly without auth
using cluster-admin permissions
mounting host paths unless required
running privileged by default
ignoring Collector CPU/memory
using latest image tags in production-style environments
putting raw customer data into telemetry
using Collector filtering as the only safety net
```

---

## Practical Rules

```text
Do not commit Collector secrets.
Use Secrets for credentials and TLS material.
Use ConfigMaps for non-secret Collector config.
Keep OTLP receivers internal by default.
Restrict Collector network access.
Use least-privilege RBAC.
Do not run privileged unless collecting privileged node data.
Avoid hostPath mounts unless truly needed.
Set CPU and memory requests/limits.
Monitor Collector internal telemetry.
Alert on Collector resource exhaustion and dropped data.
Scale Collectors horizontally when needed.
Use memory_limiter and batch processors in mature environments.
Pin Collector image versions before production-style use.
Treat the Collector as production infrastructure.
```

---

## Final Rule

```text
Secure the Collector like a gateway, monitor it like a service, and configure it like infrastructure.
```
