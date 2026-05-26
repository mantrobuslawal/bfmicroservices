# `db/order/migrations`

## 1. Purpose

This directory contains database migrations for the Order Service schema.

Order Service owns order lifecycle, checkout attempts, order item snapshots, order status history, and order-related events.

---

## 2. Owning Service

```text
order-service
```

Owned schema:

```text
bfstore_order
```

Only Order Service migrations should modify this schema.

---

## 3. Expected Migration Files

Recommended initial migrations:

```text
db/order/migrations/
├── README.md
├── 000001_create_checkout_attempts.up.sql
├── 000001_create_checkout_attempts.down.sql
├── 000002_create_orders.up.sql
├── 000002_create_orders.down.sql
├── 000003_create_order_items.up.sql
├── 000003_create_order_items.down.sql
├── 000004_create_order_status_history.up.sql
├── 000004_create_order_status_history.down.sql
└── 000005_create_outbox_events.up.sql
```

`outbox_events` is strongly recommended for serious implementation because `OrderCreated` is a critical event.

---

## 4. Candidate Tables

Initial priority:

```text
checkout_attempts
orders
order_items
outbox_events
```

Later:

```text
order_status_history
order_failures
order_adjustments
```

---

## 5. Core Data Ownership

Order Service owns:

```text
order_id
order number
order status
order totals
order item snapshots
checkout attempt status
order status history
order events
```

Order Service does not own:

```text
stock levels
payment provider state
shipment carrier state
notification delivery state
product truth
customer profile truth
```

Order may store references and snapshots where required.

---

## 6. Initial Table Design Notes

### `checkout_attempts`

Recommended fields:

```text
checkout_attempt_id
customer_id
basket_id
idempotency_key
request_hash
status
failure_reason
created_at
updated_at
completed_at
```

### `orders`

Recommended fields:

```text
order_id
order_number
customer_id
basket_id
status
total_amount_minor
currency_code
delivery_address_snapshot_json
idempotency_key
created_at
updated_at
confirmed_at
cancelled_at
```

### `order_items`

Recommended fields:

```text
order_item_id
order_id
product_id
variant_id
product_name_snapshot
sku_snapshot
unit_price_minor
currency_code
quantity
line_total_minor
created_at
```

### `outbox_events`

Recommended fields:

```text
outbox_event_id
event_id
event_type
event_version
aggregate_type
aggregate_id
payload
status
attempt_count
next_attempt_at
created_at
published_at
last_error
```

---

## 7. Indexing Guidance

Recommended indexes:

```text
uq_orders_order_number
uq_orders_idempotency_key
idx_orders_customer_created_at
idx_orders_status
idx_order_items_order_id
idx_checkout_attempts_customer_basket
uq_checkout_attempts_idempotency_key
idx_outbox_events_status_next_attempt_at
```

---

## 8. Constraints and Invariants

Order migrations should support:

```text
order_id is unique
order_number is unique
idempotency_key is unique where required
quantity > 0
amount fields are not negative
currency_code is present
status is present
order item belongs to order within order schema
```

Order must not create foreign keys to Basket, Catalogue, Inventory, Payment, Shipping, or Customer schemas.

---

## 9. Idempotency Requirements

Order Service must prevent duplicate checkout effects.

Database support should include:

```text
unique idempotency key
request hash
checkout attempt record
order result reference
```

Expected behaviour:

```text
same idempotency key + same request returns original result
same idempotency key + different request is rejected
```

---

## 10. Event and Outbox Considerations

Order Service should publish:

```text
OrderCreated
OrderFailed
OrderCancelled
OrderConfirmed
```

For `OrderCreated`, use outbox pattern where practical:

```text
create order
create order_items
create outbox event
commit transaction
publish event later
```

This prevents lost events after successful order creation.

---

## 11. Snapshot Rules

Order items should preserve historical snapshots:

```text
product_name_snapshot
sku_snapshot
unit_price_minor
currency_code
quantity
line_total_minor
```

Customer address may be snapshotted if required for order history.

Snapshots should be clearly named and should not become live source-of-truth data.

---

## 12. Migration Safety Rules

```text
do not create foreign keys to other service schemas
do not store payment card data
do not store live stock quantities
do not edit applied migrations
do not remove idempotency support
do not drop snapshot fields without migration plan
```

---

## 13. Testing Expectations

Order migrations should be validated by tests for:

```text
migrations apply cleanly
order can be created
order items can be created
idempotency uniqueness works
duplicate checkout is prevented
order item snapshots are preserved
outbox event can be inserted with order transaction
```

---

## 14. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/events/event-envelope.md
docs/events/ordering-and-idempotency.md
adr/0011-use-outbox-pattern-for-critical-events.md
proto/acme/order/v1/README.md
```

---

## 15. Summary

Order migrations define the central checkout result model.

This schema is critical to business correctness and must protect idempotency, historical snapshots, order state, and reliable event publication.
