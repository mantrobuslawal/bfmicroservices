# ADR-0007: Use OpenTelemetry for Observability

## Status

Accepted

## Date

2026-05-26

## Context

bfstore is a distributed microservice system.

A single user journey may cross:

```text
API Gateway
Order Service
Basket Service
Inventory Service
Payment Service
Shipping Service
Kafka
Notification Service
```

Without consistent observability, failures in checkout, events, retries, and downstream consumers would be difficult to diagnose.

## Decision

bfstore will use OpenTelemetry as the standard instrumentation approach for traces, metrics, and correlation context.

Structured logs will include correlation and trace identifiers.

## Drivers

This decision supports:

```text
end-to-end checkout tracing
service-level metrics
dependency visibility
Kafka producer/consumer observability
incident diagnosis
client-facing operational maturity
portable observability tooling
```

## Alternatives Considered

### Option 1: Logs Only

Benefits:

```text
simple
low initial setup
```

Costs:

```text
hard to trace distributed flows
weak dependency visibility
difficult root cause analysis
```

### Option 2: Vendor-Specific Observability SDK

Benefits:

```text
deep integration with one platform
```

Costs:

```text
vendor lock-in
less portable portfolio evidence
```

### Option 3: OpenTelemetry

Benefits:

```text
open standard
broad ecosystem support
supports traces and metrics
portable across backends
strong industry relevance
```

Costs:

```text
requires instrumentation discipline
local observability stack adds setup complexity
```

## Consequences

### Positive

```text
checkout can be traced across services
Kafka events can preserve correlation context
service metrics can support dashboards and alerts
observability design remains vendor-neutral
```

### Negative

```text
more implementation work
need consistent instrumentation wrappers
trace propagation through Kafka must be designed
local observability may require extra containers
```

## Implementation Notes

Services should emit:

```text
structured logs
request metrics
dependency metrics
Kafka publish/consume metrics
distributed traces
health and readiness status
```

Important context:

```text
correlation_id
trace_id
request_id
event_id
event_type
service
operation
```

Potential local stack:

```text
OpenTelemetry Collector
Prometheus
Grafana
Loki
Tempo
```

## Risks

| Risk | Mitigation |
|---|---|
| Inconsistent instrumentation | Provide shared telemetry package |
| Sensitive data in logs | Maintain secure logging standards |
| High cardinality metrics | Define allowed labels |
| Tracing not propagated through Kafka | Include trace/correlation metadata in event envelope |

## Review Triggers

Revisit this decision if:

```text
OpenTelemetry tooling becomes too heavy for early implementation
a target platform mandates a different approach
instrumentation overhead becomes unacceptable
```

## Related Documents

```text
docs/observability/logging.md
docs/observability/metrics.md
docs/observability/tracing.md
docs/architecture/communication-patterns.md
docs/architecture/resilience-patterns.md
docs/events/event-envelope.md
```
