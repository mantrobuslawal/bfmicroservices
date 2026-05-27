package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/acme-ltd/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// TODO:
	// Uncomment after generated Protobuf code is committed.
	//
	// catalogv1 "github.com/acme-ltd/bfstore/gen/go/acme/catalog/v1"
)

// CatalogHandler will implement the generated Catalogue Service gRPC interface.
//
// Once generated code exists, this type should embed:
//
//   catalogv1.UnimplementedCatalogServiceServer
//
// and implement:
//
//   ListProducts
//   GetProduct
//   ListCategories
type CatalogHandler struct {
	// catalogv1.UnimplementedCatalogServiceServer

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
// Replace the compile-safe placeholder methods below with generated Protobuf
// method signatures after running buf generate and committing generated code.

// Example ListProducts target implementation:
//
// func (h *CatalogHandler) ListProducts(ctx context.Context, req *catalogv1.ListProductsRequest) (*catalogv1.ListProductsResponse, error) {
// 	products, err := h.catalogService.ListProducts(ctx, catalog.ListProductsFilter{
// 		CategoryID:       req.GetCategoryId(),
// 		IncludeInactive: req.GetIncludeInactive(),
// 		Limit:           int(req.GetPage().GetPageSize()),
// 		Offset:          0,
// 	})
// 	if err != nil {
// 		h.logger.Error("failed to list products", "error", err)
// 		return nil, status.Error(codes.Internal, "failed to list products")
// 	}
//
// 	response := &catalogv1.ListProductsResponse{
// 		Products: make([]*catalogv1.Product, 0, len(products)),
// 	}
//
// 	for _, product := range products {
// 		response.Products = append(response.Products, productToProto(product))
// 	}
//
// 	return response, nil
// }
//
// Example GetProduct target implementation:
//
// func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
// 	product, err := h.catalogService.GetProduct(ctx, req.GetProductId())
// 	if errors.Is(err, catalog.ErrNotFound) {
// 		return nil, status.Error(codes.NotFound, "product not found")
// 	}
// 	if err != nil {
// 		h.logger.Error("failed to get product", "product_id", req.GetProductId(), "error", err)
// 		return nil, status.Error(codes.Internal, "failed to get product")
// 	}
//
// 	return &catalogv1.GetProductResponse{
// 		Product: productToProto(product),
// 	}, nil
// }
//
// Example ListCategories target implementation:
//
// func (h *CatalogHandler) ListCategories(ctx context.Context, req *catalogv1.ListCategoriesRequest) (*catalogv1.ListCategoriesResponse, error) {
// 	categories, err := h.catalogService.ListCategories(ctx, catalog.ListCategoriesFilter{
// 		ParentCategoryID: req.GetParentCategoryId(),
// 		IncludeInactive: req.GetIncludeInactive(),
// 		Limit:           int(req.GetPage().GetPageSize()),
// 		Offset:          0,
// 	})
// 	if err != nil {
// 		h.logger.Error("failed to list categories", "error", err)
// 		return nil, status.Error(codes.Internal, "failed to list categories")
// 	}
//
// 	response := &catalogv1.ListCategoriesResponse{
// 		Categories: make([]*catalogv1.Category, 0, len(categories)),
// 	}
//
// 	for _, category := range categories {
// 		response.Categories = append(response.Categories, categoryToProto(category))
// 	}
//
// 	return response, nil
// }

// grpcError maps domain errors to gRPC status errors.
//
// This helper can be used by the concrete generated-code handlers later.
func grpcError(logger *slog.Logger, publicMessage string, err error) error {
	if errors.Is(err, catalog.ErrNotFound) {
		return status.Error(codes.NotFound, publicMessage)
	}

	logger.Error(publicMessage, "error", err)
	return status.Error(codes.Internal, publicMessage)
}

// compileSafeHandlerCheck keeps imports meaningful until generated handlers are wired.
func compileSafeHandlerCheck(ctx context.Context, h *CatalogHandler) error {
	if h == nil || h.catalogService == nil {
		return status.Error(codes.FailedPrecondition, "catalog handler is not configured")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
