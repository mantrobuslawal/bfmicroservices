# Resource Attachments

This document defines how bfstore deployment environments attach backing service resources to applications.

It complements:

```text
docs/architecture/backing-services.md
docs/development/local-backing-services.md
docs/security/backing-service-secrets.md
```

---

## Purpose

This document explains:

```text
how resource handles are injected
ConfigMaps vs Secrets
cloud resource outputs
Kubernetes service discovery
resource attachment per service
```

---

## Core Rule

```text
Infrastructure provisions resources.
Deployment config attaches them.
Application code consumes them.
```

---

## Resource Handles

Resource handles tell a service where an attached backing service is.

Examples:

```text
MYSQL_DSN
KAFKA_BROKERS
SMTP_HOST
PAYMENT_PROVIDER_ENDPOINT
OTEL_EXPORTER_OTLP_ENDPOINT
CATALOG_SERVICE_ADDR
```

Resource handles must be injected through config.

---

## Credentials

Credentials authenticate to backing services.

Examples:

```text
MYSQL_PASSWORD
KAFKA_PASSWORD
SMTP_PASSWORD
PAYMENT_PROVIDER_API_KEY
OTEL_EXPORTER_API_KEY
TLS_PRIVATE_KEY
```

Credentials must be injected through secrets.

---

## Kubernetes Attachment Model

Use:

```text
ConfigMap:
  non-secret resource handles

Secret:
  sensitive resource handles and credentials

Deployment:
  injects both as environment variables
```

Example:

```yaml
env:
  - name: ORDER_KAFKA_BROKERS
    valueFrom:
      configMapKeyRef:
        name: order-config
        key: kafka_brokers

  - name: ORDER_MYSQL_DSN
    valueFrom:
      secretKeyRef:
        name: order-db
        key: dsn
```

---

## Cloud Resource Outputs

Infrastructure tools may create:

```text
managed MySQL endpoint
managed Kafka broker list
secret manager entries
object storage bucket name
SMTP provider credentials
```

Those outputs should be connected to application deploys through configuration and secrets.

---

## Service Discovery

In Kubernetes, internal services may be addressed through DNS.

Examples:

```text
catalog-service:50051
inventory-service:50051
payment-service:50051
shipping-service:50051
otel-collector.bfstore-observability.svc.cluster.local:4317
```

The service address should still be config.

---

## Per-service Attachments

### catalog-service

```text
CATALOG_MYSQL_DSN
OTEL_EXPORTER_OTLP_ENDPOINT
PRODUCT_IMAGE_BUCKET later
```

### order-service

```text
ORDER_MYSQL_DSN
ORDER_KAFKA_BROKERS
INVENTORY_SERVICE_ADDR
PAYMENT_SERVICE_ADDR
SHIPPING_SERVICE_ADDR
OTEL_EXPORTER_OTLP_ENDPOINT
```

### notification-service

```text
NOTIFICATION_KAFKA_BROKERS
SMTP_HOST
SMTP_PORT
SMTP_USERNAME
SMTP_PASSWORD
EMAIL_FROM_ADDRESS
OTEL_EXPORTER_OTLP_ENDPOINT
```

---

## Replacement Flow

If a backing service is replaced:

```text
provision replacement resource
verify credentials
update config/secret
redeploy or restart affected services
verify health checks
verify telemetry
```

Application code should not change.

---

## Practical Rules

```text
Use ConfigMaps for non-secret handles.
Use Secrets for sensitive handles and credentials.
Keep resource handles out of code.
Keep credentials out of plain Git.
Make service addresses configurable.
Make resource replacement a deploy/config task, not a code rewrite.
```

---

## Final Rule

```text
Resource attachment is deployment wiring.
```
