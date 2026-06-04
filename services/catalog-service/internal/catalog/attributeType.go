package catalog


// ProductAttributeType is the domain representation 
// of the a product attribute definition's data type.
type ProductAttributeDataType string

const(
	ProductAttributeString       ProductAttributeDataType = "string"
	ProductAttributeNumber      ProductAttributeDataType = "number"
	ProductAttributeBool         ProductAttributeDataType = "boolean"
	ProductAttributeOption       ProductAttributeDataType  = "option"
	ProductAttributeMultiOption   ProductAttributeDataType = "multi_option"
	ProductAttributeJson        ProductAttributeDataType = "json"
       
)

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

func (p ProductAttributeType) String() string {
	return string(p)
}


// ProductAttributeValue is the domain representation 
// of the a product attribute definition's data type.
type ProductAttributeValue string

const (
	ProductAttributeValueUnspecified ProductAttributeValue = "unspecified" 
	ProductAttributeValueString ProductAttributeValue = "string" 
	ProductAttributeValueNumber ProductAttributeValue = "number" 
	ProductAttributeValueBool ProductAttributeValue = "boolean" 
	ProductAttributeValueOption ProductAttributeValue = "option" 
	ProductAttributeValueMultiOption ProductAttributeValue = "multi_option" 
	ProductAttributeValueJson ProductAttributeValue = "json" 
)

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

func (v ProductAttributeValue) String() string {
	return string(v)
}


