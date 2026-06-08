// Package catalog contains the domain model and application service logic for
// the bfstore Catalog Service.
//
// The package owns catalog concepts such as products, categories, product
// variants, product attributes, product images, lifecycle statuses, and catalog-specific
// validation rules.
//
// The Catalog Service is the source of truth for the product identity and product
// description data. It does not own stock levels, basket contents, orders, payments,
// shipping, search indexes, or recommendations. Those concerns belong to other services
// or downstream projections.
//
// Code in this package should stay independent of transport and persistence
// details. Generated Protobuf types belong at the gRPC boundary, and SQL/database
// concerns belong in repository implementations. The service layer should work
// with catalog domain types rather than generated API or database-specific
// representations.
//
// Typical responsibilities in this package include:
//
//   - validating catalog identifiers, lifecycle statuses, prices and display
//     ordering;
//   - enforcing catalog business rules before data is persisted;
//   - coordinating repository calls for product, category, variant, attribute,
//     and image data;
//   - returning domain errors that transport adapters can map to API-specific
//     error responses.
//
// Repository implementations translate between MySQL rows the domain models
// defined here. Transport adapters, such as gRPC handlers, translate between
// generated Protobuf messages and service inputs/results.
package catalog
