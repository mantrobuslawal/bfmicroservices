# Local Backing Services

This document defines local backing services for bfstore development.

It complements:

```text
docs/architecture/backing-services.md
docs/deployment/resource-attachments.md
docs/security/backing-service-secrets.md
```

---

## Purpose

Local backing services should make bfstore easy to run, test, and understand without relying on production infrastructure.

---

## Core Rule

```text
Local services should provide the same capability interface as managed services.
```

The implementation may differ, but the application-facing config should remain stable.

---

## Local Backing Services

Recommended local backing services:

```text
MySQL
Kafka
OpenTelemetry Collector
MailHog or fake SMTP later
payment simulator
object storage emulator later
Redis later if caching is added
```

---

## Docker Compose

Docker Compose should document local resource attachments.

Example:

```yaml
services:
  catalog-service:
    environment:
      CATALOG_MYSQL_DSN: catalog_user:catalog_password@tcp(mysql:3306)/catalog_db
      OTEL_EXPORTER_OTLP_ENDPOINT: http://otel-collector:4317
    depends_on:
      - mysql
      - otel-collector

  order-service:
    environment:
      ORDER_MYSQL_DSN: order_user:order_password@tcp(mysql:3306)/order_db
      ORDER_KAFKA_BROKERS: kafka:9092
      INVENTORY_SERVICE_ADDR: inventory-service:50051
      PAYMENT_SERVICE_ADDR: payment-service:50051
      SHIPPING_SERVICE_ADDR: shipping-service:50051
```

---

## MySQL

Local MySQL may provide multiple service-owned databases.

Examples:

```text
catalog_db
basket_db
inventory_db
order_db
payment_db
shipping_db
notification_db
```

Even if these databases run on one local MySQL container, services should treat their own DB/schema as their attached resource.

---

## Kafka

Local Kafka supports event-driven flows.

Example topics:

```text
bfstore.order.orders.v1
bfstore.inventory.stock.v1
bfstore.payment.payments.v1
bfstore.shipping.shipments.v1
bfstore.notification.notifications.v1
bfstore.catalog.products.v1
```

Kafka broker config:

```text
ORDER_KAFKA_BROKERS=kafka:9092
NOTIFICATION_KAFKA_BROKERS=kafka:9092
```

---

## OpenTelemetry Collector

Local telemetry should use a local Collector.

Example:

```text
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
```

Local Collector may use debug exporters before full observability backends are added.

---

## MailHog / Fake SMTP Later

Local notification testing can use MailHog.

Example:

```text
SMTP_HOST=mailhog
SMTP_PORT=1025
EMAIL_FROM_ADDRESS=no-reply@bfstore.local
```

---

## Payment Simulator

Local payments should use a simulator.

Example:

```text
PAYMENT_PROVIDER=simulated
PAYMENT_SIMULATED_FAILURE_RATE=0.05
PAYMENT_TIMEOUT_MS=2000
```

This lets checkout behaviour be tested without real payment credentials.

---

## Testing

Local integration tests may use:

```text
Docker Compose
testcontainers
ephemeral MySQL
ephemeral Kafka
fake SMTP
payment simulator
local OTel Collector
```

The app code should not change for test resources.

---

## Practical Rules

```text
Use Docker Compose to document local attachments.
Use service-owned local databases.
Use local Kafka for event flows.
Use a local OTel Collector for telemetry.
Use fake providers for email/payment locally.
Keep local credentials safe and low-risk.
Do not hard-code local-only assumptions in application code.
```

---

## Final Rule

```text
Local backing services should make development realistic without making it dangerous.
```
