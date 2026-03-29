package repository

import (
  "context"
  "fmt"

  "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)

// Usual size of searchValue struct
const STD_QUERY_SIZE = 1

// slice implementation of ports.RepositoryPort interface for catalog service
type sliceRepo []domain.Product

type Adapter struct {
     repo sliceRepo
}

// GetProducts searches catalog repo, filtering records based on property chosen by client
func (a Adapter) GetProducts(ctx context.Context, query domain.SearchType) ([]domain.Product, error) {	
	var results []domain.Product	// Return empty slice if no match found
	field := query.Opt.String()
	searchValue1 := query.SearchValue[0]
	var searchValue2 string
	if len(query.SearchValue) > STD_QUERY_SIZE {
		searchValue2 = query.SearchValue[1]
	} else {
		searchValue2 = ""
	}
	
// Repitition exists in algorithm due to ability to search on any struct property and
// because I'm using a slice as temp db - will change implementation to real db 
	switch field {
		case "sku":
                for _, product := range a.repo {
			if product.SKU == searchValue1 {
				results = append(results, product)
			}
		}
		
		case "name":
                for _, product := range a.repo {
			if product.Name == searchValue1 {
				results = append(results, product)
			}
		}
		case "brand":
                for _, product := range a.repo {
			if product.Brand == searchValue1 {
				results = append(results, product)
			}
		}
		case "category":
                for _, product := range a.repo {
			if product.Category == searchValue1 {
		// The second case represents when there is no subcategory to filter.
		// It is always true, tested in toplevel if product.Category == searchValue
				if (product.Subcategory != "" && product.Subcategory == searchValue2) || true  { 
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

func NewAdapter(repo []domain.Product) (Adapter, error) {
	return Adapter{repo: repo}, nil
}


