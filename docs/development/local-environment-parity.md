# Local Environment Parity

This document defines how bfstore's local environment should mirror the production-style architecture at a smaller scale.

---

## Purpose

This document explains:

```text
Docker Compose stack
required local backing services
local config shape
local smoke testing
developer workflow
```

---

## Core Rule

```text
Docker Compose should be a small model of the real platform, not a parallel universe.
```

---

## Required Local Stack

The local development environment should include:

```text
api-gateway
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-worker
MySQL
Kafka
OpenTelemetry Collector
payment simulator
fake SMTP later
```

---

## Backing Services

Use the same backing service types locally where behaviour matters.

```text
MySQL for service-owned data
Kafka for events
OpenTelemetry Collector for telemetry
fake SMTP later for email
payment simulator for payment boundary
```

---

## Config Shape

Local config should use the same environment variable names as deployed environments.

Examples:

```text
CATALOG_MYSQL_DSN
ORDER_KAFKA_BROKERS
OTEL_EXPORTER_OTLP_ENDPOINT
PAYMENT_PROVIDER
PAYMENT_TIMEOUT_MS
CATALOG_SERVICE_ADDR
```

Values may differ, but names and meaning should stay consistent.

---

## Smoke Testing

Local smoke tests should verify:

```text
api-gateway starts
gRPC services start
MySQL migrations apply
Kafka topics exist
services can connect to dependencies
OpenTelemetry Collector receives telemetry
core browse -> basket -> checkout path works eventually
```

---

## Developer Workflow

Recommended local flow:

```bash
make bootstrap
make up
make migrate-up
make seed
make test
make smoke
```

Actual commands should evolve with the repository.

---

## Practical Rules

```text
Use MySQL locally, not SQLite.
Use Kafka locally, not an in-memory event bus.
Use gRPC locally, not direct package calls.
Use Protobuf locally, not JSON-only shortcuts.
Use OTel locally, even with simple exporters.
Keep local secrets fake and safe.
Document all local ports and service addresses.
```

---

## Final Rule

```text
Local development should teach the same architecture that staging and production-style environments run.
```
