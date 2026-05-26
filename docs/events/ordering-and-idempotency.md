# Ordering and Idempotency

## 1. Purpose

This document defines the ordering and idempotency strategy for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how services should handle duplicate requests, duplicate events, retries, replay, Kafka ordering limits, and business operations that must not happen more than once.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s reliability and event-driven architecture.

---

## 2. Scope

This document covers ordering and idempotency for:

```text
gRPC commands
checkout requests
stock reservations
payment authorisations
shipment creation
notification sending
Kafka event producers
Kafka event consumers
event replay
projection rebuilding
```

It should be read alongside:

```text
docs/events/event-envelope.md
docs/events/kafka-topic-design.md
docs/events/retry-and-dlq-strategy.md
docs/architecture/resilience-patterns.md
```

---

## 3. Why Ordering and Idempotency Matter

Distributed systems frequently process work more than once or out of the expected order.

This can happen because of:

```text
client retries
network timeouts
service restarts
Kafka redelivery
consumer rebalances
manual replay
database transaction conflicts
provider timeouts
partial failures
```

Without idempotency, bfstore risks:

```text
duplicate orders
duplicate payments
duplicate stock releases
duplicate customer notifications
corrupted projections
incorrect order state transitions
```

---

## 4. Key Principles

## 4.1 Do Not Assume Exactly-Once Business Behaviour

Even if infrastructure provides strong guarantees, business logic must still protect itself.

Services should assume:

```text
requests may be retried
events may be duplicated
events may be replayed
consumers may restart
timeouts may hide successful downstream processing
```

---

## 4.2 Make Critical Operations Idempotent

Idempotency is required where duplicate processing could cause business harm.

Critical operations:

```text
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
SendNotification
RefundPayment
ReleaseStockReservation
```

---

## 4.3 Ordering Is Local, Not Global

Kafka does not provide global ordering across all events.

Ordering is only guaranteed within a partition.

Therefore:

```text
events for the same aggregate should use the same partition key
consumers should avoid assuming ordering across different aggregates
business state machines should reject invalid transitions
```

---

## 5. Idempotency Definitions

## 5.1 Idempotent Request

A request is idempotent if sending it more than once has the same business effect as sending it once.

Example:

```text
CreateOrder with the same idempotency key returns the same order result.
```

## 5.2 Idempotent Consumer

A consumer is idempotent if processing the same event more than once does not create incorrect duplicate side effects.

Example:

```text
Notification Service receives the same OrderCreated event twice but records only one confirmation notification for the event.
```

---

## 6. Idempotency Keys

Idempotency keys identify a logical operation.

Example checkout input:

```text
customer_id
basket_id
idempotency_key
```

Rules:

```text
idempotency keys must be stable across retries
idempotency keys should be scoped to the operation and caller
same key with same request should return original result
same key with different request should be rejected
idempotency records should have a retention period
```

---

## 7. Idempotency Key Scope

Recommended scoping:

| Operation | Scope |
|---|---|
| `CreateOrder` | `customer_id + basket_id + idempotency_key` |
| `ReserveStock` | `order_id or checkout_attempt_id + idempotency_key` |
| `AuthorisePayment` | `order_id + payment_request_id or idempotency_key` |
| `CreateShipment` | `order_id + shipment_request_id or idempotency_key` |
| `SendNotification` | `event_id + notification_type + recipient` |
| `RefundPayment` | `payment_id + refund_request_id or idempotency_key` |

---

## 8. Idempotency Storage

Services that require idempotency should persist idempotency records.

Example fields:

```text
idempotency_key
operation_name
request_hash
response_reference
status
created_at
expires_at
```

## 8.1 Request Hash

A request hash helps detect accidental reuse of the same key with a different request.

Behaviour:

```text
same key + same request hash -> return original result
same key + different request hash -> reject as conflict
```

---

## 9. gRPC Idempotency

Idempotency should be explicit in gRPC APIs for critical commands.

Example:

```proto
message CreateOrderRequest {
  string customer_id = 1;
  string basket_id = 2;
  string idempotency_key = 3;
}
```

Operations that should include idempotency keys:

```text
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
ReleaseStockReservation
RefundPayment
```

Read operations such as `GetProduct` and `GetOrder` are naturally idempotent.

---

## 10. Checkout Idempotency

Checkout is the most important idempotency scenario.

Duplicate checkout may happen because:

```text
customer double-clicks checkout
API Gateway retries after timeout
client retries after network error
Order Service times out waiting for Payment Service
```

Expected behaviour:

```text
same checkout idempotency key returns the same result
duplicate request must not create duplicate confirmed orders
duplicate request must not authorise payment twice
duplicate request must not reserve stock twice
```

---

## 11. Stock Reservation Idempotency

Inventory Service must prevent duplicate reservation effects.

Rules:

```text
same reservation request should return the existing reservation
same idempotency key with different items should be rejected
stock must not be reserved twice for the same logical checkout step
reservation release must also be idempotent
```

Important operations:

```text
ReserveStock
ReleaseStockReservation
ExpireStockReservation
CommitStockReservation
```

---

## 12. Payment Idempotency

Payment operations must be carefully idempotent.

Rules:

```text
same payment request must not authorise payment twice
provider idempotency keys should be used where available
payment attempts must be recorded
timeouts must be reconciled safely
payment retries must be tied to a stable request ID
```

Risky operations:

```text
AuthorisePayment
CapturePayment
RefundPayment
VoidPayment
```

Payment timeouts are especially sensitive because the provider may have processed the request even if bfstore did not receive the response.

---

## 13. Shipment Idempotency

Shipping Service should prevent duplicate shipments.

Rules:

```text
same shipment request should return existing shipment
same idempotency key with different shipment details should be rejected
carrier/provider references should be recorded where relevant
shipment cancellation should be idempotent
```

---

## 14. Notification Idempotency

Notification Service should avoid duplicate customer messages where practical.

Potential deduplication inputs:

```text
event_id
notification_type
recipient
order_id
template_id
```

Rules:

```text
same OrderCreated event should not send duplicate order confirmation
same ShipmentDispatched event should not send duplicate dispatch notification
provider retries should not create duplicate records
```

Some channels may still deliver duplicates in rare provider failure cases, so the system should record attempts clearly.

---

## 15. Kafka Producer Idempotency

Event producers should avoid publishing duplicate events for the same business fact.

Recommended approaches:

```text
publish event as part of outbox workflow
use stable event_id for retries
avoid generating a new event_id for the same outbox record
record publish status
monitor publish failures
```

If the same outbox event is retried, it should retain the same `event_id`.

---

## 16. Kafka Consumer Idempotency

Consumers must handle duplicate events safely.

Recommended consumer deduplication table fields:

```text
event_id
event_type
producer
processed_at
processing_status
failure_reason
```

Example behaviour:

```text
if event_id already processed successfully:
    skip or acknowledge
else:
    process event
    record event_id as processed
```

---

## 17. Event Replay

Replay is useful but can be dangerous.

Safe replay candidates:

```text
search-service projections
recommendation-service projections
analytics projections
read model rebuilds
```

High-risk replay candidates:

```text
notification-service
payment-service
shipping-service
inventory-service
```

High-risk consumers need explicit replay-safe logic.

Example:

```text
replaying OrderCreated should rebuild a projection, not resend customer emails unless explicitly requested
```

---

## 18. Kafka Ordering

Kafka ordering is guaranteed only within a partition.

Example:

```text
OrderCreated(order_id=123)
OrderCancelled(order_id=123)
```

If both events use `order_id` as the partition key, they will be ordered relative to each other.

But:

```text
OrderCreated(order_id=123)
PaymentAuthorised(payment_id=999)
ShipmentCreated(shipment_id=555)
```

may be on different topics or partitions and should not be assumed to have global order.

---

## 19. Recommended Partition Keys

| Stream | Recommended Key |
|---|---|
| Product events | `product_id` |
| Basket events | `basket_id` |
| Stock events | `product_id` or `reservation_id` |
| Order events | `order_id` |
| Payment events | `payment_id` or `order_id` |
| Shipment events | `shipment_id` or `order_id` |
| Notification events | `notification_id` |
| Review events | `review_id` or `product_id` |
| Search events | `projection_id` or `product_id` |
| Recommendation events | `customer_id` or `product_id` |

---

## 20. State Machines and Ordering Protection

Services should protect business state transitions.

Example order states:

```text
PENDING
CONFIRMED
CANCELLED
FAILED
```

Invalid transitions should be rejected or ignored.

Examples:

```text
CANCELLED -> CONFIRMED should not be allowed
FAILED -> CONFIRMED should not be allowed without explicit recovery process
DELIVERED -> DISPATCHED should not be allowed
```

State machines reduce damage from out-of-order events and duplicate messages.

---

## 21. Handling Out-of-Order Events

Consumers should decide how to handle events that arrive unexpectedly.

Options:

```text
ignore if stale
store pending until prerequisite state arrives
fetch current state from owning service
send to DLQ if invalid
process idempotently if safe
```

Example:

```text
ShipmentDelivered arrives before ShipmentDispatched
```

Consumer options:

```text
Shipping Service may accept delivered as a later state.
Notification Service may send delivered notification only once.
Order Service may update fulfilment state if transition is valid.
```

---

## 22. Request Timeout Ambiguity

Timeouts create uncertainty.

Example:

```text
Order Service calls Payment Service.
Payment Service authorises payment.
Network timeout occurs before response reaches Order Service.
```

Order Service does not know whether payment succeeded.

Mitigation:

```text
use idempotency key
retry safely with same key
query payment status by request ID
record payment attempt
use reconciliation workflow
```

---

## 23. Retention of Idempotency Records

Idempotency records should not grow forever.

Retention depends on operation risk.

Example guidance:

| Operation | Retention Consideration |
|---|---|
| Checkout | long enough to cover client retries and support queries |
| Payment | aligned with payment reconciliation needs |
| Shipment | aligned with provider retry/reconciliation needs |
| Notification | long enough to avoid duplicate event sends |
| Event consumer deduplication | aligned with replay and topic retention |

Final retention should be documented in data retention standards.

---

## 24. Observability

Idempotency and ordering should be observable.

Metrics:

```text
idempotency_hits_total
idempotency_conflicts_total
duplicate_events_total
stale_events_total
invalid_state_transitions_total
event_replay_count
consumer_lag
dlq_messages_total
```

Logs should include:

```text
idempotency_key
event_id
event_type
correlation_id
operation
existing_result_reference
state_transition
deduplication_result
```

Avoid logging sensitive values.

---

## 25. Testing Requirements

Tests should cover:

```text
duplicate checkout request returns same order
duplicate stock reservation does not double reserve
duplicate payment request does not double authorise
duplicate shipment request does not create duplicate shipment
duplicate OrderCreated does not send duplicate notification
event replay rebuilds search projection safely
invalid order state transition is rejected
same idempotency key with different request is rejected
```

---

## 26. Initial Implementation Scope

The first version should implement idempotency for:

```text
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
Notification processing of OrderCreated
Kafka consumer deduplication for critical consumers
```

Ordering should be explicitly handled for:

```text
order lifecycle events
payment events for a payment
shipment events for a shipment
stock reservation events
```

---

## 27. Anti-Patterns to Avoid

Avoid:

```text
blind retries of payment operations
generating new event_id for every publish retry
assuming Kafka provides global ordering
using timestamps alone for ordering
sending duplicate notifications during replay
treating idempotency keys as optional for critical commands
ignoring same key with different request body
allowing invalid state transitions silently
```

---

## 28. Open Questions

| Question | Status |
|---|---|
| What ID format will be used for idempotency keys generated by clients? | To decide |
| How long should idempotency records be retained per service? | To decide |
| Should idempotency records store full responses or response references? | To decide |
| Which services need deduplication tables from the first implementation? | To decide |
| Should replay mode be explicitly flagged in event consumers? | Proposed |
| Should payment reconciliation be part of the first version? | Deferred or simulated |

---

## 29. Related Documents

This document should be read alongside:

```text
docs/events/event-envelope.md
docs/events/kafka-topic-design.md
docs/events/retry-and-dlq-strategy.md
docs/architecture/resilience-patterns.md
docs/architecture/event-driven-design.md
docs/api/error-model.md
docs/testing/testing-strategy.md
docs/data/data-ownership.md
```

---

## 30. Summary

bfstore must treat ordering and idempotency as first-class design concerns.

The most important rules are:

```text
critical commands need idempotency keys
duplicate events must be safe
event_id must stay stable during publish retries
Kafka ordering is only per partition
state machines must reject invalid transitions
payment and stock operations must not be retried blindly
event replay must not create harmful side effects
```

This approach protects bfstore from duplicate orders, duplicate payments, stock errors, and inconsistent projections.
