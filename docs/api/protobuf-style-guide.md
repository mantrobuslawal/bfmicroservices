# Protobuf Style Guide

## 1. Purpose

This document defines the Protobuf style guide for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It establishes conventions for writing `.proto` files that are clear, consistent, maintainable, versioned, and suitable for professional service-to-service contracts.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s API design discipline.

---

## 2. Scope

This guide applies to:

```text
gRPC service definitions
request messages
response messages
domain value messages
event payload messages
shared common protobuf types
error-related protobuf types
protobuf package naming
field naming
versioning conventions
Buf linting and breaking-change checks
```

It does not define every service API in detail. Individual service APIs should live in:

```text
proto/acme/
```

and be documented through:

```text
docs/api/grpc-overview.md
docs/api/error-model.md
docs/api/versioning.md
```

---

## 3. Design Goals

bfstore protobuf contracts should be:

| Goal | Description |
|---|---|
| Clear | Easy to understand without knowing service internals |
| Stable | Safe to evolve without breaking consumers unexpectedly |
| Consistent | Same naming and structure across services |
| Explicit | Business meaning is visible in names and fields |
| Typed | Avoid vague stringly typed contracts |
| Versioned | Breaking changes are handled deliberately |
| Observable | Requests support tracing and correlation |
| Secure | Contracts avoid leaking sensitive information |
| Testable | Contracts support generated clients and contract tests |

---

## 4. Protobuf Principles

## 4.1 Contract First

Protobuf files define service contracts.

The contract should be designed before implementation details are finalised.

A service should not expose internal database models directly through protobuf.

Good:

```text
CreateOrderRequest
OrderSummary
PaymentAuthorisationResult
StockReservation
```

Poor:

```text
OrderTableRow
PaymentDbRecord
InventorySqlResult
```

---

## 4.2 Business Language Over Technical Leakage

Use names that reflect the business domain.

Good:

```text
ReserveStock
AuthorisePayment
CreateShipment
OrderCreated
```

Poor:

```text
InsertReservation
UpdatePaymentTable
CreateShipmentRow
```

---

## 4.3 Prefer Explicit Messages

Avoid overusing generic messages.

Good:

```proto
message CreateOrderRequest {
  string basket_id = 1;
  string customer_id = 2;
  string idempotency_key = 3;
}
```

Poor:

```proto
message GenericRequest {
  map<string, string> values = 1;
}
```

---

## 4.4 Avoid Exposing Persistence Models

Database tables and protobuf messages may look similar, but they serve different purposes.

Persistence models optimise storage.

Protobuf messages define integration contracts.

Keep them separate.

---

## 5. Directory Layout

Recommended protobuf layout:

```text
proto/
└── acme/
    ├── common/
    │   └── v1/
    │       ├── money.proto
    │       ├── pagination.proto
    │       ├── errors.proto
    │       └── metadata.proto
    ├── catalog/
    │   └── v1/
    │       └── catalog_service.proto
    ├── inventory/
    │   └── v1/
    │       └── inventory_service.proto
    ├── basket/
    │   └── v1/
    │       └── basket_service.proto
    ├── order/
    │   └── v1/
    │       └── order_service.proto
    ├── payment/
    │   └── v1/
    │       └── payment_service.proto
    ├── shipping/
    │   └── v1/
    │       └── shipping_service.proto
    ├── notification/
    │   └── v1/
    │       └── notification_service.proto
    ├── review/
    │   └── v1/
    │       └── review_service.proto
    ├── search/
    │   └── v1/
    │       └── search_service.proto
    └── recommendation/
        └── v1/
            └── recommendation_service.proto
```

Events may either live under the relevant service package or a separate event package.

Example:

```text
proto/acme/order/events/v1/order_events.proto
proto/acme/payment/events/v1/payment_events.proto
```

---

## 6. Package Naming

Use stable package names that include company, domain, and version.

Format:

```text
acme.<domain>.v1
```

Examples:

```proto
package acme.catalog.v1;
package acme.inventory.v1;
package acme.basket.v1;
package acme.order.v1;
package acme.payment.v1;
package acme.shipping.v1;
```

For events:

```proto
package acme.order.events.v1;
package acme.payment.events.v1;
```

---

## 7. Go Package Option

Each `.proto` file should define `go_package`.

Example:

```proto
option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/order/v1;orderv1";
```

Naming convention:

```text
<domain>v<version>
```

Examples:

```text
catalogv1
inventoryv1
basketv1
orderv1
paymentv1
shippingv1
notificationv1
```

---

## 8. File Naming

Use lower snake case for `.proto` files.

Good:

```text
catalog_service.proto
inventory_service.proto
order_events.proto
payment_service.proto
common_metadata.proto
```

Avoid:

```text
CatalogService.proto
catalogService.proto
service.proto
messages.proto
```

File names should make ownership and purpose obvious.

---

## 9. Service Naming

Use PascalCase for service names.

Format:

```text
<Domain>Service
```

Examples:

```proto
service CatalogService {}
service InventoryService {}
service BasketService {}
service OrderService {}
service PaymentService {}
service ShippingService {}
service NotificationService {}
```

Avoid generic names:

```proto
service Service {}
service Api {}
service Handler {}
```

---

## 10. RPC Naming

RPC names should use PascalCase and describe the business operation.

Good:

```proto
rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
rpc GetProduct(GetProductRequest) returns (GetProductResponse);
rpc AddBasketItem(AddBasketItemRequest) returns (AddBasketItemResponse);
rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);
rpc AuthorisePayment(AuthorisePaymentRequest) returns (AuthorisePaymentResponse);
rpc CreateShipment(CreateShipmentRequest) returns (CreateShipmentResponse);
```

Avoid implementation-focused names:

```proto
rpc InsertProduct(...);
rpc UpdateDb(...);
rpc Process(...);
```

---

## 11. Message Naming

Use PascalCase for message names.

Request and response messages should follow this pattern:

```text
<RpcName>Request
<RpcName>Response
```

Examples:

```proto
message GetProductRequest {}
message GetProductResponse {}

message ReserveStockRequest {}
message ReserveStockResponse {}
```

Domain messages should be clear and reusable within the same package where appropriate:

```proto
message Product {}
message Basket {}
message Order {}
message Payment {}
message Shipment {}
```

---

## 12. Field Naming

Use lower snake case for field names.

Good:

```proto
string product_id = 1;
string customer_id = 2;
int32 quantity = 3;
string idempotency_key = 4;
```

Avoid:

```proto
string productID = 1;
string CustomerId = 2;
int32 qty = 3;
```

Field names should be descriptive.

Prefer:

```proto
reserved_quantity
available_quantity
unit_price
delivery_address
```

over:

```proto
qty
val
data
info
```

---

## 13. Field Numbering

Field numbers are part of the wire contract.

Rules:

- Do not reuse field numbers.
- Do not renumber existing fields.
- Reserve removed field numbers.
- Group related fields logically.
- Leave room for future expansion where sensible.

Example:

```proto
message Product {
  string product_id = 1;
  string name = 2;
  string description = 3;
  Money price = 4;
  ProductStatus status = 5;

  reserved 6;
  reserved "legacy_sku";
}
```

---

## 14. Reserved Fields

When removing a field, reserve the field number and name.

Example:

```proto
message Customer {
  string customer_id = 1;
  string display_name = 2;

  reserved 3;
  reserved "legacy_email";
}
```

This prevents accidental reuse and unsafe compatibility issues.

---

## 15. Enum Naming

Use PascalCase for enum types.

Use UPPER_SNAKE_CASE for enum values.

The first enum value must be an unspecified value.

Example:

```proto
enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING = 1;
  ORDER_STATUS_CONFIRMED = 2;
  ORDER_STATUS_CANCELLED = 3;
  ORDER_STATUS_FAILED = 4;
}
```

Rules:

- Prefix enum values with the enum name.
- Include an `_UNSPECIFIED = 0` value.
- Do not rely on default zero having business meaning.
- Consumers should handle unknown enum values safely.

---

## 16. Timestamps

Use `google.protobuf.Timestamp` for absolute times.

Example:

```proto
import "google/protobuf/timestamp.proto";

message Order {
  string order_id = 1;
  google.protobuf.Timestamp created_at = 2;
}
```

Use UTC for stored and transmitted timestamps unless there is a strong reason otherwise.

Do not use strings for timestamps unless representing a human-entered value that is not an instant in time.

---

## 17. Money

Use a shared `Money` type.

Example:

```proto
message Money {
  int64 amount_minor = 1;
  string currency_code = 2;
}
```

Rules:

- Store monetary amounts in minor units.
- Use ISO currency codes such as `GBP`.
- Avoid floating point types for money.

Good:

```proto
Money total_amount = 1;
```

Avoid:

```proto
double total_price = 1;
```

---

## 18. IDs

IDs should be represented as strings unless there is a strong reason otherwise.

Examples:

```proto
string product_id = 1;
string customer_id = 2;
string order_id = 3;
string payment_id = 4;
```

Benefits:

- supports UUID, ULID, KSUID, or external IDs
- avoids leaking database implementation
- easier to use across systems

Avoid exposing auto-increment database IDs as public contracts.

---

## 19. Pagination

List APIs should use consistent pagination.

Example:

```proto
message PageRequest {
  int32 page_size = 1;
  string page_token = 2;
}

message PageResponse {
  string next_page_token = 1;
  int32 total_count = 2;
}
```

Example usage:

```proto
message ListProductsRequest {
  acme.common.v1.PageRequest page = 1;
}

message ListProductsResponse {
  repeated Product products = 1;
  acme.common.v1.PageResponse page = 2;
}
```

Rules:

- Avoid unbounded list responses.
- Use tokens for cursor-style pagination where possible.
- Document maximum page sizes.

---

## 20. Repeated Fields

Use repeated fields for lists.

Example:

```proto
repeated BasketItem items = 1;
```

Rules:

- Avoid returning extremely large repeated fields.
- Use pagination where result size can grow.
- Preserve meaningful ordering only when documented.

---

## 21. Maps

Use maps only when keys and values are genuinely dynamic.

Acceptable:

```proto
map<string, string> attributes = 1;
```

Avoid maps for structured business data that deserves explicit fields.

Poor:

```proto
map<string, string> order = 1;
```

Prefer:

```proto
message Order {
  string order_id = 1;
  repeated OrderItem items = 2;
}
```

---

## 22. Optional Fields

Proto3 supports `optional`.

Use optional fields where presence matters.

Example:

```proto
optional string discount_code = 1;
```

Avoid using optional for every field by default.

Use it deliberately when the difference between “not supplied” and “empty value” matters.

---

## 23. Oneof

Use `oneof` where only one of several fields may be set.

Example:

```proto
message PaymentMethod {
  oneof method {
    CardPaymentToken card_token = 1;
    WalletPaymentToken wallet_token = 2;
  }
}
```

Rules:

- Use `oneof` for mutually exclusive choices.
- Avoid large, confusing `oneof` structures.
- Document behaviour clearly.

---

## 24. Request Metadata

Common request metadata may include:

```proto
message RequestMetadata {
  string correlation_id = 1;
  string request_id = 2;
  string idempotency_key = 3;
}
```

However, some values may be better passed through gRPC metadata rather than every request body.

Recommended approach:

| Value | Recommended Location |
|---|---|
| `correlation_id` | gRPC metadata |
| `trace_id` | OpenTelemetry context |
| `idempotency_key` | Request body or metadata, depending on operation |
| `customer_id` | Request body only where part of the business request |
| auth token | gRPC metadata |

Be consistent across services.

---

## 25. Error Handling

Do not model every error as a successful response.

Prefer standard gRPC error status for failed requests.

Good:

```text
ReserveStock returns FAILED_PRECONDITION for insufficient stock
GetProduct returns NOT_FOUND for missing product
CreateOrder returns ALREADY_EXISTS or equivalent conflict handling for duplicate idempotency where appropriate
```

Response messages should not usually include fields like:

```proto
bool success = 1;
string error = 2;
```

Detailed error handling is defined in:

```text
docs/api/error-model.md
```

---

## 26. Events in Protobuf

Kafka event payloads should also use protobuf contracts.

Event payload messages should describe facts.

Good:

```proto
message OrderCreatedEvent {
  string order_id = 1;
  string customer_id = 2;
  Money total_amount = 3;
  google.protobuf.Timestamp created_at = 4;
}
```

Avoid:

```proto
message CreateOrderEvent {}
```

Event envelope design should be documented separately:

```text
docs/events/event-envelope.md
```

---

## 27. Sensitive Data Rules

Protobuf contracts must avoid exposing sensitive data unnecessarily.

Do not include:

```text
passwords
password hashes
authentication tokens
raw card data
secret values
internal stack traces
```

Be careful with:

```text
customer email
phone number
delivery address
payment provider references
payment failure details
```

Prefer IDs where the receiving service can retrieve sensitive data through authorised APIs if required.

---

## 28. Comments and Documentation

Public protobuf contracts should include comments.

Example:

```proto
// CatalogService provides read access to active product catalogue data.
service CatalogService {
  // GetProduct returns product details for a single product.
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
}
```

Comments should explain:

```text
business purpose
important behaviour
field meaning where not obvious
idempotency requirements
consistency expectations
```

Avoid comments that only repeat the field name.

Poor:

```proto
// product id
string product_id = 1;
```

Better:

```proto
// Unique product identifier assigned by Catalog Service.
string product_id = 1;
```

---

## 29. Compatibility Rules

Compatible changes:

```text
add optional fields
add new RPCs
add new messages
add new enum values if consumers handle unknowns
add comments
```

Breaking changes:

```text
rename package
rename service
rename RPC
remove field
renumber field
change field type
change field meaning
remove enum value
change request semantics
change response semantics
```

Breaking changes require a versioning plan.

---

## 30. Buf Standards

bfstore should use Buf for protobuf quality checks.

Expected files:

```text
buf.yaml
buf.gen.yaml
```

Recommended checks:

```text
buf lint
buf breaking
buf generate
```

CI should fail on:

```text
lint violations
breaking changes against configured baseline
generation errors
```

---

## 31. Example Service Definition

Example only:

```proto
syntax = "proto3";

package acme.inventory.v1;

option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/inventory/v1;inventoryv1";

import "google/protobuf/timestamp.proto";
import "acme/common/v1/money.proto";

service InventoryService {
  // ReserveStock reserves stock for an order attempt.
  //
  // This operation must be idempotent when the same idempotency key is supplied.
  rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);

  // ReleaseStockReservation releases a previous stock reservation.
  rpc ReleaseStockReservation(ReleaseStockReservationRequest) returns (ReleaseStockReservationResponse);
}

message ReserveStockRequest {
  string order_id = 1;
  string basket_id = 2;
  string idempotency_key = 3;
  repeated StockReservationItem items = 4;
}

message ReserveStockResponse {
  string reservation_id = 1;
  ReservationStatus status = 2;
  google.protobuf.Timestamp expires_at = 3;
}

message StockReservationItem {
  string product_id = 1;
  string variant_id = 2;
  int32 quantity = 3;
}

message ReleaseStockReservationRequest {
  string reservation_id = 1;
  string reason = 2;
  string idempotency_key = 3;
}

message ReleaseStockReservationResponse {
  string reservation_id = 1;
  ReservationStatus status = 2;
}

enum ReservationStatus {
  RESERVATION_STATUS_UNSPECIFIED = 0;
  RESERVATION_STATUS_RESERVED = 1;
  RESERVATION_STATUS_RELEASED = 2;
  RESERVATION_STATUS_EXPIRED = 3;
  RESERVATION_STATUS_FAILED = 4;
}
```

---

## 32. Anti-Patterns to Avoid

Avoid:

```text
generic request and response messages
database table messages
unversioned packages
unbounded list responses
floating point money values
field renumbering
field number reuse
ambiguous enum zero values
sensitive data in contracts
business logic hidden in comments only
```

---

## 33. Review Checklist

Before approving a protobuf change, check:

```text
Is the package versioned?
Is the service name clear?
Are RPC names business-focused?
Are request and response messages explicit?
Are field names descriptive?
Are field numbers stable?
Are removed fields reserved?
Are enums prefixed and do they include UNSPECIFIED = 0?
Is money represented safely?
Are timestamps represented with google.protobuf.Timestamp?
Is sensitive data avoided?
Are comments useful?
Does buf lint pass?
Does buf breaking pass?
Are contract tests updated?
```

---

## 34. Related Documents

This document should be read alongside:

```text
docs/api/grpc-overview.md
docs/api/error-model.md
docs/api/versioning.md
docs/events/event-envelope.md
docs/events/event-catalog.md
docs/architecture/communication-patterns.md
docs/architecture/service-boundaries.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0002-use-grpc-for-service-communication.md
adr/0006-use-buf-for-protobuf.md
adr/0008-use-contract-first-service-design.md
```

---

## 35. Summary

bfstore protobuf contracts should be explicit, versioned, business-focused, and safe to evolve.

The most important rules are:

```text
contracts are not database models
packages must be versioned
field numbers must never be reused
removed fields must be reserved
money must not use floating point types
events describe facts
errors should use the standard gRPC error model
Buf should enforce linting and breaking-change checks
```

This style guide helps keep bfstore’s APIs professional, maintainable, and credible for a client-facing portfolio project.
