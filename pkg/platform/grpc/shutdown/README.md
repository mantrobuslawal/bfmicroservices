# gRPC Shutdown Package

This package provides shared helpers for graceful gRPC shutdown across **bfstore** services.

Recommended location:

```text
pkg/platform/grpc/shutdown/
```

Suggested files:

```text
pkg/platform/grpc/shutdown/
├── README.md
└── shutdown.go
```

---

## Purpose

This package helps bfstore services shut down consistently by:

```text
marking health NOT_SERVING
calling grpc.Server.GracefulStop()
forcing grpc.Server.Stop() after timeout
logging shutdown progress
supporting OS signal driven shutdown
```

It should stay small, explicit, and boring.

---

## Why This Lives Under `pkg/platform`

Graceful shutdown is a shared platform concern.

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

A shared helper avoids each service inventing a slightly different shutdown pattern.

---

## Expected Service Usage

Typical service flow:

```text
create grpc.Server
register health manager
register service handlers
start serving in goroutine
wait for SIGINT/SIGTERM
call shutdown helper
exit
```

Example:

```go
shutdown.Graceful(
    logger,
    grpcServer,
    healthManager,
    10*time.Second,
)
```

---

## Health Manager Contract

This package should use a small interface:

```go
type HealthManager interface {
    MarkNotServing()
    Shutdown()
}
```

This keeps the shutdown helper easy to test and avoids tight coupling.

---

## Example Usage in `catalog-service`

```go
const grpcShutdownTimeout = 10 * time.Second

shutdownSignal := make(chan os.Signal, 1)
signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

select {
case sig := <-shutdownSignal:
    logger.Info("shutdown signal received", "signal", sig.String())
case err := <-serverErr:
    logger.Error("grpc server stopped unexpectedly", "error", err)
}

shutdown.Graceful(
    logger,
    grpcServer,
    healthManager,
    grpcShutdownTimeout,
)
```

---

## Recommended Shutdown Sequence

The helper should:

```text
log graceful shutdown start
mark health NOT_SERVING
notify health watchers, if supported
start GracefulStop in goroutine
wait for GracefulStop or timeout
if graceful completes, log success
if timeout expires, call Stop and log forced stop
```

---

## Kafka and Background Worker Coordination

This package handles the gRPC server drain only.

Services with Kafka consumers, producers, outbox publishers, or other background workers must coordinate those separately.

The shutdown helper should not know Kafka-specific business behaviour.

---

## Testing Guidance

Recommended tests:

```text
calls healthManager.MarkNotServing before GracefulStop
calls healthManager.Shutdown when provided
returns after GracefulStop completes
calls Stop when timeout expires
does not call Stop when GracefulStop completes before timeout
handles nil health manager
logs key lifecycle events
```

---

## What This Package Should Not Do

Do not put service-specific business shutdown behaviour here.

Bad:

```text
shutdown package commits Kafka offsets
shutdown package publishes OrderCreated
shutdown package closes catalog database directly
shutdown package knows checkout state machine
```

Good:

```text
shutdown package drains gRPC server
service main coordinates dependencies and workers
```

---

## Practical Rules

```text
GracefulStop first.
Stop only after timeout.
Mark health NOT_SERVING before draining.
Keep timeout explicit.
Keep helper small.
Do not close databases inside this helper.
Do not manage Kafka lifecycle inside this helper.
Log shutdown clearly.
```

---

## Final Rule

```text
The shutdown package drains the gRPC road.
The service decides what vehicles must park first.
```
