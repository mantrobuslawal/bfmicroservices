# Telemetry Instrumentation Scopes

This document describes how `pkg/platform/telemetry` should support instrumentation scopes for **bfstore**.

It complements:

```text
docs/architecture/instrumentation-scopes.md
docs/architecture/opentelemetry-resources.md
docs/architecture/semantic-conventions.md
docs/architecture/instrumentation.md
docs/architecture/opentelemetry-components.md
docs/architecture/telemetry-signals.md
docs/architecture/context-propagation.md
```

---

## Purpose

`pkg/platform/telemetry` may provide helpers or constants for consistent instrumentation scope names across bfstore services and shared platform packages.

---

## Scope Naming Rules

Use stable fully qualified names.

Service scopes:

```text
github.com/mantrobuslawal/bfstore/services/<service-name>
```

Platform package scopes:

```text
github.com/mantrobuslawal/bfstore/pkg/platform/<package-name>
```

Avoid vague names:

```text
app
main
utils
common
tracing
telemetry
```

---

## Suggested Constants

Possible future location:

```text
pkg/platform/telemetry/scope.go
```

Example constants:

```go
package telemetry

const (
    ScopeCatalogService      = "github.com/mantrobuslawal/bfstore/services/catalog-service"
    ScopeBasketService       = "github.com/mantrobuslawal/bfstore/services/basket-service"
    ScopeInventoryService    = "github.com/mantrobuslawal/bfstore/services/inventory-service"
    ScopeOrderService        = "github.com/mantrobuslawal/bfstore/services/order-service"
    ScopePaymentService      = "github.com/mantrobuslawal/bfstore/services/payment-service"
    ScopeShippingService     = "github.com/mantrobuslawal/bfstore/services/shipping-service"
    ScopeNotificationService = "github.com/mantrobuslawal/bfstore/services/notification-service"

    ScopeTelemetryPackage = "github.com/mantrobuslawal/bfstore/pkg/platform/telemetry"
    ScopeGRPCInterceptors = "github.com/mantrobuslawal/bfstore/pkg/platform/grpc/interceptors"
    ScopeGRPCAuth         = "github.com/mantrobuslawal/bfstore/pkg/platform/grpc/auth"
    ScopeGRPCHealthcheck  = "github.com/mantrobuslawal/bfstore/pkg/platform/grpc/healthcheck"
    ScopeKafkaPackage     = "github.com/mantrobuslawal/bfstore/pkg/platform/kafka"
    ScopeMySQLPackage     = "github.com/mantrobuslawal/bfstore/pkg/platform/mysql"
)
```

Constants are helpful when they reduce typos. Do not build a giant abstraction maze around them.

---

## Tracer Helper

A helper may centralise tracer creation:

```go
func Tracer(provider trace.TracerProvider, scopeName string, version string) trace.Tracer {
    return provider.Tracer(
        scopeName,
        trace.WithInstrumentationVersion(version),
    )
}
```

Services may still call the OpenTelemetry API directly where clearer.

---

## Meter Helper

A similar helper may exist for meters:

```go
func Meter(provider metric.MeterProvider, scopeName string, version string) metric.Meter {
    return provider.Meter(
        scopeName,
        metric.WithInstrumentationVersion(version),
    )
}
```

---

## Logger Scope

If OpenTelemetry logging is used later, logger creation should follow the same scope naming policy.

---

## Version Policy

Use:

```text
service build version
module/package version
release version
```

Local fallback:

```text
dev
0.0.0-dev
```

Do not leave versions empty once release metadata exists.

---

## What This Package Should Not Do

Do not use instrumentation scope for:

```text
operation names
checkout stages
payment providers
customer tiers
order IDs
request IDs
correlation IDs
```

Correct locations:

```text
span.name
span attributes
log fields
metric attributes where safe
```

---

## Testing Guidance

Recommended tests:

```text
scope constants match documented values
scope names are not empty
scope names include github.com/mantrobuslawal/bfstore for bfstore-owned code
scope names do not use vague values like app/main/utils
tracer helper includes instrumentation version
meter helper includes instrumentation version
```

Integration checks:

```text
business span shows service scope
platform package span shows package scope
automatic instrumentation spans keep external library scope
```

---

## Practical Rules

```text
Keep scope names stable.
Keep scope names specific.
Use fully qualified package/module-style names.
Use version values where available.
Do not use request-specific data in scope names.
Use constants only where they reduce mistakes.
Do not hide normal OpenTelemetry concepts behind unnecessary wrappers.
```

---

## Final Rule

```text
pkg/platform/telemetry should help bfstore identify which code produced telemetry without making instrumentation mysterious.
```
