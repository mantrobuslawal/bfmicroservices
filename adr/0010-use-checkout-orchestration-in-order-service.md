# ADR-0010: Use Order Service for Checkout Orchestration

## Status

Accepted for initial implementation

## Date

2026-05-26

## Context

Checkout requires coordination across several services:

```text
Basket Service
Inventory Service
Payment Service
Shipping Service
Order Service
Notification Service
```

The initial checkout flow must:

```text
validate basket
reserve stock
authorise payment
create shipment
create order
publish OrderCreated
trigger notification asynchronously
```

A design choice is needed for where orchestration should live.

## Decision

For the initial implementation, `order-service` will orchestrate the checkout workflow.

Order Service will coordinate synchronous calls to Basket, Inventory, Payment, and Shipping services, then create the order and publish `OrderCreated`.

## Drivers

This decision supports:

```text
clear ownership of order outcome
simple initial vertical slice
straightforward error handling
easier reasoning than choreography for first version
no early workflow engine dependency
```

## Alternatives Considered

### Option 1: API Gateway Orchestration

Benefits:

```text
simple request entry point
fewer internal orchestration calls from Order Service
```

Costs:

```text
puts business process in gateway
weakens domain boundaries
gateway becomes business workflow owner
```

### Option 2: Order Service Orchestration

Benefits:

```text
order lifecycle owner coordinates order creation
clear initial implementation
business workflow close to order state
easy to record checkout attempts
```

Costs:

```text
Order Service may grow orchestration complexity
must avoid owning payment, stock, and shipment internals
```

### Option 3: Event Choreography

Benefits:

```text
high decoupling
services react independently
```

Costs:

```text
harder to reason about initial checkout
more eventual consistency
more complex failure handling
harder for first vertical slice
```

### Option 4: Workflow Engine

Examples:

```text
Temporal
Cadence
Camunda
```

Benefits:

```text
excellent for long-running workflows
built-in retries and visibility
```

Costs:

```text
additional platform dependency
too much complexity for first implementation
```

## Consequences

### Positive

```text
checkout flow is explicit
order state and checkout attempts are centralised
first implementation remains achievable
failure cases can be tested directly
```

### Negative

```text
Order Service needs careful boundaries
orchestration logic may grow
long-running workflows may later need stronger tooling
```

## Checkout Flow

```text
Client
    -> API Gateway
        -> Order Service CreateOrder
            -> Basket Service GetBasket
            -> Inventory Service ReserveStock
            -> Payment Service AuthorisePayment
            -> Shipping Service CreateShipment
            -> Order Service creates order
            -> Order Service publishes OrderCreated
                -> Notification Service consumes OrderCreated
```

## Rules

```text
Order Service owns order lifecycle only
Inventory Service owns stock decisions
Payment Service owns payment decisions
Shipping Service owns shipment decisions
Notification Service owns notification delivery
Notification failure does not roll back order
```

## Risks

| Risk | Mitigation |
|---|---|
| Order Service becomes too large | Keep domain ownership explicit |
| Compensation paths grow complex | Document state machine and failure handling |
| Checkout becomes long-running | Revisit workflow engine later |
| Payment timeout ambiguity | Use idempotency and reconciliation |

## Review Triggers

Revisit this decision if:

```text
checkout becomes long-running
compensation logic becomes difficult to maintain
many additional services enter checkout
workflow visibility becomes poor
event choreography becomes more suitable
a workflow engine becomes justified
```

## Related Documents

```text
docs/architecture/communication-patterns.md
docs/architecture/resilience-patterns.md
docs/requirements/business-rules.md
docs/requirements/acceptance-criteria.md
docs/events/event-catalog.md
```
