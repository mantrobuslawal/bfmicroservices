# User Journeys

## 1. Document Purpose

This document describes the key user journeys for **bfstore**, ACME Ltd’s fictional online furniture store backend.

User journeys explain how different users interact with the system to achieve their goals. They help guide:

- functional requirements
- service boundaries
- gRPC API design
- Kafka event design
- data ownership
- acceptance criteria
- end-to-end testing
- observability requirements
- operational runbooks

This document focuses on the business flows that the backend must support.

---

## 2. Product Context

bfstore is the backend platform for ACME Ltd’s online furniture store.

Customers should be able to browse furniture, search the catalogue, manage a basket, place orders, reserve stock, authorise payments, arrange delivery, receive notifications, submit reviews, and receive product recommendations.

The backend is implemented as a microservice architecture using:

- gRPC for synchronous service-to-service communication
- Kafka for asynchronous business events
- Protobuf for API and event contracts
- MySQL for service-owned databases

---

## 3. Journey Design Principles

The user journeys should follow these principles:

1. **Start with user goals**  
   Journeys should describe what the user is trying to achieve, not which service is called first.

2. **Map business flow to service behaviour**  
   Each journey should identify the services involved and the expected backend responsibilities.

3. **Separate synchronous and asynchronous work**  
   Immediate decisions use gRPC. Downstream reactions use Kafka events.

4. **Include success and failure paths**  
   Real systems must handle unavailable stock, failed payments, duplicate events, and service failures.

5. **Link journeys to tests**  
   Each journey should be testable through end-to-end, integration, contract, performance, or resilience tests.

6. **Make observability explicit**  
   Important journeys should define the logs, metrics, traces, and alerts needed to diagnose them.

---

## 4. Primary Actors

| Actor | Description |
|---|---|
| Customer | A person browsing and purchasing furniture from ACME Ltd |
| Registered Customer | A signed-in customer with saved profile, addresses, orders, and preferences |
| Guest Customer | A customer browsing before signing in or registering |
| ACME Support User | An internal user who may need to inspect orders, payments, shipments, or customer issues |
| ACME Operations Team | Team responsible for monitoring and operating the platform |
| ACME Engineering Team | Team responsible for building, testing, deploying, and maintaining services |
| External Payment Provider | Simulated or future third-party payment provider |
| External Notification Provider | Simulated or future email/SMS provider |
| External Shipping Provider | Simulated or future carrier/fulfilment provider |

---

## 5. Core User Journeys Summary

| ID | Journey | Primary Actor | Priority | Initial Version |
|---|---|---|---|---|
| `UJ-001` | Browse furniture catalogue | Customer | Must | Yes |
| `UJ-002` | View product details | Customer | Must | Yes |
| `UJ-003` | Search and filter products | Customer | Should | Later |
| `UJ-004` | Register account | Customer | Should | Later |
| `UJ-005` | Sign in | Registered Customer | Should | Later |
| `UJ-006` | Manage customer profile and addresses | Registered Customer | Should | Later |
| `UJ-007` | Add product to basket | Customer | Must | Yes |
| `UJ-008` | Update basket | Customer | Must | Yes |
| `UJ-009` | Checkout and create order | Customer | Must | Yes |
| `UJ-010` | Reserve stock during checkout | System | Must | Yes |
| `UJ-011` | Authorise payment during checkout | System | Must | Yes |
| `UJ-012` | Create shipment | System | Must | Yes |
| `UJ-013` | Send order confirmation notification | System | Must | Yes |
| `UJ-014` | View order history | Registered Customer | Should | Later |
| `UJ-015` | Track shipment | Registered Customer | Should | Later |
| `UJ-016` | Submit product review | Registered Customer | Could | Later |
| `UJ-017` | Receive product recommendations | Customer | Could | Later |
| `UJ-018` | Handle payment failure | Customer/System | Must | Yes |
| `UJ-019` | Handle insufficient stock | Customer/System | Must | Yes |
| `UJ-020` | Cancel order | Registered Customer | Should | Later |

---

# 6. Initial Vertical Slice

The first implementation should focus on a complete checkout journey.

```text
Browse product
    -> View product details
    -> Add to basket
    -> Checkout
    -> Reserve stock
    -> Authorise payment
    -> Create order
    -> Create shipment
    -> Publish OrderCreated event
    -> Send notification
```

Initial services involved:

```text
api-gateway
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

The purpose of this vertical slice is to prove the core architecture before expanding into search, reviews, recommendations, advanced account features, and production platform maturity.

---

## 7. Journey Details

### UJ-001: Browse Furniture Catalogue

**Goal**

A customer wants to browse available furniture products by category.

**Primary Actor**

Customer

**Priority**


Must

**Services Involved**

| Service             | Responsibility                                                |
| ------------------- | ------------------------------------------------------------- |
| `api-gateway`       | Receives client request and routes to catalogue               |
| `catalog-service`   | Returns product catalogue data                                |
| `inventory-service` | May provide stock availability summary                        |
| `search-service`    | Later phase: supports search/filter-backed catalogue browsing |


**Preconditions**

- Catalogue products exist.
- Products have active/inactive status.
- Product categories exist.
- API Gateway is available.
- Catalogue Service is available.

**Main Success Flow**

1. Customer opens the furniture catalogue.
2. API Gateway receives a request to list products.
3. API Gateway calls Catalogue Service using gRPC.
4. Catalogue Service retrieves active products.
5. Catalogue Service returns product summaries.
6. API Gateway returns the product list to the client.

**Example Product Summary**

```text
product_id
name
category
thumbnail_url
price
currency
material
colour
dimensions_summary
availability_summary
```

**Alternative Flows**

A1: Category Filter Applied

1. Customer selects a category, such as sofas or dining tables.
2. API Gateway sends category filter to Catalogue Service.
3. Catalogue Service returns active products in that category.

A2: No Products Found

1. Customer selects a category with no active products.
2. Catalogue Service returns an empty result set.
3. API Gateway returns an empty catalogue response.

   
**Failure Flows**

F1: Catalogue Service Unavailable
API Gateway calls Catalogue Service.
Catalogue Service is unavailable.
API Gateway returns an appropriate service unavailable error.
Error is logged with correlation ID.
Events

No event is required for simple browsing in the initial version.

Later phases may emit analytics events, for example:

```text
ProductListViewed
CategoryViewed
```

**Observability Requirements**

- Log catalogue list requests.
- Record request latency.
- Record error count.
- Trace API Gateway to Catalogue Service call.
- Track empty result responses.

**Acceptance Criteria**

```text
Given active products exist
When a customer browses the catalogue
Then the system returns a list of active product summaries

Given a category has no active products
When a customer browses that category
Then the system returns an empty product list without error
```

---

### UJ-002: View Product Details

**Goal**

A customer wants to view detailed information about a furniture product.

**Primary Actor**

Customer

**Priority**

Must

**Services Involved**

| Service             | Responsibility                          |
| ------------------- | --------------------------------------- |
| `api-gateway`       | Receives request and routes to services |
| `catalog-service`   | Returns product details                 |
| `inventory-service` | Returns stock availability              |


**Preconditions**

- Product exists.
- Product is active.
- Catalogue Service is available.

**Main Success Flow**

1. Customer selects a product.
2. API Gateway receives product detail request.
3. API Gateway calls Catalogue Service to get product details.
4. API Gateway may call Inventory Service to get availability.
5. API Gateway returns product details and availability to the client.

**Product Details Include**

```text
product_id
name
description
category
price
currency
material
colour
dimensions
weight
images
care_instructions
availability
```

**Alternative Flows**

A1: Product Exists But Is Out of Stock

1. Product exists and is active.
2. Inventory Service returns unavailable or low stock.
3. API Gateway returns product details with stock availability status.

A2: Product Is Inactive

1. Customer requests an inactive product.
2. Catalogue Service rejects or hides the product.
3. API Gateway returns not found or unavailable.

**Failure Flows**

F1: Product Not Found

1. Customer requests an unknown product ID.
2. Catalogue Service returns not found.
3. API Gateway returns not found to the client.


F2: Inventory Service Unavailable

1. Catalogue details are available.
2. Inventory Service cannot be reached.
3. API Gateway may return product details with availability marked as unknown.
4. Failure is logged and traced.

Events

No Kafka event is required in the initial version.

Later phases may emit:
```Later phases may emit:```

Observability Requirements

- Trace API Gateway to Catalogue Service.
- Trace optional Inventory Service call.
- Record product not found count.
- Record inventory availability lookup failures.

Acceptance Criteria

```text
Given an active product exists
When a customer views the product
Then the system returns product details

Given a product is out of stock
When a customer views the product
Then the system shows the product as unavailable or out of stock
```

---

### UJ-003: Search and Filter Products

**Goal**

A customer wants to search for furniture and filter results by category, price, material, colour, or dimensions.

**Primary Actor**

Customer

**Priority**

Should

**Initial Version**

No. Later phase.

**Services Involved**

| Service                   | Responsibility                               |
| ------------------------- | -------------------------------------------- |
| `api-gateway`             | Receives search request                      |
| `search-service`          | Executes product search and filtering        |
| `catalog-service`         | Source of product truth                      |
| `inventory-service`       | Optional availability data                   |
| `catalog-service` / Kafka | Publishes product change events for indexing |

Preconditions

- Product catalogue exists.
- Search index exists or can be generated.
- Product updates are reflected in search index.
  
Main Success Flow

1. Customer enters a search query.
2. Customer applies optional filters.
3. API Gateway calls Search Service.
4. Search Service queries its search index.
5. Search Service returns matching product summaries.
6. API Gateway returns search results to the client.
   
Example Filters

```text
category
price_min
price_max
material
colour
room
width_cm
height_cm
depth_cm
availability
```

Events

Search Service may consume:

```text
ProductCreated
ProductUpdated
ProductDeleted
InventoryAdjusted
```

Search Service may publish:
```SearchIndexUpdated```

**Failure Flows**

F1: Search Index Stale

1. Product was updated in Catalogue Service.
2. Search index has not yet processed the update.
3. Search results may temporarily show stale data.
4. Eventual consistency is accepted and documented.
   
F2: Search Service Unavailable

1. API Gateway calls Search Service.
2. Search Service is unavailable.
3. API Gateway returns search unavailable response.
   
Acceptance Criteria
```text
Given products exist in the search index
When a customer searches by keyword
Then matching active products are returned

Given a product is inactive
When a customer searches
Then the inactive product is not returned
```

---

### UJ-004: Register Account

**Goal**

A customer wants to create an ACME store account.

**Primary Actor**

Customer

**Priority**

Should

**Initial Version**

Optional for first vertical slice. Can be deferred if using test customers.

Services Involved

| Service                | Responsibility                             |
| ---------------------- | ------------------------------------------ |
| `api-gateway`          | Receives registration request              |
| `auth-service`         | Creates authentication identity            |
| `customer-service`     | Creates customer profile                   |
| `notification-service` | Sends welcome or verification notification |

Preconditions

- Customer does not already have an account with the same email address.
- Auth Service is available.
- Customer Service is available.

Main Success Flow

1. Customer submits registration details.
2. API Gateway validates basic request shape.
3. API Gateway calls Auth Service to create identity.
4. Auth Service hashes password and creates authentication record.
5. Auth Service calls or emits event for Customer Service to create profile.
6. Customer profile is created.
7. Registration success is returned.
8. Notification may be sent asynchronously.

Events

Possible events:
```text
CustomerRegistered
CustomerProfileCreated
NotificationRequested
```

Failure Flows

F1: Email Already Registered

1. Customer submits an email already in use.
2. Auth Service rejects the request.
3. API Gateway returns duplicate account error.

F2: Customer Profile Creation Fails

1. Auth identity is created.
2. Customer profile creation fails.
3. System must either compensate, retry, or mark registration incomplete.
4. Failure is logged and visible to operations.

Acceptance Criteria

```text
Given a new customer email
When the customer registers
Then an authentication identity and customer profile are created

Given an email is already registered
When the customer attempts to register
Then the system rejects the registration
```

---

UJ-005: Sign In
Goal

A registered customer wants to sign in to access protected actions.

Primary Actor

Registered Customer

Priority

Should

Initial Version

Optional for first vertical slice if using test identity.

Services Involved

| Service        | Responsibility                                 |
| -------------- | ---------------------------------------------- |
| `api-gateway`  | Receives login request                         |
| `auth-service` | Validates credentials and issues token/session |

Preconditions
Customer account exists.
Customer credentials are valid.
Auth Service is available.
Main Success Flow
Customer submits email and password.
API Gateway sends credentials to Auth Service.
Auth Service verifies password.
Auth Service issues access token/session.
API Gateway returns authentication response.
Failure Flows
F1: Invalid Credentials
Customer submits incorrect credentials.
Auth Service rejects the request.
API Gateway returns authentication failure.
F2: Account Locked or Disabled
Customer submits credentials for locked account.
Auth Service rejects the request.
Audit event is recorded.
Events

Potential events:

CustomerSignedIn
FailedLoginAttempted

Acceptance Criteria

Given valid credentials
When a customer signs in
Then the system returns an authentication token or session
Given invalid credentials
When a customer signs in
Then the system rejects the request

---

UJ-006: Manage Customer Profile and Addresses
Goal

A registered customer wants to manage their profile and delivery addresses.

Primary Actor

Registered Customer

Priority

Should

Initial Version

Later phase.

Services Involved

| Service            | Responsibility                                    |
| ------------------ | ------------------------------------------------- |
| `api-gateway`      | Receives profile/address requests                 |
| `customer-service` | Stores and manages customer profile and addresses |
| `auth-service`     | Provides identity and authorisation context       |


Main Success Flow

Customer signs in.
Customer opens profile settings.
Customer adds or updates delivery address.
API Gateway calls Customer Service.
Customer Service validates and stores the address.
Updated profile/address is returned.

Business Rules

A customer may have multiple delivery addresses.
One address may be marked as default.
Checkout requires a valid delivery address.
Sensitive customer data must not be logged.
Events

Possible events:

CustomerProfileUpdated
CustomerAddressAdded
CustomerAddressUpdated

Acceptance Criteria

Given a signed-in customer
When the customer adds a valid delivery address
Then the address is saved to the customer profile

---

UJ-007: Add Product to Basket
Goal

A customer wants to add a furniture product to their basket.

Primary Actor

Customer

Priority

Must

Services Involved

| Service             | Responsibility                                |
| ------------------- | --------------------------------------------- |
| `api-gateway`       | Receives add-to-basket request                |
| `basket-service`    | Manages basket state                          |
| `catalog-service`   | Validates product existence and active status |
| `inventory-service` | Optionally validates availability             |

Preconditions

Product exists.
Product is active.
Requested quantity is valid.
Basket Service is available.

Main Success Flow

Customer selects product and quantity.
API Gateway receives add-to-basket request.
API Gateway calls Basket Service.
Basket Service validates product through Catalogue Service.
Basket Service may check availability through Inventory Service.
Basket Service adds item to basket.
Updated basket is returned.

Business Rules

Quantity must be greater than zero.
Inactive products cannot be added.
Basket should store product ID, quantity, and price snapshot where needed.
Basket does not reserve stock. Stock is reserved during checkout.

Alternative Flows

A1: Product Already in Basket

Customer adds a product already in basket.
Basket Service increases quantity or updates existing basket item.
Updated basket is returned.

A2: Product Available But Low Stock

Customer requests a quantity higher than currently available.
Basket Service rejects or adjusts depending on business rule.

Failure Flows

F1: Product Not Found

Basket Service validates product.
Catalogue Service returns not found.
Basket Service rejects add-to-basket request.

F2: Product Inactive

Product exists but is inactive.
Basket Service rejects request.

Events

Potential event:

```BasketItemAdded```

This may be deferred unless needed for analytics or recommendations.

Acceptance Criteria

```text
Given an active product exists
When a customer adds the product to their basket
Then the basket contains the product and requested quantity

Given a product is inactive
When a customer adds the product to their basket
Then the system rejects the request
```

---

UJ-008: Update Basket
Goal

A customer wants to update the quantity of an item or remove an item from their basket.

Primary Actor

Customer

Priority

Must

Services Involved

| Service             | Responsibility                              |
| ------------------- | ------------------------------------------- |
| `api-gateway`       | Receives update request                     |
| `basket-service`    | Updates basket item quantities and removals |
| `catalog-service`   | Optional product validation                 |
| `inventory-service` | Optional availability check                 |

Preconditions

Basket exists.
Basket item exists.
Requested update is valid.

Main Success Flow

Customer changes item quantity.
API Gateway sends update to Basket Service.
Basket Service validates quantity.
Basket Service updates basket.
Updated basket is returned.

Alternative Flows

A1: Quantity Set to Zero

Customer sets quantity to zero.
Basket Service removes item from basket.
Updated basket is returned.

A2: Basket Becomes Empty

Customer removes final basket item.
Basket Service returns empty basket.

Business Rules

Quantity cannot be negative.
Quantity zero means remove item.
Basket updates do not reserve stock.
Basket prices may be recalculated or use stored price snapshots depending on design decision.

Events

Possible events:
```text
BasketItemUpdated
BasketItemRemoved
BasketEmptied
```

Acceptance Criteria

```text
Given a basket contains an item
When the customer updates the item quantity
Then the basket reflects the new quantity

Given a basket contains an item
When the customer removes the item
Then the item is no longer present in the basket
```

---

UJ-009: Checkout and Create Order

Goal

A customer wants to checkout and create an order.

Primary Actor

Customer

Priority

Must

Services Involved

| Service                | Responsibility                        |
| ---------------------- | ------------------------------------- |
| `api-gateway`          | Receives checkout request             |
| `basket-service`       | Provides basket contents              |
| `order-service`        | Coordinates order creation            |
| `inventory-service`    | Reserves stock                        |
| `payment-service`      | Authorises payment                    |
| `shipping-service`     | Creates shipment                      |
| `notification-service` | Sends confirmation asynchronously     |
| `Kafka`                | Carries order and notification events |

Preconditions

Customer has a non-empty basket.
Basket items are valid.
Delivery address exists or is supplied.
Stock is available.
Payment details are accepted by Payment Service or simulated provider.
Required services are available.

Main Success Flow

Customer submits checkout request.
API Gateway calls Order Service.
Order Service retrieves or validates basket.
Order Service requests stock reservation from Inventory Service.
Inventory Service reserves stock.
Order Service requests payment authorisation from Payment Service.
Payment Service authorises payment.
Order Service creates order.
Order Service requests shipment creation from Shipping Service.
Shipping Service creates shipment.
Order Service publishes OrderCreated.
Order Service returns order confirmation.
Notification Service consumes event and sends confirmation.

Synchronous Calls

```text
api-gateway -> order-service
order-service -> basket-service
order-service -> inventory-service
order-service -> payment-service
order-service -> shipping-service
```

Events Published

```text 
StockReserved
PaymentAuthorised
ShipmentCreated
OrderCreated
NotificationRequested
```

Business Rules

An order cannot be confirmed unless stock is reserved.
An order cannot be confirmed unless payment is authorised.
A basket must not be empty at checkout.
Checkout should be idempotent where possible.
Duplicate checkout submissions must not create duplicate confirmed orders.
Sensitive payment data must not be logged.

Failure Flows

F1: Basket Empty

Customer attempts checkout.
Order Service detects empty basket.
Checkout is rejected.

F2: Stock Reservation Fails

Order Service requests stock reservation.
Inventory Service returns insufficient stock.
Order Service rejects checkout.
Payment is not attempted.

F3: Payment Authorisation Fails

Stock reservation succeeds.
Payment Service rejects payment.
Order Service marks checkout failed.
Stock reservation is released or allowed to expire.
Customer receives payment failure response.

F4: Shipment Creation Fails

Stock reservation and payment authorisation succeed.
Shipping Service fails to create shipment.
Order Service must handle failure according to design decision:
fail order creation and compensate
create order in pending fulfilment state
retry shipment creation asynchronously

F5: Duplicate Checkout Request

Customer submits same checkout request twice.
Order Service detects idempotency key or duplicate request.
Existing order result is returned rather than creating a duplicate order.

Observability Requirements

End-to-end trace from API Gateway to all services.
Correlation ID across gRPC calls and Kafka events.
Metrics for checkout success and failure.
Metrics for stock reservation failures.
Metrics for payment failures.
Metrics for checkout latency.
Logs must include order ID once created.
Logs must not include sensitive payment data.

Acceptance Criteria

```text
Given a customer has a valid basket and stock is available
When the customer checks out with valid payment details
Then stock is reserved, payment is authorised, an order is created, a shipment is created, and an OrderCreated event is published

Given stock is unavailable
When the customer checks out
Then the order is not created and payment is not authorised

Given payment authorisation fails
When the customer checks out
Then the order is not confirmed and the customer receives a payment failure response
```

---







