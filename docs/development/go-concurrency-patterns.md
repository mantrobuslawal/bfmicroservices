# Go Concurrency Patterns

This document defines bfstore conventions for goroutines, channels, and in-process concurrency.

## Core Rule

```text
Goroutines are concurrency tools, not durability tools.
```

## Goroutines

Use goroutines for controlled in-process concurrency.

Good:

```go
go func() {
    if err := server.Serve(listener); err != nil {
        errCh <- err
    }
}()
```

Avoid fire-and-forget business work.

Bad:

```go
go sendOrderConfirmation(orderID)
```

Better:

```text
order-service:
  writes OrderCreated to outbox

notification-worker:
  consumes OrderCreated
  sends email
  retries safely
```

## Channels

Use channels for coordination inside a process.

```go
jobs := make(chan ProductID)
```

Do not use channels for cross-service messaging.

```text
channels:
  in-process coordination

Kafka:
  durable cross-service events
```

## Buffered Channels

Buffered channels can limit concurrency.

```go
sem := make(chan struct{}, 4)

sem <- struct{}{}
go func() {
    defer func() { <-sem }()
    // work
}()
```

## Shutdown

Long-running goroutines must be stoppable.

Use context cancellation and signal-aware service entrypoints.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

## Worker Guidelines

For workers:

```text
accept context.Context
log structured events
respect cancellation
bound concurrency
make handlers idempotent
commit Kafka offsets after success
avoid unbounded goroutine creation
```

## Kafka vs Channels

Use Kafka when work must:

```text
cross service boundaries
survive process death
support replay
support consumer groups
support durable retries
```

Use channels when work is:

```text
inside one process
short-lived
coordinated by the same runtime
safe to lose on process exit
```

## Practical Rules

```text
Do not use goroutines for durable background business work.
Use channels for in-process coordination.
Use Kafka for durable cross-service events.
Bound concurrency.
Make goroutines stoppable.
Respect context cancellation.
Avoid fire-and-forget customer-facing work.
```

## Final Rule

```text
Concurrency should be controlled, observable, and boring under failure.
```
