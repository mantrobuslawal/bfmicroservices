# `db/shipping/migrations`

## 1. Purpose

This directory contains database migrations for the Shipping Service schema.

Shipping Service owns delivery options, shipment state, tracking references, delivery address snapshots, and shipment-related events.

---

## 2. Owning Service

```text
shipping-service
```

Owned schema:

```text
bfstore_shipping
```

Only Shipping Service migrations should modify this schema.

---

## 3. Expected Migration Files

Recommended initial migrations:

```text
db/shipping/migrations/
├── README.md
├── 000001_create_delivery_options.up.sql
├── 000001_create_delivery_options.down.sql
├── 000002_create_shipments.up.sql
├── 000002_create_shipments.down.sql
├── 000003_create_shipment_status_history.up.sql
├── 000003_create_shipment_status_history.down.sql
├── 000004_create_tracking_events.up.sql
└── 000005_create_outbox_events.up.sql
```

Tracking and outbox may be deferred depending on implementation phase.

---

## 4. Candidate Tables

Initial priority:

```text
delivery_options
shipments
```

Later:

```text
shipment_status_history
tracking_events
outbox_events
```

---

## 5. Core Data Ownership

Shipping Service owns:

```text
delivery option definitions
shipment_id
shipment status
shipment carrier
tracking reference
shipment failure state
delivery address snapshot
tracking events
```

Shipping Service does not own:

```text
order lifecycle
payment state
customer saved addresses
stock state
notification delivery
```

Shipping stores `order_id` and `customer_id` as references only.

---

## 6. Initial Table Design Notes

### `delivery_options`

Recommended fields:

```text
delivery_option_id
name
description
price_minor
currency_code
estimated_days_min
estimated_days_max
status
created_at
updated_at
```

### `shipments`

Recommended fields:

```text
shipment_id
order_id
customer_id
status
delivery_option_id
carrier
tracking_reference
delivery_address_snapshot_json
idempotency_key
request_hash
created_at
updated_at
dispatched_at
delivered_at
```

### `tracking_events`

Recommended fields:

```text
tracking_event_id
shipment_id
status
description
occurred_at
created_at
```

---

## 7. Indexing Guidance

Recommended indexes:

```text
idx_shipments_order_id
idx_shipments_customer_id
idx_shipments_status
uq_shipments_idempotency_key
idx_tracking_events_shipment_id
idx_delivery_options_status
```

---

## 8. Constraints and Invariants

Shipping migrations should support:

```text
shipment_id is unique
idempotency_key is unique where required
delivery option status is present
shipment status is present
delivery option price is not negative
currency_code is present
```

Foreign keys within the Shipping schema are acceptable.

Cross-service foreign keys are not allowed.

---

## 9. Address Snapshot Rules

Shipping may store delivery address snapshots.

Rules:

```text
snapshot only what is required for fulfilment
do not update historical snapshots when customer address changes
do not log full address unnecessarily
avoid including full address in Kafka events unless justified
```

---

## 10. Event and Outbox Considerations

Shipping Service should publish:

```text
ShipmentCreated
ShipmentFailed
ShipmentDispatched
ShipmentDelivered
ShipmentCancelled
```

Outbox is recommended if shipment events become critical to order state or notification flows.

---

## 11. Seed Data

Local seed data should include:

```text
standard delivery option
express delivery option
disabled delivery option
```

Carrier integrations may be simulated initially.

---

## 12. Migration Safety Rules

```text
do not create foreign keys to order or customer schemas
do not store payment data
do not store customer profile truth
do not remove idempotency support
do not expose real addresses in seed data
do not edit applied migrations
```

---

## 13. Testing Expectations

Shipping migrations should be validated by tests for:

```text
migrations apply cleanly
delivery options can be inserted
shipment can be created
idempotency uniqueness works
shipment status is stored
address snapshot is stored safely
shipment lookup by order_id works
```

---

## 14. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/requirements/business-rules.md
proto/acme/shipping/v1/README.md
```

---

## 15. Summary

Shipping migrations define delivery and fulfilment persistence.

They must support idempotent shipment creation, clear shipment state, safe address snapshots, and service-owned data boundaries.
