package catalog

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
	switch s {
	case  ProductStatusDraft,
	      ProductStatusActive,
	      ProductStatusInactive,
	      ProductStatusArchived:
	      return true
	default:
	      return false 
	}
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
	switch c {
	case  CategoryStatusDraft,
	      CategoryStatusActive,
	      CategoryStatusInactive,
	      CategoryStatusArchived:
	      return true
	default:
	      return false 
	}
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
	switch v {
	case  ProductVariantStatusDraft,
	      ProductVariantStatusActive,
	      ProductVariantStatusInactive,
	      ProductVariantStatusArchived:
	      return true
	default:
	      return false 
	}
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
	switch s {
	case  ProducAttributeDefinitionStatusDraft,
	      ProducAttributeDefinitionStatusActive,
	      ProducAttributeDefinitionStatusInactive,
	      ProductAttributeDefinitionStatusArchived:
	      return true
	default:
	      return false 
	}
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
	switch o {
	case  ProductAttributeOptionStatusDraft,
	      ProductAttributeOptionStatusActive,
	      ProductAttributeOptionStatusInactive,
	      ProductAttributeOptionStatusArchived:
	      return true
	default:
	      return false 
	}
}


func (o  ProductAttributeOptionStatus ) String() string {
	return string(o)

}

