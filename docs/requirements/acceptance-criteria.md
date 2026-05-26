# Acceptance Criteria

## 1. Purpose

This document defines acceptance criteria for **bfstore**, ACME Ltd’s fictional online furniture store backend.

Acceptance criteria describe how to verify that functional and non-functional requirements have been satisfied. They provide testable conditions for user journeys, service behaviours, failure cases, and operational expectations.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s delivery readiness.

---

## 2. Scope

This document covers acceptance criteria for:

```text
catalogue browsing
basket management
checkout
stock reservation
payment authorisation
order creation
shipment creation
notification
error handling
events
data ownership
observability
testing
deployment readiness
```

The first priority is the checkout vertical slice.

---

## 3. Acceptance Criteria Format

Acceptance criteria use the following structure:

```text
Given <initial context>
When <action occurs>
Then <expected outcome>
And <additional expected outcomes>
```

Each criterion includes:

| Field | Description |
|---|---|
| ID | Stable acceptance criterion identifier |
| Priority | Must, Should, or Could |
| Related Requirements | Functional or non-functional requirement references |
| Verification | Suggested test type |

---

## 4. Catalogue Acceptance Criteria

## AC-CAT-001: List Active Products

**Priority:** Must  
**Related Requirements:** `FR-CAT-001`, `BR-CAT-001`  
**Verification:** API contract test, integration test, E2E test

```text
Given active products exist in the catalogue
When a client requests the product list
Then the system returns active products
And inactive products are not returned in the normal customer product list
And each product includes product ID, name, price, currency, and status
```

---

## AC-CAT-002: Get Product Details

**Priority:** Must  
**Related Requirements:** `FR-CAT-002`  
**Verification:** API contract test, integration test

```text
Given an active product exists
When a client requests product details by product ID
Then the system returns the product details
And the response includes product ID, name, description, price, currency, and status
```

---

## AC-CAT-003: Product Not Found

**Priority:** Must  
**Related Requirements:** `FR-CAT-002`, `NFR-SEC-004`  
**Verification:** API contract test

```text
Given no product exists for the requested product ID
When a client requests product details
Then the system returns a safe not found error
And the response does not expose internal database details
```

---

## AC-CAT-004: Inactive Product Cannot Be Purchased

**Priority:** Must  
**Related Requirements:** `FR-CAT-010`, `BR-CAT-001`  
**Verification:** Integration test, E2E test

```text
Given a product is inactive
When a customer attempts to add the product to a basket or checkout with it
Then the system rejects the operation
And the product is not included in a confirmed order
```

---

## 5. Basket Acceptance Criteria

## AC-BAS-001: Create Basket

**Priority:** Must  
**Related Requirements:** `FR-BAS-001`  
**Verification:** API contract test

```text
Given a customer or session has no active basket
When a basket is created
Then the system returns a basket ID
And the basket status is active
```

---

## AC-BAS-002: Add Product to Basket

**Priority:** Must  
**Related Requirements:** `FR-BAS-002`, `FR-BAS-006`, `BR-BAS-002`  
**Verification:** API contract test, integration test

```text
Given an active product exists
And the requested quantity is greater than zero
When the customer adds the product to a basket
Then the basket contains the product item
And the item quantity matches the request
And the updated basket is returned
```

---

## AC-BAS-003: Reject Invalid Basket Quantity

**Priority:** Must  
**Related Requirements:** `FR-BAS-007`, `BR-BAS-001`  
**Verification:** API contract test

```text
Given a customer has an active basket
When the customer adds or updates an item with quantity zero or less
Then the system rejects the request
And returns a validation error
And the basket remains unchanged
```

---

## AC-BAS-004: Update Basket Item

**Priority:** Must  
**Related Requirements:** `FR-BAS-003`  
**Verification:** API contract test, integration test

```text
Given a basket contains an item
When the customer updates the item quantity to a valid value
Then the system updates the item quantity
And the basket reflects the new quantity
```

---

## AC-BAS-005: Remove Basket Item

**Priority:** Must  
**Related Requirements:** `FR-BAS-004`  
**Verification:** API contract test, integration test

```text
Given a basket contains an item
When the customer removes the item
Then the item is no longer present in the basket
And the basket remains active
```

---

## 6. Inventory Acceptance Criteria

## AC-INV-001: Reserve Available Stock

**Priority:** Must  
**Related Requirements:** `FR-INV-003`, `BR-INV-003`  
**Verification:** API contract test, repository integration test

```text
Given stock is available for all requested items
When stock reservation is requested
Then Inventory Service creates a reservation
And reserved quantity increases
And available quantity decreases or is otherwise accounted for
And the response includes a reservation ID
```

---

## AC-INV-002: Reject Insufficient Stock

**Priority:** Must  
**Related Requirements:** `FR-INV-004`, `BR-INV-002`  
**Verification:** API contract test, integration test, E2E test

```text
Given requested quantity exceeds available stock
When stock reservation is requested
Then Inventory Service rejects the reservation
And no payment is attempted by checkout
And no confirmed order is created
```

---

## AC-INV-003: Idempotent Stock Reservation

**Priority:** Must  
**Related Requirements:** `NFR-REL-003`, `BR-INV-003`  
**Verification:** Integration test

```text
Given a stock reservation request has already succeeded for an idempotency key
When the same request is repeated with the same idempotency key
Then Inventory Service returns the existing reservation result
And stock is not reserved a second time
```

---

## AC-INV-004: Release Stock Reservation

**Priority:** Must  
**Related Requirements:** `FR-INV-005`, `BR-CHK-009`  
**Verification:** Integration test

```text
Given stock has been reserved for a checkout attempt
And payment subsequently fails
When the checkout failure is handled
Then the stock reservation is released or left to expire according to the documented policy
And stock is not permanently locked
```

---

## 7. Checkout and Order Acceptance Criteria

## AC-CHK-001: Successful Checkout

**Priority:** Must  
**Related Requirements:** `FR-ORD-001`, `FR-ORD-003`, `FR-ORD-004`, `FR-ORD-006`  
**Verification:** E2E test

```text
Given a customer has a basket containing active products
And sufficient stock exists
And payment authorisation succeeds
And shipment creation succeeds
When checkout is submitted
Then stock is reserved
And payment is authorised
And a shipment is created
And an order is created
And the order status is confirmed or equivalent successful state
And an OrderCreated event is published
And the customer receives or is scheduled to receive an order confirmation notification
```

---

## AC-CHK-002: Empty Basket Checkout Fails

**Priority:** Must  
**Related Requirements:** `BR-CHK-001`, `FR-ORD-002`  
**Verification:** API contract test, E2E test

```text
Given a customer has an empty basket
When checkout is submitted
Then checkout is rejected
And no stock reservation is created
And no payment is attempted
And no order is confirmed
```

---

## AC-CHK-003: Inactive Product Checkout Fails

**Priority:** Must  
**Related Requirements:** `BR-CHK-002`, `FR-CAT-010`  
**Verification:** Integration test, E2E test

```text
Given a basket contains a product that is no longer active
When checkout is submitted
Then checkout is rejected
And no payment is attempted
And no confirmed order is created
```

---

## AC-CHK-004: Insufficient Stock Checkout Fails Safely

**Priority:** Must  
**Related Requirements:** `BR-CHK-007`, `FR-INV-004`  
**Verification:** E2E test

```text
Given a basket contains a product with insufficient stock
When checkout is submitted
Then stock reservation fails
And payment is not attempted
And no confirmed order is created
And a safe error is returned to the client
And the failure is logged with a correlation ID
```

---

## AC-CHK-005: Payment Failure Checkout Fails Safely

**Priority:** Must  
**Related Requirements:** `BR-CHK-008`, `FR-PAY-003`  
**Verification:** E2E test

```text
Given stock has been reserved for checkout
And payment authorisation is declined
When checkout continues
Then no confirmed order is created
And the stock reservation is released or allowed to expire according to policy
And the payment failure is recorded
And a safe payment failure response is returned
```

---

## AC-CHK-006: Duplicate Checkout Request

**Priority:** Must  
**Related Requirements:** `FR-ORD-012`, `NFR-REL-001`, `BR-CHK-006`  
**Verification:** Integration test, E2E test

```text
Given a checkout request has already completed with an idempotency key
When the same checkout request is submitted again with the same idempotency key
Then the system returns the original checkout result
And a duplicate confirmed order is not created
And payment is not authorised a second time
And stock is not reserved a second time
```

---

## AC-ORD-001: Retrieve Order

**Priority:** Must  
**Related Requirements:** `FR-ORD-010`  
**Verification:** API contract test, integration test

```text
Given an order exists
When the order is requested by order ID
Then the system returns order details
And the response includes order ID, status, items, totals, and created timestamp
```

---

## AC-ORD-002: Order Item Snapshot

**Priority:** Must  
**Related Requirements:** `FR-ORD-007`, `BR-ORD-002`  
**Verification:** Repository integration test

```text
Given an order is created for a product
When the product name or price changes later
Then the historical order item still shows the product name and price captured at checkout time
```

---

## 8. Payment Acceptance Criteria

## AC-PAY-001: Payment Authorised

**Priority:** Must  
**Related Requirements:** `FR-PAY-001`, `FR-PAY-004`  
**Verification:** API contract test, integration test

```text
Given a valid payment authorisation request
And the simulated payment provider approves the request
When Payment Service authorises payment
Then a payment record is created
And a payment attempt is recorded
And PaymentAuthorised is published
```

---

## AC-PAY-002: Payment Declined

**Priority:** Must  
**Related Requirements:** `FR-PAY-003`, `FR-PAY-005`  
**Verification:** API contract test, integration test

```text
Given a payment request uses a declined payment scenario
When Payment Service attempts authorisation
Then payment is rejected
And the failed attempt is recorded
And PaymentFailed is published
And raw provider-sensitive details are not exposed to the client
```

---

## AC-PAY-003: Idempotent Payment Authorisation

**Priority:** Must  
**Related Requirements:** `NFR-REL-002`, `BR-PAY-003`  
**Verification:** Integration test

```text
Given payment authorisation has already succeeded for an idempotency key
When the same payment request is repeated with the same idempotency key
Then Payment Service returns the original result
And does not create a duplicate provider authorisation
```

---

## AC-PAY-004: Raw Payment Data Not Stored

**Priority:** Must  
**Related Requirements:** `FR-PAY-007`, `NFR-SEC-002`  
**Verification:** Security test, repository inspection test

```text
Given payment authorisation is processed
When payment records are stored
Then raw card data is not stored
And logs do not contain raw card data
And only safe payment references are retained
```

---

## 9. Shipping Acceptance Criteria

## AC-SHP-001: Create Shipment

**Priority:** Must  
**Related Requirements:** `FR-SHP-002`, `FR-SHP-004`  
**Verification:** API contract test, integration test

```text
Given a valid order checkout flow requires shipment creation
When Shipping Service creates a shipment
Then a shipment record is created
And the response includes shipment ID and status
And ShipmentCreated is published
```

---

## AC-SHP-002: Idempotent Shipment Creation

**Priority:** Must  
**Related Requirements:** `FR-SHP-008`, `BR-SHP-002`  
**Verification:** Integration test

```text
Given a shipment has already been created for an idempotency key
When the same shipment request is repeated
Then Shipping Service returns the existing shipment
And does not create a duplicate shipment
```

---

## AC-SHP-003: Shipment Creation Failure

**Priority:** Must  
**Related Requirements:** `FR-SHP-005`, `BR-SHP-005`  
**Verification:** Integration test, E2E test

```text
Given shipment creation fails during checkout
When the checkout flow handles the failure
Then the order outcome follows the documented fulfilment decision
And ShipmentFailed is published
And the failure is observable through logs and metrics
```

---

## 10. Notification Acceptance Criteria

## AC-NOT-001: OrderCreated Triggers Notification

**Priority:** Must  
**Related Requirements:** `FR-NOT-001`, `FR-NOT-003`  
**Verification:** Kafka integration test, E2E test

```text
Given an OrderCreated event is published
When Notification Service consumes the event
Then a notification record is created
And an order confirmation notification is sent or simulated
And the processing result is logged with the event ID and correlation ID
```

---

## AC-NOT-002: Duplicate OrderCreated Does Not Duplicate Notification

**Priority:** Must  
**Related Requirements:** `FR-NOT-006`, `BR-NOT-004`  
**Verification:** Kafka integration test

```text
Given Notification Service has already processed an OrderCreated event
When the same event is consumed again
Then Notification Service does not send a duplicate order confirmation
And the duplicate event is handled idempotently
```

---

## AC-NOT-003: Notification Failure Does Not Roll Back Order

**Priority:** Must  
**Related Requirements:** `FR-NOT-009`, `BR-CHK-010`  
**Verification:** E2E test, resilience test

```text
Given an order has been created
And Notification Service fails to send confirmation
When notification processing fails
Then the order remains created
And the notification failure is recorded
And the failure can be retried or sent to DLQ according to policy
```

---

## 11. Event Acceptance Criteria

## AC-EVT-001: Event Envelope Present

**Priority:** Must  
**Related Requirements:** `FR-EVT-002`, `BR-EVT-003`  
**Verification:** Event contract test

```text
Given a service publishes a Kafka event
When the event is inspected
Then it includes event ID, event type, event version, occurred timestamp, producer, correlation ID, and data payload
```

---

## AC-EVT-002: Unsupported Event Version Handled Safely

**Priority:** Should  
**Related Requirements:** `FR-EVT-006`, `NFR-COMP-004`  
**Verification:** Consumer contract test

```text
Given a consumer receives an unsupported event version
When the event is processed
Then the consumer does not crash
And the event is rejected, ignored, or sent to DLQ according to policy
And the failure is observable
```

---

## AC-EVT-003: Invalid Event Sent to DLQ

**Priority:** Should  
**Related Requirements:** `FR-EVT-004`, `NFR-RES-004`  
**Verification:** Kafka integration test

```text
Given a consumer receives an invalid event payload
When retries are exhausted or the error is non-retryable
Then the event is sent to the appropriate DLQ
And DLQ metadata includes the original topic, event ID, failure reason, and consumer name
```

---

## 12. Data Ownership Acceptance Criteria

## AC-DATA-001: Service-Owned Schema Access

**Priority:** Must  
**Related Requirements:** `NFR-DATA-001`, `NFR-DATA-002`  
**Verification:** Repository test, configuration review

```text
Given a service is configured for database access
When the service starts
Then it uses credentials for its own schema only
And it cannot read or write another service's schema
```

---

## AC-DATA-002: No Cross-Service Database Joins

**Priority:** Must  
**Related Requirements:** `BR-DATA-002`, `BR-DATA-003`  
**Verification:** Code review, static analysis where practical

```text
Given service repository code is reviewed
When SQL queries are inspected
Then no query joins across service-owned schemas
And cross-service data is accessed through APIs, events, snapshots, or projections
```

---

## AC-DATA-003: Monetary Values Use Minor Units

**Priority:** Must  
**Related Requirements:** `NFR-DATA-004`, `BR-CAT-004`  
**Verification:** Schema review, repository test

```text
Given product, order, or payment monetary values are stored
When database schemas are inspected
Then monetary values are represented in minor units
And currency code is stored separately
And floating point types are not used for money
```

---

## 13. Error Handling Acceptance Criteria

## AC-ERR-001: Safe Client Error Response

**Priority:** Must  
**Related Requirements:** `NFR-SEC-004`  
**Verification:** API contract test

```text
Given a service returns an internal failure
When the API Gateway sends an error response to the client
Then the response contains a safe error message
And includes a correlation ID
And does not expose stack traces, SQL errors, secrets, or internal hostnames
```

---

## AC-ERR-002: Validation Error

**Priority:** Must  
**Related Requirements:** `docs/api/error-model.md`  
**Verification:** API contract test

```text
Given a client sends an invalid request
When the request is validated
Then the system returns a validation error
And identifies the invalid field where safe
And does not change business state
```

---

## 14. Observability Acceptance Criteria

## AC-OBS-001: Correlation ID Propagation

**Priority:** Must  
**Related Requirements:** `NFR-OBS-002`, `FR-GW-005`  
**Verification:** Integration test, trace inspection

```text
Given a client request includes or receives a correlation ID
When the request flows through multiple services
Then logs from each participating service include the same correlation ID
And events produced during the flow include the correlation ID
```

---

## AC-OBS-002: Checkout Trace

**Priority:** Should  
**Related Requirements:** `NFR-OBS-006`, `NFR-OBS-007`  
**Verification:** E2E test, observability inspection

```text
Given distributed tracing is enabled
When a checkout is submitted
Then the trace shows API Gateway, Order Service, Basket Service, Inventory Service, Payment Service, Shipping Service, and notification event processing where applicable
```

---

## AC-OBS-003: Health and Readiness

**Priority:** Must  
**Related Requirements:** `NFR-AVL-001`, `NFR-AVL-002`  
**Verification:** Smoke test

```text
Given a service is running
When health and readiness endpoints are requested
Then health reports process status
And readiness reports whether the service can serve traffic
```

---

## 15. Deployment Acceptance Criteria

## AC-DEP-001: Local Environment Starts

**Priority:** Must  
**Related Requirements:** `NFR-DEP-002`, `NFR-DEV-003`  
**Verification:** Local smoke test

```text
Given a developer has Docker available
When the local development environment is started
Then MySQL, Kafka, and the initial bfstore services start successfully
And health checks pass
```

---

## AC-DEP-002: Database Migrations Apply

**Priority:** Must  
**Related Requirements:** `NFR-DATA-007`, `docs/data/migrations.md`  
**Verification:** CI test

```text
Given an empty local or CI MySQL database
When migrations are applied
Then all required schemas and tables are created
And repository tests can run successfully
```

---

## AC-DEP-003: Service Container Builds

**Priority:** Must  
**Related Requirements:** `NFR-DEP-001`  
**Verification:** CI build

```text
Given service source code is available
When the container build runs
Then the service image builds successfully
And the image does not include development-only secrets
```

---

## 16. Security Acceptance Criteria

## AC-SEC-001: Secret Scanning

**Priority:** Should  
**Related Requirements:** `NFR-SEC-007`  
**Verification:** CI security check

```text
Given a pull request is opened
When secret scanning runs
Then committed secrets are detected
And the pipeline fails for confirmed secret findings
```

---

## AC-SEC-002: Container Image Scan

**Priority:** Should  
**Related Requirements:** `NFR-SEC-009`  
**Verification:** CI security check

```text
Given a service container image is built
When container scanning runs
Then vulnerabilities are reported
And critical findings are handled according to the project policy
```

---

## AC-SEC-003: Sensitive Data Not Logged

**Priority:** Must  
**Related Requirements:** `NFR-SEC-001`, `BR-SEC-002`  
**Verification:** Unit/integration test, log inspection

```text
Given a payment or authentication-related operation is processed
When service logs are inspected
Then logs do not contain passwords, tokens, raw card data, or secret values
```

---

## 17. Initial Acceptance Test Set

The first implementation should include automated acceptance coverage for:

```text
AC-CAT-001
AC-CAT-002
AC-CAT-004
AC-BAS-002
AC-BAS-003
AC-INV-001
AC-INV-002
AC-INV-003
AC-CHK-001
AC-CHK-004
AC-CHK-005
AC-CHK-006
AC-PAY-001
AC-PAY-002
AC-PAY-003
AC-SHP-001
AC-SHP-002
AC-NOT-001
AC-NOT-002
AC-EVT-001
AC-DATA-001
AC-DATA-003
AC-ERR-001
AC-OBS-001
AC-DEP-001
AC-DEP-002
```

---

## 18. Definition of Done for the Vertical Slice

The checkout vertical slice is considered accepted when:

```text
active products can be browsed
active products can be added to basket
invalid basket quantities are rejected
checkout succeeds when basket, stock, payment, and shipment are valid
checkout fails safely for insufficient stock
checkout fails safely for payment decline
duplicate checkout requests do not create duplicate orders
OrderCreated is published after successful order creation
Notification Service consumes OrderCreated idempotently
database schemas are service-owned
migrations apply cleanly
basic logs contain correlation IDs
basic tests run in CI
local environment starts successfully
```

---

## 19. Traceability Matrix

| Acceptance Criterion | Related Requirement | Suggested Test |
|---|---|---|
| `AC-CHK-001` | `FR-ORD-001`, `FR-ORD-003`, `FR-ORD-004` | E2E checkout success |
| `AC-CHK-004` | `FR-INV-004`, `BR-CHK-007` | Insufficient stock E2E |
| `AC-CHK-005` | `FR-PAY-003`, `BR-CHK-008` | Payment decline E2E |
| `AC-CHK-006` | `FR-ORD-012`, `NFR-REL-001` | Idempotency integration test |
| `AC-NOT-002` | `FR-NOT-006`, `BR-NOT-004` | Duplicate event consumer test |
| `AC-DATA-001` | `NFR-DATA-001`, `NFR-DATA-002` | Database permission test |
| `AC-EVT-001` | `FR-EVT-002`, `BR-EVT-003` | Event contract test |
| `AC-OBS-001` | `NFR-OBS-002` | Correlation propagation test |

---

## 20. Open Questions

| Question | Status |
|---|---|
| Should shipment failure fail checkout or create pending fulfilment? | To decide |
| Should authentication be required for initial checkout acceptance? | To decide |
| What exact local performance thresholds should be accepted? | To refine |
| Which CI system will produce acceptance test evidence first? | To decide |
| Which security scanners are mandatory in the first CI pipeline? | To decide |
| Should acceptance tests be written in Go, k6, or a separate E2E framework? | To decide |

---

## 21. Related Documents

This document should be read alongside:

```text
docs/requirements/functional-requirements.md
docs/requirements/non-functional-requirements.md
docs/requirements/business-rules.md
docs/requirements/user-journeys.md
docs/architecture/communication-patterns.md
docs/architecture/resilience-patterns.md
docs/events/event-catalog.md
docs/data/data-ownership.md
docs/testing/testing-strategy.md
```

---

## 22. Summary

bfstore acceptance criteria define how the project proves its requirements.

The first acceptance target is the checkout vertical slice:

```text
Browse product
    -> Add to basket
    -> Checkout
    -> Reserve stock
    -> Authorise payment
    -> Create shipment
    -> Create order
    -> Publish OrderCreated
    -> Send notification
```

The criteria are designed to verify not only happy-path behaviour, but also failure handling, idempotency, data ownership, event contracts, observability, and deployment readiness.
