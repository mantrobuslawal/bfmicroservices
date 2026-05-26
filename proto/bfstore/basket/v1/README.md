# `basket/v1`

## 1. Purpose

This document describes the `basket/v1` Protobuf package for bfstore.

The package defines shopping basket contracts for current shopping intent before checkout.

These contracts are client-visible engineering artefacts and should be treated as stable integration boundaries.

---

## 2. Package Name

```proto
package acme.basket.v1;
```

Recommended Go package option:

```proto
option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/basket/v1;basketv1";
```

---

## 3. Ownership

| Area | Owner |
|---|---|
| Package owner | `basket-service` |
| Contract review | API governance / service owner |
| Backward compatibility | Service owner |
| Generated clients | Build/tooling pipeline |

Ownership rule:

> The service that owns the business capability owns the contract. Other services integrate through APIs and events, not database tables.

---

## 4. Expected Files

```text
proto/acme/basket/v1/
├── README.md
├── basket_service.proto
├── basket.proto
├── basket_events.proto
```

Event contracts may later move to a dedicated package such as:

```text
proto/acme/<domain>/events/v1/
```

if that separation improves clarity.

---

## 5. Primary Service Contract

```proto
service BasketService {
  rpc CreateBasket(... ) returns (...);
  rpc GetBasket(... ) returns (...);
  rpc AddBasketItem(... ) returns (...);
  rpc UpdateBasketItem(... ) returns (...);
  rpc RemoveBasketItem(... ) returns (...);
  rpc MarkBasketCheckedOut(... ) returns (...);
}
```

---

## 6. Core Message Types

Recommended message types:

- `Basket`
- `BasketItem`
- `BasketStatus`

These messages should describe business concepts and API contracts, not database rows.

---

## Recommended Status Enum

```proto
enum BasketStatus {
  BASKET_STATUS_UNSPECIFIED = 0;
  BASKET_STATUS_ACTIVE = 1;
  BASKET_STATUS_CHECKED_OUT = 2;
  BASKET_STATUS_EXPIRED = 3;
  BASKET_STATUS_ABANDONED = 4;
}
```

---

## 7. Contract Rules

- Basket Service owns baskets and basket items.
- Basket Service does not reserve stock.
- Only active baskets may be changed.
- Quantities must be greater than zero.
- Basket prices are not final order prices unless explicitly snapshotted.

---

## 8. Event Contracts

Expected or potential events:

- `BasketCreated`
- `BasketItemAdded`
- `BasketItemUpdated`
- `BasketItemRemoved`
- `BasketCheckedOut`
- `BasketExpired`

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

- NOT_FOUND for missing basket
- FAILED_PRECONDITION for checked-out basket or inactive product
- INVALID_ARGUMENT for invalid quantity
- UNAVAILABLE for database outage

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
package acme.basket.v2;
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

The `acme.basket.v1` package is part of bfstore's contract-first service design.

It should remain business-focused, versioned, testable, and aligned with the owning service boundary.
