# gRPC Graceful Shutdown

This document defines how **bfstore** services should shut down gRPC servers safely.

Graceful shutdown allows services to stop accepting new RPCs while giving in-flight RPCs time to complete before a forceful stop fallback.

---

## Purpose

Graceful shutdown protects bfstore during:

```text
Docker Compose restarts
local development Ctrl+C
Kubernetes rolling deployments
pod termination
dependency failure handling
incident recovery
```

It reduces avoidable request failures during normal operational events.

---

## Standard Policy

bfstore services should:

```text
handle SIGINT and SIGTERM
mark health NOT_SERVING before draining
call grpcServer.GracefulStop()
use grpcServer.Stop() only as a timeout fallback
respect context cancellation in handlers
log each shutdown stage
```

---

## GracefulStop vs Stop

| Method | Behaviour | Policy |
|---|---|---|
| `GracefulStop()` | Stops accepting new RPCs and waits for in-flight RPCs to finish. | Normal shutdown path. |
| `Stop()` | Immediately closes connections and cancels active RPCs. | Timeout fallback only. |

---

## Startup and Shutdown Flow

```text
service starts
health = NOT_SERVING
dependencies become ready
health = SERVING
shutdown signal received
health = NOT_SERVING
gRPC server drains in-flight RPCs
force stop if timeout expires
process exits
```

---

## Health Check Interaction

Health checking and graceful shutdown should work together.

Health check:

```text
Do not send me new traffic.
```

Graceful shutdown:

```text
Let me finish current traffic.
```

Before calling `GracefulStop`, services should mark health as `NOT_SERVING`.

---

## Timeout Policy

Every graceful shutdown must have a timeout.

Starting guidance:

```text
catalog-service:       5–10 seconds
basket-service:        5–10 seconds
inventory-service:     10 seconds
order-service:         15–25 seconds
payment-service:       15–25 seconds
shipping-service:      10–15 seconds
notification-service:  10 seconds
```

Timeouts must remain below Kubernetes `terminationGracePeriodSeconds`.

---

## Kubernetes Policy

Application shutdown timeout must be lower than pod termination budget.

Example:

```yaml
terminationGracePeriodSeconds: 30
```

Application:

```go
const grpcShutdownTimeout = 20 * time.Second
```

---

## Context Cancellation

Handlers must pass `context.Context` through service and database layers.

Good:

```go
result, err := h.service.DoWork(ctx, req)
```

Good database usage:

```go
rows, err := db.QueryContext(ctx, query, args...)
```

Avoid using `context.Background()` inside request handling paths when the request context is already available.

---

## Streaming RPCs

Long-lived streaming RPCs must monitor stream context cancellation.

```go
select {
case <-stream.Context().Done():
    return stream.Context().Err()
case <-ticker.C:
    // send update
}
```

---

## Kafka and Background Workers

Services with Kafka consumers/producers should coordinate shutdown:

```text
mark gRPC NOT_SERVING
stop accepting new gRPC work
stop Kafka consumers from polling new messages
finish/commit current message if safe
flush producers/outbox where appropriate
gracefully stop gRPC server
close database connections
exit
```

---

## Outbox Pattern

For critical events such as `OrderCreated`, graceful shutdown should be paired with durable event publishing.

Recommended approach:

```text
write business state and outbox event in same DB transaction
publish asynchronously from outbox
retry unpublished events after restart
```

Graceful shutdown reduces interruption. Durable patterns prevent data loss.

---

## Logging Requirements

Shutdown logs should include:

```text
shutdown signal received
health marked NOT_SERVING
graceful shutdown started
graceful shutdown completed
timeout exceeded
force stop triggered
```

---

## What Not To Do

Avoid:

```text
calling os.Exit immediately on SIGTERM
using Stop as the normal shutdown path
omitting a shutdown timeout
leaving health SERVING during shutdown
closing database connections before RPCs drain
ignoring context cancellation
assuming Kubernetes waits forever
using sleep as the primary shutdown strategy
```

---

## Practical Rules

```text
Use GracefulStop as the normal path.
Use Stop only as a timeout fallback.
Mark health NOT_SERVING before draining.
Always use a shutdown timeout.
Keep timeout below terminationGracePeriodSeconds.
Pass request context through all layers.
Coordinate gRPC shutdown with Kafka/background workers.
Use durable outbox patterns for critical events.
```

---

## Final Rule

```text
Graceful shutdown is a traffic drain, not a data safety mechanism by itself.
```

It should be combined with health checks, context-aware handlers, durable storage, and clear service shutdown order.
