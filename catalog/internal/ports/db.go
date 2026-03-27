package ports

import (
   "context"
   "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)

type DBPort interface {
     Get(ctx context.Context, searchType domain.SearchType) ([]domain.Product, error)
}


