# Retry and Dead-Letter Queue Strategy

## 1. Purpose

This document defines the retry and dead-letter queue strategy for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how services should handle temporary failures, permanent failures, poison messages, retry limits, dead-letter topics, replay, alerting, and operational support.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s reliability and operational readiness.

---

## 2. Scope

This document covers retry and DLQ handling for:

```text
Kafka event consumers
Kafka event producers
gRPC client calls
database operations
external provider integrations
notification delivery
checkout-related workflows
projection updates
```

It focuses on strategy and standards. Service-specific implementation details should be documented in service README files and runbooks.

---

## 3. Design Goals

The retry and DLQ strategy should ensure:

| Goal | Description |
|---|---|
| Reliability | Temporary failures can recover without manual action |
| Safety | Retries do not create duplicate harmful side effects |
| Visibility | Failed processing is observable and alertable |
| Isolation | Poison messages do not block healthy messages forever |
| Recoverability | Failed events can be inspected and replayed where safe |
| Business correctness | Checkout, payment, stock, and notification flows fail safely |
| Operational clarity | Runbooks explain what to do when retries or DLQs grow |

---

## 4. Core Principles

## 4.1 Retry Only When Safe

A retry is safe only when the operation is idempotent or protected against duplicate side effects.

Safe examples:

```text
read-only queries
idempotent stock reservation
idempotent payment authorisation with provider key
idempotent shipment creation
projection update by event_id
```

Unsafe examples:

```text
blind payment capture
blind refund
blind notification send
blind stock release
```

---

## 4.2 Classify Failures

Every retry decision should begin with failure classification.

Failure types:

```text
retryable temporary failure
non-retryable permanent failure
business rule failure
poison message
unknown failure
```

---

## 4.3 Do Not Retry Forever

Unbounded retries can hide problems and block progress.

Every retry policy should define:

```text
maximum attempts
backoff strategy
failure classification
DLQ behaviour
alerting threshold
```

---

## 4.4 DLQs Are Not a Fix

A DLQ is a safety mechanism, not a resolution.

Every DLQ needs:

```text
owner
dashboard
alert
runbook
inspection process
replay or discard decision
```

---

## 5. Failure Classification

## 5.1 Retryable Failures

Examples:

```text
temporary database outage
Kafka broker temporarily unavailable
downstream gRPC service unavailable
network timeout
provider timeout
temporary provider rate limit
consumer rebalance interruption
```

Handling:

```text
retry with backoff
preserve correlation context
record retry count
emit metrics
eventually DLQ if unresolved
```

---

## 5.2 Non-Retryable Failures

Examples:

```text
invalid event schema
unsupported event version
missing required event_id
invalid enum value with no safe fallback
malformed payload
business entity permanently invalid
```

Handling:

```text
do not repeatedly retry
send to DLQ
log safe failure reason
alert if volume exceeds threshold
```

---

## 5.3 Business Rule Failures

Examples:

```text
insufficient stock
payment declined
order already cancelled
product inactive
delivery option invalid
```

Handling depends on context.

These are usually not DLQ cases unless they appear unexpectedly in an event consumer.

---

## 5.4 Poison Messages

A poison message repeatedly fails because of the message content or unsupported assumptions.

Handling:

```text
stop retrying after configured attempts
send to DLQ
alert owner
inspect and decide repair, replay, or discard
```

---

## 6. Kafka Consumer Retry Strategy

Kafka consumers should use controlled retries.

Recommended flow:

```text
consume event
validate envelope
validate event version
process event idempotently
commit offset after successful processing
if retryable failure:
    retry with backoff
if retries exhausted:
    send to DLQ
    commit or skip according to consumer design
if non-retryable failure:
    send to DLQ
    commit or skip according to consumer design
```

## 6.1 Retry Attempts

Initial guidance:

```text
3 to 5 immediate or short retries for transient errors
longer delayed retries only where justified
DLQ after retry exhaustion
```

Exact numbers should be tuned through testing and operational requirements.

---

## 7. Backoff Strategy

Use exponential backoff with jitter where practical.

Example conceptual policy:

```text
attempt 1: immediate
attempt 2: 1 second
attempt 3: 5 seconds
attempt 4: 30 seconds
attempt 5: DLQ
```

Jitter helps avoid retry storms when many consumers fail at once.

---

## 8. Retry Topics

Retry topics may be used for delayed retries.

Example:

```text
bfstore.order.order-events.retry.5m.v1
bfstore.order.order-events.retry.30m.v1
bfstore.order.order-events.dlq.v1
```

## 8.1 Initial Recommendation

For the first implementation:

```text
use simple consumer-level retry
use DLQ after exhaustion
defer retry topics unless needed
```

Retry topics can be introduced once operational requirements justify the added complexity.

---

## 9. Dead-Letter Queue Design

## 9.1 DLQ Naming

Recommended DLQ format:

```text
bfstore.<domain>.<stream>.dlq.v<major>
```

Examples:

```text
bfstore.order.order-events.dlq.v1
bfstore.payment.payment-events.dlq.v1
bfstore.notification.notification-events.dlq.v1
bfstore.search.search-events.dlq.v1
```

## 9.2 DLQ Payload

A DLQ message should preserve the original event and include failure metadata.

Example conceptual structure:

```json
{
  "original_topic": "bfstore.order.order-events.v1",
  "original_partition": 2,
  "original_offset": 9182,
  "failure_reason": "unsupported_event_version",
  "failure_message": "event version 3.0 is not supported",
  "failed_at": "2026-05-26T10:20:00Z",
  "consumer": "notification-service",
  "retry_count": 5,
  "correlation_id": "corr_123",
  "original_event": {}
}
```

## 9.3 DLQ Metadata

DLQ messages should include:

```text
original topic
original partition
original offset
event_id
event_type
event_version
producer
consumer
correlation_id
failure category
failure reason
retry count
failed_at
```

Do not include sensitive data beyond what already exists in the original event.

---

## 10. DLQ Ownership

Each DLQ must have an owner.

Ownership usually belongs to the consuming service that failed to process the event.

Examples:

| DLQ | Owner |
|---|---|
| `bfstore.order.order-events.dlq.v1` for Notification consumer failures | `notification-service` team/owner |
| `bfstore.catalog.product-events.dlq.v1` for Search consumer failures | `search-service` team/owner |
| `bfstore.payment.payment-events.dlq.v1` for Order consumer failures | `order-service` team/owner |

If multiple consumers need DLQs for the same source topic, consumer-specific DLQs may be clearer.

Example:

```text
bfstore.notification.order-events.dlq.v1
bfstore.search.product-events.dlq.v1
```

The final naming model should be documented once consumer ownership is implemented.

---

## 11. Producer Retry Strategy

Kafka producers may fail to publish because of:

```text
broker unavailable
network error
timeout
authentication issue
serialisation failure
configuration error
```

## 11.1 Retryable Producer Failures

Temporary broker or network failures may be retried.

Rules:

```text
preserve event_id across retries
do not regenerate event payload for the same business fact
log publish attempts
emit publish failure metrics
```

## 11.2 Non-Retryable Producer Failures

Examples:

```text
serialisation failure
invalid topic
invalid credentials
unsupported schema
```

These should fail fast and alert.

## 11.3 Outbox Recommendation

For important domain events, use the outbox pattern.

Recommended services:

```text
order-service
payment-service
inventory-service
shipping-service
catalog-service
review-service
```

---

## 12. gRPC Retry Strategy

gRPC retries should be used carefully.

## 12.1 Safe to Retry

Usually safe:

```text
GetProduct
ListProducts
GetBasket
GetOrder
```

Conditionally safe with idempotency:

```text
ReserveStock
AuthorisePayment
CreateShipment
CreateOrder
RefundPayment
```

## 12.2 Unsafe to Retry Blindly

Avoid blind retries for:

```text
payment authorisation
payment capture
refund
stock reservation
notification send
shipment creation
```

Use idempotency keys and status checks.

---

## 13. Database Retry Strategy

Database failures may be retryable when caused by temporary conditions.

Examples:

```text
deadlock
connection pool exhaustion
temporary connection loss
transaction conflict
```

Retry rules:

```text
retry only safe transactions
use bounded retry count
log retry reason
do not hide persistent failures
surface readiness issues when database unavailable
```

---

## 14. External Provider Retry Strategy

External providers may include:

```text
payment provider
email provider
SMS provider
shipping carrier
```

Rules:

```text
use provider idempotency keys where available
classify provider errors
respect provider rate limits
do not retry permanent declines
record attempts
support reconciliation for ambiguous outcomes
```

Payment provider timeouts require special care because the provider may have processed the request even if bfstore timed out.

---

## 15. Notification Retry Strategy

Notification delivery should be retried when failure is temporary.

Retryable:

```text
email provider unavailable
SMS provider timeout
temporary rate limit
network failure
```

Non-retryable:

```text
invalid recipient address
unsupported channel
template missing
permanent provider rejection
```

Notification failure should not roll back the original business event.

---

## 16. Search and Projection Retry Strategy

Search projection consumers should retry temporary failures.

Examples:

```text
search projection database unavailable
temporary indexing failure
consumer restart
```

Projection consumers are often good candidates for replay.

Rules:

```text
updates must be idempotent
offsets should be tracked
projection lag should be observable
rebuild process should be documented
```

---

## 17. DLQ Replay

DLQ replay should be controlled.

Before replay:

```text
inspect failure reason
confirm bug or dependency issue is resolved
confirm event is safe to replay
confirm consumer idempotency
choose replay scope
monitor replay
```

Replay modes:

```text
single message replay
filtered replay by event_type
filtered replay by correlation_id
bulk replay after bug fix
projection rebuild replay
```

High-risk consumers such as notifications and payments should not be bulk replayed without explicit safeguards.

---

## 18. Alerting

Alerts should be configured for:

```text
DLQ count above threshold
consumer lag above threshold
retry rate above threshold
Kafka publish failures
outbox backlog
notification failure spike
payment provider retry spike
search projection lag
```

Alert severity should consider business impact.

Example:

| Alert | Severity |
|---|---|
| Payment event DLQ growing | High |
| Order event consumer lag growing | High |
| Search projection lag | Medium |
| Recommendation consumer lag | Low to Medium |
| Notification DLQ growing | Medium to High depending on volume |

---

## 19. Runbook Requirements

Each DLQ or retry-heavy consumer should have a runbook.

Runbook sections:

```text
symptoms
dashboard links
common causes
diagnostic steps
safe replay procedure
discard procedure
escalation criteria
rollback options
related services
```

Example runbooks:

```text
Order event DLQ runbook
Payment event DLQ runbook
Notification failure runbook
Search projection lag runbook
Outbox backlog runbook
```

---

## 20. Observability

Metrics:

```text
consumer_retries_total
consumer_retry_exhausted_total
dlq_messages_total
dlq_publish_failures_total
event_processing_failures_total
event_processing_duration_seconds
outbox_pending_events
outbox_publish_failures_total
grpc_client_retries_total
provider_retries_total
```

Logs:

```text
event_id
event_type
correlation_id
consumer
producer
retry_count
failure_category
failure_reason
dlq_topic
original_topic
original_offset
```

Traces:

```text
record retry attempts as span events
mark failed spans clearly
link consumer span to producer trace where possible
```

---

## 21. Testing Requirements

Tests should cover:

```text
retryable consumer failure is retried
non-retryable consumer failure goes to DLQ
retry exhaustion sends event to DLQ
DLQ message contains original event and failure metadata
duplicate replay does not cause duplicate side effect
payment timeout is handled with idempotency
notification provider failure is retried
invalid event version goes to DLQ
outbox publish failure is retried later
```

---

## 22. Initial Implementation Scope

The first implementation should include retry and DLQ handling for:

```text
OrderCreated consumed by notification-service
PaymentFailed consumed by order-service where applicable
StockReservationFailed consumed by order-service where applicable
ShipmentFailed consumed by order-service where applicable
ProductUpdated consumed by search-service if search is implemented early
```

Minimum initial features:

```text
bounded retry
basic backoff
DLQ topic
failure logging
retry metrics
DLQ metrics
manual replay notes
```

---

## 23. Anti-Patterns to Avoid

Avoid:

```text
infinite retries
retrying payment capture blindly
DLQs with no owner
DLQs with no alerting
discarding failed events silently
regenerating event_id during producer retry
retrying invalid schema messages repeatedly
bulk replaying notifications without safeguards
using DLQ as normal control flow
```

---

## 24. Open Questions

| Question | Status |
|---|---|
| Will retry topics be implemented or deferred? | Proposed: defer initially |
| Will DLQs be source-topic based or consumer-specific? | To decide |
| What retry library or framework will Go services use? | To decide |
| What is the default retry count for consumers? | To decide |
| How will DLQ replay be authorised and audited? | To decide |
| Will outbox publishing be implemented from the first event producer? | Proposed for order-service |
| How long should DLQ topics retain messages? | To decide |

---

## 25. Related Documents

This document should be read alongside:

```text
docs/events/event-envelope.md
docs/events/kafka-topic-design.md
docs/events/ordering-and-idempotency.md
docs/events/event-catalog.md
docs/architecture/event-driven-design.md
docs/architecture/resilience-patterns.md
docs/api/error-model.md
docs/testing/testing-strategy.md
docs/operations/runbooks.md
docs/operations/incident-response.md
```

---

## 26. Summary

bfstore’s retry and DLQ strategy is designed to recover from temporary failures without hiding permanent failures or creating duplicate business side effects.

The most important rules are:

```text
classify failures before retrying
retry only when safe
use bounded retries
preserve event_id across retries
send poison messages to DLQ
monitor DLQs and consumer lag
provide runbooks and replay procedures
protect payment, stock, shipment, and notification operations with idempotency
```

This strategy makes bfstore’s event-driven workflows safer, more observable, and more operationally credible.
