# Telemetry Resources

This document describes how `pkg/platform/telemetry` should support OpenTelemetry resources for **bfstore**.

It complements:

```text
docs/architecture/opentelemetry-resources.md
docs/architecture/semantic-conventions.md
docs/architecture/opentelemetry.md
docs/architecture/opentelemetry-components.md
docs/architecture/instrumentation.md
docs/architecture/telemetry-signals.md
docs/architecture/context-propagation.md
```

---

## Purpose

`pkg/platform/telemetry` should centralise resource creation so every bfstore service emits telemetry with a consistent identity.

---

## Responsibilities

This package should help configure:

```text
service.name
service.namespace
service.version
deployment.environment.name
host/process/container attributes where useful
resource detectors where appropriate
```

It should not contain business-flow attributes.

---

## ResourceConfig

Suggested configuration shape:

```go
type ResourceConfig struct {
    ServiceName string
    Namespace   string
    Version     string
    Environment string
}
```

Recommended values:

```text
ServiceName:
  bfstore-order-service

Namespace:
  bfstore

Version:
  build or release version

Environment:
  local, dev, staging, production
```

---

## NewResource Helper

Possible helper shape:

```go
func NewResource(ctx context.Context, cfg ResourceConfig) (*resource.Resource, error) {
    return resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(cfg.ServiceName),
            semconv.ServiceNamespace(cfg.Namespace),
            semconv.ServiceVersion(cfg.Version),
            attribute.String("deployment.environment.name", cfg.Environment),
        ),
        resource.WithTelemetrySDK(),
        resource.WithHost(),
        resource.WithProcess(),
    )
}
```

This helper should be called during service startup before tracer and meter providers are created.

---

## Service Startup Usage

Example:

```go
res, err := telemetry.NewResource(ctx, telemetry.ResourceConfig{
    ServiceName: "bfstore-order-service",
    Namespace: "bfstore",
    Version: version,
    Environment: cfg.Environment,
})
if err != nil {
    return err
}

tp := trace.NewTracerProvider(
    trace.WithResource(res),
)
```

---

## Environment Variable Support

The package may support or document:

```text
OTEL_SERVICE_NAME
OTEL_RESOURCE_ATTRIBUTES
```

Example:

```bash
OTEL_SERVICE_NAME=bfstore-order-service
OTEL_RESOURCE_ATTRIBUTES=service.namespace=bfstore,deployment.environment.name=local,service.version=0.1.0
```

Deployment-specific values should come from configuration/environment.

Code defaults should be safe for local development only.

---

## Defaults

Safe defaults:

```text
service.namespace = bfstore
deployment.environment.name = local
```

Unsafe defaults:

```text
service.name = unknown_service
service.name = app
deployment.environment.name = production
```

Do not silently default to production.

Do not allow unknown service names in real runtime config.

---

## Resource Detectors

The package may enable:

```text
resource.WithTelemetrySDK()
resource.WithHost()
resource.WithProcess()
```

Later, for Kubernetes/cloud:

```text
container detector
Kubernetes detector
cloud provider detector
```

Manual identity attributes should remain authoritative for service naming.

---

## What Not To Include

Do not include these as resource attributes:

```text
order.id
basket.id
customer.email
customer.tier
bfstore.checkout.stage
bfstore.payment.provider
shipping.address
correlation_id
trace_id
```

Resources describe the producer, not the individual request.

---

## Testing Guidance

Recommended tests:

```text
NewResource sets service.name
NewResource sets service.namespace
NewResource sets service.version
NewResource sets deployment.environment.name
NewResource rejects empty service name
NewResource rejects unknown_service where appropriate
local environment default is safe
resource creation includes telemetry SDK information
```

Integration checks:

```text
Collector receives spans with expected service.name
metrics are grouped under expected service.name
logs can be filtered by service.name
```

---

## Practical Rules

```text
Centralise resource creation.
Set service.name explicitly.
Use service.namespace=bfstore.
Use deployment.environment.name from config.
Use service.version from build metadata where possible.
Attach resources before creating providers.
Do not put request/business/customer data into resources.
Use detectors for runtime detail, not core service identity.
Test resource attributes early.
```

---

## Final Rule

```text
pkg/platform/telemetry should give every bfstore service a clear telemetry identity card.
```
