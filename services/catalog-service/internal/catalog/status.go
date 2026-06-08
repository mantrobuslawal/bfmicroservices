package catalog

import (
	"errors"
	"strings"
)

// LifecycleStatus interface represents life phases of catalog objects such as
// products, product categories and product variants.
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

// Valid ProductStatus values.
const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusArchived ProductStatus = "archived"
)

// IsValid returns true if the status is valid and false otherwise.
func (s ProductStatus) IsValid() bool {
	return isKnownLifecycleStatus(s)
}

// String converts a status from a LifecycleStatus to a string.
func (s ProductStatus) String() string {
	return string(s)
}

// CategoryStatus is the domain representation of the
// product category lifecycle.
type CategoryStatus string

// Valid CategoryStatus values.
const (
	CategoryStatusDraft     CategoryStatus = "draft"
	CategoryStatusActive    CategoryStatus = "active"
	CategortyStatusInactive CategoryStatus = "inactive"
	CategoryStatusArchived  CategoryStatus = "archived"
)

// IsValid returns true if the status is valid and false otherwise.
func (c CategoryStatus) IsValid() bool {
	return isKnownLifecycleStatus(c)
}

// String converts a status from a LifecycleStatus to a string.
func (c CategoryStatus) String() string {
	return string(c)
}

// ProductVariantStatus is the domain representation of a
// products variant's lifecycle.
type ProductVariantStatus string

// Valid ProductVariantStatus values.
const (
	ProductVariantStatusDraft    ProductVariantStatus = "draft"
	ProductVariantStatusActive   ProductVariantStatus = "active"
	ProductVariantStatusInactive ProductVariantStatus = "inactive"
	ProductVariantStatusArchived ProductVariantStatus = "archived"
)

// IsValid returns true if the status is valid and false otherwise.
func (v ProductVariantStatus) IsValid() bool {
	return isKnownLifecycleStatus(v)
}

// String converts a status from a LifecycleStatus to a string.
func (v ProductVariantStatus) String() string {
	return string(v)
}

// ProductAttributeDefinitionStatus is the domain representation of a
// product attribute's lifecycle.
type ProductAttributeDefinitionStatus string

// Valid ProductAttributeDefintionStatus values.
const (
	ProductAttributeDefinitionStatusDraft    ProductAttributeDefinitionStatus = "draft"
	ProductAttributeDefinitionStatusActive   ProductAttributeDefinitionStatus = "active"
	ProductAttributeDefinitionStatusInactive ProductAttributeDefinitionStatus = "inactive"
	ProductAttributeDefinitionStatusArchived ProductAttributeDefinitionStatus = "archived"
)

// IsValid returns true if the status is valid and false otherwise.
func (s ProductAttributeDefinitionStatus) IsValid() bool {
}

// String converts a status from a LifecycleStatus to a string.
func (s ProductAtrributeDefinitionStatus) String() string {
	return string(s)
}

// ProductAttributeOptionStatus is the domain representation of a
// product attribute option's lifecycle.
type ProductAttributeOptionStatus string

// Valid ProductAttrubuteOptions values.
const (
	ProductAttributeOptionsDraft    ProductAttributeOptionStatus = "draft"
	ProductAttributeOptionsActive   ProductAttributeOptionStatus = "active"
	ProductAttributeOptionsInactive ProductAttributeOptionStatus = "inactive"
	ProductAttributeOptionsArchived ProductAttributeOptionStatus = "archived"
)

// IsValid returns true if the status is valid and false otherwise.
func (o ProductAttributeOptionStatus) IsValid() bool {
	return isKnownLifecycleStatus(o)
}

// String converts a status from a LifecycleStatus to a string.
func (o ProductAttributeOptionStatus) String() string {
	return string(o)
}

// ParseLifecycleStatus converts a valid string value to a LifecycleStatus.
// Returns an ErrInvalidLifecycleStatus error when the input string cannot
// be matched to a valid LifecycleStatus.
func ParseLifecycleStatus(status string) (LifecycleStatus, error) {
	status = strings.Trim(status)
	switch status {
	case "draft":
		return ProductStatusDraft, nil
	case "active":
		return ProductStatusActive, nil
	case "inactive":
		return ProductStatusInactive, nil
	case "archived":
		return ProductStatusArchived, nil

	default:
		return nil, ErrInvalidLifecycleStatus
	}
}
