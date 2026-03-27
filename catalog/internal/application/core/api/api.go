package api

import (
   "context"
   "github.com/mantrobuslawal/bfmircoservices/catalog/internal/application/core/domain"
   "github.com/mantrobuslawal/bfmicroservices/catalog/internal/ports" 
)

type Application struct {
     db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
     return &Application{
	db: db,
     }
}

func (a Application) Get(ctx context.Context, searchOpt domain.SearchType) ([]domain.Product, error) {
	products, err := db.Get(ctx, searchOpt)
	if err != nil {
		return []domain.Product(nil), err
	}

	return products, nil
}

