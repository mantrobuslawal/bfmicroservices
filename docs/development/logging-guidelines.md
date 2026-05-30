# Logging Guidelines

This document defines developer logging guidelines for bfstore services.

---

## Purpose

This document explains:

```text
log levels
field naming
correlation IDs
examples for Go slog
what to log and not log
local logging workflow
```

---

## Core Rule

```text
Log the right event, at the right level, with the right fields.
```

---

## Log Levels

```text
DEBUG:
  detailed troubleshooting

INFO:
  meaningful lifecycle and business events

WARN:
  abnormal but recoverable situations

ERROR:
  failed operations
```

---

## Common Fields

Use consistent names across services:

```text
timestamp
level
service
environment
version
event
message
trace_id
span_id
correlation_id
request_id
error_code
duration_ms
```

---

## Business Context Fields

Use safe IDs where useful:

```text
order_id
basket_id
product_id
payment_attempt_id
event_id
notification_type
checkout_id
```

Do not log sensitive payloads.

---

## Go slog Example

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

logger.Info(
    "order created",
    "event", "order_created",
    "order_id", orderID,
    "checkout_id", checkoutID,
    "trace_id", traceID,
)
```

---

## Local Workflow

Useful local commands:

```bash
docker compose logs -f catalog-service
docker compose logs -f order-service
docker compose logs -f notification-worker
```

Optional local pretty logging may be allowed, but JSON should remain the production-style default.

---

## What To Log

Log:

```text
service startup/readiness/shutdown
important business events
external calls starting/completing/failing
retry events
Kafka publish/consume milestones
unexpected failures
slow operations where useful
```

---

## What Not To Log

Avoid:

```text
full request/response bodies by default
secrets
tokens
payment details
customer private data
huge payloads
low-value function entry/exit noise at INFO
```

---

## Practical Rules

```text
Use INFO to tell the business/runtime story.
Use WARN for risk and degradation.
Use ERROR for failed operations.
Use DEBUG for investigation.
Keep fields consistent.
Make logs joinable through trace_id/correlation_id.
Prefer IDs to full objects.
```

---

## Final Rule

```text
Every log should earn its place.
```
