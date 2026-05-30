# Secrets Management

This document defines bfstore's secrets management expectations.

It complements:

```text
docs/architecture/configuration.md
docs/development/local-configuration.md
docs/deployment/configuration-management.md
docs/security/dependency-management.md
```

---

## Purpose

This document explains:

```text
what counts as a secret
where secrets live
what must never be committed
rotation expectations later
External Secrets/cloud secret manager options
```

---

## Core Rule

```text
Secrets belong outside code and outside plain Git.
```

---

## What Counts as a Secret

Examples:

```text
database passwords
full DSNs containing passwords
JWT signing keys
payment provider API keys
OTLP exporter API keys
TLS private keys
Kafka passwords
cloud access keys
bearer tokens
session signing keys
```

If it authenticates, decrypts, signs, or grants access, treat it as a secret.

---

## What Must Never Be Committed

Do not commit:

```text
real .env files
production credentials
private keys
API tokens
cloud credentials
database passwords
JWT signing keys
payment provider secrets
```

Commit only safe examples:

```text
.env.example
.env.local.example
secret.example.yaml
```

---

## Local Secrets

Local secrets may live in ignored files such as:

```text
.env.local
.env
```

These must be listed in `.gitignore`.

Local development should use fake/simulated credentials where possible.

---

## Kubernetes Secrets

Kubernetes may inject secrets as environment variables or mounted files.

Example:

```yaml
env:
  - name: CATALOG_MYSQL_DSN
    valueFrom:
      secretKeyRef:
        name: catalog-db
        key: dsn
```

TLS material should be mounted read-only where practical.

---

## Cloud Secret Managers

Later production-style environments may use:

```text
AWS Secrets Manager
Azure Key Vault
Google Secret Manager
External Secrets Operator
SOPS-encrypted files
```

Recommended later pattern:

```text
cloud secret manager
  -> ExternalSecret
  -> Kubernetes Secret
  -> service environment variable or mounted file
```

---

## Rotation

Later, secrets should have documented rotation expectations.

Examples:

```text
database passwords
JWT signing keys
payment provider API keys
OTLP exporter tokens
TLS certificates
```

Rotation should avoid unnecessary downtime where possible.

---

## Logging and Redaction

Services must not log secrets.

Avoid logging:

```text
full DSNs
API keys
bearer tokens
JWTs
private keys
passwords
authorisation headers
```

Config logs should redact sensitive values.

---

## Secret Scanning

CI should eventually include secret scanning.

Possible tools:

```text
gitleaks
trufflehog
GitHub secret scanning
```

If a secret is committed:

```text
revoke it
rotate it
remove it from history where appropriate
document the incident if needed
```

---

## Practical Rules

```text
Never commit real secrets.
Use safe example files.
Keep local secrets ignored.
Use Kubernetes Secrets for sensitive config.
Use cloud secret managers later.
Do not log secrets.
Redact sensitive config in diagnostics.
Scan for committed secrets.
Rotate exposed secrets immediately.
```

---

## Final Rule

```text
A secret in Git is no longer a secret.
```
