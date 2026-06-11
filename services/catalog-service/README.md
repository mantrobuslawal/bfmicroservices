# Catalog Service

The `catalog-service` owns product catalog data for bfstore.

It provides gRPC APIs for reading catalog information such as products, categories, product variants, images, and product attribute definitions.

The service is implemented in Go and uses:

* gRPC for service APIs;
* Protobuf-generated contracts from the root `proto` directory;
* MySQL for catalog persistence;
* shared platform gRPC interceptors;
* standard gRPC health checks;
* structured logging with `log/slog`;
* graceful shutdown for local and container runtime behaviour.

## Responsibilities

The catalog service is responsible for:

* serving product catalog read APIs;
* retrieving product and category data from MySQL;
* mapping catalog domain models to Protobuf responses;
* exposing standard gRPC health status;
* participating in shared platform runtime behaviour such as logging, recovery, and correlation ID propagation.

The catalog service is not responsible for:

* basket management;
* inventory reservation;
* order orchestration;
* payment processing;
* shipping;
* notification delivery;
* search indexing;
* recommendation logic.

Those responsibilities belong to separate bfstore services.

## Package layout

```text
services/catalog-service/
├── cmd/
│   └── catalog-service/
│       └── main.go
├── internal/
│   ├── catalog/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── mysql.go
│   ├── grpcadapter/
│   │   ├── catalog_handler.go
│   │   ├── errors.go
│   │   ├── mappers.go
│   │   └── server.go
│   └── health/
│       ├── doc.go
│       └── health.go
├── test/
│   └── integration/
├── Dockerfile
├── README.md
└── go.mod
```

## Runtime architecture

At runtime, the catalog service starts like this:

```text
load configuration
-> create logger
-> open MySQL connection
-> run catalog readiness check
-> create catalog repository
-> create catalog service
-> create gRPC server
-> register platform interceptors
-> register catalog gRPC handler
-> register gRPC health service
-> optionally register gRPC reflection
-> start serving requests
```

During shutdown:

```text
receive SIGINT or SIGTERM
-> mark service NOT_SERVING
-> stop accepting new gRPC traffic
-> allow in-flight requests to finish
-> force stop if graceful shutdown times out
-> close database connection
```

## gRPC server behaviour

The catalog service uses shared platform interceptors from:

```text
pkg/platform/grpc/interceptors
```

Recommended chain:

```go
grpc.NewServer(
	grpc.ChainUnaryInterceptor(
		platforminterceptors.UnaryRecoveryInterceptor(logger),
		platforminterceptors.UnaryCorrelationIDInterceptor(),
		platforminterceptors.UnaryLoggingInterceptor(logger),
	),
)
```

The order matters:

```text
Recovery
-> Correlation ID
-> Logging
-> Handler
```

### Recovery interceptor

The recovery interceptor catches panics from later interceptors and handlers.

It:

* prevents one bad request from crashing the whole process;
* logs panic details server-side;
* returns a safe `codes.Internal` response to the client;
* avoids exposing stack traces or implementation details to callers.

### Correlation ID interceptor

The correlation ID interceptor ensures every request has a correlation ID.

It:

* reads `x-correlation-id` from incoming gRPC metadata;
* reuses the incoming value when present;
* generates a new correlation ID when missing;
* stores the correlation ID in request context;
* returns the correlation ID in response metadata.

Metadata key:

```text
x-correlation-id
```

This helps connect logs across service boundaries.

### Logging interceptor

The logging interceptor writes structured logs for each unary gRPC request.

Current fields include:

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
duration_ms=4
correlation_id=local-dev-123
error="rpc error: code = InvalidArgument desc = invalid page size"
```

## Health checks

The catalog service exposes the standard gRPC health API:

```text
grpc.health.v1.Health
```

Health status is managed by the shared platform health manager:

```text
pkg/platform/healthcheck
```

The platform health manager owns:

```text
gRPC health service registration
whole-server health status
service-specific health status
SERVING / NOT_SERVING transitions
shutdown health status
```

The catalog service owns its own readiness truth through:

```text
services/catalog-service/internal/health
```

The catalog health checker currently verifies whether the catalog database is reachable.

Practical split:

```text
pkg/platform/healthcheck
  shared health status plumbing

services/catalog-service/internal/health
  catalog-specific dependency checks
```

## gRPC reflection

gRPC reflection can be enabled for local development and testing.

Reflection allows tools such as `grpcurl` to discover services and methods without needing local `.proto` files.

Enable reflection with:

```bash
GRPC_REFLECTION_ENABLED=true go run ./cmd/catalog-service
```

Reflection should be used for local development and testing. It should not be enabled by default in production.

## Running locally

From the repository root, start dependencies:

```bash
make up
```

Then start the catalog service:

```bash
make catalog-run
```

Or run directly from the service directory:

```bash
cd services/catalog-service
GRPC_REFLECTION_ENABLED=true go run ./cmd/catalog-service
```

## Running tests

Run catalog service tests:

```bash
make catalog-test
```

Run catalog integration tests:

```bash
make catalog-integration-test
```

Run all Go tests from the repository root:

```bash
make test
```

Run platform interceptor tests:

```bash
go test ./pkg/platform/grpc/interceptors -v
```

Run platform healthcheck tests:

```bash
go test ./pkg/platform/healthcheck -v
```

## Makefile targets

Useful root Makefile targets:

```bash
make catalog-run
make catalog-test
make catalog-integration-test
make catalog-build
make catalog-docker-build
make catalog-grpc-list
make catalog-health
make catalog-list-products
make catalog-list-categories
make catalog-list-products-with-correlation
```

Example target:

```makefile
.PHONY: catalog-list-products-with-correlation
catalog-list-products-with-correlation:
	grpcurl -plaintext \
		-H 'x-correlation-id: local-dev-123' \
		-d '{"page":{"page_size":5}}' \
		localhost:50051 \
		bfstore.catalog.v1.CatalogService/ListProducts
```

## Smoke testing with grpcurl

List available gRPC services:

```bash
make catalog-grpc-list
```

Or directly:

```bash
grpcurl -plaintext localhost:50051 list
```

Expected services include:

```text
bfstore.catalog.v1.CatalogService
grpc.health.v1.Health
grpc.reflection.v1.ServerReflection
```

List catalog service methods:

```bash
grpcurl -plaintext \
  localhost:50051 \
  list bfstore.catalog.v1.CatalogService
```

## Health check smoke test

Check overall service health:

```bash
make catalog-health
```

Or directly:

```bash
grpcurl -plaintext \
  -d '{}' \
  localhost:50051 \
  grpc.health.v1.Health/Check
```

Expected response:

```json
{
  "status": "SERVING"
}
```

Check catalog-specific health:

```bash
grpcurl -plaintext \
  -d '{"service":"bfstore.catalog.v1.CatalogService"}' \
  localhost:50051 \
  grpc.health.v1.Health/Check
```

Expected response:

```json
{
  "status": "SERVING"
}
```

## Catalog API smoke tests

List products:

```bash
make catalog-list-products
```

Or directly:

```bash
grpcurl -plaintext \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

List categories:

```bash
make catalog-list-categories
```

Or directly:

```bash
grpcurl -plaintext \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListCategories
```

Get a product by ID:

```bash
grpcurl -plaintext \
  -d '{"product_id":"prod_gopher_lamp"}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/GetProduct
```

Use a product ID that exists in the local seed data.

## Correlation ID smoke test

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

Send a request without an explicit correlation ID:

```bash
grpcurl -plaintext \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

Expected behaviour:

```text
request succeeds or fails normally
correlation interceptor generates a correlation ID
logs include the generated correlation_id
```

## Graceful shutdown behaviour

The service handles:

```text
SIGINT
SIGTERM
```

During shutdown it should:

```text
mark service NOT_SERVING
stop accepting new gRPC traffic
allow in-flight requests to finish
force stop if graceful shutdown times out
close the database connection
log shutdown progress
```

This behaviour is important for Docker and Kubernetes because containers are normally terminated using `SIGTERM`.

## Configuration

Configuration is loaded by:

```text
services/catalog-service/internal/config
```

Common local settings include:

```text
GRPC_PORT
GRPC_REFLECTION_ENABLED
DATABASE_HOST
DATABASE_PORT
DATABASE_USER
DATABASE_PASSWORD
DATABASE_NAME
```

Check `.env.example` and `internal/config/config.go` for the exact supported environment variables.

## Troubleshooting

### `grpcurl` cannot list services

Check that the service is running and reflection is enabled:

```bash
GRPC_REFLECTION_ENABLED=true go run ./cmd/catalog-service
```

Then retry:

```bash
grpcurl -plaintext localhost:50051 list
```

### Health check returns `SERVICE_UNKNOWN`

Make sure the service name is exactly:

```text
bfstore.catalog.v1.CatalogService
```

Service names are case-sensitive.

### Makefile gives `missing separator`

Makefile recipe lines must start with a tab.

Correct:

```makefile
catalog-grpc-list:
	grpcurl -plaintext localhost:50051 list
```

Incorrect:

```makefile
catalog-grpc-list:
    grpcurl -plaintext localhost:50051 list
```

### Makefile gives `target pattern contains no '%'`

Check for an extra colon in `.PHONY` declarations.

Incorrect:

```makefile
.PHONY: catalog-list-categories:
```

Correct:

```makefile
.PHONY: catalog-list-categories
```

Also check that commands are not accidentally written on the same line as targets.

Incorrect:

```makefile
catalog-health: grpcurl -plaintext -d '{}' localhost:50051 grpc.health.v1.Health/Check
```

Correct:

```makefile
catalog-health:
	grpcurl -plaintext -d '{}' localhost:50051 grpc.health.v1.Health/Check
```

### `command-line-arguments [command-line-arguments.test]`

This usually means a test file was run directly instead of testing the package.

Avoid:

```bash
go test catalog_handler_test.go
```

Use:

```bash
go test ./internal/grpcadapter
```

or from the package directory:

```bash
go test .
```

### Stale generated Protobuf code

If generated Protobuf accessors appear missing, confirm the workspace is using the local root module.

From the repository root:

```bash
go env GOWORK
go list -f '{{.Dir}}' github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1
```

The package should resolve to the local repository, not an old module cache version.

## Design notes

The catalog service follows these rules:

```text
service-specific dependency checks stay inside the service
shared runtime plumbing lives in pkg/platform
gRPC transport adaptation stays in grpcadapter
domain logic stays in catalog
database access stays in repository code
startup and shutdown wiring stays in cmd/catalog-service/main.go
```

This keeps package boundaries clear and avoids turning shared platform packages into a service framework.

## Current runtime foundation

The service currently demonstrates:

```text
gRPC API serving
standard gRPC health checks
gRPC reflection for local development
structured request logging
correlation ID propagation
panic recovery
graceful shutdown
database readiness checks
Makefile-driven local smoke tests
```

## Practical rule

```text
Platform packages provide reusable plumbing.
Service packages own service-specific truth.
```
