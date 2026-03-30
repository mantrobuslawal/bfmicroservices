package repository

import "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"

// Data used to test slice repository implementation
var sliceCatalog = []domain.Product{
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

