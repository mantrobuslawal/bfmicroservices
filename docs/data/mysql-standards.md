# MySQL Standards

## 1. Purpose

This document defines MySQL standards for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It establishes conventions for schema naming, table design, column naming, data types, indexing, constraints, timestamps, transactions, permissions, and operational behaviour.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s database design discipline.

---

## 2. Scope

This document applies to all MySQL schemas owned by bfstore services, including:

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
schema naming
table naming
column naming
data types
primary keys
foreign keys
indexes
timestamps
money representation
JSON usage
transactions
database users
local development
security
testing
```

---

## 3. Design Goals

bfstore MySQL usage should be:

| Goal | Description |
|---|---|
| Consistent | Schemas and tables follow predictable conventions |
| Service-owned | Each service owns its own schema |
| Secure | Database permissions follow least privilege |
| Auditable | Important changes can be traced |
| Performant | Indexes support known access patterns |
| Portable | Avoid unnecessary vendor-specific complexity |
| Maintainable | Migrations and naming are easy to review |
| Testable | Database behaviour can be validated in CI |

---

## 4. Schema Naming

Use the format:

```text
bfstore_<service_domain>
```

Examples:

```text
bfstore_catalog
bfstore_inventory
bfstore_basket
bfstore_order
bfstore_payment
bfstore_shipping
bfstore_notification
```

Rules:

```text
use lower snake case
include bfstore prefix
use business domain name
avoid environment names inside schema names unless platform requires it
```

Environment separation should ideally be handled by separate databases, instances, clusters, or configuration.

---

## 5. Table Naming

Use plural lower snake case table names.

Good:

```text
products
product_variants
stock_reservations
basket_items
orders
order_items
payment_attempts
shipments
notification_attempts
```

Avoid:

```text
Product
ProductVariant
tbl_products
orderItem
```

Tables should be named after business concepts, not implementation classes.

---

## 6. Column Naming

Use lower snake case.

Good:

```text
product_id
created_at
updated_at
idempotency_key
currency_code
amount_minor
```

Avoid:

```text
productID
CreatedAt
currencyCode
amountInPence
```

Column names should be explicit and business-readable.

---

## 7. Primary Keys

Each table should have a stable primary key.

Recommended:

```text
product_id
order_id
payment_id
shipment_id
notification_id
```

For internal child records:

```text
order_item_id
payment_attempt_id
notification_attempt_id
```

Rules:

```text
primary keys should be immutable
public contract IDs should not expose auto-increment database IDs
string IDs are preferred for service contracts
indexes should support lookup by public IDs
```

---

## 8. ID Data Types

Recommended ID column type:

```sql
VARCHAR(64)
```

Example:

```sql
product_id VARCHAR(64) NOT NULL PRIMARY KEY
```

This supports:

```text
ULID
UUID
UUIDv7
prefixed IDs
external provider references where appropriate
```

If binary UUIDs are later chosen for storage efficiency, the decision should be documented and hidden behind service APIs.

---

## 9. Timestamps

Use UTC timestamps.

Recommended columns:

```text
created_at
updated_at
deleted_at
```

Column type:

```sql
TIMESTAMP(6)
```

or:

```sql
DATETIME(6)
```

The final choice should be consistent across the project.

Recommended default pattern:

```sql
created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
```

Rules:

```text
store UTC
use consistent precision
include created_at on most tables
include updated_at on mutable tables
use deleted_at only if soft delete is needed
```

---

## 10. Soft Deletes

Soft deletes should be used only where there is a business need.

Potential use cases:

```text
products
customer addresses
notification templates
reviews
```

Avoid soft deletes on every table by default.

If soft delete is used:

```text
include deleted_at
ensure queries filter deleted records correctly
consider uniqueness constraints carefully
document retention behaviour
```

---

## 11. Money Representation

Do not use floating point types for money.

Use minor units.

Example:

```sql
amount_minor BIGINT NOT NULL,
currency_code CHAR(3) NOT NULL
```

Examples:

```text
129900 GBP = £1,299.00
4999 GBP = £49.99
```

Rules:

```text
use BIGINT for monetary amounts
use ISO currency code
avoid FLOAT and DOUBLE for money
document tax and discount handling separately
```

---

## 12. Quantities

Use integer types for quantities.

Example:

```sql
quantity INT NOT NULL
```

Rules:

```text
quantity must be greater than zero for basket/order items
available stock must not be negative
reserved stock must not be negative
use constraints where practical
```

Example:

```sql
CHECK (quantity > 0)
```

---

## 13. Status Columns

Status columns should use stable string values or small enums.

Example:

```sql
status VARCHAR(32) NOT NULL
```

Examples:

```text
PENDING
CONFIRMED
CANCELLED
FAILED
```

Rules:

```text
status values must be documented
state transitions should be enforced in application logic
important state changes should be recorded in history tables
avoid ambiguous statuses such as ACTIVE2 or TEMP
```

---

## 14. Boolean Columns

Use boolean columns for true or false values.

Example:

```sql
is_active BOOLEAN NOT NULL DEFAULT TRUE
```

Rules:

```text
use is_ or has_ prefix
avoid nullable booleans unless three states are genuinely needed
```

---

## 15. JSON Columns

JSON columns are allowed, but should be used deliberately.

Acceptable uses:

```text
snapshots
provider metadata
flexible product attributes
event outbox payloads
audit context
```

Avoid using JSON for core relational data that needs frequent filtering, joining, or constraints.

Good:

```text
delivery_address_snapshot_json
provider_response_summary_json
```

Risky:

```text
order_items_json
payment_state_json
```

Rules:

```text
JSON fields must have clear purpose
do not hide critical business fields inside JSON
avoid storing sensitive data unnecessarily
index generated columns only where justified
```

---

## 16. Constraints

Use constraints to protect core invariants where practical.

Examples:

```sql
CHECK (quantity > 0)
CHECK (amount_minor >= 0)
UNIQUE (idempotency_key)
UNIQUE (order_number)
```

Rules:

```text
use NOT NULL for required fields
use UNIQUE for natural uniqueness where required
use CHECK constraints for simple invariants
do not rely only on application validation for critical data integrity
```

---

## 17. Foreign Keys

## 17.1 Within Schema

Foreign keys are allowed within the same service-owned schema.

Example:

```sql
order_items.order_id -> orders.order_id
```

## 17.2 Across Schemas

Foreign keys across service schemas are not allowed.

Avoid:

```sql
bfstore_order.order_items.product_id
    -> bfstore_catalog.products.product_id
```

Use service APIs, events, or snapshots instead.

---

## 18. Indexing Standards

Indexes should be created for known access patterns.

Common indexes:

```text
created_at
updated_at
customer_id
order_id
basket_id
product_id
status
idempotency_key
event_id
```

Examples:

```sql
CREATE INDEX idx_orders_customer_created_at
ON orders (customer_id, created_at);

CREATE UNIQUE INDEX uq_orders_idempotency_key
ON orders (idempotency_key);
```

Rules:

```text
index names should be clear
avoid unused indexes
avoid over-indexing write-heavy tables
review query plans for important queries
include indexes in migrations
```

---

## 19. Index Naming

Use prefixes:

```text
pk_ for primary key where manually named
uq_ for unique indexes
idx_ for non-unique indexes
fk_ for foreign keys
```

Examples:

```text
pk_orders
uq_orders_order_number
uq_orders_idempotency_key
idx_orders_customer_created_at
fk_order_items_order_id
```

---

## 20. Uniqueness and Idempotency

Critical operations should use uniqueness constraints where practical.

Examples:

```text
orders.idempotency_key
stock_reservations.idempotency_key
payments.idempotency_key
shipments.idempotency_key
processed_events.event_id
```

Rules:

```text
same idempotency key should not create duplicate business effects
same event_id should not be processed twice by a consumer
same key with different request hash should be rejected
```

---

## 21. Transactions

Use transactions for multi-step changes within one service schema.

Examples:

```text
create order + order items + outbox event
reserve stock + stock reservation items + stock level update
record payment + payment attempt + outbox event
record notification + notification attempt
```

Rules:

```text
transactions must stay within one service schema
keep transactions short
avoid network calls inside database transactions
handle deadlocks and conflicts explicitly
```

---

## 22. Isolation and Concurrency

Inventory and checkout flows require careful concurrency handling.

Potential techniques:

```text
row-level locking
optimistic concurrency
unique constraints
idempotency records
transaction retries for deadlocks
```

Inventory must prevent overselling.

Example invariant:

```text
available_quantity >= 0
reserved_quantity >= 0
```

---

## 23. Outbox Tables

Services that publish events after database writes should consider an outbox table.

Candidate table:

```text
outbox_events
```

Candidate fields:

```text
outbox_event_id
event_id
event_type
event_version
aggregate_type
aggregate_id
payload
status
attempt_count
next_attempt_at
created_at
published_at
last_error
```

Rules:

```text
outbox row should be written in same transaction as business data
event_id must remain stable across publish retries
publisher should update status after successful publish
outbox backlog should be observable
```

---

## 24. Audit and History Tables

Important stateful entities should record history.

Examples:

```text
order_status_history
payment_status_history
shipment_status_history
stock_adjustments
notification_attempts
```

History tables should include:

```text
previous_status
new_status
reason
changed_at
changed_by or source where applicable
correlation_id
```

---

## 25. Database Users and Permissions

Each service should have its own database user.

Examples:

```text
bfstore_catalog_user
bfstore_inventory_user
bfstore_order_user
bfstore_payment_user
```

Rules:

```text
service users access only their own schema
no broad application user for all schemas
migration users may have elevated permissions but should be controlled
read-only users should be separate where needed
```

Example:

```sql
GRANT SELECT, INSERT, UPDATE, DELETE
ON bfstore_order.*
TO 'bfstore_order_user'@'%';
```

---

## 26. Local Development Standards

Local development may use:

```text
one MySQL container
multiple schemas
service-specific users
seed data
migration scripts
```

Local setup should create:

```text
schemas
users
permissions
initial seed data where needed
```

Expected scripts:

```text
db/mysql-init/001-create-databases.sql
db/mysql-init/002-create-users.sql
db/mysql-init/003-grant-permissions.sql
```

---

## 27. Migration Standards

Migrations should be:

```text
versioned
reviewed
repeatable in CI
owned by the service
safe for deployment
documented when destructive
```

Migration details are defined in:

```text
docs/data/migrations.md
```

---

## 28. Backup and Restore Considerations

Data criticality differs by service.

High criticality:

```text
orders
payments
customers
inventory
shipments
```

Potentially rebuildable:

```text
search projections
recommendation projections
analytics projections
```

Backup and restore expectations should be documented per service.

---

## 29. Security and Privacy

Do not store:

```text
raw payment card data
plaintext passwords
secret values
authentication tokens unless specifically required and protected
```

Minimise storage of:

```text
customer email
phone number
delivery address
provider references
notification recipient details
```

Sensitive columns should be identified in data classification documents.

---

## 30. Testing Standards

Database tests should cover:

```text
migrations apply cleanly
constraints enforce expected rules
repository queries work
transactions commit and roll back correctly
idempotency uniqueness works
concurrency rules protect inventory
soft delete behaviour works where used
```

Tests should use:

```text
Docker Compose MySQL
Testcontainers
isolated schemas
transaction rollbacks or database reset
```

---

## 31. Anti-Patterns to Avoid

Avoid:

```text
using FLOAT or DOUBLE for money
cross-service foreign keys
cross-service joins
nullable fields without reason
generic JSON blobs for core data
shared database user for all services
migration scripts edited after release
unindexed high-volume lookup fields
storing raw payment details
secrets in seed data
```

---

## 32. Initial Implementation Standards

For the first checkout vertical slice:

```text
create separate schemas for active services
create service-specific users
use string IDs
use amount_minor and currency_code for money
use created_at and updated_at consistently
add idempotency_key columns for critical commands
add processed_events table for Notification Service
add outbox_events for Order Service if outbox is implemented
```

---

## 33. Open Questions

| Question | Status |
|---|---|
| Will standard IDs use ULID, UUIDv7, or UUIDv4? | To decide |
| Will timestamps use TIMESTAMP(6) or DATETIME(6)? | To decide |
| Should outbox payload be JSON or protobuf binary? | To decide |
| Should address snapshots be JSON or structured columns? | To decide |
| Will local dev users exactly mirror production-style permissions? | Proposed |
| Should migrations be run by service startup or a separate job? | To decide |

---

## 34. Related Documents

This document should be read alongside:

```text
docs/data/data-ownership.md
docs/data/service-database-design.md
docs/data/migrations.md
docs/data/pii-handling.md
docs/data/retention.md
docs/architecture/service-boundaries.md
docs/architecture/resilience-patterns.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
```

---

## 35. Summary

bfstore MySQL usage should be consistent, service-owned, secure, and migration-driven.

The most important standards are:

```text
one schema per service
one database user per service
no cross-service joins
no cross-service foreign keys
string IDs for contracts
minor units for money
UTC timestamps
clear indexes
safe constraints
service-owned migrations
least privilege permissions
```

These standards support a professional and maintainable data layer for bfstore’s microservice architecture.
