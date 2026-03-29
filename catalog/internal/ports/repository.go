package ports

import (
   "context"
   "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)

// RepositoryPort Interface describes expected method signatures for the catalog service repository implementation.
type RepositoryPort interface {
     GetProducts(ctx context.Context, searchType domain.SearchType) ([]domain.Product, error)
}


