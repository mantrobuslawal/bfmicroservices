# ADR-0001: Use Microservices for bfstore

## Status

Accepted

## Date

2026-05-26

## Context

bfstore is a cloud-native backend for ACME Ltd’s fictional online furniture store.

The project is intended to demonstrate senior-level capability across:

```text
backend service design
platform engineering
DevSecOps
Kubernetes
service ownership
observability
resilience
contract-first APIs
event-driven architecture
```

The application domain includes multiple business capabilities:

```text
catalogue
basket
inventory
orders
payments
shipping
notifications
reviews
search
recommendations
customers
authentication
```

A single monolithic application would be simpler to build initially, but would not demonstrate the same level of platform, service-boundary, operational, and distributed systems thinking.

## Decision

bfstore will use a microservice architecture.

Each major business capability will be represented as an independently owned service boundary.

Initial services include:

```text
api-gateway
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

Later services include:

```text
auth-service
customer-service
review-service
search-service
recommendation-service
```

## Drivers

This decision is driven by the need to demonstrate:

```text
clear service boundaries
service-owned data
gRPC service contracts
Kafka event-driven workflows
Kubernetes deployment patterns
observability across distributed services
resilience and failure handling
CI/CD and platform automation
security controls per service
```

## Alternatives Considered

### Option 1: Modular Monolith

A single deployable application organised into internal modules.

Benefits:

```text
simpler local development
simpler transactions
faster initial implementation
less infrastructure overhead
```

Costs:

```text
less evidence of distributed systems skill
fewer Kubernetes/platform engineering signals
less realistic service-to-service communication
less opportunity to demonstrate event-driven architecture
```

### Option 2: Microservices

Multiple independently deployable services aligned to business capabilities.

Benefits:

```text
clear ownership boundaries
strong platform engineering demonstration
realistic distributed communication
service-owned database design
resilience and observability opportunities
```

Costs:

```text
more operational complexity
more testing complexity
more local development complexity
risk of distributed monolith if boundaries are weak
```

### Option 3: Shared Database Distributed Application

Multiple services sharing one database.

Benefits:

```text
simpler reporting
simpler joins
less migration overhead
```

Costs:

```text
weak service ownership
tight coupling through database schema
unsafe for serious microservice boundaries
harder independent service evolution
```

## Consequences

### Positive

```text
bfstore can demonstrate professional microservice design
service boundaries are explicit and reviewable
each service can own its data and contracts
the architecture supports Kubernetes and GitOps demonstrations
distributed tracing, metrics, and resilience patterns become meaningful
```

### Negative

```text
implementation effort increases
local development requires more tooling
cross-service testing becomes more complex
data consistency requires careful design
debugging requires good observability
```

## Risks

| Risk | Mitigation |
|---|---|
| Services become too small or artificial | Align services to business capabilities, not database tables |
| Distributed monolith | Enforce service-owned databases and contract-first APIs |
| Too much complexity before useful behaviour exists | Start with checkout vertical slice |
| Local development becomes heavy | Use Docker Compose profiles and staged implementation |
| Documentation drifts from implementation | Treat docs as living artefacts and update through PRs |

## Implementation Notes

The first implementation should focus on:

```text
Browse product
Add to basket
Checkout
Reserve stock
Authorise payment
Create shipment
Create order
Publish OrderCreated
Consume OrderCreated in Notification Service
```

This proves the architecture without implementing every service immediately.

## Review Triggers

Revisit this decision if:

```text
most changes require coordinated releases across many services
service boundaries prove artificial
local development becomes unmanageable
the project fails to deliver working vertical slices
a modular monolith would better serve the project goals
```

## Related Documents

```text
docs/architecture/service-boundaries.md
docs/architecture/domain-model.md
docs/architecture/communication-patterns.md
docs/data/data-ownership.md
docs/requirements/functional-requirements.md
docs/testing/testing-strategy.md
```
