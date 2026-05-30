# Configuration

This document defines bfstore's runtime configuration strategy.

It is based on the 12-Factor Config principle:

```text
Store deploy-varying config outside code.
```

---

## Purpose

This document defines:

```text
what config is
what config is not
environment variables
granular config
config vs code
config vs secrets
bfstore naming policy
validation expectations
```

---

## Core Rule

```text
Code defines what the app does.
Config defines where and how this deploy runs.
```

The same service artefact should run in multiple environments with different configuration.

---

## What Counts as Config

Config is anything likely to vary between deploys.

Examples:

```text
BFSTORE_ENV
GRPC_PORT
MYSQL_DSN
KAFKA_BROKERS
OTEL_EXPORTER_OTLP_ENDPOINT
LOG_LEVEL
PAYMENT_PROVIDER
PAYMENT_TIMEOUT_MS
ORDER_CHECKOUT_DEADLINE_MS
```

Secrets are also config:

```text
MYSQL_PASSWORD
JWT_SIGNING_KEY
PAYMENT_PROVIDER_API_KEY
OTEL_EXPORTER_API_KEY
TLS_PRIVATE_KEY
```

---

## What Is Not Deploy Config

The following belong in code:

```text
gRPC service definitions
route registration
dependency injection wiring
business rules
checkout orchestration flow
Protobuf message definitions
validation rules
```

Rule:

```text
Business behaviour belongs in code.
Deploy-specific values belong in config.
```

---

## Environment Variables

bfstore services should consume runtime config through environment variables.

Example:

```bash
BFSTORE_ENV=local
GRPC_PORT=50051
MYSQL_DSN=catalog_user:catalog_password@tcp(mysql:3306)/catalog_db
KAFKA_BROKERS=kafka:9092
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
LOG_LEVEL=debug
```

---

## Granular Config

Do not use `BFSTORE_ENV` as a giant switch statement.

Bad:

```text
if BFSTORE_ENV == production:
  set database, Kafka, log level, payment provider, timeout
```

Better:

```text
BFSTORE_ENV=staging
MYSQL_DSN=...
KAFKA_BROKERS=...
LOG_LEVEL=info
PAYMENT_PROVIDER=simulated
```

Rule:

```text
Use BFSTORE_ENV for deploy labelling, not hidden configuration.
```

---

## Config vs Secrets

Non-secret config may live in:

```text
.env.example
ConfigMaps
Helm values
Kustomize overlays
Terraform variables
```

Secrets should live in:

```text
Kubernetes Secrets
External Secrets Operator
AWS Secrets Manager
Azure Key Vault
Google Secret Manager
SOPS-encrypted files
CI/CD secret store
```

Rule:

```text
Config belongs outside code.
Secrets belong outside code and outside plain Git.
```

---

## Naming Policy

Global values:

```text
BFSTORE_ENV
LOG_LEVEL
OTEL_EXPORTER_OTLP_ENDPOINT
OTEL_RESOURCE_ATTRIBUTES
```

Service-specific values:

```text
CATALOG_GRPC_PORT
CATALOG_MYSQL_DSN
ORDER_GRPC_PORT
ORDER_MYSQL_DSN
ORDER_KAFKA_BROKERS
PAYMENT_PROVIDER
PAYMENT_TIMEOUT_MS
```

Dependency addresses:

```text
CATALOG_SERVICE_ADDR
INVENTORY_SERVICE_ADDR
PAYMENT_SERVICE_ADDR
SHIPPING_SERVICE_ADDR
```

Secrets:

```text
MYSQL_PASSWORD
JWT_SIGNING_KEY
PAYMENT_PROVIDER_API_KEY
```

Rule:

```text
Config names should be boring, predictable, and documented.
```

---

## Validation

Services must validate config at startup.

Validate:

```text
required env vars are present
ports are valid
URLs are valid
timeouts parse correctly
ratios are within allowed ranges
provider names are recognised
broker lists are not empty
```

Rule:

```text
Invalid config should fail fast before serving traffic.
```

---

## Logging

Do not log secret config values.

Safe:

```text
configured environment name
configured service name
configured non-secret port
whether optional config is enabled
```

Unsafe:

```text
database passwords
API keys
JWT signing keys
full DSNs containing passwords
bearer tokens
```

---

## Testing Expectations

Tests should cover:

```text
required values missing
defaults applied
invalid values rejected
secret values redacted
OpenTelemetry resource attributes populated
local examples work
Kubernetes examples render
```

---

## Practical Rules

```text
Store deploy-varying config outside code.
Use environment variables as the application-facing interface.
Keep real secrets out of Git.
Use granular config values.
Validate config at startup.
Fail fast on missing required config.
Do not log secrets.
Keep business rules in code.
Use consistent names.
Document required variables.
```

---

## Final Rule

```text
The same artefact should run everywhere; configuration tells it where it is.
```
