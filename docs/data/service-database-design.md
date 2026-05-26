# Service Database Design

## 1. Purpose

This document defines the service database design approach for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how each service owns its data, how MySQL schemas are separated, how tables should be designed, how cross-service data references should work, and how the catalogue data model supports flexible product types such as homeware, curtains, bed frames, rugs, lamps, sofas, and wardrobes.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s data architecture and microservice design maturity.

---

## 2. Scope

This document covers:

```text
service-owned database design
logical MySQL schema separation
initial service schemas
table ownership
cross-service references
snapshots
projections
flexible catalogue attributes
transaction boundaries
database users
migration ownership
local development database layout
initial checkout vertical slice data model
```

It should be read alongside:

```text
docs/data/data-ownership.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/architecture/service-boundaries.md
docs/architecture/domain-model.md
```

---

## 3. Core Design Principle

The central rule is:

> Each service owns its own database schema. Other services must access data through APIs, events, or projections, not direct database access.

This means:

```text
catalog-service owns bfstore_catalog
inventory-service owns bfstore_inventory
basket-service owns bfstore_basket
order-service owns bfstore_order
payment-service owns bfstore_payment
shipping-service owns bfstore_shipping
notification-service owns bfstore_notification
```

A service must not directly query another service’s tables.

---

## 4. Source of Truth vs Projection

bfstore separates the product source of truth from optimised read/search projections.

```text
Catalogue Service
  owns governed product data in MySQL

Search Service
  owns denormalised product search documents and facets

Recommendation Service
  owns recommendation signals and outputs
```

This allows the catalogue to remain governed and auditable while allowing the search representation to be flexible, denormalised, and optimised for browse/filter behaviour.

---

## 5. Logical Database Model

In local development, bfstore may run one MySQL instance with multiple schemas.

Example:

```text
MySQL instance
├── bfstore_auth
├── bfstore_customer
├── bfstore_catalog
├── bfstore_inventory
├── bfstore_basket
├── bfstore_order
├── bfstore_payment
├── bfstore_shipping
├── bfstore_notification
├── bfstore_review
├── bfstore_search
└── bfstore_recommendation
```

In production-style environments, these schemas may later become separate databases or managed instances.

The logical ownership rule remains the same.

---

## 6. Schema Ownership Summary

| Service | Schema | Purpose |
|---|---|---|
| `auth-service` | `bfstore_auth` | Identity, credentials, sessions, roles |
| `customer-service` | `bfstore_customer` | Customer profiles, addresses, preferences |
| `catalog-service` | `bfstore_catalog` | Products, categories, variants, category-scoped attributes, catalogue pricing |
| `inventory-service` | `bfstore_inventory` | Stock levels, stock reservations, warehouse availability |
| `basket-service` | `bfstore_basket` | Baskets and basket items |
| `order-service` | `bfstore_order` | Orders, order items, order state |
| `payment-service` | `bfstore_payment` | Payments, attempts, refunds |
| `shipping-service` | `bfstore_shipping` | Shipments, delivery options, tracking |
| `notification-service` | `bfstore_notification` | Notifications, attempts, templates |
| `review-service` | `bfstore_review` | Reviews, ratings, moderation |
| `search-service` | `bfstore_search` | Search projections, facets, search index state |
| `recommendation-service` | `bfstore_recommendation` | Recommendation signals and results |

---

## 7. Database Access Rules

## 7.1 Allowed

A service may:

```text
read its own schema
write its own schema
run migrations for its own schema
manage seed data for its own schema
create indexes in its own schema
own transaction boundaries inside its own schema
```

## 7.2 Forbidden

A service must not:

```text
read another service's tables
write another service's tables
join across service schemas
share ORM models with another service
depend on another service's internal table design
```

Forbidden example:

```sql
SELECT *
FROM bfstore_order.orders o
JOIN bfstore_customer.customers c ON c.customer_id = o.customer_id;
```

Correct alternatives:

```text
use Customer Service API
consume CustomerUpdated event into a projection
store an order-time customer/address snapshot where required
```

---

## 8. Cross-Service References

Services may store references to entities owned by other services.

Example:

```text
order-service stores customer_id
order-service stores product_id
payment-service stores order_id
shipping-service stores order_id
basket-service stores product_id
inventory-service stores product_id
```

Rules:

```text
a foreign ID is not ownership
do not create database-level foreign keys across schemas
validate cross-service references through APIs where needed
use events to keep projections updated
store snapshots where historical accuracy is required
```

---

## 9. Foreign Key Strategy

## 9.1 Within a Service Schema

Foreign keys are allowed within a service-owned schema.

Example:

```text
bfstore_order.orders
bfstore_order.order_items
```

`order_items.order_id` may reference `orders.order_id`.

## 9.2 Across Service Schemas

Foreign keys across service schemas are not allowed.

Avoid:

```text
bfstore_order.order_items.product_id -> bfstore_catalog.products.product_id
```

Reason:

```text
cross-schema foreign keys create deployment coupling
service schemas cannot evolve independently
database availability becomes shared
ownership becomes unclear
```

---

## 10. Identifier Strategy

Use string IDs for service contracts and stored cross-service references.

Recommended ID formats:

```text
ULID
UUIDv7
UUIDv4
```

Example IDs:

```text
prd_01HX...
cat_01HX...
attr_01HX...
cus_01HX...
bas_01HX...
ord_01HX...
pay_01HX...
shp_01HX...
```

Rules:

```text
do not expose database auto-increment IDs as public contracts
use stable business/entity IDs in events and APIs
indexes should support frequent lookup by entity ID
```

---

## 11. Initial Implementation Schemas

The first implementation should focus on the checkout vertical slice.

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

Search Service may be introduced after catalogue events are stable.

---

## 12. Catalogue Schema

## 12.1 Purpose

The Catalogue schema stores governed product information owned by `catalog-service`.

It must support many product types without creating one large `products` table containing hundreds of nullable product-type-specific columns.

Examples of future product types:

```text
curtains
bed frames
mattresses
sofas
rugs
lamps
wardrobes
tables
mirrors
cushions
homeware
```

## 12.2 Design Approach

Use a relational core product model plus category-scoped flexible attributes.

```text
Core product tables
  store fields common to most products

Category and attribute tables
  define which attributes apply to which product categories

Attribute value tables
  store product-specific values

Search projection
  denormalises product and attribute data for browse/filter/search
```

This avoids both extremes:

```text
one huge products table with many nullable columns
uncontrolled schemaless product documents with weak governance
```

## 12.3 Candidate Tables

Initial priority:

```text
products
product_variants
categories
product_attribute_definitions
product_attribute_values
```

Later:

```text
product_attribute_options
product_images
product_price_history
outbox_events
```

## 12.4 `products`

Candidate fields:

```text
product_id
category_id
name
description
status
base_price_minor
currency_code
brand
created_at
updated_at
```

The `products` table should contain common product data only.

Avoid placing highly specific fields here, such as:

```text
curtain_drop_cm
bed_size
bulb_type
rug_shape
mattress_firmness
sofa_orientation
```

These belong in product attributes.

## 12.5 `product_variants`

Candidate fields:

```text
variant_id
product_id
sku
variant_name
price_minor
currency_code
status
created_at
updated_at
```

Variant-specific attribute values may be supported later if required.

Examples:

```text
same bed frame in double, king, and super king
same curtains in multiple drops and widths
same sofa in different fabrics or orientations
```

## 12.6 `categories`

Candidate fields:

```text
category_id
parent_category_id
name
slug
description
status
created_at
updated_at
```

Categories define the taxonomy for products and attribute rules.

Example categories:

```text
curtains
bed-frames
rugs
lamps
sofas
wardrobes
```

## 12.7 `product_attribute_definitions`

This table defines attributes that are valid for a category.

Candidate fields:

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

Example rows:

| Category | Code | Data Type | Unit | Filterable |
|---|---|---|---|---|
| curtains | `drop_cm` | number | cm | yes |
| curtains | `heading_type` | string | none | yes |
| bed-frames | `bed_size` | string | none | yes |
| bed-frames | `storage_type` | string | none | yes |
| rugs | `shape` | string | none | yes |
| lamps | `bulb_type` | string | none | yes |

## 12.8 `product_attribute_values`

This table stores product-specific attribute values.

Candidate fields:

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
only one value column should be populated according to attribute data type
attribute_id determines expected type and unit
variant_id is optional and used only when value differs per variant
```

Example values:

| Product Type | Attribute | Value |
|---|---|---|
| Curtain | `drop_cm` | `228` |
| Curtain | `heading_type` | `eyelet` |
| Bed frame | `bed_size` | `king` |
| Bed frame | `storage_type` | `ottoman` |
| Rug | `shape` | `round` |
| Lamp | `bulb_type` | `E27` |

## 12.9 `product_attribute_options`

Optional table for controlled values.

Candidate fields:

```text
attribute_option_id
attribute_id
value
display_name
display_order
status
created_at
updated_at
```

This is useful where values should be governed rather than free text.

Examples:

```text
bed_size: single, double, king, super king
heading_type: eyelet, pencil pleat, tab top
storage_type: none, drawer, ottoman
```

## 12.10 Ownership Rules

```text
catalog-service owns product truth
catalog-service owns category taxonomy
catalog-service owns product attribute definitions
catalog-service owns product attribute values
catalog-service owns product status
catalog-service owns catalogue price in the initial version
inventory-service owns stock levels
search-service owns denormalised search documents and facets
order-service may store product snapshots
```

---

## 13. Search Projection Relationship

Search Service should consume catalogue events and build a denormalised search document.

Example search document:

```json
{
  "product_id": "prd_123",
  "title": "Blackout Eyelet Curtains",
  "category": "curtains",
  "price_minor": 8999,
  "currency_code": "GBP",
  "attributes": {
    "colour": "navy",
    "drop_cm": 228,
    "width_cm": 167,
    "lining": "blackout",
    "heading_type": "eyelet"
  },
  "filterable": {
    "colour": ["navy"],
    "lining": ["blackout"],
    "heading_type": ["eyelet"]
  }
}
```

This allows customer-facing browse and filter behaviour to be optimised without making Search Service the product source of truth.

---

## 14. Inventory Schema

Inventory Service stores stock, warehouse, and reservation data.

Candidate tables:

```text
inventory_items
warehouses
stock_levels
stock_reservations
stock_reservation_items
stock_adjustments
```

Inventory stores `product_id` and `variant_id` as references only.

It does not own product names, descriptions, or attribute definitions.

---

## 15. Basket Schema

Basket Service stores shopping basket data.

Candidate tables:

```text
baskets
basket_items
basket_events_outbox
```

Basket may store display price snapshots for convenience, but final checkout price should be confirmed or snapshotted by Order Service.

---

## 16. Order Schema

Order Service stores order lifecycle data.

Candidate tables:

```text
orders
order_items
order_status_history
checkout_attempts
order_failures
outbox_events
```

Order item snapshots should preserve product details at checkout time.

Examples:

```text
product_name_snapshot
sku_snapshot
unit_price_minor
currency_code
selected_attribute_summary_json
```

The optional `selected_attribute_summary_json` can record customer-relevant selections such as curtain drop or bed size without making Order Service own catalogue attributes.

---

## 17. Payment Schema

Payment Service stores payment state.

Candidate tables:

```text
payments
payment_attempts
refunds
payment_status_history
outbox_events
```

Payment Service must not store raw payment card data.

---

## 18. Shipping Schema

Shipping Service stores shipment and fulfilment data.

Candidate tables:

```text
delivery_options
shipments
shipment_status_history
tracking_events
outbox_events
```

Shipping may store delivery address snapshots for fulfilment history.

---

## 19. Notification Schema

Notification Service stores notification request, attempt, and delivery state.

Candidate tables:

```text
notification_templates
notifications
notification_attempts
processed_events
```

Notification consumers must be idempotent.

---

## 20. Review Schema

Review Service may later store product reviews, ratings, and moderation state.

Candidate tables:

```text
reviews
rating_summaries
moderation_decisions
review_reports
outbox_events
```

Review Service stores `product_id` as a reference only.

---

## 21. Search Schema

Search Service stores search projections and index state.

Candidate tables:

```text
search_index_entries
search_facets
projection_offsets
search_query_logs
```

Search owns denormalised read/search documents only.

Catalogue remains product source of truth.

Search projections must be rebuildable where practical.

---

## 22. Recommendation Schema

Recommendation Service stores recommendation signals and outputs.

Candidate tables:

```text
recommendation_signals
recommendation_rules
recommendation_results
recommendation_feedback
```

Recommendation Service does not own product truth or order truth.

---

## 23. Snapshot Design

Snapshots are allowed when historical accuracy is required.

Important snapshots:

```text
order item product snapshot
order item selected attribute summary
order delivery address snapshot
shipment delivery address snapshot
payment provider reference
```

Snapshot rules:

```text
snapshot fields must be clearly named
snapshots are not source-of-truth live data
snapshots should contain only necessary information
PII snapshots should be minimised and protected
```

---

## 24. Projection Design

Projections are allowed for optimised reads.

Examples:

```text
search-service stores product search documents
recommendation-service stores event-derived signals
notification-service stores processed event IDs
```

Projection rules:

```text
projection owner must be clear
projection lag must be observable
projection should be rebuildable where practical
projection must not become hidden source of truth
```

---

## 25. Outbox Table Design

Services that write data and publish events should consider an outbox table.

Candidate fields:

```text
outbox_event_id
event_id
event_type
event_version
aggregate_type
aggregate_id
payload_json_or_binary
status
attempt_count
next_attempt_at
created_at
published_at
last_error
```

Recommended initial services:

```text
order-service
payment-service
inventory-service
shipping-service
catalog-service
```

Catalogue Service should use outbox if catalogue events become critical to search or recommendations.

---

## 26. Transaction Boundaries

Transactions should stay within one service schema.

Allowed:

```text
catalog-service transaction creates product and attribute values
order-service transaction creates order, order_items, outbox_event
inventory-service transaction reserves stock and records reservation
payment-service transaction records payment and payment_attempt
```

Avoid:

```text
single transaction across order, payment, inventory, shipping, and catalogue schemas
```

---

## 27. Indexing Strategy

Indexes should support service-owned access patterns.

Catalogue examples:

```text
products(status, category_id)
products(category_id, status)
product_variants(product_id)
product_attribute_definitions(category_id, code)
product_attribute_values(product_id)
product_attribute_values(attribute_id)
```

Search/filter-heavy queries should generally be served by Search Service rather than complex Catalogue Service joins.

---

## 28. Audit and Status History

Important lifecycle entities should keep history.

Candidate history tables:

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

## 29. Local Development Seed Data

Local seed data should support:

```text
active products
inactive products
curtains with curtain-specific attributes
bed frames with bed-specific attributes
rugs with rug-specific attributes
stocked products
out-of-stock products
test baskets
successful payment token
declined payment token
simulated delivery options
notification templates
```

Seed data must not contain real personal data.

---

## 30. Security and Privacy

Database design must protect sensitive data.

Sensitive areas:

```text
auth credentials
customer PII
delivery addresses
payment provider references
notification recipient details
```

Catalogue product attributes are usually not sensitive, but supplier-only fields, cost prices, margin data, and unpublished commercial details should not be exposed through customer-facing APIs or search projections.

---

## 31. Initial Checkout Data Flow

Successful checkout data flow:

```text
1. Catalogue Service owns product and category-scoped attribute truth.
2. Basket Service owns basket and basket items.
3. Order Service retrieves basket through API.
4. Inventory Service creates stock reservation.
5. Payment Service records payment authorisation.
6. Shipping Service records shipment.
7. Order Service records order and order item snapshots.
8. Order Service publishes OrderCreated.
9. Notification Service records notification attempt.
```

Each service writes only to its own schema.

---

## 32. Anti-Patterns to Avoid

Avoid:

```text
shared application database for all services
cross-service joins
foreign keys across service schemas
shared ORM models
using database tables as integration contracts
one giant products table with hundreds of nullable type-specific columns
uncontrolled JSON product blobs with no attribute governance
storing raw payment data
projections becoming hidden source of truth
one service running another service's migrations
```

---

## 33. Initial Implementation Checklist

Before implementing a service schema, confirm:

```text
Does this service own the data?
Is the schema separate?
Is the database user service-specific?
Are cross-service references stored as IDs only?
Are snapshots clearly named?
Are migrations owned by the service?
Are indexes based on known access patterns?
Are sensitive fields minimised?
Are idempotency requirements represented?
Are event outbox needs considered?
Are tests planned for repository behaviour?
```

For Catalogue Service, also confirm:

```text
Are variable product characteristics modelled as category-scoped attributes?
Are filterable attributes identified?
Are required attributes defined per category?
Are controlled values represented where needed?
Will Search Service receive enough data to build browse/filter projections?
```

---

## 34. Open Questions

| Question | Status |
|---|---|
| Which ID format will be standard across services? | To decide |
| Should outbox be implemented from the first service release? | Proposed for `order-service`; later for catalogue |
| Should order address snapshots be JSON or structured columns? | To decide |
| Should selected product attributes be snapshotted into order items? | Proposed |
| Should product attributes support variant-level values in version one? | To decide |
| Should catalogue attribute options be a separate table or JSON field first? | To decide |
| Should search use MySQL projection initially or external search later? | To decide |

---

## 35. Related Documents

```text
docs/data/data-ownership.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/events/event-catalog.md
docs/events/event-envelope.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
adr/0003-use-kafka-for-events.md
```

---

## 36. Summary

bfstore’s service database design supports clear microservice ownership by giving each service its own schema, migrations, database user, and data model.

The Catalogue Service remains on MySQL as the governed product source of truth, but supports varied product types through category-scoped attribute definitions and values.

Search Service owns a denormalised product projection optimised for browse, filtering, and faceted search.

This separates catalogue governance from search performance without prematurely moving the product source of truth to NoSQL.
