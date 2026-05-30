# Disposability

This document defines how **bfstore** services are designed to start quickly, stop safely, and tolerate sudden failure.

---

## Purpose

This document explains:

```text
fast startup
graceful shutdown
sudden death tolerance
stateless processes
worker reentrancy
service-specific examples
```

---

## Core Rule

```text
Processes should be replaceable.
Business state should be recoverable.
```

A bfstore service should be safe to restart, replace, scale, or roll during deployment.

---

## Fast Startup

Services should become ready quickly.

Startup should include:

```text
load config
validate required config
initialise logging/telemetry
connect to required backing services where appropriate
start HTTP/gRPC listener
register health checks
become ready
```

Startup should not include:

```text
long database migrations
large backfills
full search reindex
large cache warmup required for correctness
protobuf generation
data repair jobs
```

Rule:

```text
Expensive work belongs in explicit jobs, not ordinary service startup.
```

---

## Graceful Shutdown

Services should handle shutdown signals.

Expected behaviour:

```text
receive SIGTERM
fail readiness / stop receiving new traffic
stop accepting new work
allow in-flight work to finish within timeout
cancel work that exceeds shutdown budget
flush telemetry
close backing-service connections
exit cleanly
```

For gRPC services, use graceful server shutdown where possible.

---

## Sudden Death Tolerance

Services must also tolerate unexpected failure.

Examples:

```text
pod killed
node failure
OOMKilled
process panic
network partition
```

Important business state must already be durable:

```text
orders in MySQL
payment attempts in MySQL
events in outbox/Kafka
Kafka offsets committed after success
notification attempts recorded where needed
```

Rule:

```text
Do not rely on shutdown hooks for correctness.
```

---

## Service-specific Guidance

### catalog-service

Avoid:

```text
loading all products into memory at startup
rebuilding search indexes during normal startup
```

State should remain in catalogue MySQL.

### basket-service

Basket state should survive pod restarts.

Source of truth:

```text
basket MySQL
```

### order-service

Order operations need:

```text
short transactions
idempotency keys
outbox for event publication
retry-safe external calls
```

### payment-service

Payment operations need:

```text
payment attempt records
provider references
idempotency keys
safe timeout handling
```

### notification-worker

Worker shutdown needs:

```text
stop polling
finish current message where safe
commit offset only after success
idempotent handler
```

---

## Practical Rules

```text
Start quickly.
Stop safely.
Handle SIGTERM.
Expect SIGKILL.
Use readiness probes.
Use graceful gRPC shutdown.
Make Kafka workers replay-safe.
Make payment/order operations retry-safe.
Do not run migrations as surprise startup behaviour.
Use telemetry around startup and shutdown.
```

---

## Final Rule

```text
A bfstore process may disappear; bfstore business state must not.
```
