# Catalogue Database

## 1. Purpose

This directory contains the local database artefacts for the **Catalogue Service**.

Catalogue Service owns the product catalogue domain for bfstore.

It is responsible for:

```text
products
categories
product variants
category-scoped product attributes
product imagery
catalogue outbox events
```

It is not responsible for:

```text
stock levels
basket state
orders
payments
shipping
search indexes
recommendation models
```

---

## 2. Directory Layout

```text
db/catalog/
├── README.md
├── migrations/
│   ├── 000001_create_catalog_schema.up.sql
│   └── 000001_create_catalog_schema.down.sql
└── seeds/
    └── 001_seed_borough_products.sql
```

---

## 3. Design Summary

bfstore sells varied developer-themed homeware.

Different product categories need different attributes.

Examples:

```text
lamps need bulb type and max wattage
wall art needs size and material
lockboxes need lock type and security rating
rugs need shape and pile height
soft furnishings need fabric type and care instructions
```

The catalogue schema uses:

```text
relational product core
category-scoped attribute definitions
typed product attribute values
variant support
controlled attribute options
```

This avoids:

```text
one giant products table with hundreds of nullable columns
uncontrolled schemaless JSON blobs
mixing stock/order/payment data into the catalogue database
```

---

## 4. Migration

Initial schema migration:

```text
db/catalog/migrations/000001_create_catalog_schema.up.sql
```

Rollback migration:

```text
db/catalog/migrations/000001_create_catalog_schema.down.sql
```

The schema includes:

```text
categories
products
product_variants
product_attribute_definitions
product_attribute_options
product_attribute_values
product_images
catalogue_outbox_events
```

---

## 5. Seed Data

Seed data:

```text
db/catalog/seeds/001_seed_borough_products.sql
```

Example products:

```text
Gopher Desk Lamp
Gopher Cushion Set
Rob Pike Wall Tapestry
Rivest Super-Secure Lockbox
Dijkstra Pathfinding Rug
Grace Hopper Debugging Blanket
```

The seed data is intentionally memorable, but it exercises serious catalogue modelling concerns.

---

## 6. Event Outbox

The `catalogue_outbox_events` table supports reliable event publishing.

Catalogue events should be serialised as Protobuf payloads.

Recommended content type:

```text
application/x-protobuf
```

Potential events:

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
CategoryUpdated
ProductAttributeDefinitionUpdated
```

---

## 7. Client-Facing Engineering Evidence

This database foundation demonstrates:

```text
service-owned data design
least-privilege database thinking
repeatable migrations
realistic seed data
catalogue modelling for varied product types
event-driven outbox readiness
```

This is the bridge from architecture documentation into implementation.
