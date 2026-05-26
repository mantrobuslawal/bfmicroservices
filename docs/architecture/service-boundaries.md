# Service Boundaries

## 1. Purpose

This document defines the service boundaries for **bfstore**, ACME Ltd’s fictional online furniture store backend.

The purpose of this document is to make clear:

- which business capability each service owns
- what data each service owns
- which APIs each service exposes
- which events each service publishes or consumes
- what each service must not own
- how services communicate without sharing databases
- where synchronous gRPC calls are appropriate
- where asynchronous Kafka events are appropriate
- how the system avoids becoming a distributed monolith

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating the architecture of bfstore.

---

## 2. Architecture Context

bfstore is designed as a cloud-native microservice backend for an online furniture store.

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
| OpenTelemetry | Logs, metrics, traces, and correlation |
| Buf | Protobuf linting, generation, and breaking-change checks |

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

```text
database-service
validation-service
utility-service
common-service
crud-service
```

A service should exist because the business capability it owns can change, scale, and be operated independently.

---

### 3.2 Each Service Owns Its Own Data

Each service owns its own MySQL database or schema.

```text
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

```text
order-service -> inventory-service ReserveStock gRPC API
inventory-service -> StockReserved Kafka event
```

Incorrect:

```text
order-service -> inventory_db.stock_reservations table
```

The API is the contract. The database is a private implementation detail.

---

### 3.3 APIs and Events Are the Integration Boundary

Services integrate through:

```text
gRPC APIs
Kafka events
protobuf contracts
```

Services must not integrate through:

```text
shared tables
shared ORM models
shared business logic packages
direct file access
implicit database joins
```

This keeps service ownership clear and prevents tight runtime and data coupling.

---

### 3.4 Shared Packages Must Not Contain Business Ownership

Shared packages may contain technical utilities.

Acceptable shared packages:

```text
logger
config
gRPC middleware
Kafka client wrapper
OpenTelemetry helpers
error handling
health checks
test helpers
```

Shared packages should not contain business rules such as:

```text
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

| Concept | Correct Owner |
|---|---|
| Product description | `catalog-service` |
| Product price | `catalog-service` |
| Stock quantity | `inventory-service` |
| Stock reservation | `inventory-service` |
| Basket contents | `basket-service` |
| Order lifecycle | `order-service` |
| Payment state | `payment-service` |
| Shipment state | `shipping-service` |
| Notification delivery status | `notification-service` |
| Product review | `review-service` |
| Search index | `search-service` |
| Recommendation output | `recommendation-service` |

Other services may store references, snapshots, or projections, but they must not become the source of truth.

---

### 3.6 Design for Independent Deployment

A well-bounded service should be independently:

```text
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

### 3.7 Keep the Critical Path Small

The checkout path is business-critical. It should involve only the services required to complete checkout safely.

Critical path services:

```text
api-gateway
basket-service
order-service
inventory-service
payment-service
shipping-service
```

Asynchronous downstream services:

```text
notification-service
search-service
recommendation-service
analytics or reporting consumers
```

Search, recommendations, analytics, and notifications should not block order creation unless there is a deliberate business reason.

---

## 4. Service Landscape

bfstore’s target service landscape is:

| Service | Primary Capability |
|---|---|
| `api-gateway` | Client-facing entry point and request routing |
| `auth-service` | Authentication, authorisation, sessions, identity |
| `customer-service` | Customer profiles, addresses, preferences |
| `catalog-service` | Products, categories, furniture attributes, pricing |
| `inventory-service` | Stock levels, warehouses, stock reservations |
| `basket-service` | Customer baskets and basket items |
| `order-service` | Order creation, order lifecycle, order history |
| `payment-service` | Payment authorisation, capture, refunds, payment attempts |
| `shipping-service` | Delivery options, shipments, fulfilment status |
| `notification-service` | Customer notifications and delivery status |
| `review-service` | Product reviews, ratings, moderation |
| `search-service` | Product search index and query model |
| `recommendation-service` | Product recommendations and recommendation signals |

The initial implementation should focus on the checkout vertical slice:

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

---

## 5. Service Boundary Summary

| Service | Owns | Does Not Own |
|---|---|---|
| `api-gateway` | External API routing, request shaping, edge concerns | Business data, order lifecycle, product catalogue, payment state |
| `auth-service` | Identities, credentials, tokens, sessions, roles | Customer profile, addresses, baskets, orders |
| `customer-service` | Customer profile, addresses, preferences | Credentials, sessions, orders, payments |
| `catalog-service` | Products, categories, variants, product metadata, pricing | Stock levels, baskets, orders, reviews |
| `inventory-service` | Stock levels, warehouses, reservations, stock adjustments | Product descriptions, baskets, orders, payments |
| `basket-service` | Baskets, basket items, basket lifecycle | Stock reservation, payment, order lifecycle |
| `order-service` | Orders, order items, order status, order orchestration | Product catalogue, payment internals, stock source of truth |
| `payment-service` | Payments, payment attempts, refunds, payment state | Order lifecycle, customer profile, stock, shipment fulfilment |
| `shipping-service` | Delivery options, shipments, tracking status | Order creation, payment, stock ownership |
| `notification-service` | Notification requests, delivery attempts, delivery status | Order state, customer profile source of truth |
| `review-service` | Reviews, ratings, moderation decisions | Product catalogue source of truth, order lifecycle |
| `search-service` | Search index, search query model, search projections | Product catalogue source of truth |
| `recommendation-service` | Recommendation rules, signals, calculated recommendations | Product catalogue, orders, payments, reviews as source of truth |

---

## 6. Detailed Service Boundaries

### 6.1 API Gateway

#### Purpose

The `api-gateway` is the public entry point for client applications.

It hides internal service topology from external clients and provides a stable client-facing API surface.

#### Owns

```text
client-facing API routes
request validation at the edge
authentication enforcement at the edge
request correlation IDs
response mapping
error mapping
rate limiting design
client protocol adaptation
```

#### Does Not Own

```text
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

#### Inbound Interfaces

The API Gateway receives external client requests.

The external API may be REST, GraphQL, gRPC-Web, or another client-friendly interface. The final protocol choice should be documented in an ADR.

#### Outbound Calls

The API Gateway may call:

```text
auth-service
customer-service
catalog-service
basket-service
order-service
search-service
recommendation-service
review-service
```

#### Boundary Rules

- The gateway must not contain core business rules.
- The gateway must not directly access service databases.
- The gateway may perform request shape validation, but domain validation belongs in the owning service.
- The gateway should propagate correlation IDs across all downstream calls.
- The gateway should map internal service errors to safe client-facing errors.
- The gateway should not become a monolith in front of the microservices.

#### Correct Responsibilities

```text
Validate request format.
Check authentication token.
Add or propagate correlation ID.
Call order-service CreateOrder.
Map service response to client response.
```

#### Incorrect Responsibilities

```text
Calculate final order state.
Reserve stock directly.
Write to order database.
Authorise payment directly.
Send order confirmation directly.
```

---

### 6.2 Auth Service

#### Purpose

The `auth-service` owns authentication and authorisation identity concerns.

It answers:

```text
Who is the user?
Can this user be issued a token?
What roles or permissions does this identity have?
```

#### Owns

```text
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

#### Does Not Own

```text
customer profile
delivery addresses
orders
payments
basket contents
reviews
```

#### Inbound APIs

Potential gRPC APIs:

```text
RegisterIdentity
Authenticate
ValidateToken
RefreshToken
RevokeSession
GetIdentity
```

#### Events Published

Potential events:

```text
IdentityCreated
CustomerSignedIn
FailedLoginAttempted
SessionRevoked
```

#### Events Consumed

Potential events:

```text
CustomerDeleted
CustomerDisabled
```

#### Data Owned

```text
identities
credentials
sessions
roles
permissions
login_attempts
```

#### Boundary Rules

- Auth Service owns credentials.
- Customer Service owns customer profile data.
- Auth Service must not store customer addresses or order history.
- Passwords must be hashed and never stored in plain text.
- Authentication tokens must not be logged.
- Customer-facing account data should be split carefully between identity and profile concerns.

---

### 6.3 Customer Service

#### Purpose

The `customer-service` owns customer profile information, delivery addresses, and preferences.

It answers:

```text
Who is this customer from a business profile perspective?
Where can orders be delivered?
What preferences has the customer set?
```

#### Owns

```text
customer profile
customer addresses
default delivery address
customer preferences
contact preferences
profile status
```

#### Does Not Own

```text
authentication credentials
tokens
basket contents
orders
payments
shipments
reviews
```

#### Inbound APIs

Potential gRPC APIs:

```text
CreateCustomerProfile
GetCustomerProfile
UpdateCustomerProfile
AddAddress
UpdateAddress
DeleteAddress
ListAddresses
SetDefaultAddress
```

#### Events Published

```text
CustomerProfileCreated
CustomerProfileUpdated
CustomerAddressAdded
CustomerAddressUpdated
CustomerAddressDeleted
```

#### Events Consumed

```text
IdentityCreated
CustomerRegistered
```

#### Data Owned

```text
customers
customer_profiles
addresses
customer_preferences
```

#### Boundary Rules

- Customer Service owns addresses, but orders and shipments may store address snapshots for historical accuracy.
- Customer Service must not own credentials or sessions.
- Customer PII must be protected and should not be logged unnecessarily.
- Other services should request customer profile data via API or consume events where appropriate.
- Customer Service should not own order history.

---

### 6.4 Catalog Service

#### Purpose

The `catalog-service` owns the product catalogue.

It answers:

```text
What products does ACME sell?
What are their details?
Which products are active and purchasable?
How are products categorised?
```

#### Owns

```text
products
product variants
categories
product image metadata
product attributes
materials
colours
dimensions
product pricing
product active/inactive state
```

#### Does Not Own

```text
stock quantity
stock reservations
basket state
order history
customer reviews
search index projections
recommendation results
```

#### Inbound APIs

Potential gRPC APIs:

```text
GetProduct
ListProducts
CreateProduct
UpdateProduct
ActivateProduct
DeactivateProduct
ListCategories
GetCategory
```

#### Events Published

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
ProductDiscontinued
CategoryCreated
CategoryUpdated
```

#### Events Consumed

Usually none for the initial version.

Potential later events:

```text
ReviewApproved
InventoryAdjusted
```

These should only update local display summaries if a deliberate projection is introduced.

#### Data Owned

```text
products
product_variants
categories
product_images
product_attributes
product_prices
```

#### Boundary Rules

- Catalog Service is the source of truth for product details.
- Inventory Service is the source of truth for stock.
- Search Service may index product data but must not become the product source of truth.
- Recommendation Service may use product data but must not own product facts.
- Order Service should store product snapshots for historical order accuracy.
- Product status must be checked before products are added to baskets or shown in customer-facing flows.

---

### 6.5 Inventory Service

#### Purpose

The `inventory-service` owns stock levels, warehouse stock, and stock reservations.

It answers:

```text
Is stock available?
Can stock be reserved?
Has reserved stock been committed, released, or expired?
```

#### Owns

```text
inventory items
stock levels
warehouses
stock reservations
reservation expiry
stock adjustments
stock commitment
```

#### Does Not Own

```text
product descriptions
product pricing
basket contents
order lifecycle
payment state
shipment tracking
```

#### Inbound APIs

Potential gRPC APIs:

```text
CheckAvailability
ReserveStock
ReleaseReservation
CommitReservation
GetStockLevel
AdjustStock
```

#### Events Published

```text
InventoryAdjusted
StockReserved
StockReservationFailed
StockReservationReleased
StockReservationExpired
StockCommitted
```

#### Events Consumed

Potential events:

```text
ProductCreated
ProductDiscontinued
OrderCancelled
PaymentFailed
```

#### Data Owned

```text
inventory_items
warehouses
stock_reservations
stock_adjustments
```

#### Boundary Rules

- Inventory Service is the only service that can change stock levels.
- Basket Service may check availability but must not reserve stock.
- Order Service requests reservations; it does not write stock records.
- Stock reservations should be idempotent where possible.
- Reservation expiry should be explicit and observable.
- Inventory Service must prevent stock from being reserved below zero.

---

### 6.6 Basket Service

#### Purpose

The `basket-service` owns customer baskets and basket items.

It answers:

```text
What is the customer currently intending to buy?
What items and quantities are in the basket?
```

#### Owns

```text
baskets
basket items
basket lifecycle
guest basket reference
basket expiry
```

#### Does Not Own

```text
stock reservation
order creation
payment state
product source of truth
customer profile source of truth
```

#### Inbound APIs

Potential gRPC APIs:

```text
CreateBasket
GetBasket
AddItem
UpdateItem
RemoveItem
ClearBasket
MarkBasketCheckedOut
```

#### Outbound Calls

May call:

```text
catalog-service
inventory-service
```

#### Events Published

```text
BasketCreated
BasketItemAdded
BasketItemUpdated
BasketItemRemoved
BasketCheckedOut
BasketAbandoned
```

#### Events Consumed

Potential events:

```text
OrderCreated
ProductDeactivated
```

#### Data Owned

```text
baskets
basket_items
```

#### Boundary Rules

- Basket Service does not reserve stock.
- Basket Service should validate that products are active before adding them.
- Basket Service may show availability, but Inventory Service owns availability truth.
- Basket price data should be treated carefully. Final price should be confirmed during checkout or stored as an explicit snapshot.
- Basket Service should not create orders.
- Basket Service should support idempotent item updates where practical.

---

### 6.7 Order Service

#### Purpose

The `order-service` owns order creation, order lifecycle, and order history.

It coordinates the checkout workflow.

It answers:

```text
Has an order been created?
What is the current order state?
What items were purchased?
What happened during checkout?
```

#### Owns

```text
orders
order items
order status
order lifecycle
order status history
checkout orchestration
order failure state
order cancellation state
```

#### Does Not Own

```text
product catalogue source of truth
stock source of truth
payment processing internals
shipment fulfilment internals
notification delivery
customer identity
```

#### Inbound APIs

Potential gRPC APIs:

```text
CreateOrder
GetOrder
ListCustomerOrders
CancelOrder
GetOrderStatus
```

#### Outbound Calls

May call:

```text
basket-service
inventory-service
payment-service
shipping-service
customer-service
```

#### Events Published

```text
OrderCreated
OrderConfirmed
OrderFailed
OrderCancelled
OrderFulfilled
OrderRefunded
```

#### Events Consumed

Potential events:

```text
PaymentAuthorised
PaymentFailed
ShipmentCreated
ShipmentFailed
StockReserved
StockReservationFailed
```

The first version may use synchronous responses for checkout and consume fewer events.

#### Data Owned

```text
orders
order_items
order_status_history
checkout_attempts
```

#### Boundary Rules

- Order Service owns order lifecycle but does not own payment state.
- Order Service owns order items and historical snapshots.
- Order Service should not directly modify stock.
- Order Service should not directly send notifications.
- Order creation must be idempotent where possible.
- Duplicate checkout requests must not create duplicate confirmed orders.
- Order Service may coordinate checkout because order creation is the business process being completed.
- Long-running downstream side effects should be event-driven where possible.

#### Important Design Decision

The initial architecture proposes that checkout orchestration lives in `order-service`.

This is acceptable because order creation is the business process being coordinated. If orchestration grows significantly, a later ADR may consider a dedicated workflow or orchestration component.

---

### 6.8 Payment Service

#### Purpose

The `payment-service` owns payment attempts, payment state, authorisation, capture, and refunds.

It answers:

```text
Was payment authorised?
Was payment captured?
Did payment fail?
Was a refund issued?
```

#### Owns

```text
payments
payment attempts
payment authorisation state
payment capture state
refunds
provider references
payment audit records
```

#### Does Not Own

```text
order lifecycle
stock reservation
shipment creation
customer profile
raw card storage
```

#### Inbound APIs

Potential gRPC APIs:

```text
AuthorisePayment
CapturePayment
RefundPayment
GetPayment
ListPaymentAttempts
```

#### Events Published

```text
PaymentAuthorised
PaymentFailed
PaymentCaptured
PaymentRefunded
PaymentCancelled
```

#### Events Consumed

Potential events:

```text
OrderCreated
OrderCancelled
```

#### Data Owned

```text
payments
payment_attempts
refunds
payment_provider_references
```

#### Boundary Rules

- Payment Service owns payment state.
- Order Service owns order state.
- Raw card data must not be stored.
- Sensitive payment details must not be logged.
- Payment requests should be idempotent where possible.
- Payment failures must be auditable.
- Real payment provider integration is out of scope for the initial version; simulation is acceptable.

---

### 6.9 Shipping Service

#### Purpose

The `shipping-service` owns delivery options, shipment creation, and shipment lifecycle.

It answers:

```text
Can this order be shipped?
What shipment was created?
What is the fulfilment status?
What is the tracking reference?
```

#### Owns

```text
delivery options
shipments
shipment status
tracking references
shipment events
delivery state
```

#### Does Not Own

```text
order creation
payment state
stock ownership
customer profile source of truth
notification delivery
```

#### Inbound APIs

Potential gRPC APIs:

```text
GetDeliveryOptions
CreateShipment
GetShipment
UpdateShipmentStatus
CancelShipment
```

#### Events Published

```text
ShipmentCreated
ShipmentDispatched
ShipmentDelivered
ShipmentDelayed
ShipmentFailed
ShipmentCancelled
```

#### Events Consumed

Potential events:

```text
OrderCreated
OrderCancelled
PaymentAuthorised
```

#### Data Owned

```text
delivery_options
shipments
shipment_status_history
tracking_events
```

#### Boundary Rules

- Shipping Service owns shipment state, not order state.
- Order Service may request shipment creation.
- Shipment failure must be visible to Order Service and Operations.
- Live carrier integration is out of scope initially.
- Shipment records should use delivery address snapshots, not depend on mutable customer address records.
- Shipment creation should be idempotent where possible.

---

### 6.10 Notification Service

#### Purpose

The `notification-service` owns customer-facing notifications and delivery status.

It answers:

```text
Was a notification requested?
Was it sent?
Did it fail?
Should it be retried?
```

#### Owns

```text
notification requests
notification templates
notification delivery attempts
notification status
retry state
provider references
```

#### Does Not Own

```text
order lifecycle
payment state
shipment state
customer profile source of truth
business decision to create an order
```

#### Inbound APIs

Potential gRPC APIs:

```text
RequestNotification
GetNotificationStatus
```

The initial version may primarily consume Kafka events rather than expose many synchronous APIs.

#### Events Published

```text
NotificationSent
NotificationFailed
NotificationRetryScheduled
```

#### Events Consumed

```text
OrderCreated
OrderCancelled
PaymentFailed
ShipmentCreated
ShipmentDispatched
ShipmentDelivered
NotificationRequested
```

#### Data Owned

```text
notifications
notification_attempts
notification_templates
```

#### Boundary Rules

- Notification failure must not roll back order creation.
- Notification processing must be idempotent.
- Duplicate events should not cause duplicate customer messages where avoidable.
- Notification Service may request customer contact details, but Customer Service owns the profile.
- Real email/SMS provider integration is out of scope initially; simulation is acceptable.

---

### 6.11 Review Service

#### Purpose

The `review-service` owns product reviews, ratings, and moderation state.

It answers:

```text
What reviews has this product received?
What rating did a customer submit?
Is the review approved, rejected, or pending moderation?
```

#### Owns

```text
reviews
ratings
rating summaries
moderation status
review reports
```

#### Does Not Own

```text
product catalogue source of truth
customer profile source of truth
order lifecycle
search index source of truth
```

#### Inbound APIs

Potential gRPC APIs:

```text
SubmitReview
GetReview
ListProductReviews
ModerateReview
GetRatingSummary
```

#### Outbound Calls

May call:

```text
catalog-service
order-service
customer-service
```

#### Events Published

```text
ReviewCreated
ReviewApproved
ReviewRejected
ReviewDeleted
RatingSummaryUpdated
```

#### Events Consumed

Potential events:

```text
ProductDeactivated
OrderFulfilled
```

#### Data Owned

```text
reviews
rating_summaries
moderation_decisions
```

#### Boundary Rules

- Review Service owns review data, not product data.
- Product existence should be validated against Catalog Service.
- Purchase eligibility may be validated against Order Service.
- Review summaries may be eventually consistent.
- Search and Recommendation Services may consume approved review events.
- Reviews should not be allowed to alter catalogue product ownership.

---

### 6.12 Search Service

#### Purpose

The `search-service` owns the product search projection and query model.

It answers:

```text
Which products match this search query and filter set?
```

#### Owns

```text
search index
search documents
facets
search query logs
index update status
```

#### Does Not Own

```text
product source of truth
stock source of truth
review source of truth
recommendation logic
```

#### Inbound APIs

Potential gRPC APIs:

```text
SearchProducts
SuggestProducts
GetSearchFacets
```

#### Events Published

```text
SearchIndexUpdated
SearchIndexUpdateFailed
```

#### Events Consumed

```text
ProductCreated
ProductUpdated
ProductDeactivated
InventoryAdjusted
ReviewApproved
RatingSummaryUpdated
```

#### Data Owned

```text
search_index_entries
search_facets
search_query_logs
index_update_offsets
```

#### Boundary Rules

- Search Service owns an optimised projection, not the source catalogue.
- Search results may be eventually consistent.
- Inactive products must not appear in customer-facing search results.
- Search index rebuilds must be possible.
- Search Service should tolerate duplicate product update events.
- Search Service must not become a dependency for successful checkout.

---

### 6.13 Recommendation Service

#### Purpose

The `recommendation-service` owns recommendation logic, recommendation signals, and calculated recommendation outputs.

It answers:

```text
Which products should be recommended to this customer or alongside this product?
```

#### Owns

```text
recommendation rules
recommendation signals
calculated recommendations
recommendation result cache
recommendation feedback
```

#### Does Not Own

```text
product catalogue
orders
payments
reviews
search index
customer profile source of truth
```

#### Inbound APIs

Potential gRPC APIs:

```text
GetProductRecommendations
GetBasketRecommendations
GetCustomerRecommendations
RecordRecommendationFeedback
```

#### Events Published

```text
RecommendationGenerated
RecommendationFeedbackRecorded
```

#### Events Consumed

```text
ProductViewed
BasketItemAdded
OrderCreated
ReviewCreated
ProductUpdated
ProductDeactivated
```

#### Data Owned

```text
recommendation_rules
recommendation_signals
recommendation_results
recommendation_feedback
```

#### Boundary Rules

- Recommendations must not include inactive products.
- Recommendation outputs should degrade gracefully when data is limited.
- Initial implementation may be rules-based.
- Recommendation Service must not become the source of truth for product, order, or customer data.
- Recommendation projections may be eventually consistent.
- Recommendation Service must not block catalogue browsing or checkout.

---

## 7. Boundary Rules for the Checkout Flow

The checkout flow is the most important initial service interaction.

```text
Customer
    -> API Gateway
    -> Order Service
        -> Basket Service
        -> Inventory Service
        -> Payment Service
        -> Shipping Service
        -> Kafka
            -> Notification Service
```

### 7.1 Ownership During Checkout

| Step | Owner | Notes |
|---|---|---|
| Retrieve basket | `basket-service` | Basket Service owns current basket state |
| Validate product status | `catalog-service` | Product active/inactive state belongs to Catalogue |
| Reserve stock | `inventory-service` | Inventory owns stock and reservations |
| Authorise payment | `payment-service` | Payment owns authorisation result and payment attempts |
| Create order | `order-service` | Order owns order lifecycle and item snapshots |
| Create shipment | `shipping-service` | Shipping owns shipment state |
| Send notification | `notification-service` | Notification owns delivery status |

### 7.2 Checkout Boundary Rules

- Basket Service does not create orders.
- Order Service does not update stock tables.
- Payment Service does not mark orders as confirmed.
- Shipping Service does not decide order success.
- Notification Service does not affect whether an order is created.
- Search and Recommendation Services are not part of the critical checkout path.
- Downstream events must include correlation IDs.
- Idempotency is required for checkout, payment, stock reservation, shipment creation, and notification sending.

---

## 8. Data Ownership Rules

### 8.1 Direct Database Access

Direct database access across services is forbidden.

Allowed:

```text
order-service -> inventory-service ReserveStock gRPC
inventory-service -> StockReserved Kafka event
```

Forbidden:

```text
order-service -> inventory_db.stock_reservations
```

### 8.2 Snapshots Are Allowed

Services may store snapshots for historical or operational reasons.

Examples:

```text
order-service stores product_name_snapshot
order-service stores unit_price_snapshot
order-service stores delivery_address_snapshot
shipping-service stores delivery_address_snapshot
payment-service stores provider_reference
```

Snapshots must be clearly labelled as snapshots and must not be treated as the source of truth for future business decisions.

### 8.3 Projections Are Allowed

Services may store projections optimised for their own use.

Examples:

```text
search-service stores product search documents
recommendation-service stores recommendation signals
notification-service stores customer notification status
```

Projection owners must be clear that source-of-truth data lives elsewhere.

---

## 9. API Boundary Rules

### 9.1 gRPC APIs

gRPC APIs are used when a service needs an immediate response.

Examples:

```text
GetProduct
AddBasketItem
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
```

### 9.2 API Ownership

Each service owns its own gRPC API.

Rules:

- The owning service defines the API contract.
- Protobuf definitions must be versioned.
- Breaking changes must be detected and controlled.
- Consumers must not infer internal database structure from API messages.
- API messages should represent service behaviour, not raw table models.

### 9.3 API Gateway Boundary

The API Gateway may expose a different client-facing shape from internal gRPC APIs.

For example:

```text
external REST/JSON request
    -> internal gRPC service call
```

The API Gateway should not become a business logic monolith.

---

## 10. Event Boundary Rules

### 10.1 Kafka Events

Kafka events represent facts that have already happened.

Good event names:

```text
OrderCreated
PaymentFailed
ShipmentCreated
StockReserved
```

Poor event names:

```text
CreateOrder
SendNotification
ReserveStock
```

Commands belong in APIs. Facts belong in events.

### 10.2 Event Ownership

The service that owns the business fact publishes the event.

| Event | Producer |
|---|---|
| `ProductUpdated` | `catalog-service` |
| `StockReserved` | `inventory-service` |
| `BasketCheckedOut` | `basket-service` |
| `OrderCreated` | `order-service` |
| `PaymentAuthorised` | `payment-service` |
| `ShipmentCreated` | `shipping-service` |
| `NotificationSent` | `notification-service` |
| `ReviewCreated` | `review-service` |

### 10.3 Consumer Responsibilities

Consumers must:

- be idempotent
- tolerate duplicate events
- handle out-of-order events where relevant
- record processing failures
- use dead-letter queues where appropriate
- not assume synchronous completion
- not make the producer responsible for consumer behaviour

---

## 11. Shared Package Boundary

Shared packages live under:

```text
packages/go/
```

Allowed shared package examples:

```text
logger
config
grpc
kafka
telemetry
auth helpers
errors
health
middleware
testkit
```

Rules:

- Shared packages must be technical, not business-domain heavy.
- Business rules must remain inside owning services.
- Shared packages should avoid becoming a hidden platform framework.
- Shared packages should remain small and well tested.
- If a shared package grows too much business logic, move the logic back into the owning service.

---

## 12. Anti-Patterns to Avoid

### 12.1 Shared Database

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

### 12.2 Distributed Monolith

Avoid:

```text
every request requires every service
every deployment requires all services
business rules are spread across many services
shared packages contain core domain logic
```

A microservice architecture should allow services to change independently. If every change requires coordinated releases across many services, the architecture has become a distributed monolith.

---

### 12.3 Chatty Service Design

Avoid designs where simple operations require many synchronous calls.

Example risk:

```text
api-gateway -> catalog-service -> inventory-service -> review-service -> recommendation-service
```

Prefer:

- API composition where needed
- read models or projections where justified
- asynchronous updates for non-critical data
- clear ownership of critical synchronous paths

---

### 12.4 Kafka as RPC

Avoid publishing an event and waiting for another service to treat it like a synchronous command.

If an immediate result is required, use gRPC.

---

### 12.5 API Gateway as a Monolith

Avoid putting order, payment, stock, or shipping business logic in the gateway.

The gateway should route, authenticate, shape requests, map responses, and propagate context.

---

## 13. Boundary Decision Checklist

Before creating a new service, ask:

```text
Does it own a clear business capability?
Does it have its own data?
Can it expose a clear API?
Does it publish or consume meaningful events?
Can it be tested independently?
Can it be deployed independently?
Can it fail independently?
Can it be observed independently?
Does it avoid duplicating another service’s business ownership?
```

Before adding functionality to an existing service, ask:

```text
Does this capability belong to this service?
Is this service the source of truth for the data involved?
Would another team expect to own this capability?
Will this introduce cross-service coupling?
Does this change require an ADR?
```

---

## 14. Boundary Change Process

Service boundaries may evolve as the project matures.

A boundary change should be documented when it affects:

```text
service ownership
database ownership
gRPC contracts
Kafka events
checkout orchestration
security responsibilities
operational responsibilities
deployment responsibilities
```

Significant changes should be recorded as ADRs in:

```text
adr/
```

Example ADRs:

```text
0001-use-microservices.md
0002-use-grpc-for-service-communication.md
0003-use-kafka-for-events.md
0004-use-service-owned-databases.md
0005-use-mysql.md
0008-use-contract-first-service-design.md
```

---

## 15. Initial Implementation Boundaries

The first implementation should focus on the minimum useful service set for checkout.

| Service | Include Initially? | Reason |
|---|---:|---|
| `api-gateway` | Yes | Entry point for client workflows |
| `catalog-service` | Yes | Product browsing and product validation |
| `inventory-service` | Yes | Stock reservation |
| `basket-service` | Yes | Basket management |
| `order-service` | Yes | Checkout orchestration and order creation |
| `payment-service` | Yes | Payment authorisation simulation |
| `shipping-service` | Yes | Shipment creation |
| `notification-service` | Yes | Event-driven notification |
| `auth-service` | Optional | Can be simplified or mocked for first slice |
| `customer-service` | Optional | Can use test customer/address initially |
| `review-service` | Later | Not required for checkout |
| `search-service` | Later | Basic catalogue list can come from Catalog Service first |
| `recommendation-service` | Later | Not required for checkout |

This staged approach avoids building many incomplete services before the core flow works.

---

## 16. Target Boundary Maturity

As the project matures, each service should have:

```text
owned database/schema
protobuf API contract
Kafka event contract
service README
service requirements document
service architecture notes
database migrations
unit tests
integration tests
contract tests
health check
readiness check
structured logs
metrics
traces
dashboard
runbook
deployment manifest/chart
resource requests and limits
security notes
production readiness evidence
```

This is the level of evidence expected for a senior platform engineering portfolio.

---

## 17. Open Questions

| Question | Status |
|---|---|
| Should checkout orchestration remain in `order-service` long term? | Proposed for initial version |
| Should shipment creation block order confirmation? | To decide |
| Should notifications consume `OrderCreated` directly or a dedicated `NotificationRequested` event? | To decide |
| Should `auth-service` be included in the first vertical slice or mocked? | To decide |
| Should `customer-service` be required for initial checkout or should a test customer/address be used? | To decide |
| Should search initially be handled by `catalog-service` before introducing `search-service`? | Proposed |
| Should recommendations be rules-based before introducing event-driven personalised recommendations? | Proposed |
| Should service mesh be introduced later for service-to-service identity and traffic policy? | Deferred |
| Should payment authorisation and capture be separate flows? | To decide |
| Should guest checkout be supported in the first implementation? | To decide |

---

## 18. Related Documents

This document should be read alongside:

```text
docs/requirements/product-vision.md
docs/requirements/scope.md
docs/requirements/user-journeys.md
docs/architecture/domain-model.md
docs/architecture/communication-patterns.md
docs/architecture/event-driven-design.md
docs/data/data-ownership.md
docs/data/service-database-design.md
docs/events/event-catalog.md
docs/api/grpc-overview.md
```

Relevant ADRs:

```text
adr/0001-use-microservices.md
adr/0002-use-grpc-for-service-communication.md
adr/0003-use-kafka-for-events.md
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
adr/0008-use-contract-first-service-design.md
```

---

## 19. Summary

bfstore’s service boundaries are designed around business capabilities.

The most important boundary rules are:

```text
Catalog owns products.
Inventory owns stock.
Basket owns basket state.
Order owns order lifecycle.
Payment owns payment state.
Shipping owns shipment state.
Notification owns notification delivery.
Review owns reviews.
Search owns search projections.
Recommendation owns recommendation outputs.
```

Services must not share databases or silently own each other’s business rules.

The checkout flow coordinates multiple services, but each service remains responsible for its own domain:

```text
Basket -> current buying intent
Inventory -> stock reservation
Payment -> payment authorisation
Order -> order lifecycle
Shipping -> shipment creation
Notification -> customer communication
```

This service boundary model supports a professional microservice architecture that can be tested, deployed, secured, observed, and operated as part of the wider ACME platform engineering estate.
