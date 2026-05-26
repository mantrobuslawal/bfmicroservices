# Event-Driven Design

## 1. Purpose

This document defines the event-driven design approach for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how Kafka is used, what events mean in this system, how event ownership works, how consumers should behave, and how event-driven workflows are made reliable, observable, testable, and safe.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s asynchronous architecture.

---

## 2. Architecture Context

bfstore uses a hybrid communication model:

| Pattern | Technology | Purpose |
|---|---|---|
| Synchronous APIs | gRPC | Commands and queries requiring immediate response |
| Asynchronous messaging | Kafka | Business facts and downstream reactions |
| Schema definition | Protobuf | API and event contracts |
| Persistence | MySQL | Service-owned data |
| Observability | OpenTelemetry | Logs, metrics, traces, correlation |

Event-driven design is used where services need to react to business facts without tightly coupling the producer to every consumer.

Example:

```text
Order Service creates an order
    -> publishes OrderCreated
        -> Notification Service sends confirmation
        -> Recommendation Service updates signals
        -> Analytics consumer records sale
```

Order Service should not need to know every downstream consumer.

---

## 3. Event-Driven Principles

## 3.1 Events Are Facts

Events describe something that has already happened.

Good event names:

```text
OrderCreated
PaymentAuthorised
StockReserved
ShipmentCreated
NotificationFailed
ProductUpdated
```

Poor event names:

```text
CreateOrder
AuthorisePayment
ReserveStock
SendNotification
UpdateProduct
```

Commands request work. Events describe completed facts.

---

## 3.2 Event Producers Own the Facts

The service that owns the business fact publishes the event.

| Event | Producer |
|---|---|
| `ProductUpdated` | `catalog-service` |
| `StockReserved` | `inventory-service` |
| `BasketCheckedOut` | `basket-service` |
| `OrderCreated` | `order-service` |
| `PaymentAuthorised` | `payment-service` |
| `ShipmentCreated` | `shipping-service` |
| `NotificationSent` | `notification-service` |
| `ReviewApproved` | `review-service` |

A service should not publish an event for a business fact it does not own.

---

## 3.3 Consumers Are Independent

A producer should not need to know all consumers.

For example, Order Service may publish `OrderCreated`.

Consumers may include:

```text
notification-service
recommendation-service
analytics consumer
audit projection
```

Adding a new consumer should not require changing Order Service unless the event contract itself needs to evolve.

---

## 3.4 Consumers Must Be Idempotent

Kafka consumers must safely handle duplicate delivery and replay.

Examples:

- Notification Service must avoid duplicate customer messages where practical.
- Search Service must handle duplicate `ProductUpdated` events.
- Inventory Service must handle duplicate stock release events safely.
- Order Service must avoid duplicate state transitions.

---

## 3.5 Eventual Consistency Must Be Explicit

Event-driven workflows often introduce temporary inconsistency.

Acceptable eventual consistency examples:

```text
search index updates after product changes
recommendations after orders
review summaries after review approval
notification delivery after order creation
```

Eventual consistency is acceptable only when:

- the business can tolerate it
- the behaviour is documented
- consumer lag is observable
- replay or repair is possible

---

## 4. Where bfstore Uses Events

bfstore uses events for:

| Use Case | Reason |
|---|---|
| Notifications | Order creation should not block on email/SMS sending |
| Search indexing | Product changes should update search projections asynchronously |
| Recommendations | Behavioural signals can be processed later |
| Review summaries | Rating projections can update asynchronously |
| Operational projections | Dashboards and audit projections can consume events |
| Decoupled side effects | Producers avoid knowing every downstream workflow |

---

## 5. Where bfstore Should Not Use Events

Events should not be used when the caller needs an immediate result.

Use gRPC instead for:

```text
GetProduct
AddBasketItem
ReserveStock
AuthorisePayment
CreateShipment
CreateOrder
```

Avoid using Kafka as a hidden request/response mechanism.

If a service publishes an event and waits for a consumer to respond before continuing, the design may be using Kafka as RPC.

---

## 6. Event Envelope

All events should use a consistent envelope.

Example conceptual envelope:

```json
{
  "event_id": "evt_01HX...",
  "event_type": "OrderCreated",
  "event_version": "1.0",
  "occurred_at": "2026-05-26T10:15:30Z",
  "producer": "order-service",
  "correlation_id": "corr_01HX...",
  "causation_id": "cmd_01HX...",
  "trace_id": "trace_01HX...",
  "data": {}
}
```

## 6.1 Required Metadata

| Field | Required | Purpose |
|---|---:|---|
| `event_id` | Yes | Deduplication and audit |
| `event_type` | Yes | Identifies business event |
| `event_version` | Yes | Schema compatibility |
| `occurred_at` | Yes | Business event time |
| `producer` | Yes | Publishing service |
| `correlation_id` | Yes | End-to-end business flow tracing |
| `causation_id` | Recommended | What caused this event |
| `trace_id` | Recommended | Distributed tracing |
| `data` | Yes | Event payload |

Detailed envelope design should be documented in:

```text
docs/events/event-envelope.md
```

---

## 7. Topic Design

Recommended topic naming pattern:

```text
bfstore.<domain>.<stream>.v<version>
```

Examples:

```text
bfstore.catalog.product-events.v1
bfstore.inventory.stock-events.v1
bfstore.basket.basket-events.v1
bfstore.order.order-events.v1
bfstore.payment.payment-events.v1
bfstore.shipping.shipment-events.v1
bfstore.notification.notification-events.v1
bfstore.review.review-events.v1
bfstore.search.search-events.v1
bfstore.recommendation.recommendation-events.v1
```

## 7.1 Topic Design Rules

- Keep topics aligned with domain ownership.
- Avoid one topic for every tiny event unless justified.
- Avoid one giant topic for all events.
- Use versioned topics for incompatible changes.
- Choose partition keys based on ordering needs.

Detailed topic guidance should be documented in:

```text
docs/events/kafka-topic-design.md
```

---

## 8. Partitioning and Ordering

Kafka preserves ordering only within a partition.

Recommended partition keys:

| Event Stream | Suggested Key |
|---|---|
| Product events | `product_id` |
| Stock events | `product_id` or `reservation_id` |
| Basket events | `basket_id` |
| Order events | `order_id` |
| Payment events | `payment_id` or `order_id` |
| Shipment events | `shipment_id` or `order_id` |
| Notification events | `notification_id` |
| Review events | `review_id` or `product_id` |

## 8.1 Ordering Rules

- Do not assume global ordering across all events.
- Keep events for the same aggregate on the same partition where ordering matters.
- Consumers should handle missing or out-of-order events where practical.
- Projections should track offsets and event versions.

---

## 9. Idempotency

Consumers must be idempotent.

## 9.1 Why Idempotency Matters

Kafka consumers may reprocess events because of:

```text
consumer restart
offset reset
retry
rebalance
manual replay
producer duplicate
network failure
```

## 9.2 Deduplication Inputs

Consumers may deduplicate using:

```text
event_id
event_type
producer
business_entity_id
event_version
```

## 9.3 High-Risk Idempotency Areas

| Consumer | Risk |
|---|---|
| `notification-service` | Duplicate customer messages |
| `payment-service` | Duplicate refund or capture behaviour |
| `inventory-service` | Duplicate stock release |
| `order-service` | Duplicate order transitions |
| `search-service` | Projection corruption |
| `recommendation-service` | Duplicate signal weighting |

---

## 10. Retry and Dead-Letter Design

Consumers should distinguish retryable and non-retryable failures.

## 10.1 Retryable Failures

Examples:

```text
temporary database outage
temporary network issue
temporary downstream service failure
temporary provider outage
```

Handling:

```text
retry with backoff
record retry count
emit metrics
alert if persistent
```

## 10.2 Non-Retryable Failures

Examples:

```text
invalid event payload
unsupported event version
missing required identifier
schema mismatch
poison message
```

Handling:

```text
send to DLQ
record failure reason
alert operations if needed
provide replay or repair process
```

## 10.3 DLQ Naming

Example DLQ names:

```text
bfstore.order.order-events.dlq.v1
bfstore.payment.payment-events.dlq.v1
bfstore.notification.notification-events.dlq.v1
bfstore.search.search-events.dlq.v1
```

DLQ handling should be documented in runbooks.

---

## 11. Outbox Pattern

## 11.1 Problem

A service may update its database successfully but fail to publish the matching event.

Example:

```text
Order is saved in MySQL
OrderCreated event fails to publish
Notification Service never sees the order
```

This creates data/event inconsistency.

## 11.2 Proposed Pattern

Use an outbox table owned by the service.

Example:

```text
order-service transaction:
    insert order
    insert order_items
    insert outbox_event OrderCreated
```

A separate publisher then reads the outbox and publishes to Kafka.

## 11.3 Recommended Use

For a serious implementation, consider the outbox pattern for:

```text
order-service
payment-service
inventory-service
shipping-service
notification-service
catalog-service
review-service
```

## 11.4 Trade-Off

The outbox pattern improves reliability but adds operational and implementation complexity.

This decision should be captured in an ADR if implemented.

---

## 12. Event Replay

Event replay may be required for:

```text
rebuilding search indexes
repairing recommendation projections
reprocessing failed notifications
rebuilding reporting projections
testing consumer idempotency
```

## 12.1 Replay Rules

- Consumers must be idempotent.
- Replayed events must preserve original event metadata.
- Replays should be observable.
- Replays should be controlled and documented.
- Replays should not cause duplicate harmful side effects.

## 12.2 Side-Effect Caution

Not all consumers are safe to replay blindly.

High caution:

```text
notification-service
payment-service
shipping-service
```

Lower risk:

```text
search-service
recommendation-service
analytics projections
```

---

## 13. Event Versioning

Events are contracts and must be versioned.

## 13.1 Backward-Compatible Changes

Usually safe:

```text
adding optional fields
adding new event types
adding new enum values if consumers handle unknowns
```

## 13.2 Breaking Changes

Usually breaking:

```text
removing fields
renaming fields
changing field meaning
changing field type
changing required behaviour
renaming event types
changing topic semantics
```

Breaking changes should use a new version and migration plan.

---

## 14. Event Security and Privacy

Events must not leak sensitive data.

Do not include:

```text
passwords
authentication tokens
raw card data
secret values
full stack traces
unnecessary personal data
```

Be careful with:

```text
customer email
phone number
delivery address
payment failure details
provider references
```

Prefer IDs and allow authorised services to retrieve sensitive details through APIs when required.

---

## 15. Observability

Event-driven systems must be highly observable.

## 15.1 Producer Metrics

Producers should expose:

```text
events published
publish failures
publish latency
outbox pending count
outbox publish lag
event type count
```

## 15.2 Consumer Metrics

Consumers should expose:

```text
events consumed
processing success count
processing failure count
retry count
DLQ count
consumer lag
processing latency
duplicate event count
```

## 15.3 Logs

Event logs should include:

```text
event_id
event_type
event_version
topic
partition where available
offset where available
correlation_id
producer
consumer
```

## 15.4 Traces

Event consumption should link to the originating trace where practical.

This allows a business flow such as checkout to be traced across synchronous and asynchronous boundaries.

---

## 16. Testing Requirements

Event-driven workflows require dedicated tests.

## 16.1 Producer Tests

Validate:

```text
event emitted when business action succeeds
event has correct envelope
event payload matches schema
event uses correct topic
event includes correlation ID
```

## 16.2 Consumer Tests

Validate:

```text
consumer handles valid event
consumer ignores duplicate event safely
consumer handles unsupported version
consumer sends invalid event to DLQ
consumer updates projection correctly
```

## 16.3 Integration Tests

Validate:

```text
Kafka publish and consume flow
consumer group behaviour
retry behaviour
DLQ behaviour
outbox publisher behaviour
```

## 16.4 End-to-End Tests

Critical event E2E tests:

```text
OrderCreated triggers notification
ProductUpdated updates search projection
PaymentFailed prevents confirmed order
Duplicate OrderCreated does not duplicate notification
```

---

## 17. Initial Event-Driven Scope

The first implementation should focus on checkout-related events.

Initial events:

```text
StockReserved
StockReservationFailed
PaymentAuthorised
PaymentFailed
ShipmentCreated
ShipmentFailed
OrderCreated
OrderFailed
NotificationSent
NotificationFailed
```

Deferred events:

```text
ProductViewed
ReviewCreated
ReviewApproved
SearchIndexUpdated
RecommendationGenerated
BasketAbandoned
```

This staged approach proves event-driven architecture without overloading the first implementation.

---

## 18. Event-Driven Checkout Flow

Example:

```text
1. Customer submits checkout.
2. Order Service retrieves basket.
3. Order Service requests stock reservation.
4. Inventory Service reserves stock and publishes StockReserved.
5. Order Service requests payment authorisation.
6. Payment Service authorises payment and publishes PaymentAuthorised.
7. Order Service requests shipment creation.
8. Shipping Service creates shipment and publishes ShipmentCreated.
9. Order Service creates order and publishes OrderCreated.
10. Notification Service consumes OrderCreated.
11. Notification Service sends or simulates confirmation.
12. Notification Service publishes NotificationSent.
```

## 18.1 Failure Events

Failure examples:

```text
StockReservationFailed
PaymentFailed
ShipmentFailed
OrderFailed
NotificationFailed
```

Rules:

- If stock reservation fails, payment should not be attempted.
- If payment fails, order should not be confirmed.
- If notification fails, order should remain created.
- If shipment creation fails, behaviour must follow the documented order/fulfilment decision.

---

## 19. Anti-Patterns to Avoid

## 19.1 Event Soup

Avoid publishing many poorly defined events without ownership or purpose.

Every event should have:

```text
business meaning
owning producer
documented consumers
schema
version
topic
observability
test coverage
```

## 19.2 Kafka as RPC

Avoid using Kafka when the caller needs an immediate result.

## 19.3 Hidden Coupling

Avoid consumers depending on undocumented event fields or producer implementation details.

## 19.4 Unbounded Event Payloads

Avoid large payloads containing unnecessary data.

Prefer meaningful IDs and business fields.

## 19.5 Non-Idempotent Consumers

Avoid consumers that create duplicate side effects when events are redelivered.

---

## 20. Operational Considerations

Event-driven operations should include:

```text
topic dashboards
consumer lag alerts
DLQ alerts
outbox lag alerts
replay runbooks
schema compatibility checks
consumer failure dashboards
event throughput metrics
```

Runbooks should cover:

```text
consumer lag
DLQ growth
outbox backlog
poison messages
replay procedure
duplicate event handling
Kafka unavailable
```

---

## 21. Open Questions

| Question | Status |
|---|---|
| Should the outbox pattern be implemented from the first event-producing service? | Proposed for serious implementation |
| Should Notification Service consume `OrderCreated` directly or use `NotificationRequested`? | To decide |
| What protobuf schema registry approach will be used with Kafka? | To decide |
| What DLQ retention period should be used? | To decide |
| Should event replay tooling be built early or deferred? | To decide |
| Which consumers are safe for replay by default? | To document |
| Should all event payloads use full business snapshots or minimal IDs? | To decide per event |

---

## 22. Related Documents

This document should be read alongside:

```text
docs/events/event-catalog.md
docs/events/event-envelope.md
docs/events/kafka-topic-design.md
docs/events/ordering-and-idempotency.md
docs/events/retry-and-dlq-strategy.md
docs/architecture/communication-patterns.md
docs/architecture/service-boundaries.md
docs/architecture/resilience-patterns.md
docs/data/data-ownership.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
adr/0008-use-contract-first-service-design.md
```

---

## 23. Summary

bfstore uses event-driven design to decouple services and support asynchronous workflows.

The most important principles are:

```text
events are facts
event producers own the facts
consumers are independent
consumers are idempotent
eventual consistency is explicit
events are versioned contracts
event flows are observable and testable
```

The first implementation should focus on checkout events, then expand into search, reviews, recommendations, analytics, and operational projections as the platform matures.
