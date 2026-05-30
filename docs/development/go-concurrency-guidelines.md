# Go Concurrency Guidelines

This document defines Go concurrency guidance for bfstore services and workers.

## Core Rule

```text
Use goroutines for in-process concurrency.
Use Kafka for durable cross-service asynchronous work.
```

## Goroutines

Use goroutines for bounded in-process work, not for hiding durable business tasks.

Bad:

```go
func Checkout(ctx context.Context, cmd CheckoutCommand) error {
    go sendConfirmationEmail(cmd.OrderID)
    return nil
}
```

Better:

```text
order-service writes OrderCreated event/outbox record
notification-worker consumes event
notification-worker sends email with retry/idempotency
```

## Channels

Channels are useful for in-process coordination.

```go
jobs := make(chan ProductID)
```

Do not use channels for cross-service communication.

## Context Cancellation

Long-running functions should accept `context.Context`.

```go
func (w *NotificationWorker) Run(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // poll and process
        }
    }
}
```

## Shutdown

Make goroutines stoppable.

Use signal-aware context in service entrypoints:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

## Maps and Locks

Maps are not safe for concurrent writes.

Use:

```text
sync.RWMutex
sync.Map for specific cases
external backing service/cache where appropriate
```

## Worker Design

For Kafka workers:

```text
poll message
handle message
commit offset only after success
respect context cancellation
log structured events
make handler idempotent
```

## Practical Rules

```text
Do not use goroutines for durable background business work.
Use Kafka/workers for cross-service events.
Use channels for in-process coordination.
Always give long-running goroutines a way to stop.
Protect shared maps with locks.
Commit Kafka offsets only after success.
Respect context cancellation.
```

## Final Rule

```text
Concurrency is useful only when it is controlled, observable, and stoppable.
```
