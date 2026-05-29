# Instrumentation Scopes

This document defines the instrumentation scope policy for **bfstore** telemetry.

Instrumentation scope identifies the code, package, module, framework, or instrumentation library that produced a telemetry item.

---

## Purpose

Instrumentation scope helps bfstore distinguish telemetry emitted by:

```text
service business code
shared platform packages
OpenTelemetry instrumentation libraries
Kafka helpers
gRPC interceptors
database helpers
```

It answers:

```text
Which code produced this span, metric, or log?
Which instrumentation library emitted this telemetry?
Which version of the instrumentation emitted it?
```

---

## Core Distinction

```text
Resource:
  Which service/runtime produced this telemetry?

Instrumentation scope:
  Which code or instrumentation library produced this telemetry item?

Attributes:
  What operation happened?
```

Example:

```text
Resource:
  service.name = bfstore-order-service

Instrumentation scope:
  scope.name = github.com/mantrobuslawal/bfstore/services/order-service
  scope.version = 0.3.1

Attributes:
  rpc.method = Checkout
  bfstore.checkout.stage = payment_authorisation
```

---

## Scope Naming Policy

Use stable, fully qualified package/module-style names.

Application service business code:

```text
github.com/mantrobuslawal/bfstore/services/<service-name>
```

Shared platform packages:

```text
github.com/mantrobuslawal/bfstore/pkg/platform/<package-name>
```

External instrumentation:

```text
use the instrumentation library’s own scope name
```

Avoid vague names:

```text
app
main
tracing
utils
common
telemetry
```

---

## Recommended bfstore Scope Names

Service scopes:

```text
github.com/mantrobuslawal/bfstore/services/catalog-service
github.com/mantrobuslawal/bfstore/services/basket-service
github.com/mantrobuslawal/bfstore/services/inventory-service
github.com/mantrobuslawal/bfstore/services/order-service
github.com/mantrobuslawal/bfstore/services/payment-service
github.com/mantrobuslawal/bfstore/services/shipping-service
github.com/mantrobuslawal/bfstore/services/notification-service
```

Platform package scopes:

```text
github.com/mantrobuslawal/bfstore/pkg/platform/telemetry
github.com/mantrobuslawal/bfstore/pkg/platform/grpc/interceptors
github.com/mantrobuslawal/bfstore/pkg/platform/grpc/auth
github.com/mantrobuslawal/bfstore/pkg/platform/grpc/healthcheck
github.com/mantrobuslawal/bfstore/pkg/platform/kafka
github.com/mantrobuslawal/bfstore/pkg/platform/mysql
```

---

## Scope Version Policy

Use a version where available.

Recommended:

```text
service build version
module/package version
release version
```

Local development fallback:

```text
dev
0.0.0-dev
```

Avoid empty versions once release metadata exists.

---

## Manual vs Automatic Instrumentation

Manual business instrumentation should use bfstore scopes.

Example:

```go
tracer := tp.Tracer(
    "github.com/mantrobuslawal/bfstore/services/order-service",
    trace.WithInstrumentationVersion(version),
)
```

Automatic instrumentation should use the instrumentation library’s own scope.

Example:

```text
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
```

Rule:

```text
Business spans use bfstore scopes.
Auto-generated spans use instrumentation library scopes.
```

---

## Example Checkout Trace

```text
Span: POST /checkout
Resource:
  service.name = bfstore-api-gateway
Scope:
  go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp

Span: OrderService/Checkout
Resource:
  service.name = bfstore-order-service
Scope:
  github.com/mantrobuslawal/bfstore/services/order-service

Span: InventoryService/ReserveStock client call
Resource:
  service.name = bfstore-order-service
Scope:
  go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc

Span: Kafka publish OrderCreated
Resource:
  service.name = bfstore-order-service
Scope:
  github.com/mantrobuslawal/bfstore/pkg/platform/kafka
```

---

## What Not To Put In Scope Names

Do not put request or business-flow data in scope names.

Bad:

```text
OrderService/Checkout
checkout-stage-payment
payment-provider-simulated
customer-tier-premium
order-id-123
```

Correct locations:

```text
span.name = OrderService/Checkout
bfstore.checkout.stage = payment_authorisation
bfstore.payment.provider = simulated
```

Rule:

```text
If the value changes per request, it is not an instrumentation scope.
```

---

## Debugging Uses

Instrumentation scopes help answer:

```text
Which package produced this span?
Did pkg/platform/grpc/interceptors v0.2.0 increase latency?
Are Kafka helper metrics behaving differently after an upgrade?
Which instrumentation library produced this noisy span?
Are spans from automatic gRPC instrumentation missing?
```

Useful filters:

```text
scope.name = github.com/mantrobuslawal/bfstore/pkg/platform/kafka
```

```text
scope.name = github.com/mantrobuslawal/bfstore/services/payment-service
scope.version = 0.4.0
```

---

## Testing Expectations

Unit tests should verify:

```text
service tracer has expected scope name
platform package tracer has expected scope name
scope version is included where available
scope name is not empty
scope name is not app, main, utils, or tracing
```

Integration checks should verify:

```text
test span from order-service has order-service scope
test span from pkg/platform/kafka has kafka package scope
automatic gRPC spans use instrumentation library scope
```

---

## Practical Rules

```text
Use stable fully qualified scope names.
Use service scopes for business instrumentation.
Use platform package scopes for shared package telemetry.
Let external instrumentation libraries use their own scope names.
Include scope version where available.
Do not use request-specific values as scope names.
Do not use business stages as scope names.
Do not confuse scope.name with service.name.
Keep scope names boring and consistent.
```

---

## Final Rule

```text
Instrumentation scope identifies the code that produced the telemetry.
```
