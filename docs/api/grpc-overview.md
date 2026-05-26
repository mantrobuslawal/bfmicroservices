# gRPC Overview

## 1. Purpose

This document defines the gRPC communication model for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how synchronous service-to-service communication should be designed, documented, implemented, tested, secured, observed, and evolved across the bfstore microservice architecture.

This document is intended for engineers, technical reviewers, platform teams, and potential clients evaluating the quality of the bfstore backend architecture.

---

## 2. Scope

This document covers:

- when gRPC should be used
- when gRPC should not be used
- how protobuf contracts should be organised
- expected service API design principles
- request and response conventions
- error handling
- authentication and authorisation considerations
- timeout and retry behaviour
- observability requirements
- contract testing expectations
- versioning and compatibility rules
- client and server implementation standards

This document does not define every RPC method in detail. Service-specific protobuf definitions live under:

```text
proto/acme/
```

Detailed API documentation for individual services should live in:

```text
docs/api/
services/<service-name>/docs/
proto/acme/<domain>/
```

---

## 3. Architecture Context

bfstore is a cloud-native microservice backend for an online furniture store.

The backend uses:

| Technology | Purpose |
|---|---|
| Go | Primary service implementation language |
| gRPC | Synchronous internal service-to-service communication |
| Protobuf | API and event contract definition |
| Kafka | Asynchronous event-driven communication |
| MySQL | Service-owned relational persistence |
| OpenTelemetry | Distributed tracing, metrics, and logging correlation |
| Buf | Protobuf linting, generation, and breaking-change detection |

The core communication principle is:

> Commands and queries that need an immediate response use gRPC. Facts that have already happened are published as Kafka events.

---

## 4. gRPC Design Principles

## 4.1 Use gRPC for Immediate Decisions

gRPC should be used when a caller needs a direct response before it can continue.

Examples:

```text
GetProduct
ListProducts
AddBasketItem
GetBasket
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
ValidateToken
```

These operations need clear success, failure, or validation responses.

---

## 4.2 Do Not Use gRPC for Downstream Notifications

gRPC should not be used for every downstream reaction after a business event.

For example, after an order is created, the Order Service should not synchronously call every interested downstream service.

Avoid:

```text
order-service -> notification-service
order-service -> recommendation-service
order-service -> search-service
order-service -> analytics-service
```

Prefer:

```text
order-service -> Kafka: OrderCreated
notification-service <- Kafka: OrderCreated
recommendation-service <- Kafka: OrderCreated
analytics-service <- Kafka: OrderCreated
```

This keeps the checkout path smaller and avoids unnecessary runtime coupling.

---

## 4.3 APIs Should Represent Service Behaviour

Protobuf messages should describe service behaviour, not database tables.

Good:

```text
CreateOrderRequest
ReserveStockRequest
AuthorisePaymentRequest
GetProductResponse
```

Avoid exposing internal database implementation details such as:

```text
OrderTableRow
InventoryDbRecord
PaymentEntity
```

The API is the service contract. The database is a private implementation detail.

---

## 4.4 Service APIs Are Owned by the Service

Each service owns its own protobuf API contract.

For example:

| Service | Owns API Contract |
|---|---|
| `catalog-service` | `proto/acme/catalog/v1/` |
| `inventory-service` | `proto/acme/inventory/v1/` |
| `basket-service` | `proto/acme/basket/v1/` |
| `order-service` | `proto/acme/order/v1/` |
| `payment-service` | `proto/acme/payment/v1/` |
| `shipping-service` | `proto/acme/shipping/v1/` |
| `notification-service` | `proto/acme/notification/v1/` |

Consumers may depend on the contract, but they do not own it.

---

## 4.5 Keep APIs Small and Purposeful

RPC methods should map to meaningful service behaviours.

Prefer:

```text
ReserveStock
ReleaseStockReservation
CommitStockReservation
```

Avoid vague generic methods such as:

```text
Process
Handle
Execute
UpdateData
DoAction
```

A method name should make the business operation clear.

---

## 4.6 Avoid Chatty APIs

Avoid designing flows where a simple request requires many sequential gRPC calls.

High-risk pattern:

```text
api-gateway
    -> catalog-service
    -> inventory-service
    -> review-service
    -> recommendation-service
    -> search-service
```

Better options:

- compose only what is needed for the user journey
- use read models where appropriate
- use Kafka events for downstream projections
- keep the critical path small

---

## 5. When to Use gRPC vs Kafka

## 5.1 Use gRPC When

Use gRPC when:

```text
the caller needs an immediate answer
the caller is making a command or query
the caller needs validation before continuing
the operation is part of a synchronous user journey
the result affects the next step in the workflow
```

Examples:

| Operation | Reason gRPC Is Appropriate |
|---|---|
| `GetProduct` | Client needs product details immediately |
| `AddBasketItem` | Client needs updated basket immediately |
| `ReserveStock` | Checkout cannot continue without stock decision |
| `AuthorisePayment` | Order cannot be confirmed without payment result |
| `CreateShipment` | Order flow needs shipment result or defined failure state |
| `ValidateToken` | Gateway needs identity validation before routing request |

---

## 5.2 Use Kafka When

Use Kafka when:

```text
something has already happened
multiple consumers may be interested
the producer should not block on downstream consumers
eventual consistency is acceptable
the event is useful for projections, notifications, search, recommendations, or audit
```

Examples:

| Event | Reason Kafka Is Appropriate |
|---|---|
| `OrderCreated` | Notifications, recommendations, and reporting may react later |
| `ProductUpdated` | Search index can update asynchronously |
| `ReviewApproved` | Rating summaries and search projections can update later |
| `ShipmentDispatched` | Notifications can be sent asynchronously |
| `PaymentFailed` | Order and notification workflows may react to the fact |

---

## 5.3 Rule of Thumb

```text
Need an answer now? Use gRPC.
Need to announce a fact? Use Kafka.
```

---

## 6. Protobuf Repository Layout

The target protobuf layout is:

```text
proto/
└── acme/
    ├── common/
    │   └── v1/
    ├── auth/
    │   └── v1/
    ├── customer/
    │   └── v1/
    ├── catalog/
    │   └── v1/
    ├── inventory/
    │   └── v1/
    ├── basket/
    │   └── v1/
    ├── order/
    │   └── v1/
    ├── payment/
    │   └── v1/
    ├── shipping/
    │   └── v1/
    ├── notification/
    │   └── v1/
    ├── review/
    │   └── v1/
    ├── search/
    │   └── v1/
    └── recommendation/
        └── v1/
```

Each service package should contain:

```text
<domain>_service.proto
<domain>_messages.proto
<domain>_events.proto
```

Example:

```text
proto/acme/order/v1/
├── order_service.proto
├── order_messages.proto
└── order_events.proto
```

---

## 7. Protobuf Package Naming

Use stable, versioned protobuf packages.

Example:

```proto
syntax = "proto3";

package acme.order.v1;

option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/order/v1;orderv1";
```

Package naming rules:

- use lowercase package names
- include the domain name
- include a major version
- avoid implementation-specific names
- do not include environment names

Good:

```text
acme.order.v1
acme.inventory.v1
acme.payment.v1
```

Avoid:

```text
acme.order_service_db.v1
acme.prod.order.v1
acme.internal.temp.v1
```

---

## 8. Service Definition Conventions

Service names should end with `Service`.

Example:

```proto
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc ListCustomerOrders(ListCustomerOrdersRequest) returns (ListCustomerOrdersResponse);
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
}
```

Rules:

- use clear method names
- use request and response messages for every method
- avoid returning raw domain entities without context
- avoid overly generic RPC methods
- avoid leaking database model names

---

## 9. Request and Response Conventions

## 9.1 Request Message Naming

Request messages should use the pattern:

```text
<MethodName>Request
```

Examples:

```text
CreateOrderRequest
ReserveStockRequest
AuthorisePaymentRequest
GetProductRequest
```

---

## 9.2 Response Message Naming

Response messages should use the pattern:

```text
<MethodName>Response
```

Examples:

```text
CreateOrderResponse
ReserveStockResponse
AuthorisePaymentResponse
GetProductResponse
```

---

## 9.3 Include IDs and Stable References

Requests should use stable IDs rather than database-specific identifiers.

Good:

```proto
message GetOrderRequest {
  string order_id = 1;
}
```

Avoid:

```proto
message GetOrderRequest {
  int64 order_table_primary_key = 1;
}
```

---

## 9.4 Use Explicit Fields

Prefer explicit fields over vague maps or unstructured payloads.

Good:

```proto
message AddBasketItemRequest {
  string basket_id = 1;
  string product_id = 2;
  string variant_id = 3;
  int32 quantity = 4;
}
```

Avoid:

```proto
message AddBasketItemRequest {
  map<string, string> data = 1;
}
```

Maps may be appropriate for metadata, but they should not replace well-defined contracts.

---

## 9.5 Use Common Metadata Carefully

Common request metadata may include:

```text
correlation_id
request_id
idempotency_key
actor_id
actor_type
```

Prefer transport metadata for cross-cutting request context where appropriate, and explicit message fields where the data is part of the business operation.

---

## 10. Common Protobuf Types

Shared common types should live under:

```text
proto/acme/common/v1/
```

Potential common types:

```text
Money
PaginationRequest
PaginationResponse
AddressSnapshot
AuditMetadata
TimeRange
```

Example:

```proto
message Money {
  string currency = 1;
  int64 amount_minor = 2;
}
```

Rules for common types:

- keep common types stable and minimal
- do not place service-specific business rules in common packages
- do not create a large shared domain model
- avoid making all services depend on unnecessary shared objects

---

## 11. Initial Service API Candidates

The following APIs are likely candidates for the initial checkout vertical slice.

## 11.1 Catalog Service

```proto
service CatalogService {
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc ValidateProductForPurchase(ValidateProductForPurchaseRequest) returns (ValidateProductForPurchaseResponse);
}
```

Purpose:

- browse product catalogue
- view product details
- validate product is active and purchasable

---

## 11.2 Basket Service

```proto
service BasketService {
  rpc CreateBasket(CreateBasketRequest) returns (CreateBasketResponse);
  rpc GetBasket(GetBasketRequest) returns (GetBasketResponse);
  rpc AddItem(AddItemRequest) returns (AddItemResponse);
  rpc UpdateItem(UpdateItemRequest) returns (UpdateItemResponse);
  rpc RemoveItem(RemoveItemRequest) returns (RemoveItemResponse);
  rpc MarkBasketCheckedOut(MarkBasketCheckedOutRequest) returns (MarkBasketCheckedOutResponse);
}
```

Purpose:

- manage basket state
- add, update, and remove basket items
- provide basket contents during checkout

---

## 11.3 Inventory Service

```proto
service InventoryService {
  rpc CheckAvailability(CheckAvailabilityRequest) returns (CheckAvailabilityResponse);
  rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);
  rpc ReleaseReservation(ReleaseReservationRequest) returns (ReleaseReservationResponse);
  rpc CommitReservation(CommitReservationRequest) returns (CommitReservationResponse);
}
```

Purpose:

- check stock availability
- reserve stock during checkout
- release stock after failed checkout
- commit stock after successful order confirmation

---

## 11.4 Order Service

```proto
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc ListCustomerOrders(ListCustomerOrdersRequest) returns (ListCustomerOrdersResponse);
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
}
```

Purpose:

- coordinate checkout
- own order lifecycle
- expose order history
- support cancellation later

---

## 11.5 Payment Service

```proto
service PaymentService {
  rpc AuthorisePayment(AuthorisePaymentRequest) returns (AuthorisePaymentResponse);
  rpc CapturePayment(CapturePaymentRequest) returns (CapturePaymentResponse);
  rpc RefundPayment(RefundPaymentRequest) returns (RefundPaymentResponse);
  rpc GetPayment(GetPaymentRequest) returns (GetPaymentResponse);
}
```

Purpose:

- simulate or perform payment authorisation
- record payment attempts
- support capture and refund later

---

## 11.6 Shipping Service

```proto
service ShippingService {
  rpc GetDeliveryOptions(GetDeliveryOptionsRequest) returns (GetDeliveryOptionsResponse);
  rpc CreateShipment(CreateShipmentRequest) returns (CreateShipmentResponse);
  rpc GetShipment(GetShipmentRequest) returns (GetShipmentResponse);
  rpc CancelShipment(CancelShipmentRequest) returns (CancelShipmentResponse);
}
```

Purpose:

- provide delivery options
- create shipment records
- expose fulfilment status

---

## 11.7 Notification Service

```proto
service NotificationService {
  rpc RequestNotification(RequestNotificationRequest) returns (RequestNotificationResponse);
  rpc GetNotificationStatus(GetNotificationStatusRequest) returns (GetNotificationStatusResponse);
}
```

Purpose:

- support direct notification requests where needed
- primarily consume Kafka events in the initial architecture

---

## 12. Checkout gRPC Flow

The initial checkout flow is expected to use gRPC for immediate decisions.

```text
API Gateway
    -> Order Service: CreateOrder
        -> Basket Service: GetBasket
        -> Inventory Service: ReserveStock
        -> Payment Service: AuthorisePayment
        -> Shipping Service: CreateShipment
    <- Order Service: CreateOrderResponse
```

After the order is created, asynchronous events are published.

```text
Order Service
    -> Kafka: OrderCreated

Notification Service
    <- Kafka: OrderCreated or NotificationRequested
```

The checkout gRPC path should remain focused on the steps required to determine whether an order can be created.

---

## 13. Idempotency

Idempotency is required for operations where duplicate requests could cause business harm.

High-priority idempotent operations:

```text
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
RequestNotification
CancelOrder
RefundPayment
```

## 13.1 Idempotency Key

Client-initiated or workflow-initiated requests should support an idempotency key where appropriate.

Example:

```proto
message CreateOrderRequest {
  string customer_id = 1;
  string basket_id = 2;
  string idempotency_key = 3;
  PaymentInput payment = 4;
  DeliveryInput delivery = 5;
}
```

## 13.2 Expected Behaviour

If the same idempotency key is received more than once for the same operation:

- the service should not create duplicate business records
- the service should return the original result where possible
- repeated requests should be safe
- conflicts should be explicit and observable

---

## 14. Timeouts

Every gRPC call must have a timeout.

Timeouts should be set by the caller based on the user journey and dependency criticality.

Example guidance:

| Call | Example Timeout Behaviour |
|---|---|
| API Gateway to Catalog Service | Short timeout suitable for user-facing read |
| Order Service to Inventory Service | Timeout must fail checkout safely |
| Order Service to Payment Service | Timeout must avoid duplicate unsafe payment attempts |
| Order Service to Shipping Service | Timeout behaviour depends on shipment/order design decision |
| API Gateway to Search Service | Timeout may return search unavailable or fallback response |

Rules:

- do not allow unbounded gRPC calls
- propagate context cancellation
- record timeout metrics
- ensure timeouts produce clear errors
- design compensation where timeout uncertainty can cause duplicated side effects

---

## 15. Retries

Retries must be used carefully.

Retries are safer for:

```text
read-only requests
idempotent write requests
transient network failures
unavailable responses where retry policy is explicitly defined
```

Retries are risky for:

```text
payment authorisation
stock reservation
shipment creation
order creation
```

For risky operations, retries require:

- idempotency keys
- clear timeout handling
- duplicate detection
- observable retry attempts
- bounded retry limits

Retry behaviour should be implemented in gRPC client middleware where appropriate, but service-specific business safety rules must still live in the owning service.

---

## 16. Error Handling

## 16.1 Use Standard gRPC Status Codes

Services should use standard gRPC status codes consistently.

| gRPC Code | Typical Use |
|---|---|
| `INVALID_ARGUMENT` | Request shape or validation error |
| `NOT_FOUND` | Requested entity does not exist |
| `ALREADY_EXISTS` | Duplicate entity or conflicting creation request |
| `FAILED_PRECONDITION` | Business precondition not satisfied |
| `PERMISSION_DENIED` | Authenticated caller lacks permission |
| `UNAUTHENTICATED` | Missing or invalid identity |
| `RESOURCE_EXHAUSTED` | Rate limit, quota, or capacity exceeded |
| `UNAVAILABLE` | Dependency or service temporarily unavailable |
| `DEADLINE_EXCEEDED` | Timeout exceeded |
| `INTERNAL` | Unexpected server error |

## 16.2 Business Errors

Business errors should be explicit and safe to expose through controlled responses.

Examples:

```text
PRODUCT_NOT_ACTIVE
INSUFFICIENT_STOCK
BASKET_EMPTY
PAYMENT_DECLINED
SHIPMENT_CREATION_FAILED
ORDER_ALREADY_EXISTS
```

These may be represented using structured error details.

## 16.3 Error Response Rules

- do not expose stack traces to clients
- do not expose internal database errors directly
- do not include secrets, tokens, or payment details in errors
- include correlation IDs in logs and traces
- map internal errors to safe API Gateway responses

---

## 17. Authentication and Authorisation

## 17.1 Authentication Context

The API Gateway should validate external client authentication where applicable.

Internal services should receive identity context through trusted metadata or service-to-service authentication mechanisms.

Potential metadata:

```text
authorization
x-correlation-id
x-request-id
x-actor-id
x-actor-type
```

## 17.2 Service-to-Service Trust

The target architecture should support service-to-service identity.

Potential approaches:

- mTLS through a service mesh later
- workload identity in Kubernetes
- signed internal tokens
- network policies and service accounts

The first local version may use simplified trust, but the production architecture should not rely on implicit network trust alone.

## 17.3 Authorisation Rules

Authorisation belongs to the service that owns the protected resource.

Examples:

- Order Service ensures customers can only access their own orders.
- Customer Service ensures users can only update their own profile.
- Review Service checks whether a customer can submit a review.
- Payment Service restricts sensitive payment operations.

The API Gateway may enforce broad access checks, but resource-level authorisation must remain with the owning service.

---

## 18. Observability Requirements

Every gRPC service must be observable by default.

## 18.1 Required Signals

Each service should emit:

```text
structured logs
request count
request latency
error count
timeout count
retry count
status code metrics
trace spans
correlation IDs
service name
method name
```

## 18.2 Tracing

Every gRPC call should create or propagate a trace span.

A checkout trace should show:

```text
api-gateway CreateOrder request
order-service CreateOrder
basket-service GetBasket
inventory-service ReserveStock
payment-service AuthorisePayment
shipping-service CreateShipment
Kafka OrderCreated publish
notification-service event consume
```

## 18.3 Logging

Logs should include:

```text
timestamp
level
service
method
correlation_id
request_id
trace_id
span_id
status
latency_ms
error_code where applicable
```

Logs must not include:

```text
raw card details
passwords
authentication tokens
full sensitive customer details
unredacted secrets
```

---

## 19. Health and Readiness

Each gRPC service should expose health and readiness information.

Recommended checks:

| Check | Purpose |
|---|---|
| Liveness | Is the process running? |
| Readiness | Can the service safely receive traffic? |
| Dependency readiness | Can required dependencies be reached? |
| gRPC health check | Allows Kubernetes and clients to check service health |

Readiness should fail when critical dependencies are unavailable.

Example:

- `order-service` readiness may depend on MySQL and required configuration.
- `inventory-service` readiness may depend on MySQL.
- `notification-service` readiness may depend on Kafka if it is primarily event-driven.

---

## 20. Contract Testing

gRPC contracts must be tested to reduce integration risk.

Contract testing should validate:

```text
protobuf compatibility
required fields and validation behaviour
expected status codes
business error mapping
backwards compatibility
consumer expectations
```

Recommended test types:

| Test Type | Purpose |
|---|---|
| Protobuf linting | Enforce style and consistency |
| Breaking change checks | Prevent incompatible contract changes |
| Server contract tests | Verify service implements expected behaviour |
| Client integration tests | Verify clients call APIs correctly |
| End-to-end tests | Verify full journey across services |

Buf should be used for protobuf linting and breaking-change detection.

---

## 21. Versioning and Compatibility

## 21.1 Package Versioning

Protobuf packages should include a major version.

Example:

```text
acme.order.v1
```

A breaking API change should require a new version.

Example:

```text
acme.order.v2
```

## 21.2 Backwards-Compatible Changes

Generally safe changes:

- adding optional fields
- adding new RPC methods
- adding new enum values with care
- adding new messages

Risky or breaking changes:

- removing fields
- renaming fields
- changing field numbers
- changing field types
- changing RPC request or response types
- changing semantics of existing fields
- reusing old field numbers

## 21.3 Reserved Fields

When removing protobuf fields, reserve field numbers and names.

Example:

```proto
message Order {
  reserved 8;
  reserved "legacy_status";
}
```

This prevents accidental reuse and unsafe compatibility issues.

---

## 22. API Documentation Expectations

Each service API should document:

```text
purpose
RPC methods
request fields
response fields
business validation rules
error cases
timeout expectations
idempotency requirements
authorisation requirements
observability requirements
example requests
example responses
```

Service-specific API docs may live in:

```text
docs/api/<service-name>.md
services/<service-name>/docs/api.md
```

---

## 23. Security Considerations

gRPC APIs must be designed with security in mind.

Rules:

- authenticate callers where required
- enforce resource-level authorisation in owning services
- do not trust caller-supplied user IDs without validated identity context
- do not expose internal error details
- do not log sensitive payloads
- validate all request fields
- apply least privilege to service accounts
- use TLS or mTLS in production-like environments
- propagate audit metadata for sensitive operations

Sensitive operations include:

```text
CreateOrder
AuthorisePayment
RefundPayment
CancelOrder
UpdateCustomerProfile
AddAddress
SubmitReview
```

---

## 24. Performance Considerations

gRPC provides efficient binary communication, but good API design is still required.

Guidelines:

- avoid unnecessary large responses
- paginate list APIs
- avoid excessive service fan-out
- use deadlines on every call
- use streaming only where there is a clear need
- measure p95 and p99 latency for important methods
- avoid putting non-critical services in critical request paths

Important methods for performance monitoring:

```text
ListProducts
GetProduct
AddBasketItem
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
SearchProducts
```

---

## 25. Local Development

Local development should support gRPC service execution through Docker Compose and local processes.

Expected developer tasks:

```sh
make proto
make test
make dev-up
make run
make dev-down
```

`make proto` should generate service stubs from protobuf definitions.

The local environment should allow developers to:

- run services locally
- call gRPC APIs
- inspect logs
- run contract tests
- run integration tests
- test the checkout flow end to end

Useful local tools may include:

```text
grpcurl
grpcui
buf
Docker Compose
```

---

## 26. CI/CD Expectations

The CI pipeline should validate gRPC contracts before merging changes.

Recommended CI checks:

```text
protobuf formatting
protobuf linting
protobuf breaking-change detection
generated code freshness
unit tests
contract tests
integration tests
security scanning
container build
```

Protobuf changes should be reviewed carefully because they affect multiple services.

---

## 27. Example Protobuf Style

Example service definition:

```proto
syntax = "proto3";

package acme.order.v1;

option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/order/v1;orderv1";

import "acme/common/v1/money.proto";

service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
}

message CreateOrderRequest {
  string customer_id = 1;
  string basket_id = 2;
  string idempotency_key = 3;
  PaymentInput payment = 4;
  DeliveryInput delivery = 5;
}

message CreateOrderResponse {
  string order_id = 1;
  string order_number = 2;
  OrderStatus status = 3;
  acme.common.v1.Money total = 4;
}

message GetOrderRequest {
  string order_id = 1;
}

message GetOrderResponse {
  Order order = 1;
}

message PaymentInput {
  string payment_method_token = 1;
}

message DeliveryInput {
  string address_id = 1;
  string delivery_option_id = 2;
}

message Order {
  string order_id = 1;
  string order_number = 2;
  OrderStatus status = 3;
  repeated OrderItem items = 4;
  acme.common.v1.Money total = 5;
}

message OrderItem {
  string product_id = 1;
  string product_name_snapshot = 2;
  int32 quantity = 3;
  acme.common.v1.Money unit_price = 4;
}

enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING = 1;
  ORDER_STATUS_CONFIRMED = 2;
  ORDER_STATUS_FAILED = 3;
  ORDER_STATUS_CANCELLED = 4;
}
```

This example is illustrative and should be refined as the service contracts are implemented.

---

## 28. Anti-Patterns to Avoid

## 28.1 Database-Shaped APIs

Avoid APIs that simply expose database rows.

Poor:

```text
GetOrderTableRow
UpdateOrderTableColumn
```

Better:

```text
GetOrder
CancelOrder
CreateOrder
```

---

## 28.2 Generic RPC Methods

Avoid generic methods that hide meaning.

Poor:

```text
ProcessRequest
HandleAction
ExecuteCommand
```

Better:

```text
ReserveStock
AuthorisePayment
CreateShipment
```

---

## 28.3 No Deadlines

Every call must have a deadline or timeout.

Unbounded calls cause poor failure behaviour and can exhaust resources.

---

## 28.4 Business Logic in API Gateway

The API Gateway should not own business workflows.

The gateway may route and map requests, but order creation belongs to Order Service, stock reservation belongs to Inventory Service, and payment state belongs to Payment Service.

---

## 28.5 Kafka Used as Synchronous RPC

Do not publish an event and wait for another service to process it as if it were an immediate command.

Use gRPC when the caller needs an immediate result.

---

## 29. Open Questions

| Question | Status |
|---|---|
| Will the API Gateway expose REST, GraphQL, or gRPC-Web externally? | To decide |
| Should all internal APIs use gRPC from the first version? | Proposed |
| Should payment authorisation and capture be separate RPCs in the first version? | To decide |
| Should shipment creation be synchronous during checkout or event-driven after order creation? | To decide |
| Should customer profile lookups be synchronous during checkout or based on address snapshots? | To decide |
| Which gRPC retry policies should be standardised in shared middleware? | To decide |
| Should service-to-service identity use mTLS via service mesh later? | Deferred |

---

## 30. Related Documents

This document should be read alongside:

```text
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/architecture/event-driven-design.md
docs/architecture/resilience-patterns.md
docs/api/protobuf-style-guide.md
docs/api/error-model.md
docs/api/authentication.md
docs/api/versioning.md
docs/events/event-catalog.md
docs/testing/contract-testing.md
docs/observability/tracing.md
```

Relevant ADRs:

```text
adr/0002-use-grpc-for-service-communication.md
adr/0006-use-buf-for-protobuf.md
adr/0008-use-contract-first-service-design.md
```

---

## 31. Summary

gRPC is the standard synchronous communication mechanism for internal bfstore services.

It should be used for commands and queries that require an immediate response, such as:

```text
GetProduct
AddBasketItem
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
```

Kafka should be used for facts that have already happened, such as:

```text
OrderCreated
PaymentAuthorised
StockReserved
ShipmentCreated
NotificationSent
```

The gRPC API model should remain:

```text
contract-first
service-owned
versioned
observable
secure
idempotent where necessary
tested through CI
safe to evolve
```

This approach gives bfstore a professional, maintainable service communication model suitable for a senior platform engineering portfolio.
