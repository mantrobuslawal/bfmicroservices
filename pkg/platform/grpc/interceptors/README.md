# gRPC Interceptors

The `interceptors` package contains shared gRPC interceptors for bfstore services.

Interceptors provide reusable middleware for cross-cutting service concerns such as:

* panic recovery;
* correlation ID propagation;
* request logging;
* metrics;
* OpenTelemetry tracing;
* authentication;
* authorisation;
* client-side metadata propagation.

They keep service handlers focused on transport adaptation and business logic instead of operational plumbing.

## Package location

```text
pkg/platform/grpc/interceptors
```

Current package structure:

```text
pkg/platform/grpc/interceptors/
├── README.md
├── correlation.go
├── correlation_test.go
├── logging.go
├── logging_test.go
├── recovery.go
└── recovery_test.go
```

Future additions may include:

```text
metrics.go
tracing.go
auth.go
client_correlation.go
client_logging.go
```

## Why this lives under `pkg/platform`

These interceptors are shared platform concerns, not catalog-service-specific business logic.

This package is appropriate because it:

* can be reused by multiple bfstore services;
* standardises gRPC behaviour across the platform;
* keeps service handlers clean;
* makes runtime behaviour easier to test;
* demonstrates mature platform engineering practice.

Use this package for behaviour that should be consistent across services such as `catalog-service`, `basket-service`, `inventory-service`, `order-service`, `payment-service`, `shipping-service`, and `notification-service`.

If an interceptor becomes specific to one service, keep it inside that service instead.

Example:

```text
services/payment-service/internal/grpcadapter/interceptors/
```

## Practical rule

```text
Interceptors handle cross-cutting concerns.
Handlers and services handle transport and business concerns.
```

Do not put checkout orchestration, payment authorisation rules, inventory reservation rules, catalog business logic, database setup, Kafka consumer lifecycle, service startup, or background worker logic inside interceptors.

## Current interceptors

The current server-side unary interceptors are:

```go
UnaryRecoveryInterceptor(logger)
UnaryCorrelationIDInterceptor()
UnaryLoggingInterceptor(logger)
```

## Recommended server chain

Recommended server-side chain:

```go
grpc.NewServer(
	grpc.ChainUnaryInterceptor(
		interceptors.UnaryRecoveryInterceptor(logger),
		interceptors.UnaryCorrelationIDInterceptor(),
		interceptors.UnaryLoggingInterceptor(logger),
	),
)
```

## Interceptor order

Order matters.

Given this chain:

```go
grpc.ChainUnaryInterceptor(
	A,
	B,
	C,
)
```

The call flow is:

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

For bfstore services, use this order:

```text
Recovery
-> Correlation ID
-> Logging
-> Handler
```

Why:

```text
Recovery catches panics from later interceptors and handlers.
Correlation ID ensures the request context has a stable request identifier.
Logging records method, status, duration, error, and correlation ID.
The handler runs service-specific request logic.
```

When metrics, tracing, and auth are added, a likely order is:

```go
grpc.NewServer(
	grpc.ChainUnaryInterceptor(
		interceptors.UnaryRecoveryInterceptor(logger),
		interceptors.UnaryCorrelationIDInterceptor(),
		interceptors.UnaryLoggingInterceptor(logger),
		interceptors.UnaryMetricsInterceptor(metrics),
		interceptors.UnaryTracingInterceptor(tracer),
		interceptors.UnaryAuthInterceptor(authenticator),
	),
)
```

The exact order may change when auth and tracing are implemented, but recovery should remain near the outer edge of the chain.

## Recovery interceptor

```go
UnaryRecoveryInterceptor(logger)
```

Purpose:

```text
catch panics
log panic details server-side
return codes.Internal
avoid exposing stack traces to clients
prevent one bad request from crashing the service process
```

Expected behaviour:

```text
normal response -> returned unchanged
normal error    -> returned unchanged
panic           -> recovered and converted to codes.Internal
nil logger      -> falls back to slog.Default()
```

Recovery is a safety net, not a substitute for fixing bugs.

## Correlation ID interceptor

```go
UnaryCorrelationIDInterceptor()
```

Purpose:

```text
read x-correlation-id from incoming gRPC metadata
generate a correlation ID if absent
store the correlation ID in context
send the correlation ID back as response metadata
make the correlation ID available to logging and later outbound calls
```

Metadata key:

```text
x-correlation-id
```

Context helpers:

```go
ContextWithCorrelationID(ctx, correlationID)
CorrelationIDFromContext(ctx)
```

Example request flow:

```text
incoming request
-> read x-correlation-id metadata
-> reuse incoming value or generate a new one
-> store value in context
-> send value as response header
-> call handler
```

Example distributed flow:

```text
api-gateway
-> catalog-service
-> inventory-service
-> order-service
```

The same correlation ID should appear in logs across the whole request path.

## Logging interceptor

```go
UnaryLoggingInterceptor(logger)
```

Purpose:

```text
log every unary RPC consistently
include gRPC method name
include gRPC status code
include duration
include error details on failure
include correlation ID when available
```

Current log fields:

```text
grpc.method
grpc.code
duration_ms
correlation_id
error
```

Example successful request log fields:

```text
grpc.method=/bfstore.catalog.v1.CatalogService/ListProducts
grpc.code=OK
duration_ms=12
correlation_id=local-dev-123
```

Example failed request log fields:

```text
grpc.method=/bfstore.catalog.v1.CatalogService/ListProducts
grpc.code=InvalidArgument
duration_ms=3
correlation_id=local-dev-123
error="rpc error: code = InvalidArgument desc = invalid page size"
```

The logging interceptor should not alter handler responses or errors. It observes the request outcome and records useful operational context.

## Suggested usage in `catalog-service`

In `services/catalog-service/internal/grpcadapter/server.go`:

```go
package grpcadapter

import (
	"log/slog"

	"github.com/mantrobuslawal/bfstore/pkg/platform/grpc/interceptors"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	"google.golang.org/grpc"
)

func NewServer(catalogService *catalog.Service, logger *slog.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryRecoveryInterceptor(logger),
			interceptors.UnaryCorrelationIDInterceptor(),
			interceptors.UnaryLoggingInterceptor(logger),
		),
	)

	handler := NewCatalogHandler(catalogService, logger)

	catalogv1.RegisterCatalogServiceServer(server, handler)

	return server
}
```

If you alias the import to make call sites clearer:

```go
platforminterceptors "github.com/mantrobuslawal/bfstore/pkg/platform/grpc/interceptors"
```

then use:

```go
grpc.NewServer(
	grpc.ChainUnaryInterceptor(
		platforminterceptors.UnaryRecoveryInterceptor(logger),
		platforminterceptors.UnaryCorrelationIDInterceptor(),
		platforminterceptors.UnaryLoggingInterceptor(logger),
	),
)
```

## Smoke test with grpcurl

Start the catalog service with reflection enabled, then run:

```bash
grpcurl -plaintext localhost:50051 list
```

Expected services should include:

```text
bfstore.catalog.v1.CatalogService
grpc.health.v1.Health
grpc.reflection.v1.ServerReflection
```

Send a request with an explicit correlation ID:

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Expected behaviour:

```text
request succeeds or fails normally
logs include correlation_id=local-dev-123
logs include grpc.method
logs include grpc.code
logs include duration_ms
```

Send a request without a correlation ID:

```bash
grpcurl -plaintext \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Expected behaviour:

```text
request succeeds or fails normally
correlation interceptor generates a new ID
logs include the generated correlation_id
```

## Suggested Makefile target

```makefile
.PHONY: catalog-list-products-with-correlation
catalog-list-products-with-correlation:
	grpcurl -plaintext \
		-H 'x-correlation-id: local-dev-123' \
		-d '{"page":{"page_size":5}}' \
		localhost:50051 \
		bfstore.catalog.v1.CatalogService/ListProducts
```

## Testing guidance

Each interceptor should have focused unit tests.

Recommended test coverage:

```text
recovery interceptor:
  successful response passes through
  normal error passes through
  panic is recovered
  panic returns codes.Internal
  nil logger is safe
  handler is called once

correlation ID interceptor:
  incoming x-correlation-id is preserved
  missing x-correlation-id generates a new ID
  whitespace-only x-correlation-id is treated as missing
  correlation ID is stored in context
  handler response passes through
  handler error passes through
  handler is called once

logging interceptor:
  successful response passes through
  normal error passes through
  non-gRPC error passes through
  nil logger is safe
  handler is called once
  correlation_id is logged when present
  correlation_id is not logged when absent
  correlation_id is logged on error
```

Useful Go packages for tests:

```go
context
bytes
errors
io
log/slog
strings
testing

google.golang.org/grpc
google.golang.org/grpc/codes
google.golang.org/grpc/metadata
google.golang.org/grpc/status
```

## What not to put in interceptors

Do not use interceptors for:

```text
database connection setup
repository construction
Kafka consumer lifecycle
configuration loading
port selection
TLS certificate loading
service startup
background workers
business workflows
checkout orchestration
inventory reservation rules
payment authorisation rules
catalog domain logic
```

Those belong in application bootstrap code, service packages, or domain/application layers.

## Future work

Likely future additions:

```text
UnaryMetricsInterceptor
UnaryTracingInterceptor
UnaryAuthInterceptor
UnaryAuthorisationInterceptor
UnaryClientCorrelationInterceptor
UnaryClientLoggingInterceptor
streaming interceptor variants
```

Add new interceptors only when a real service need appears. Keep the package small, tested, and boring.

## Design summary

```text
Recovery keeps the process alive.
Correlation connects logs across a request path.
Logging records what happened.
Handlers execute service-specific request logic.
```


