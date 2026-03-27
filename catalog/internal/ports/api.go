package ports

import (
	"context"
	"github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)

// Interface for catalog api
type APIPort interface {
    Get(ctx context.Context, searchOpt  domain.SearchType) ([]domain.Product, error)   
}
