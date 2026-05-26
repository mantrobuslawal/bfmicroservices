# ADR-0002: Use gRPC for Service-to-Service Communication

## Status

Accepted

## Date

2026-05-26

## Context

bfstore services need to communicate for commands and queries that require an immediate response.

Examples include:

```text
API Gateway retrieving products
Order Service retrieving a basket
Order Service reserving stock
Order Service authorising payment
Order Service creating a shipment
```

These interactions require strongly typed request and response contracts, clear error behaviour, deadline propagation, and generated clients.

## Decision

bfstore will use gRPC for internal synchronous service-to-service communication.

Protobuf will be used to define service contracts and generate Go clients/servers.

## Drivers

This decision supports:

```text
strongly typed contracts
code generation
contract-first service design
consistent internal APIs
efficient service-to-service communication
clear versioned packages
Buf linting and breaking-change checks
```

## Alternatives Considered

### Option 1: REST/JSON Internally

Benefits:

```text
easy to inspect manually
widely familiar
simple browser tooling
```

Costs:

```text
weaker type safety
more manual client code
less strict compatibility checking
less efficient for internal service calls
```

### Option 2: gRPC Internally

Benefits:

```text
strong service contracts
generated clients
good Go ecosystem support
efficient binary protocol
well suited to internal APIs
```

Costs:

```text
requires protobuf tooling
less convenient for browser clients
requires deliberate error model
learning curve for contributors
```

### Option 3: Kafka for Most Communication

Benefits:

```text
high decoupling
asynchronous processing
natural event-driven behaviour
```

Costs:

```text
poor fit for immediate request/response behaviour
more difficult checkout orchestration
higher eventual consistency complexity
```

## Consequences

### Positive

```text
internal APIs are explicit and reviewable
service contracts can be linted and tested
generated clients reduce manual integration mistakes
gRPC metadata can carry correlation and auth context
```

### Negative

```text
developers need protobuf/gRPC tooling
client-facing browser APIs need API Gateway translation
testing requires generated clients or gRPC test tooling
error mapping must be standardised
```

## Implementation Notes

gRPC should be used for:

```text
GetProduct
ListProducts
GetBasket
ReserveStock
AuthorisePayment
CreateShipment
CreateOrder
GetOrder
```

Kafka should be used for facts that have already happened, such as:

```text
OrderCreated
PaymentAuthorised
StockReserved
ShipmentCreated
NotificationSent
```

## Error Handling

Services must use a consistent gRPC error model.

Examples:

```text
NOT_FOUND for missing products or orders
INVALID_ARGUMENT for validation failures
FAILED_PRECONDITION for business rule failures
UNAVAILABLE for unavailable dependencies
DEADLINE_EXCEEDED for timeouts
```

## Risks

| Risk | Mitigation |
|---|---|
| gRPC error handling becomes inconsistent | Maintain `docs/api/error-model.md` |
| Contracts break consumers | Use Buf breaking-change checks |
| API Gateway becomes too complex | Keep business logic in services |
| Generated code becomes messy | Standardise Buf and output paths |

## Review Triggers

Revisit this decision if:

```text
external clients need direct service access
team/tooling support for gRPC becomes a blocker
REST becomes more suitable for most internal calls
API Gateway protocol translation becomes excessive
```

## Related Documents

```text
docs/api/grpc-overview.md
docs/api/protobuf-style-guide.md
docs/api/error-model.md
docs/api/versioning.md
docs/architecture/communication-patterns.md
```
