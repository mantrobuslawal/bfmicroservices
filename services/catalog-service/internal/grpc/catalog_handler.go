package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
)

// CatalogHandler implements the generated Catalogue Service gRPC interface.
type CatalogHandler struct {
	catalogv1.UnimplementedCatalogServiceServer

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

// ListProducts returns collection of products matching filter criteria.
func (h *CatalogHandler) ListProducts(ctx context.Context, req *catalogv1.ListProductsRequest) (*catalogv1.ListProductsResponse, error) {
	products, err := h.catalogService.ListProducts(ctx, catalog.ListProductsFilter{
    	CategoryID:       req.GetCategoryId(),
    	IncludeInactive: req.GetIncludeInactive(),
 		Limit:           int(req.GetPage().GetPageSize()),
 		Offset:          0, 	})
 	if err != nil {
 		h.logger.Error("failed to list products", "error", err)
 		return nil, status.Error(codes.Internal, "failed to list products")
 	}

 	response := &catalogv1.ListProductsResponse{
 		Products: make([]*catalogv1.Product, 0, len(products)),
 	}

 	for _, product := range products {
 		response.Products = append(response.Products, productToProto(product))
 	}

 	return response, nil
 }

// GetProduct returns single product matching requested productID.
func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
 	product, err := h.catalogService.GetProduct(ctx, req.GetProductId())
 	if errors.Is(err, catalog.ErrNotFound) {
 		return nil, status.Error(codes.NotFound, "product not found")
 	}
 	if err != nil {
 		h.logger.Error("failed to get product", "product_id", req.GetProductId(), "error", err)
 		return nil, status.Error(codes.Internal, "failed to get product")
 	}

 	return &catalogv1.GetProductResponse{
 		Product: productToProto(product),
 	}, nil
 }

// ListCategories returns list of product categories matching filter criteria.
func (h *CatalogHandler) ListCategories(ctx context.Context, req *catalogv1.ListCategoriesRequest) (*catalogv1.ListCategoriesResponse, error) {
 	categories, err := h.catalogService.ListCategories(ctx, catalog.ListCategoriesFilter{
 		ParentCategoryID: req.GetParentCategoryId(),
 		IncludeInactive: req.GetIncludeInactive(),
 		Limit:           int(req.GetPage().GetPageSize()),
 		Offset:          0,
 	})
 	if err != nil {
 		h.logger.Error("failed to list categories", "error", err)
 		return nil, status.Error(codes.Internal, "failed to list categories")
 	}

 	response := &catalogv1.ListCategoriesResponse{
 		Categories: make([]*catalogv1.Category, 0, len(categories)),
 	}

 	for _, category := range categories {
 		response.Categories = append(response.Categories, categoryToProto(category))
 	}

 	return response, nil
 }


// grpcError maps domain errors to gRPC status errors.
func grpcError(logger *slog.Logger, publicMessage string, err error) error {
	if errors.Is(err, catalog.ErrNotFound) {
		return status.Error(codes.NotFound, publicMessage)
	}

	logger.Error(publicMessage, "error", err)
	return status.Error(codes.Internal, publicMessage)
}


