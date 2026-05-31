# Go Code Organisation

This document defines bfstore's Go code organisation conventions.

## Purpose

```text
repo layout
service layout
cmd/internal/pkg roles
generated code placement
dependency direction
```

## Core Rule

```text
The repo layout should tell a reviewer where things live before they open the files.
```

## Recommended Layout

```text
bfstore/
  go.mod
  go.sum

  cmd/
    bfstore-admin/

  services/
    catalog/
      cmd/catalog-service/
      internal/
        catalog/
        config/
        repository/
        transport/

    order/
      cmd/order-service/
      internal/
        order/
        checkout/
        outbox/
        config/
        repository/
        transport/

  pkg/
    platform/
      config/
      logging/
      telemetry/
      grpc/
      shutdown/

  proto/
  db/
  docs/
```

## Directory Roles

```text
cmd/:
  repo-level commands such as bfstore-admin

services/<service>/cmd:
  service entrypoints

services/<service>/internal:
  service-private code

pkg/platform:
  genuinely shared platform code

proto:
  API and event contracts

db:
  migrations and seed data

docs:
  architecture and operational docs
```

## main.go

`main.go` should wire the service together.

It should handle:

```text
config loading
logger setup
telemetry setup
database/client construction
handler construction
server/worker startup
shutdown handling
```

It should not contain business rules, SQL details, payment orchestration, Kafka event semantics, or domain validation.

## Dependency Direction

Recommended direction:

```text
main
  -> transport/repository/config/platform
  -> application/domain
```

Domain/application packages should not depend on transport or infrastructure details.

## Practical Rules

```text
Use one service entrypoint per command.
Keep service-private code under internal.
Use pkg/platform only for genuine reuse.
Keep domain logic away from transport details.
Keep generated code out of core domain packages.
Avoid premature shared packages.
```

## Final Rule

```text
A clear Go layout makes the architecture easier to understand and harder to accidentally break.
```
