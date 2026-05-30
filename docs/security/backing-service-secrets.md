# Backing Service Secrets

This document defines how bfstore handles secrets used to connect to backing services.

It complements:

```text
docs/architecture/backing-services.md
docs/deployment/resource-attachments.md
docs/security/secrets-management.md
```

---

## Purpose

Backing services often require credentials.

This document explains how bfstore handles:

```text
database credentials
Kafka credentials
SMTP credentials
payment provider keys
OTLP exporter keys
secret rotation later
```

---

## Core Rule

```text
Resource handles are config.
Credentials are secret config.
Both stay out of code.
```

---

## Common Backing Service Secrets

Examples:

```text
MYSQL_PASSWORD
MYSQL_DSN containing password
KAFKA_USERNAME
KAFKA_PASSWORD
SMTP_USERNAME
SMTP_PASSWORD
PAYMENT_PROVIDER_API_KEY
OTEL_EXPORTER_API_KEY
TLS_PRIVATE_KEY
```

---

## Where Secrets Live

Local:

```text
ignored .env.local files
fake/simulated credentials where possible
```

Kubernetes:

```text
Kubernetes Secrets
mounted secret files
environment variables from secrets
```

Cloud production-style later:

```text
AWS Secrets Manager
Azure Key Vault
Google Secret Manager
External Secrets Operator
SOPS-encrypted files
```

---

## What Must Never Be Committed

Do not commit:

```text
real database passwords
real DSNs with passwords
payment provider API keys
SMTP passwords
OTLP exporter API keys
TLS private keys
cloud access keys
JWT signing keys
```

Commit only safe examples:

```text
secret.example.yaml
.env.example
.env.local.example
```

---

## Injection Pattern

Applications should consume secrets through config interfaces.

Example:

```yaml
env:
  - name: ORDER_MYSQL_DSN
    valueFrom:
      secretKeyRef:
        name: order-db
        key: dsn
```

The application should not know whether the secret came from Kubernetes Secret, cloud secret manager, or local environment.

---

## Rotation Expectations Later

Secrets should eventually have rotation guidance.

Examples:

```text
database passwords
Kafka credentials
SMTP credentials
payment provider keys
OTLP exporter tokens
TLS certificates
```

Rotation should avoid code changes.

---

## Logging

Never log:

```text
full DSNs containing passwords
API keys
bearer tokens
SMTP passwords
private keys
authorisation headers
```

Redact secrets in diagnostics.

---

## Practical Rules

```text
Keep backing service credentials out of code.
Keep real secrets out of plain Git.
Use fake credentials locally where possible.
Inject secrets through environment variables or mounted files.
Use Kubernetes Secrets or cloud secret managers.
Do not log secrets.
Rotate exposed secrets immediately.
```

---

## Final Rule

```text
A backing service secret should be replaceable without rewriting or rebuilding the app.
```
