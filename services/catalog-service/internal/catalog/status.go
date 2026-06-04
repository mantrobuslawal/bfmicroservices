package catalog

type LifecycleStatus interface {
	ProductStatus |
	CategoryStatus |
	ProductVariantStatus |
	ProductAttributeDefinitionStatus |
	ProductAttributeOptionStatus
}

func isKnownLifecycleStatus[S LifecycleStatus](status S) bool {
	switch string(status) {
	case "draft", "active", "inactive", "archived":
		return true
	default:
		return false
	}
} 


// ProductStatus is the domain representation of a
// product's lifecycle.
type ProductStatus string

const(
	 ProductStatusDraft    ProductStatus = "draft" 
	 ProductStatusActive   ProductStatus = "active"
	 ProductStatusInactive ProductStatus = "inactive"
	 ProductStatusArchived ProductStatus = "archived"
)

func (s ProductStatus) IsValid() bool {
	return isKnownLifecycleStatus(s)
}

func (s ProductStatus) String() string {
	return string(s)
}


// CategoryStatus is the domain representation of the
// product category lifecycle.
type CategoryStatus string

const(
	 CategoryStatusDraft     CategoryStatus = "draft" 
	 CategoryStatusActive    CategoryStatus = "active"
	 CategortyStatusInactive CategoryStatus = "inactive"
	 CategoryStatusArchived  CategoryStatus = "archived"
)

func (c CategoryStatus) IsValid() bool {
	return isKnownLifecycleStatus(c)
}


func (c CategoryStatus) String() string {
	return string(c)
}


// ProductVariantStatus is the domain representation of a
// products variant's lifecycle.
type ProductVariantStatus string

const(
	 ProductVariantStatusDraft    ProductVariantStatus = "draft" 
	 ProductVariantStatusActive   ProductVariantStatus = "active"
	 ProductVariantStatusInactive ProductVariantStatus = "inactive"
	 ProductVariantStatusArchived ProductVariantStatus = "archived"
)

func (v ProductVariantStatus) IsValid() bool {
	return isKnownLifecycleStatus(v)
}


func (v ProductVariantStatus) String() string {
	return string(v)
}


// ProductAttributeDefinitionStatus is the domain representation of a
// product attribute's lifecycle.
type ProductAttributeDefinitionStatus string

const(
	 ProductAttributeDefinitionStatusDraft  ProductAttributeDefinitionStatus = "draft" 
	 ProductAttributeDefinitionStatusActive  ProductAttributeDefinitionStatus = "active"
	 ProductAttributeDefinitionStatusInactive ProductAttributeDefinitionStatus = "inactive"
	 ProductAttributeDefinitionStatusArchived ProductAttributeDefinitionStatus = "archived"
)

func ( s ProductAttributeDefinitionStatus) IsValid() bool {
}


func (s  ProductAtrributeDefinitionStatus) String() string {
	return string(s)
}


// ProductAttributeOptionStatus is the domain representation of a
// product attribute option's lifecycle.
type ProductAttributeOptionStatus string

const(
	 ProductAttributeOptionsDraft    ProductAttributeOptionStatus   = "draft" 
	 ProductAttributeOptionsActive   ProductAttributeOptionStatus   = "active"
	 ProductAttributeOptionsInactive  ProductAttributeOptionStatus  = "inactive"
	 ProductAttributeOptionsArchived  ProductAttributeOptionStatus  = "archived"
)

func (o ProductAttributeOptionStatus ) IsValid() bool {
	return isKnownLifecycleStatus(o)
}


func (o  ProductAttributeOptionStatus ) String() string {
	return string(o)
}

