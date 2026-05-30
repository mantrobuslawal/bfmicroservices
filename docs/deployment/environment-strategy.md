# Environment Strategy

This document defines how **bfstore** handles multiple deploy environments while preserving one source of truth for application code.

---

## Purpose

This document explains:

```text
local/dev/staging/prod deploys
same code, different config
environment variables
secrets
Docker Compose vs Kubernetes
build once, deploy many
```

---

## Core Rule

```text
Same artefact, different config.
```

The same service code should be deployable to many environments without copying or modifying source code per environment.

---

## Environments

bfstore may use:

```text
local
dev
test
staging
production-style
```

Each environment is a deploy target, not a separate codebase.

---

## Configuration Differences

Environment-specific values include:

```text
MySQL endpoint
Kafka broker endpoint
OTLP Collector endpoint
log level
feature flags
timeouts
replica count
CPU/memory resources
external provider endpoints
secrets
```

These should be supplied through configuration, not copied source code.

---

## Local Example

```text
BFSTORE_ENV=local
MYSQL_HOST=localhost
KAFKA_BROKERS=localhost:9092
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
PAYMENT_PROVIDER=simulated
```

---

## Dev Example

```text
BFSTORE_ENV=dev
MYSQL_HOST=dev-mysql.bfstore.internal
KAFKA_BROKERS=dev-kafka:9092
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector.dev:4317
PAYMENT_PROVIDER=simulated
```

---

## Production-style Example

```text
BFSTORE_ENV=prod
MYSQL_HOST=<managed-mysql-endpoint>
KAFKA_BROKERS=<managed-kafka-endpoint>
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector.prod:4317
PAYMENT_PROVIDER=<approved-provider>
```

Secrets should come from a secret manager or Kubernetes Secrets.

---

## Configuration Locations

Possible locations:

```text
.env.example
Docker Compose override files
Kubernetes Helm values
Kustomize overlays
Terraform variables
GitOps environment folders
Secret manager entries
```

Do not use:

```text
environment-specific source folders
environment-specific branches
copied service code
hard-coded production values
```

---

## Docker Compose

Local Docker Compose should support:

```text
fast local development
local MySQL
local Kafka
local OpenTelemetry Collector
debug exporters
safe defaults
```

No production secrets should be used locally.

---

## Kubernetes

Kubernetes environments may use:

```text
base manifests
environment overlays
ConfigMaps for non-secret config
Secrets for sensitive config
NetworkPolicies
resource requests/limits
horizontal scaling
```

Example structure:

```text
deploy/kubernetes/
├── base/
└── overlays/
    ├── local/
    ├── dev/
    ├── staging/
    └── prod/
```

---

## Build Once, Deploy Many

Recommended flow:

```text
commit
  -> test
  -> build image
  -> tag image with commit SHA
  -> deploy to dev
  -> promote same image to staging
  -> promote same image to prod
```

Avoid rebuilding different images from different branches for each environment.

---

## Secrets

Secrets must not be committed to Git.

Use:

```text
Kubernetes Secrets
cloud secret manager
CI/CD secret store
ignored local .env files
```

Examples:

```text
database password
Kafka credentials
OTLP exporter token
TLS private key
payment provider credentials
```

---

## Traceability

Each deploy should record:

```text
service name
environment
git commit
image tag
release version
deployment time
```

Rule:

```text
If you cannot trace a running deploy back to a commit, the delivery story is weak.
```

---

## Practical Rules

```text
Use config for environment differences.
Never copy source code per environment.
Never use long-lived environment branches.
Build once, deploy many.
Tag images with commit SHA.
Keep secrets out of Git.
Use ConfigMaps for non-secret config.
Use Secrets for sensitive config.
Make each deploy traceable.
```

---

## Final Rule

```text
Environments change configuration, not the application codebase.
```
