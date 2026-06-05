package catalog

import "time"

// Money represents a monetary value in minor units.
//
// This mirrors the Protobuf Money contract.
type Money struct {
	AmountMinor  int64
	CurrencyCode string
}

// CatalogID models the product,Attribute, Sku, category and partentCategory Id fields.
type CatalogID string

// Product is the internal domain representation of a catalogue product.
type Product struct {
	ProductID   CatalogID
	CategoryID  CatalogID
	Name        string
	Slug        string
	Description string
	Brand       string
	Status      ProductStatus
	BasePrice   Money
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Variants    []*ProductVariant
	Attributes  []*ProductAttributeValue
        Images      []*ProductImage
}

// Category is the internal domain representation of a product category.
type Category struct {
	CategoryID        CatalogID
	ParentCategoryID  CatalogID
	Name              string
	Slug              string
	Description       string
	Status            CategoryStatus
	DisplayOrder     int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}


// ProductVariant represents a purchaseable product variation.
type ProductVariant struct{
	VariantID CatalogID
        ProductID CatalogID
        Sku       CatalogID
        VariantName string
        Status  ProductVariantStatus
        Price   Money
	Attributes []*ProductAttributeValue
	CreatedAt time.Time
	UpdatedAt time.Time 
}

// ProductAttributeDefintion defines an attribute that is
// valid for a category
type ProductAttributeDefinition struct {
	AttributeID CatalogID
	CategoryID CatalogID
	Code string
	DisplayName string
	Description string
	DataType ProductAttributeDataType
	Unit string
	IsRequired bool
	IsFilterable bool
	IsVariantDefining bool
	Options []*ProductAttributeOption
	Status ProductAttributeDefinitionStatus
}

// ProductAttributeOption represents a controlled allowed
// value for an attribute definition.
type ProductAttributeOption struct {
	OptionID CatalogID
	Value	string
	DisplayName string
	DisplayOrder int
	Status 	ProductAttributeOptionStatus
}

// ProductAttributeValue represents a product-specific or
// a variant-specific value for a defined attribute.
type ProductAttributeValue struct {
	AttributeID  CatalogID
	Code string
	DisplayName string
	DataType ProductAttributeDataType
	Value *ProductAttributeValue
	ValueOptions []string
	Unit string
}


// ProductImage represents customer-facing catalogue imagery.
type ProductImage struct {
	ImageID CatalogID
	ProductID string
	Url string
	AltText string
	DisplayOrder int
}


// ListProductsFilter defines filter for product listing.
type ListProductsFilter struct {
	CategoryID      CatalogID
	IncludeInactive bool
	PageSize        int
	PageToken       string
}

// ListCategoriesFilter defines filter for category listing.
type ListCategoriesFilter struct {
	ParentCategoryID CatalogID
	IncludeInactive  bool
	PageSize         int
	PageToken        string
}

// ListProductAttributeDefinitionFilter defines filter for 
// product attribute definitions.
type ListProductAttributeDefinitionFilter struct {
	// Required category ID.
	CategoryID CatalogID
	IsFilterable bool
	IncludeInactive bool
        PageSize int
	PageToken string	
}
