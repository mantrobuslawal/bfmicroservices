# Scaling and Restarts

This document defines how bfstore services should behave during scaling, restarts, rollouts, and failures.

---

## Purpose

This document explains:

```text
pod restarts
horizontal scaling
rolling deployments
readiness/liveness
graceful shutdown
idempotent consumers
```

---

## Core Rule

```text
A bfstore replica should be disposable.
The service state should be recoverable.
```

---

## Pod Restarts

Services should tolerate pod restarts.

Expected behaviour:

```text
pod exits
new pod starts
config loads
backing services connect
readiness passes
traffic resumes
business data remains intact
```

No basket, order, payment, or notification state should be lost because a pod restarted.

---

## Horizontal Scaling

A service should scale horizontally if its state is externalised.

Example:

```text
catalog-service replicas: 1 -> 3
```

Works if:

```text
each replica reads from catalog MySQL
each replica has same config
no replica owns unique local state
load balancer can route to any replica
```

---

## Rolling Deployments

During rollout:

```text
new pods start
readiness checks pass
old pods drain/stop
traffic shifts safely
```

Services should support graceful shutdown so in-flight gRPC requests and Kafka processing are handled cleanly.

---

## Readiness and Liveness

Readiness should indicate whether the process can serve traffic.

Possible readiness checks:

```text
config loaded
gRPC server listening
critical backing services reachable where appropriate
```

Liveness should indicate whether the process is stuck and should be restarted.

Do not make liveness too dependent on transient backing-service blips.

---

## Graceful Shutdown

Services should:

```text
stop accepting new work
finish or cancel in-flight requests safely
commit or avoid committing Kafka offsets correctly
flush telemetry
close database connections
exit within timeout
```

---

## Kafka Consumers

Kafka consumers should be restart-safe.

Rules:

```text
make handlers idempotent
commit offsets only after successful processing
record delivery attempts where needed
handle duplicate events safely
use consumer groups for scaling
```

---

## Observability

Monitor:

```text
pod restarts
readiness failures
liveness restarts
gRPC error rates
Kafka consumer lag
duplicate processing
payment retry/idempotency behaviour
notification failures
```

Telemetry should include:

```text
service.name
service.version
deployment.environment.name
```

---

## Practical Rules

```text
Treat pod deletion as normal.
Keep durable state in backing services.
Make replicas interchangeable.
Use readiness before serving traffic.
Use graceful shutdown.
Make consumers idempotent.
Observe restarts and replay behaviour.
```

---

## Final Rule

```text
Scaling should add capacity, not create correctness bugs.
```
