# Platform Healthcheck

The `healthcheck` package provides shared gRPC health-check wiring for bfstore services.

It wraps the standard gRPC health service and gives services a small, consistent API for:

* registering service health status;
* marking services as `SERVING`;
* marking services as `NOT_SERVING`;
* exposing the standard `grpc.health.v1.Health` API;
* coordinating health status during startup, readiness checks, and graceful shutdown.

This package is intended for reusable platform-level health behaviour. It should not contain service-specific dependency checks such as MySQL queries, Kafka connectivity checks, cache checks, or domain-specific readiness rules. Those checks belong inside each service.

For example, the catalog service can own the logic for checking whether its MySQL database is reachable, while this package owns the standard gRPC health service registration and status management.

## Responsibilities

This package is responsible for:

* registering the standard gRPC health server;
* tracking whole-server and service-specific health status;
* marking services as `SERVING`;
* marking services as `NOT_SERVING`;
* supporting graceful shutdown by moving services to `NOT_SERVING`;
* giving all bfstore services a consistent health-check integration pattern.

## Non-responsibilities

This package should not:

* connect to service databases;
* know about service repositories;
* run SQL queries;
* know about Kafka topics, consumers, or producers;
* decide whether a specific service is ready;
* contain business logic;
* replace service-owned readiness checks.

Service-specific readiness checks should live in service-local packages such as:

```text
services/catalog-service/internal/health
services/basket-service/internal/health
services/inventory-service/internal/health
```

## Package boundary

The intended split is:

```text
pkg/platform/healthcheck
  Shared gRPC health server registration and status management.

services/<service-name>/internal/health
  Service-specific dependency checks and readiness decisions.
```

For example:

```text
Catalog service readiness:
  - Can the process run?
  - Can it reach MySQL?
  - Can it safely serve catalog requests?

Platform healthcheck:
  - Register grpc.health.v1.Health.
  - Mark the catalog service SERVING.
  - Mark the catalog service NOT_SERVING during shutdown.
```

## Example usage

A service can create a health manager during startup:

```go
healthManager := healthcheck.NewManager(grpcServer)

healthManager.RegisterService("bfstore.catalog.v1.CatalogService")
```

After service dependencies have passed readiness checks:

```go
if err := catalogHealthChecker.Ready(ctx); err != nil {
	logger.Error("catalog service is not ready", "error", err)
	os.Exit(1)
}

healthManager.MarkServing()
```

During graceful shutdown:

```go
healthManager.MarkNotServing()
healthManager.Shutdown()
```

## Example startup flow

A typical service startup flow should look like this:

```text
1. Load configuration.
2. Create logger.
3. Open database or external dependencies.
4. Create service-specific health checker.
5. Run readiness checks.
6. Create gRPC server.
7. Register service handlers.
8. Register platform health manager.
9. Mark service SERVING.
10. Start serving gRPC traffic.
```

## Example shutdown flow

A typical graceful shutdown flow should look like this:

```text
1. Receive SIGINT or SIGTERM.
2. Mark service NOT_SERVING.
3. Stop accepting new traffic.
4. Allow in-flight gRPC requests to finish.
5. Force stop if graceful shutdown times out.
6. Close database connections and other resources.
```

## Health statuses

The package uses the standard gRPC health statuses:

```text
SERVING
NOT_SERVING
UNKNOWN
SERVICE_UNKNOWN
```

For bfstore services, the normal lifecycle should be:

```text
startup      -> NOT_SERVING
ready        -> SERVING
shutting down -> NOT_SERVING
stopped      -> NOT_SERVING
```

## Checking health with grpcurl

Once a service has registered the health server, check overall health with:

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

Check a specific service:

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

## Development guidance

Use this package when adding health support to bfstore services.

Keep the platform manager small and boring. It should provide shared health-check plumbing, not become a service orchestration framework.

A good rule:

```text
If the code manages gRPC health status, it belongs here.
If the code checks whether a service dependency is available, it belongs in that service.
```

## Future improvements

Possible future additions include:

* health status change logging;
* service-specific status helpers;
* readiness monitor helpers;
* test utilities for asserting health status;
* Kubernetes probe documentation;
* integration with OpenTelemetry metrics.

Keep it boring where production matters.
