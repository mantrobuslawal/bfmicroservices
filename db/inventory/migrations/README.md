# `db/inventory/migrations`

## 1. Purpose

This directory contains database migrations for the Inventory Service schema.

Inventory Service owns stock levels, stock reservations, reservation release, reservation expiry, and inventory adjustment history.

---

## 2. Owning Service

```text
inventory-service
```

Owned schema:

```text
bfstore_inventory
```

Only Inventory Service migrations should modify this schema.

---

## 3. Expected Migration Files

Recommended initial migrations:

```text
db/inventory/migrations/
├── README.md
├── 000001_create_warehouses.up.sql
├── 000001_create_warehouses.down.sql
├── 000002_create_stock_levels.up.sql
├── 000002_create_stock_levels.down.sql
├── 000003_create_stock_reservations.up.sql
├── 000003_create_stock_reservations.down.sql
├── 000004_create_stock_reservation_items.up.sql
├── 000004_create_stock_reservation_items.down.sql
└── 000005_create_stock_adjustments.up.sql
```

`warehouses` may be simplified or deferred for the first version.

---

## 4. Candidate Tables

Initial priority:

```text
stock_levels
stock_reservations
stock_reservation_items
```

Later:

```text
warehouses
stock_adjustments
outbox_events
```

---

## 5. Core Data Ownership

Inventory Service owns:

```text
available stock quantity
reserved stock quantity
stock reservations
reservation status
reservation expiry
stock release
stock adjustment history
```

Inventory Service does not own:

```text
product name
product description
basket contents
order lifecycle
payment state
shipment state
```

Inventory may store `product_id` and `variant_id` as references only.

---

## 6. Initial Table Design Notes

### `stock_levels`

Recommended fields:

```text
stock_level_id
product_id
variant_id
warehouse_id
available_quantity
reserved_quantity
created_at
updated_at
```

### `stock_reservations`

Recommended fields:

```text
reservation_id
order_id
basket_id
customer_id
status
idempotency_key
request_hash
expires_at
created_at
updated_at
```

### `stock_reservation_items`

Recommended fields:

```text
reservation_item_id
reservation_id
product_id
variant_id
quantity
created_at
```

---

## 7. Indexing Guidance

Recommended indexes:

```text
idx_stock_levels_product_variant
idx_stock_reservations_order_id
idx_stock_reservations_basket_id
idx_stock_reservations_status_expires_at
uq_stock_reservations_idempotency_key
idx_stock_reservation_items_reservation_id
```

---

## 8. Constraints and Invariants

Inventory migrations should support:

```text
available_quantity >= 0
reserved_quantity >= 0
reservation_id is unique
idempotency_key is unique for reservation operations
quantity > 0
reservation status is present
expires_at is present for active reservations
```

Stock must not become negative.

Concurrency control should be enforced through a combination of:

```text
transactions
row-level locking
constraints
idempotency keys
careful repository logic
```

---

## 9. Event and Outbox Considerations

Inventory Service should publish:

```text
StockReserved
StockReservationFailed
StockReservationReleased
StockReservationExpired
StockCommitted
InventoryAdjusted
```

If reliable event publication is required, add:

```text
outbox_events
```

Recommended for serious implementation once reservation events become critical.

---

## 10. Seed Data

Local seed data should include:

```text
products with available stock
products with low stock
products with zero stock
optional warehouse records
```

The product IDs should align with local Catalogue seed data, but this must not be enforced through cross-service foreign keys.

---

## 11. Migration Safety Rules

```text
do not create foreign keys to bfstore_catalog
do not create order or payment tables
do not allow negative stock values
do not remove idempotency support from reservation tables
do not edit applied migrations
```

---

## 12. Testing Expectations

Inventory migrations should be validated by tests for:

```text
migrations apply cleanly
stock levels can be inserted and queried
reservation uniqueness works
quantity constraints work
insufficient stock cannot oversell
concurrent reservation behaviour is safe
reservation release is idempotent
```

---

## 13. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/events/ordering-and-idempotency.md
proto/acme/inventory/v1/README.md
```

---

## 14. Summary

Inventory migrations define the stock control model for bfstore.

This schema is critical to checkout safety and must protect against overselling, duplicate reservations, and permanently locked stock.
