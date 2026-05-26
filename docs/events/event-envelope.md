# Event Envelope

## 1. Purpose

This document defines the standard event envelope for **bfstore**, ACME Ltd’s fictional online furniture store backend.

The event envelope provides a consistent structure for Kafka events so that all services can publish, consume, trace, version, validate, and troubleshoot events in a predictable way.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s event-driven architecture.

---

## 2. Scope

This document applies to all Kafka events produced and consumed by bfstore services, including:

```text
catalogue events
inventory events
basket events
order events
payment events
shipping events
notification events
review events
search projection events
recommendation events
```

It covers:

```text
event metadata
event identity
event type
event version
correlation and causation IDs
trace propagation
producer ownership
event timestamps
payload structure
security and privacy expectations
observability requirements
testing expectations
```

The event catalogue itself is documented in:

```text
docs/events/event-catalog.md
```

---

## 3. Design Goals

The bfstore event envelope should be:

| Goal | Description |
|---|---|
| Consistent | All services use the same metadata structure |
| Traceable | Events can be linked to business flows and distributed traces |
| Versioned | Consumers can identify event schema versions |
| Auditable | Events include identity, timing, and producer information |
| Idempotency-friendly | Consumers can detect duplicate events |
| Secure | Sensitive data is avoided or minimised |
| Operationally useful | Logs, metrics, and DLQs can reference common fields |
| Contract-friendly | Event payloads can be validated against protobuf schemas |

---

## 4. Envelope Principle

The event envelope separates event metadata from event data.

```text
Envelope = information about the event
Data     = the business payload of the event
```

Example:

```json
{
  "event_id": "evt_01HX9H9YF4S7Z9J7J8J8A1B2C3",
  "event_type": "OrderCreated",
  "event_version": "1.0",
  "occurred_at": "2026-05-26T10:15:30Z",
  "producer": "order-service",
  "correlation_id": "corr_01HX9H9YF4S7Z9J7J8J8A1B2C3",
  "causation_id": "cmd_01HX9H9YF4S7Z9J7J8J8A1B2C3",
  "trace_id": "trace_01HX9H9YF4S7Z9J7J8J8A1B2C3",
  "data": {
    "order_id": "ord_123",
    "customer_id": "cus_456",
    "total_amount": {
      "amount_minor": 129900,
      "currency_code": "GBP"
    }
  }
}
```

---

## 5. Required Envelope Fields

| Field | Required | Description |
|---|---:|---|
| `event_id` | Yes | Globally unique identifier for this event instance |
| `event_type` | Yes | Stable business event name |
| `event_version` | Yes | Version of the event schema or payload |
| `occurred_at` | Yes | UTC timestamp when the business fact occurred |
| `producer` | Yes | Service that produced the event |
| `correlation_id` | Yes | Identifier shared across a business flow |
| `data` | Yes | Event-specific business payload |

---

## 6. Recommended Envelope Fields

| Field | Required | Description |
|---|---:|---|
| `causation_id` | Recommended | Identifier for the command or event that caused this event |
| `trace_id` | Recommended | Distributed tracing identifier |
| `span_id` | Optional | Span identifier where tracing context is propagated |
| `subject` | Optional | Primary business entity affected by the event |
| `partition_key` | Optional | Logical partitioning key used for Kafka |
| `schema_name` | Optional | Protobuf schema or message name |
| `schema_version` | Optional | Schema version when separate from event version |

---

## 7. Field Definitions

## 7.1 `event_id`

A unique identifier for a single event instance.

Rules:

```text
must be globally unique
must not be reused
must be stable during retries
must be logged by producers and consumers
should be used for consumer deduplication
```

Example:

```text
evt_01HX9H9YF4S7Z9J7J8J8A1B2C3
```

Recommended formats:

```text
ULID
UUIDv7
UUIDv4
```

---

## 7.2 `event_type`

A stable business event name.

Good examples:

```text
OrderCreated
PaymentAuthorised
StockReserved
ShipmentCreated
ProductUpdated
NotificationFailed
```

Avoid command-style names:

```text
CreateOrder
AuthorisePayment
ReserveStock
SendNotification
```

Events describe facts that have already happened.

---

## 7.3 `event_version`

The version of the event contract.

Example:

```text
1.0
```

Rules:

```text
major version changes indicate incompatible changes
minor version changes indicate compatible additions
consumers must check the version they can process
unsupported versions should be handled explicitly
```

Event versioning should align with:

```text
protobuf package version
Kafka topic version
event payload compatibility
```

---

## 7.4 `occurred_at`

The timestamp when the business event occurred.

Rules:

```text
use UTC
use ISO 8601 if represented as JSON
use google.protobuf.Timestamp if represented in protobuf
represent business occurrence time, not consumer processing time
```

Example:

```text
2026-05-26T10:15:30Z
```

---

## 7.5 `producer`

The service that produced the event.

Examples:

```text
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
review-service
search-service
recommendation-service
```

The producer must be the service that owns the business fact.

---

## 7.6 `correlation_id`

The identifier shared across a business flow.

Example checkout flow:

```text
API Gateway receives checkout request
    -> Order Service creates checkout attempt
    -> Inventory Service reserves stock
    -> Payment Service authorises payment
    -> Shipping Service creates shipment
    -> Order Service publishes OrderCreated
    -> Notification Service sends confirmation
```

All related logs, traces, gRPC calls, and Kafka events should carry the same correlation ID.

---

## 7.7 `causation_id`

The identifier for what caused this event.

Examples:

```text
command ID
request ID
previous event ID
checkout attempt ID
```

Causation helps answer:

```text
Why was this event produced?
Which command or event caused it?
```

---

## 7.8 `trace_id`

The distributed trace identifier.

This allows event publishing and consumption to be linked to application traces.

Trace context should be propagated through:

```text
gRPC metadata
Kafka headers
event envelope
structured logs
OpenTelemetry spans
```

---

## 7.9 `subject`

The primary business entity affected by the event.

Examples:

```text
order:ord_123
product:prd_456
payment:pay_789
shipment:shp_321
```

This is optional but useful for filtering, auditing, and troubleshooting.

---

## 8. Envelope and Payload Structure

The envelope should contain metadata and one event-specific payload.

Conceptual structure:

```text
EventEnvelope
    metadata fields
    data payload
```

For protobuf, this may be modelled using a common envelope message and typed payload messages, or by carrying envelope metadata in Kafka headers with protobuf payloads.

The final implementation choice should be documented once the Kafka/protobuf integration is selected.

---

## 9. Kafka Headers vs Payload Envelope

There are two common approaches.

## 9.1 Metadata in Kafka Headers

Example headers:

```text
event_id
event_type
event_version
correlation_id
causation_id
trace_id
producer
occurred_at
```

Payload contains only typed business data.

Benefits:

```text
easier routing and filtering
metadata available without deserialising full payload
common in event platforms
```

Costs:

```text
metadata and payload may become separated
headers need consistent tooling support
```

## 9.2 Metadata in Payload Envelope

Payload contains both metadata and data.

Benefits:

```text
self-contained event
easier replay and storage
simpler for some consumers
```

Costs:

```text
consumer must deserialise payload to read metadata
envelope repeated across every event schema
```

## 9.3 Initial Recommendation

For bfstore, use a consistent logical envelope and decide implementation detail later.

A practical approach is:

```text
Kafka headers contain routing and tracing metadata.
Protobuf payload contains required metadata and typed event data.
```

This provides operational flexibility while keeping replayed events self-contained.

---

## 10. Protobuf Representation

Example conceptual protobuf structure:

```proto
syntax = "proto3";

package acme.events.v1;

import "google/protobuf/timestamp.proto";

message EventMetadata {
  string event_id = 1;
  string event_type = 2;
  string event_version = 3;
  google.protobuf.Timestamp occurred_at = 4;
  string producer = 5;
  string correlation_id = 6;
  string causation_id = 7;
  string trace_id = 8;
  string subject = 9;
}
```

Example event payload:

```proto
syntax = "proto3";

package acme.order.events.v1;

import "acme/events/v1/event_metadata.proto";
import "acme/common/v1/money.proto";

message OrderCreatedEvent {
  acme.events.v1.EventMetadata metadata = 1;
  string order_id = 2;
  string customer_id = 3;
  acme.common.v1.Money total_amount = 4;
  repeated OrderCreatedItem items = 5;
}
```

---

## 11. Event Payload Rules

Event payloads should include enough information for consumers to react safely.

Payloads should avoid unnecessary data duplication.

Good payload fields:

```text
business entity ID
customer ID where needed
safe snapshot values
state transition result
amounts in minor units
timestamps
reason codes
```

Avoid:

```text
raw card data
passwords
tokens
large unbounded objects
internal database row dumps
stack traces
secret values
```

---

## 12. Event Identity and Idempotency

Consumers should use `event_id` for deduplication where duplicate processing may be harmful.

High-risk duplicates:

```text
sending customer notifications twice
capturing or refunding payment twice
releasing stock twice
creating duplicate projections with side effects
```

Recommended deduplication key:

```text
producer + event_type + event_id
```

Where business semantics matter, consumers may also track:

```text
order_id
payment_id
reservation_id
shipment_id
notification_id
```

---

## 13. Event Time Semantics

There are multiple timestamps in event-driven systems.

| Timestamp | Meaning |
|---|---|
| `occurred_at` | When the business fact happened |
| `published_at` | When the producer published the event |
| `consumed_at` | When a consumer processed the event |
| `created_at` | Entity creation time inside the business payload |

The envelope must include `occurred_at`.

`published_at` may be added by producer instrumentation or outbox publishing.

Consumers should not treat processing time as business occurrence time.

---

## 14. Correlation Example

Example failed checkout:

```text
correlation_id: corr_123
```

Related events and logs:

```text
CheckoutStarted
StockReserved
PaymentFailed
StockReservationReleased
OrderFailed
NotificationRequested
NotificationSent
```

Correlation allows the full business flow to be reconstructed.

---

## 15. Security and Privacy

The envelope itself should not contain sensitive data.

Avoid placing personal or sensitive information in metadata fields.

Do not use:

```text
customer email as correlation_id
phone number as subject
payment token as causation_id
secret name as producer detail
```

Prefer opaque IDs:

```text
customer_id
order_id
payment_id
shipment_id
correlation_id
```

---

## 16. Observability Requirements

Producers should log:

```text
event_id
event_type
event_version
producer
topic
partition where available
correlation_id
publish result
publish latency
```

Consumers should log:

```text
event_id
event_type
event_version
consumer
topic
partition
offset
correlation_id
processing result
processing latency
retry count
DLQ status
```

Metrics should include:

```text
events_published_total
event_publish_failures_total
events_consumed_total
event_processing_failures_total
event_processing_duration_seconds
consumer_lag
dlq_messages_total
duplicate_events_total
```

---

## 17. Testing Requirements

Event envelope tests should verify:

```text
required metadata is present
event_id is unique
event_type matches documented event
event_version is set
occurred_at is set in UTC
producer matches owning service
correlation_id is propagated
sensitive data is not present
consumer handles duplicate event_id safely
unsupported event_version is handled
```

Contract tests should fail if required envelope fields are missing.

---

## 18. Anti-Patterns to Avoid

Avoid:

```text
events without event_id
events without version
events without producer
using event names as commands
putting secrets in metadata
putting raw database rows in payloads
using processing time as occurred_at
generating a new event_id on every retry of the same event
assuming correlation_id is optional
```

---

## 19. Initial Implementation Scope

The first version should enforce envelope standards for:

```text
OrderCreated
OrderFailed
StockReserved
StockReservationFailed
PaymentAuthorised
PaymentFailed
ShipmentCreated
ShipmentFailed
NotificationSent
NotificationFailed
```

These events support the checkout vertical slice.

---

## 20. Open Questions

| Question | Status |
|---|---|
| Will envelope metadata be duplicated in Kafka headers and protobuf payload? | To decide |
| Which ID format will be standard for event IDs? | To decide |
| Will `published_at` be included in the envelope or only in logs/metrics? | To decide |
| Will a schema registry be used alongside protobuf and Buf? | To decide |
| Will `subject` be required for all events? | To decide |
| How will trace context be represented in Kafka headers? | To decide |

---

## 21. Related Documents

This document should be read alongside:

```text
docs/events/event-catalog.md
docs/events/kafka-topic-design.md
docs/events/ordering-and-idempotency.md
docs/events/retry-and-dlq-strategy.md
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/architecture/event-driven-design.md
docs/architecture/communication-patterns.md
docs/testing/testing-strategy.md
```

---

## 22. Summary

The bfstore event envelope gives every event a consistent identity, version, producer, timestamp, correlation context, and payload structure.

The most important rules are:

```text
every event must have a unique event_id
every event must have an event_type
every event must have an event_version
every event must have an occurred_at timestamp
every event must identify its producer
every event must carry correlation context
event payloads must avoid unnecessary sensitive data
```

A consistent event envelope makes bfstore’s event-driven architecture easier to test, observe, operate, and evolve.
