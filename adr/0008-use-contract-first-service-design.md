# ADR-0008: Use Contract-First Service Design

## Status

Accepted

## Date

2026-05-26

## Context

bfstore has multiple services that communicate through gRPC and Kafka events.

To avoid accidental coupling, unclear APIs, and implementation-led contracts, service interfaces should be designed before service internals are finalised.

Contract-first design helps ensure that service boundaries are explicit, reviewable, testable, and aligned with business capabilities.

## Decision

bfstore will use contract-first service design.

Service APIs and event contracts should be designed and reviewed before implementation.

## Drivers

This decision supports:

```text
clear service boundaries
protobuf API consistency
event contract governance
consumer-driven thinking
testable service interactions
professional client-facing documentation
safer API evolution
```

## Alternatives Considered

### Option 1: Code-First Design

Benefits:

```text
fast initial coding
contracts emerge naturally from implementation
```

Costs:

```text
risk of exposing internal models
inconsistent APIs
harder consumer testing
late discovery of boundary problems
```

### Option 2: Contract-First Design

Benefits:

```text
clear APIs before implementation
better service boundary discipline
supports contract tests
supports generated clients
improves documentation quality
```

Costs:

```text
more upfront design work
contracts may need revision once implementation begins
requires discipline to keep contracts current
```

## Consequences

### Positive

```text
services can be built against stable contracts
API and event reviews happen early
tests can be planned from contracts
clients do not depend on database models
```

### Negative

```text
initial progress may feel slower
over-design is possible
contracts must be updated when implementation reveals better choices
```

## Implementation Notes

Design order:

```text
requirements
business rules
service boundaries
API contracts
event contracts
data ownership
database design
implementation
tests
deployment
operations
```

Contract artefacts:

```text
proto/acme/<domain>/v1/*.proto
docs/api/*
docs/events/*
docs/requirements/acceptance-criteria.md
```

## Rules

```text
do not expose database models as API contracts
do not publish events without documented ownership
do not add service APIs without clear business purpose
do not introduce breaking contract changes without versioning
```

## Risks

| Risk | Mitigation |
|---|---|
| Too much documentation before implementation | Move to vertical slice after foundation docs |
| Contracts become stale | Enforce docs/contracts review in PRs |
| Overly generic contracts | Use business language and examples |
| Hidden implementation leakage | Review protobuf and event payloads carefully |

## Review Triggers

Revisit this decision if:

```text
contract design delays all implementation progress
contracts repeatedly fail to match real service needs
a different workflow becomes more productive
```

## Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/events/event-catalog.md
docs/architecture/service-boundaries.md
docs/testing/testing-strategy.md
```
