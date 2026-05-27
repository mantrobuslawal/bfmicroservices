# ADR-0003: Use Kafka for Domain Events

## Status

Accepted

## Date

2026-05-26

## Context

bfstore needs asynchronous communication between services for facts that have already happened.

Examples:

```text
OrderCreated
StockReserved
PaymentAuthorised
ShipmentCreated
NotificationSent
ProductUpdated
```

Synchronous gRPC calls are appropriate for commands and queries requiring an immediate response. Events are appropriate for decoupled follow-up work, projections, notifications, recommendations, analytics, and operational workflows.

bfstore will use Protocol Buffers for Kafka event payloads to align event contracts with gRPC contracts.

---

## Decision

bfstore will use Kafka for domain events.

Kafka event payloads will be encoded using Protocol Buffers.

```text
Kafka key     = stable business identifier
Kafka value   = Protobuf event message
Kafka headers = lightweight metadata
```

Recommended content type:

```text
application/x-protobuf
```

---

## Alternatives Considered

## Synchronous gRPC Only

Simpler, but too tightly couples asynchronous side effects and projections.

## Kafka with JSON Events

Readable and easy to inspect, but weaker type safety and less consistent with gRPC Protobuf contracts.

## Kafka with Protobuf Events

Typed, generated, versioned, and compatible with Buf checks. Less human-readable, but stronger for contract-first service design.

---

## Consequences

Positive:

```text
services can publish facts without knowing all consumers
Search Service can build product projections
Notification Service can react to OrderCreated
Recommendation Service can consume behavioural signals later
events have strong typed contracts
schema evolution can be tested in CI
```

Negative:

```text
Kafka adds operational complexity
consumers must be idempotent
Protobuf payloads require decoding tools
DLQ handling must preserve binary payloads safely
versioning discipline is mandatory
```

---

## Implementation Notes

Initial topics:

```text
bfstore.order.orders.v1
bfstore.inventory.stock.v1
bfstore.payment.payments.v1
bfstore.shipping.shipments.v1
bfstore.notification.notifications.v1
```

Recommended proto layout:

```text
proto/acme/events/v1/event_metadata.proto
proto/acme/order/events/v1/order_events.proto
proto/acme/payment/events/v1/payment_events.proto
proto/acme/inventory/events/v1/inventory_events.proto
proto/acme/shipping/events/v1/shipping_events.proto
proto/acme/notification/events/v1/notification_events.proto
proto/acme/catalog/events/v1/catalog_events.proto
```

---

## Summary

bfstore uses Kafka for domain events and Protobuf for event payloads. This gives the system a consistent contract-first model across both gRPC APIs and asynchronous event-driven communication.


---

## Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/events/event-envelope.md
docs/events/event-catalog.md
docs/events/kafka-topic-design.md
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
```
