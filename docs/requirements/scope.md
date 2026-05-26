# Scope

## 1. Purpose

This document defines the scope of **bfstore**, the backend platform for **Borough Furniture Store**.

It clarifies what is included in the initial implementation, what is deferred, and what is deliberately excluded.

---

## 2. Product Context

Borough Furniture Store is a fictional ecommerce company selling developer-themed furniture and homeware.

The company began with Golang-themed homeware inspired by the Go gopher mascot and later expanded into programming-language mascot-themed and computer-science-inspired products.

Example products:

```text
Gopher desk lamp
Gopher cushion set
Rob Pike wall tapestry
Rivest super-secure lockbox
Dijkstra pathfinding rug
Grace Hopper debugging blanket
Turing machine wall clock
Rust crab coat hooks
Python plush beanbag
Kubernetes helm bookends
```

The platform should be realistic enough to support a broader furniture and homeware catalogue over time.

---

## 3. Scope Strategy

bfstore is intentionally built in stages.

The first stage should prove the core architecture through one complete customer journey rather than trying to implement every ecommerce feature.

The first implementation should focus on:

```text
small but complete
well documented
testable
observable
architecturally representative
```

---

## 4. In Scope for Initial Version

## 4.1 Product Catalogue

Included:

```text
basic product listing
product detail retrieval
categories
product variants
category-scoped product attributes
active/inactive product status
fictional seed products
```

The catalogue should support different product types without a rigid product table full of nullable fields.

Initial product examples should include multiple categories, such as:

```text
wall tapestries
lockboxes
desk lamps
rugs
cushions
bookends
```

## 4.2 Basket

Included:

```text
create basket
add item to basket
update item quantity
remove item
get basket
mark basket as checked out
```

## 4.3 Inventory

Included:

```text
basic stock level storage
check availability
reserve stock
release reservation on failure
prevent obvious overselling
```

## 4.4 Order and Checkout

Included:

```text
create order from basket
coordinate checkout flow
store order item snapshots
store selected product attribute summaries where useful
publish OrderCreated event
return order result
```

The first checkout orchestration is owned by Order Service.

## 4.5 Payment

Included:

```text
simulated payment authorisation
successful payment path
declined payment path
payment attempt recording
idempotency support
```

Real payment provider integration is out of scope for the first version.

## 4.6 Shipping

Included:

```text
static delivery options
simulated shipment creation
shipment status
shipment record
```

Real carrier integration is out of scope for the first version.

## 4.7 Notification

Included:

```text
consume OrderCreated
create notification record
simulate order confirmation
record notification attempt
avoid duplicate notifications for duplicate events
```

## 4.8 Events

Included:

```text
OrderCreated
PaymentAuthorised
PaymentFailed
StockReserved
StockReservationFailed
ShipmentCreated
ShipmentFailed
NotificationSent
NotificationFailed
```

Catalogue events may be documented and partially implemented later.

## 4.9 Database

Included:

```text
service-owned MySQL schemas
migration directories
local database initialisation
catalogue flexible attribute tables
order snapshots
payment attempts
stock reservations
notification processed events
```

## 4.10 Documentation

Included:

```text
requirements
architecture
ADRs
API design
event design
data ownership
testing strategy
database design
service README files
```

---

## 5. Deferred Scope

## 5.1 Authentication and Customer Accounts

Deferred:

```text
full login and registration
password reset
customer profile management
saved customer addresses
role-based admin access
```

These can be added once the checkout vertical slice works.

## 5.2 Search Service

Deferred or light initial version:

```text
full-text search
faceted filtering
ranking
search index rebuilds
query analytics
```

Search Service should eventually consume catalogue events and build denormalised product documents.

## 5.3 Recommendation Service

Deferred:

```text
related products
personalised recommendations
event-derived recommendation signals
recommendation feedback
```

## 5.4 Review Service

Deferred:

```text
customer reviews
ratings
moderation
review reporting
rating summaries
```

## 5.5 Advanced Catalogue Management

Deferred:

```text
admin product creation UI
bulk imports
supplier feeds
product approval workflow
complex pricing
promotions
discounts
tax rules
```

## 5.6 Real External Providers

Deferred:

```text
real payment provider
real shipping provider
email/SMS provider
fraud provider
tax provider
```

Simulated providers are sufficient for the first implementation.

## 5.7 Advanced Platform Capabilities

Deferred:

```text
multi-region deployment
service mesh
canary releases
advanced chaos engineering
full FinOps automation
advanced policy-as-code
production-grade secrets rotation
```

---

## 6. Explicitly Out of Scope

The project will not initially include:

```text
real customer payment data
real customer addresses
production secrets
real supplier integrations
marketplace seller onboarding
returns and exchanges
warehouse management system
ERP integration
mobile app
frontend ecommerce site
```

A lightweight frontend may be added later, but the first priority is backend and platform evidence.

---

## 7. First Vertical Slice

The first end-to-end slice is:

```text
1. Customer lists products.
2. Customer views a product.
3. Customer adds product to basket.
4. Customer checks out.
5. Order Service retrieves basket.
6. Inventory Service reserves stock.
7. Payment Service authorises payment.
8. Shipping Service creates shipment.
9. Order Service creates order.
10. Order Service publishes OrderCreated.
11. Notification Service consumes OrderCreated.
12. Notification Service records simulated confirmation.
```

This slice proves the major architectural decisions.

---

## 8. Product Data Scope

The initial product catalogue should include varied product types to prove the flexible catalogue model.

Required initial product categories:

```text
wall art
desk accessories
soft furnishings
storage
lighting
rugs
```

Required attribute examples:

```text
dimensions
material
colour
theme
mascot
programming_language
security_rating
bulb_type
fabric_type
mounting_style
```

This ensures the catalogue model is tested against realistic variation.

---

## 9. Technical Scope

Initial technical artefacts should include:

```text
buf.yaml
buf.gen.yaml
gRPC proto contracts
service skeletons
database migrations
Docker Compose
Makefile
Kafka topics
basic structured logging
basic OpenTelemetry wiring
health checks
contract tests
integration tests
```

---

## 10. Scope Boundaries by Service

| Service | Initial Scope |
|---|---|
| `catalog-service` | List/get products, categories, attributes |
| `basket-service` | Manage basket |
| `inventory-service` | Check/reserve/release stock |
| `order-service` | Checkout orchestration and order creation |
| `payment-service` | Simulated payment authorisation |
| `shipping-service` | Simulated shipment creation |
| `notification-service` | Consume order events and simulate confirmation |
| `api-gateway` | Optional initial routing layer |
| `search-service` | Deferred |
| `recommendation-service` | Deferred |
| `review-service` | Deferred |
| `auth-service` | Deferred |
| `customer-service` | Deferred |

---

## 11. Non-Functional Scope

Initial non-functional scope includes:

```text
clear local development workflow
repeatable migrations
structured logs
correlation IDs
basic tracing
health endpoints
idempotency for checkout-critical operations
test coverage for key flows
safe error handling
no raw payment data
```

---

## 12. Success Criteria

The initial scope is complete when:

```text
a developer can run the system locally
seed data creates Borough Furniture Store products
Catalogue Service returns products with varied attributes
Basket Service manages basket items
Order Service performs checkout orchestration
Inventory Service reserves stock
Payment Service simulates authorisation
Shipping Service simulates shipment creation
OrderCreated is published
Notification Service consumes OrderCreated
tests prove the main flow
documentation matches implementation
```

---

## 13. Risks

| Risk | Mitigation |
|---|---|
| Scope becomes too large | Focus on first vertical slice |
| Product model becomes overcomplicated | Start with a small set of attribute types |
| Catalogue attributes become ungoverned | Use category-scoped definitions |
| Search is implemented too early | Defer until catalogue events are stable |
| Too many services before working flow | Skeleton only what the vertical slice needs |
| Documentation drifts from code | Update docs after implementation lessons |

---

## 14. Related Documents

```text
docs/requirements/product-vision.md
docs/requirements/functional-requirements.md
docs/requirements/non-functional-requirements.md
docs/requirements/business-rules.md
docs/requirements/acceptance-criteria.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/data/service-database-design.md
```

---

## 15. Summary

The initial bfstore scope is intentionally focused.

The project should first prove a complete checkout journey for Borough Furniture Store using a small but varied developer-themed product catalogue.

The scope includes flexible catalogue attributes, service-owned MySQL schemas, gRPC APIs, Kafka events, and enough operational/testing evidence to demonstrate senior engineering judgement.
