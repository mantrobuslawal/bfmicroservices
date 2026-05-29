# Collector Configuration Security

This document defines secure OpenTelemetry Collector configuration principles for **bfstore**.

Collector hosting security defines where and how the Collector runs. Collector configuration security defines what the Collector is allowed to do.

---

## Purpose

This document defines bfstore policy for:

```text
secure Collector configuration
secrets handling
encryption and authentication
component minimisation
receiver binding
permissions
redaction
resource utilisation
local and Kubernetes configuration examples
```

---

## Core Rule

```text
Collector config is security architecture written in YAML.
```

A Collector can be securely hosted but dangerously configured. Both hosting and configuration must be treated as platform concerns.

---

## Secrets Policy

Collector YAML must not contain real secrets.

Do not commit:

```text
API keys
bearer tokens
TLS private keys
vendor credentials
cloud credentials
production endpoints with embedded credentials
```

Use:

```text
environment variables
Kubernetes Secrets
cloud secret managers
CI/CD secret stores
ignored local .env files
```

Safe exporter example:

```yaml
exporters:
  otlp/vendor:
    endpoint: ${env:BFSTORE_OTEL_EXPORTER_ENDPOINT}
    headers:
      api-key: ${env:BFSTORE_OTEL_EXPORTER_API_KEY}
```

---

## Encryption and Authentication

Use TLS/authentication where telemetry crosses trust boundaries.

Local development may use internal Docker networking and debug exporters.

Kubernetes and production-style environments should consider:

```text
TLS to telemetry backends
authenticated exporter connections
authenticated receiver endpoints where exposed beyond trusted workloads
mTLS for stronger workload identity later
```

TLS material should come from Secrets and be mounted read-only.

---

## Component Minimisation

Enable only the Collector components bfstore needs.

Initial local Collector:

```yaml
receivers:
  otlp:

processors:
  batch:

exporters:
  debug:
```

Later additions must have a documented reason:

```text
memory_limiter
tail_sampling
prometheus exporter
otlp exporter
health_check extension
redaction processor
```

Rule:

```text
Every enabled component is part of the attack surface.
```

---

## Receiver Binding Policy

Avoid casual broad binding.

Risky:

```yaml
endpoint: 0.0.0.0:4317
```

Prefer environment-appropriate binding:

```text
local-only:
  127.0.0.1

Docker Compose:
  service name such as otel-collector:4317

Kubernetes:
  pod IP / ClusterIP service / internal-only receiver
```

Rule:

```text
Receivers are ingestion doors. Open only the doors bfstore intends to guard.
```

---

## Exporter Policy

Exporter endpoints must be deliberate and documented.

Avoid:

```text
unknown external destinations
hard-coded credentials
vendor-specific assumptions in service code
unreviewed exporters
```

Prefer:

```text
OTLP from services to Collector
Collector routes to approved backends
credentials from Secrets
backend choices documented
```

---

## Queue and Buffer Policy

Queue and buffer settings must be tuned carefully.

Avoid huge queue values without memory planning.

Risky:

```yaml
sending_queue:
  queue_size: 1000000
```

Safer mature controls:

```yaml
processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 400

  batch:
    timeout: 5s
```

Rule:

```text
Tune Collector queues and buffers with evidence, not vibes.
```

---

## Permission Policy

The Collector should run with least privilege.

A gateway-style Collector should not need:

```text
root user
privileged container
host filesystem access
cluster-admin RBAC
host network
host PID namespace
```

Kubernetes security intent:

```yaml
securityContext:
  runAsNonRoot: true
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
```

Components that perform discovery, such as observers, may require extra permissions and must be reviewed deliberately.

---

## Redaction Policy

First rule:

```text
Do not emit sensitive data in the first place.
```

Never emit:

```text
raw JWTs
passwords
API keys
payment card numbers
CVV
full shipping address
customer email
full basket JSON
raw Kafka payloads
full SQL with customer values
```

Redaction processors may be used as a safety net, especially in mature environments.

---

## Resource Utilisation Policy

Collector config must protect the Collector from excessive CPU and memory use.

Watch carefully when enabling:

```text
tail_sampling
large queues
large batch sizes
many exporters
many receivers
heavy redaction/transformation rules
logs ingestion
```

Monitor:

```text
accepted spans/sec
refused spans/sec
dropped spans
exporter failures
queue size
memory usage
CPU usage
processor latency
```

---

## Local Development Policy

Local Collector should be small and safe.

Recommended:

```text
OTLP receiver
batch processor
debug exporter
no real vendor credentials
internal Docker network
read-only config mount
version-pinned image where practical
```

---

## Kubernetes Policy

Kubernetes Collector config should use:

```text
ConfigMap for non-secret config
Secret for credentials and TLS material
ClusterIP service
NetworkPolicy restrictions
non-root security context
batch and memory_limiter processors
redaction where needed
internal metrics monitoring
```

---

## What Not To Do

Avoid:

```text
hard-coded API keys
public OTLP receivers
casual 0.0.0.0 binding
unused receivers/exporters
root Collector containers
cluster-admin permissions
huge queues without memory planning
redaction as the only protection
latest image tags in production-style environments
```

---

## Practical Rules

```text
Keep Collector config free of hard-coded secrets.
Use secret stores for sensitive values.
Use TLS/authentication across trust boundaries.
Enable only required components.
Avoid broad receiver binding.
Prefer internal-only receivers.
Use least privilege.
Use redaction as a safety net.
Protect Collector CPU and memory.
Monitor Collector internal telemetry.
Pin image versions before production-style use.
Document configuration decisions.
```

---

## Final Rule

```text
Make telemetry work without turning the telemetry system into the incident.
```
