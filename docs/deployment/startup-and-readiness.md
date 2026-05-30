# Startup and Readiness

This document defines startup, readiness, and liveness guidance for bfstore services.

---

## Purpose

This document explains:

```text
startup probe
readiness probe
liveness probe
dependency checks
fast startup policy
slow startup anti-patterns
```

---

## Core Rule

```text
Started is not the same as ready.
```

---

## States

```text
started:
  process exists

ready:
  safe to receive traffic

live:
  process is healthy enough to keep running
```

---

## Startup Policy

Startup should be fast and predictable.

Startup should do:

```text
load config
validate config
initialise telemetry
open listeners
connect/check critical resources where appropriate
register health services
```

Startup should not do:

```text
long migrations
full reindex
large backfills
huge cache warmups required for correctness
runtime code generation
blocking optional dependency checks
```

---

## Readiness Probes

Readiness protects users.

A service should become ready only when it can serve traffic safely.

For gRPC services:

```text
gRPC health status SERVING
listener active
config loaded
critical dependencies reachable enough
```

For HTTP gateway:

```text
/readyz
```

Readiness should fail during graceful shutdown.

---

## Liveness Probes

Liveness helps recover stuck processes.

Liveness should check whether the process should be restarted.

Avoid making liveness too sensitive to transient downstream issues.

Rule:

```text
Use readiness for dependency availability.
Use liveness for stuck/dead process recovery.
```

---

## Startup Probes

Startup probes can give slow-starting services a longer initial window before liveness begins.

Use startup probes if a service needs legitimate initialisation time.

Do not use startup probes to hide avoidable slow startup.

---

## Dependency Checks

Be careful with dependency checks.

Critical dependency examples:

```text
service database
Kafka for workers/producers where required
```

Optional dependency examples:

```text
search service later
non-critical telemetry backend
debug tooling
```

Readiness should reflect whether the service can safely do its job.

---

## Practical Rules

```text
Keep startup fast.
Use readiness before receiving traffic.
Fail readiness during shutdown.
Use liveness carefully.
Do not use sleep hacks.
Do not run migrations during normal startup.
Measure startup and readiness duration.
```

---

## Final Rule

```text
Traffic should only reach a service that is ready to serve it.
```
