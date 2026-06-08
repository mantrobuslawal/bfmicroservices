package catalog

// ProductAttributeDataType is the domain representation of the a product attribute definition's data type.
type ProductAttributeDataType string

// ProductAttributeDataType defined constants.
const (
	ProductAttributeTypeString      ProductAttributeDataType = "string"
	ProductAttributeTypeNumber      ProductAttributeDataType = "number"
	ProductAttributeTypeBool        ProductAttributeDataType = "boolean"
	ProductAttributeTypeOption      ProductAttributeDataType = "option"
	ProductAttributeTypeMultiOption ProductAttributeDataType = "multi_option"
	ProductAttributeTypeJson        ProductAttributeDataType = "json"
)

// IsValid validates whether a ProductAttributeDataType is belongs to the
// list of predefined constants. It returns true for valid data types and
// false otherwise.
func (p ProductAttributeDataType) IsValid() bool {
	switch p {
	case ProductAttributeTypeString,
		ProductAttributeTypeNumber,
		ProductAttributeTypeBool,
		ProductAttributeTypeOption,
		ProductAttributeTypeMultiOption,
		ProductAttributeTypeJson:
		return true
	default:
		return false
	}
}

// String converts a ProductAttributeDataType to a string.
func (p ProductAttributeDataType) String() string {
	return string(p)
}

// ProductAttributeValue is the domain representation
// of the a product attribute definition's data type.
type ProductAttributeValueKind string

// ProductAttributeValueType defined constants.
const (
	ProductAttributeValueKindUnspecified ProductAttributeValueKind = "unspecified"
	ProductAttributeValueKindNumber      ProductAttributeValueKind = "number"
	ProductAttributeValueKindString      ProductAttributeValueKind = "string"
	ProductAttributeValueKindBool        ProductAttributeValueKind = "boolean"
	ProductAttributeValueKindOption      ProductAttributeValueKind = "option"
	ProductAttributeValueKindMultiOption ProductAttributeValueKind = "multi_option"
	ProductAttributeValueKindJson        ProductAttributeValueKind = "json"
)

// IsValid validates whether a ProductAttributeValueType is belongs to the
// list of predefined constants. It returns true for valid data types and
// false otherwise.
func (v ProductAttributeValueKind) IsValid() bool {
	switch v {
	case ProductAttributeValueKindString,
		ProductAttributeValueKindNumber,
		ProductAttributeValueKindBool,
		ProductAttributeValueKindOption,
		ProductAttributeValueKindMultiOption,
		ProductAttributeValueKindJson:
		return true
	default:
		return false
	}
}

// String converts a ProductAttributeValueType to a string.
func (v ProductAttributeValueKind) String() string {
	return string(v)
}
