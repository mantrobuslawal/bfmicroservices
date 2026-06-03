package grpc

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	common "github.com/mantrobuslawal/bfstore/gen/go/bfstore/common/v1"
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
	result, err := h.catalogService.ListProducts(ctx, catalog.ListProductsFilter{
    	CategoryID:       req.GetCategoryId(),
    	IncludeInactive: req.GetIncludeInactive(),
 		PageSize:           int(req.GetPage().GetPageSize()),
 		PageToken:          req.GetPage().GetPageToken(), 	})
 
	if err != nil {
 		h.logger.Error("failed to list products", "error", err)
 		return nil, mapServiceError(err)
 	}

 	response := &catalogv1.ListProductsResponse{
 		Products: make([]*catalogv1.Product, 0, len(result.Products)),
		Page: &common.PageResponse{
			NextPageToken: result.NextToken,
			TotalCount: 0, // Not calculated. Set to 0 as default
		}
 	}

 	for _, product := range result.Products {
 		response.Products = append(response.Products, productToProto(product))
 	}

 	return response, nil
 }

// GetProduct returns single product matching requested productID.
func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
	productID := req.GetProductId()
	if strings.Trim(productID) == "" {
		return nil, status.Error(code.INVALID_ARGUMENT, "missing product id")  

	}
	product, err := h.catalogService.GetProduct(ctx, productID)
 	if err != nil {
 		h.logger.Error("failed to get product", "product_id", productID, "error", err)
 		return nil, mapServiceError(err)
 	}

 	return &catalogv1.GetProductResponse{
 		Product: productToProto(product),
 	}, nil
 }

// ListCategories returns list of product categories matching filter criteria.
func (h *CatalogHandler) ListCategories(ctx context.Context, req *catalogv1.ListCategoriesRequest) (*catalogv1.ListCategoriesResponse, error) {
	result, err := h.catalogService.ListCategories(ctx, catalog.ListCategoriesFilter{
 		ParentCategoryID: req.GetParentCategoryId(),
 		IncludeInactive: req.GetIncludeInactive(),
 		PageSize:           int(req.GetPage().GetPageSize()),
 		PageToken:          req.GetPage().GetPageToken(),
 	})

 	if err != nil {
 		h.logger.Error("failed to list categories", "error", err)
 		return nil, mapServiceError(err)
 	}

 	response := &catalogv1.ListCategoriesResponse{
 		Categories: make([]*catalogv1.Category, 0, len(result.Categories)),
		Page: &common.PageResponse{
			NextPageToken: result.NextPageToken,
			TotalCount: 0, // Not calculated - set to 0. 
		}
 	}

 	for _, category := range result.Categories {
 		response.Categories = append(response.Categories, categoryToProto(category))
 	}

 	return response, nil
 }

// ListProductAttributeDefinitions returns list of product attribute definitions
// matching filter criteria.
func (h *CatalogHandler)  ListProductAttributeDefinitions(ctx context.Context, 
	  req *catalogv1.ListProductAttributeDefinitionsRequest) 
	(*catalogv1.ListProductAttributeDefinitionsResponse, error) {
	
	// category ID required
	catID := req.GetCategoryId()
	if strings.trim(catId) == "" {
		return nil, status.Error(code.INVALID_ARGUMENT, "missing category id")  
	}
	
	result, err := h.catalogService.ListProductAttributeDefinitons(ctx, 
				catalog.ListProductAttributeDefinitionsFilter{
					CategoryID: catId,
                                        FilterableOnly: req.GetFilterableOnly(),
                                        IncludeInactive: req.GetIncludeInactive(),
					PageSize: int(req.GetPage().GetPageSize()),	
					PageToken: req.GetPage().GetPageToken(),
				})
	if err != nil {
		h.logger.Error("failed to list product attribute definitions", "error", err)
		return nil, mapServiceError(err)
	}

	response := &catalogv1.ListAttributeDefinitionsResponse {
		AttributeDefinitions: make([]*catalogv1.ProductAttributeDefinition,
					   0, len(result.ProductAttributeDefinitions),
		Page: &common.PageResponse{
			NextPageToken: result.NextPageToken,
			TotalCount: 0, // Not calculated - set to 0. 
		}),
	}

	for _, productAttributeDefinition := range result.ProductAttributeDefinitions {
		response.AttributeDefinitions = append(response.AttributeDefinitions, 
			productAttributeDefinitionToProto(productAttributeDefinition))
	}

	return response, nil
}



