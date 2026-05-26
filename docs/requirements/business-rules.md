# Business Rules

## 1. Purpose

This document defines the core business rules for **bfstore**, ACME Ltd’s fictional online furniture store backend.

Business rules describe constraints, decisions, and invariants that must hold true regardless of implementation detail. They guide service logic, API behaviour, database design, event handling, and tests.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s domain modelling and business correctness.

---

## 2. Scope

This document covers business rules for:

```text
catalogue
basket
inventory
checkout
orders
payments
shipping
notifications
customers
reviews
search
recommendations
events
data ownership
security and privacy
```

---

## 3. Business Rule Format

Each rule includes:

| Field | Description |
|---|---|
| ID | Stable business rule identifier |
| Rule | Business rule statement |
| Priority | Must, Should, Could |
| Owner | Primary service or domain owner |
| Notes | Additional clarification |

---

## 4. Catalogue Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-CAT-001 | Only active products may be shown in normal customer purchase flows. | Must | `catalog-service` | Inactive products must not be purchasable |
| BR-CAT-002 | A product must have a stable product ID before it can be referenced by another service. | Must | `catalog-service` | Used by basket/order/inventory |
| BR-CAT-003 | A product must have a name, status, price, and currency before it can be active. | Must | `catalog-service` | Minimum catalogue integrity |
| BR-CAT-004 | Product price must be represented in minor units and a currency code. | Must | `catalog-service` | Avoid floating point money |
| BR-CAT-005 | Product status changes should be published as events when downstream projections exist. | Should | `catalog-service` | Search/recommendation |
| BR-CAT-006 | Catalogue Service owns product truth; Search Service may only own projections. | Must | `catalog-service`, `search-service` | Data ownership rule |
| BR-CAT-007 | A product may be visible but unavailable if stock is zero. | Should | `catalog-service`, `inventory-service` | Availability belongs to Inventory |

---

## 5. Basket Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-BAS-001 | A basket item quantity must be greater than zero. | Must | `basket-service` | Validation |
| BR-BAS-002 | A basket must not contain inactive products. | Must | `basket-service`, `catalog-service` | Validate before add/checkout |
| BR-BAS-003 | Basket Service does not reserve stock. | Must | `basket-service` | Stock reservation belongs to Inventory |
| BR-BAS-004 | Adding an item already in the basket should update or merge quantity according to documented behaviour. | Should | `basket-service` | Avoid duplicate lines unless variant differs |
| BR-BAS-005 | A checked-out basket should not be mutated. | Should | `basket-service` | Prevent post-checkout changes |
| BR-BAS-006 | Basket prices are not final order prices unless explicitly snapshotted at checkout. | Must | `basket-service`, `order-service` | Avoid stale pricing assumptions |
| BR-BAS-007 | Abandoned baskets may expire after a defined period. | Could | `basket-service` | Retention rule |

---

## 6. Inventory Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-INV-001 | Inventory Service is the source of truth for stock levels. | Must | `inventory-service` | No other service writes stock |
| BR-INV-002 | Stock must not be reserved below zero. | Must | `inventory-service` | Prevent overselling |
| BR-INV-003 | Stock reservation must be idempotent for retried checkout operations. | Must | `inventory-service` | Avoid duplicate reservation |
| BR-INV-004 | A stock reservation must have an expiry time or explicit release path. | Must | `inventory-service` | Avoid locked stock |
| BR-INV-005 | Stock must be released if checkout fails after reservation and before final commitment. | Must | `inventory-service`, `order-service` | Compensation |
| BR-INV-006 | Inventory events should be published for stock reservation, reservation failure, and reservation release. | Should | `inventory-service` | Event visibility |
| BR-INV-007 | Inventory availability may be eventually consistent in search or recommendation projections. | Should | `inventory-service`, `search-service` | Projection rule |

---

## 7. Checkout Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-CHK-001 | Checkout may only proceed for a non-empty basket. | Must | `order-service`, `basket-service` | Core validation |
| BR-CHK-002 | Checkout must validate that basket products are still purchasable. | Must | `order-service`, `catalog-service` | Product may become inactive |
| BR-CHK-003 | Stock must be reserved before payment is authorised. | Must | `order-service`, `inventory-service` | Avoid charging for unavailable stock |
| BR-CHK-004 | Payment must be authorised before an order is confirmed. | Must | `order-service`, `payment-service` | Avoid unpaid confirmed orders |
| BR-CHK-005 | Shipment creation or fulfilment state must follow a documented decision before order confirmation. | Must | `order-service`, `shipping-service` | Pending fulfilment vs checkout failure |
| BR-CHK-006 | A duplicate checkout request with the same idempotency key must not create a duplicate confirmed order. | Must | `order-service` | Critical idempotency |
| BR-CHK-007 | If stock reservation fails, payment must not be attempted. | Must | `order-service` | Failure ordering |
| BR-CHK-008 | If payment fails, the order must not be confirmed. | Must | `order-service`, `payment-service` | Failure safety |
| BR-CHK-009 | If payment fails after stock reservation, stock must be released or allowed to expire. | Must | `order-service`, `inventory-service` | Compensation |
| BR-CHK-010 | Notification failure must not roll back order creation. | Must | `order-service`, `notification-service` | Async side effect |
| BR-CHK-011 | Checkout state and failure reason should be recorded for operational diagnosis. | Should | `order-service` | Supportability |

---

## 8. Order Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-ORD-001 | Order Service owns order lifecycle. | Must | `order-service` | Data ownership |
| BR-ORD-002 | Order items must preserve product and price snapshots from checkout time. | Must | `order-service` | Historical accuracy |
| BR-ORD-003 | Order totals must be calculated from order item snapshots and adjustments. | Must | `order-service` | Integrity |
| BR-ORD-004 | An order must have a stable order ID before events referencing it are published. | Must | `order-service` | Event contract |
| BR-ORD-005 | `OrderCreated` must only be published after the order exists in Order Service. | Must | `order-service` | Event correctness |
| BR-ORD-006 | Order cancellation must follow valid status transitions. | Should | `order-service` | Lifecycle |
| BR-ORD-007 | Cancelled or failed orders must not be treated as confirmed. | Must | `order-service` | State integrity |
| BR-ORD-008 | Order status changes should be recorded in status history. | Should | `order-service` | Audit |
| BR-ORD-009 | Orders must not directly own payment, stock, shipment, or notification internals. | Must | `order-service` | Boundary rule |

---

## 9. Payment Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-PAY-001 | Payment Service owns payment state. | Must | `payment-service` | Data ownership |
| BR-PAY-002 | Raw payment card data must not be stored. | Must | `payment-service` | Security |
| BR-PAY-003 | Payment authorisation must be idempotent. | Must | `payment-service` | Avoid duplicate charges |
| BR-PAY-004 | A declined payment must not create a confirmed order. | Must | `payment-service`, `order-service` | Checkout safety |
| BR-PAY-005 | Payment attempts must be recorded. | Must | `payment-service` | Audit/reconciliation |
| BR-PAY-006 | Payment provider references must be treated as sensitive operational data. | Should | `payment-service` | Security |
| BR-PAY-007 | Payment timeout outcomes must be reconciled safely. | Should | `payment-service` | Provider ambiguity |
| BR-PAY-008 | Payment failure details exposed externally must be safe and non-sensitive. | Must | `payment-service`, `api-gateway` | Error model |

---

## 10. Shipping Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-SHP-001 | Shipping Service owns shipment state. | Must | `shipping-service` | Data ownership |
| BR-SHP-002 | Shipment creation must be idempotent. | Must | `shipping-service` | Avoid duplicate shipments |
| BR-SHP-003 | Shipments should store delivery address snapshots. | Should | `shipping-service` | Historical fulfilment |
| BR-SHP-004 | Shipment status transitions must be valid. | Should | `shipping-service` | Lifecycle integrity |
| BR-SHP-005 | Shipment failure behaviour during checkout must be explicitly documented. | Must | `shipping-service`, `order-service` | Pending vs fail |
| BR-SHP-006 | Shipment events should be published for created, dispatched, delivered, and failed states. | Should | `shipping-service` | Event visibility |

---

## 11. Notification Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-NOT-001 | Notification Service owns notification delivery state. | Must | `notification-service` | Data ownership |
| BR-NOT-002 | Notification failure must not roll back the business action that caused it. | Must | `notification-service` | Async side effect |
| BR-NOT-003 | Notification consumers must be idempotent. | Must | `notification-service` | Avoid duplicate sends |
| BR-NOT-004 | The same `OrderCreated` event should not produce duplicate order confirmations. | Must | `notification-service` | Event deduplication |
| BR-NOT-005 | Notification attempts should be recorded. | Should | `notification-service` | Audit/support |
| BR-NOT-006 | Notification payloads should minimise PII. | Must | `notification-service` | Privacy |
| BR-NOT-007 | Notification templates should be versioned or managed carefully where implemented. | Could | `notification-service` | Later maturity |

---

## 12. Customer and Authentication Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-AUTH-001 | Auth Service owns authentication identity and credentials. | Should | `auth-service` | Deferred if auth is stubbed |
| BR-AUTH-002 | Customer Service owns customer profile and address data. | Should | `customer-service` | Data ownership |
| BR-AUTH-003 | Customers must not access other customers’ private orders or profiles. | Must | `auth-service`, `order-service`, `customer-service` | Authorisation |
| BR-AUTH-004 | Passwords must never be stored in plain text. | Must | `auth-service` | Security |
| BR-AUTH-005 | Tokens and credentials must not be logged. | Must | `auth-service`, all services | Secure logging |
| BR-CUS-001 | Address updates must not change historical order delivery snapshots. | Must | `customer-service`, `order-service`, `shipping-service` | Historical accuracy |

---

## 13. Review Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-REV-001 | Review Service owns review content and moderation state. | Could | `review-service` | Later service |
| BR-REV-002 | Only approved reviews should be visible publicly. | Could | `review-service` | Moderation |
| BR-REV-003 | Reviews must be associated with a valid product ID. | Could | `review-service`, `catalog-service` | Cross-service validation |
| BR-REV-004 | Review summaries may be eventually consistent. | Could | `review-service` | Projection |
| BR-REV-005 | Review events should be published when reviews are approved. | Could | `review-service` | Search/recommendation |

---

## 14. Search Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-SRCH-001 | Search Service owns search projections, not product truth. | Should | `search-service` | Catalogue remains source of truth |
| BR-SRCH-002 | Search results must exclude inactive products. | Must | `search-service`, `catalog-service` | Customer safety |
| BR-SRCH-003 | Search results may be eventually consistent with catalogue updates. | Should | `search-service` | Projection behaviour |
| BR-SRCH-004 | Search projections should be rebuildable. | Should | `search-service` | Operations |
| BR-SRCH-005 | Search failure must not block checkout. | Must | `search-service` | Non-critical path |

---

## 15. Recommendation Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-REC-001 | Recommendation Service owns recommendation outputs. | Could | `recommendation-service` | Later service |
| BR-REC-002 | Recommendations must not include inactive products. | Must | `recommendation-service`, `catalog-service` | Customer safety |
| BR-REC-003 | Recommendation failure must not block browsing or checkout. | Must | `recommendation-service` | Resilience |
| BR-REC-004 | Initial recommendations may be rules-based. | Could | `recommendation-service` | Avoid premature ML |
| BR-REC-005 | Recommendation signals may be eventually consistent. | Could | `recommendation-service` | Event-driven |

---

## 16. Event Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-EVT-001 | Events describe facts that have already happened. | Must | All producers | Not commands |
| BR-EVT-002 | The service that owns the business fact publishes the event. | Must | All producers | Ownership |
| BR-EVT-003 | Events must include a unique event ID. | Must | All producers | Deduplication |
| BR-EVT-004 | Events must include a version. | Must | All producers | Compatibility |
| BR-EVT-005 | Events must include correlation context. | Must | All producers | Traceability |
| BR-EVT-006 | Event consumers must handle duplicate events safely. | Must | All consumers | Kafka reliability |
| BR-EVT-007 | Unsupported or invalid events should be handled through retry/DLQ strategy. | Should | All consumers | Operations |
| BR-EVT-008 | Event payloads must avoid unnecessary sensitive data. | Must | All producers | Privacy/security |

---

## 17. Data Ownership Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-DATA-001 | Each service owns its own data schema. | Must | All services | Service-owned database |
| BR-DATA-002 | Services must not directly access another service’s database. | Must | All services | Boundary rule |
| BR-DATA-003 | Cross-service references should use IDs, not database foreign keys. | Must | All services | Decoupling |
| BR-DATA-004 | Snapshots are allowed when historical accuracy is required. | Must | Order/Shipping/Payment | Product/address/payment reference snapshots |
| BR-DATA-005 | Projections must not become hidden sources of truth. | Must | Search/Recommendation/etc. | Projection rule |
| BR-DATA-006 | Database migrations are owned by the service that owns the schema. | Must | All services | Migration governance |

---

## 18. Security and Privacy Rules

| ID | Rule | Priority | Owner | Notes |
|---|---|---:|---|---|
| BR-SEC-001 | Secrets must not be committed to source control. | Must | All services | DevSecOps |
| BR-SEC-002 | Tokens, passwords, and raw payment data must not be logged. | Must | All services | Secure logging |
| BR-SEC-003 | Service database users must follow least privilege. | Must | All services | One service, one schema |
| BR-SEC-004 | External error responses must not expose internal implementation details. | Must | `api-gateway`, all services | Error model |
| BR-SEC-005 | Events must not contain unnecessary PII. | Must | All producers | Privacy |
| BR-SEC-006 | Customer data copied into snapshots must be justified by business need. | Must | `order-service`, `shipping-service` | Data minimisation |

---

## 19. Initial Vertical Slice Business Rules

The first implementation must enforce:

```text
BR-CAT-001
BR-BAS-001
BR-BAS-002
BR-BAS-003
BR-INV-001
BR-INV-002
BR-INV-003
BR-CHK-001
BR-CHK-003
BR-CHK-004
BR-CHK-006
BR-CHK-007
BR-CHK-008
BR-CHK-010
BR-ORD-001
BR-ORD-002
BR-ORD-005
BR-PAY-002
BR-PAY-003
BR-PAY-004
BR-SHP-002
BR-NOT-002
BR-NOT-003
BR-EVT-001
BR-EVT-003
BR-DATA-001
BR-DATA-002
BR-SEC-001
BR-SEC-002
```

---

## 20. Business Rule Traceability

Business rules should be linked to:

```text
functional requirements
acceptance criteria
service tests
contract tests
database constraints
event contracts
runbooks
```

Example:

```text
BR-CHK-007: If stock reservation fails, payment must not be attempted.
    -> FR-INV-004
    -> FR-ORD-003
    -> AC-CHK-004
    -> order-service integration test
    -> insufficient stock E2E test
```

---

## 21. Open Questions

| Question | Status |
|---|---|
| Should shipment creation failure fail checkout or create pending fulfilment? | To decide |
| Should payment capture be separate from authorisation in version one? | To decide |
| Should customer authentication be mandatory for checkout in the first slice? | To decide |
| How long should stock reservations remain valid? | To decide |
| Should basket prices be refreshed at checkout or preserved from add-to-basket time? | To decide |
| What refund and cancellation rules apply after dispatch? | To decide later |
| What review eligibility rules should apply? | To decide later |

---

## 22. Related Documents

This document should be read alongside:

```text
docs/requirements/functional-requirements.md
docs/requirements/non-functional-requirements.md
docs/requirements/acceptance-criteria.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/architecture/resilience-patterns.md
docs/data/data-ownership.md
docs/events/event-catalog.md
docs/testing/testing-strategy.md
```

---

## 23. Summary

bfstore’s business rules protect the core integrity of the commerce domain.

The most important rules are:

```text
only active products can be purchased
basket quantities must be valid
stock must be reserved before payment
payment must be authorised before order confirmation
duplicate checkout must not create duplicate orders
payment failure must not create confirmed orders
notification failure must not roll back orders
each service owns its own data
events describe facts
sensitive data must be protected
```

These rules should guide implementation, tests, service contracts, database constraints, and operational runbooks.
