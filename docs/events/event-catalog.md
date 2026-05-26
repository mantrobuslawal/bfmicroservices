# Event Catalog

## 1. Purpose

This document defines the event catalog for **bfstore**, ACME Ltd’s fictional online furniture store backend.

The event catalog describes the business events published and consumed by bfstore services. It provides a shared reference for engineers, reviewers, and potential clients to understand how asynchronous communication works across the platform.

This document explains:

- which events exist
- which service owns each event
- why each event is published
- which services consume each event
- what the event means from a business perspective
- what consistency and idempotency expectations apply
- how events support the wider architecture

This document is intentionally business-facing and architecture-facing. Detailed protobuf schemas should live in the `proto/` directory.

---

## 2. Architecture Context

bfstore uses a hybrid communication model:

| Communication Type | Technology | Used For |
|---|---|---|
| Synchronous communication | gRPC | Commands and queries requiring an immediate response |
| Asynchronous communication | Kafka | Business facts that have already happened |
| Contract definition | Protobuf | Typed API and event payload contracts |

The guiding principle is:

> Commands that require an immediate result use gRPC. Facts that have already happened are published as Kafka events.

Examples:

| Type | Example | Transport |
|---|---|---|
| Command | `CreateOrder` | gRPC |
| Command | `ReserveStock` | gRPC |
| Command | `AuthorisePayment` | gRPC |
| Event | `OrderCreated` | Kafka |
| Event | `StockReserved` | Kafka |
| Event | `PaymentAuthorised` | Kafka |

---

## 3. Event Design Principles

## 3.1 Events Represent Facts

Events should describe something that has already happened.

Good event names:

```text
OrderCreated
PaymentAuthorised
StockReserved
ShipmentCreated
NotificationSent
```

Poor event names:

```text
CreateOrder
AuthorisePayment
ReserveStock
SendNotification
```

Commands belong in APIs. Facts belong in events.

---

## 3.2 The Owning Service Publishes the Event

The service that owns the business fact is responsible for publishing the event.

Examples:

| Event | Producer |
|---|---|
| `ProductUpdated` | `catalog-service` |
| `StockReserved` | `inventory-service` |
| `OrderCreated` | `order-service` |
| `PaymentAuthorised` | `payment-service` |
| `ShipmentCreated` | `shipping-service` |
| `NotificationSent` | `notification-service` |

A consumer must not publish events pretending to own facts from another service.

---

## 3.3 Events Are Part of the Contract

Events are integration contracts.

They must be:

- named consistently
- versioned
- documented
- observable
- tested
- backward compatible where possible

Breaking event changes should be treated with the same seriousness as breaking API changes.

---

## 3.4 Consumers Must Be Idempotent

Kafka may deliver duplicate messages or consumers may reprocess events after failures.

Consumers must be able to handle duplicate events safely.

Important examples:

- duplicate `OrderCreated` must not send duplicate customer emails where avoidable
- duplicate `PaymentAuthorised` must not confirm the same payment twice
- duplicate `ProductUpdated` must not corrupt the search index
- duplicate `StockReservationReleased` must not release stock more than once

---

## 3.5 Events Should Carry Correlation Context

Events should include enough metadata to trace a business flow across services.

Each event should include:

```text
event_id
event_type
event_version
occurred_at
producer
correlation_id
causation_id
trace_id
```

This is especially important for checkout, where a request crosses several services and later produces asynchronous side effects.

---

## 4. Event Envelope

All bfstore events should use a consistent envelope.

Example conceptual structure:

```json
{
  "event_id": "evt_01HXYZ...",
  "event_type": "OrderCreated",
  "event_version": "1.0",
  "occurred_at": "2026-05-26T10:15:30Z",
  "producer": "order-service",
  "correlation_id": "corr_01HXYZ...",
  "causation_id": "cmd_01HXYZ...",
  "trace_id": "trace_01HXYZ...",
  "data": {}
}
```

## 4.1 Envelope Fields

| Field | Required | Description |
|---|---:|---|
| `event_id` | Yes | Globally unique event identifier |
| `event_type` | Yes | Stable business event name |
| `event_version` | Yes | Event schema version |
| `occurred_at` | Yes | Time the business event occurred |
| `producer` | Yes | Service that produced the event |
| `correlation_id` | Yes | Identifier shared across a business flow |
| `causation_id` | Recommended | Identifier of the command or event that caused this event |
| `trace_id` | Recommended | Distributed trace identifier |
| `data` | Yes | Event-specific payload |

---

## 5. Topic Naming

Kafka topic names should be consistent and predictable.

Recommended topic format:

```text
bfstore.<domain>.<event-stream>.v<version>
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

## 5.1 Topic Design Notes

- Topic names should identify the domain area and stream purpose.
- Major incompatible schema changes should use a new topic version.
- Minor backward-compatible schema changes should keep the same topic.
- Avoid one topic per tiny event unless there is a strong operational reason.
- Avoid one giant topic for unrelated business domains.

---

## 6. Event Catalog Summary

| Event | Producer | Typical Consumers | Priority |
|---|---|---|---|
| `ProductCreated` | `catalog-service` | `search-service`, `recommendation-service` | Should |
| `ProductUpdated` | `catalog-service` | `search-service`, `recommendation-service` | Should |
| `ProductActivated` | `catalog-service` | `search-service`, `recommendation-service` | Should |
| `ProductDeactivated` | `catalog-service` | `basket-service`, `search-service`, `recommendation-service` | Should |
| `InventoryAdjusted` | `inventory-service` | `search-service`, `recommendation-service` | Should |
| `StockReserved` | `inventory-service` | `order-service` | Must |
| `StockReservationFailed` | `inventory-service` | `order-service` | Must |
| `StockReservationReleased` | `inventory-service` | `order-service` | Must |
| `BasketCreated` | `basket-service` | `recommendation-service` | Could |
| `BasketItemAdded` | `basket-service` | `recommendation-service` | Could |
| `BasketCheckedOut` | `basket-service` | `order-service`, `recommendation-service` | Should |
| `OrderCreated` | `order-service` | `notification-service`, `recommendation-service` | Must |
| `OrderConfirmed` | `order-service` | `notification-service`, `analytics consumers` | Should |
| `OrderFailed` | `order-service` | `notification-service`, `operations consumers` | Must |
| `OrderCancelled` | `order-service` | `payment-service`, `inventory-service`, `shipping-service`, `notification-service` | Should |
| `PaymentAuthorised` | `payment-service` | `order-service`, `notification-service` | Must |
| `PaymentFailed` | `payment-service` | `order-service`, `notification-service`, `inventory-service` | Must |
| `PaymentCaptured` | `payment-service` | `order-service`, `notification-service` | Should |
| `PaymentRefunded` | `payment-service` | `order-service`, `notification-service` | Should |
| `ShipmentCreated` | `shipping-service` | `order-service`, `notification-service` | Must |
| `ShipmentDispatched` | `shipping-service` | `notification-service`, `order-service` | Should |
| `ShipmentDelivered` | `shipping-service` | `notification-service`, `order-service` | Should |
| `ShipmentFailed` | `shipping-service` | `order-service`, `notification-service`, `operations consumers` | Must |
| `NotificationRequested` | `order-service` or domain service | `notification-service` | Should |
| `NotificationSent` | `notification-service` | `operations consumers` | Should |
| `NotificationFailed` | `notification-service` | `operations consumers` | Should |
| `ReviewCreated` | `review-service` | `search-service`, `recommendation-service` | Could |
| `ReviewApproved` | `review-service` | `search-service`, `recommendation-service`, `catalog-service` | Could |
| `SearchIndexUpdated` | `search-service` | `operations consumers` | Could |
| `RecommendationGenerated` | `recommendation-service` | `operations consumers` | Could |

---

## 7. Catalogue Events

## 7.1 ProductCreated

### Meaning

A new product has been created in the product catalogue.

### Producer

```text
catalog-service
```

### Topic

```text
bfstore.catalog.product-events.v1
```

### Typical Consumers

```text
search-service
recommendation-service
```

### Example Payload Fields

```text
product_id
name
category_id
status
price
currency
created_at
```

### Business Use

Consumers use this event to create product projections, search documents, or recommendation inputs.

### Consumer Expectations

- Consumers should treat Catalogue Service as the product source of truth.
- Consumers should support duplicate delivery.
- Consumers should be able to rebuild projections if needed.

---

## 7.2 ProductUpdated

### Meaning

An existing product has changed.

### Producer

```text
catalog-service
```

### Topic

```text
bfstore.catalog.product-events.v1
```

### Typical Consumers

```text
search-service
recommendation-service
```

### Example Payload Fields

```text
product_id
name
description
category_id
price
currency
material
colour
dimensions
updated_at
```

### Business Use

Used to keep search and recommendation projections aligned with catalogue data.

### Consistency

Search and recommendation updates are eventually consistent. Product updates do not need to block on downstream consumers.

---

## 7.3 ProductDeactivated

### Meaning

A product has been marked inactive and should no longer appear in customer-facing purchase flows.

### Producer

```text
catalog-service
```

### Topic

```text
bfstore.catalog.product-events.v1
```

### Typical Consumers

```text
basket-service
search-service
recommendation-service
```

### Business Use

- Search Service should remove or hide the product from search results.
- Recommendation Service should stop recommending the product.
- Basket Service may flag baskets containing the product.

### Important Rule

An inactive product must not be purchasable.

---

## 8. Inventory Events

## 8.1 InventoryAdjusted

### Meaning

The available or recorded stock level for a product has changed.

### Producer

```text
inventory-service
```

### Topic

```text
bfstore.inventory.stock-events.v1
```

### Typical Consumers

```text
search-service
recommendation-service
catalogue availability projections
```

### Example Payload Fields

```text
inventory_item_id
product_id
variant_id
warehouse_id
available_quantity
reserved_quantity
adjustment_reason
adjusted_at
```

### Business Use

Consumers may update availability displays, search filters, or recommendation eligibility.

---

## 8.2 StockReserved

### Meaning

Stock has been successfully reserved for a checkout or order.

### Producer

```text
inventory-service
```

### Topic

```text
bfstore.inventory.stock-events.v1
```

### Typical Consumers

```text
order-service
operations consumers
```

### Example Payload Fields

```text
reservation_id
order_id
basket_id
items
reserved_at
expires_at
```

### Business Use

Order Service uses this fact to continue or update the checkout process.

### Important Rules

- The event should be idempotent for consumers.
- Reservation expiry must be clear.
- The reservation should be linked to checkout or order context.

---

## 8.3 StockReservationFailed

### Meaning

Inventory Service could not reserve the requested stock.

### Producer

```text
inventory-service
```

### Topic

```text
bfstore.inventory.stock-events.v1
```

### Typical Consumers

```text
order-service
operations consumers
```

### Example Payload Fields

```text
reservation_request_id
order_id
basket_id
failed_items
reason
failed_at
```

### Business Use

Order Service should reject or fail checkout and must not attempt payment if stock reservation fails before payment.

### Important Rule

Payment must not be attempted when stock reservation fails before payment authorisation.

---

## 8.4 StockReservationReleased

### Meaning

A previous stock reservation has been released.

### Producer

```text
inventory-service
```

### Topic

```text
bfstore.inventory.stock-events.v1
```

### Typical Consumers

```text
order-service
operations consumers
```

### Business Use

Used when checkout fails, payment fails, or order cancellation releases reserved stock.

---

## 9. Basket Events

## 9.1 BasketItemAdded

### Meaning

A customer added an item to their basket.

### Producer

```text
basket-service
```

### Topic

```text
bfstore.basket.basket-events.v1
```

### Typical Consumers

```text
recommendation-service
analytics consumers
```

### Business Use

Can be used for recommendation signals and behaviour analytics.

### Initial Version

Optional. This event is not required for the first checkout vertical slice unless recommendation or analytics behaviour is implemented.

---

## 9.2 BasketCheckedOut

### Meaning

A basket has entered or completed the checkout process.

### Producer

```text
basket-service
```

### Topic

```text
bfstore.basket.basket-events.v1
```

### Typical Consumers

```text
order-service
recommendation-service
```

### Business Use

May be used to mark the basket as checked out and prevent further mutation.

---

## 10. Order Events

## 10.1 OrderCreated

### Meaning

An order has been created.

### Producer

```text
order-service
```

### Topic

```text
bfstore.order.order-events.v1
```

### Typical Consumers

```text
notification-service
recommendation-service
analytics consumers
operations consumers
```

### Example Payload Fields

```text
order_id
order_number
customer_id
basket_id
order_status
total_amount
currency
items
created_at
```

### Business Use

This is one of the most important events in the platform.

It may trigger:

- order confirmation notification
- recommendation updates
- analytics projections
- operational dashboards

### Important Rules

- `OrderCreated` must only be published after the order exists in Order Service.
- Consumers must handle duplicate events.
- Notification failure must not roll back order creation.

---

## 10.2 OrderFailed

### Meaning

Order creation or checkout failed.

### Producer

```text
order-service
```

### Topic

```text
bfstore.order.order-events.v1
```

### Typical Consumers

```text
notification-service
operations consumers
```

### Example Payload Fields

```text
checkout_attempt_id
order_id
customer_id
failure_reason
failed_at
```

### Business Use

Used for customer communication, operations dashboards, and failure analysis.

---

## 10.3 OrderCancelled

### Meaning

An order has been cancelled.

### Producer

```text
order-service
```

### Topic

```text
bfstore.order.order-events.v1
```

### Typical Consumers

```text
payment-service
inventory-service
shipping-service
notification-service
```

### Business Use

May trigger refund, stock release, shipment cancellation, and cancellation notification workflows.

### Important Rule

Cancellation consumers must be idempotent because cancellation events may be replayed or redelivered.

---

## 11. Payment Events

## 11.1 PaymentAuthorised

### Meaning

Payment authorisation succeeded.

### Producer

```text
payment-service
```

### Topic

```text
bfstore.payment.payment-events.v1
```

### Typical Consumers

```text
order-service
notification-service
operations consumers
```

### Example Payload Fields

```text
payment_id
order_id
customer_id
amount
currency
provider_reference
authorised_at
```

### Business Use

Order Service may use this event to update order state if payment handling is asynchronous.

### Security Notes

- Do not include raw card data.
- Do not include authentication credentials.
- Provider references should be safe for operational use.

---

## 11.2 PaymentFailed

### Meaning

Payment authorisation, capture, or refund failed.

### Producer

```text
payment-service
```

### Topic

```text
bfstore.payment.payment-events.v1
```

### Typical Consumers

```text
order-service
inventory-service
notification-service
operations consumers
```

### Business Use

May cause checkout failure, stock reservation release, and customer notification.

### Important Rule

Payment failure must not produce a confirmed order.

---

## 11.3 PaymentCaptured

### Meaning

Previously authorised payment has been captured.

### Producer

```text
payment-service
```

### Topic

```text
bfstore.payment.payment-events.v1
```

### Typical Consumers

```text
order-service
notification-service
operations consumers
```

### Initial Version

Optional. The first version may treat payment authorisation as sufficient and defer separate capture.

---

## 12. Shipping Events

## 12.1 ShipmentCreated

### Meaning

A shipment has been created for an order.

### Producer

```text
shipping-service
```

### Topic

```text
bfstore.shipping.shipment-events.v1
```

### Typical Consumers

```text
order-service
notification-service
operations consumers
```

### Example Payload Fields

```text
shipment_id
order_id
customer_id
tracking_reference
delivery_option
shipment_status
created_at
```

### Business Use

Order Service may update order fulfilment status. Notification Service may send a shipment confirmation.

---

## 12.2 ShipmentDispatched

### Meaning

A shipment has been dispatched.

### Producer

```text
shipping-service
```

### Topic

```text
bfstore.shipping.shipment-events.v1
```

### Typical Consumers

```text
notification-service
order-service
```

### Business Use

Used to notify customers and update order/shipment status views.

---

## 12.3 ShipmentDelivered

### Meaning

A shipment has been delivered.

### Producer

```text
shipping-service
```

### Topic

```text
bfstore.shipping.shipment-events.v1
```

### Typical Consumers

```text
notification-service
order-service
review-service
```

### Business Use

May trigger delivery notification and later review eligibility.

---

## 12.4 ShipmentFailed

### Meaning

Shipment creation or fulfilment failed.

### Producer

```text
shipping-service
```

### Topic

```text
bfstore.shipping.shipment-events.v1
```

### Typical Consumers

```text
order-service
notification-service
operations consumers
```

### Business Use

Used for operational visibility and customer communication.

---

## 13. Notification Events

## 13.1 NotificationRequested

### Meaning

A notification has been requested by a domain service.

### Producer

```text
order-service
payment-service
shipping-service
```

or another service that owns the triggering business event.

### Topic

```text
bfstore.notification.notification-events.v1
```

### Consumer

```text
notification-service
```

### Business Use

Provides a clear command-like request to the Notification Service while still using an event stream. This should be used carefully and documented because it is closer to a work-request pattern than a pure domain event.

### Design Note

For simple workflows, Notification Service may consume domain events directly, such as `OrderCreated`. For more complex workflows, a dedicated `NotificationRequested` event can make notification intent explicit.

---

## 13.2 NotificationSent

### Meaning

A notification was successfully sent or simulated.

### Producer

```text
notification-service
```

### Topic

```text
bfstore.notification.notification-events.v1
```

### Typical Consumers

```text
operations consumers
analytics consumers
```

### Business Use

Used for audit and operational visibility.

---

## 13.3 NotificationFailed

### Meaning

A notification failed after an attempted send.

### Producer

```text
notification-service
```

### Topic

```text
bfstore.notification.notification-events.v1
```

### Typical Consumers

```text
operations consumers
```

### Business Use

Used for retry handling, alerts, and operational reporting.

### Important Rule

Notification failure must not roll back the original business action.

---

## 14. Review Events

## 14.1 ReviewCreated

### Meaning

A customer submitted a product review.

### Producer

```text
review-service
```

### Topic

```text
bfstore.review.review-events.v1
```

### Typical Consumers

```text
search-service
recommendation-service
operations consumers
```

### Initial Version

Later phase. Not required for the checkout vertical slice.

---

## 14.2 ReviewApproved

### Meaning

A submitted review was approved for customer-facing display.

### Producer

```text
review-service
```

### Topic

```text
bfstore.review.review-events.v1
```

### Typical Consumers

```text
search-service
recommendation-service
catalog-service
```

### Business Use

Can update rating summaries and search/recommendation projections.

---

## 15. Search Events

## 15.1 SearchIndexUpdated

### Meaning

Search Service successfully updated its product search projection.

### Producer

```text
search-service
```

### Topic

```text
bfstore.search.search-events.v1
```

### Typical Consumers

```text
operations consumers
```

### Business Use

Used for monitoring projection health.

---

## 15.2 SearchIndexUpdateFailed

### Meaning

Search Service failed to update a search projection.

### Producer

```text
search-service
```

### Topic

```text
bfstore.search.search-events.v1
```

### Typical Consumers

```text
operations consumers
```

### Business Use

Used for alerting, replay, and operational investigation.

---

## 16. Recommendation Events

## 16.1 RecommendationGenerated

### Meaning

Recommendation Service generated recommendation results.

### Producer

```text
recommendation-service
```

### Topic

```text
bfstore.recommendation.recommendation-events.v1
```

### Typical Consumers

```text
operations consumers
analytics consumers
```

### Business Use

Used to monitor recommendation generation and support analysis of recommendation behaviour.

---

## 17. Checkout Event Flow

The main checkout flow may produce the following events:

```text
StockReserved
PaymentAuthorised
ShipmentCreated
OrderCreated
NotificationSent
```

Conceptual flow:

```text
Customer submits checkout
    -> Order Service coordinates checkout
    -> Inventory Service publishes StockReserved
    -> Payment Service publishes PaymentAuthorised
    -> Shipping Service publishes ShipmentCreated
    -> Order Service publishes OrderCreated
    -> Notification Service consumes OrderCreated
    -> Notification Service publishes NotificationSent
```

## 17.1 Checkout Failure Events

Possible failure events:

```text
StockReservationFailed
PaymentFailed
ShipmentFailed
OrderFailed
NotificationFailed
```

Failure handling rules:

- If stock reservation fails, payment should not be attempted.
- If payment fails, order should not be confirmed.
- If shipment creation fails, the order may fail or enter a pending fulfilment state depending on the chosen design.
- Notification failure must not roll back the order.

---

## 18. Event Versioning

## 18.1 Versioning Rules

Events must be versioned.

Recommended format:

```text
event_version: "1.0"
```

or through protobuf package versioning:

```text
acme.order.events.v1
```

## 18.2 Backward-Compatible Changes

Usually safe:

- adding optional fields
- adding new event types
- adding new enum values when consumers handle unknown values safely

## 18.3 Breaking Changes

Usually breaking:

- removing fields
- changing field meaning
- changing field type
- changing required semantics
- renaming event types
- changing topic meaning

Breaking changes should use a new version and be documented.

---

## 19. Idempotency and Deduplication

Consumers should track processed event IDs where duplicate processing could cause harm.

High-risk consumers:

| Consumer | Risk |
|---|---|
| `notification-service` | Duplicate emails or SMS messages |
| `payment-service` | Duplicate payment operations |
| `inventory-service` | Duplicate stock release or commit |
| `order-service` | Duplicate order state transitions |
| `search-service` | Projection inconsistency |
| `recommendation-service` | Duplicate recommendation signals |

Recommended deduplication inputs:

```text
event_id
event_type
producer
business_entity_id
event_version
```

---

## 20. Ordering Considerations

Kafka ordering is only guaranteed within a partition.

Partition keys should be chosen deliberately.

Recommended keys:

| Event Stream | Suggested Partition Key |
|---|---|
| Product events | `product_id` |
| Inventory events | `product_id` or `reservation_id` |
| Basket events | `basket_id` |
| Order events | `order_id` |
| Payment events | `payment_id` or `order_id` |
| Shipment events | `shipment_id` or `order_id` |
| Notification events | `notification_id` |
| Review events | `review_id` or `product_id` |

Ordering should not be assumed across unrelated aggregates.

---

## 21. Retry and Dead-Letter Handling

Consumers should use retry and dead-letter strategies appropriate to the event.

## 21.1 Retryable Failures

Examples:

- temporary database outage
- temporary downstream service failure
- temporary network issue
- temporary provider failure

## 21.2 Non-Retryable Failures

Examples:

- invalid event payload
- unsupported event version
- missing required business identifier
- malformed data

## 21.3 Dead-Letter Queues

Events that cannot be processed after retries should be sent to a DLQ.

Example DLQ naming:

```text
bfstore.order.order-events.dlq.v1
bfstore.payment.payment-events.dlq.v1
bfstore.notification.notification-events.dlq.v1
```

DLQ events should be observable and have runbook coverage.

---

## 22. Observability Requirements

Event-driven flows must be observable.

Each producer should record:

- events published
- publish failures
- publish latency
- event type
- topic
- partition where available
- correlation ID

Each consumer should record:

- events consumed
- events processed successfully
- processing failures
- retry count
- DLQ count
- consumer lag
- processing latency
- duplicate event count where measurable

Important dashboards:

```text
Kafka topic throughput
consumer lag
event processing errors
DLQ count
checkout event flow health
notification event health
search projection lag
```

---

## 23. Testing Requirements

Event contracts should be tested.

Test types:

| Test Type | Purpose |
|---|---|
| Producer contract tests | Confirm producer emits valid event payloads |
| Consumer contract tests | Confirm consumer handles expected event versions |
| Integration tests | Confirm Kafka publish/consume behaviour |
| Idempotency tests | Confirm duplicate events are safe |
| Replay tests | Confirm projections can be rebuilt |
| Failure tests | Confirm retry and DLQ behaviour |
| End-to-end tests | Confirm business journey across events |

Critical event-driven test cases:

```text
OrderCreated triggers notification
Duplicate OrderCreated does not send duplicate notification
ProductUpdated updates search projection
PaymentFailed prevents confirmed order
StockReservationFailed prevents payment attempt
```

---

## 24. Security and Privacy Considerations

Events must not leak sensitive data.

Do not include:

```text
raw card data
passwords
authentication tokens
full secret values
unnecessary personal data
internal stack traces
```

Be careful with:

```text
customer email
phone number
delivery address
provider references
payment failure details
```

Where possible, use IDs and allow the owning service to provide sensitive details through authorised APIs.

---

## 25. Initial Implementation Scope

The first version should focus on checkout-related events:

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

Events for reviews, search, recommendations, and advanced analytics can be added later.

This staged approach keeps the initial implementation focused while preserving a professional event-driven architecture.

---

## 26. Open Questions

| Question | Status |
|---|---|
| Should Notification Service consume `OrderCreated` directly or use `NotificationRequested`? | To decide |
| Should payment authorisation and capture be separate event flows in the first version? | To decide |
| Should shipment creation failure fail the order or place it in pending fulfilment? | To decide |
| Should search use Kafka projections from the start or begin with Catalogue Service queries? | To decide |
| What event schema registry approach will be used with protobuf and Kafka? | To decide |
| How long should DLQ events be retained? | To decide |
| What replay tooling should be provided for projections? | To decide |
| Should all events use outbox pattern from the first implementation? | Proposed for serious version |

---

## 27. Related Documents

This document should be read alongside:

```text
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/architecture/event-driven-design.md
docs/architecture/resilience-patterns.md
docs/data/data-ownership.md
docs/testing/testing-strategy.md
docs/observability/kafka-consumer-lag.md
docs/api/grpc-overview.md
```

Relevant ADRs:

```text
adr/0002-use-grpc-for-service-communication.md
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
adr/0008-use-contract-first-service-design.md
```

---

## 28. Summary

The bfstore event catalog defines how services communicate asynchronously without sharing databases or tightly coupling downstream workflows.

The most important event principles are:

```text
Events are facts, not commands.
The owning service publishes the event.
Consumers are idempotent.
Events carry correlation context.
Events are versioned contracts.
Event-driven flows are observable and testable.
```

The initial implementation should focus on checkout events first, then expand into search, reviews, recommendations, analytics, and operational projections as the platform matures.
