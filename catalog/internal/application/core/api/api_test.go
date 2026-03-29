package api

import (
   "context"
   "testing"

   "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"

   "github.com/stretchr/testify/assert"  
)

type testRepo struct { store []domain.Product }

func (tr testRepo) GetProducts(context.Context, domain.SearchType) ([]domain.Product, error) {
	return tr.store, nil
}

func TestGetProducts(t *testing.T) {
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

		opt := domain.SearchType{}
		got, err := app.GetProducts(context.Background(), opt)
		assert.ErrorIs(err, nil)
        	assert.Equal(got, mock.store)
	})
}
