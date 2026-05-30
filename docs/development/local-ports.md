# Local Ports

This document defines local development port conventions for bfstore.

---

## Purpose

This document explains:

```text
local Docker Compose port map
which ports are published to host
how to use grpcurl/curl
how to avoid collisions
```

---

## Core Rule

```text
Expose ports to the host only when humans or tools need them.
Expose ports inside the Docker network when services need them.
```

---

## Suggested Local Port Map

Simple early convention:

```text
api-gateway:
  host 8080 -> container 8080

catalog-service:
  host 50051 -> container 50051

basket-service:
  host 50052 -> container 50052

inventory-service:
  host 50053 -> container 50053

order-service:
  host 50054 -> container 50054

payment-service:
  host 50055 -> container 50055

shipping-service:
  host 50056 -> container 50056

notification-service:
  host 50057 -> container 50057, if it exposes APIs
```

---

## Docker Compose Example

```yaml
services:
  api-gateway:
    ports:
      - "8080:8080"
    environment:
      HTTP_PORT: "8080"

  catalog-service:
    ports:
      - "50051:50051"
    environment:
      CATALOG_GRPC_PORT: "50051"
```

---

## Internal Network Access

Not every internal service needs to be published to the host.

Internal calls can use Docker service names:

```text
order-service -> catalog-service:50051
order-service -> payment-service:50055
```

This keeps local networking cleaner and closer to production-style service discovery.

---

## Useful Commands

HTTP gateway:

```bash
curl http://localhost:8080/health
```

gRPC service:

```bash
grpcurl -plaintext localhost:50051 list
```

Specific service method examples should be added once proto services are stable.

---

## Avoiding Port Collisions

Inside Docker, services can use the same internal port if they are addressed by different service names.

Example:

```text
catalog-service:50051
order-service:50051
payment-service:50051
```

But host ports must be unique:

```yaml
catalog-service:
  ports:
    - "50051:50051"

order-service:
  ports:
    - "50054:50051"
```

Rule:

```text
Container port can be standardised.
Host port must be unique.
```

---

## Practical Rules

```text
Document local ports.
Only publish useful ports to the host.
Use Docker network names for internal calls.
Avoid random undocumented ports.
Keep service addresses configurable.
Use grpcurl and curl for local smoke testing.
```

---

## Final Rule

```text
Local ports should make development easy without exposing everything by default.
```
