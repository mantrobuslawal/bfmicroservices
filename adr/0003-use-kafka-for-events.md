# ADR-0003: Use Kafka for Event-Driven Communication

## Status

Accepted

## Date

2026-05-26

## Context

bfstore needs asynchronous communication for workflows where downstream services react to business facts after they occur.

Examples:

```text
OrderCreated triggers notification
ProductUpdated updates search projection
ReviewApproved updates rating and recommendation signals
ShipmentDispatched triggers customer notification
```

The producer should not need to know every downstream consumer.

## Decision

bfstore will use Kafka for asynchronous event-driven communication.

Services will publish domain events to Kafka topics when important business facts occur.

## Drivers

This decision supports:

```text
decoupled downstream processing
multiple independent consumers
event-driven projections
notification workflows
search indexing
recommendation signals
consumer lag and DLQ observability
realistic platform engineering patterns
```

## Alternatives Considered

### Option 1: Direct Synchronous Calls

Benefits:

```text
simple to reason about initially
immediate response from downstream service
```

Costs:

```text
tight coupling
cascading failures
producers must know consumers
poor fit for non-critical side effects
```

### Option 2: Job Queue

Benefits:

```text
simpler than Kafka for background work
good for task processing
```

Costs:

```text
less suitable for multiple independent event consumers
weaker event stream and replay story
```

### Option 3: Kafka

Benefits:

```text
durable event streams
multiple consumer groups
good replay/projection model
industry-relevant platform skill
```

Costs:

```text
operational complexity
eventual consistency
consumer idempotency required
schema/versioning discipline required
```

## Consequences

### Positive

```text
Order Service can publish OrderCreated without calling every consumer
Notification Service can process asynchronously
Search and Recommendation can build projections from events
Kafka consumer lag and DLQs provide operational signals
```

### Negative

```text
eventual consistency must be understood and documented
events require versioning and schema governance
consumers must be idempotent
debugging requires strong observability
Kafka adds local and platform complexity
```

## Implementation Notes

Initial Kafka topics should include:

```text
bfstore.inventory.stock-events.v1
bfstore.payment.payment-events.v1
bfstore.shipping.shipment-events.v1
bfstore.order.order-events.v1
bfstore.notification.notification-events.v1
```

Initial events should include:

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

## Event Rules

```text
events describe facts, not commands
event producers must own the business fact
events must include an event envelope
events must be versioned
events must carry correlation context
consumers must be idempotent
poison events should go to DLQ after retry exhaustion
```

## Risks

| Risk | Mitigation |
|---|---|
| Kafka used as RPC | Use gRPC when immediate response is required |
| Event soup | Maintain event catalogue and topic design |
| Duplicate processing | Implement idempotent consumers |
| Lost events | Consider outbox pattern for critical producers |
| Hard debugging | Use correlation IDs, logs, metrics, traces |

## Review Triggers

Revisit this decision if:

```text
Kafka complexity outweighs value for the implementation stage
events are mostly used as request/response messages
consumer lag and DLQs become unmanageable
a simpler queue would satisfy the project goals
```

## Related Documents

```text
docs/architecture/event-driven-design.md
docs/events/event-catalog.md
docs/events/event-envelope.md
docs/events/kafka-topic-design.md
docs/events/ordering-and-idempotency.md
docs/events/retry-and-dlq-strategy.md
```
