# ADR-0005: Use MySQL as the Primary Relational Database

## Status

Accepted

## Date

2026-05-26

## Context

bfstore requires durable transactional storage for commerce data.

Core data includes:

```text
products
baskets
stock levels
stock reservations
orders
payments
shipments
notifications
reviews
```

The user has chosen MySQL for the project. MySQL is suitable for transactional commerce-style workloads and is widely used in production systems.

## Decision

bfstore will use MySQL as the primary relational database.

Each service will own its own logical MySQL schema.

## Drivers

This decision supports:

```text
transactional consistency inside service boundaries
familiar relational modelling
straightforward local development
widely available managed cloud options
good fit for orders, payments, inventory, and catalogue data
```

## Alternatives Considered

### Option 1: PostgreSQL

Benefits:

```text
rich feature set
strong SQL capabilities
excellent JSON and indexing options
popular for modern backend systems
```

Costs:

```text
not the selected database for this project
would not materially change the main architecture principles
```

### Option 2: MySQL

Benefits:

```text
widely used
good transactional support
simple local setup
suitable for commerce data
cloud managed options available
```

Costs:

```text
some advanced PostgreSQL features are not available
care needed with data types, constraints, and migrations
```

### Option 3: NoSQL Database

Benefits:

```text
flexible document modelling
can be useful for selected high-scale access patterns
```

Costs:

```text
less suitable for initial transactional order/payment/inventory model
would complicate relational integrity expectations
```

## Consequences

### Positive

```text
clear relational model for core services
works well with service-owned schemas
easy to run in Docker Compose
good fit for migration-based schema management
```

### Negative

```text
search may require a separate projection or search engine later
advanced analytical queries should not be forced into service databases
care is required for money, timestamps, constraints, and concurrency
```

## Implementation Notes

Use MySQL for:

```text
catalogue data
inventory stock and reservations
basket data
orders and order items
payment attempts
shipments
notification attempts
reviews
```

Potentially rebuildable projections may also use MySQL initially:

```text
search projections
recommendation signals
```

These can later move to more specialised stores if needed.

## Standards

bfstore should use:

```text
service-owned schemas
service-specific database users
versioned migrations
minor units for money
UTC timestamps
string IDs for service contracts
constraints for core invariants
indexes based on known access patterns
```

## Risks

| Risk | Mitigation |
|---|---|
| Cross-service database coupling | Use schema ownership and least privilege |
| Money stored incorrectly | Use `amount_minor` and `currency_code` |
| Migration failures | Use tested, versioned migrations |
| Search requirements outgrow MySQL | Add Search Service projection/search engine later |
| Inventory concurrency issues | Use transactions and locking/constraints carefully |

## Review Triggers

Revisit this decision if:

```text
search or analytics requirements dominate
cloud/platform constraints favour another database
MySQL limitations block important behaviour
specialised storage is needed for a specific service
```

## Related Documents

```text
docs/data/mysql-standards.md
docs/data/service-database-design.md
docs/data/migrations.md
docs/data/data-ownership.md
```
