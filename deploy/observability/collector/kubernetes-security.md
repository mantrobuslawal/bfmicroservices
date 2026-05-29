# Collector Kubernetes Security

This document defines Kubernetes-specific security expectations for running the OpenTelemetry Collector in **bfstore**.

It complements:

```text
docs/architecture/collector-hosting-security.md
deploy/observability/collector/security.md
```

---

## Purpose

When bfstore runs on Kubernetes, the Collector should be deployed with least privilege, internal network access, controlled secrets, and clear resource limits.

---

## Deployment Pattern

Default recommended pattern:

```text
Collector as Deployment
ClusterIP Service
ConfigMap for non-secret config
Secret for credentials/TLS
minimal ServiceAccount
NetworkPolicy restricting ingress
```

Use DaemonSet only when intentionally collecting node-level telemetry such as host logs or node metrics.

---

## ServiceAccount and RBAC

Use a dedicated ServiceAccount:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: otel-collector
  namespace: bfstore-observability
```

Grant only required permissions.

Avoid:

```text
cluster-admin
wildcard verbs
wildcard resources
unnecessary node access
```

---

## NetworkPolicy Expectations

The Collector receiver should only be reachable by trusted workloads.

Policy intent:

```text
allow bfstore services -> otel-collector
deny public internet -> otel-collector
deny unrelated namespaces -> otel-collector
```

Example future file:

```text
deploy/observability/collector/networkpolicy.yaml
```

---

## Service Exposure

Default service type:

```text
ClusterIP
```

Avoid by default:

```text
LoadBalancer
NodePort
public ingress
```

OTLP receivers should not be public unless there is a deliberate authenticated design.

---

## ConfigMap vs Secret

Use ConfigMap for:

```text
receivers
processors
exporters without credentials
extensions
pipeline structure
```

Use Secret for:

```text
API tokens
bearer tokens
TLS certs
TLS keys
CA bundles
vendor credentials
```

Secrets should be mounted read-only where practical.

---

## Security Context

Use hardened container settings where possible:

```yaml
securityContext:
  runAsNonRoot: true
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
```

Pod-level settings may include:

```text
non-root user
dropped Linux capabilities
seccomp profile
read-only filesystem
```

---

## hostPath Policy

Avoid `hostPath` mounts unless required.

Gateway-style Collector:

```text
should not need hostPath mounts
```

DaemonSet node log Collector:

```text
may need specific read-only hostPath mounts
```

Bad:

```text
mount entire host filesystem
run privileged by default
write access to host logs
```

Better:

```text
specific read-only mounts
documented reason
minimal RBAC
separate DaemonSet if needed
```

---

## Resource Requests and Limits

Set requests and limits.

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

Tune values through load testing and Collector internal telemetry.

---

## Scaling and Availability

For production-style environments, consider:

```text
multiple Collector replicas
horizontal scaling
PodDisruptionBudget
readiness/liveness probes
internal telemetry alerts
```

Tail sampling may require special scaling considerations because traces must be grouped for decisions.

---

## Collector Internal Metrics

Scrape and alert on:

```text
CPU
memory
receiver accepted/refused telemetry
processor dropped telemetry
exporter send failures
queue size
batch processor behaviour
tail sampling processor health
```

---

## Image Policy

Production-style deployments should:

```text
pin image versions
avoid latest
scan images
generate SBOMs where appropriate
sign custom images if built
```

---

## Practical Rules

```text
Use ClusterIP by default.
Restrict receiver access with NetworkPolicy.
Use a dedicated ServiceAccount.
Avoid cluster-admin.
Use ConfigMaps for config and Secrets for secrets.
Mount secrets read-only.
Avoid hostPath unless deliberately collecting node-level data.
Use securityContext hardening.
Set resource requests and limits.
Monitor the Collector.
Pin images before production-style use.
```

---

## Final Rule

```text
In Kubernetes, the Collector should be internal, least-privileged, monitored, and deliberately configured.
```
