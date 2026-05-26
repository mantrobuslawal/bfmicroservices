# ADR-0009: Use an API Gateway

## Status

Accepted

## Date

2026-05-26

## Context

bfstore has multiple backend services. External clients should not need to know the internal service topology or call each service directly.

The system needs a stable client-facing entry point that can handle routing, authentication integration, request shaping, correlation ID propagation, and safe error mapping.

## Decision

bfstore will use an API Gateway as the client-facing entry point.

The API Gateway will call backend services using gRPC.

## Drivers

This decision supports:

```text
single external entry point
internal service topology hiding
client-friendly API shape
auth integration point
correlation ID propagation
safe external error mapping
rate limiting and edge controls later
```

## Alternatives Considered

### Option 1: Expose Services Directly

Benefits:

```text
fewer components
simple for internal testing
```

Costs:

```text
clients depend on internal topology
harder auth and error consistency
more exposed attack surface
harder API evolution
```

### Option 2: API Gateway

Benefits:

```text
stable external API
central routing
central edge controls
safe error translation
hides internal services
```

Costs:

```text
additional component
risk of gateway becoming a monolith
must avoid domain logic in gateway
```

### Option 3: Backend-for-Frontend

Benefits:

```text
tailored APIs per client type
good for complex frontend needs
```

Costs:

```text
more components
not required for first implementation
```

## Consequences

### Positive

```text
clients use one entry point
internal services remain private
gateway can map REST/JSON to internal gRPC if required
correlation IDs can start at the edge
client error responses can be standardised
```

### Negative

```text
gateway must be deployed and observed
gateway can become a bottleneck
domain logic may drift into gateway if not controlled
```

## Gateway Responsibilities

The API Gateway may:

```text
route requests
perform request shape validation
propagate correlation IDs
enforce authentication at the edge
map gRPC errors to client-safe responses
compose simple read responses where appropriate
```

The API Gateway must not:

```text
own order lifecycle
own payment logic
own stock reservation rules
write service databases
become a business rules monolith
```

## Risks

| Risk | Mitigation |
|---|---|
| Gateway accumulates domain logic | Keep business decisions in owning services |
| Gateway becomes single bottleneck | Scale independently and monitor |
| Error mapping inconsistent | Use shared error model |
| Gateway hides service failures | Emit dependency metrics and traces |

## Review Triggers

Revisit this decision if:

```text
multiple client types need distinct API shapes
gateway becomes overloaded with business logic
direct gRPC clients become the main usage model
```

## Related Documents

```text
docs/architecture/communication-patterns.md
docs/api/error-model.md
docs/api/grpc-overview.md
docs/architecture/service-boundaries.md
```
