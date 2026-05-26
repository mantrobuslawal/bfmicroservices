# ADR-0011: Use Outbox Pattern for Critical Events

## Status

Proposed for serious implementation

## Date

2026-05-26

## Context

bfstore services update MySQL and publish Kafka events.

A common failure scenario is:

```text
service writes business data successfully
service fails before publishing Kafka event
downstream consumers never receive the business fact
```

Example:

```text
Order Service creates order in MySQL
Kafka publish of OrderCreated fails
Notification Service never sends confirmation
Search/analytics/recommendation consumers miss the event
```

This creates inconsistency between service state and event streams.

## Decision

bfstore will use the outbox pattern for critical event-producing services.

For the first serious implementation, `order-service` should use an outbox table for `OrderCreated`.

The pattern may later be expanded to:

```text
inventory-service
payment-service
shipping-service
catalog-service
review-service
notification-service
```

## Drivers

This decision supports:

```text
reliable event publication
database/event consistency
retryable Kafka publishing
observable event backlog
safer recovery from Kafka outages
professional event-driven architecture
```

## Alternatives Considered

### Option 1: Publish Directly After Database Commit

Benefits:

```text
simple implementation
fewer tables and workers
```

Costs:

```text
events can be lost after successful database writes
failure recovery is harder
less reliable event-driven architecture
```

### Option 2: Publish Before Database Commit

Benefits:

```text
event exists before data commit
```

Costs:

```text
consumers may see event for data that does not exist
unsafe and difficult to reason about
```

### Option 3: Outbox Pattern

Benefits:

```text
business data and outbox event written in one transaction
event publishing can be retried
event_id remains stable
outbox backlog can be monitored
```

Costs:

```text
additional table
publisher process required
more operational metrics
more tests
```

## Consequences

### Positive

```text
critical events are less likely to be lost
Kafka outages can be tolerated better
event publication becomes observable
retries preserve event identity
```

### Negative

```text
implementation complexity increases
outbox publisher needs monitoring
outbox table retention must be managed
duplicate publication is still possible and consumers must be idempotent
```

## Implementation Notes

A service transaction should write:

```text
business data
outbox event row
```

Example:

```text
Order Service transaction:
    insert order
    insert order_items
    insert outbox_events row for OrderCreated
```

A publisher then:

```text
reads pending outbox events
publishes to Kafka
marks events as published
retries failed publishes
exposes backlog metrics
```

## Candidate Outbox Fields

```text
outbox_event_id
event_id
event_type
event_version
aggregate_type
aggregate_id
payload
status
attempt_count
next_attempt_at
created_at
published_at
last_error
```

## Rules

```text
event_id must remain stable across publish retries
outbox row must be written in same transaction as business state
publisher must be idempotent where practical
consumers must still be idempotent
outbox backlog must be observable
```

## Risks

| Risk | Mitigation |
|---|---|
| Outbox publisher duplicates events | Consumers deduplicate by event_id |
| Outbox table grows forever | Add retention/cleanup policy |
| Publisher silently fails | Add metrics and alerts |
| Too much complexity too early | Start with Order Service only |
| Payload versioning unclear | Use event envelope and protobuf standards |

## Review Triggers

Revisit this decision if:

```text
direct publishing proves sufficient for portfolio scope
outbox implementation delays core vertical slice too much
a managed eventing platform provides equivalent guarantees
event loss risk is accepted for early demos only
```

## Related Documents

```text
docs/architecture/event-driven-design.md
docs/architecture/resilience-patterns.md
docs/events/event-envelope.md
docs/events/retry-and-dlq-strategy.md
docs/data/service-database-design.md
docs/data/migrations.md
```
