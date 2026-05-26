# `db/catalog/migrations`

## 1. Purpose

This directory contains database migrations for the Catalogue Service schema.

Catalogue Service owns governed product catalogue data, including products, variants, categories, category-scoped attribute definitions, product attribute values, product status, and product price data in the initial version.

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

## 3. Design Direction

The catalogue model must support varied future product types such as:

```text
curtains
bed frames
mattresses
sofas
rugs
lamps
wardrobes
tables
homeware
```

The schema should avoid:

```text
one huge products table with hundreds of nullable type-specific columns
uncontrolled JSON blobs with weak governance
```

Recommended approach:

```text
relational core product tables
category taxonomy
category-scoped attribute definitions
product/variant attribute values
denormalised Search Service projection for browse/filter/search
```

---

## 4. Expected Migration Files

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
├── 000004_create_product_attribute_definitions.up.sql
├── 000004_create_product_attribute_definitions.down.sql
├── 000005_create_product_attribute_values.up.sql
├── 000005_create_product_attribute_values.down.sql
├── 000006_create_product_attribute_options.up.sql
├── 000006_create_product_attribute_options.down.sql
├── 000007_create_product_images.up.sql
├── 000007_create_product_images.down.sql
├── 000008_create_product_price_history.up.sql
└── 000008_create_product_price_history.down.sql
```

`product_attribute_options`, `product_images`, and `product_price_history` may be deferred if not required for the first implementation.

---

## 5. Candidate Tables

Initial priority:

```text
categories
products
product_variants
product_attribute_definitions
product_attribute_values
```

Later:

```text
product_attribute_options
product_images
product_price_history
product_status_history
outbox_events
```

---

## 6. Core Data Ownership

Catalogue Service owns:

```text
product_id
product name
product description
product category
category taxonomy
product variant metadata
product active/inactive status
category-scoped attribute definitions
product attribute values
catalogue price
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

## 7. Initial Table Design Notes

## 7.1 `categories`

Recommended fields:

```text
category_id
parent_category_id
name
slug
description
status
created_at
updated_at
```

Rules:

```text
category_id is stable
slug should be unique where used publicly
categories define the scope for product attributes
```

## 7.2 `products`

Recommended fields:

```text
product_id
category_id
name
description
status
base_price_minor
currency_code
brand
created_at
updated_at
```

Rules:

```text
products contain common product data only
type-specific characteristics belong in product attributes
only ACTIVE products are purchasable
```

Avoid columns such as:

```text
curtain_drop_cm
bed_size
bulb_type
rug_shape
mattress_firmness
sofa_orientation
```

## 7.3 `product_variants`

Recommended fields:

```text
variant_id
product_id
sku
variant_name
price_minor
currency_code
status
created_at
updated_at
```

Examples:

```text
curtain width/drop variants
bed size variants
sofa fabric variants
rug size variants
```

## 7.4 `product_attribute_definitions`

Recommended fields:

```text
attribute_id
category_id
code
display_name
description
data_type
unit
is_required
is_filterable
is_variant_defining
allowed_values_json
display_order
status
created_at
updated_at
```

Rules:

```text
attribute codes should be stable
attribute definitions are scoped to a category
filterable attributes should be identified for Search Service
required attributes should be validated by Catalogue Service
```

Example rows:

| Category | Code | Data Type | Unit | Filterable |
|---|---|---|---|---|
| curtains | `drop_cm` | number | cm | yes |
| curtains | `heading_type` | option | none | yes |
| bed-frames | `bed_size` | option | none | yes |
| bed-frames | `storage_type` | option | none | yes |
| rugs | `shape` | option | none | yes |
| lamps | `bulb_type` | option | none | yes |

## 7.5 `product_attribute_values`

Recommended fields:

```text
product_attribute_value_id
product_id
variant_id
attribute_id
value_string
value_number
value_boolean
value_json
unit
created_at
updated_at
```

Rules:

```text
only one typed value column should be populated
attribute definition determines the expected type
variant_id is nullable unless the value differs per variant
```

## 7.6 `product_attribute_options`

Recommended fields:

```text
attribute_option_id
attribute_id
value
display_name
display_order
status
created_at
updated_at
```

This table is useful for controlled values such as:

```text
single
double
king
super_king
eyelet
pencil_pleat
ottoman
```

For early implementation, controlled values may temporarily live in `allowed_values_json` if clearly documented.

---

## 8. Indexing Guidance

Recommended indexes:

```text
idx_products_status
idx_products_category_id
idx_products_category_status
idx_product_variants_product_id
uq_product_variants_sku
idx_categories_parent_category_id
uq_categories_slug
uq_product_attribute_definitions_category_code
idx_product_attribute_definitions_category_id
idx_product_attribute_definitions_filterable
idx_product_attribute_values_product_id
idx_product_attribute_values_variant_id
idx_product_attribute_values_attribute_id
idx_product_attribute_options_attribute_id
```

Search/filter-heavy browse should eventually use Search Service rather than complex Catalogue SQL queries.

---

## 9. Constraints and Invariants

Catalogue migrations should enforce:

```text
product_id is unique
variant_id is unique
sku is unique where required
amount fields are not negative
currency_code is present
product status is present
category_id is present
attribute definition code is unique per category
attribute value references an attribute definition in the catalogue schema
created_at is present
```

Catalogue should not create foreign keys into other service schemas.

---

## 10. Event and Outbox Considerations

Catalogue Service may publish:

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
ProductArchived
CategoryCreated
CategoryUpdated
ProductAttributeDefinitionCreated
ProductAttributeDefinitionUpdated
ProductAttributeDefinitionDeprecated
```

If Search Service relies on catalogue events, consider:

```text
outbox_events
```

for reliable product event publication.

---

## 11. Seed Data

Local seed data should include multiple product types with different attributes.

Examples:

```text
active curtain product with drop_cm, width_cm, lining, heading_type
active bed frame with bed_size, material, storage_type, slat_type
active rug with shape, material, pile_height_cm
active lamp with bulb_type, wattage, fitting_type
inactive product
out-of-stock product reference
```

Seed data must be fictional and safe for public repositories.

---

## 12. Migration Safety Rules

```text
do not edit migrations after they have been applied
do not reference other service schemas
do not store stock in catalogue tables
do not use FLOAT or DOUBLE for money
do not add many nullable product-type-specific columns to products
do not hide all product data in ungoverned JSON
reserve destructive changes for explicit reviewed migrations
```

---

## 13. Testing Expectations

Catalogue migrations should be validated by tests for:

```text
migrations apply cleanly
products can be inserted and queried
inactive products can be filtered out
money uses minor units
required constraints are enforced
attribute definitions can be created per category
attribute values can be attached to products
attribute type validation works in service logic
filterable attributes can be identified for Search Service
repository queries use indexes where appropriate
```

---

## 14. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/requirements/business-rules.md
proto/acme/catalog/v1/README.md
```

---

## 15. Summary

Catalogue migrations define the governed product data model.

The updated design keeps MySQL as the Catalogue Service source of truth while supporting varied product types through category-scoped attributes.

Search Service should consume catalogue events and maintain denormalised search documents for browse, filters, and facets.
