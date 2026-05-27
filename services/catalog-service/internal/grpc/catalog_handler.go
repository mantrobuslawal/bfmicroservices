package grpc

import (
	"log/slog"

	"github.com/acme-ltd/bfstore/services/catalog-service/internal/catalog"
)

// CatalogHandler will implement the generated Catalogue Service gRPC interface.
//
// This file is intentionally kept as a placeholder until generated Protobuf
// code is committed. Once generated code exists, this type should implement:
//
//   catalogv1.CatalogServiceServer
//
// Initial RPCs:
//
//   ListProducts
//   GetProduct
//   ListCategories
type CatalogHandler struct {
	catalogService *catalog.Service
	logger         *slog.Logger
}

// NewCatalogHandler creates a Catalogue Service gRPC handler.
func NewCatalogHandler(catalogService *catalog.Service, logger *slog.Logger) *CatalogHandler {
	return &CatalogHandler{
		catalogService: catalogService,
		logger:         logger,
	}
}

// TODO:
// Add methods once generated Protobuf contracts are available.
//
// Expected methods:
//
//   func (h *CatalogHandler) ListProducts(ctx context.Context, req *catalogv1.ListProductsRequest) (*catalogv1.ListProductsResponse, error)
//   func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error)
//   func (h *CatalogHandler) ListCategories(ctx context.Context, req *catalogv1.ListCategoriesRequest) (*catalogv1.ListCategoriesResponse, error)
