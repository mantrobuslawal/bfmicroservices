# Network Exposure

This document defines bfstore network exposure and port security guidance.

---

## Purpose

This document explains:

```text
public vs private ports
NetworkPolicy
debug/admin ports
least exposure
service mesh policy later
```

---

## Core Rule

```text
Open only the ports that need to be reachable by the actors that need to reach them.
```

---

## Public vs Internal Ports

Public-facing:

```text
api-gateway HTTP port
```

Internal-only:

```text
catalog-service gRPC port
basket-service gRPC port
inventory-service gRPC port
order-service gRPC port
payment-service gRPC port
shipping-service gRPC port
notification-service gRPC port if used
OpenTelemetry Collector receiver ports
database ports
Kafka ports
```

---

## Least Exposure

Avoid:

```text
public LoadBalancer for every service
debug/admin ports exposed publicly
database ports exposed outside private network
Kafka exposed unnecessarily
```

Prefer:

```text
ClusterIP internal services
Ingress/Gateway for public entry
NetworkPolicy for allowed service paths
private subnets/VPC controls in cloud
```

---

## NetworkPolicy Later

NetworkPolicy can restrict which pods can talk to which services.

Examples:

```text
api-gateway may call catalog-service and basket-service
order-service may call inventory-service, payment-service, shipping-service
notification-service may consume Kafka
only app services may send telemetry to OTel Collector
```

This supports zero-trust-style service networking.

---

## Debug and Admin Ports

Debug/admin ports are high risk.

Examples:

```text
pprof
admin HTTP server
debug endpoints
metrics endpoints
```

Rules:

```text
do not expose publicly
disable unless needed
protect with network policy/auth
document purpose and owner
```

---

## Metrics Ports

If services expose `/metrics`, keep them internal.

Current bfstore recommendation:

```text
services export telemetry to OpenTelemetry Collector
avoid extra metrics ports unless deliberately needed
```

---

## Bind Address and Exposure

Binding to `0.0.0.0` inside a pod does not automatically mean public internet exposure.

Actual exposure depends on:

```text
Kubernetes Service type
Ingress/Gateway
NetworkPolicy
cloud security groups/firewalls
service mesh policy
```

Rule:

```text
Bind so the platform can reach the app.
Restrict exposure at the platform/network layer.
```

---

## Practical Rules

```text
Expose api-gateway deliberately.
Keep internal gRPC services internal by default.
Do not expose databases or Kafka publicly.
Avoid unnecessary admin/debug ports.
Use NetworkPolicy later.
Use private networking for backing services.
Document all exposed ports.
Treat every port as part of the attack surface.
```

---

## Final Rule

```text
A port is not just connectivity; it is also exposure.
```
