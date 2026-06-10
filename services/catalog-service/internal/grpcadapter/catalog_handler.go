package grpcadapter

import (
	"context"
	"log/slog"
	"strings"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	commonv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/common/v1"
)

// CatalogHandler implements the generated Catalog Service gRPC interface.
type CatalogHandler struct {
	catalogv1.UnimplementedCatalogServiceServer

	catalogService *catalog.Service
	logger         *slog.Logger
}

// NewCatalogHandler creates a Catalog Service gRPC handler.
func NewCatalogHandler(catalogService *catalog.Service, logger *slog.Logger) *CatalogHandler {
	return &CatalogHandler{
		catalogService: catalogService,
		logger:         logger,
	}
}

// ListProducts returns collection of products matching filter criteria.
// Products are not fully hydrated - images, variant data, product attributes etc
// - are not provided by this endpoint.
func (h *CatalogHandler) ListProducts(ctx context.Context, req *catalogv1.ListProductsRequest) (*catalogv1.ListProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil ListProductsRequest")

	}

	result, err := h.catalogService.ListProducts(ctx, catalog.ListProductsFilter{
		CategoryID:      catalog.CategoryID(req.GetCategoryId()),
		IncludeInactive: req.GetIncludeInactive(),
		PageSize:        int(req.GetPage().GetPageSize()),
		PageToken:       req.GetPage().GetPageToken()})

	if err != nil {
		h.logger.Error("failed to list products", "error", err)
		return nil, mapServiceError(err)
	}

	products := result.Result

	response := &catalogv1.ListProductsResponse{
		Products: make([]*catalogv1.Product, 0, len(products)),
		Page: &commonv1.PageResponse{
			NextPageToken: result.NextPageToken,
			TotalCount:    0, // Not calculated. Set to 0 as default
		},
	}

	for _, product := range products {
		protoProduct, err := listProductToProto(&product)
		if err != nil {
			return nil, mapServiceError(err)
		}
		response.Products = append(response.Products, protoProduct)
	}

	return response, nil
}

// GetProduct returns single product matching requested productID.
// This Product instance is fully hydrated - including image, variant, attributes etc.
func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil GetProductRequest")

	}

	productID := req.GetProductId()
	if strings.TrimSpace(productID) == "" {
		return nil, status.Error(codes.InvalidArgument, "missing product id")

	}

	id := catalog.ProductID(productID)

	product, err := h.catalogService.GetProduct(ctx, id)
	if err != nil {
		h.logger.Error("failed to get product", "product_id", productID, "error", err)
		return nil, mapServiceError(err)
	}

	// TODO: Need to update GetProduct request proto to include attribute definition filter options !!
	response, err := h.catalogService.ListProductAttributeDefinitions(ctx, catalog.ListProductAttributeDefinitionsFilter{CategoryID: product.CategoryID})
	if err != nil {
		h.logger.Error("failed to get product attribute definitions", "product_id", productID, "error", err)
		return nil, mapServiceError(err)
	}
	definitions := make([]catalog.ProductAttributeDefinition, 0, len(response.Result))
	definitions = append(definitions, response.Result...)

	// This endpoint uses pagination, therefore we need to ensure we have all the definitions required
	// to hydrate the product structure.
	for response.NextPageToken != "" {
		next := response.NextPageToken
		response, err := h.catalogService.ListProductAttributeDefinitions(ctx, catalog.ListProductAttributeDefinitionsFilter{
			CategoryID: product.CategoryID,
			PageToken:  next})
		if err != nil {
			h.logger.Error("failed to get product attribute definitions", "product_id", productID, "error", err)
			return nil, mapServiceError(err)
		}
		definitions = append(definitions, response.Result...)
	}

	hydratedProduct, err := productToProto(&product, definitions)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &catalogv1.GetProductResponse{
		Product: hydratedProduct,
	}, nil
}

// ListCategories returns list of product categories matching filter criteria.
func (h *CatalogHandler) ListCategories(ctx context.Context, req *catalogv1.ListCategoriesRequest) (*catalogv1.ListCategoriesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil ListCategoriesRequest")

	}

	result, err := h.catalogService.ListCategories(ctx, catalog.ListCategoriesFilter{
		ParentCategoryID: catalog.CategoryID(req.GetParentCategoryId()),
		IncludeInactive:  req.GetIncludeInactive(),
		PageSize:         int(req.GetPage().GetPageSize()),
		PageToken:        req.GetPage().GetPageToken(),
	})

	if err != nil {
		h.logger.Error("failed to list categories", "error", err)
		return nil, mapServiceError(err)
	}

	categories := result.Result
	token := result.NextPageToken

	response := &catalogv1.ListCategoriesResponse{
		Categories: make([]*catalogv1.Category, 0, len(categories)),
		Page: &commonv1.PageResponse{
			NextPageToken: token,
			TotalCount:    0, // Not calculated - set to 0.
		},
	}

	for _, category := range categories {
		categoryProto, err := categoryToProto(&category)
		if err != nil {
			return nil, mapServiceError(err)
		}
		response.Categories = append(response.Categories, categoryProto)
	}

	return response, nil
}

// ListProductAttributeDefinitions returns list of product attribute definitions
// matching filter criteria.
func (h *CatalogHandler) ListProductAttributeDefinitions(ctx context.Context,
	req *catalogv1.ListProductAttributeDefinitionsRequest) (*catalogv1.ListProductAttributeDefinitionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil ListProductAttributeDefinitions")

	}

	// category ID required
	catID := req.GetCategoryId()
	if strings.TrimSpace(catID) == "" {
		return nil, status.Error(codes.InvalidArgument, "missing category id")
	}

	result, err := h.catalogService.ListProductAttributeDefinitions(ctx,
		catalog.ListProductAttributeDefinitionsFilter{
			CategoryID:      catalog.CategoryID(catID),
			IsFilterable:    req.GetFilterableOnly(),
			IncludeInactive: req.GetIncludeInactive(),
			PageSize:        int(req.GetPage().GetPageSize()),
			PageToken:       req.GetPage().GetPageToken(),
		})
	if err != nil {
		h.logger.Error("failed to list product attribute definitions", "error", err)
		return nil, mapServiceError(err)
	}

	productAttributeDefinitions := result.Result

	response := &catalogv1.ListProductAttributeDefinitionsResponse{
		AttributeDefinitions: make([]*catalogv1.ProductAttributeDefinition, 0, len(productAttributeDefinitions)),
		Page: &commonv1.PageResponse{
			NextPageToken: result.NextPageToken,
			TotalCount:    0, // Not calculated - set to 0.
		},
	}

	for _, pad := range productAttributeDefinitions {
		padProto, err := productAttributeDefinitionToProto(&pad)
		if err != nil {
			return nil, mapServiceError(err)
		}

		response.AttributeDefinitions = append(response.AttributeDefinitions, padProto)
	}

	return response, nil
}
