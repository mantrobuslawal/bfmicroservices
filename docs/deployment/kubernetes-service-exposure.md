# Kubernetes Service Exposure

This document defines how bfstore services are exposed inside Kubernetes.

---

## Purpose

This document explains:

```text
containerPort
Service port
targetPort
ClusterIP
Ingress/Gateway
internal vs external exposure
```

---

## Core Rule

```text
Pod listens.
Service routes.
Ingress/Gateway exposes public traffic where appropriate.
```

---

## Container Port

A container port documents the port the application listens on.

Example:

```yaml
containers:
  - name: catalog-service
    ports:
      - name: grpc
        containerPort: 50051
```

The application should bind to the configured port.

---

## Kubernetes Service

A Kubernetes Service gives stable network access to pods.

Example:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: catalog-service
spec:
  selector:
    app: catalog-service
  ports:
    - name: grpc
      port: 50051
      targetPort: grpc
```

Other services can call:

```text
catalog-service:50051
```

---

## Service Types

Use service types deliberately.

```text
ClusterIP:
  internal cluster access

LoadBalancer:
  external load balancer, usually public or VPC-facing depending on cloud config

NodePort:
  direct node port exposure, use carefully

ExternalName:
  DNS alias for external services
```

Default bfstore recommendation:

```text
internal services -> ClusterIP
api-gateway -> exposed through Ingress/Gateway/LoadBalancer
```

---

## Ingress / Gateway

Public traffic should enter through a deliberate entry point.

Example:

```text
internet/client
  -> Ingress or Gateway
  -> api-gateway Service
  -> internal gRPC services
```

Do not expose every internal gRPC service publicly.

---

## targetPort

`targetPort` routes Service traffic to the container port.

Using named ports is clearer:

```yaml
ports:
  - name: grpc
    port: 50051
    targetPort: grpc
```

This avoids accidental mismatch if container ports change later.

---

## Service Discovery

Kubernetes DNS enables service discovery:

```text
catalog-service
catalog-service.default.svc.cluster.local
catalog-service.bfstore.svc.cluster.local
```

Callers should receive service addresses through config:

```text
CATALOG_SERVICE_ADDR=catalog-service:50051
PAYMENT_SERVICE_ADDR=payment-service:50055
```

---

## Health Checks

Kubernetes probes should target the appropriate service health endpoint.

Examples:

```text
gRPC health check
HTTP /healthz
HTTP /readyz
```

Readiness controls whether traffic should be sent to the pod.

Liveness controls whether the pod should be restarted.

---

## Practical Rules

```text
Use ClusterIP for internal services.
Expose api-gateway deliberately.
Use named ports.
Keep service addresses configurable.
Do not expose internal gRPC services publicly by default.
Use readiness and liveness probes.
Restrict traffic later with NetworkPolicy.
```

---

## Final Rule

```text
Kubernetes networking should make service access stable, explicit, and least-exposed.
```
