# ADR-0006: Use Buf for Protobuf Management

## Status

Accepted

## Date

2026-05-26

## Context

bfstore uses Protocol Buffers for:

```text
gRPC service APIs
Kafka event payloads
shared common messages
event metadata
```

The project needs consistent linting, code generation, and compatibility checks.

---

## Decision

bfstore will use **Buf** to manage Protobuf definitions.

Buf will be used for:

```text
linting
breaking-change detection
code generation
contract governance
```

This applies to both gRPC API contracts and Kafka event payload contracts.

---

## Drivers

```text
contract-first service design
consistent Protobuf style
generated Go code
safe schema evolution
event compatibility checks
CI quality gates
developer workflow consistency
```

---

## Implementation Notes

Recommended files:

```text
buf.yaml
buf.gen.yaml
```

Recommended commands:

```sh
buf lint
buf breaking
buf generate
```

Recommended proto layout:

```text
proto/acme/common/v1/
proto/acme/events/v1/
proto/acme/catalog/v1/
proto/acme/catalog/events/v1/
proto/acme/order/v1/
proto/acme/order/events/v1/
proto/acme/payment/v1/
proto/acme/payment/events/v1/
```

---

## Event Contract Governance

Kafka event payloads are Protobuf contracts and should be treated as seriously as gRPC APIs.

Buf should check:

```text
event metadata messages
domain event messages
shared common types
service API messages
```

CI should fail when incompatible event changes are introduced without an intentional versioning plan.

---

## Summary

bfstore uses Buf to manage Protobuf contracts for both gRPC APIs and Kafka event payloads. This strengthens the project’s contract-first architecture.


---

## Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/events/event-envelope.md
docs/events/event-catalog.md
docs/events/kafka-topic-design.md
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
```
