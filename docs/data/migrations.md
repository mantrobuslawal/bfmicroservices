# Database Migrations

## 1. Purpose

This document defines the database migration strategy for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how schema changes should be created, reviewed, tested, applied, rolled back, and coordinated with service deployments.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s database delivery and operational maturity.

---

## 2. Scope

This document applies to migrations for all service-owned MySQL schemas:

```text
bfstore_auth
bfstore_customer
bfstore_catalog
bfstore_inventory
bfstore_basket
bfstore_order
bfstore_payment
bfstore_shipping
bfstore_notification
bfstore_review
bfstore_search
bfstore_recommendation
```

It covers:

```text
migration ownership
migration file layout
versioning
up and down migrations
CI testing
local development
deployment coordination
rollback strategy
expand-and-contract changes
destructive changes
seed data
migration tooling expectations
```

---

## 3. Migration Principles

## 3.1 Each Service Owns Its Migrations

A service owns migrations for its own schema.

Example:

```text
catalog-service owns bfstore_catalog migrations
order-service owns bfstore_order migrations
payment-service owns bfstore_payment migrations
```

No service should run migrations for another service’s schema.

---

## 3.2 Migrations Are Code

Migrations should be treated as production code.

They must be:

```text
versioned
reviewed
tested
repeatable
auditable
documented where risky
```

---

## 3.3 Prefer Backward-Compatible Changes

Database changes should support safe deployment and rollback.

Prefer:

```text
add nullable column
add table
add index carefully
add new optional field
write both old and new fields temporarily
```

Avoid:

```text
drop column immediately
rename column without transition
change type in place
tighten constraints before data is clean
remove table used by current service version
```

---

## 3.4 Migrations Must Respect Service Boundaries

Migrations must not create cross-service coupling.

Avoid:

```text
foreign keys across service schemas
views joining service schemas
stored procedures that update multiple service schemas
shared lookup tables used by multiple services
```

---

## 4. Recommended Directory Layout

Recommended layout:

```text
db/
├── README.md
├── mysql-init/
│   ├── 001-create-databases.sql
│   ├── 002-create-users.sql
│   └── 003-grant-permissions.sql
├── catalog/
│   └── migrations/
├── inventory/
│   └── migrations/
├── basket/
│   └── migrations/
├── order/
│   └── migrations/
├── payment/
│   └── migrations/
├── shipping/
│   └── migrations/
├── notification/
│   └── migrations/
├── review/
│   └── migrations/
├── search/
│   └── migrations/
└── recommendation/
    └── migrations/
```

Alternative service-local layout:

```text
services/
└── order-service/
    └── migrations/
```

## 4.1 Initial Recommendation

For the bfstore repo, keep migrations under:

```text
db/<service-domain>/migrations/
```

This makes service schema ownership visible at the repository level.

---

## 5. Migration File Naming

Use sequential and descriptive migration names.

Example:

```text
000001_create_products.up.sql
000001_create_products.down.sql
000002_create_product_variants.up.sql
000002_create_product_variants.down.sql
000003_add_product_status_index.up.sql
000003_add_product_status_index.down.sql
```

Rules:

```text
use zero-padded numbers
use lower snake case
include action and target
separate up and down migrations if tooling supports it
do not edit released migrations
create new migrations for changes
```

---

## 6. Up and Down Migrations

## 6.1 Up Migration

An up migration applies the schema change.

Example:

```sql
CREATE TABLE orders (
  order_id VARCHAR(64) NOT NULL PRIMARY KEY,
  customer_id VARCHAR(64) NOT NULL,
  status VARCHAR(32) NOT NULL,
  total_amount_minor BIGINT NOT NULL,
  currency_code CHAR(3) NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
);
```

## 6.2 Down Migration

A down migration reverses the schema change where safe.

Example:

```sql
DROP TABLE orders;
```

## 6.3 Down Migration Caution

Some migrations cannot be fully reversed safely.

Examples:

```text
data deletion
column drop
data type conversion
data backfill with lossy transformation
```

When a down migration is unsafe, document it clearly.

---

## 7. Migration Tooling

The final migration tool is to be decided.

Suitable options include:

```text
golang-migrate
goose
Atlas
Flyway
Liquibase
```

## 7.1 Tool Selection Criteria

The tool should support:

```text
MySQL
versioned migrations
CI execution
local development
up and down migrations
clear failure behaviour
automation through Makefile or scripts
```

## 7.2 Initial Recommendation

Use a Go-friendly migration tool such as:

```text
golang-migrate
```

or:

```text
goose
```

Either is suitable for a Go microservice portfolio project.

---

## 8. Migration Execution Model

There are several ways to run migrations.

## 8.1 Local Development

Run through Makefile or scripts.

Example:

```sh
make migrate-up
make migrate-down
```

## 8.2 CI

Run migrations against a temporary MySQL instance.

CI should verify:

```text
migrations apply from empty schema
migrations can run before integration tests
repository tests pass after migrations
down migrations work where required
```

## 8.3 Kubernetes

Possible approaches:

```text
run migrations as Kubernetes Jobs
run migrations in CI/CD before deployment
run migrations manually for controlled environments
run migrations at service startup only if carefully controlled
```

## 8.4 Recommendation

For production-style environments, prefer:

```text
migration job or CI/CD controlled migration step
```

Avoid uncontrolled migrations from every service pod at startup.

---

## 9. Migration Ordering

Migrations are ordered within each service schema.

There should be no ordering dependency between service schemas unless documented.

Example:

```text
catalog migrations should not depend on order migrations
order migrations should not depend on payment migrations
```

If application behaviour requires multiple services to change, coordinate through API/event versioning and deployment sequencing, not cross-schema migrations.

---

## 10. Expand and Contract Pattern

Use expand and contract for safer schema changes.

## 10.1 Example: Renaming a Column

Unsafe direct change:

```text
rename column product_name to display_name
deploy service expecting display_name
rollback becomes difficult
```

Safer approach:

```text
1. Add display_name column.
2. Deploy service writing both product_name and display_name.
3. Backfill display_name.
4. Deploy service reading display_name.
5. Stop writing product_name.
6. Drop product_name in later release.
```

## 10.2 Benefits

```text
supports rolling deployments
supports rollback
reduces release risk
allows old and new service versions to coexist
```

---

## 11. Destructive Changes

Destructive changes include:

```text
drop table
drop column
change column type with data loss
tighten nullable column to NOT NULL
delete large volumes of data
remove indexes needed by current queries
```

Rules:

```text
avoid destructive changes in the same release as application changes
use deprecation period
take backup where appropriate
document rollback limitations
test in staging-like environment
```

---

## 12. Index Migrations

Adding indexes can affect production databases.

Rules:

```text
create indexes for known access patterns
consider table size
avoid unnecessary indexes
test query plans
avoid blocking operations where possible
monitor migration duration
```

Index naming:

```text
idx_orders_customer_created_at
uq_orders_idempotency_key
```

---

## 13. Data Backfills

Backfills should be treated carefully.

Rules:

```text
avoid large unbounded updates in one transaction
batch large backfills
make backfills resumable where practical
log progress
test on realistic data volume
monitor database load
```

Example:

```text
backfill order item product snapshots
backfill new status field
backfill search projection table
```

---

## 14. Seed Data

Seed data is useful for local development and tests.

Seed data should include:

```text
test products
test stock levels
test delivery options
notification templates
simulated payment cases
```

Seed data must not contain real personal data.

Separate:

```text
local development seed data
test fixtures
production reference data
```

Do not mix them.

---

## 15. Migration Testing

Migrations should be tested in CI.

Minimum tests:

```text
apply all migrations from empty database
run repository tests after migrations
validate expected tables and indexes exist
validate constraints enforce expected rules
run down migrations where supported
re-apply migrations after down where supported
```

Advanced tests:

```text
migration from previous release snapshot
rollback compatibility test
performance test for large migration
backfill test
```

---

## 16. Migration Review Checklist

Before approving a migration, check:

```text
Which service owns this migration?
Does it affect only the service-owned schema?
Is the file name clear and versioned?
Is the migration backward-compatible?
Is there a down migration?
Are destructive changes avoided or documented?
Are indexes justified by access patterns?
Are constraints correct?
Could this block deployment?
Could this break rollback?
Are tests updated?
Are seed data changes safe?
Is sensitive data handled correctly?
```

---

## 17. Rollback Strategy

Rollback has two parts:

```text
application rollback
database rollback
```

Application rollback is usually easier than database rollback.

Database rollback may be unsafe if data has changed.

## 17.1 Safer Rollback Approach

Prefer schema changes that allow old and new application versions to run.

Examples:

```text
adding nullable columns
adding new tables
adding optional indexes
adding new status values carefully
```

## 17.2 Risky Rollback Scenarios

```text
dropped column needed by old service
renamed column without compatibility
changed enum/status assumptions
changed data format in place
removed index required by old query path
```

---

## 18. Deployment Coordination

Database migrations and service deployments must be coordinated.

Recommended order for backward-compatible changes:

```text
1. Apply migration.
2. Deploy service version that can use new schema.
3. Verify health and smoke tests.
4. Later remove old schema elements if needed.
```

For breaking changes:

```text
1. Avoid if possible.
2. Introduce new version.
3. Migrate consumers.
4. Remove old version later.
```

---

## 19. Kubernetes Migration Jobs

If using Kubernetes, migrations may run as Jobs.

Example responsibilities:

```text
connect to service schema
apply pending migrations
fail clearly on error
emit logs
exit successfully when complete
```

Rules:

```text
only one migration job should run for a schema at a time
migration job should use controlled credentials
migration job should be observable
failed migration should block deployment promotion
```

---

## 20. Migration Permissions

Application runtime users and migration users may differ.

## 20.1 Runtime User

Runtime user permissions:

```text
SELECT
INSERT
UPDATE
DELETE
```

on the service schema.

## 20.2 Migration User

Migration user may require:

```text
CREATE
ALTER
DROP
INDEX
REFERENCES within schema
```

Migration credentials should be more tightly controlled.

---

## 21. Local Development Commands

Expected commands:

```sh
make migrate-up
make migrate-down
make migrate-status
make migrate-reset
```

Service-specific commands may include:

```sh
make migrate-up SERVICE=order
make migrate-down SERVICE=payment
```

Reset commands should be clearly marked as destructive.

---

## 22. Migration Failure Handling

If a migration fails:

```text
stop deployment
inspect migration logs
check partial application state
restore from backup if required
repair migration state carefully
do not blindly rerun destructive migrations
document incident if environment is shared
```

Migration tools usually track applied versions. Manual edits to migration tracking tables should be avoided unless clearly understood.

---

## 23. Observability

Migration execution should produce:

```text
migration version
service/schema
start time
end time
duration
success/failure
error details
operator or pipeline run ID
```

For production-style environments, migration results should be visible in CI/CD logs and deployment records.

---

## 24. Initial Migration Scope

The first migration set should create schemas for the checkout vertical slice:

```text
catalog
inventory
basket
order
payment
shipping
notification
```

Initial tables:

```text
products
product_variants
stock_levels
stock_reservations
baskets
basket_items
orders
order_items
payments
payment_attempts
shipments
notifications
notification_attempts
processed_events
```

Optional early table:

```text
outbox_events
```

especially for `order-service`.

---

## 25. Anti-Patterns to Avoid

Avoid:

```text
editing migrations after they have been applied
one giant migration for all services
cross-service schema changes
runtime pods racing to apply migrations
destructive changes without transition
manual production changes not represented in Git
seed data containing real personal data
using migrations to move data between service-owned schemas
```

---

## 26. Open Questions

| Question | Status |
|---|---|
| Which migration tool will be used? | To decide |
| Will migrations live under `db/<service>/migrations` or service directories? | Proposed: `db/<service>/migrations` |
| Will Kubernetes migrations run as Jobs or through CI/CD? | To decide |
| Will application startup ever run migrations? | Proposed: local only, not production-style |
| Should outbox tables be included from the first migration set? | Proposed for `order-service` |
| What migration rollback policy applies to staging and production-style environments? | To decide |
| Will migration users be separate from runtime users locally? | Proposed |

---

## 27. Related Documents

This document should be read alongside:

```text
docs/data/data-ownership.md
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/architecture/deployment-view.md
docs/architecture/resilience-patterns.md
docs/testing/testing-strategy.md
docs/operations/database-migration-strategy.md
docs/operations/rollback-strategy.md
```

Relevant ADRs:

```text
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
```

---

## 28. Summary

bfstore database migrations should be service-owned, versioned, tested, and deployment-aware.

The most important rules are:

```text
each service owns its migrations
migrations affect only service-owned schemas
prefer backward-compatible changes
do not edit released migrations
test migrations in CI
avoid destructive changes without transition
coordinate migrations with deployments
use expand-and-contract for risky changes
control migration permissions
```

This migration strategy supports safer releases, clearer rollback planning, and professional database change management.
