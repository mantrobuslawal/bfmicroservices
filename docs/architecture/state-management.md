# State Management

This document defines where bfstore state should live.

---

## Purpose

This document explains:

```text
where product, basket, order, payment, notification, event, and image state lives
what must never live only in memory
what may be cached temporarily
how future backing services fit in
```

---

## Core Rule

```text
The system has state.
The process should not be the state.
```

---

## Product State

Source of truth:

```text
catalog MySQL
```

Includes:

```text
products
categories
variants
attributes
product image metadata
```

Future binary image storage:

```text
object storage
```

Do not store product truth only in process memory.

---

## Basket State

Source of truth:

```text
basket MySQL
```

Possible later optimisation:

```text
Redis with TTL if deliberately chosen
```

Do not store basket contents only in a Go map or process memory.

---

## Order State

Source of truth:

```text
order MySQL
```

Includes:

```text
orders
order items
checkout state
order status transitions
outbox events
```

Important transitions should be persisted.

---

## Payment State

Source of truth:

```text
payment MySQL
```

Includes:

```text
payment attempts
authorisation results
provider references
idempotency keys
failure reasons where safe
```

Payment state must survive restarts and retries.

---

## Inventory State

Source of truth:

```text
inventory MySQL
```

Includes:

```text
stock levels
reservations
reservation status
reservation expiry if used
```

Inventory changes must be transactional and recoverable.

---

## Shipping State

Source of truth:

```text
shipping MySQL
```

Includes:

```text
shipment records
shipment status
carrier references later
```

---

## Notification State

Possible source of truth:

```text
Kafka offsets
notification MySQL if delivery tracking is needed
```

Includes:

```text
event_id
order_id
notification_type
delivery_status
sent_at
attempt_count
```

Do not rely only on in-memory sent/deduplication maps.

---

## Event State

Event stream:

```text
Kafka
```

Critical event production may use:

```text
service-owned outbox table
Kafka publisher
idempotent producer/consumer patterns
```

Events should be replay-safe where practical.

---

## Telemetry State

Telemetry should be sent to:

```text
OpenTelemetry Collector
metrics backend
trace backend
log backend
```

Do not store important diagnostic history only inside a process.

---

## Temporary State

Acceptable temporary state:

```text
request DTOs
parsed config
DB connection pools
short-lived local variables
temporary files during image processing
safe-to-lose cache
```

Temporary state must be rebuildable or disposable.

---

## What Must Never Live Only in Memory

```text
basket contents
orders
payment attempts
inventory reservations
shipment status
notification sent status
customer session state if needed
product truth
uploaded images
Kafka processing progress that affects correctness
```

---

## Practical Rules

```text
Use service-owned databases for durable business state.
Use Kafka for event streams and replay.
Use object storage later for product assets.
Use Redis later only as explicit backing service.
Keep process memory temporary.
Keep container filesystem temporary.
Use idempotency for retries and replays.
```

---

## Final Rule

```text
If the business must remember it, store it in a backing service.
```
