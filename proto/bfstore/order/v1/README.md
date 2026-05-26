# `order/v1`

## 1. Purpose

This document describes the `order/v1` Protobuf package for bfstore.

The package defines order lifecycle contracts and initial checkout orchestration contracts.

These contracts are client-visible engineering artefacts and should be treated as stable integration boundaries.

---

## 2. Package Name

```proto
package acme.order.v1;
```

Recommended Go package option:

```proto
option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/order/v1;orderv1";
```

---

## 3. Ownership

| Area | Owner |
|---|---|
| Package owner | `order-service` |
| Contract review | API governance / service owner |
| Backward compatibility | Service owner |
| Generated clients | Build/tooling pipeline |

Ownership rule:

> The service that owns the business capability owns the contract. Other services integrate through APIs and events, not database tables.

---

## 4. Expected Files

```text
proto/acme/order/v1/
├── README.md
├── order_service.proto
├── order.proto
├── checkout.proto
├── order_events.proto
```

Event contracts may later move to a dedicated package such as:

```text
proto/acme/<domain>/events/v1/
```

if that separation improves clarity.

---

## 5. Primary Service Contract

```proto
service OrderService {
  rpc CreateOrder(... ) returns (...);
  rpc GetOrder(... ) returns (...);
  rpc ListOrders(... ) returns (...);
  rpc CancelOrder(... ) returns (...);
}
```

---

## 6. Core Message Types

Recommended message types:

- `Order`
- `OrderItem`
- `OrderStatus`
- `CheckoutAttempt`
- `CreateOrderRequest`
- `CreateOrderResponse`

These messages should describe business concepts and API contracts, not database rows.

---

## Recommended Status Enum

```proto
enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING = 1;
  ORDER_STATUS_CONFIRMED = 2;
  ORDER_STATUS_CANCELLED = 3;
  ORDER_STATUS_FAILED = 4;
}
```

---

## 7. Contract Rules

- Order Service owns order lifecycle.
- Order Service coordinates initial checkout orchestration.
- Stock, payment, shipment, and notification internals remain owned by their services.
- CreateOrder must be idempotent.
- Order items must preserve product and price snapshots.

---

## 8. Event Contracts

Expected or potential events:

- `OrderCreated`
- `OrderFailed`
- `OrderCancelled`
- `OrderConfirmed`

Event rules:

```text
events describe facts that have already happened
events must include the standard event envelope
events must be versioned
events must carry correlation context
consumers must be idempotent
```

---

## 9. Error Behaviour

Expected error behaviour:

- INVALID_ARGUMENT for missing request fields
- FAILED_PRECONDITION for empty basket, inactive product, insufficient stock, or payment decline
- DEADLINE_EXCEEDED for downstream timeout
- UNAVAILABLE for downstream outage
- NOT_FOUND for missing order
- OK for duplicate idempotent retry with same request

All gRPC errors should follow:

```text
docs/api/error-model.md
```

---

## 10. Versioning

This package is versioned as `v1`.

Compatible changes include:

```text
adding optional fields
adding new messages
adding new RPCs
adding new event types
adding comments
```

Breaking changes include:

```text
renaming services or RPCs
renaming packages
removing fields
renumbering fields
changing field types
changing field meaning
changing idempotency behaviour
changing error semantics
```

Breaking changes require a new package version, such as:

```proto
package acme.order.v2;
```

Removed fields must be reserved.

---

## 11. Security and Privacy

Contracts must avoid unnecessary sensitive data.

Do not expose:

```text
passwords
tokens
raw payment card data
secret values
internal stack traces
internal database IDs
unnecessary customer PII
```

Prefer opaque IDs and service-owned lookups where sensitive details are required.

---

## 12. Testing Expectations

This package should be covered by:

```text
buf lint
buf breaking
protobuf generation checks
gRPC contract tests
event contract tests where events are defined
integration tests for critical service behaviours
```

---

## 13. Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/error-model.md
docs/api/versioning.md
docs/architecture/communication-patterns.md
docs/architecture/service-boundaries.md
docs/events/event-catalog.md
docs/testing/testing-strategy.md
```

---

## 14. Summary

The `acme.order.v1` package is part of bfstore's contract-first service design.

It should remain business-focused, versioned, testable, and aligned with the owning service boundary.
