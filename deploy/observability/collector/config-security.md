# Collector Config Security

This document provides practical secure configuration guidance for bfstore OpenTelemetry Collector YAML.

It complements:

```text
docs/architecture/collector-configuration-security.md
deploy/observability/collector/security.md
deploy/observability/collector/kubernetes-security.md
deploy/observability/collector/redaction-policy.md
```

---

## Purpose

This document defines safe Collector YAML patterns for bfstore.

It covers:

```text
environment variable usage
ConfigMap vs Secret split
safe receiver endpoints
safe exporter patterns
queue and buffer tuning
component minimisation
resource protection
```

---

## ConfigMap vs Secret Split

Use ConfigMaps for non-secret Collector config:

```text
receivers
processors
exporters without credentials
extensions
pipeline definitions
```

Use Secrets for:

```text
API keys
bearer tokens
TLS private keys
CA certificates
vendor credentials
cloud credentials
```

---

## Environment Variable Usage

Collector config may reference environment variables for sensitive values.

Example:

```yaml
exporters:
  otlp/vendor:
    endpoint: ${env:BFSTORE_OTEL_EXPORTER_ENDPOINT}
    headers:
      api-key: ${env:BFSTORE_OTEL_EXPORTER_API_KEY}
```

Values should be injected from:

```text
Kubernetes Secrets
cloud secret manager
CI/CD secret store
ignored local .env
```

---

## Safe Local Collector Pattern

Recommended local shape:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: otel-collector:4317
      http:
        endpoint: otel-collector:4318

processors:
  batch:

exporters:
  debug:

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]

    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug]
```

Local rules:

```text
no real vendor credentials
debug exporter only
internal Docker network
config mounted read-only
```

---

## Receiver Endpoint Rules

Avoid broad binding unless required and protected.

Risky:

```yaml
endpoint: 0.0.0.0:4317
```

Prefer:

```text
127.0.0.1 for local-only
otel-collector:4317 for Docker Compose network
pod IP / ClusterIP for Kubernetes
```

Production-style environments should not expose OTLP receivers publicly.

---

## Exporter Rules

Exporter configuration must be intentional.

Good:

```text
approved backend
documented endpoint
credentials from Secret
TLS where appropriate
```

Bad:

```text
hard-coded token
unknown external endpoint
vendor lock-in hidden in service code
unreviewed exporter
```

---

## Queue and Buffer Rules

Do not increase queue sizes blindly.

Bad:

```yaml
sending_queue:
  queue_size: 1000000
```

Better:

```yaml
processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 400

  batch:
    timeout: 5s
```

Queue and buffer settings should be based on testing, backend behaviour, and Collector internal telemetry.

---

## Component Minimisation

Start with:

```text
otlp receiver
batch processor
debug exporter
```

Add components only when needed:

```text
memory_limiter for resource protection
redaction for safety net filtering
tail_sampling for mature trace sampling
prometheus exporter for metrics backend integration
otlp exporter for trace/log backend forwarding
```

Do not enable unused receivers/exporters/extensions.

---

## Resource Protection

Mature Collector deployments should include:

```text
memory_limiter processor
batch processor
Kubernetes resource requests/limits
Collector internal telemetry alerts
```

Monitor:

```text
CPU
memory
accepted telemetry
refused telemetry
dropped telemetry
queue size
exporter failures
processor latency
```

---

## Testing Expectations

Test Collector config by checking:

```text
Collector starts successfully
config validates
OTLP gRPC receiver works
OTLP HTTP receiver works
debug exporter receives test spans
no real secrets exist in config
no unexpected receivers are enabled
resource limits exist in Kubernetes manifests
```

---

## Practical Rules

```text
Keep YAML small.
Keep secrets out of YAML.
Bind receivers deliberately.
Use approved exporters only.
Tune queues with evidence.
Use memory_limiter in mature configs.
Use redaction as a safety net.
Validate config in CI.
Document component choices.
```

---

## Final Rule

```text
Collector YAML should be boring, reviewable, and safe.
```
