# MySQL Standards

## 1. Purpose

This document defines MySQL standards for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It establishes conventions for schema naming, table design, column naming, data types, indexing, constraints, timestamps, transactions, permissions, flexible product attributes, and operational behaviour.

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
| Flexible | Product types can vary without uncontrolled schema sprawl |
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

---

## 5. Table Naming

Use plural lower snake case table names.

Good:

```text
products
product_variants
product_attribute_definitions
product_attribute_values
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
attribute_id
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
category_id
attribute_id
order_id
payment_id
shipment_id
notification_id
```

For internal child records:

```text
product_attribute_value_id
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

---

## 9. Timestamps

Use UTC timestamps.

Recommended columns:

```text
created_at
updated_at
deleted_at
```

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
categories
attribute definitions
customer addresses
notification templates
reviews
```

Avoid soft deletes on every table by default.

---

## 11. Money Representation

Do not use floating point types for money.

Use minor units.

Example:

```sql
amount_minor BIGINT NOT NULL,
currency_code CHAR(3) NOT NULL
```

Rules:

```text
use BIGINT for monetary amounts
use ISO-style currency code
avoid FLOAT and DOUBLE for money
document tax and discount handling separately
```

---

## 12. Flexible Product Attribute Standards

## 12.1 Principle

Catalogue Service should use relational core product tables plus category-scoped flexible attributes.

This supports varied product types without creating one large products table with many nullable columns.

Examples:

```text
curtains need drop, width, lining, heading type
bed frames need bed size, frame material, storage type
rugs need shape, pile height, weave, material
lamps need bulb type, wattage, fitting type
```

## 12.2 Recommended Tables

```text
categories
products
product_variants
product_attribute_definitions
product_attribute_values
product_attribute_options
```

## 12.3 Attribute Definitions

`product_attribute_definitions` should define which attributes apply to which categories.

Recommended fields:

```text
attribute_id
category_id
code
display_name
description
data_type
unit
is_required
is_filterable
is_variant_defining
allowed_values_json
display_order
status
created_at
updated_at
```

## 12.4 Attribute Values

`product_attribute_values` should store product-specific values.

Recommended fields:

```text
product_attribute_value_id
product_id
variant_id
attribute_id
value_string
value_number
value_boolean
value_json
unit
created_at
updated_at
```

Rules:

```text
only one typed value column should be populated per row
attribute definition determines expected type
variant_id should be nullable unless value differs by variant
attribute_id should reference an attribute definition within the catalogue schema
```

## 12.5 Attribute Data Types

Supported attribute data types should be controlled.

Recommended initial values:

```text
string
number
boolean
json
option
multi_option
```

Avoid uncontrolled free-form data types.

## 12.6 Attribute Codes

Attribute codes should be stable and machine-readable.

Good:

```text
drop_cm
width_cm
heading_type
bed_size
storage_type
bulb_type
rug_shape
material
colour
```

Avoid:

```text
Drop CM
Curtain Drop!!!
misc1
field_27
```

## 12.7 Controlled Attribute Values

For attributes used as filters, controlled values are preferred.

Examples:

```text
bed_size: single, double, king, super_king
heading_type: eyelet, pencil_pleat, tab_top
storage_type: none, drawer, ottoman
```

Controlled values may be represented using:

```text
product_attribute_options table
allowed_values_json field
```

For early implementation, `allowed_values_json` is acceptable if documented and tested.

## 12.8 Search Projection

Do not force Search Service behaviour into the Catalogue database.

Catalogue should govern product attributes.

Search Service should denormalise product data into browse/search documents.

Rules:

```text
Catalogue owns attribute truth
Search owns search documents and facets
filterable attributes should be clearly identified
search projections should be rebuildable
```

## 12.9 Anti-Patterns

Avoid:

```text
one products table with hundreds of nullable type-specific columns
uncontrolled JSON blobs for all product data
attributes with no category ownership
filterable fields hidden in ungoverned JSON
Search Service becoming the product source of truth
```

---

## 13. Quantities

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

---

## 14. Status Columns

Status columns should use stable string values or small enums.

Example:

```sql
status VARCHAR(32) NOT NULL
```

Examples:

```text
ACTIVE
INACTIVE
PENDING
CONFIRMED
CANCELLED
FAILED
```

---

## 15. Boolean Columns

Use boolean columns for true or false values.

Example:

```sql
is_filterable BOOLEAN NOT NULL DEFAULT FALSE
```

Rules:

```text
use is_ or has_ prefix
avoid nullable booleans unless three states are genuinely needed
```

---

## 16. JSON Columns

JSON columns are allowed, but should be used deliberately.

Acceptable uses:

```text
snapshots
provider metadata
allowed product attribute values in early implementation
event outbox payloads
audit context
selected attribute summaries
```

Avoid using JSON for core relational data that needs frequent filtering, joining, or constraints.

Good:

```text
delivery_address_snapshot_json
provider_response_summary_json
selected_attribute_summary_json
allowed_values_json
```

Risky:

```text
entire_product_json
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

## 17. Constraints

Use constraints to protect core invariants where practical.

Examples:

```sql
CHECK (quantity > 0)
CHECK (amount_minor >= 0)
UNIQUE (idempotency_key)
UNIQUE (order_number)
UNIQUE (category_id, code)
```

For product attributes, enforce uniqueness such as:

```text
one attribute definition code per category
one product-level value per product and attribute where appropriate
one variant-level value per variant and attribute where appropriate
```

---

## 18. Foreign Keys

## 18.1 Within Schema

Foreign keys are allowed within the same service-owned schema.

Examples:

```text
product_attribute_values.attribute_id -> product_attribute_definitions.attribute_id
product_variants.product_id -> products.product_id
order_items.order_id -> orders.order_id
```

## 18.2 Across Schemas

Foreign keys across service schemas are not allowed.

Avoid:

```text
bfstore_order.order_items.product_id -> bfstore_catalog.products.product_id
```

Use service APIs, events, snapshots, or projections instead.

---

## 19. Indexing Standards

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

Catalogue-specific indexes:

```text
idx_products_category_status
idx_product_variants_product_id
idx_attribute_definitions_category_code
idx_attribute_definitions_filterable
idx_attribute_values_product_id
idx_attribute_values_attribute_id
```

Search/filter-heavy customer queries should generally be served by Search Service rather than complex Catalogue SQL joins.

---

## 20. Index Naming

Use prefixes:

```text
pk_ for primary key where manually named
uq_ for unique indexes
idx_ for non-unique indexes
fk_ for foreign keys
```

Examples:

```text
uq_attribute_definitions_category_code
idx_attribute_values_product_id
idx_products_category_status
uq_orders_idempotency_key
```

---

## 21. Uniqueness and Idempotency

Critical operations should use uniqueness constraints where practical.

Examples:

```text
orders.idempotency_key
stock_reservations.idempotency_key
payments.idempotency_key
shipments.idempotency_key
processed_events.event_id
```

Catalogue examples:

```text
categories.slug
product_variants.sku
product_attribute_definitions(category_id, code)
```

---

## 22. Transactions

Use transactions for multi-step changes within one service schema.

Examples:

```text
create product + product variants + product attribute values
create order + order items + outbox event
reserve stock + reservation items + stock level update
record payment + payment attempt + outbox event
```

Rules:

```text
transactions must stay within one service schema
keep transactions short
avoid network calls inside database transactions
handle deadlocks and conflicts explicitly
```

---

## 23. Outbox Tables

Services that publish events after database writes should consider an outbox table.

Candidate services:

```text
order-service
payment-service
inventory-service
shipping-service
catalog-service
```

Catalogue outbox becomes useful when ProductUpdated events drive Search Service projections.

---

## 24. Audit and History Tables

Important stateful entities should record history.

Examples:

```text
order_status_history
payment_status_history
shipment_status_history
stock_adjustments
product_price_history
product_status_history
notification_attempts
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

## 28. Security and Privacy

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

Catalogue product attributes are generally not sensitive, but supplier-only information, margin data, cost prices, and unpublished product metadata should not be exposed through customer-facing APIs or search projections.

---

## 29. Testing Standards

Database tests should cover:

```text
migrations apply cleanly
constraints enforce expected rules
repository queries work
transactions commit and roll back correctly
idempotency uniqueness works
concurrency rules protect inventory
soft delete behaviour works where used
product attribute definitions constrain attribute values
filterable catalogue attributes can be projected to Search Service
```

---

## 30. Anti-Patterns to Avoid

Avoid:

```text
using FLOAT or DOUBLE for money
cross-service foreign keys
cross-service joins
nullable fields without reason
one products table with many nullable type-specific columns
generic JSON blobs for all catalogue data
shared database user for all services
migration scripts edited after release
unindexed high-volume lookup fields
storing raw payment details
secrets in seed data
```

---

## 31. Initial Implementation Standards

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

For Catalogue Service:

```text
create products
create categories
create product_variants
create product_attribute_definitions
create product_attribute_values
seed at least two different product types with different attributes
```

---

## 32. Open Questions

| Question | Status |
|---|---|
| Will standard IDs use ULID, UUIDv7, or UUIDv4? | To decide |
| Will timestamps use TIMESTAMP(6) or DATETIME(6)? | To decide |
| Should outbox payload be JSON or protobuf binary? | To decide |
| Should product attribute options use a table or JSON initially? | To decide |
| Should variant-specific attributes be supported in version one? | To decide |
| Should catalogue use generated columns for selected attributes? | Defer |
| Should Search Service use MySQL projection initially or external search later? | To decide |

---

## 33. Related Documents

```text
docs/data/data-ownership.md
docs/data/service-database-design.md
docs/data/migrations.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/events/event-catalog.md
```

Relevant ADRs:

```text
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
```

---

## 34. Summary

bfstore MySQL usage should be consistent, service-owned, secure, and migration-driven.

Catalogue Service remains on MySQL as the governed source of truth for product data. It supports varied product types using category-scoped attribute definitions and values, while Search Service owns denormalised browse/search projections.

This gives bfstore product flexibility without losing relational governance, data ownership, and client-reviewable design discipline.
