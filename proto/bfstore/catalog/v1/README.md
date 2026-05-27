# Catalogue Protobuf Package

## Package

```proto
package acme.catalog.v1;
```

## Purpose

This package defines the public gRPC API and catalogue domain messages for the bfstore Catalogue Service.

Catalogue Service owns:

```text
products
categories
product variants
category-scoped product attributes
product images
catalogue lifecycle state
```

It does not own:

```text
stock levels
basket state
orders
payments
shipping
search ranking
recommendations
```

---

## Field Presence Convention

bfstore uses proto3.

Use `optional` for scalar fields when the service must distinguish between:

```text
field omitted
field explicitly set to the scalar default value
field explicitly set to a non-default value
```

This is especially important for Catalogue admin/update APIs.

Examples:

```text
description omitted        means do not change description
description set to ""      means intentionally clear description
active omitted             means do not change status
active set to false        means intentionally deactivate
base_price_minor omitted   means do not change price
base_price_minor set to 0  means intentionally set to zero, if business rules allow it
```

Do not use `optional` blindly for normal response models.

---

## Recommended Usage

### Plain scalar fields

Use plain scalar fields for response messages and required-by-validation identifiers.

Example:

```proto
message Product {
  string product_id = 1;
  string category_id = 2;
  string name = 3;
  string slug = 4;
  string description = 5;
  string brand = 6;
  ProductStatus status = 7;
  acme.common.v1.Money base_price = 8;
}
```

Although proto3 cannot enforce required fields, the service should always populate these values in responses.

### Optional scalar fields

Use `optional` for patch/update request fields.

Example:

```proto
message UpdateProductRequest {
  string product_id = 1;

  optional string name = 2;
  optional string description = 3;
  optional string brand = 4;
  optional int64 base_price_minor = 5;
  optional string currency_code = 6;
  optional ProductStatus status = 7;
}
```

The `product_id` is not optional because the command must identify the target product.

The update fields are optional because the service must know which fields the caller intends to modify.

---

## List and Search Filters

For list/search APIs, use `optional` when absence has a different meaning from the default value.

Example:

```proto
message ListProductsRequest {
  optional string category_id = 1;
  optional string search_query = 2;
  optional bool include_inactive = 3;
  acme.common.v1.PageRequest page = 4;
}
```

Interpretation:

```text
category_id omitted   = do not filter by category
category_id set to "" = caller explicitly supplied an empty category filter; validate or reject
include_inactive omitted = use service default
include_inactive false   = caller explicitly requests active-only results
```

If the API treats omitted and default values the same, plain scalar fields are also acceptable.

Be deliberate and document the behaviour.

---

## Product Attributes

Catalogue products support category-scoped flexible attributes.

Examples:

```text
lamps: bulb_type, max_wattage
wall art: wall_art_size
lockboxes: security_rating, lock_type
rugs: rug_shape, pile_height
soft furnishings: fabric_type, care_instructions
```

For attribute values, field presence may matter because an attribute can intentionally contain a default-like value.

Recommended pattern:

```proto
message ProductAttributeValue {
  string attribute_id = 1;

  optional string value_string = 2;
  optional double value_number = 3;
  optional bool value_boolean = 4;
  optional string value_json = 5;
  optional string unit = 6;
}
```

The service should validate that only the correct value field is set for the attribute definition's data type.

---

## Validation

The Catalogue Service must validate required-by-convention fields.

Examples:

```text
product_id must be present for update/get commands
name must not be empty when creating a product
slug must be valid and unique
currency_code must be valid
base_price_minor must satisfy price rules
status transitions must be allowed
attribute values must match their definitions
```

Do not rely on Protobuf field presence alone for business validation.

---

## Implementation Note for Go

When a scalar field is declared `optional`, generated Go code allows the handler to check whether the field was present.

The handler should map presence into domain update commands rather than passing Protobuf messages throughout the service.

Recommended boundary:

```text
generated request
→ gRPC handler
→ domain command with explicit optional values
→ service layer
→ repository
```

Avoid leaking generated Protobuf types into repository code.
