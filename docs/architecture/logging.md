# Logging

This document defines bfstore's logging architecture.

---

## Purpose

This document explains:

```text
logs as event streams
stdout/stderr policy
structured JSON logs
relationship to traces and metrics
service-specific log events
```

---

## Core Rule

```text
The app emits log events.
The platform captures and routes them.
```

Applications must not manage log files, rotation, retention, or vendor-specific routing.

---

## Logs as Event Streams

Logs are continuous runtime events.

Examples:

```text
service_started
checkout_started
payment_authorised
kafka_event_published
notification_sent
shutdown_started
```

Each event should be useful, safe, and structured.

---

## stdout/stderr Policy

Services should write logs to stdout/stderr.

Do not write application logs to:

```text
/var/log
local container files
service-specific mounted log directories
vendor-specific SDK destinations by default
```

The runtime/platform is responsible for capture and routing.

---

## Structured JSON

Production-style logs should be JSON.

Recommended common fields:

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

Service-specific fields may include:

```text
order_id
payment_attempt_id
basket_id
product_id
kafka_topic
kafka_partition
kafka_offset
grpc_method
grpc_status
http_route
```

---

## Logs, Traces, and Metrics

```text
logs:
  explain discrete events

traces:
  show request flow across services

metrics:
  show rates, trends, and thresholds
```

All three should be correlated where practical.

---

## Practical Rules

```text
Emit structured JSON logs.
Use stdout/stderr.
Use consistent field names.
Include correlation and trace IDs.
Avoid sensitive payloads.
Use metrics for high-volume trends.
Use logs for explanation and investigation.
```

---

## Final Rule

```text
bfstore logs should describe runtime behaviour without owning logging infrastructure.
```
