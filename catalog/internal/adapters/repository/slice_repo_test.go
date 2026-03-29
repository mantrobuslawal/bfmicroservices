package repository

import (
   "testing"
   "context"
    
    "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"

    "github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	sliceCatalog := []domain.Product{
		{
			SKU: "abdcdegh12345",
		        Name: "gopher desk",
                        Brand: "the golang furniture company",
                        UnitPrice: 59.99,
			Sizes: []string{"standard"},
                        Description: "gopher desk 3000",
 		        Category: "office furniture",
			Subcategory: "desks",
		},

		{
			SKU: "abdcdzzz78945",
		        Name: "a gopher's day out hanging",
                        Brand: "rob pike tapestry",
                        UnitPrice: 96.99,
			Sizes: []string{"standard"},
                        Description: "a gopher's day out tapestry wall hanging",
 		        Category: "home decor",
			Subcategory: "wall decor",
		},
		
		{
			SKU: "abdcdooo33654",
		        Name: "rust desk",
                        Brand: "mr. karbs office furnishings",
                        UnitPrice: 79.99,
			Sizes: []string{"standard"},
                        Description: "crustacean home office desk",
 		        Category: "office furniture",
			Subcategory: "desks",
		},

	}

	tests := map[string]struct{
	   query       domain.SearchType
           products    []domain.Product
           expectedErr error
	}{
		"sku in repo": {
			query: domain.SearchType{domain.SKU, []string{"abdcdegh12345"}},
			products: sliceCatalog[0:1],
			expectedErr: nil,
                       
		},
		
		"sku not repo": {
			query: domain.SearchType{domain.SKU, []string{"xxxxxxxxx"}},
			products: nil,
			expectedErr: nil,
                       
		},
		
		"sku empty string": {
			query: domain.SearchType{domain.SKU, []string{""}},
			products: nil,
			expectedErr: nil,
                       
		},
               
		"name in repo": {
			query: domain.SearchType{domain.ProductName, []string{"gopher desk"}},
			products: sliceCatalog[0:1],
			expectedErr: nil,
                       
		},
		
		"name not in repo": {
			query: domain.SearchType{domain.ProductName, []string{"spam and eggs"}},
			products: nil,
			expectedErr: nil,
                       
		},
		
		"name as empty string": {
			query: domain.SearchType{domain.ProductName, []string{""}},
			products: nil,
			expectedErr: nil,
                       
		},  
		
		"brand in repo": {
			query: domain.SearchType{domain.Brand, []string{"rob pike tapestry"}},
			products: sliceCatalog[1:2],
			expectedErr: nil,
                       
		},
		
		"brand not in repo": {
			query: domain.SearchType{domain.Brand, []string{"foo"}},
			products: nil,
			expectedErr: nil,
                       
		},
		
		"brand as empty string": {
			query: domain.SearchType{domain.Brand, []string{""}},
			products: nil,
			expectedErr: nil,
                       
		},
		  
		"category in repo": {
			query: domain.SearchType{domain.Category, []string{"office furniture"}},
			products: []domain.Product{sliceCatalog[0], sliceCatalog[2]},
			expectedErr: nil,
                       
		},
		
		"category and subcategory in repo": {
			query: domain.SearchType{domain.Category, []string{"home decor", "wall decor"}},
			products: sliceCatalog[1:2],
			expectedErr: nil,
                       
		},
		
		"category in repo, but subcategory not in repo": {
			query: domain.SearchType{domain.Category, []string{"home decor", "foo"}},
			products: nil,
			expectedErr: nil,
                       
		}, 
		 
		"category not in repo, but subcategory in repo": {
			query: domain.SearchType{domain.Category, []string{"fuzz", "wall decor"}},
			products: nil,
			expectedErr: nil,
                       
		},
		
		"category empty string and subcategory non-empty string": {
			query: domain.SearchType{domain.Category, []string{"", "wall decor"}},
			products: nil,
			expectedErr: nil,
                       
		},
		
		"category and subcategory not in repo": {
			query: domain.SearchType{domain.Category, []string{"bar", "foo"}},
			products: nil,
			expectedErr: nil,
                       
		},  
	}

	repo, _ := NewAdapter(sliceCatalog)

	for name, tc := range tests {
		name, tc := name, tc
		
		t.Run(name, func(t *testing.T){
			got, err := repo.GetProducts(context.Background(), tc.query)
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				assert.Equal(t, tc.products, got)	
			}
		})
	}
}
