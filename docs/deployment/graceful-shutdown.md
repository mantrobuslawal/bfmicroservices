# Graceful Shutdown

This document defines graceful shutdown behaviour for bfstore services and workers.

---

## Purpose

This document explains:

```text
SIGTERM handling
gRPC GracefulStop
worker shutdown
Kafka offset handling
terminationGracePeriodSeconds
telemetry flushing
```

---

## Core Rule

```text
Stop accepting new work, finish or safely abandon current work, then exit before the grace period ends.
```

---

## Shutdown Signal

Kubernetes sends:

```text
SIGTERM
waits for terminationGracePeriodSeconds
SIGKILL if the process has not exited
```

The service should listen for `SIGTERM` and start shutdown immediately.

---

## gRPC Services

For gRPC services:

```text
receive SIGTERM
fail readiness
call GracefulStop
allow in-flight RPCs to finish
fallback to Stop after timeout
flush telemetry
close resources
exit
```

Use a bounded timeout.

Rule:

```text
Graceful shutdown should be kind, but not infinite.
```

---

## HTTP Gateway

For `api-gateway`:

```text
receive SIGTERM
fail readiness
stop accepting new HTTP requests
allow in-flight requests to finish
shutdown HTTP server with timeout
flush telemetry
exit
```

---

## Kafka Workers

For Kafka workers:

```text
receive SIGTERM
stop polling for new messages
finish current message if safe
commit offset only after durable success
leave consumer group cleanly
flush telemetry
exit
```

If processing cannot complete safely:

```text
do not commit offset
allow message to be redelivered
ensure handler is idempotent
```

---

## terminationGracePeriodSeconds

Suggested starting values:

```text
request-serving services:
  30 seconds

notification-worker:
  60 seconds if message processing needs more time
```

Tune based on measured shutdown behaviour.

---

## Telemetry Flushing

Before exit, services should attempt to flush:

```text
logs
metrics
traces
```

Shutdown telemetry should include:

```text
shutdown duration
cancelled request count
in-flight request count
worker message completion/retry
flush success/failure
```

---

## What Not To Do

Avoid:

```text
ignoring SIGTERM
exiting immediately while requests are in-flight
committing Kafka offsets before success
waiting forever for shutdown
not flushing telemetry
using arbitrary sleep instead of readiness/shutdown logic
```

---

## Practical Rules

```text
Handle SIGTERM.
Fail readiness during shutdown.
Stop new work.
Bound graceful shutdown with a timeout.
Commit Kafka offsets only after success.
Make interrupted work retry-safe.
Flush telemetry.
Exit within grace period.
```

---

## Final Rule

```text
Graceful shutdown is controlled disappearing.
```
