package domain


// Avoiding modelling attributes in initial stage
// TODO update attribute types to domain types 

// Item being sold in store
type Product struct {
     SKU 	  string /*sku*/		   `json: "sku"`
     Name 	  string /*productName*/           `json: "name"`
     Brand 	  string /*manufacturer*/	   `json: "brand"`
     UnitPrice 	  float64 /*price*/		   `json: "unit_price"`
     Sizes 	  []string /*size*/		   `json: "sizes"`
     Description  string /*productDescription*/    `json: "description"`
     Category 	  string /*category*/              `json: "category"`
     Subcategory  string /*subCategory*/	   `json: "sub_category"` 
}


// Type used to query catalog database
// on different fields i.e. brand, sku, product name etc.
type SearchType struct {
     Opt         SearchOption
     SearchValue []string
}

type SearchOption int

// Create enum representing search categories
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
type rating float32
type review string // TODO - CREATE STRUCT{}
*/
