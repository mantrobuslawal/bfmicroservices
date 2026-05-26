# Service Database Design

## 1. Purpose

This document defines the service database design approach for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how each service owns its data, how MySQL schemas are separated, how tables should be designed, how cross-service data references should work, and how the initial database model supports the checkout vertical slice.

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

## 4. Logical Database Model

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

## 5. Schema Ownership Summary

| Service | Schema | Purpose |
|---|---|---|
| `auth-service` | `bfstore_auth` | Identity, credentials, sessions, roles |
| `customer-service` | `bfstore_customer` | Customer profiles, addresses, preferences |
| `catalog-service` | `bfstore_catalog` | Products, variants, categories, prices |
| `inventory-service` | `bfstore_inventory` | Stock, warehouses, reservations |
| `basket-service` | `bfstore_basket` | Baskets and basket items |
| `order-service` | `bfstore_order` | Orders, order items, order state |
| `payment-service` | `bfstore_payment` | Payments, attempts, refunds |
| `shipping-service` | `bfstore_shipping` | Shipments, delivery options, tracking |
| `notification-service` | `bfstore_notification` | Notifications, attempts, templates |
| `review-service` | `bfstore_review` | Reviews, ratings, moderation |
| `search-service` | `bfstore_search` | Search projections and index state |
| `recommendation-service` | `bfstore_recommendation` | Recommendation signals and results |

---

## 6. Database Access Rules

## 6.1 Allowed

A service may:

```text
read its own schema
write its own schema
run migrations for its own schema
manage seed data for its own schema
create indexes in its own schema
own transaction boundaries inside its own schema
```

## 6.2 Forbidden

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

## 7. Cross-Service References

Services may store references to entities owned by other services.

Example:

```text
order-service stores customer_id
order-service stores product_id
payment-service stores order_id
shipping-service stores order_id
basket-service stores product_id
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

## 8. Foreign Key Strategy

## 8.1 Within a Service Schema

Foreign keys are allowed within a service-owned schema.

Example:

```text
bfstore_order.orders
bfstore_order.order_items
```

`order_items.order_id` may reference `orders.order_id`.

## 8.2 Across Service Schemas

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

## 9. Identifier Strategy

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

## 10. Initial Implementation Schemas

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

This allows the project to prove the core business flow before expanding into secondary services.

---

## 11. Catalogue Schema

## 11.1 Purpose

The Catalogue schema stores product information owned by `catalog-service`.

## 11.2 Candidate Tables

```text
products
product_variants
categories
product_images
product_attributes
product_price_history
```

## 11.3 `products`

Candidate fields:

```text
product_id
category_id
name
description
status
base_price_minor
currency_code
material
colour
created_at
updated_at
```

## 11.4 `product_variants`

Candidate fields:

```text
variant_id
product_id
sku
size
colour
material
price_minor
currency_code
status
created_at
updated_at
```

## 11.5 Ownership Rules

```text
catalog-service owns product truth
catalog-service owns product status
catalog-service owns product pricing in the initial version
inventory-service owns stock levels
search-service may later store product projections
order-service may store product snapshots
```

---

## 12. Inventory Schema

## 12.1 Purpose

The Inventory schema stores stock, warehouse, and reservation data owned by `inventory-service`.

## 12.2 Candidate Tables

```text
inventory_items
warehouses
stock_levels
stock_reservations
stock_reservation_items
stock_adjustments
```

## 12.3 `stock_levels`

Candidate fields:

```text
stock_level_id
product_id
variant_id
warehouse_id
available_quantity
reserved_quantity
created_at
updated_at
```

## 12.4 `stock_reservations`

Candidate fields:

```text
reservation_id
order_id
basket_id
customer_id
status
idempotency_key
expires_at
created_at
updated_at
```

## 12.5 `stock_reservation_items`

Candidate fields:

```text
reservation_item_id
reservation_id
product_id
variant_id
quantity
created_at
```

## 12.6 Ownership Rules

```text
inventory-service is the only service that changes stock levels
stock reservations must be idempotent
stock release must be idempotent
stock cannot be reserved below zero
reservation expiry must be observable
```

---

## 13. Basket Schema

## 13.1 Purpose

The Basket schema stores shopping basket data owned by `basket-service`.

## 13.2 Candidate Tables

```text
baskets
basket_items
basket_events_outbox
```

## 13.3 `baskets`

Candidate fields:

```text
basket_id
customer_id
session_id
status
created_at
updated_at
expires_at
checked_out_at
```

## 13.4 `basket_items`

Candidate fields:

```text
basket_item_id
basket_id
product_id
variant_id
quantity
unit_price_snapshot_minor
currency_code
created_at
updated_at
```

## 13.5 Ownership Rules

```text
basket-service owns current shopping intent
basket-service does not reserve stock
basket-service may store price snapshots for display
final checkout price should be confirmed or snapshotted by order-service
```

---

## 14. Order Schema

## 14.1 Purpose

The Order schema stores order lifecycle data owned by `order-service`.

## 14.2 Candidate Tables

```text
orders
order_items
order_status_history
checkout_attempts
order_failures
outbox_events
```

## 14.3 `orders`

Candidate fields:

```text
order_id
order_number
customer_id
basket_id
status
total_amount_minor
currency_code
delivery_address_snapshot_json
idempotency_key
created_at
updated_at
confirmed_at
cancelled_at
```

## 14.4 `order_items`

Candidate fields:

```text
order_item_id
order_id
product_id
variant_id
product_name_snapshot
sku_snapshot
unit_price_minor
currency_code
quantity
line_total_minor
created_at
```

## 14.5 `checkout_attempts`

Candidate fields:

```text
checkout_attempt_id
customer_id
basket_id
idempotency_key
status
failure_reason
created_at
updated_at
completed_at
```

## 14.6 Ownership Rules

```text
order-service owns order lifecycle
order-service stores historical product and address snapshots
order-service does not own stock, payment, or shipment state
order-service may coordinate checkout but must respect service ownership
order creation must be idempotent
```

---

## 15. Payment Schema

## 15.1 Purpose

The Payment schema stores payment data owned by `payment-service`.

## 15.2 Candidate Tables

```text
payments
payment_attempts
refunds
payment_status_history
outbox_events
```

## 15.3 `payments`

Candidate fields:

```text
payment_id
order_id
customer_id
status
amount_minor
currency_code
provider
provider_reference
idempotency_key
created_at
updated_at
authorised_at
captured_at
```

## 15.4 `payment_attempts`

Candidate fields:

```text
payment_attempt_id
payment_id
order_id
attempt_type
status
failure_reason
provider_reference
created_at
completed_at
```

## 15.5 Ownership Rules

```text
payment-service owns payment state
payment-service must not store raw card data
payment operations must be idempotent
payment attempts must be auditable
payment provider references should be protected
```

---

## 16. Shipping Schema

## 16.1 Purpose

The Shipping schema stores shipment and fulfilment data owned by `shipping-service`.

## 16.2 Candidate Tables

```text
delivery_options
shipments
shipment_status_history
tracking_events
outbox_events
```

## 16.3 `shipments`

Candidate fields:

```text
shipment_id
order_id
customer_id
status
delivery_option_id
tracking_reference
carrier
delivery_address_snapshot_json
idempotency_key
created_at
updated_at
dispatched_at
delivered_at
```

## 16.4 Ownership Rules

```text
shipping-service owns shipment state
shipments should be created idempotently
delivery address should be snapshotted
carrier integrations may be simulated initially
shipment failure behaviour must be documented
```

---

## 17. Notification Schema

## 17.1 Purpose

The Notification schema stores notification request, attempt, and delivery state owned by `notification-service`.

## 17.2 Candidate Tables

```text
notification_templates
notifications
notification_attempts
processed_events
```

## 17.3 `notifications`

Candidate fields:

```text
notification_id
event_id
notification_type
recipient_reference
channel
status
template_id
created_at
updated_at
sent_at
```

## 17.4 `notification_attempts`

Candidate fields:

```text
notification_attempt_id
notification_id
provider
status
failure_reason
attempt_number
created_at
completed_at
```

## 17.5 `processed_events`

Candidate fields:

```text
event_id
event_type
producer
processed_at
processing_status
```

## 17.6 Ownership Rules

```text
notification-service owns delivery state
notification failure must not roll back order creation
event processing must be idempotent
notification data should minimise PII
```

---

## 18. Review Schema

## 18.1 Purpose

The Review schema stores product reviews, ratings, and moderation data owned by `review-service`.

## 18.2 Candidate Tables

```text
reviews
rating_summaries
moderation_decisions
review_reports
outbox_events
```

## 18.3 Ownership Rules

```text
review-service owns review content
review-service owns moderation state
review-service may validate purchase eligibility through order-service
rating summaries may be eventually consistent
```

This schema may be deferred until the checkout vertical slice is complete.

---

## 19. Search Schema

## 19.1 Purpose

The Search schema stores search projections and index state owned by `search-service`.

## 19.2 Candidate Tables

```text
search_index_entries
search_facets
projection_offsets
search_query_logs
```

## 19.3 Ownership Rules

```text
search-service owns search projection only
catalog-service remains product source of truth
search projection must be rebuildable
search index lag should be observable
```

This schema may be deferred until product browsing and catalogue events are stable.

---

## 20. Recommendation Schema

## 20.1 Purpose

The Recommendation schema stores recommendation signals and outputs owned by `recommendation-service`.

## 20.2 Candidate Tables

```text
recommendation_signals
recommendation_rules
recommendation_results
recommendation_feedback
```

## 20.3 Ownership Rules

```text
recommendation-service owns recommendation outputs
recommendation-service does not own product or order truth
recommendations must exclude inactive products
recommendation failure must not block browsing or checkout
```

This schema may be deferred until core commerce events are available.

---

## 21. Snapshot Design

Snapshots are allowed when historical accuracy is required.

Important snapshots:

```text
order item product snapshot
order delivery address snapshot
shipment delivery address snapshot
payment provider reference
```

## 21.1 Snapshot Rules

```text
snapshot fields must be clearly named
snapshots are not source-of-truth live data
snapshots should contain only necessary information
PII snapshots should be minimised and protected
```

Example order item snapshot fields:

```text
product_name_snapshot
sku_snapshot
unit_price_minor
currency_code
```

---

## 22. Projection Design

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

## 23. Outbox Table Design

Services that write data and publish events should consider an outbox table.

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
review-service
```

The first serious implementation should consider outbox for `order-service` because `OrderCreated` is critical.

---

## 24. Transaction Boundaries

Transactions should stay within one service schema.

Allowed:

```text
order-service transaction creates order, order_items, outbox_event
inventory-service transaction reserves stock and records reservation
payment-service transaction records payment and payment_attempt
```

Avoid:

```text
single transaction across order, payment, inventory, and shipping schemas
```

Use:

```text
idempotency
state machines
events
compensating actions
outbox
reconciliation
```

instead of distributed database transactions.

---

## 25. Indexing Strategy

Indexes should support service-owned access patterns.

Examples:

```text
products(status, category_id)
basket_items(basket_id)
orders(customer_id, created_at)
orders(idempotency_key)
stock_levels(product_id, variant_id, warehouse_id)
stock_reservations(idempotency_key)
payments(order_id)
payments(idempotency_key)
shipments(order_id)
notifications(event_id)
processed_events(event_id)
```

Indexing should be reviewed after real query patterns emerge.

Avoid adding indexes speculatively without a known access pattern.

---

## 26. Audit and Status History

Important lifecycle entities should keep history.

Candidate history tables:

```text
order_status_history
payment_status_history
shipment_status_history
stock_adjustments
notification_attempts
```

History tables support:

```text
debugging
auditability
incident response
customer support
reconciliation
```

---

## 27. Local Development Seed Data

Local seed data should support:

```text
active products
inactive products
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

## 28. Security and Privacy

Database design must protect sensitive data.

Sensitive areas:

```text
auth credentials
customer PII
delivery addresses
payment provider references
notification recipient details
```

Rules:

```text
do not store raw payment card data
do not store secrets in tables unnecessarily
avoid logging sensitive fields
minimise PII duplication
protect database users with least privilege
```

---

## 29. Initial Checkout Data Flow

Successful checkout data flow:

```text
1. Basket Service owns basket and basket items.
2. Order Service retrieves basket through API.
3. Inventory Service creates stock reservation.
4. Payment Service records payment authorisation.
5. Shipping Service records shipment.
6. Order Service records order and order items.
7. Order Service publishes OrderCreated.
8. Notification Service records notification attempt.
```

Each service writes only to its own schema.

---

## 30. Anti-Patterns to Avoid

Avoid:

```text
shared application database for all services
cross-service joins
foreign keys across service schemas
shared ORM models
using database tables as integration contracts
storing raw payment data
unbounded JSON blobs for core relational data
projections becoming hidden source of truth
one service running another service's migrations
```

---

## 31. Initial Implementation Checklist

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

---

## 32. Open Questions

| Question | Status |
|---|---|
| Which ID format will be standard across services? | To decide |
| Should outbox be implemented from the first service release? | Proposed for `order-service` |
| Should order address snapshots be JSON or structured columns? | To decide |
| Should price history be stored from version one? | To decide |
| Should local development create all schemas upfront or only active service schemas? | To decide |
| Should search use MySQL projection initially or external search later? | To decide |
| How long should idempotency records be retained? | To decide |

---

## 33. Related Documents

This document should be read alongside:

```text
docs/data/data-ownership.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/events/event-catalog.md
docs/events/event-envelope.md
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/architecture/resilience-patterns.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
adr/0003-use-kafka-for-events.md
```

---

## 34. Summary

bfstore’s service database design supports clear microservice ownership by giving each service its own schema, migrations, database user, and data model.

The most important rules are:

```text
each service owns its own schema
no cross-service database access
no cross-service joins
cross-service references are IDs only
snapshots are allowed for history
projections are allowed for optimised reads
transactions stay inside one service boundary
outbox should be considered for reliable event publishing
```

This design provides a professional foundation for reliable, secure, and maintainable service-owned persistence.
