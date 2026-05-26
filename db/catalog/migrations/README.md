# `db/catalog/migrations`

## 1. Purpose

This directory contains database migrations for the Catalogue Service schema.

Catalogue Service owns product catalogue data, including products, variants, categories, product status, and product price data in the initial version.

---

## 2. Owning Service

```text
catalog-service
```

Owned schema:

```text
bfstore_catalog
```

Only Catalogue Service migrations should modify this schema.

---

## 3. Expected Migration Files

Recommended initial migrations:

```text
db/catalog/migrations/
├── README.md
├── 000001_create_categories.up.sql
├── 000001_create_categories.down.sql
├── 000002_create_products.up.sql
├── 000002_create_products.down.sql
├── 000003_create_product_variants.up.sql
├── 000003_create_product_variants.down.sql
├── 000004_create_product_images.up.sql
├── 000004_create_product_images.down.sql
└── 000005_create_product_price_history.up.sql
```

Product images and price history may be deferred if not needed for the first implementation.

---

## 4. Candidate Tables

Initial priority:

```text
categories
products
product_variants
```

Later:

```text
product_images
product_attributes
product_price_history
outbox_events
```

---

## 5. Core Data Ownership

Catalogue Service owns:

```text
product_id
product name
product description
product category
product variant metadata
product active/inactive status
initial product price
```

Catalogue Service does not own:

```text
stock quantity
stock reservation
basket contents
order item history
search ranking
recommendation outputs
```

---

## 6. Initial Table Design Notes

### `products`

Recommended fields:

```text
product_id
category_id
name
description
status
base_price_minor
currency_code
material
colour
created_at
updated_at
```

### `product_variants`

Recommended fields:

```text
variant_id
product_id
sku
size
colour
material
price_minor
currency_code
status
created_at
updated_at
```

### `categories`

Recommended fields:

```text
category_id
parent_category_id
name
slug
status
created_at
updated_at
```

---

## 7. Indexing Guidance

Recommended indexes:

```text
idx_products_status
idx_products_category_id
idx_products_category_status
idx_product_variants_product_id
uq_product_variants_sku
idx_categories_parent_category_id
```

Indexes should be based on real access patterns and reviewed as queries evolve.

---

## 8. Constraints and Invariants

Catalogue migrations should enforce:

```text
product_id is unique
variant_id is unique
sku is unique where required
amount fields are not negative
currency_code is present
product status is present
created_at is present
```

Catalogue should not create foreign keys into other service schemas.

---

## 9. Event and Outbox Considerations

Catalogue Service may later publish:

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
ProductArchived
```

If reliable publication is required, add:

```text
outbox_events
```

to this schema.

The first checkout vertical slice may not require catalogue event publishing immediately.

---

## 10. Seed Data

Local seed data should include:

```text
active products
inactive product
out-of-stock product reference
furniture categories
product variants
```

Seed data must be fictional and safe for public repositories.

---

## 11. Migration Safety Rules

```text
do not edit migrations after they have been applied
do not reference other service schemas
do not store stock in catalogue tables
do not use FLOAT or DOUBLE for money
reserve destructive changes for explicit reviewed migrations
```

---

## 12. Testing Expectations

Catalogue migrations should be validated by tests for:

```text
migrations apply cleanly
products can be inserted and queried
inactive products can be filtered out
money uses minor units
required constraints are enforced
repository queries use indexes where appropriate
```

---

## 13. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/requirements/business-rules.md
proto/acme/catalog/v1/README.md
```

---

## 14. Summary

Catalogue migrations define the product catalogue data model.

They must preserve Catalogue Service ownership and avoid leaking stock, order, basket, or search responsibilities into the catalogue schema.
