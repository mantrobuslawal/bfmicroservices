# Data Ownership

## 1. Purpose

This document defines the data ownership model for **bfstore**, ACME Ltd’s fictional online furniture store backend.

The purpose of this document is to make clear:

- which service owns which data
- which service is the source of truth for each domain concept
- how services may access data owned by other services
- where snapshots and projections are acceptable
- how MySQL schemas are separated
- how cross-service data coupling is avoided
- how data ownership supports security, testing, observability, and operations

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s microservice architecture.

---

## 2. Data Ownership Principle

The core data ownership rule is:

> Each service owns its own data. Other services must communicate through APIs or events, not direct database access.

This means:

- each service has its own logical MySQL database or schema
- each service has its own database user
- each service controls its own migrations
- each service defines its own data model
- no service directly reads or writes another service’s tables
- cross-service reads use gRPC APIs, Kafka events, or local projections

The database is an implementation detail of the owning service.

---

## 3. Architecture Context

bfstore is a microservice backend using:

| Technology | Purpose |
|---|---|
| MySQL | Primary relational persistence |
| gRPC | Synchronous access to service-owned behaviour |
| Kafka | Asynchronous propagation of business facts |
| Protobuf | Contract definition for APIs and event payloads |
| Docker Compose | Local development databases |
| Kubernetes | Target runtime platform |

In local development, a single MySQL container may host multiple logical schemas. In production-style environments, these could later be separated into distinct managed databases, clusters, or tenancy boundaries.

The logical ownership model remains the same in both cases.

---

## 4. Service Database Ownership

Each service owns one logical database or schema.

| Service | Database / Schema | Ownership |
|---|---|---|
| `auth-service` | `bfstore_auth` | Identities, credentials, sessions, roles |
| `customer-service` | `bfstore_customer` | Customer profiles, addresses, preferences |
| `catalog-service` | `bfstore_catalog` | Products, categories, variants, pricing |
| `inventory-service` | `bfstore_inventory` | Stock levels, warehouses, reservations |
| `basket-service` | `bfstore_basket` | Baskets and basket items |
| `order-service` | `bfstore_order` | Orders, order items, order state |
| `payment-service` | `bfstore_payment` | Payments, payment attempts, refunds |
| `shipping-service` | `bfstore_shipping` | Shipments, delivery options, tracking |
| `notification-service` | `bfstore_notification` | Notifications, templates, delivery attempts |
| `review-service` | `bfstore_review` | Reviews, ratings, moderation |
| `search-service` | `bfstore_search` | Search projections and query logs |
| `recommendation-service` | `bfstore_recommendation` | Recommendation rules, signals, outputs |

---

## 5. Database Access Rules

## 5.1 Allowed

A service may access:

```text
its own database/schema
its own migrations
its own seed data
its own read models
its own projections
```

Example:

```text
order-service -> bfstore_order
payment-service -> bfstore_payment
catalog-service -> bfstore_catalog
```

## 5.2 Forbidden

A service must not access another service’s database directly.

Forbidden examples:

```text
order-service -> bfstore_inventory.stock_reservations
basket-service -> bfstore_catalog.products
notification-service -> bfstore_customer.addresses
search-service -> bfstore_catalog.products
```

## 5.3 Correct Alternatives

Use gRPC when an immediate answer is required:

```text
basket-service -> catalog-service GetProduct
order-service -> inventory-service ReserveStock
order-service -> payment-service AuthorisePayment
```

Use Kafka events when a service needs to react to a business fact:

```text
catalog-service -> ProductUpdated -> search-service
order-service -> OrderCreated -> notification-service
review-service -> ReviewApproved -> search-service
```

Use projections when a service needs an optimised local read model:

```text
search-service stores product search documents
recommendation-service stores recommendation signals
notification-service stores notification delivery history
```

---

## 6. Source of Truth by Domain Concept

| Domain Concept | Source of Truth | Notes |
|---|---|---|
| Identity | `auth-service` | Authentication identity, credentials, sessions |
| Customer profile | `customer-service` | Customer name, contact data, preferences |
| Address | `customer-service` | Saved customer addresses |
| Product | `catalog-service` | Product details, descriptions, status |
| Product price | `catalog-service` | Base product price and currency |
| Product category | `catalog-service` | Category structure and labels |
| Product variant | `catalog-service` | SKU, size, colour, material variants |
| Stock level | `inventory-service` | Available and reserved stock |
| Stock reservation | `inventory-service` | Reservation lifecycle and expiry |
| Basket | `basket-service` | Current shopping basket |
| Basket item | `basket-service` | Basket product reference and quantity |
| Order | `order-service` | Order lifecycle and state |
| Order item | `order-service` | Historical product and price snapshot |
| Payment | `payment-service` | Payment state and provider references |
| Payment attempt | `payment-service` | Authorisation, capture, refund attempts |
| Shipment | `shipping-service` | Shipment and fulfilment state |
| Notification | `notification-service` | Notification request and delivery state |
| Review | `review-service` | Review content and moderation state |
| Rating summary | `review-service` | Product review aggregates |
| Search document | `search-service` | Search-optimised projection |
| Recommendation | `recommendation-service` | Recommendation output and signals |

---

## 7. Detailed Data Ownership

---

## 7.1 Auth Service Data

### Owns

```text
identities
credentials
sessions
roles
permissions
login_attempts
refresh_tokens
```

### Does Not Own

```text
customer profile
customer addresses
orders
payments
basket contents
```

### Sensitive Data

```text
password hashes
tokens
session identifiers
login attempt metadata
```

### Ownership Rules

- Passwords must never be stored in plain text.
- Tokens must not be logged.
- Auth Service owns identity, not customer profile.
- Customer Service may link to identity by `identity_id`.

---

## 7.2 Customer Service Data

### Owns

```text
customers
customer_profiles
addresses
customer_preferences
contact_preferences
```

### Does Not Own

```text
credentials
sessions
order lifecycle
payment state
shipment state
```

### Sensitive Data

```text
name
email
phone number
delivery address
preferences
```

### Ownership Rules

- Customer Service is the source of truth for customer profile data.
- Orders and shipments may store address snapshots for historical accuracy.
- Customer PII should not be copied unless there is a clear business reason.
- Customer data access should be authorised and auditable.

---

## 7.3 Catalog Service Data

### Owns

```text
products
product_variants
categories
product_images
product_attributes
product_prices
```

### Does Not Own

```text
stock quantity
stock reservations
basket items
orders
reviews
search index
recommendations
```

### Ownership Rules

- Catalog Service owns product facts.
- Product status determines whether products are customer-visible and purchasable.
- Search Service may store product projections but is not the source of truth.
- Order Service stores product snapshots only for historical order accuracy.
- Product changes should be published as events where downstream services need projections.

---

## 7.4 Inventory Service Data

### Owns

```text
inventory_items
warehouses
stock_levels
stock_reservations
stock_adjustments
reservation_expiry
```

### Does Not Own

```text
product descriptions
product pricing
basket contents
order lifecycle
payment state
```

### Ownership Rules

- Inventory Service is the only service allowed to change stock levels.
- Stock cannot be reserved below zero.
- Stock reservations must be explicit and traceable.
- Reservation expiry must be observable.
- Other services may request availability or reservation through APIs.

---

## 7.5 Basket Service Data

### Owns

```text
baskets
basket_items
basket_status
basket_expiry
guest_session_reference
```

### Does Not Own

```text
product source of truth
stock reservation
payment state
order lifecycle
customer profile
```

### Ownership Rules

- Basket Service owns current shopping intent.
- Basket Service does not reserve stock.
- Basket items reference product IDs and quantities.
- Product validity should be checked through Catalog Service.
- Basket may store price snapshots, but final order price must be confirmed at checkout or explicitly snapshotted in the order.

---

## 7.6 Order Service Data

### Owns

```text
orders
order_items
order_status_history
checkout_attempts
order_failures
order_cancellations
```

### Does Not Own

```text
product catalogue source of truth
stock source of truth
payment internals
shipment internals
notification delivery
customer identity
```

### Ownership Rules

- Order Service owns order lifecycle.
- Order items should store historical product and price snapshots.
- Orders should store delivery address snapshots.
- Order Service should not directly modify stock or payment records.
- Order creation should be idempotent where possible.

### Snapshot Data

Order Service may store:

```text
product_name_snapshot
sku_snapshot
unit_price_snapshot
currency
delivery_address_snapshot
```

Reason:

Historical orders must remain accurate even if products, prices, or customer addresses change later.

---

## 7.7 Payment Service Data

### Owns

```text
payments
payment_attempts
refunds
provider_references
payment_status_history
```

### Does Not Own

```text
order lifecycle
stock
shipments
customer profile
raw card data
```

### Sensitive Data

```text
provider references
payment failure reasons
transaction identifiers
```

### Ownership Rules

- Raw card data must not be stored.
- Sensitive payment data must not be logged.
- Payment attempts must be auditable.
- Payment Service owns payment state; Order Service owns order state.
- Payment operations should be idempotent where practical.

---

## 7.8 Shipping Service Data

### Owns

```text
delivery_options
shipments
shipment_status_history
tracking_events
carrier_references
```

### Does Not Own

```text
order creation
payment state
stock source of truth
customer profile source of truth
```

### Ownership Rules

- Shipping Service owns shipment state.
- Shipments should store delivery address snapshots.
- Shipment state should be visible to Order Service through APIs or events.
- Live carrier integration is out of scope initially and may be simulated.

---

## 7.9 Notification Service Data

### Owns

```text
notifications
notification_templates
notification_attempts
delivery_status
retry_state
provider_references
```

### Does Not Own

```text
order state
payment state
shipment state
customer profile source of truth
```

### Ownership Rules

- Notification Service owns delivery status, not the underlying business event.
- Notification failures must not roll back order creation.
- Notification consumers must be idempotent.
- Notification data should avoid unnecessary PII.
- Customer contact details should come from Customer Service or approved event payloads.

---

## 7.10 Review Service Data

### Owns

```text
reviews
ratings
rating_summaries
moderation_decisions
review_reports
```

### Does Not Own

```text
product source of truth
customer profile source of truth
order lifecycle
search projection
```

### Ownership Rules

- Review Service owns review content and moderation state.
- Product existence should be validated against Catalog Service.
- Purchase eligibility may be validated against Order Service.
- Rating summaries may be eventually consistent.

---

## 7.11 Search Service Data

### Owns

```text
search_index_entries
search_facets
search_query_logs
index_update_offsets
projection_state
```

### Does Not Own

```text
product source of truth
stock source of truth
review source of truth
recommendation logic
```

### Ownership Rules

- Search Service owns an optimised projection.
- Search results may be eventually consistent.
- Search Service must support projection rebuild.
- Search Service must not become the product catalogue source of truth.
- Inactive products must not appear in customer-facing search results.

---

## 7.12 Recommendation Service Data

### Owns

```text
recommendation_rules
recommendation_signals
recommendation_results
recommendation_feedback
```

### Does Not Own

```text
product catalogue
orders
payments
reviews
search index
customer profile source of truth
```

### Ownership Rules

- Recommendation Service owns recommendation outputs and signals.
- Recommendations must not include inactive products.
- Recommendation outputs may be eventually consistent.
- Initial recommendations may be rules-based.
- Recommendation Service should degrade gracefully when signal data is limited.

---

## 8. Snapshots

Snapshots are allowed when historical accuracy is required.

## 8.1 Order Product Snapshot

Order items should store product details from the time of purchase.

Example fields:

```text
product_id
variant_id
product_name_snapshot
sku_snapshot
unit_price_snapshot
currency
quantity
line_total
```

Reason:

A customer’s historical order should not change if the product is renamed or repriced later.

---

## 8.2 Delivery Address Snapshot

Orders and shipments should store the delivery address used at checkout.

Reason:

A customer may later update their saved address, but historical orders must preserve the address used at the time.

---

## 8.3 Payment Reference Snapshot

Payment Service may store provider references and payment state.

Reason:

Payment records need auditability without storing raw sensitive payment data.

---

## 8.4 Snapshot Rules

- Snapshots must be clearly named.
- Snapshots must not be treated as live source-of-truth data.
- Snapshots should be limited to fields needed for history, audit, or operations.
- Sensitive snapshot data should be minimised.

---

## 9. Projections

Projections are local read models derived from another service’s source-of-truth data.

Examples:

| Projection | Owner | Source |
|---|---|---|
| Product search document | `search-service` | `catalog-service` events |
| Product availability summary | `search-service` | `inventory-service` events |
| Recommendation signal | `recommendation-service` | basket, order, review events |
| Rating summary | `review-service` | review records |
| Notification history | `notification-service` | domain events and delivery attempts |

## 9.1 Projection Rules

- Projections must have a clear owner.
- Projections must be rebuildable where practical.
- Projections may be eventually consistent.
- Projections must not become hidden sources of truth.
- Projection lag should be observable.

---

## 10. Cross-Service Data Access Patterns

## 10.1 Synchronous Access

Use gRPC when a service needs an immediate answer.

Examples:

```text
basket-service -> catalog-service GetProduct
order-service -> basket-service GetBasket
order-service -> inventory-service ReserveStock
order-service -> payment-service AuthorisePayment
order-service -> shipping-service CreateShipment
```

## 10.2 Asynchronous Access

Use Kafka events when a service needs to react to a business fact.

Examples:

```text
catalog-service -> ProductUpdated -> search-service
order-service -> OrderCreated -> notification-service
review-service -> ReviewApproved -> recommendation-service
```

## 10.3 Local Projection

Use a projection when repeated synchronous calls would make the system too chatty or fragile.

Examples:

```text
search-service maintains product search documents
recommendation-service maintains recommendation signals
```

---

## 11. Consistency Model

## 11.1 Stronger Consistency Required

The checkout path requires careful coordination.

Important concepts:

```text
basket
stock reservation
payment authorisation
order creation
shipment creation
```

The system must avoid:

- confirming orders without stock
- confirming orders without payment authorisation
- creating duplicate confirmed orders
- leaving stock reserved indefinitely after payment failure
- losing visibility into failed shipment creation

This does not require a distributed transaction, but it does require explicit state transitions, idempotency, retries, and compensation.

---

## 11.2 Eventual Consistency Accepted

The following areas may be eventually consistent:

```text
search index
recommendations
review summaries
notification delivery status
availability summaries
analytics projections
```

Eventual consistency must be documented, observable, and testable.

---

## 12. Data Classification

bfstore data should be classified so that storage, logging, and access controls are appropriate.

| Classification | Description | Examples |
|---|---|---|
| Public | Safe for customer-facing display | product name, product description, public price |
| Internal | Operational or system data | service status, non-sensitive IDs, event metadata |
| Confidential | Business-sensitive data | order details, stock levels, provider references |
| Personal | Personally identifiable information | name, email, phone number, address |
| Highly Sensitive | Data requiring strongest controls | password hashes, tokens, payment-related sensitive data |

## 12.1 Logging Rules

Do not log:

```text
passwords
tokens
raw payment details
full customer addresses
secret values
```

Use:

```text
customer_id
order_id
payment_id
shipment_id
correlation_id
trace_id
```

instead of sensitive values where possible.

---

## 13. Database User Model

Each service should have a dedicated database user.

Example:

```text
bfstore_catalog_user
bfstore_inventory_user
bfstore_basket_user
bfstore_order_user
bfstore_payment_user
```

## 13.1 Permission Rules

A service database user should have access only to its own schema.

Example:

```text
catalog-service user -> bfstore_catalog only
order-service user   -> bfstore_order only
payment-service user -> bfstore_payment only
```

No service user should have broad access to all schemas.

---

## 14. Migrations

Each service owns its own migrations.

Example layout:

```text
db/
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
└── shipping/
    └── migrations/
```

## 14.1 Migration Rules

- Migrations should be versioned.
- Migrations should be reviewed.
- Destructive migrations should be handled carefully.
- Rollback strategy should be documented.
- Service deployments and schema changes should be coordinated safely.
- Migrations should be tested in CI where practical.

---

## 15. Backup and Restore Ownership

Each service’s data should have a clear restore expectation.

| Service | Restore Concern |
|---|---|
| `auth-service` | Identity and session recovery |
| `customer-service` | Customer profile and address recovery |
| `catalog-service` | Product catalogue recovery |
| `inventory-service` | Stock and reservation consistency |
| `basket-service` | Active basket recovery |
| `order-service` | Order history and checkout state |
| `payment-service` | Payment audit and reconciliation |
| `shipping-service` | Shipment state and tracking |
| `notification-service` | Notification delivery history |
| `review-service` | Review and moderation history |
| `search-service` | Projection rebuild may be preferred over restore |
| `recommendation-service` | Projection rebuild may be preferred over restore |

Search and recommendation data may be rebuildable from source events, while order and payment data require stronger restore guarantees.

---

## 16. Data Retention

Retention requirements should be defined per domain.

Initial guidance:

| Data Type | Retention Approach |
|---|---|
| Product catalogue | Retain while active; archive discontinued products |
| Basket data | Expire abandoned baskets after defined period |
| Orders | Retain for business and audit needs |
| Payments | Retain payment audit records; minimise sensitive fields |
| Notifications | Retain delivery status for operational period |
| Search logs | Retain for limited analytics period |
| Recommendation signals | Retain only while useful |
| Customer data | Support deletion/anonymisation where required |
| Reviews | Retain while visible; support moderation and deletion |

Detailed retention rules should be documented in:

```text
docs/data/retention.md
```

---

## 17. Privacy and PII Handling

Personal data appears in:

```text
customer profile
customer address
orders
payments
shipments
notifications
reviews
```

## 17.1 PII Rules

- Minimise personal data copied across services.
- Prefer IDs over personal details in events.
- Use snapshots only when required for historical or operational accuracy.
- Avoid logging PII.
- Document where PII is stored.
- Consider anonymisation or deletion workflows for customer data.

Detailed PII handling should be documented in:

```text
docs/data/pii-handling.md
```

---

## 18. Data Ownership During Checkout

The checkout flow touches multiple data owners.

| Step | Data Owner | Data Created or Updated |
|---|---|---|
| Basket retrieved | `basket-service` | Basket, basket items |
| Product validated | `catalog-service` | Product status checked |
| Stock reserved | `inventory-service` | Stock reservation |
| Payment authorised | `payment-service` | Payment, payment attempt |
| Order created | `order-service` | Order, order items |
| Shipment created | `shipping-service` | Shipment |
| Notification sent | `notification-service` | Notification attempt |

## 18.1 Checkout Rule

Order Service may coordinate checkout, but it does not take ownership of stock, payment, shipment, or notification data.

---

## 19. Outbox Pattern

For services that update a database and publish Kafka events, the outbox pattern should be considered.

## 19.1 Problem

A service may successfully write to MySQL but fail to publish the matching Kafka event.

Example:

```text
Order created in database
Kafka publish fails
Notification Service never receives OrderCreated
```

## 19.2 Proposed Solution

Use an outbox table owned by the service:

```text
orders
order_items
outbox_events
```

The service writes business data and outbox event in the same local transaction. A publisher process then publishes the outbox event to Kafka.

## 19.3 Initial Recommendation

For a serious implementation, use the outbox pattern for events from:

```text
order-service
payment-service
inventory-service
shipping-service
notification-service
```

This should be captured in an ADR if implemented.

---

## 20. Anti-Patterns to Avoid

## 20.1 Shared Database

Avoid:

```text
all services -> bfstore_db
```

Prefer:

```text
catalog-service -> bfstore_catalog
order-service   -> bfstore_order
payment-service -> bfstore_payment
```

---

## 20.2 Cross-Service Joins

Avoid:

```sql
SELECT *
FROM bfstore_order.orders o
JOIN bfstore_customer.customers c ON c.customer_id = o.customer_id;
```

Prefer service APIs, events, or reporting projections.

---

## 20.3 Shared ORM Models

Avoid sharing database models across services.

Each service should own its own persistence model.

---

## 20.4 Hidden Source-of-Truth Projections

Avoid allowing projections to become hidden sources of truth.

Example risk:

```text
search-service product document becomes more trusted than catalog-service product record
```

---

## 21. Data Ownership Checklist

Before adding data to a service, ask:

```text
Which service owns this data?
Is this source-of-truth data or a projection?
Does another service already own this concept?
Who is allowed to write this data?
Who is allowed to read this data directly?
Does this data contain PII?
Does it need retention rules?
Does it need audit logging?
Does it need to be exposed through an API or event?
Can it be rebuilt from events?
```

---

## 22. Initial Implementation Scope

The first version should focus on data ownership for the checkout vertical slice.

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

Initial ownership focus:

```text
products
stock reservations
baskets
orders
payments
shipments
notifications
```

Deferred data areas:

```text
auth identity
customer profile management
reviews
search projections
recommendations
advanced analytics
```

---

## 23. Open Questions

| Question | Status |
|---|---|
| Should local development use one MySQL container with separate schemas? | Proposed |
| Should production-style environments separate databases physically or logically? | To decide |
| Should outbox pattern be implemented from the first service release? | Proposed for serious version |
| Should Basket Service store price snapshots before checkout? | To decide |
| Should Order Service store full delivery address snapshots from the first version? | Proposed |
| Should Search Service initially use a MySQL projection or a dedicated search engine later? | To decide |
| Should recommendation data be rebuilt entirely from events? | To decide |
| What customer deletion/anonymisation workflow is required? | To decide |

---

## 24. Related Documents

This document should be read alongside:

```text
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/events/event-catalog.md
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/pii-handling.md
docs/testing/testing-strategy.md
docs/security/privacy-and-pii.md
```

Relevant ADRs:

```text
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
adr/0003-use-kafka-for-events.md
adr/0008-use-contract-first-service-design.md
```

---

## 25. Summary

bfstore’s data ownership model is designed to avoid shared database coupling and preserve clear microservice boundaries.

The most important rules are:

```text
Each service owns its own data.
No service directly accesses another service’s database.
APIs and events are the integration boundaries.
Snapshots are allowed for history.
Projections are allowed for query optimisation.
Source-of-truth ownership must always be clear.
```

This model supports independent service development, stronger security boundaries, clearer testing, safer deployments, and a more professional cloud-native architecture.
