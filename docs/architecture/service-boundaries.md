# Service Boundaries

## 1. Purpose

This document defines service boundaries for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains what each service owns, what each service does not own, how services communicate, and how the product catalogue is separated from search and recommendation projections.

---

## 2. Boundary Principles

bfstore service boundaries follow these principles:

```text
services are aligned to business capabilities
each service owns its own data
services communicate through APIs and events
services do not read another service's database
shared packages must not contain business ownership logic
projections are not sources of truth
```

---

## 3. Service Landscape

Initial vertical slice services:

```text
api-gateway
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

Later services:

```text
auth-service
customer-service
review-service
search-service
recommendation-service
```

---

## 4. Boundary Summary

| Service | Owns | Does Not Own |
|---|---|---|
| `api-gateway` | Client entry point and routing | Domain business rules |
| `catalog-service` | Product truth, category taxonomy, product attributes | Stock, search ranking, recommendations |
| `basket-service` | Current shopping intent | Stock reservation, order lifecycle |
| `inventory-service` | Stock and reservations | Product details, order lifecycle |
| `order-service` | Order lifecycle and checkout orchestration | Payment internals, stock internals, shipment internals |
| `payment-service` | Payment state and attempts | Order lifecycle, raw card storage |
| `shipping-service` | Shipment state and delivery options | Order lifecycle, customer profile truth |
| `notification-service` | Notification delivery state | Order lifecycle |
| `search-service` | Search projections and facets | Product truth |
| `recommendation-service` | Recommendation signals and outputs | Product/order truth |

---

## 5. Catalogue Service Boundary

## 5.1 Owns

Catalogue Service owns governed product data.

It owns:

```text
product identity
product name and description
product status
category taxonomy
product variants
category-scoped product attribute definitions
product attribute values
product attribute options where implemented
catalogue price in the initial version
catalogue product events
```

This supports varied product types such as:

```text
curtains
bed frames
rugs
lamps
sofas
mattresses
wardrobes
homeware
```

without creating a single table with hundreds of nullable columns.

## 5.2 Does Not Own

Catalogue Service does not own:

```text
stock levels
stock reservations
basket state
order history
payment state
shipment state
notification delivery
search ranking
recommendation outputs
```

## 5.3 Communication

Catalogue Service exposes gRPC APIs such as:

```text
ListProducts
GetProduct
ListCategories
```

It may publish events such as:

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
ProductAttributeDefinitionCreated
ProductAttributeDefinitionUpdated
CategoryUpdated
```

Search Service consumes catalogue events to update denormalised search documents.

---

## 6. Search Service Boundary

## 6.1 Owns

Search Service owns:

```text
denormalised product search documents
facets
filterable fields
search index state
projection offsets
search query behaviour
```

## 6.2 Does Not Own

Search Service does not own:

```text
product truth
category truth
product attribute truth
product status truth
pricing truth
stock truth
```

## 6.3 Boundary Rule

Search Service may store a product document like:

```json
{
  "product_id": "prd_123",
  "title": "Blackout Eyelet Curtains",
  "category": "curtains",
  "attributes": {
    "drop_cm": 228,
    "heading_type": "eyelet"
  },
  "filterable": {
    "heading_type": ["eyelet"]
  }
}
```

But this is a projection.

If Catalogue Service changes the product, Search Service must update or rebuild the projection.

---

## 7. Basket Service Boundary

Basket Service owns:

```text
basket lifecycle
basket items
basket item quantities
customer/session basket association
basket expiry
checked-out marker
```

Basket Service does not own:

```text
stock reservation
order lifecycle
payment state
shipment state
product truth
```

Basket items store `product_id` and `variant_id` as references.

---

## 8. Inventory Service Boundary

Inventory Service owns:

```text
stock levels
stock reservations
reservation expiry
reservation release
stock adjustments
inventory events
```

Inventory Service does not own:

```text
product name
product description
product attributes
basket contents
order lifecycle
payment state
shipment state
```

Inventory stores product references only.

---

## 9. Order Service Boundary

Order Service owns:

```text
checkout attempts
order lifecycle
order status
order item snapshots
order totals
order events
```

For the initial implementation, Order Service orchestrates checkout.

It coordinates:

```text
Basket Service
Inventory Service
Payment Service
Shipping Service
```

but does not own their internal state.

Order item snapshots may include selected product attributes, such as:

```text
bed size
curtain drop
curtain heading type
sofa fabric
```

These are historical snapshots, not live catalogue truth.

---

## 10. Payment Service Boundary

Payment Service owns:

```text
payment state
payment attempts
payment authorisation
payment provider references
refunds where implemented
payment events
```

Payment Service must not store raw card data.

---

## 11. Shipping Service Boundary

Shipping Service owns:

```text
delivery options
shipment creation
shipment status
tracking references
delivery address snapshots
shipping events
```

Shipping Service does not own customer saved address truth.

---

## 12. Notification Service Boundary

Notification Service owns:

```text
notification records
notification attempts
processed event IDs
notification templates
notification delivery events
```

Notification failure must not roll back order creation.

---

## 13. Recommendation Service Boundary

Recommendation Service owns:

```text
recommendation signals
recommendation rules
recommendation outputs
recommendation feedback
```

Recommendation Service does not own:

```text
product truth
order truth
review truth
customer truth
```

It consumes events and stores derived data.

---

## 14. API Gateway Boundary

API Gateway owns:

```text
external routing
request shaping
safe error mapping
correlation ID propagation
authentication integration point
```

API Gateway must not own:

```text
checkout business rules
order lifecycle
payment decisions
stock decisions
product catalogue truth
```

---

## 15. Communication Rules

Use gRPC for commands and queries requiring an immediate response.

Examples:

```text
GetProduct
GetBasket
ReserveStock
AuthorisePayment
CreateShipment
CreateOrder
```

Use Kafka for facts that have already happened.

Examples:

```text
ProductUpdated
StockReserved
PaymentAuthorised
ShipmentCreated
OrderCreated
NotificationSent
```

---

## 16. Data Boundary Rules

```text
each service owns its own schema
no cross-service database joins
no cross-service foreign keys
cross-service references are IDs only
snapshots are clearly named
projections are rebuildable where practical
```

---

## 17. Catalogue vs Search Boundary

This boundary is especially important.

Catalogue Service answers:

```text
What is this product?
Which category does it belong to?
Which attributes are valid for this category?
What are this product's governed attribute values?
Is this product active?
```

Search Service answers:

```text
How should products be found?
Which filters should appear?
What denormalised document supports fast browse?
How should results be ranked or faceted?
```

Catalogue is authoritative.

Search is optimised.

---

## 18. Anti-Patterns to Avoid

Avoid:

```text
Search Service becoming product source of truth
Inventory Service storing product descriptions
Order Service depending on live catalogue data for old orders
Basket Service reserving stock
API Gateway owning checkout logic
shared database tables across services
shared product ORM model across services
one giant product table with many nullable type-specific columns
uncontrolled product JSON with no category governance
```

---

## 19. Boundary Decision Checklist

Before adding behaviour to a service, ask:

```text
Which business capability owns this?
Which service owns the data?
Is this source-of-truth data, a snapshot, or a projection?
Does this require immediate response or an event?
Would this create a cross-service database dependency?
Can this service evolve independently after this change?
```

For product-related changes, also ask:

```text
Is this product truth or search projection?
Is this a category-scoped attribute?
Should this be filterable?
Should this be snapshotted into orders?
Should this be included in ProductUpdated events?
```

---

## 20. Open Questions

| Question | Status |
|---|---|
| Should Search Service be implemented immediately after Catalogue, or deferred? | Proposed: defer until catalogue events are stable |
| Should product attribute options use a table or JSON initially? | To decide |
| Should variant-specific attribute values be supported in version one? | To decide |
| Should Recommendation Service consume catalogue attributes directly or via search projection? | To decide |
| Should API Gateway expose catalogue filters before Search Service exists? | To decide |

---

## 21. Related Documents

```text
docs/architecture/domain-model.md
docs/architecture/communication-patterns.md
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/events/event-catalog.md
proto/acme/catalog/v1/README.md
```

---

## 22. Summary

bfstore service boundaries remain capability-based.

The new catalogue refinement strengthens the architecture:

```text
Catalogue Service owns governed product data and category-scoped attributes.
Search Service owns denormalised search and filter projections.
Recommendation Service owns derived recommendation outputs.
```

This gives bfstore flexibility for varied product types while preserving clean service ownership.
