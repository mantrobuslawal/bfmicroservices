# gRPC Interceptors Package

This package contains shared gRPC interceptors for **bfstore** services.

Interceptors provide reusable middleware for cross-cutting service concerns such as:

```text
panic recovery
correlation IDs
request logging
metrics
OpenTelemetry tracing
authentication
authorisation
client-side metadata propagation
```

They keep service handlers focused on business logic.

---

## Recommended Package Location

```text
pkg/platform/grpc/interceptors/
```

Suggested structure:

```text
pkg/platform/grpc/interceptors/
├── README.md
├── correlation.go
├── logging.go
├── recovery.go
├── metrics.go
├── auth.go
└── client.go
```

This package can be reused by:

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

---

## Why This Lives Under `pkg/platform`

The interceptors are shared platform concerns, not business logic for one service.

`pkg/platform/grpc/interceptors` is appropriate because:

```text
multiple services can reuse it
it standardises behaviour across the platform
it keeps handlers clean
it demonstrates mature platform engineering practice
```

If an interceptor becomes specific to one service, keep that service-specific interceptor inside that service instead.

Example:

```text
services/payment-service/internal/grpc/interceptors/
```

Only use the shared package for genuinely reusable behaviour.

---

## Practical Rule

```text
Interceptors handle cross-cutting concerns.
Handlers and services handle business concerns.
```

Do not put checkout orchestration, payment authorisation rules, inventory reservation rules, or catalogue business logic inside interceptors.

---

## Recommended Initial Interceptors

For early bfstore gRPC services, start with:

```text
UnaryRecovery
UnaryCorrelationID
UnaryLogging
UnaryClientCorrelation
UnaryClientLogging
```

Add later:

```text
UnaryMetrics
UnaryTracing
UnaryAuth
UnaryAuthorisation
UnaryRateLimit
UnaryFaultInjection for local/dev/test only
streaming variants when streaming RPCs are introduced
```

Keep the first implementation small and reliable.

---

## Server Interceptor Chain

Recommended local/dev server chain:

```go
grpcServer := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        interceptors.UnaryRecovery(logger),
        interceptors.UnaryCorrelationID(logger),
        interceptors.UnaryLogging(logger),
    ),
)
```

Later, when metrics and auth are added:

```go
grpcServer := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        interceptors.UnaryRecovery(logger),
        interceptors.UnaryCorrelationID(logger),
        interceptors.UnaryLogging(logger),
        interceptors.UnaryMetrics(metrics),
        interceptors.UnaryAuth(authenticator),
    ),
)
```

---

## Interceptor Order

Order matters.

```go
grpc.ChainUnaryInterceptor(
    A,
    B,
    C,
)
```

Call flow:

```text
request
-> A before
  -> B before
    -> C before
      -> handler
    <- C after
  <- B after
<- A after
response
```

Recommended order:

```text
Recovery
-> Correlation ID
-> Logging
-> Metrics
-> Auth
-> Handler
```

Why:

```text
Recovery catches panics from later interceptors and handlers.
Correlation runs before logging so logs include correlation ID.
Logging records method, status, duration, and correlation ID.
Metrics records status and duration.
Auth blocks unauthenticated calls before business logic.
```

---

## Recovery Interceptor

Purpose:

```text
catch panics
log safely
return codes.Internal
avoid exposing stack traces to clients
```

Example shape:

```go
func UnaryRecovery(logger *slog.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req any,
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (resp any, err error) {
        defer func() {
            if recovered := recover(); recovered != nil {
                logger.Error(
                    "panic recovered in grpc handler",
                    "grpc.method", info.FullMethod,
                    "panic", recovered,
                )

                err = status.Error(codes.Internal, "internal server error")
            }
        }()

        return handler(ctx, req)
    }
}
```

Recovery is a safety net, not a substitute for fixing bugs.

---

## Correlation ID Interceptor

Purpose:

```text
read x-correlation-id from incoming metadata
generate one if absent
store it in context
make it available to logs and outbound calls
```

Example metadata key:

```text
x-correlation-id
```

Checkout example:

```text
api-gateway
-> order-service
-> inventory-service
-> payment-service
-> shipping-service
```

The same correlation ID should appear in logs across the whole flow.

---

## Logging Interceptor

Purpose:

```text
log every RPC consistently
include method name
include status code
include duration
include correlation ID where available
```

Example log fields:

```text
grpc.method
grpc.status_code
duration_ms
correlation_id
```

This makes debugging service behaviour far easier.

---

## Metrics Interceptor

Purpose:

```text
record RPC count
record RPC duration
record status code
record method name
```

Useful dashboard views:

```text
CatalogService/GetProduct p95 latency
OrderService/Checkout error rate
PaymentService/AuthorisePayment DeadlineExceeded count
InventoryService/ReserveStock Aborted count
```

Metrics turns “it feels slow” into evidence.

---

## Auth Interceptor

Purpose:

```text
read authorization metadata
validate credentials
attach principal to context
return codes.Unauthenticated for invalid/missing credentials
return codes.PermissionDenied when identity is valid but not allowed
```

Initial local development may skip auth. Document the intended position now and implement later.

Do not log raw auth tokens.

---

## Client Correlation Interceptor

Purpose:

```text
copy current correlation ID from context
append it to outgoing gRPC metadata
```

Example:

```go
ctx = metadata.AppendToOutgoingContext(
    ctx,
    "x-correlation-id",
    correlationID,
)
```

This should be used by services making outbound calls, especially:

```text
order-service -> inventory-service
order-service -> payment-service
order-service -> shipping-service
api-gateway -> catalog-service
api-gateway -> basket-service
```

---

## Client Logging Interceptor

Purpose:

```text
log outbound RPCs consistently
include method
status code
duration
correlation ID
target service where available
```

This is useful when debugging orchestration flows such as checkout.

---

## What Not To Put In Interceptors

Do not put business workflows in interceptors.

Bad:

```text
Checkout interceptor reserves stock and authorises payment.
```

Good:

```text
Order service application layer orchestrates checkout.
Interceptors log, trace, authenticate, and recover.
```

Do not use interceptors for:

```text
database connection setup
Kafka consumer lifecycle
loading configuration
choosing TCP ports
TLS certificate setup
service startup
background workers
```

Those belong in application bootstrap code.

---

## Suggested Usage in `catalog-service`

```go
package main

import (
    "log/slog"

    "google.golang.org/grpc"

    "github.com/mantrobuslawal/bfstore/pkg/platform/grpc/interceptors"
)

func newGRPCServer(logger *slog.Logger) *grpc.Server {
    return grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            interceptors.UnaryRecovery(logger),
            interceptors.UnaryCorrelationID(logger),
            interceptors.UnaryLogging(logger),
        ),
    )
}
```

---

## Suggested Usage in `order-service`

`order-service` is an orchestrator and will make outbound calls to other services.

It should use both server-side and client-side interceptors.

Server-side:

```go
grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        interceptors.UnaryRecovery(logger),
        interceptors.UnaryCorrelationID(logger),
        interceptors.UnaryLogging(logger),
    ),
)
```

Client-side:

```go
conn, err := grpc.NewClient(
    target,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithChainUnaryInterceptor(
        interceptors.UnaryClientCorrelation(),
        interceptors.UnaryClientLogging(logger),
    ),
)
```

---

## Testing Guidance

Each interceptor should have focused unit tests.

Recommended tests:

```text
recovery interceptor converts panic to codes.Internal
correlation interceptor preserves incoming x-correlation-id
correlation interceptor generates ID when absent
client correlation interceptor appends outgoing metadata
logging interceptor calls handler exactly once
auth interceptor rejects missing authorization metadata
auth interceptor passes valid principal to handler
```

Use:

```go
metadata.NewIncomingContext
metadata.AppendToOutgoingContext
metadata.FromOutgoingContext
status.Code
```

for tests.

---

## Portfolio Value

This package demonstrates that bfstore is not just a set of handlers.

It shows:

```text
consistent service behaviour
clean separation of concerns
operability thinking
observability readiness
auth readiness
platform reuse
senior engineering judgement
```

Client/interview phrasing:

```text
I use shared gRPC interceptors for cross-cutting platform concerns such as
recovery, correlation IDs, logging, metrics, authentication, and outbound
metadata propagation, keeping service handlers focused on domain behaviour.
```

---

## Final Rule

```text
Keep boring platform behaviour centralised.
Keep business behaviour explicit in services.
```

That is how bfstore stays clean as the number of services grows.
