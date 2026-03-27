package db

import (
  "context"
  "fmt"

  "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)

// Usual size of searchValue struct
const STD_QUERY_SIZE = 1

type tempDB []domain.Product

type Adapter struct {
     db *tempDB
}

func (a Adapter) Get(_ context.Context, search domain.SearchType) ([]domain.Product, error) {	
	results := nil
	field := query.Opt.String()
	searchValue1 := query.SearchValue[0]
	var searchValue2 string
	query_size := len(query.SearchValue)
	if len(query_size > STD_QUERY_SIZE) {
		searchValue2 = query.SearchValue[1]
	}
       
// Repitition exists in algorithm due to ability to search on any struct property and
// because I'm using a slice as temp db - will change implementation to real db 
	switch field {
		case "sku":
                for product := range *a.db {
			if product.SKU == searchValue {
				if results == nil { make([]domain.Product) }
				results = append(results, product)
			}
		}
		
		case "name":
                for product := range *a.db {
			if product.Name == searchValue {
				if results == nil { make([]domain.Product) }
				results = append(results, product)
			}
		}
		case "brand":
                for product := range *a.db {
			if product.Brand == searchValue {
				if results == nil { make([]domain.Product) }
				results = append(results, product)
			}
		}
		case "category":
                for product := range *a.db {
			if product.Category == searchValue {
				if product.Subcategory != "" && product.Subcategory == searchValue2 || query_size == STD_QUERY_SIZE  {
					if results == nil { make([]domain.Product) }
				        results = append(results, product)
				}
			}
		}
		default:
			err := fmt.Errorf("Unknown product property: %s", field)
			return []domain.Product(results), err
	}


	return results, nil
}

func NewAdapater(db ports.DBPort) (*Adapter, error) {
	return &Adapter{db: db}, nil
}


