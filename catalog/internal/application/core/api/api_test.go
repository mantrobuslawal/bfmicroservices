package api

import (
   "context"
   "testing"

   "github.com/mantrobuslawal/bfmicroservices/catalog.git/domain"
   "github.com/mantrobuslawal/bfmicroservices/catalog.git/ports"

   "github.com/stretchr/testify/assert"  
)

func TestGetProducts(t *testing.T) {
	type testRepo struct {store []domain.Product}	

	func (tr testRepo) GetProduct(ctx context.Context, opt domain.SearchType) []domain.Product {
		return tr.store
	}

        var mock testRepo
        mock.store = []domain.Product{{
                         SKU: "abcfghd12345",
                         Name: "gopher desk",
		         Brand: "the golang furniture company",
                         UnitPrice: 59.99,
                         Sizes: []string{"standard"},
                         Description: "gopher desk 3000",
                         Category: "office furniture",
                         Subcategory: "desks",
                },
	}

        
	t.Run("nominal", func(t *testing.T) {
		app := NewApplication(mock)
        	assert := assert.New(t)

		got, err := app.GetProducts(context.Background(), _ domain.SearchType)
		assert.ErrorIs(err, nil)
        	assert.Equal(items, mock.store)
	}
}
