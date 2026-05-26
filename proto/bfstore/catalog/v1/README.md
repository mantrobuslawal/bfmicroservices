# `proto/acme/catalog/v1`

## 1. Purpose

This package defines the Catalogue Service Protobuf contracts for bfstore.

Catalogue Service owns product catalogue truth, including product identity, product details, product status, category taxonomy, product variants, category-scoped product attribute definitions, product attribute values, and product pricing in the initial version.

Inventory Service owns stock. Search Service owns search projections. Recommendation Service owns recommendation outputs. Catalogue remains the source of truth for product data.

---

## 2. Package

```proto
package acme.catalog.v1;
```

Recommended Go package option:

```proto
option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/catalog/v1;catalogv1";
```

---

## 3. Ownership

| Responsibility | Owner |
|---|---|
| Product identity | `catalog-service` |
| Product name and description | `catalog-service` |
| Category taxonomy | `catalog-service` |
| Product variants | `catalog-service` |
| Product status | `catalog-service` |
| Product attribute definitions | `catalog-service` |
| Product attribute values | `catalog-service` |
| Product price in initial version | `catalog-service` |
| Stock quantity | `inventory-service` |
| Search ranking/projection | `search-service` |
| Recommendation outputs | `recommendation-service` |

---

## 4. Expected Files

```text
proto/acme/catalog/v1/
├── README.md
├── catalog_service.proto
├── product.proto
├── category.proto
├── product_attribute.proto
└── catalog_events.proto
```

Events may later move to:

```text
proto/acme/catalog/events/v1/
```

---

## 5. Primary Service

```proto
service CatalogService {
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc ListCategories(ListCategoriesRequest) returns (ListCategoriesResponse);
  rpc ListProductAttributeDefinitions(ListProductAttributeDefinitionsRequest) returns (ListProductAttributeDefinitionsResponse);
}
```

Possible later administrative APIs:

```proto
rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
rpc ActivateProduct(ActivateProductRequest) returns (ActivateProductResponse);
rpc DeactivateProduct(DeactivateProductRequest) returns (DeactivateProductResponse);
rpc CreateProductAttributeDefinition(CreateProductAttributeDefinitionRequest) returns (CreateProductAttributeDefinitionResponse);
rpc UpdateProductAttributeDefinition(UpdateProductAttributeDefinitionRequest) returns (UpdateProductAttributeDefinitionResponse);
```

Admin APIs may be deferred until the core customer journey is implemented.

---

## 6. Core Messages

Recommended messages:

```text
Product
ProductVariant
Category
ProductAttributeDefinition
ProductAttributeValue
ProductAttributeOption
ProductAttributeDataType
ProductImage
ProductStatus
```

---

## 7. Product

Conceptual example:

```proto
message Product {
  string product_id = 1;
  string name = 2;
  string description = 3;
  acme.common.v1.Money price = 4;
  ProductStatus status = 5;
  string category_id = 6;
  repeated ProductVariant variants = 7;
  repeated ProductAttributeValue attributes = 8;
}
```

The `Product` message should contain common product data and associated attribute values.

Highly type-specific fields should not be hardcoded into the base product message.

Avoid:

```proto
string curtain_drop = 50;
string bed_size = 51;
string bulb_type = 52;
string rug_shape = 53;
```

Use `ProductAttributeValue` instead.

---

## 8. Product Status

```proto
enum ProductStatus {
  PRODUCT_STATUS_UNSPECIFIED = 0;
  PRODUCT_STATUS_DRAFT = 1;
  PRODUCT_STATUS_ACTIVE = 2;
  PRODUCT_STATUS_INACTIVE = 3;
  PRODUCT_STATUS_ARCHIVED = 4;
}
```

Rules:

```text
only ACTIVE products are purchasable
inactive products must not be added to basket
archived products are retained for history but not normal purchase flows
```

---

## 9. Category

Conceptual example:

```proto
message Category {
  string category_id = 1;
  string parent_category_id = 2;
  string name = 3;
  string slug = 4;
  ProductStatus status = 5;
}
```

Categories define the product taxonomy and provide the scope for attribute definitions.

Examples:

```text
curtains
bed-frames
rugs
lamps
sofas
wardrobes
```

---

## 10. Product Attribute Definition

Product attribute definitions describe which attributes are valid for a category.

Conceptual example:

```proto
message ProductAttributeDefinition {
  string attribute_id = 1;
  string category_id = 2;
  string code = 3;
  string display_name = 4;
  string description = 5;
  ProductAttributeDataType data_type = 6;
  string unit = 7;
  bool is_required = 8;
  bool is_filterable = 9;
  bool is_variant_defining = 10;
  repeated ProductAttributeOption options = 11;
}
```

Example definitions:

| Category | Code | Data Type | Unit | Filterable |
|---|---|---|---|---|
| curtains | `drop_cm` | number | cm | yes |
| curtains | `heading_type` | string/option | none | yes |
| bed-frames | `bed_size` | option | none | yes |
| bed-frames | `storage_type` | option | none | yes |
| rugs | `shape` | option | none | yes |
| lamps | `bulb_type` | string/option | none | yes |

---

## 11. Product Attribute Data Type

Recommended enum:

```proto
enum ProductAttributeDataType {
  PRODUCT_ATTRIBUTE_DATA_TYPE_UNSPECIFIED = 0;
  PRODUCT_ATTRIBUTE_DATA_TYPE_STRING = 1;
  PRODUCT_ATTRIBUTE_DATA_TYPE_NUMBER = 2;
  PRODUCT_ATTRIBUTE_DATA_TYPE_BOOLEAN = 3;
  PRODUCT_ATTRIBUTE_DATA_TYPE_OPTION = 4;
  PRODUCT_ATTRIBUTE_DATA_TYPE_MULTI_OPTION = 5;
  PRODUCT_ATTRIBUTE_DATA_TYPE_JSON = 6;
}
```

Rules:

```text
data types must be stable
new data types should be introduced carefully
consumers must handle unsupported values safely
```

---

## 12. Product Attribute Value

Conceptual example:

```proto
message ProductAttributeValue {
  string attribute_id = 1;
  string code = 2;
  string display_name = 3;
  ProductAttributeDataType data_type = 4;
  string value_string = 5;
  double value_number = 6;
  bool value_boolean = 7;
  repeated string value_options = 8;
  string value_json = 9;
  string unit = 10;
}
```

Only the field matching the attribute data type should be populated.

Examples:

```text
drop_cm = 228
heading_type = eyelet
bed_size = king
storage_type = ottoman
rug_shape = round
bulb_type = E27
```

---

## 13. Required RPC Behaviour

## 13.1 `ListProducts`

Expected behaviour:

```text
returns active products by default
supports pagination
may support category filtering
must not return inactive products in normal browse flow
may include summary attributes suitable for display
```

Advanced filtering should eventually belong to Search Service.

## 13.2 `GetProduct`

Expected behaviour:

```text
returns product by product_id
returns common fields and relevant attributes
returns NOT_FOUND when product does not exist or is not visible
does not expose internal database details
```

## 13.3 `ListCategories`

Expected behaviour:

```text
returns customer-browsable categories
may later support hierarchy
```

## 13.4 `ListProductAttributeDefinitions`

Expected behaviour:

```text
returns attribute definitions for a category
supports UI/filter construction
helps clients understand which attributes apply to which product type
```

---

## 14. Event Contracts

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

Initial priority:

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
CategoryUpdated
ProductAttributeDefinitionUpdated
```

These events support Search Service and Recommendation Service projections.

---

## 15. Search Projection Relationship

Catalogue Service should provide enough event data for Search Service to build denormalised product search documents.

Catalogue event payloads should include or allow reconstruction of:

```text
product_id
category
status
price
display attributes
filterable attributes
variant summary
```

Search Service owns the denormalised search shape.

Catalogue Service owns the governed source data.

---

## 16. Error Behaviour

| Scenario | Code |
|---|---|
| Missing product ID | `INVALID_ARGUMENT` |
| Product not found | `NOT_FOUND` |
| Category not found | `NOT_FOUND` |
| Attribute definition not found | `NOT_FOUND` |
| Invalid page size | `INVALID_ARGUMENT` |
| Invalid attribute value for definition | `INVALID_ARGUMENT` |
| Required category attribute missing | `FAILED_PRECONDITION` |
| Catalogue database unavailable | `UNAVAILABLE` |
| Unexpected failure | `INTERNAL` |

---

## 17. Security and Privacy

Catalogue data is mostly public product data, but APIs should not expose:

```text
internal supplier notes
cost price
margin data
unpublished draft metadata
internal database IDs
```

Admin APIs should require authentication and authorisation once implemented.

---

## 18. Testing Expectations

Tests should cover:

```text
list active products
exclude inactive products
get product by ID
product not found
category attribute definitions returned
required attribute validation
attribute data type validation
filterable attributes present for search projection
product update event emitted where implemented
protobuf compatibility checks
```

---

## 19. Related Documents

```text
docs/requirements/functional-requirements.md
docs/requirements/business-rules.md
docs/api/protobuf-style-guide.md
docs/api/error-model.md
docs/events/event-catalog.md
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
```

---

## 20. Summary

`acme.catalog.v1` defines the product catalogue contract.

Catalogue Service owns governed product truth and supports varied product types through category-scoped product attribute definitions and values.

Search Service consumes catalogue events and owns denormalised product search projections.
