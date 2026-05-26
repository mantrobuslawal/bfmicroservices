# Service Boundaries

## 1. Purpose

This document defines the service boundaries for **bfstore**, ACME Ltd’s fictional online furniture store backend.

The purpose of this document is to make clear:

- which business capability each service owns
- what data each service owns
- which APIs each service exposes
- which events each service publishes or consumes
- what each service must not own
- where synchronous gRPC calls are appropriate
- where asynchronous Kafka events are appropriate
- how to avoid shared database coupling and distributed monolith design

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating the architecture of bfstore.

---

## 2. Architecture Context

bfstore is designed as a microservice backend for an online furniture store.

The system uses:

| Technology | Purpose |
|---|---|
| Go | Primary backend implementation language |
| gRPC | Synchronous service-to-service communication |
| Protobuf | API and event contract definition |
| Kafka | Asynchronous event-driven messaging |
| MySQL | Service-owned relational persistence |
| Docker Compose | Local development environment |
| Kubernetes | Target runtime platform |

The service boundary model follows this principle:

> Each service owns a business capability, its own data, its own API, and its own operational responsibilities.

Services must communicate through documented contracts, not through shared databases or hidden implementation coupling.

---

## 3. Boundary Design Principles

### 3.1 Services Are Organised Around Business Capabilities

Services should represent business responsibilities, not technical layers.

Good service boundaries:

```text
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

Poor service boundaries:

```
database-service
validation-service
utility-service
common-service
crud-service
```

A service should exist because the business concept it owns can change, scale, and be operated independently.

---

### 3.2 Each Service Owns Its Own Data

Each service owns its own MySQL database/schema.

```
catalog-service          -> bfstore_catalog
inventory-service        -> bfstore_inventory
basket-service           -> bfstore_basket
order-service            -> bfstore_order
payment-service          -> bfstore_payment
customer-service         -> bfstore_customer
shipping-service         -> bfstore_shipping
notification-service     -> bfstore_notification
review-service           -> bfstore_review
search-service           -> bfstore_search
recommendation-service   -> bfstore_recommendation
```

Other services must not directly read from or write to another service’s database.

Correct:

```
order-service -> inventory-service gRPC API
order-service <- StockReserved Kafka event
```

Incorrect:

```order-service -> inventory_db.stock_reservations table```

---

### 3.3 APIs and Events Are the Integration Boundary

Services integrate through:

```
gRPC APIs
Kafka events
protobuf contracts
```

They must not integrate through:

```
shared tables
shared ORM models
shared business logic packages
direct file access
implicit database joins
```

The API is the contract.
The database is a private implementation detail.

---

### 3.4 Shared Packages Must Not Contain Business Ownership

Shared packages may contain technical utilities.

Acceptable shared packages:

```
logger
config
grpc middleware
Kafka client wrapper
OpenTelemetry helpers
error handling
health checks
test helpers
```

Avoid shared packages containing business rules such as:

```
order validation
stock reservation rules
payment state transitions
catalogue pricing logic
customer eligibility rules
```

Business rules belong inside the service that owns the capability.

---

### 3.5 Prefer Explicit Ownership Over Convenience

If more than one service appears to own the same concept, the boundary is unclear and must be resolved.

For example:

| Concept                      | Correct Owner          |
| ---------------------------- | ---------------------- |
| Product description          | `catalog-service`      |
| Stock quantity               | `inventory-service`    |
| Basket contents              | `basket-service`       |
| Order lifecycle              | `order-service`        |
| Payment state                | `payment-service`      |
| Shipment state               | `shipping-service`     |
| Notification delivery status | `notification-service` |


Other services may store references, snapshots, or projections, but not become the source of truth.

---

### 3.6 Design for Independent Deployment

A well-bounded service should be independently:

```
built
tested
deployed
scaled
observed
rolled back
operated
```

If a service cannot change without coordinated database or code changes in many other services, the boundary may be too weak.

---

## 4. Service Landscape

bfstore’s target service landscape is:

| Service                  | Primary Capability                                        |
| ------------------------ | --------------------------------------------------------- |
| `api-gateway`            | Client-facing entry point and request routing             |
| `auth-service`           | Authentication, authorisation, sessions, identity         |
| `customer-service`       | Customer profiles, addresses, preferences                 |
| `catalog-service`        | Products, categories, furniture attributes, pricing       |
| `inventory-service`      | Stock levels, warehouses, stock reservations              |
| `basket-service`         | Customer baskets and basket items                         |
| `order-service`          | Order creation, order lifecycle, order history            |
| `payment-service`        | Payment authorisation, capture, refunds, payment attempts |
| `shipping-service`       | Delivery options, shipments, fulfilment status            |
| `notification-service`   | Customer notifications and delivery status                |
| `review-service`         | Product reviews, ratings, moderation                      |
| `search-service`         | Product search index and query model                      |
| `recommendation-service` | Product recommendations and recommendation signals        |

The initial implementation should focus on the checkout vertical slice:

```
api-gateway
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

---

##  5. Service Boundary Summary

| Service                  | Owns                                                      | Does Not Own                                                           |
| ------------------------ | --------------------------------------------------------- | ---------------------------------------------------------------------- |
| `api-gateway`            | External API routing, request shaping, edge concerns      | Business data, order lifecycle, product catalogue, payment state       |
| `auth-service`           | Identities, credentials, tokens, sessions, roles          | Customer profile, addresses, baskets, orders                           |
| `customer-service`       | Customer profile, addresses, preferences                  | Credentials, sessions, orders, payments                                |
| `catalog-service`        | Products, categories, variants, product metadata, pricing | Stock levels, baskets, orders, reviews                                 |
| `inventory-service`      | Stock levels, warehouses, reservations, stock adjustments | Product descriptions, baskets, orders, payments                        |
| `basket-service`         | Baskets, basket items, basket lifecycle                   | Stock reservation, payment, order lifecycle                            |
| `order-service`          | Orders, order items, order status, order orchestration    | Product catalogue, payment processing internals, stock source of truth |
| `payment-service`        | Payments, payment attempts, refunds, payment state        | Order lifecycle, customer profile, stock, shipment fulfilment          |
| `shipping-service`       | Delivery options, shipments, tracking status              | Order creation, payment, stock ownership                               |
| `notification-service`   | Notification requests, delivery attempts, delivery status | Order state, customer profile source of truth                          |
| `review-service`         | Reviews, ratings, moderation decisions                    | Product catalogue source of truth, order lifecycle                     |
| `search-service`         | Search index, search query model, search projections      | Product catalogue source of truth                                      |
| `recommendation-service` | Recommendation rules, signals, calculated recommendations | Product catalogue, orders, payments, reviews as source of truth        |

---

## 6. Detailed Service Boundaries

### 6.1 API Gateway

### Purpose

The ```api-gateway``` is the public entry point for client applications.

It hides internal service topology from external clients and provides a stable client-facing API surface.

### Owns

```
client-facing API routes
request validation at the edge
authentication enforcement at the edge
request correlation IDs
response mapping
error mapping
rate limiting design
client protocol adaptation
```

### Does Not Own

```
business entities
product catalogue
stock levels
basket state
order lifecycle
payment state
customer profile source of truth
shipping state
notification state
```

### Inbound Interfaces

External client requests.

The external API may be REST, GraphQL, gRPC-Web, or another client-friendly interface. The final protocol choice should be documented in an ADR.

### Outbound Calls

The API Gateway may call:

```
auth-service
customer-service
catalog-service
basket-service
order-service
search-service
recommendation-service
review-service
```

Boundary Rules
The gateway must not contain core business rules.
The gateway must not directly access service databases.
The gateway may perform request shape validation, but domain validation belongs in the owning service.
The gateway should propagate correlation IDs across all downstream calls.
The gateway should map internal errors to safe client-facing errors.

Example Responsibilities

Correct:

```
Validate request format.
Check authentication token.
Call order-service CreateOrder.
Return standardised client response.
```

Incorrect:

```
Calculate final order state.
Reserve stock directly.
Write to order database.
Authorise payment directly.
```

---


### 6.2 Auth Service

### Purpose

The auth-service owns authentication and authorisation identity concerns.

It answers:

```
Who is the user?
Can this user be issued a token?
What roles or permissions does this identity have?
```

### Owns

```
identity
credentials
password hashes
sessions
access tokens
refresh tokens
roles
permissions
login attempts
account lock status
```

Does Not Own

```
customer profile
delivery addresses
orders
payments
basket contents
reviews
Inbound APIs
```

### Inbound APIs

Potential gRPC APIs:

```
RegisterIdentity
Authenticate
ValidateToken
RefreshToken
RevokeSession
GetIdentity
```

### Events Published

Potential events:

```
IdentityCreated
CustomerSignedIn
FailedLoginAttempted
SessionRevoked
```

### Events Consumed

Potential events:

```
CustomerDeleted
CustomerDisabled
```

### Data Owned

```
identities
credentials
sessions
roles
permissions
login_attempts
```

Boundary Rules
Auth Service owns credentials. Customer Service owns customer profile.
Auth Service must not store customer addresses or order history.
Passwords must be hashed and never stored in plain text.
Authentication tokens must not be logged.
Customer-facing account data should be split carefully between identity and profile concerns.

---

### 6.3 Customer Service

### Purpose

The customer-service owns customer profile information, delivery addresses, and preferences.

It answers:

```
Who is this customer from a business profile perspective?
Where can orders be delivered?
What preferences has the customer set?
```

### Owns

```
customer profile
customer addresses
default delivery address
customer preferences
contact preferences
profile status
```

### Does Not Own

```
authentication credentials
tokens
basket contents
orders
payments
shipments
reviews
```

### Inbound APIs

Potential gRPC APIs:

```
CreateCustomerProfile
GetCustomerProfile
UpdateCustomerProfile
AddAddress
UpdateAddress
DeleteAddress
ListAddresses
SetDefaultAddress
```

### Events Published

```
CustomerProfileCreated
CustomerProfileUpdated
CustomerAddressAdded
CustomerAddressUpdated
CustomerAddressDeleted
```

### Events Consumed

```
IdentityCreated
CustomerRegistered
```

### Data Owned

```
customers
customer_profiles
addresses
customer_preferences
```

### Boundary Rules

Customer Service owns addresses, but orders and shipments may store address snapshots for historical accuracy.
Customer Service must not own credentials or sessions.
Customer PII must be protected and should not be logged unnecessarily.
Other services should request customer profile data via API or consume events where appropriate.

---

### 6.4 Catalog Service

### Purpose

The catalog-service owns the product catalogue.

It answers:

```
What products does ACME sell?
What are their details?
Which products are active and purchasable?
How are products categorised?
```

### Owns

```
products
product variants
categories
product images metadata
product attributes
materials
colours
dimensions
product pricing
product active/inactive state
```

### Does Not Own

```
stock quantity
stock reservations
basket state
order history
customer reviews
search index projections
recommendation results
```

### Inbound APIs

Potential gRPC APIs:

```
GetProduct
ListProducts
CreateProduct
UpdateProduct
ActivateProduct
DeactivateProduct
ListCategories
GetCategory
```

### Events Published

```
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
ProductDiscontinued
CategoryCreated
CategoryUpdated
```

### Events Consumed

Usually none for the initial version.

Potential later events:

```
ReviewApproved
InventoryAdjusted
```

These should only update local display summaries if a deliberate projection is introduced.

### Data Owned

```
products
product_variants
categories
product_images
product_attributes
product_prices
```

### Boundary Rules
Catalog Service is the source of truth for product details.
Inventory Service is the source of truth for stock.
Search Service may index product data but must not become the product source of truth.
Recommendation Service may use product data but must not own product facts.
Order Service should store product snapshots for historical order accuracy, not call Catalog Service for historical product values.

---

###  6.5 Inventory Service

### Purpose

The inventory-service owns stock levels, warehouse stock, and stock reservations.

It answers:

```
Is stock available?
Can stock be reserved?
Has reserved stock been committed, released, or expired?
```

### Owns

```
inventory items
stock levels
warehouses
stock reservations
reservation expiry
stock adjustments
stock commitment
```

### Does Not Own

```
product descriptions
product pricing
basket contents
order lifecycle
payment state
shipment tracking
```

### Inbound APIs

Potential gRPC APIs:

```
CheckAvailability
ReserveStock
ReleaseReservation
CommitReservation
GetStockLevel
AdjustStock
```

### Events Published

```
InventoryAdjusted
StockReserved
StockReservationFailed
StockReservationReleased
StockReservationExpired
StockCommitted
```

### Events Consumed

Potential events:

```
ProductCreated
ProductDiscontinued
OrderCancelled
PaymentFailed
```

### Data Owned

```
inventory_items
warehouses
stock_reservations
stock_adjustments
```

### Boundary Rules

Inventory Service is the only service that can change stock levels.
Basket Service may check availability but must not reserve stock.
Order Service requests reservations; it does not write stock records.
Stock reservations should be idempotent where possible.
Reservation expiry should be explicit and observable.

---

### 6.6 Basket Service

### Purpose

The basket-service owns customer baskets and basket items.

It answers:

```
What is the customer currently intending to buy?
What items and quantities are in the basket?
```

### Owns

```
baskets
basket items
basket lifecycle
guest basket reference
basket expiry
```

### Does Not Own

```
stock reservation
order creation
payment state
product source of truth
customer profile source of truth
```

### Inbound APIs

Potential gRPC APIs:

```
CreateBasket
GetBasket
AddItem
UpdateItem
RemoveItem
ClearBasket
MarkBasketCheckedOut
```

### Outbound Calls

May call:

```
catalog-service
inventory-service
Events Published
BasketCreated
BasketItemAdded
BasketItemUpdated
BasketItemRemoved
BasketCheckedOut
BasketAbandoned
```

### Events Consumed

Potential events:

```
OrderCreated
ProductDeactivated
```

### Data Owned

```
baskets
basket_items
```

### Boundary Rules
Basket Service does not reserve stock.
Basket Service should validate that products are active before adding them.
Basket Service may show availability, but Inventory Service owns availability truth.
Basket price data should be treated carefully. Final price should be confirmed during checkout or stored as an explicit snapshot.
Basket Service should not create orders.

---









 




