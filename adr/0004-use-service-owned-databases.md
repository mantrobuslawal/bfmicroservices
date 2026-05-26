# ADR-0004: Use Service-Owned Databases

## Status

Accepted

## Date

2026-05-26

## Context

bfstore is designed as a microservice system. Each service must own its own business capability and data.

A shared database would simplify joins and local development, but would create hidden coupling between services and weaken the microservice architecture.

## Decision

Each bfstore service will own its own logical MySQL schema.

Services must not directly read or write another service’s schema.

Local development may use one MySQL instance with multiple schemas, but ownership boundaries remain strict.

## Drivers

This decision supports:

```text
clear service ownership
independent schema evolution
least privilege database access
strong microservice boundaries
better alignment with service contracts
realistic platform engineering design
```

## Alternatives Considered

### Option 1: Shared Database

Benefits:

```text
simple joins
easy reporting
less setup
```

Costs:

```text
tight coupling
weak ownership
hard independent deployment
schema changes can break unrelated services
services become database clients rather than domain owners
```

### Option 2: Schema per Service

Benefits:

```text
clear logical ownership
manageable local development
least privilege possible
good stepping stone to separate databases
```

Costs:

```text
more migrations
no cross-service joins
cross-service data requires APIs, events, projections, or snapshots
```

### Option 3: Physically Separate Database per Service

Benefits:

```text
strongest isolation
clear operational boundaries
```

Costs:

```text
more infrastructure overhead
more complex local development
higher operational cost
```

## Consequences

### Positive

```text
service boundaries are enforceable
database permissions can follow least privilege
schemas can evolve independently
data ownership is clear to reviewers and clients
```

### Negative

```text
reporting and cross-service queries require projections
more database migration management
data duplication through snapshots/projections must be governed
eventual consistency must be accepted in some areas
```

## Implementation Notes

Initial schemas:

```text
bfstore_catalog
bfstore_inventory
bfstore_basket
bfstore_order
bfstore_payment
bfstore_shipping
bfstore_notification
```

Deferred schemas:

```text
bfstore_auth
bfstore_customer
bfstore_review
bfstore_search
bfstore_recommendation
```

Each service should have:

```text
own schema
own database user
own migrations
own repository/data access layer
own seed data where required
```

## Rules

```text
no cross-service database joins
no cross-service foreign keys
no shared ORM models across services
cross-service references use IDs only
snapshots are allowed for historical accuracy
projections are allowed for optimised reads
```

## Risks

| Risk | Mitigation |
|---|---|
| Developers bypass APIs and query schemas directly | Enforce least privilege users |
| Too much duplicated data | Use snapshots only where justified |
| Reporting becomes hard | Build projections or analytics later |
| Migration overhead increases | Use consistent migration standards |

## Review Triggers

Revisit this decision if:

```text
service boundaries prove incorrect
most features require cross-service joins
operational overhead becomes too high
a modular monolith becomes more appropriate
```

## Related Documents

```text
docs/data/data-ownership.md
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/architecture/service-boundaries.md
```
