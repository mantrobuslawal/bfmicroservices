package repository

import (
  "context"
  "fmt"

  "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
  "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/ports"
)

// Usual size of searchValue struct
const STD_QUERY_SIZE = 1

// slice implementation of ports.RepositoryPort interface for catalog service
type sliceRepo []domain.Product

type Adapter struct {
     repo *sliceRepo
}

// GetProducts searches catalog repo, filtering records based on property chosen by client
func (a Adapter) GetProducts(ctx context.Context, query domain.SearchType) ([]domain.Product, error) {	
	results := nil	// Return empty slice if no match found
	field := query.Opt.String()
	searchValue1 := query.SearchValue[0]
	
// Repitition exists in algorithm due to ability to search on any struct property and
// because I'm using a slice as temp db - will change implementation to real db 
	switch field {
		case "sku":
                for product := range *a.repo {
			if product.SKU == searchValue {
				if results == nil { make([]domain.Product) }
				results = append(results, product)
			}
		}
		
		case "name":
                for product := range *a.repo {
			if product.Name == searchValue {
				if results == nil { make([]domain.Product) }
				results = append(results, product)
			}
		}
		case "brand":
                for product := range *a.repo {
			if product.Brand == searchValue {
				if results == nil { make([]domain.Product) }
				results = append(results, product)
			}
		}
		case "category":
                for product := range *a.repo {
			if product.Category == searchValue {
		// The second case represents when there is no subcategory to filter.
		// It is always true, tested in toplevel if product.Category == searchValue
				if (product.Subcategory != "" && product.Subcategory == searchValue2) || true  { 
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

func NewAdapater(repo ports.RepositoryPort) (*Adapter, error) {
	return &Adapter{repo: repo}, nil
}


