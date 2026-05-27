package catalog

import "time"

// Money represents a monetary value in minor units.
//
// This mirrors the Protobuf Money contract.
type Money struct {
	AmountMinor  int64
	CurrencyCode string
}

// Product is the internal domain representation of a catalogue product.
type Product struct {
	ProductID   string
	CategoryID  string
	Name        string
	Slug        string
	Description string
	Brand       string
	Status      string
	BasePrice   Money
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Category is the internal domain representation of a product category.
type Category struct {
	CategoryID       string
	ParentCategoryID *string
	Name             string
	Slug             string
	Description      string
	Status           string
	DisplayOrder     int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ListProductsFilter defines filters for product listing.
type ListProductsFilter struct {
	CategoryID      string
	IncludeInactive bool
	Limit           int
	Offset          int
}

// ListCategoriesFilter defines filters for category listing.
type ListCategoriesFilter struct {
	ParentCategoryID string
	IncludeInactive  bool
	Limit            int
	Offset           int
}
