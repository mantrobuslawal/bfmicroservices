# `shipping/v1`

## 1. Purpose

This document describes the `shipping/v1` Protobuf package for bfstore.

The package defines delivery option, shipment creation, shipment state, tracking reference, and shipping event contracts.

These contracts are client-visible engineering artefacts and should be treated as stable integration boundaries.

---

## 2. Package Name

```proto
package acme.shipping.v1;
```

Recommended Go package option:

```proto
option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/shipping/v1;shippingv1";
```

---

## 3. Ownership

| Area | Owner |
|---|---|
| Package owner | `shipping-service` |
| Contract review | API governance / service owner |
| Backward compatibility | Service owner |
| Generated clients | Build/tooling pipeline |

Ownership rule:

> The service that owns the business capability owns the contract. Other services integrate through APIs and events, not database tables.

---

## 4. Expected Files

```text
proto/acme/shipping/v1/
├── README.md
├── shipping_service.proto
├── shipment.proto
├── delivery_option.proto
├── shipping_events.proto
```

Event contracts may later move to a dedicated package such as:

```text
proto/acme/<domain>/events/v1/
```

if that separation improves clarity.

---

## 5. Primary Service Contract

```proto
service ShippingService {
  rpc ListDeliveryOptions(... ) returns (...);
  rpc CreateShipment(... ) returns (...);
  rpc GetShipment(... ) returns (...);
  rpc CancelShipment(... ) returns (...);
}
```

---

## 6. Core Message Types

Recommended message types:

- `DeliveryOption`
- `Shipment`
- `ShipmentStatus`
- `TrackingEvent`

These messages should describe business concepts and API contracts, not database rows.

---

## Recommended Status Enum

```proto
enum ShipmentStatus {
  SHIPMENT_STATUS_UNSPECIFIED = 0;
  SHIPMENT_STATUS_CREATED = 1;
  SHIPMENT_STATUS_PENDING_FULFILMENT = 2;
  SHIPMENT_STATUS_DISPATCHED = 3;
  SHIPMENT_STATUS_DELIVERED = 4;
  SHIPMENT_STATUS_CANCELLED = 5;
  SHIPMENT_STATUS_FAILED = 6;
}
```

---

## 7. Contract Rules

- Shipping Service owns shipment state.
- CreateShipment must be idempotent.
- Delivery address snapshots must be minimised and protected.
- Shipment failure behaviour during checkout must be documented.
- Shipping failure must be observable.

---

## 8. Event Contracts

Expected or potential events:

- `ShipmentCreated`
- `ShipmentFailed`
- `ShipmentDispatched`
- `ShipmentDelivered`
- `ShipmentCancelled`

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

- INVALID_ARGUMENT for invalid delivery option or address
- NOT_FOUND for missing shipment
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
package acme.shipping.v2;
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

The `acme.shipping.v1` package is part of bfstore's contract-first service design.

It should remain business-focused, versioned, testable, and aligned with the owning service boundary.
