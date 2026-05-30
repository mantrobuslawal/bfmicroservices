# Repository Strategy

This document defines the repository strategy for **bfstore**.

bfstore may use a monorepo for the application platform while keeping service boundaries explicit.

---

## Purpose

This document explains:

```text
repo structure
service boundaries
shared package rules
what belongs in bfstore vs other portfolio repos
naming conventions
branching expectations
```

---

## Current Recommendation

Use a monorepo for the primary bfstore application platform:

```text
github.com/mantrobuslawal/bfstore
```

This supports:

```text
clear portfolio review
shared Protobuf contracts
consistent platform packages
single architecture story
easier local development
coordinated CI/CD
```

---

## Recommended Layout

```text
bfstore/
├── README.md
├── docs/
├── adr/
├── proto/
├── services/
│   ├── catalog-service/
│   ├── basket-service/
│   ├── inventory-service/
│   ├── order-service/
│   ├── payment-service/
│   ├── shipping-service/
│   └── notification-service/
├── pkg/
│   └── platform/
│       ├── telemetry/
│       ├── grpc/
│       ├── kafka/
│       └── mysql/
├── deploy/
│   ├── docker-compose/
│   ├── kubernetes/
│   └── observability/
├── db/
└── .github/
```

---

## Service Boundaries

Each service should be independently buildable and deployable.

Example:

```text
services/catalog-service/
├── cmd/catalog-service/main.go
├── internal/
├── Dockerfile
├── README.md
└── go.mod
```

or, if using a root Go module:

```text
services/catalog-service/
├── cmd/catalog-service/main.go
├── internal/
├── Dockerfile
└── README.md
```

The chosen module strategy should be documented.

---

## Shared Package Rules

Shared code belongs in capability-specific packages.

Good:

```text
pkg/platform/telemetry
pkg/platform/grpc/interceptors
pkg/platform/grpc/auth
pkg/platform/kafka
pkg/platform/mysql
```

Avoid vague dumping grounds:

```text
pkg/common
pkg/utils
pkg/helpers
```

Rule:

```text
If a shared package name does not explain its capability, it is probably too vague.
```

---

## Protobuf Strategy

Contracts live under:

```text
proto/bfstore/
```

The source of truth is the `.proto` file.

Generated code policy must be explicit:

```text
generated during CI/build
or committed and verified
```

Do not manually edit generated files.

---

## Documentation Strategy

Docs should live with the system.

Recommended:

```text
docs/architecture/
docs/requirements/
docs/testing/
docs/security/
docs/deployment/
docs/development/
adr/
```

Docs should explain real engineering decisions and evolve with the repo.

---

## Portfolio Repo Boundaries

The wider portfolio may include separate repos.

Potential boundaries:

```text
bfstore:
  application/microservices platform code

bfstore-infra:
  cloud infrastructure modules and environment stacks

bfstore-platform:
  Kubernetes platform/GitOps/observability/security layer

bfstore-performance:
  load, stress, soak, and performance testing

bfstore-security:
  DevSecOps/security validation, hardening, policy-as-code

bfstore-docs-or-portfolio:
  case studies, diagrams, write-ups, portfolio site material
```

Each repo should have a clear reason to exist.

---

## Branching Expectations

Branches are for code change, not environment identity.

Avoid:

```text
dev branch
staging branch
production branch
```

Prefer:

```text
main
short-lived feature branches
pull requests
tagged releases
environment config selected during deploy
```

Rule:

```text
Branches are for code change.
Environments are for deploy configuration.
```

---

## CI/CD Expectations

CI should support:

```text
linting
unit tests
integration tests
buf lint
buf breaking checks
security scans
Docker image build
image tagging with commit SHA
changed-service detection where useful
```

Deployments should be traceable to:

```text
git commit
image tag
environment
service name
release/version
```

---

## Practical Rules

```text
Use the monorepo to tell a coherent platform story.
Keep service boundaries clear.
Build services independently.
Deploy services independently.
Keep shared packages capability-specific.
Avoid common/utils dumping grounds.
Keep docs versioned with code.
Use branches for code change, not environments.
Make deploys traceable to commits.
```

---

## Final Rule

```text
A repo should have a clear reason to exist.
```
