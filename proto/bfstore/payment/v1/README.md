# `payment/v1`

## 1. Purpose

This document describes the `payment/v1` Protobuf package for bfstore.

The package defines payment state, attempts, authorisation, provider references, refunds where implemented, and payment event contracts.

These contracts are client-visible engineering artefacts and should be treated as stable integration boundaries.

---

## 2. Package Name

```proto
package acme.payment.v1;
```

Recommended Go package option:

```proto
option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/payment/v1;paymentv1";
```

---

## 3. Ownership

| Area | Owner |
|---|---|
| Package owner | `payment-service` |
| Contract review | API governance / service owner |
| Backward compatibility | Service owner |
| Generated clients | Build/tooling pipeline |

Ownership rule:

> The service that owns the business capability owns the contract. Other services integrate through APIs and events, not database tables.

---

## 4. Expected Files

```text
proto/acme/payment/v1/
├── README.md
├── payment_service.proto
├── payment.proto
├── payment_events.proto
```

Event contracts may later move to a dedicated package such as:

```text
proto/acme/<domain>/events/v1/
```

if that separation improves clarity.

---

## 5. Primary Service Contract

```proto
service PaymentService {
  rpc AuthorisePayment(... ) returns (...);
  rpc GetPayment(... ) returns (...);
  rpc CapturePayment (later)(... ) returns (...);
  rpc RefundPayment (later)(... ) returns (...);
}
```

---

## 6. Core Message Types

Recommended message types:

- `Payment`
- `PaymentAttempt`
- `PaymentStatus`
- `PaymentMethod`
- `PaymentAuthorisationResult`

These messages should describe business concepts and API contracts, not database rows.

---

## Recommended Status Enum

```proto
enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED = 0;
  PAYMENT_STATUS_PENDING = 1;
  PAYMENT_STATUS_AUTHORISED = 2;
  PAYMENT_STATUS_DECLINED = 3;
  PAYMENT_STATUS_FAILED = 4;
  PAYMENT_STATUS_CAPTURED = 5;
  PAYMENT_STATUS_REFUNDED = 6;
  PAYMENT_STATUS_CANCELLED = 7;
}
```

---

## 7. Contract Rules

- Payment Service owns payment state.
- Raw card data must not be stored or logged.
- AuthorisePayment must be idempotent.
- Payment decline must not create a confirmed order.
- Provider timeouts must be reconciled safely.

---

## 8. Event Contracts

Expected or potential events:

- `PaymentAuthorised`
- `PaymentFailed`
- `PaymentCaptured`
- `PaymentRefunded`
- `PaymentCancelled`

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

- INVALID_ARGUMENT for invalid amount, currency, or missing idempotency key
- FAILED_PRECONDITION for declined payment
- DEADLINE_EXCEEDED for provider timeout
- UNAVAILABLE for provider outage
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
package acme.payment.v2;
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

The `acme.payment.v1` package is part of bfstore's contract-first service design.

It should remain business-focused, versioned, testable, and aligned with the owning service boundary.
