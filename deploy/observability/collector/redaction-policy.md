# Collector Redaction Policy

This document defines bfstore guidance for redacting sensitive telemetry in the OpenTelemetry Collector.

It complements:

```text
docs/architecture/collector-configuration-security.md
deploy/observability/collector/config-security.md
```

---

## Purpose

Redaction protects bfstore from accidentally forwarding sensitive telemetry to observability backends.

Redaction is a safety net, not the first line of defence.

---

## First Rule

```text
Do not emit sensitive data in the first place.
```

Services should avoid creating telemetry that contains sensitive or personal data.

---

## Never Emit

bfstore telemetry must not include:

```text
raw JWTs
passwords
API keys
payment card numbers
CVV
full shipping address
customer email
full basket JSON
raw Kafka payloads
full SQL with customer values
session cookies
authorisation headers
```

---

## Safe Attribute Examples

Allowed-style attributes:

```text
service.name
deployment.environment.name
rpc.service
rpc.method
grpc.status_code
db.system
db.operation
db.collection.name
messaging.system
messaging.destination.name
event.type
error.type
bfstore.checkout.stage
bfstore.payment.provider
bfstore.stock.reservation_result
```

These describe behaviour without exposing private data.

---

## Redaction Processor Role

A redaction processor may:

```text
remove attributes not on an allow list
mask blocked patterns
drop unsafe fields
protect against accidental leakage
```

Example policy shape:

```yaml
processors:
  redaction:
    allow_all_keys: false
    allowed_keys:
      - service.name
      - deployment.environment.name
      - rpc.service
      - rpc.method
      - grpc.status_code
      - db.system
      - db.operation
      - db.collection.name
      - messaging.system
      - messaging.destination.name
      - event.type
      - error.type
      - bfstore.checkout.stage
      - bfstore.payment.provider
      - bfstore.stock.reservation_result
```

This is a conceptual starting point and should be tested carefully before production-style use.

---

## Blocked Patterns

Consider blocking patterns that look like:

```text
credit card numbers
bearer tokens
API keys
email addresses
private keys
session cookies
authorisation headers
```

---

## Where Redaction Belongs

Best:

```text
avoid sensitive telemetry in service instrumentation
```

Good safety net:

```text
Collector redaction processor
```

Also useful:

```text
backend-side controls
log access controls
retention policies
```

Do not rely only on Collector redaction.

---

## Testing Expectations

Test redaction rules with sample telemetry containing:

```text
fake JWT
fake API key
fake email address
fake card-like number
safe rpc.method
safe service.name
safe bfstore.checkout.stage
```

Verify:

```text
unsafe values are removed or masked
safe attributes remain
redaction does not break dashboards
redaction does not remove required trace correlation
```

---

## Practical Rules

```text
Do not emit secrets.
Do not emit personal data.
Do not emit raw payloads.
Use redaction as a safety net.
Prefer allow lists for mature configs.
Test redaction with fake sensitive values.
Keep dashboards aligned with allowed attributes.
Review new telemetry attributes before rollout.
```

---

## Final Rule

```text
The best sensitive telemetry is the telemetry that never leaves the service.
```
