package api

import (
   "context"
   "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
   "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/ports" 
)

// Struct to hold repository
type Application struct {
     repo ports.RepositoryPort
}

func NewApplication(repo ports.RepositoryPort) *Application {
     return &Application{
	repo: repo,
     }
}

// GetProducts retrieves products from repository filtered by SearchOpt
func (a Application) GetProducts(ctx context.Context, searchOpt domain.SearchType) ([]domain.Product, error) {
	products, err := a.repo.GetProducts(ctx, searchOpt)
	if err != nil {
		return []domain.Product(nil), err
	}

	return products, nil
}

