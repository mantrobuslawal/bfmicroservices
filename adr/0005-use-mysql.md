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
categories
category-scoped product attributes
baskets
stock levels
stock reservations
orders
payments
shipments
notifications
reviews
```

The catalogue must support varied product types such as homeware, curtains, bed frames, sofas, rugs, lamps, and wardrobes.

The user has chosen MySQL for the project. MySQL is suitable for transactional commerce-style workloads and is widely used in production systems.

The design should not force all product-specific fields into one large rigid `products` table. It should also avoid uncontrolled schemaless product blobs with weak governance.

## Decision

bfstore will use MySQL as the primary relational database.

Each service will own its own logical MySQL schema.

Catalogue Service will use MySQL as the governed product source of truth, using:

```text
core product tables
category taxonomy
product variants
category-scoped product attribute definitions
product attribute values
```

Search Service will own denormalised product search projections for browse, filters, and faceted search.

## Drivers

This decision supports:

```text
transactional consistency inside service boundaries
familiar relational modelling
straightforward local development
widely available managed cloud options
clear service-owned schemas
governed product catalogue data
flexible product attributes without uncontrolled schema sprawl
```

## Alternatives Considered

## Option 1: PostgreSQL

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

## Option 2: MySQL with Relational Core and Flexible Attributes

Benefits:

```text
widely used
good transactional support
simple local setup
suitable for commerce data
supports governed product attributes
works well with service-owned schemas
```

Costs:

```text
attribute modelling requires discipline
filter-heavy browse should eventually move to Search Service
care needed with constraints, migrations, and indexing
```

## Option 3: Document Database for Catalogue Source of Truth

Benefits:

```text
natural fit for varied product document shapes
flexible product attributes
whole-product reads can be simple
```

Costs:

```text
schema governance must be enforced outside the database
category-specific validation still required
faceted search still usually needs dedicated indexing/projection
less aligned with current MySQL project decision
```

## Option 4: One Rigid Product Table

Benefits:

```text
simple initial SQL
easy to understand for very small catalogues
```

Costs:

```text
many nullable product-type-specific columns
poor long-term maintainability
weak fit for varied furniture/homeware catalogue
hard to evolve categories and attributes
```

## Consequences

## Positive

```text
clear relational model for core services
works well with service-owned schemas
easy to run in Docker Compose
good fit for migration-based schema management
Catalogue Service remains governed and auditable
Search Service can optimise browse/filter behaviour separately
```

## Negative

```text
catalogue attribute modelling is more complex than a simple products table
search may require a separate projection or search engine later
advanced analytical queries should not be forced into service databases
care is required for money, timestamps, constraints, and concurrency
```

## Implementation Notes

Use MySQL for:

```text
catalogue data
category taxonomy
product attribute definitions
product attribute values
inventory stock and reservations
basket data
orders and order items
payment attempts
shipments
notification attempts
reviews
```

Catalogue tables should include:

```text
categories
products
product_variants
product_attribute_definitions
product_attribute_values
product_attribute_options
product_images
product_price_history
```

Potentially rebuildable projections may also use MySQL initially:

```text
search projections
recommendation signals
```

These can later move to specialised stores if needed.

## Catalogue Modelling Rule

Use:

```text
relational core product data
category-scoped attribute definitions
typed attribute values
search projection for customer-facing browse/filter/search
```

Avoid:

```text
one products table with hundreds of nullable type-specific columns
uncontrolled product JSON as the only model
Search Service becoming the product source of truth
```

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
category-scoped product attributes
denormalised search projections for browse/filter/search
```

## Risks

| Risk | Mitigation |
|---|---|
| Cross-service database coupling | Use schema ownership and least privilege |
| Money stored incorrectly | Use `amount_minor` and `currency_code` |
| Migration failures | Use tested, versioned migrations |
| Product attributes become ungoverned | Use attribute definitions, data types, and category ownership |
| Product table becomes too wide | Keep type-specific fields in attribute values |
| Search requirements outgrow MySQL | Add Search Service projection/search engine later |
| Inventory concurrency issues | Use transactions and locking/constraints carefully |

## Review Triggers

Revisit this decision if:

```text
catalogue product data becomes mostly document-shaped with limited relational value
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
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
proto/acme/catalog/v1/README.md
db/catalog/migrations/README.md
```

## Summary

bfstore will continue to use MySQL as the primary relational database.

For Catalogue Service, MySQL remains the governed source of truth, with flexible category-scoped product attributes. Search Service owns denormalised product projections for browse, filtering, and facets.

This gives bfstore product flexibility without prematurely moving the catalogue source of truth to NoSQL.
