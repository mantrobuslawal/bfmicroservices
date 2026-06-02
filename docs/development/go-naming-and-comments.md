# Go Naming and Comments

This document defines bfstore naming and commenting conventions for Go code.

## Core Rule

```text
Names should be short, clear, package-aware, and idiomatic.
```

## Package Names

Use short, lower-case package names.

Good:

```text
catalog
checkout
outbox
payment
logging
telemetry
```

Avoid:

```text
catalogue_service_business_logic
catalogService
catalog_utils
```

## Exported Names

Names beginning with an upper-case letter are exported.

Export only what other packages should rely on.

```go
type Product struct {
    ID ProductID
}

func NewService(products ProductRepository) *Service {
    return &Service{products: products}
}
```

Kuti rule:

```text
Exporting is a commitment.
```

## Avoid Repetition

The package name is part of the public API name.

Prefer:

```go
catalog.Product
catalog.Service
catalog.NewService
```

Avoid:

```go
catalog.CatalogProduct
catalog.CatalogProductService
catalog.NewCatalogProductService
```

## Getters

Avoid unnecessary `Get` prefixes.

Prefer:

```go
func (o Order) Status() OrderStatus
```

Instead of:

```go
func (o Order) GetStatus() OrderStatus
```

## MixedCaps

Use `MixedCaps` or `mixedCaps`, not snake_case.

```go
type PaymentAttemptID string

func calculateOrderTotal() Money {
    // ...
}
```

## Comments

Comment exported packages, types, functions, methods, variables, and constants.

Good:

```go
// ProductID uniquely identifies a product in the catalogue.
type ProductID string
```

Bad:

```go
// ProductID is a ProductID.
type ProductID string
```

## Package Comments

```go
// Package catalog contains catalogue domain and application logic.
//
// It owns product retrieval, catalogue filtering, and product availability
// rules used by catalog-service.
package catalog
```

## Practical Rules

```text
Use short package names.
Use package names to reduce repetition.
Use MixedCaps.
Comment exported API.
Explain meaning, not obvious syntax.
Avoid unnecessary Get prefixes.
Do not export implementation details.
```

## Final Rule

```text
A good Go name should make surrounding code easier to read.
```
