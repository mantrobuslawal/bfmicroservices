# Parity Testing

This document defines testing practices that support dev/prod parity for bfstore.

---

## Purpose

This document explains:

```text
integration tests with real MySQL/Kafka
contract tests with Buf
gRPC smoke tests
environment parity checks
```

---

## Core Rule

```text
Mocks are useful for speed.
Real backing services are useful for truth.
```

---

## Unit Tests

Unit tests should cover:

```text
domain logic
validation
mapping
small pure functions
error handling
```

Mocks are appropriate here.

---

## Integration Tests

Integration tests should use real backing service types where behaviour matters.

Examples:

```text
MySQL repository tests
Kafka producer/consumer tests
gRPC handler/client tests
outbox publishing tests
payment simulator boundary tests
```

---

## Contract Tests

Contract checks should include:

```text
buf lint
buf breaking
protobuf generation checks
gRPC compatibility checks where practical
```

This protects Protobuf/gRPC parity.

---

## Smoke Tests

Smoke tests should verify deployed environment basics.

Examples:

```text
services are healthy
api-gateway reachable
catalog-service gRPC health passes
MySQL migrations applied
Kafka topics available
notification-worker can consume
OTel telemetry is emitted
```

---

## End-to-end Tests

Eventually test the core business flow:

```text
browse products
add item to basket
checkout
reserve inventory
authorise payment
create shipment
create order
publish OrderCreated
consume notification event
```

---

## Parity Checks

Check that local/staging/prod-style environments use:

```text
same config names
same service protocols
same backing service types
same event serialization
same health check model
same telemetry model
```

---

## Practical Rules

```text
Test with MySQL, not just mocks.
Test with Kafka, not just fake queues.
Run Buf checks in CI.
Use grpcurl or automated gRPC smoke tests.
Check telemetry locally and in staging.
Document any test environment shortcuts.
```

---

## Final Rule

```text
Parity testing makes environment differences visible before they become incidents.
```
