# Local Configuration

This document explains how bfstore should manage local development configuration.

It complements:

```text
docs/architecture/configuration.md
docs/deployment/configuration-management.md
docs/security/secrets-management.md
```

---

## Purpose

Local configuration should make bfstore easy to run without committing secrets or hard-coding environment-specific values.

---

## Core Rule

```text
Local config should be easy to create, safe to share, and impossible to confuse with production secrets.
```

---

## Files

Recommended files:

```text
.env.example
.env.local.example
.env.local
.gitignore
```

Commit:

```text
.env.example
.env.local.example
```

Do not commit:

```text
.env
.env.local
.env.*.local
```

---

## .env.example

`.env.example` documents required variables with safe placeholder values.

Example:

```bash
BFSTORE_ENV=local
LOG_LEVEL=debug
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
```

---

## Service Local Config

Example catalogue config:

```bash
CATALOG_GRPC_PORT=50051
CATALOG_MYSQL_DSN=catalog_user:catalog_password@tcp(mysql:3306)/catalog_db
```

Example order config:

```bash
ORDER_GRPC_PORT=50052
ORDER_MYSQL_DSN=order_user:order_password@tcp(mysql:3306)/order_db
ORDER_KAFKA_BROKERS=kafka:9092
INVENTORY_SERVICE_ADDR=inventory-service:50053
PAYMENT_SERVICE_ADDR=payment-service:50054
SHIPPING_SERVICE_ADDR=shipping-service:50055
```

Example payment config:

```bash
PAYMENT_PROVIDER=simulated
PAYMENT_TIMEOUT_MS=2000
PAYMENT_SIMULATED_FAILURE_RATE=0.05
```

---

## Docker Compose

Docker Compose may inject config through:

```text
environment:
env_file:
```

Example:

```yaml
services:
  catalog-service:
    environment:
      BFSTORE_ENV: local
      CATALOG_GRPC_PORT: "50051"
      CATALOG_MYSQL_DSN: catalog_user:catalog_password@tcp(mysql:3306)/catalog_db
      OTEL_EXPORTER_OTLP_ENDPOINT: http://otel-collector:4317
```

Local Compose config should use safe local credentials only.

---

## make doctor / config-check

A future `make doctor` or `make config-check` target should verify:

```text
required local env files exist
required variables are documented
Docker Compose config renders
ports are not obviously conflicting
secrets are not committed
```

---

## Safe Defaults

Safe local defaults may include:

```text
LOG_LEVEL=debug
PAYMENT_PROVIDER=simulated
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
```

Do not use production credentials locally.

---

## Practical Rules

```text
Commit examples, not secrets.
Keep local .env files ignored.
Use Docker Compose to inject local config.
Use simulated providers locally.
Fail fast if required config is missing.
Keep local config boring and reproducible.
```

---

## Final Rule

```text
Local configuration should help future you run bfstore quickly without secret-leak anxiety.
```
