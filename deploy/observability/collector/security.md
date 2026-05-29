# Collector Security

This document defines the practical security rules for running the OpenTelemetry Collector in **bfstore**.

It complements:

```text
docs/architecture/collector-hosting-security.md
deploy/observability/collector/README.md
deploy/observability/collector/kubernetes-security.md
```

---

## Purpose

The Collector receives telemetry from bfstore services and exports it to observability backends.

It should be configured securely because it may handle:

```text
operational telemetry
exporter credentials
TLS certificates
authentication headers
high-volume traffic
```

---

## Safe Configuration Layout

Recommended directory:

```text
deploy/observability/collector/
├── README.md
├── security.md
├── kubernetes-security.md
├── otel-collector.local.yaml
├── otel-collector.k8s.yaml
└── otel-collector.example.yaml
```

Config files may define:

```text
receivers
processors
exporters
extensions
pipelines
```

But must not include real secrets.

---

## Secrets Handling

Do not commit:

```text
API keys
TLS private keys
bearer tokens
vendor credentials
cloud credentials
```

Use:

```text
Kubernetes Secrets
cloud secret manager
ignored local .env files
CI/CD secret store
```

Example Kubernetes environment variable:

```yaml
env:
  - name: OTEL_EXPORTER_OTLP_HEADERS
    valueFrom:
      secretKeyRef:
        name: otel-exporter-credentials
        key: headers
```

Example TLS mount:

```yaml
volumeMounts:
  - name: otel-collector-tls
    mountPath: /etc/otel/tls
    readOnly: true

volumes:
  - name: otel-collector-tls
    secret:
      secretName: otel-collector-tls
```

---

## Receiver Security

Default OTLP ports:

```text
4317 gRPC
4318 HTTP
```

Default policy:

```text
internal only
not publicly exposed
restricted to trusted bfstore services
```

Local development may expose ports for debugging, but production-style environments should not expose Collector receivers publicly.

---

## Exporter Security

Exporter credentials must come from secure sources.

Exporter endpoints should be deliberate and documented.

Avoid:

```text
unknown vendor endpoints
unreviewed external exporters
hard-coded credentials
plain-text secrets in ConfigMaps
```

---

## Resource Controls

Set resource requests and limits.

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

Mature Collector configs should consider:

```text
batch processor
memory_limiter processor
queue monitoring
exporter retry policy
```

---

## Internal Monitoring

Monitor:

```text
Collector CPU
Collector memory
received telemetry rate
exported telemetry rate
dropped telemetry
queue length
exporter failures
processor latency
receiver errors
```

Alert on:

```text
high memory
dropped spans
exporter failures
queue saturation
receiver errors
```

---

## Local Development Rules

Local Collector:

```text
may use debug exporter
must not use real production credentials
should use version-pinned image where practical
should mount config read-only
should run with minimal privileges
```

---

## Production-style Rules

Production-style Collector:

```text
internal receiver access only
least-privilege ServiceAccount/RBAC
Secrets for credentials
ConfigMap for non-secret config
resource requests/limits
internal telemetry monitored
version-pinned image
NetworkPolicy or equivalent restrictions
```

---

## What Not To Do

Avoid:

```text
hard-coded exporter tokens
public OTLP receiver
cluster-admin permissions
privileged container by default
unnecessary hostPath mounts
latest image tags
no resource limits
no Collector monitoring
```

---

## Practical Rules

```text
Keep secrets out of Collector config.
Keep receiver access internal.
Use least privilege.
Pin images before production-style use.
Monitor the Collector.
Protect resource usage.
Document exporter choices.
Treat the Collector as a platform service.
```

---

## Final Rule

```text
The Collector is bfstore’s telemetry gateway; do not let the gateway become the weak point.
```
