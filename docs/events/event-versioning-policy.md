# Event Versioning Policy

This document defines versioning rules for bfstore Kafka events.

## Core Rule

```text
Kafka consumers are API consumers too.
```

## Event contracts

Kafka event contracts include:

```text
topic name
message key semantics
Protobuf message schema
event type
event version
event meaning
ordering assumptions
idempotency fields
```

## Topic versioning

```text
bfstore.order.orders.v1
bfstore.order.orders.v2
```

## PATCH changes

```text
producer bug fix
consumer bug fix
logging fix
no schema or meaning change
```

## MINOR changes

```text
add optional field
add new event type
add new enum value carefully
add new metadata where old consumers ignore it
```

## MAJOR changes

```text
remove required field
change field type
change event meaning
change key semantics
change ordering assumptions
change topic naming/partitioning contract
```

## Event key compatibility

Changing the message key can break:

```text
ordering
partitioning
consumer assumptions
idempotency
replay behaviour
```

Treat key changes as major unless proven otherwise.

## Practical rules

```text
Version event topics deliberately.
Use Protobuf compatibility rules.
Keep event meanings stable.
Do not silently change keys.
Preserve idempotency fields.
Document consumer impact.
```

## Final rule

```text
Event versioning protects asynchronous consumers from surprise breakage.
```
