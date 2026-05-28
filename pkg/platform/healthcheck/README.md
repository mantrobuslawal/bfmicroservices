# gRPC Health Check Package

This package provides shared helpers for registering and managing the standard gRPC health service across **bfstore** services.

Recommended location:

```text
pkg/platform/grpc/healthcheck/
```

Suggested files:

```text
pkg/platform/grpc/healthcheck/
├── README.md
└── server.go
```

---

## Purpose

This package should make it easy for each bfstore service to:

```text
register the standard gRPC health service
track whole-server health
track service-specific health
mark services SERVING
mark services NOT_SERVING
support graceful shutdown
keep health behaviour consistent across services
```

It should stay small and boring.

Do not build a complex health framework before bfstore needs one.

---

## Why This Lives Under `pkg/platform`

Health checking is a shared platform concern.

It should be consistent across:

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

A shared package avoids each service inventing a slightly different health implementation.

---

## Package Responsibilities

This package may provide:

```text
Manager type
NewManager function
RegisterService method
MarkServing method
MarkNotServing method
MarkServiceServing method
MarkServiceNotServing method
Shutdown method
```

It should wrap:

```go
google.golang.org/grpc/health
google.golang.org/grpc/health/grpc_health_v1
```

---

## Example API Shape

Possible implementation shape:

```go
package healthcheck

import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/health"
    healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Manager struct {
    server   *health.Server
    services []string
}

func NewManager(grpcServer *grpc.Server) *Manager {
    healthServer := health.NewServer()
    healthpb.RegisterHealthServer(grpcServer, healthServer)

    return &Manager{
        server: healthServer,
    }
}

func (m *Manager) RegisterService(serviceName string) {
    m.services = append(m.services, serviceName)
    m.server.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
}

func (m *Manager) MarkServing() {
    m.server.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

    for _, serviceName := range m.services {
        m.server.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
    }
}

func (m *Manager) MarkNotServing() {
    m.server.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

    for _, serviceName := range m.services {
        m.server.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
    }
}

func (m *Manager) MarkServiceServing(serviceName string) {
    m.server.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
}

func (m *Manager) MarkServiceNotServing(serviceName string) {
    m.server.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
}

func (m *Manager) Shutdown() {
    m.MarkNotServing()
    m.server.Shutdown()
}
```

This is intentionally small.

Keep it boring where production matters.

---

## Example Usage in `catalog-service`

```go
const catalogServiceName = "bfstore.catalog.v1.CatalogService"

grpcServer := grpc.NewServer()

healthManager := healthcheck.NewManager(grpcServer)
healthManager.RegisterService(catalogServiceName)

// Start unavailable until dependencies are ready.
healthManager.MarkNotServing()

catalogv1.RegisterCatalogServiceServer(
    grpcServer,
    catalogHandler,
)

if err := catalogDependencies.Ready(ctx); err != nil {
    logger.Error("catalog-service dependencies not ready", "error", err)
} else {
    healthManager.MarkServing()
}
```

---

## Example Runtime Dependency Monitor

```go
func monitorHealth(
    ctx context.Context,
    logger *slog.Logger,
    manager *healthcheck.Manager,
    db *sql.DB,
) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            manager.MarkNotServing()
            return

        case <-ticker.C:
            pingCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
            err := db.PingContext(pingCtx)
            cancel()

            if err != nil {
                logger.Error("dependency health check failed", "error", err)
                manager.MarkNotServing()
                continue
            }

            manager.MarkServing()
        }
    }
}
```

This keeps the actual gRPC health `Check` call cheap because dependency checks happen separately.

---

## Example Graceful Shutdown

```go
func shutdown(
    logger *slog.Logger,
    grpcServer *grpc.Server,
    healthManager *healthcheck.Manager,
) {
    logger.Info("marking service not serving")
    healthManager.Shutdown()

    done := make(chan struct{})

    go func() {
        grpcServer.GracefulStop()
        close(done)
    }()

    select {
    case <-done:
        logger.Info("grpc server stopped gracefully")

    case <-time.After(10 * time.Second):
        logger.Warn("grpc graceful stop timed out; forcing stop")
        grpcServer.Stop()
    }
}
```

---

## Service Names

Use fully-qualified Protobuf service names.

Examples:

```text
bfstore.catalog.v1.CatalogService
bfstore.basket.v1.BasketService
bfstore.inventory.v1.InventoryService
bfstore.order.v1.OrderService
bfstore.payment.v1.PaymentService
bfstore.shipping.v1.ShippingService
bfstore.notification.v1.NotificationService
```

These should match the names used in the `.proto` service definitions.

---

## Kubernetes Probe Examples

Readiness should use service-specific health:

```yaml
readinessProbe:
  grpc:
    port: 50051
    service: bfstore.catalog.v1.CatalogService
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 2
  failureThreshold: 3
```

Liveness can use whole-server health:

```yaml
livenessProbe:
  grpc:
    port: 50051
    service: ""
  initialDelaySeconds: 10
  periodSeconds: 20
  timeoutSeconds: 2
  failureThreshold: 3
```

---

## Testing Guidance

Recommended tests for this package:

```text
NewManager registers health service
RegisterService starts service as NOT_SERVING
MarkServing marks whole-server and registered services SERVING
MarkNotServing marks whole-server and registered services NOT_SERVING
MarkServiceServing affects one service
MarkServiceNotServing affects one service
Shutdown marks services NOT_SERVING and shuts down health server
```

Recommended service tests:

```text
service starts NOT_SERVING before dependencies are ready
service becomes SERVING after dependency readiness passes
service becomes NOT_SERVING when dependency monitor fails
service becomes NOT_SERVING during shutdown
```

---

## What This Package Should Not Do

Do not put business readiness rules directly into this package.

Bad:

```text
healthcheck package knows catalog database schema details
healthcheck package knows payment provider logic
healthcheck package runs checkout simulations
```

Good:

```text
catalog-service checks catalogue dependencies
payment-service checks payment dependencies
healthcheck package only exposes status management helpers
```

The service owns its readiness logic.

The shared package owns the standard gRPC health plumbing.

---

## Practical Rules

```text
Keep the package small.
Use the standard gRPC health service.
Start NOT_SERVING.
Mark SERVING only after service dependencies are ready.
Mark NOT_SERVING before graceful shutdown.
Expose both whole-server and service-specific health.
Do not put service-specific business dependency logic in the shared package.
Do not perform expensive checks inside health RPC handlers.
```

---

## Final Rule

```text
The healthcheck package should provide the traffic signal.
Each service decides when the light should be green.
```
