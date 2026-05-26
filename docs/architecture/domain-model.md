# Domain Model

## 1. Purpose

This document defines the high-level domain model for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It describes the core business concepts, their relationships, and which service owns each concept.

This document is intended for engineers, reviewers, technical leads, and clients evaluating bfstore’s business modelling and service boundary design.

---

## 2. Domain Overview

bfstore supports an online furniture and homeware shopping journey.

Core domains:

```text
Catalogue
Basket
Inventory
Order
Payment
Shipping
Notification
Customer
Review
Search
Recommendation
```

The first implementation focuses on the checkout vertical slice:

```text
Browse product
→ Add to basket
→ Checkout
→ Reserve stock
→ Authorise payment
→ Create shipment
→ Create order
→ Publish OrderCreated
→ Send notification
```

---

## 3. Key Modelling Principle

The domain model separates:

```text
source-of-truth models
snapshots
references
projections
```

Examples:

```text
Catalogue Service owns Product truth.
Order Service stores product snapshots for historical order accuracy.
Search Service stores denormalised product projections for browse and filtering.
Inventory Service stores product IDs as references while owning stock.
```

---

## 4. Catalogue Domain

## 4.1 Purpose

The Catalogue domain owns governed product information.

It must support varied future product types such as:

```text
curtains
bed frames
mattresses
sofas
rugs
lamps
tables
wardrobes
mirrors
cushions
homeware
```

## 4.2 Core Concepts

```text
Product
Category
ProductVariant
ProductAttributeDefinition
ProductAttributeValue
ProductAttributeOption
ProductImage
ProductPriceHistory
```

## 4.3 Concept Definitions

### Product

A sellable catalogue item.

Common fields:

```text
product_id
category_id
name
description
status
base_price
brand
created_at
updated_at
```

A Product contains only common fields shared across product types.

Product-type-specific details belong in attributes.

### Category

A taxonomy node used to organise products and define relevant attributes.

Examples:

```text
curtains
bed-frames
rugs
lamps
sofas
```

Categories may form a hierarchy.

### ProductVariant

A purchasable variation of a product.

Examples:

```text
same curtain in different widths and drops
same bed frame in different sizes
same sofa in different fabrics
```

### ProductAttributeDefinition

A category-scoped definition of a product attribute.

Examples:

| Category | Attribute |
|---|---|
| curtains | drop_cm |
| curtains | heading_type |
| bed-frames | bed_size |
| bed-frames | storage_type |
| rugs | shape |
| lamps | bulb_type |

Definitions describe:

```text
code
display name
data type
unit
required flag
filterable flag
allowed values
```

### ProductAttributeValue

A product-specific or variant-specific value for a defined attribute.

Examples:

```text
drop_cm = 228
heading_type = eyelet
bed_size = king
storage_type = ottoman
rug_shape = round
bulb_type = E27
```

### ProductAttributeOption

A controlled allowed value for an attribute.

Examples:

```text
bed_size: single, double, king, super king
heading_type: eyelet, pencil pleat, tab top
storage_type: none, drawer, ottoman
```

## 4.4 Catalogue Ownership

Catalogue Service owns:

```text
product identity
product name and description
category taxonomy
variant metadata
attribute definitions
attribute values
product status
initial product price
```

Catalogue Service does not own:

```text
stock quantity
stock reservation
basket state
order history
payment state
shipment state
search ranking
recommendation outputs
```

---

## 5. Search Projection Domain

## 5.1 Purpose

Search Service owns denormalised product search documents and facets.

It is optimised for:

```text
keyword search
browse pages
filtering
faceted navigation
sorting
projection rebuilds
```

## 5.2 Relationship to Catalogue

Catalogue is the source of truth.

Search is a projection.

Flow:

```text
Catalogue Service updates product
Catalogue Service publishes ProductUpdated
Search Service consumes event
Search Service updates denormalised product search document
```

## 5.3 Example Search Document

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

Search Service must not become the hidden product source of truth.

---

## 6. Basket Domain

Basket Service owns current shopping intent.

Core concepts:

```text
Basket
BasketItem
BasketStatus
```

Basket items reference products and variants by ID.

Basket Service does not reserve stock and does not own final order state.

---

## 7. Inventory Domain

Inventory Service owns stock and reservations.

Core concepts:

```text
StockLevel
StockReservation
StockReservationItem
ReservationStatus
StockAdjustment
```

Inventory stores `product_id` and `variant_id` as references only.

It does not own product details or catalogue attributes.

---

## 8. Order Domain

Order Service owns order lifecycle and checkout orchestration in the initial implementation.

Core concepts:

```text
CheckoutAttempt
Order
OrderItem
OrderStatus
OrderStatusHistory
```

Order item snapshots may include:

```text
product_name_snapshot
sku_snapshot
unit_price
selected_attribute_summary
quantity
line_total
```

`selected_attribute_summary` records customer-relevant product selections at checkout time, such as bed size or curtain drop, without making Order Service the owner of catalogue attributes.

---

## 9. Payment Domain

Payment Service owns payment state.

Core concepts:

```text
Payment
PaymentAttempt
PaymentStatus
Refund
```

Payment Service must not store raw payment card data.

---

## 10. Shipping Domain

Shipping Service owns shipment state and delivery options.

Core concepts:

```text
DeliveryOption
Shipment
ShipmentStatus
TrackingEvent
```

Shipping may store delivery address snapshots for fulfilment history.

---

## 11. Notification Domain

Notification Service owns notification delivery state.

Core concepts:

```text
Notification
NotificationAttempt
NotificationStatus
NotificationTemplate
ProcessedEvent
```

Notification Service must process events idempotently.

---

## 12. Customer Domain

Customer Service owns customer profile and saved addresses.

Core concepts:

```text
Customer
CustomerAddress
CustomerPreference
```

Orders and shipments may store address snapshots for historical accuracy.

---

## 13. Review Domain

Review Service owns product reviews and moderation state.

Core concepts:

```text
Review
RatingSummary
ModerationDecision
ReviewReport
```

Review Service references products by ID.

---

## 14. Recommendation Domain

Recommendation Service owns recommendation signals and outputs.

Core concepts:

```text
RecommendationSignal
RecommendationRule
RecommendationResult
RecommendationFeedback
```

Recommendation Service consumes product, order, review, and basket events but does not own those domains.

---

## 15. Relationship Summary

| Concept | Source of Truth | Notes |
|---|---|---|
| Product | Catalogue Service | Includes common product data |
| Category | Catalogue Service | Defines product taxonomy |
| Product Attribute Definition | Catalogue Service | Defines category-specific attributes |
| Product Attribute Value | Catalogue Service | Stores product-specific values |
| Product Search Document | Search Service | Denormalised projection |
| Stock Level | Inventory Service | Product ID reference only |
| Basket | Basket Service | Current shopping intent |
| Order | Order Service | Order lifecycle |
| Order Item Snapshot | Order Service | Historical product/price/attribute summary |
| Payment | Payment Service | Payment state |
| Shipment | Shipping Service | Fulfilment state |
| Notification | Notification Service | Delivery state |
| Review | Review Service | Review content and moderation |
| Recommendation Result | Recommendation Service | Derived output |

---

## 16. Domain Events

Important events include:

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
ProductAttributeDefinitionCreated
ProductAttributeDefinitionUpdated
StockReserved
StockReservationFailed
PaymentAuthorised
PaymentFailed
ShipmentCreated
ShipmentFailed
OrderCreated
OrderFailed
NotificationSent
NotificationFailed
```

Catalogue events are especially important for Search Service projections.

---

## 17. Anti-Patterns to Avoid

Avoid:

```text
one product table with hundreds of nullable product-specific columns
uncontrolled product JSON with no category governance
Search Service owning product truth
Inventory Service owning product details
Order Service depending on live product data for historical orders
Basket Service reserving stock
API Gateway owning checkout business rules
```

---

## 18. Open Questions

| Question | Status |
|---|---|
| Should variant-level attributes be implemented in version one? | To decide |
| Should attribute options use a table or JSON initially? | To decide |
| Should Search Service be MySQL-backed first or use a dedicated search engine later? | To decide |
| Which product attributes are required for first seed categories? | To decide |
| Should order item snapshots include selected attributes as JSON? | Proposed |

---

## 19. Related Documents

```text
docs/architecture/service-boundaries.md
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/events/event-catalog.md
docs/requirements/business-rules.md
proto/acme/catalog/v1/README.md
```

---

## 20. Summary

The bfstore domain model keeps Catalogue Service as the governed product source of truth while allowing varied product types through category-scoped attributes.

Search Service owns denormalised product search projections, not product truth.

This gives bfstore flexibility for homeware, curtains, bed frames, and other product types without compromising service ownership or data governance.
