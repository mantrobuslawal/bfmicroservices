# OpenTelemetry Resources

This document defines the resource attribute policy for **bfstore** telemetry.

OpenTelemetry resources describe the entity producing telemetry, such as a service, process, container, pod, deployment, host, or cloud runtime.

---

## Purpose

Resources answer:

```text
Who produced this telemetry?
Where was it running?
Which version produced it?
Which environment produced it?
```

For bfstore, this allows traces, metrics, and logs to be grouped by service and environment.

---

## Core Rule

Every bfstore service must set resource attributes consistently.

Required baseline:

```text
service.name
service.namespace
service.version
deployment.environment.name
```

Example:

```text
service.name = bfstore-order-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment.name = local
```

---

## Service Name Policy

`service.name` must be explicit.

Good:

```text
bfstore-catalog-service
bfstore-basket-service
bfstore-order-service
bfstore-payment-service
bfstore-notification-service
```

Bad:

```text
unknown_service
app
main
api
service
```

Rule:

```text
Do not allow unknown_service in real telemetry.
```

---

## Baseline Attributes

For every service:

```text
service.name = bfstore-<service-name>
service.namespace = bfstore
service.version = <build or release version>
deployment.environment.name = local|dev|staging|production
```

---

## Docker Compose Attributes

For local Docker Compose, useful additional attributes may include:

```text
container.name
host.name
```

Example:

```text
container.name = bfstore-order-service
host.name = local-dev
```

Do not rely on container names as a replacement for `service.name`.

---

## Kubernetes Attributes

When bfstore runs on Kubernetes, resource detectors or platform configuration should add:

```text
k8s.namespace.name
k8s.deployment.name
k8s.pod.name
k8s.container.name
k8s.cluster.name
```

Example:

```text
k8s.namespace.name = bfstore
k8s.deployment.name = order-service
k8s.pod.name = order-service-7f8d9c9d9c-abc12
k8s.container.name = order-service
k8s.cluster.name = bfstore-dev
```

---

## Environment Variable Policy

Use environment variables for deployment-specific values.

Local example:

```bash
OTEL_SERVICE_NAME=bfstore-order-service
OTEL_RESOURCE_ATTRIBUTES=service.namespace=bfstore,deployment.environment.name=local,service.version=0.1.0
```

Docker Compose example:

```yaml
environment:
  OTEL_SERVICE_NAME: bfstore-order-service
  OTEL_RESOURCE_ATTRIBUTES: service.namespace=bfstore,deployment.environment.name=local,service.version=0.1.0
```

Rule:

```text
Deployment-specific values should come from configuration or environment.
Code defaults should be safe local defaults only.
```

---

## Resource Detectors

Resource detectors may add runtime/platform attributes.

Useful detectors later:

```text
host
process
container
kubernetes
cloud provider
```

Manual attributes give intentional identity.

Detectors fill in runtime detail.

Do not rely on detectors to guess `service.name` correctly.

---

## Resource Attributes vs Operation Attributes

Resource attributes describe the producer:

```text
service.name = bfstore-payment-service
service.version = 0.2.0
deployment.environment.name = local
```

Operation attributes describe the work:

```text
rpc.system = grpc
rpc.service = bfstore.payment.v1.PaymentService
rpc.method = AuthorisePayment
grpc.status_code = DEADLINE_EXCEEDED
```

Rule:

```text
Resources describe who/where.
Spans, logs, and metrics describe what happened.
```

---

## Resource Attributes vs Correlation IDs

Resource:

```text
service.name = bfstore-order-service
deployment.environment.name = staging
```

Trace/request context:

```text
trace_id = a0892f3577b34da6a3ce929d0e0e4736
correlation_id = checkout-abc-123
```

Rule:

```text
Resources identify the actor.
Trace/correlation IDs identify the story.
```

---

## What Not To Put In Resources

Do not put request/customer/business-flow details in resources.

Bad resource attributes:

```text
order.id
basket.id
customer.email
customer.tier
bfstore.checkout.stage
payment.provider
shipping.address
```

Business facts belong on spans, logs, events, or carefully controlled metrics.

Resources describe the producer, not the user journey.

---

## Service Examples

### Catalog Service

```text
service.name = bfstore-catalog-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment.name = local
```

### Order Service

```text
service.name = bfstore-order-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment.name = local
```

### Payment Service

```text
service.name = bfstore-payment-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment.name = local
```

### Notification Service

```text
service.name = bfstore-notification-service
service.namespace = bfstore
service.version = 0.1.0
deployment.environment.name = local
```

---

## Testing Expectations

Unit tests should verify:

```text
service.name is set
service.namespace is set
service.version is set
deployment.environment.name is set
service.name is not unknown_service
service.name is not empty
local defaults are safe
```

Integration tests should verify:

```text
Collector receives telemetry with expected service.name
trace backend groups spans under expected service name
metrics contain expected resource attributes
logs include expected resource context
```

---

## Practical Rules

```text
Set service.name explicitly.
Use service.namespace=bfstore.
Set service.version from build/release metadata where possible.
Set deployment.environment.name from config.
Attach resources during provider initialisation.
Use resource detectors for runtime/platform detail.
Do not put request IDs, order IDs, customer data, or checkout stages into resources.
Keep resource shape consistent across services.
Test resource identity early.
```

---

## Final Rule

```text
Resources are bfstore’s telemetry identity card.
```
