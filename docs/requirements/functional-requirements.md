# Functional Requirements

## 1. Purpose

This document defines the functional requirements for **bfstore**, ACME Ltd’s fictional online furniture store backend.

Functional requirements describe what the system must do from a business and user journey perspective. They focus on observable behaviours, service responsibilities, and expected system capabilities.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s product and engineering scope.

---

## 2. Product Context

bfstore is a cloud-native microservice backend for an online furniture store.

The system supports a customer journey where a user can:

```text
browse furniture products
view product details
add items to a basket
checkout
reserve stock
authorise payment
create an order
create a shipment
receive notification
view order details
```

The first implementation should focus on the checkout vertical slice.

---

## 3. Scope

This document covers functional requirements for:

```text
catalogue browsing
product details
basket management
checkout
inventory reservation
payment authorisation
order creation
shipping and fulfilment
notifications
customer profile support
authentication and authorisation
reviews
search
recommendations
administrative and operational behaviours
```

Not every requirement must be implemented in the first version. Each requirement includes a priority.

---

## 4. Requirement Priority

Priorities use the following scale:

| Priority | Meaning |
|---|---|
| Must | Required for the initial credible vertical slice |
| Should | Important for a mature version, but may follow the first slice |
| Could | Valuable enhancement once core flows are stable |
| Won't Yet | Explicitly deferred from current scope |

---

## 5. Functional Requirement Format

Each requirement includes:

| Field | Description |
|---|---|
| ID | Stable requirement identifier |
| Requirement | The expected system behaviour |
| Priority | Must, Should, Could, or Won't Yet |
| Primary Service | Service most responsible for the behaviour |
| Notes | Clarifications, dependencies, or design considerations |

---

## 6. Catalogue Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-CAT-001 | The system shall allow clients to list active furniture products. | Must | `catalog-service` | Initial browse journey |
| FR-CAT-002 | The system shall allow clients to retrieve product details by product ID. | Must | `catalog-service` | Required for basket add and checkout validation |
| FR-CAT-003 | The system shall distinguish active and inactive products. | Must | `catalog-service` | Inactive products must not be purchasable |
| FR-CAT-004 | The system shall support product categories. | Should | `catalog-service` | Example: sofa, table, wardrobe |
| FR-CAT-005 | The system shall support product variants such as colour, size, or material. | Should | `catalog-service` | Useful for furniture options |
| FR-CAT-006 | The system shall expose product price and currency. | Must | `catalog-service` | Prices should use minor units internally |
| FR-CAT-007 | The system shall publish an event when a product is created, updated, activated, or deactivated. | Should | `catalog-service` | Supports search and recommendations |
| FR-CAT-008 | The system shall support product images or image references. | Could | `catalog-service` | May be simplified initially |
| FR-CAT-009 | The system shall allow product filtering by category or status. | Should | `catalog-service` or `search-service` | Basic catalogue filtering first |
| FR-CAT-010 | The system shall prevent inactive products from appearing in customer purchase flows. | Must | `catalog-service`, `basket-service`, `order-service` | Enforced before checkout |

---

## 7. Basket Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-BAS-001 | The system shall allow a customer or session to create a basket. | Must | `basket-service` | Supports guest or authenticated flow |
| FR-BAS-002 | The system shall allow a product item to be added to a basket. | Must | `basket-service` | Product validity should be checked |
| FR-BAS-003 | The system shall allow basket item quantity to be updated. | Must | `basket-service` | Quantity must be positive |
| FR-BAS-004 | The system shall allow a basket item to be removed. | Must | `basket-service` | Required for normal shopping behaviour |
| FR-BAS-005 | The system shall allow a basket to be retrieved. | Must | `basket-service` | Used by clients and checkout |
| FR-BAS-006 | The system shall prevent inactive products from being added to the basket. | Must | `basket-service`, `catalog-service` | May call Catalogue Service |
| FR-BAS-007 | The system shall prevent invalid quantities such as zero or negative values. | Must | `basket-service` | Validation requirement |
| FR-BAS-008 | The system shall mark a basket as checked out when an order is successfully created. | Should | `basket-service`, `order-service` | Prevents further mutation |
| FR-BAS-009 | The system shall support basket expiry for abandoned baskets. | Could | `basket-service` | Useful operational behaviour |
| FR-BAS-010 | The system shall publish basket lifecycle events where required. | Could | `basket-service` | Useful for recommendations |

---

## 8. Inventory Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-INV-001 | The system shall maintain stock levels for products or variants. | Must | `inventory-service` | Inventory owns stock truth |
| FR-INV-002 | The system shall allow stock availability to be checked. | Must | `inventory-service` | Required for checkout |
| FR-INV-003 | The system shall reserve stock during checkout. | Must | `inventory-service` | Prevents overselling |
| FR-INV-004 | The system shall reject reservation requests when stock is insufficient. | Must | `inventory-service` | Checkout must fail safely |
| FR-INV-005 | The system shall release stock reservations when checkout fails or is cancelled. | Must | `inventory-service` | Compensation path |
| FR-INV-006 | The system shall expire stale stock reservations. | Should | `inventory-service` | Prevents permanent stock lock |
| FR-INV-007 | The system shall publish `StockReserved` when stock is reserved. | Must | `inventory-service` | Event-driven observability |
| FR-INV-008 | The system shall publish `StockReservationFailed` when reservation fails. | Must | `inventory-service` | Supports checkout failure handling |
| FR-INV-009 | The system shall publish `StockReservationReleased` when reserved stock is released. | Should | `inventory-service` | Supports operations and order updates |
| FR-INV-010 | The system shall prevent stock levels from becoming negative. | Must | `inventory-service` | Data integrity invariant |

---

## 9. Checkout and Order Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-ORD-001 | The system shall allow checkout to be submitted for a valid basket. | Must | `order-service` | Core vertical slice |
| FR-ORD-002 | The system shall validate basket contents before creating an order. | Must | `order-service`, `basket-service` | Basket must not be empty |
| FR-ORD-003 | The system shall reserve stock before confirming an order. | Must | `order-service`, `inventory-service` | Prevents orders without stock |
| FR-ORD-004 | The system shall authorise payment before confirming an order. | Must | `order-service`, `payment-service` | Prevents unpaid confirmed orders |
| FR-ORD-005 | The system shall create a shipment or fulfilment record as part of the checkout flow. | Must | `order-service`, `shipping-service` | Initial implementation may simulate |
| FR-ORD-006 | The system shall create an order after required checkout steps succeed. | Must | `order-service` | Order Service owns order lifecycle |
| FR-ORD-007 | The system shall store order item snapshots. | Must | `order-service` | Product and price history |
| FR-ORD-008 | The system shall publish `OrderCreated` when an order is created. | Must | `order-service` | Notifications consume this |
| FR-ORD-009 | The system shall publish `OrderFailed` when checkout fails after an order attempt. | Should | `order-service` | Operational visibility |
| FR-ORD-010 | The system shall allow an order to be retrieved by order ID. | Must | `order-service` | Customer support and UI |
| FR-ORD-011 | The system shall allow customer order history to be listed. | Should | `order-service` | Useful after identity/customer support |
| FR-ORD-012 | The system shall prevent duplicate confirmed orders for the same checkout idempotency key. | Must | `order-service` | Critical idempotency requirement |
| FR-ORD-013 | The system shall support order cancellation where business rules allow. | Should | `order-service` | Later order lifecycle |
| FR-ORD-014 | The system shall record order status history. | Should | `order-service` | Audit and support |

---

## 10. Payment Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-PAY-001 | The system shall authorise payment for a checkout. | Must | `payment-service` | May be simulated initially |
| FR-PAY-002 | The system shall record payment attempts. | Must | `payment-service` | Audit and troubleshooting |
| FR-PAY-003 | The system shall reject declined payment attempts. | Must | `payment-service` | Checkout failure case |
| FR-PAY-004 | The system shall publish `PaymentAuthorised` when payment succeeds. | Must | `payment-service` | Event catalogue |
| FR-PAY-005 | The system shall publish `PaymentFailed` when payment fails. | Must | `payment-service` | Compensation path |
| FR-PAY-006 | The system shall avoid duplicate payment authorisation for retried requests. | Must | `payment-service` | Idempotency |
| FR-PAY-007 | The system shall not store raw payment card data. | Must | `payment-service` | Security requirement |
| FR-PAY-008 | The system shall support payment capture separately from authorisation. | Could | `payment-service` | May be deferred |
| FR-PAY-009 | The system shall support refunds. | Could | `payment-service` | Later order lifecycle |
| FR-PAY-010 | The system shall support payment provider simulation for local development and tests. | Must | `payment-service` | Keeps first implementation achievable |

---

## 11. Shipping Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-SHP-001 | The system shall provide available delivery options. | Must | `shipping-service` | Can be static initially |
| FR-SHP-002 | The system shall create a shipment for a confirmed checkout flow. | Must | `shipping-service` | Initial provider can be simulated |
| FR-SHP-003 | The system shall record shipment status. | Must | `shipping-service` | Fulfilment lifecycle |
| FR-SHP-004 | The system shall publish `ShipmentCreated` when a shipment is created. | Must | `shipping-service` | Event catalogue |
| FR-SHP-005 | The system shall publish `ShipmentFailed` when shipment creation fails. | Must | `shipping-service` | Checkout failure/pending fulfilment |
| FR-SHP-006 | The system shall support shipment tracking references. | Should | `shipping-service` | Can be simulated |
| FR-SHP-007 | The system shall support shipment dispatched and delivered states. | Should | `shipping-service` | Later lifecycle |
| FR-SHP-008 | The system shall create shipments idempotently. | Must | `shipping-service` | Prevent duplicate shipments |
| FR-SHP-009 | The system shall store delivery address snapshots. | Should | `shipping-service` | Historical fulfilment accuracy |

---

## 12. Notification Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-NOT-001 | The system shall consume `OrderCreated` events. | Must | `notification-service` | Core async side effect |
| FR-NOT-002 | The system shall create a notification record when an order confirmation is required. | Must | `notification-service` | Delivery tracking |
| FR-NOT-003 | The system shall simulate sending an order confirmation notification. | Must | `notification-service` | Initial implementation |
| FR-NOT-004 | The system shall publish `NotificationSent` when a notification is sent. | Should | `notification-service` | Operational visibility |
| FR-NOT-005 | The system shall publish `NotificationFailed` when notification delivery fails. | Should | `notification-service` | DLQ/retry support |
| FR-NOT-006 | The system shall avoid duplicate notifications for the same event. | Must | `notification-service` | Idempotency |
| FR-NOT-007 | The system shall record notification attempts. | Should | `notification-service` | Audit and support |
| FR-NOT-008 | The system shall support notification templates. | Could | `notification-service` | Later maturity |
| FR-NOT-009 | Notification failure shall not roll back order creation. | Must | `notification-service`, `order-service` | Important resilience rule |

---

## 13. Authentication and Customer Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-AUTH-001 | The system shall support customer authentication. | Should | `auth-service` | May be deferred for first slice |
| FR-AUTH-002 | The system shall issue tokens for authenticated customers. | Should | `auth-service` | Token model to be documented |
| FR-AUTH-003 | The system shall protect customer-specific order access. | Should | `auth-service`, `order-service` | Prevent cross-customer access |
| FR-CUS-001 | The system shall store customer profile data. | Should | `customer-service` | Name and contact details |
| FR-CUS-002 | The system shall store customer delivery addresses. | Should | `customer-service` | Used by checkout |
| FR-CUS-003 | The system shall allow a customer to update saved addresses. | Could | `customer-service` | Later customer account feature |
| FR-CUS-004 | The system shall avoid exposing unnecessary customer PII in events. | Must | All services | Security/data rule |

---

## 14. Search Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-SRCH-001 | The system shall support product search. | Should | `search-service` | May follow catalogue filtering |
| FR-SRCH-002 | The system shall support search by product name or description. | Should | `search-service` | Initial search behaviour |
| FR-SRCH-003 | The system shall support filters such as category, material, colour, or price range. | Could | `search-service` | Enhanced search |
| FR-SRCH-004 | The system shall consume catalogue events to maintain search projections. | Should | `search-service` | Event-driven projection |
| FR-SRCH-005 | The system shall exclude inactive products from search results. | Must | `search-service`, `catalog-service` | Purchase safety |
| FR-SRCH-006 | The system shall support projection rebuild where practical. | Should | `search-service` | Operational maturity |

---

## 15. Review Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-REV-001 | The system shall allow customers to submit product reviews. | Could | `review-service` | Later feature |
| FR-REV-002 | The system shall associate reviews with products. | Could | `review-service` | Product ID reference |
| FR-REV-003 | The system shall support review moderation. | Could | `review-service` | Approved/rejected states |
| FR-REV-004 | The system shall publish `ReviewApproved` when a review is approved. | Could | `review-service` | Search/recommendation input |
| FR-REV-005 | The system shall calculate product rating summaries. | Could | `review-service` | Projection/aggregate |

---

## 16. Recommendation Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-REC-001 | The system shall provide product recommendations. | Could | `recommendation-service` | Later feature |
| FR-REC-002 | The system shall support rule-based recommendations initially. | Could | `recommendation-service` | Same category/material |
| FR-REC-003 | The system shall consume product and order events as recommendation signals. | Could | `recommendation-service` | Event-driven signals |
| FR-REC-004 | The system shall not recommend inactive products. | Must | `recommendation-service`, `catalog-service` | Safety rule |
| FR-REC-005 | Recommendation failure shall not block browsing or checkout. | Must | `recommendation-service` | Resilience rule |

---

## 17. API Gateway Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-GW-001 | The system shall expose a client-facing API through an API Gateway. | Must | `api-gateway` | External entry point |
| FR-GW-002 | The API Gateway shall route catalogue requests to Catalogue Service. | Must | `api-gateway` | Browse flow |
| FR-GW-003 | The API Gateway shall route basket requests to Basket Service. | Must | `api-gateway` | Basket flow |
| FR-GW-004 | The API Gateway shall route checkout requests to Order Service. | Must | `api-gateway` | Checkout flow |
| FR-GW-005 | The API Gateway shall propagate correlation IDs. | Must | `api-gateway` | Observability |
| FR-GW-006 | The API Gateway shall map internal errors to safe client-facing responses. | Must | `api-gateway` | Error model |
| FR-GW-007 | The API Gateway shall not own domain business rules. | Must | `api-gateway` | Boundary rule |

---

## 18. Event and Integration Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-EVT-001 | The system shall use Kafka for asynchronous business events. | Must | Event producers | Core architecture |
| FR-EVT-002 | Events shall use a consistent event envelope. | Must | All producers | See event envelope doc |
| FR-EVT-003 | Event consumers shall be idempotent. | Must | All consumers | Duplicate safety |
| FR-EVT-004 | Consumers shall send poison messages to a DLQ after retry exhaustion. | Should | All consumers | Operational maturity |
| FR-EVT-005 | Events shall carry correlation context. | Must | All producers | Traceability |
| FR-EVT-006 | Event contracts shall be versioned. | Must | All producers | Compatibility |
| FR-EVT-007 | Event producers shall publish only events for facts they own. | Must | All producers | Ownership rule |

---

## 19. Observability Functional Requirements

| ID | Requirement | Priority | Primary Service | Notes |
|---|---|---:|---|---|
| FR-OBS-001 | Services shall emit structured logs. | Must | All services | Include correlation ID |
| FR-OBS-002 | Services shall expose health checks. | Must | All services | Local and Kubernetes readiness |
| FR-OBS-003 | Services shall expose readiness checks. | Must | All services | Traffic safety |
| FR-OBS-004 | Services shall emit basic metrics. | Should | All services | Request count, latency, errors |
| FR-OBS-005 | Services shall propagate trace context. | Should | All services | OpenTelemetry |
| FR-OBS-006 | Kafka consumers shall expose consumer lag metrics where possible. | Should | Event consumers | Operational maturity |

---

## 20. Initial Vertical Slice Requirements

The first implementation should prioritise:

```text
FR-CAT-001
FR-CAT-002
FR-CAT-003
FR-CAT-006
FR-BAS-001
FR-BAS-002
FR-BAS-003
FR-BAS-005
FR-BAS-006
FR-INV-001
FR-INV-003
FR-INV-004
FR-ORD-001
FR-ORD-003
FR-ORD-004
FR-ORD-006
FR-ORD-008
FR-ORD-012
FR-PAY-001
FR-PAY-002
FR-PAY-003
FR-SHP-002
FR-SHP-004
FR-NOT-001
FR-NOT-003
FR-NOT-006
FR-GW-001
FR-GW-005
```

This proves the core system behaviour:

```text
Browse product
    -> Add to basket
    -> Checkout
    -> Reserve stock
    -> Authorise payment
    -> Create order
    -> Create shipment
    -> Publish event
    -> Send notification
```

---

## 21. Traceability

Functional requirements should be traceable to:

```text
user journeys
service requirements
API contracts
event contracts
database design
tests
acceptance criteria
```

Example:

```text
FR-ORD-003: The system shall reserve stock before confirming an order.
    -> Inventory Service ReserveStock API
    -> StockReserved event
    -> order-service checkout test
    -> insufficient stock E2E test
```

---

## 22. Open Questions

| Question | Status |
|---|---|
| Will the first client-facing API be REST/JSON or gRPC-Web? | To decide |
| Will customer authentication be included in the first vertical slice? | Proposed: defer or stub |
| Will shipment creation failure fail checkout or create pending fulfilment? | To decide |
| Will payment capture be separate from authorisation in the first version? | To decide |
| Should Search Service be implemented before or after checkout? | Proposed: after checkout |
| Should Recommendation Service be implemented after order events exist? | Proposed |
| What level of admin product management is required? | To decide |

---

## 23. Related Documents

This document should be read alongside:

```text
docs/requirements/product-vision.md
docs/requirements/scope.md
docs/requirements/user-journeys.md
docs/requirements/non-functional-requirements.md
docs/requirements/business-rules.md
docs/requirements/acceptance-criteria.md
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/events/event-catalog.md
docs/data/data-ownership.md
docs/testing/testing-strategy.md
```

---

## 24. Summary

bfstore’s functional requirements define what the system must do from a business perspective.

The most important initial capability is the checkout vertical slice:

```text
product browse
basket management
stock reservation
payment authorisation
order creation
shipment creation
notification
```

This provides a strong foundation for later capabilities such as authentication, customer profiles, reviews, search, recommendations, advanced fulfilment, and operational tooling.
