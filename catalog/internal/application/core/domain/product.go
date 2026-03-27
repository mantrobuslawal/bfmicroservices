package domain

// TODO imports

// Avoiding modelling attributes in initial stage
// TODO update attribute types to domain types 
type Product struct {
     SKU 	  string /*sku*/		        `json: "sku"`
     Name 	  string /*productName*/           `json: "name"`
     Brand 	  string /*manufacturer*/	        `json: "brand"`
     UnitPrice 	  float64 /*price*/		        `json: "unit_price"`
     Sizes 	  []string /*size*/		        `json: "sizes"`
     Description  string /*productDescription*/    `json: "description"`
     Category 	  string /*category*/              `json: "category"`
     Subcategory  string /*subCategory*/	        `json: "sub_category"` 
}


// Type used to query catalog database
type SearchType struct {
     opt         SearchOption
     searchValue []string
}

type SearchOption int

const (
    SKU SearchOption = iota
    ProductName
    Brand
    Category 
)

var optionName = map[SearchOption]string{
    SKU: "sku",
    ProductName: "name",
    Brand: "brand",
    Category: "category",
}

func (opt SearchOption) String() string {
	return optionName[opt]
}

/*
type productName string
type manufacturer string
type price float64 // TODO - create custom type
type size []string // TODO - create custom type
type productDescription string
type sku string
type category string // TODO - CREATE ENUM
type subCategory string // TODO - CREATE ENUM
type rating float32
type review string // TODO - CREATE STRUCT{}
*/
