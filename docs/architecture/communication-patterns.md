# Communication Patterns

## 1. Purpose

This document defines the communication patterns used by **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how services communicate, when to use synchronous gRPC, when to use asynchronous Kafka events, how the API Gateway fits into the architecture, and how communication patterns support reliability, observability, testing, and clear service ownership.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s architecture.

---

## 2. Architecture Context

bfstore is a cloud-native microservice backend for an online furniture store.

The system uses:

| Technology | Purpose |
|---|---|
| API Gateway | Client-facing entry point |
| gRPC | Synchronous internal service-to-service communication |
| Protobuf | Strongly typed API and event contracts |
| Kafka | Asynchronous event-driven communication |
| MySQL | Service-owned relational persistence |
| OpenTelemetry | Logs, metrics, traces, and correlation |
| Kubernetes | Target runtime platform |

bfstore follows a hybrid communication model:

> Use gRPC for commands and queries that require an immediate response. Use Kafka for facts that have already happened and can be processed asynchronously.

---

## 3. Communication Principles

## 3.1 Commands Use gRPC

Use gRPC when the caller needs an immediate response before it can continue.

Examples:

```text
GetProduct
ListProducts
AddBasketItem
GetBasket
ReserveStock
AuthorisePayment
CreateOrder
CreateShipment
```

These operations require a clear response such as:

```text
success
failure
not found
validation error
insufficient stock
payment declined
```

---

## 3.2 Events Use Kafka

Use Kafka when a service needs to publish that something has already happened.

Examples:

```text
ProductUpdated
StockReserved
PaymentAuthorised
ShipmentCreated
OrderCreated
NotificationSent
ReviewApproved
```

Events should be consumed by services that need to react later.

Examples:

```text
OrderCreated -> notification-service sends confirmation
ProductUpdated -> search-service updates projection
ReviewApproved -> recommendation-service updates signals
```

---

## 3.3 The Database Is Not a Communication Channel

Services must not communicate by reading or writing each other’s databases.

Allowed:

```text
order-service -> inventory-service ReserveStock gRPC
inventory-service -> StockReserved Kafka event
```

Forbidden:

```text
order-service -> bfstore_inventory.stock_reservations table
```

The API is the contract. The database is private to the owning service.

---

## 3.4 Communication Must Preserve Service Ownership

The service that owns the business capability owns the decision.

Examples:

| Decision | Owning Service |
|---|---|
| Is product active? | `catalog-service` |
| Is stock available? | `inventory-service` |
| Is basket valid? | `basket-service` |
| Is order confirmed? | `order-service` |
| Was payment authorised? | `payment-service` |
| Was shipment created? | `shipping-service` |
| Was notification sent? | `notification-service` |

Other services may request decisions or consume events, but they should not duplicate ownership.

---

## 3.5 Communication Must Be Observable

Every cross-service interaction should support:

```text
correlation ID
trace ID
structured logs
metrics
timeout behaviour
error classification
service name
operation name
```

A checkout request should be traceable across:

```text
api-gateway
order-service
basket-service
inventory-service
payment-service
shipping-service
Kafka
notification-service
```

---

## 4. Communication Types

## 4.1 Client to API Gateway

External clients communicate with the API Gateway.

```text
Client
    -> API Gateway
```

The client-facing protocol may be:

```text
REST/JSON
GraphQL
gRPC-Web
```

The final external protocol should be documented in an ADR.

The API Gateway should:

- authenticate requests where required
- validate request shape
- apply rate limiting where appropriate
- add or propagate correlation IDs
- call backend services
- map internal errors to safe external responses
- avoid owning business rules

The API Gateway should not:

- own business entities
- write service databases
- coordinate complex domain workflows directly
- become a monolith in front of the services

---

## 4.2 API Gateway to Backend Services

The API Gateway communicates with internal services using gRPC.

Example:

```text
Client
    -> API Gateway
        -> catalog-service ListProducts
        -> basket-service AddItem
        -> order-service CreateOrder
```

The gateway may compose responses from multiple services, but should avoid excessive synchronous fan-out.

---

## 4.3 Service to Service gRPC

Services use gRPC for immediate business operations.

Example checkout calls:

```text
order-service -> basket-service GetBasket
order-service -> inventory-service ReserveStock
order-service -> payment-service AuthorisePayment
order-service -> shipping-service CreateShipment
```

gRPC calls should have:

```text
timeouts
retries where safe
clear error codes
request validation
correlation IDs
trace propagation
idempotency keys where required
```

---

## 4.4 Service to Kafka

Services publish Kafka events when business facts occur.

Example:

```text
order-service -> Kafka -> OrderCreated
payment-service -> Kafka -> PaymentAuthorised
shipping-service -> Kafka -> ShipmentCreated
```

Kafka events allow downstream services to react without blocking the original business operation.

---

## 4.5 Kafka to Service Consumers

Services consume Kafka events when they need to update projections, trigger side effects, or maintain operational records.

Example:

```text
OrderCreated
    -> notification-service
    -> recommendation-service
    -> analytics consumers
```

Consumers must be:

```text
idempotent
observable
retryable
DLQ-aware
able to tolerate duplicate events
able to tolerate eventual consistency
```

---

## 5. When to Use gRPC

Use gRPC when:

```text
the caller needs an immediate result
the operation is a command or query
the caller must know success or failure before continuing
the response affects the current request path
the interaction requires validation
the operation is part of the critical checkout path
```

Examples:

| Operation | Reason for gRPC |
|---|---|
| `GetProduct` | Client needs product details immediately |
| `AddBasketItem` | Client needs updated basket immediately |
| `ReserveStock` | Order Service must know whether stock is reserved |
| `AuthorisePayment` | Order Service must know whether payment succeeded |
| `CreateShipment` | Order Service needs shipment result or failure |
| `GetOrder` | Client needs order details immediately |

---

## 6. When to Use Kafka

Use Kafka when:

```text
something has already happened
multiple consumers may care
the producer should not block on consumers
eventual consistency is acceptable
the event feeds projections or side effects
the action is not required to complete the current request
```

Examples:

| Event | Reason for Kafka |
|---|---|
| `OrderCreated` | Notifications, recommendations, and analytics can react asynchronously |
| `ProductUpdated` | Search index can update asynchronously |
| `ReviewApproved` | Rating summaries and search projections can update later |
| `ShipmentDispatched` | Customer notification can be sent asynchronously |
| `NotificationFailed` | Operations can monitor and alert asynchronously |

---

## 7. Decision Matrix

| Question | Use gRPC | Use Kafka |
|---|---:|---:|
| Does the caller need an immediate response? | Yes | No |
| Is this a command or query? | Yes | No |
| Is this a fact that already happened? | No | Yes |
| Can downstream processing happen later? | No | Yes |
| Can there be multiple independent consumers? | Sometimes | Yes |
| Should the producer block until all consumers finish? | No | No |
| Is eventual consistency acceptable? | Sometimes | Yes |
| Is this on the critical checkout path? | Usually | Only for downstream side effects |

---

## 8. Checkout Communication Pattern

The initial checkout flow is the most important communication path.

```text
Client
    -> API Gateway
        -> order-service CreateOrder
            -> basket-service GetBasket
            -> inventory-service ReserveStock
            -> payment-service AuthorisePayment
            -> shipping-service CreateShipment
            -> Kafka OrderCreated
                -> notification-service
```

## 8.1 Synchronous Checkout Calls

| Call | Purpose |
|---|---|
| `order-service -> basket-service` | Retrieve and validate basket contents |
| `order-service -> inventory-service` | Reserve stock |
| `order-service -> payment-service` | Authorise payment |
| `order-service -> shipping-service` | Create shipment |
| `api-gateway -> order-service` | Submit checkout and return result |

## 8.2 Asynchronous Checkout Events

| Event | Producer | Consumer |
|---|---|---|
| `StockReserved` | `inventory-service` | `order-service`, operations consumers |
| `PaymentAuthorised` | `payment-service` | `order-service`, notification consumers |
| `ShipmentCreated` | `shipping-service` | `order-service`, notification consumers |
| `OrderCreated` | `order-service` | `notification-service`, recommendation consumers |
| `NotificationSent` | `notification-service` | operations consumers |

## 8.3 Checkout Communication Rules

- Stock must be reserved before the order is confirmed.
- Payment must be authorised before the order is confirmed.
- Notification must not block order creation.
- Search and recommendations must not be part of the critical checkout path.
- Duplicate checkout requests must not create duplicate confirmed orders.
- All calls and events must carry correlation context.

---

## 9. Product Browse Communication Pattern

Product browsing should be simple and fast.

Initial version:

```text
Client
    -> API Gateway
        -> catalog-service ListProducts
```

Optional availability enrichment:

```text
Client
    -> API Gateway
        -> catalog-service ListProducts
        -> inventory-service CheckAvailability
```

Later search version:

```text
Client
    -> API Gateway
        -> search-service SearchProducts
```

Search Service may consume product and inventory events:

```text
catalog-service -> ProductUpdated -> search-service
inventory-service -> InventoryAdjusted -> search-service
```

## 9.1 Browse Rules

- Catalog Service owns product truth.
- Search Service owns search projection only.
- Inventory Service owns availability truth.
- Search may be eventually consistent.
- Product browsing should not depend on recommendation services.

---

## 10. Notification Communication Pattern

Notification processing should be asynchronous.

Option A: Notification Service consumes domain events directly.

```text
order-service -> OrderCreated -> notification-service
shipping-service -> ShipmentDispatched -> notification-service
payment-service -> PaymentFailed -> notification-service
```

Option B: Domain services publish explicit notification requests.

```text
order-service -> NotificationRequested -> notification-service
```

## 10.1 Recommendation

For the initial version, consuming `OrderCreated` directly is simple and acceptable.

For more complex notification workflows, `NotificationRequested` can make notification intent explicit.

This decision should be captured in an ADR if it affects implementation.

---

## 11. Search and Projection Communication Pattern

Search should use projections rather than direct database access.

```text
catalog-service -> ProductCreated -> search-service
catalog-service -> ProductUpdated -> search-service
catalog-service -> ProductDeactivated -> search-service
inventory-service -> InventoryAdjusted -> search-service
review-service -> ReviewApproved -> search-service
```

Search Service stores:

```text
search_index_entries
search_facets
projection_offsets
```

## 11.1 Projection Rules

- Projection updates may be eventually consistent.
- Projection consumers must be idempotent.
- Projection rebuilds should be possible.
- Projection lag should be observable.
- Search Service must not become the product source of truth.

---

## 12. Recommendation Communication Pattern

Recommendation Service may consume behavioural and product events.

```text
ProductViewed
BasketItemAdded
OrderCreated
ReviewApproved
ProductUpdated
```

Initial recommendations may be rules-based:

```text
same category
same material
same colour
popular products
frequently bought together
```

## 12.1 Recommendation Rules

- Recommendations must not include inactive products.
- Recommendation failure must not block browsing or checkout.
- Recommendation data may be eventually consistent.
- Recommendation Service must not own product, order, or review source-of-truth data.

---

## 13. Error Handling

## 13.1 gRPC Errors

gRPC errors should be standardised.

Common error categories:

```text
INVALID_ARGUMENT
NOT_FOUND
FAILED_PRECONDITION
UNAUTHENTICATED
PERMISSION_DENIED
CONFLICT
UNAVAILABLE
DEADLINE_EXCEEDED
INTERNAL
```

Examples:

| Scenario | Error Category |
|---|---|
| Product does not exist | `NOT_FOUND` |
| Product inactive | `FAILED_PRECONDITION` |
| Basket quantity invalid | `INVALID_ARGUMENT` |
| Insufficient stock | `FAILED_PRECONDITION` |
| Payment declined | `FAILED_PRECONDITION` |
| Service unavailable | `UNAVAILABLE` |
| Downstream timeout | `DEADLINE_EXCEEDED` |

Detailed error design should be documented in:

```text
docs/api/error-model.md
```

---

## 13.2 Kafka Consumer Errors

Kafka consumer failures should be classified as:

| Error Type | Example | Handling |
|---|---|---|
| Retryable | temporary database outage | retry |
| Retryable | temporary network error | retry |
| Non-retryable | invalid event payload | DLQ |
| Non-retryable | unsupported event version | DLQ |
| Business conflict | duplicate event | ignore or idempotent update |

Consumers should expose metrics for:

```text
processed events
failed events
retry count
DLQ count
consumer lag
duplicate events
```

---

## 14. Timeouts and Retries

## 14.1 gRPC Timeouts

Every gRPC call should have a timeout.

Example conceptual timeout policy:

| Call | Timeout Sensitivity |
|---|---|
| API Gateway to Catalog | low to medium |
| API Gateway to Order | medium |
| Order to Inventory | high |
| Order to Payment | high |
| Order to Shipping | high |
| Notification to Customer | asynchronous, not in critical path |

Timeout values should be tuned through testing, not guessed permanently.

## 14.2 Retry Rules

Retries are only safe when the operation is idempotent or has idempotency protection.

Safe to retry when designed properly:

```text
GetProduct
ListProducts
GetBasket
ReserveStock with idempotency key
AuthorisePayment with idempotency key
CreateShipment with idempotency key
```

Dangerous to retry blindly:

```text
CreateOrder
AuthorisePayment
CommitStock
SendNotification
```

These require idempotency keys or duplicate suppression.

---

## 15. Idempotency

Idempotency is required for operations where duplicates could cause business harm.

Important idempotency points:

```text
checkout submission
order creation
stock reservation
payment authorisation
shipment creation
notification sending
event consumption
```

Example idempotency inputs:

```text
idempotency_key
customer_id
basket_id
order_id
payment_request_id
shipment_request_id
event_id
```

---

## 16. Correlation and Tracing

Every request should carry correlation context.

## 16.1 Required Context

```text
correlation_id
trace_id
request_id
causation_id where relevant
```

## 16.2 Propagation

Correlation context should flow through:

```text
client request
API Gateway
gRPC metadata
service logs
Kafka event envelope
consumer logs
distributed traces
```

## 16.3 Why It Matters

A failed checkout should be diagnosable from one correlation ID.

Example:

```text
correlation_id=corr-123
    API Gateway received checkout
    Order Service started checkout
    Basket Service returned basket
    Inventory Service reserved stock
    Payment Service declined payment
    Order Service failed checkout
    Inventory reservation released
```

---

## 17. Security Considerations

Communication security should consider:

```text
authentication at the edge
authorisation for protected actions
service-to-service identity
TLS or mTLS where appropriate
least privilege service accounts
no secrets in messages or logs
PII minimisation in events
safe error responses
```

Events and APIs must not expose:

```text
passwords
tokens
raw payment data
secret values
unnecessary customer PII
```

---

## 18. Testing Requirements

Communication patterns should be tested.

## 18.1 gRPC Tests

```text
contract tests
timeout behaviour tests
error mapping tests
idempotency tests
client/server integration tests
```

## 18.2 Kafka Tests

```text
producer contract tests
consumer contract tests
duplicate event tests
DLQ tests
consumer lag tests
event replay tests
```

## 18.3 End-to-End Tests

Important E2E tests:

```text
successful checkout
insufficient stock checkout failure
payment failure checkout failure
OrderCreated triggers notification
duplicate OrderCreated does not duplicate notification
```

---

## 19. Anti-Patterns to Avoid

## 19.1 Chatty Synchronous Design

Avoid requiring many synchronous calls for simple reads.

Risk:

```text
API Gateway -> Catalog -> Inventory -> Review -> Recommendation
```

Prefer:

- local projections
- asynchronous updates
- careful API composition
- avoiding non-critical services in critical paths

---

## 19.2 Kafka as RPC

Avoid publishing an event and waiting for a specific consumer response.

If an immediate response is required, use gRPC.

---

## 19.3 Shared Database Communication

Avoid using database tables as integration points.

This creates tight coupling and weak service ownership.

---

## 19.4 API Gateway Business Logic

Avoid placing domain decisions in the API Gateway.

The gateway should route and shape traffic, not own checkout, payment, stock, or fulfilment rules.

---

## 20. Initial Implementation Scope

The first implementation should focus on:

```text
API Gateway to Catalog Service
API Gateway to Basket Service
API Gateway to Order Service
Order Service to Basket Service
Order Service to Inventory Service
Order Service to Payment Service
Order Service to Shipping Service
Order Service publishing OrderCreated
Notification Service consuming OrderCreated
```

This proves the core communication model before adding search, reviews, recommendations, and advanced event projections.

---

## 21. Open Questions

| Question | Status |
|---|---|
| Will the API Gateway expose REST, GraphQL, or gRPC-Web? | To decide |
| Should Notification Service consume `OrderCreated` directly or use `NotificationRequested`? | To decide |
| Should shipment creation block order confirmation? | To decide |
| Should payment capture be separate from authorisation in the first version? | To decide |
| Should search initially use Catalog Service directly or a projection-based Search Service? | Proposed: start with Catalog Service, add Search Service later |
| Should service mesh be introduced for mTLS and traffic policy? | Deferred |
| Should outbox pattern be implemented from the first event-producing service? | Proposed for serious implementation |

---

## 22. Related Documents

This document should be read alongside:

```text
docs/architecture/service-boundaries.md
docs/architecture/domain-model.md
docs/architecture/event-driven-design.md
docs/architecture/resilience-patterns.md
docs/api/grpc-overview.md
docs/api/error-model.md
docs/events/event-catalog.md
docs/events/event-envelope.md
docs/data/data-ownership.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0002-use-grpc-for-service-communication.md
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
adr/0008-use-contract-first-service-design.md
```

---

## 23. Summary

bfstore uses a deliberate hybrid communication model:

```text
gRPC for immediate commands and queries
Kafka for asynchronous business facts
API Gateway for client-facing access
Protobuf for contracts
OpenTelemetry for correlation and tracing
```

The communication model supports clear service ownership, safe checkout coordination, asynchronous downstream processing, observability, and professional microservice boundaries.

The first implementation should prove this model through the checkout vertical slice, then expand into search, reviews, recommendations, and more advanced event-driven workflows.
