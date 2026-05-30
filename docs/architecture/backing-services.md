# Backing Services

This document defines how **bfstore** treats backing services as attached resources.

A backing service is any external network service that a bfstore app consumes during normal operation.

---

## Purpose

This document explains:

```text
what backing services are
the attached resource model
local vs managed resources
service-owned databases
Kafka and OpenTelemetry as backing services
service-to-service endpoints
replaceability expectations
```

---

## Core Rule

```text
The app should depend on a capability, not a hard-coded place.
```

A service should access backing services through configuration, not source-code constants.

---

## Backing Services in bfstore

Backing services include:

```text
MySQL
Kafka
OpenTelemetry Collector
SMTP/email provider
payment provider
object storage later
Redis later
search engine later
```

---

## Attached Resource Model

A backing service is attached to a deploy through config.

Example:

```text
order-service:
  ORDER_MYSQL_DSN
  ORDER_KAFKA_BROKERS
  INVENTORY_SERVICE_ADDR
  PAYMENT_SERVICE_ADDR
  SHIPPING_SERVICE_ADDR
  OTEL_EXPORTER_OTLP_ENDPOINT
```

Changing the resource handle should rewire the app without changing code.

---

## Backing Service vs Dependency vs Config

Example:

```text
Dependency:
  github.com/go-sql-driver/mysql

Backing service:
  MySQL database

Config:
  MYSQL_DSN
```

Rule:

```text
The driver is a dependency.
The database is a backing service.
The DSN is config.
```

---

## Local vs Managed Resources

Local and managed services should look like the same capability to the app.

Example:

```text
local:
  CATALOG_MYSQL_DSN=catalog_user:catalog_password@tcp(mysql:3306)/catalog_db

production-style:
  CATALOG_MYSQL_DSN=catalog_user:<secret>@tcp(managed-mysql.example:3306)/catalog_db
```

The `catalog-service` code should not change.

---

## Service-owned Databases

Each service should own its data resource.

```text
catalog-service -> catalog database
basket-service -> basket database
order-service -> order database
payment-service -> payment database
shipping-service -> shipping database
```

Avoid shared-database coupling.

Rule:

```text
Service-owned data means service-owned backing resources.
```

---

## Kafka as a Backing Service

Kafka is an attached resource.

Services should receive broker details through config:

```text
ORDER_KAFKA_BROKERS
NOTIFICATION_KAFKA_BROKERS
PAYMENT_KAFKA_BROKERS
```

The event contracts belong in Protobuf and docs.

The broker location belongs in config.

---

## OpenTelemetry Collector as a Backing Service

The OpenTelemetry Collector is an attached telemetry resource.

Services should receive:

```text
OTEL_EXPORTER_OTLP_ENDPOINT
```

Instrumentation code should remain portable across local Collector, Kubernetes Collector, and future managed backends.

---

## Service-to-service Endpoints

Service endpoints are network resources from the caller’s perspective.

Examples:

```text
CATALOG_SERVICE_ADDR
INVENTORY_SERVICE_ADDR
PAYMENT_SERVICE_ADDR
SHIPPING_SERVICE_ADDR
```

The API contract belongs in code/protobuf.

The network address belongs in config.

---

## Replaceability

Backing services should be replaceable by changing configuration.

Examples:

```text
local Kafka -> Strimzi Kafka
MailHog -> managed SMTP provider
local MySQL -> managed MySQL
local Collector -> Kubernetes Collector
```

No application code should change.

---

## What Not To Do

Avoid:

```text
hard-coded hosts
hard-coded brokers
hard-coded provider endpoints
services accessing another service’s database
environment names controlling infrastructure decisions in code
credentials committed to Git
```

---

## Practical Rules

```text
Treat network dependencies as backing services.
Inject resource handles through config.
Keep credentials in secrets.
Keep application code provider-neutral where practical.
Keep service-owned resources separate.
Make service endpoints configurable.
Make replacements possible without code changes.
```

---

## Final Rule

```text
Backing services are attached to deploys; they are not baked into the application.
```
