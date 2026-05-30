# Codebase and Deploys

This document defines how **bfstore** interprets the 12-Factor Codebase principle.

The principle is:

```text
One codebase, many deploys.
```

For bfstore, this means each deployable app should have one source of truth in version control, while being deployable to multiple environments using environment-specific configuration.

---

## Purpose

This document explains:

```text
how bfstore interprets 12-Factor codebase
what counts as an app in the distributed system
how monorepo boundaries work
how many deploys relate to one codebase
how shared code is handled
how generated code is handled
how deployments trace back to commits
```

---

## Core Rule

```text
The code should be the same.
The environment should change.
```

Do not create environment-specific copies of source code.

---

## bfstore as a Distributed System

bfstore is a distributed system made up of deployable apps:

```text
api-gateway
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

Each service should follow:

```text
one service codebase
many deploys
```

Example:

```text
catalog-service:
  local deploy
  dev deploy
  staging deploy
  production-style deploy
```

---

## Monorepo Interpretation

bfstore may use a monorepo while still respecting the codebase principle.

Recommended structure:

```text
bfstore/
├── services/
│   ├── catalog-service/
│   ├── basket-service/
│   ├── inventory-service/
│   ├── order-service/
│   ├── payment-service/
│   ├── shipping-service/
│   └── notification-service/
├── proto/
├── pkg/
├── deploy/
├── docs/
└── db/
```

The monorepo is acceptable if:

```text
service boundaries are explicit
each service can be built independently
each service can be deployed independently
shared code is intentional
environment differences live in config
```

---

## Deploys

A deploy is a running instance of an app.

Examples:

```text
catalog-service on a laptop
catalog-service in Docker Compose
catalog-service in Kubernetes dev
catalog-service in Kubernetes staging
catalog-service in production-style Kubernetes
```

Different deploys may run different commits, but they should come from the same codebase.

---

## Environment Configuration

Environment differences should be expressed through configuration.

Examples:

```text
database endpoints
Kafka broker endpoints
OTLP collector endpoints
replica counts
resource limits
feature flags
timeouts
secrets
```

Do not create separate source folders or repos per environment.

Bad:

```text
services/payment-service-dev
services/payment-service-prod
```

Good:

```text
services/payment-service
deploy/kubernetes/overlays/dev
deploy/kubernetes/overlays/prod
```

---

## Shared Code

Shared code should live in named packages or versioned libraries.

Good:

```text
pkg/platform/telemetry
pkg/platform/grpc/interceptors
pkg/platform/grpc/auth
pkg/platform/kafka
pkg/platform/mysql
```

Avoid:

```text
pkg/common
pkg/utils
pkg/helpers
```

Rule:

```text
Shared code should have a name that explains the capability it provides.
```

---

## Generated Protobuf Code

The source of truth for contracts is:

```text
.proto files
```

Generated code is an artefact.

Allowed approaches:

```text
generate during CI/build
commit generated code and verify it is up to date
```

Do not manually edit generated files.

---

## Build Once, Deploy Many

Recommended chain:

```text
commit
  -> CI checks
  -> build Docker image
  -> tag image with commit SHA
  -> deploy to dev
  -> promote same image to staging/prod
```

Example:

```text
ghcr.io/mantrobuslawal/bfstore/order-service:abc123
```

Rule:

```text
A deployment should point back to a commit.
A running container should point back to a build.
```

---

## Violations

Avoid:

```text
separate repos per environment
long-lived environment branches
copied shared code
manual production patches
untraceable deployed artefacts
environment-specific source folders
```

---

## Practical Rules

```text
Keep one source of truth for each app.
Do not create environment-specific app repos.
Treat microservices as deployable app boundaries.
Use config for environment differences.
Keep shared code intentional and named.
Keep generated code traceable to source.
Build once, deploy many.
Tag artefacts with commit SHA.
Make deploys traceable.
```

---

## Final Rule

```text
The repo can be one.
The deploys can be many.
The boundaries must be clear.
```
