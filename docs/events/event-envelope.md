# Event Envelope

## 1. Purpose

This document defines the standard Kafka event envelope for **bfstore**.

bfstore uses **Protocol Buffers** as the primary serialisation format for Kafka event payloads. JSON may be used for documentation examples, local debugging, or inspection tooling, but it is not the production event payload contract.

---

## 2. Event Message Standard

```text
Kafka message key   = stable business identifier
Kafka message value = Protobuf event message
Kafka headers       = lightweight routing/tracing metadata
```

Recommended Kafka header:

```text
content_type = application/x-protobuf
```

---

## 3. Why Protobuf for Events?

Using Protobuf for both gRPC and Kafka gives bfstore a single contract-first model.

Benefits:

```text
typed event contracts
generated producer and consumer code
smaller payloads than verbose JSON
clear schema evolution rules
Buf linting and breaking-change checks
better contract testing
less ambiguity between services
```

Trade-offs:

```text
events are less human-readable than JSON
debugging requires generated decoders or tooling
schema discipline is mandatory
field numbers must never be reused
removed fields must be reserved
```

---

## 4. Recommended Proto Layout

```text
proto/acme/events/v1/
├── event_metadata.proto
└── event_envelope.proto

proto/acme/catalog/events/v1/catalog_events.proto
proto/acme/inventory/events/v1/inventory_events.proto
proto/acme/basket/events/v1/basket_events.proto
proto/acme/order/events/v1/order_events.proto
proto/acme/payment/events/v1/payment_events.proto
proto/acme/shipping/events/v1/shipping_events.proto
proto/acme/notification/events/v1/notification_events.proto
```

Shared metadata belongs in `acme.events.v1`.

Domain event payloads belong in domain-specific event packages.

---

## 5. Recommended Metadata Message

```proto
message EventMetadata {
  string event_id = 1;
  string event_type = 2;
  string event_version = 3;
  google.protobuf.Timestamp occurred_at = 4;
  string producer = 5;
  string subject = 6;
  string correlation_id = 7;
  string causation_id = 8;
  string trace_id = 9;
  string idempotency_key = 10;
}
```

`subject` should identify the primary business entity, for example `order_id`, `product_id`, `payment_id`, `shipment_id`, or `reservation_id`.

---

## 6. Typed Event Wrapper Pattern

Recommended pattern:

```proto
message OrderCreatedEvent {
  acme.events.v1.EventMetadata metadata = 1;
  OrderCreated payload = 2;
}

message OrderCreated {
  string order_id = 1;
  string order_number = 2;
  string customer_id = 3;
  acme.common.v1.Money total_amount = 4;
  repeated OrderItemSnapshot items = 5;
  google.protobuf.Timestamp created_at = 6;
}
```

This is preferred over raw `bytes`, arbitrary JSON strings, or undocumented payload maps.

---

## 7. Kafka Headers

Recommended headers:

```text
event_type
event_version
correlation_id
traceparent
producer
content_type
```

Headers are useful for routing, tracing, and operational diagnostics. The Protobuf value remains the source of the event contract.

---

## 8. Security and Privacy

Events must not include:

```text
raw payment card data
CVV
passwords
access tokens
provider secrets
unnecessary customer PII
large unbounded payloads
internal database implementation details
```

---

## 9. Testing Expectations

Event envelope tests should include:

```text
Buf lint
Buf breaking checks
producer serialisation tests
consumer deserialisation tests
contract tests for each event type
unsupported event version handling
invalid binary payload handling
DLQ behaviour for decode failures
```

---

## 10. Summary

bfstore Kafka messages use Protobuf event payloads with standard metadata. This gives the project typed, versioned, contract-first asynchronous communication that aligns with its gRPC API design.


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
