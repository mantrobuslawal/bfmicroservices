# Port Binding

This document defines how **bfstore** services expose network APIs by binding to ports.

---

## Purpose

This document explains:

```text
what port binding means
HTTP vs gRPC ports
service-to-service exposure
API gateway/public exposure
health check ports
metrics/admin port policy
```

---

## Core Rule

```text
The app should start itself, bind to a configured port, and serve requests directly.
```

The platform should route traffic to the app. The app should serve traffic.

---

## Service Port Binding

Each deployable service should own its listener.

Examples:

```text
api-gateway -> HTTP_PORT
catalog-service -> CATALOG_GRPC_PORT
basket-service -> BASKET_GRPC_PORT
order-service -> ORDER_GRPC_PORT
payment-service -> PAYMENT_GRPC_PORT
```

Ports should come from config, not hard-coded business logic.

---

## HTTP and gRPC

bfstore uses different protocols for different purposes.

```text
HTTP:
  api-gateway public API

gRPC:
  internal service-to-service APIs
```

Example:

```text
api-gateway:
  HTTP :8080

catalog-service:
  gRPC :50051
```

---

## Public vs Internal Exposure

Public traffic should enter through deliberate entry points.

Recommended public entry:

```text
Ingress / Gateway / LoadBalancer -> api-gateway
```

Internal gRPC services should remain internal by default.

Examples:

```text
api-gateway -> catalog-service:50051
order-service -> payment-service:50055
order-service -> shipping-service:50056
```

---

## Health Checks

Services should expose health checks.

For gRPC services:

```text
grpc.health.v1.Health
same gRPC port where practical
```

For HTTP gateway:

```text
/healthz
/readyz
```

Health checks support readiness, liveness, rollouts, and restarts.

---

## Metrics and Admin Ports

Additional ports should be deliberate.

Possible examples:

```text
metrics port
debug/admin port
profiling port
```

Early bfstore recommendation:

```text
export telemetry to OpenTelemetry Collector
avoid unnecessary extra service ports
```

Rule:

```text
Every extra port is another surface area.
```

---

## Bind Address

In containers, bind to an address reachable by the platform.

Usually:

```text
0.0.0.0
```

or:

```text
:50051
```

Avoid binding only to `127.0.0.1` when other containers or Kubernetes Services need to reach the process.

---

## Practical Rules

```text
Bind service APIs to configured ports.
Use HTTP for public gateway traffic.
Use gRPC for internal service-to-service traffic.
Keep internal services private by default.
Expose public traffic through gateway/ingress.
Use health checks.
Avoid unnecessary admin/debug ports.
Document port conventions.
```

---

## Final Rule

```text
Port binding makes each bfstore service a self-contained network process.
```
