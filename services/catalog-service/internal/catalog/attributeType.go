package catalog

// ProductAttributeType is the domain representation of the a product attribute definition's data type.
type ProductAttributeDataType string

// ProductAttributeDataType defined constants.
const (
	ProductAttributeString      ProductAttributeDataType = "string"
	ProductAttributeNumber      ProductAttributeDataType = "number"
	ProductAttributeBool        ProductAttributeDataType = "boolean"
	ProductAttributeOption      ProductAttributeDataType = "option"
	ProductAttributeMultiOption ProductAttributeDataType = "multi_option"
	ProductAttributeJson        ProductAttributeDataType = "json"
)

// IsValid validates whether a ProductAttributeDataType is belongs to the
// list of predefined constants. It returns true for valid data types and
// false otherwise.
func (p ProductAttributeType) IsValid() bool {
	switch p {
	case ProductAttributeString,
		ProductAttributeNumber,
		ProductAttributeBool,
		ProductAttributeOption,
		ProductAttributeMultiOption,
		ProductAttributeJson:
		return true
	default:
		return false
	}
}

// String converts a ProductAttributeDataType to a string.
func (p ProductAttributeType) String() string {
	return string(p)
}

// ProductAttributeValue is the domain representation
// of the a product attribute definition's data type.
type ProductAttributeValue string

// ProductAttributeValue defined constants.
const (
	ProductAttributeValueUnspecified ProductAttributeValue = "unspecified"
	ProductAttributeValueString      ProductAttributeValue = "string"
	ProductAttributeValueNumber      ProductAttributeValue = "number"
	ProductAttributeValueBool        ProductAttributeValue = "boolean"
	ProductAttributeValueOption      ProductAttributeValue = "option"
	ProductAttributeValueMultiOption ProductAttributeValue = "multi_option"
	ProductAttributeValueJson        ProductAttributeValue = "json"
)

// IsValid validates whether a ProductAttributeValue is belongs to the
// list of predefined constants. It returns true for valid data types and
// false otherwise.
func (v ProductAttributeValue) IsValid() bool {
	switch v {
	case ProductAttributeValueString,
		ProductAttributeValueNumber,
		ProductAttributeValueBool,
		ProductAttributeValueOption,
		ProductAttributeValueMultiOption,
		ProductAttributeValueJson:
		return true
	default:
		return false
	}
}

// String converts a ProductAttributeValue to a string.
func (v ProductAttributeValue) String() string {
	return string(v)
}
