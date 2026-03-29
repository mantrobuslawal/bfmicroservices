package ports

import (
	"context"
	"github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)

// APIPort interface of the catalog service describes expected endpoints available to
// search catalog repository for products 
type APIPort interface {
    GetProducts(ctx context.Context, searchOpt  domain.SearchType) ([]domain.Product, error)   
}
