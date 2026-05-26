# `db/basket/migrations`

## 1. Purpose

This directory contains database migrations for the Basket Service schema.

Basket Service owns shopping basket state before checkout, including baskets and basket items.

---

## 2. Owning Service

```text
basket-service
```

Owned schema:

```text
bfstore_basket
```

Only Basket Service migrations should modify this schema.

---

## 3. Expected Migration Files

Recommended initial migrations:

```text
db/basket/migrations/
├── README.md
├── 000001_create_baskets.up.sql
├── 000001_create_baskets.down.sql
├── 000002_create_basket_items.up.sql
├── 000002_create_basket_items.down.sql
└── 000003_create_basket_events_outbox.up.sql
```

The outbox table may be deferred unless basket events are implemented early.

---

## 4. Candidate Tables

Initial priority:

```text
baskets
basket_items
```

Later:

```text
basket_events_outbox
basket_status_history
```

---

## 5. Core Data Ownership

Basket Service owns:

```text
basket_id
basket status
basket items
basket item quantities
customer/session basket association
basket expiry
checked-out marker
```

Basket Service does not own:

```text
stock reservation
final order state
payment state
shipment state
product truth
```

Basket may store product IDs and price display snapshots, but Catalogue remains product truth.

---

## 6. Initial Table Design Notes

### `baskets`

Recommended fields:

```text
basket_id
customer_id
session_id
status
created_at
updated_at
expires_at
checked_out_at
```

### `basket_items`

Recommended fields:

```text
basket_item_id
basket_id
product_id
variant_id
quantity
unit_price_snapshot_minor
currency_code
created_at
updated_at
```

Price snapshots are for basket display only. Final order prices should be confirmed or snapshotted by Order Service.

---

## 7. Indexing Guidance

Recommended indexes:

```text
idx_baskets_customer_id
idx_baskets_session_id
idx_baskets_status
idx_baskets_expires_at
idx_basket_items_basket_id
uq_basket_items_basket_product_variant
```

---

## 8. Constraints and Invariants

Basket migrations should support:

```text
basket_id is unique
basket item quantity > 0
basket status is present
basket item belongs to a basket within the basket schema
duplicate product/variant lines are controlled where required
```

Foreign keys within the Basket schema are acceptable.

Cross-service foreign keys are not allowed.

---

## 9. Event and Outbox Considerations

Basket Service may publish:

```text
BasketCreated
BasketItemAdded
BasketItemUpdated
BasketItemRemoved
BasketCheckedOut
BasketExpired
```

These may be useful later for recommendation signals and abandoned basket analytics.

They are not required for the first checkout implementation unless specifically chosen.

---

## 10. Seed Data

Local seed data may include:

```text
active basket
basket with one item
basket with multiple items
checked-out basket
expired basket
```

Product IDs should match local Catalogue seed data but without cross-schema constraints.

---

## 11. Migration Safety Rules

```text
do not reserve stock in basket tables
do not store payment data
do not store final order state
do not create foreign keys to catalogue or inventory schemas
do not allow zero or negative item quantities
```

---

## 12. Testing Expectations

Basket migrations should be validated by tests for:

```text
migrations apply cleanly
basket can be created
item can be added
quantity constraint works
checked-out basket cannot be mutated at service logic level
basket item uniqueness behaves as expected
```

---

## 13. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/requirements/business-rules.md
proto/acme/basket/v1/README.md
```

---

## 14. Summary

Basket migrations define the pre-checkout shopping intent model.

They must keep Basket Service focused on baskets and avoid leaking stock reservation, order, or payment responsibilities into the basket schema.
