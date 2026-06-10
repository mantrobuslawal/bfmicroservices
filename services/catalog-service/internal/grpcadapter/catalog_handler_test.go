package grpcadapter

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	commonv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/common/v1"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeCatalogRepository struct {
	listProductsFunc                    func(context.Context, catalog.ListQuery) ([]catalog.Product, error)
	getProductFunc                      func(context.Context, catalog.ProductID) (catalog.Product, error)
	listCategoriesFunc                  func(context.Context, catalog.ListQuery) ([]catalog.Category, error)
	listProductAttributeDefinitionsFunc func(context.Context, catalog.ListQuery) ([]catalog.ProductAttributeDefinition, error)
}

func (f fakeCatalogRepository) ListProducts(ctx context.Context, query catalog.ListQuery) ([]catalog.Product, error) {
	if f.listProductsFunc == nil {
		return nil, errors.New("listProductsFunc not configured")
	}

	return f.listProductsFunc(ctx, query)
}

func (f fakeCatalogRepository) GetProduct(ctx context.Context, productID catalog.ProductID) (catalog.Product, error) {
	if f.getProductFunc == nil {
		return catalog.Product{}, errors.New("getProductFunc not configured")
	}

	return f.getProductFunc(ctx, productID)
}

func (f fakeCatalogRepository) ListCategories(ctx context.Context, query catalog.ListQuery) ([]catalog.Category, error) {
	if f.listCategoriesFunc == nil {
		return nil, errors.New("listCategoriesFunc not configured")
	}

	return f.listCategoriesFunc(ctx, query)
}

func (f fakeCatalogRepository) ListProductAttributeDefinitions(
	ctx context.Context,
	query catalog.ListQuery,
) ([]catalog.ProductAttributeDefinition, error) {
	if f.listProductAttributeDefinitionsFunc == nil {
		return nil, errors.New("listProductAttributeDefinitionsFunc not configured")
	}

	return f.listProductAttributeDefinitionsFunc(ctx, query)
}

func newTestCatalogHandler(repository catalog.Repository) *CatalogHandler {
	return NewCatalogHandler(catalog.NewService(repository), slog.Default())
}

func TestNewCatalogHandlerPanicsForNilCatalogService(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("NewCatalogHandler() did not panic for nil catalog service")
		}
	}()

	_ = NewCatalogHandler(nil, slog.Default())
}

func TestNewCatalogHandlerDefaultsNilLogger(t *testing.T) {
	t.Parallel()

	handler := NewCatalogHandler(catalog.NewService(fakeCatalogRepository{}), nil)

	if handler.logger == nil {
		t.Fatal("logger = nil, want default logger")
	}
}

func TestCatalogHandler_ListProductsRejectsNilRequest(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{})

	_, err := handler.ListProducts(context.Background(), nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestCatalogHandler_ListProductsReturnsProducts(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)

	handler := newTestCatalogHandler(fakeCatalogRepository{
		listProductsFunc: func(ctx context.Context, query catalog.ListQuery) ([]catalog.Product, error) {
			if query.ID != "cat_lighting" {
				t.Fatalf("query.ID = %q, want cat_lighting", query.ID)
			}

			if len(query.FilterOptions) != 1 || !query.FilterOptions[0] {
				t.Fatalf("query.FilterOptions = %+v, want includeInactive=true", query.FilterOptions)
			}

			if query.Limit != 3 {
				t.Fatalf("query.Limit = %d, want 3", query.Limit)
			}

			return []catalog.Product{
				{
					ProductID:   catalog.ProductID("prod_gopher_lamp"),
					CategoryID:  catalog.CategoryID("cat_lighting"),
					Name:        "Gopher Desk Lamp",
					Slug:        "gopher-desk-lamp",
					Description: "A cheerful lamp.",
					Brand:       "Borough",
					Status:      catalog.ProductStatusActive,
					BasePrice:   catalog.Money{AmountMinor: 4999, CurrencyCode: "GBP"},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			}, nil
		},
	})

	got, err := handler.ListProducts(context.Background(), &catalogv1.ListProductsRequest{
		Page: &commonv1.PageRequest{
			PageSize: 2,
		},
		CategoryId:      "cat_lighting",
		IncludeInactive: true,
	})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}

	if len(got.GetProducts()) != 1 {
		t.Fatalf("len(Products) = %d, want 1", len(got.GetProducts()))
	}

	if got.GetProducts()[0].GetProductId() != "prod_gopher_lamp" {
		t.Fatalf("ProductId = %q, want prod_gopher_lamp", got.GetProducts()[0].GetProductId())
	}

	if got.GetPage() == nil {
		t.Fatal("Page = nil, want page response")
	}
}

func TestCatalogHandler_ListProductsMapsServiceError(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{
		listProductsFunc: func(ctx context.Context, query catalog.ListQuery) ([]catalog.Product, error) {
			return nil, catalog.ErrInvalidPageSize
		},
	})

	_, err := handler.ListProducts(context.Background(), &catalogv1.ListProductsRequest{
		Page: &commonv1.PageRequest{
			PageSize: 2,
		},
	})

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestCatalogHandler_GetProductRejectsNilRequest(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{})

	_, err := handler.GetProduct(context.Background(), nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestCatalogHandler_GetProductRejectsMissingProductID(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{})

	_, err := handler.GetProduct(context.Background(), &catalogv1.GetProductRequest{
		ProductId: "   ",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestCatalogHandler_GetProductReturnsHydratedProduct(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)

	handler := newTestCatalogHandler(fakeCatalogRepository{
		getProductFunc: func(ctx context.Context, productID catalog.ProductID) (catalog.Product, error) {
			if productID != catalog.ProductID("prod_gopher_lamp") {
				t.Fatalf("productID = %q, want prod_gopher_lamp", productID)
			}

			return catalog.Product{
				ProductID:   catalog.ProductID("prod_gopher_lamp"),
				CategoryID:  catalog.CategoryID("cat_lighting"),
				Name:        "Gopher Desk Lamp",
				Slug:        "gopher-desk-lamp",
				Description: "A cheerful lamp.",
				Brand:       "Borough",
				Status:      catalog.ProductStatusActive,
				BasePrice:   catalog.Money{AmountMinor: 4999, CurrencyCode: "GBP"},
				CreatedAt:   now,
				UpdatedAt:   now,
				Attributes: []*catalog.ProductAttributeValue{
					{
						ProductAttributeValueID: catalog.ProductAttributeValueID("pav_material"),
						ProductID:              catalog.ProductID("prod_gopher_lamp"),
						AttributeID:            catalog.AttributeID("attr_material"),
						ValueString:            "steel",
						CreatedAt:              now,
						UpdatedAt:              now,
					},
				},
			}, nil
		},
		listProductAttributeDefinitionsFunc: func(ctx context.Context, query catalog.ListQuery) ([]catalog.ProductAttributeDefinition, error) {
			if query.ID != "cat_lighting" {
				t.Fatalf("query.ID = %q, want cat_lighting", query.ID)
			}

			return []catalog.ProductAttributeDefinition{
				{
					AttributeID: catalog.AttributeID("attr_material"),
					CategoryID:  catalog.CategoryID("cat_lighting"),
					Code:        "material",
					DisplayName: "Material",
					DataType:    catalog.ProductAttributeDataTypeString,
					Status:      catalog.ProductAttributeDefinitionStatusActive,
					CreatedAt:   now,
				},
			}, nil
		},
	})

	got, err := handler.GetProduct(context.Background(), &catalogv1.GetProductRequest{
		ProductId: "prod_gopher_lamp",
	})
	if err != nil {
		t.Fatalf("GetProduct() error = %v, want nil", err)
	}

	if got.GetProduct().GetProductId() != "prod_gopher_lamp" {
		t.Fatalf("ProductId = %q, want prod_gopher_lamp", got.GetProduct().GetProductId())
	}

	if len(got.GetProduct().GetAttributes()) != 1 {
		t.Fatalf("len(Attributes) = %d, want 1", len(got.GetProduct().GetAttributes()))
	}
}

func TestCatalogHandler_GetProductMapsNotFound(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{
		getProductFunc: func(ctx context.Context, productID catalog.ProductID) (catalog.Product, error) {
			return catalog.Product{}, catalog.ErrProductNotFound
		},
	})

	_, err := handler.GetProduct(context.Background(), &catalogv1.GetProductRequest{
		ProductId: "prod_missing",
	})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.NotFound)
	}
}

func TestCatalogHandler_ListCategoriesRejectsNilRequest(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{})

	_, err := handler.ListCategories(context.Background(), nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestCatalogHandler_ListCategoriesReturnsCategories(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	parentID := catalog.CategoryID("cat_lighting")

	handler := newTestCatalogHandler(fakeCatalogRepository{
		listCategoriesFunc: func(ctx context.Context, query catalog.ListQuery) ([]catalog.Category, error) {
			if query.ID != "cat_lighting" {
				t.Fatalf("query.ID = %q, want cat_lighting", query.ID)
			}

			if query.Limit != 3 {
				t.Fatalf("query.Limit = %d, want 3", query.Limit)
			}

			return []catalog.Category{
				{
					CategoryID:       catalog.CategoryID("cat_desk_lamps"),
					ParentCategoryID: &parentID,
					Name:             "Desk Lamps",
					Slug:             "desk-lamps",
					Description:      "Desk lamps.",
					Status:           catalog.CategoryStatusActive,
					DisplayOrder:     10,
					CreatedAt:        now,
					UpdatedAt:        now,
				},
			}, nil
		},
	})

	got, err := handler.ListCategories(context.Background(), &catalogv1.ListCategoriesRequest{
		Page: &commonv1.PageRequest{
			PageSize: 2,
		},
		ParentCategoryId: "cat_lighting",
	})
	if err != nil {
		t.Fatalf("ListCategories() error = %v, want nil", err)
	}

	if len(got.GetCategories()) != 1 {
		t.Fatalf("len(Categories) = %d, want 1", len(got.GetCategories()))
	}

	if got.GetCategories()[0].GetParentCategoryId() != "cat_lighting" {
		t.Fatalf("ParentCategoryId = %q, want cat_lighting", got.GetCategories()[0].GetParentCategoryId())
	}
}

func TestCatalogHandler_ListProductAttributeDefinitionsRejectsNilRequest(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{})

	_, err := handler.ListProductAttributeDefinitions(context.Background(), nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestCatalogHandler_ListProductAttributeDefinitionsRejectsMissingCategoryID(t *testing.T) {
	t.Parallel()

	handler := newTestCatalogHandler(fakeCatalogRepository{})

	_, err := handler.ListProductAttributeDefinitions(
		context.Background(),
		&catalogv1.ListProductAttributeDefinitionsRequest{
			CategoryId: "   ",
		},
	)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status.Code() = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestCatalogHandler_ListProductAttributeDefinitionsReturnsDefinitions(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)

	handler := newTestCatalogHandler(fakeCatalogRepository{
		listProductAttributeDefinitionsFunc: func(ctx context.Context, query catalog.ListQuery) ([]catalog.ProductAttributeDefinition, error) {
			if query.ID != "cat_lighting" {
				t.Fatalf("query.ID = %q, want cat_lighting", query.ID)
			}

			if len(query.FilterOptions) != 2 {
				t.Fatalf("len(FilterOptions) = %d, want 2", len(query.FilterOptions))
			}

			if !query.FilterOptions[0] {
				t.Fatal("includeInactive = false, want true")
			}

			if !query.FilterOptions[1] {
				t.Fatal("isFilterable = false, want true")
			}

			return []catalog.ProductAttributeDefinition{
				{
					AttributeID:  catalog.AttributeID("attr_material"),
					CategoryID:   catalog.CategoryID("cat_lighting"),
					Code:         "material",
					DisplayName:  "Material",
					Description:  "Primary material.",
					DataType:     catalog.ProductAttributeDataTypeString,
					IsRequired:   true,
					IsFilterable: true,
					Status:       catalog.ProductAttributeDefinitionStatusActive,
					CreatedAt:    now,
				},
			}, nil
		},
	})

	got, err := handler.ListProductAttributeDefinitions(
		context.Background(),
		&catalogv1.ListProductAttributeDefinitionsRequest{
			Page: &commonv1.PageRequest{
				PageSize: 2,
			},
			CategoryId:      "cat_lighting",
			FilterableOnly:  true,
			IncludeInactive: true,
		},
	)
	if err != nil {
		t.Fatalf("ListProductAttributeDefinitions() error = %v, want nil", err)
	}

	if len(got.GetAttributeDefinitions()) != 1 {
		t.Fatalf("len(AttributeDefinitions) = %d, want 1", len(got.GetAttributeDefinitions()))
	}

	if got.GetAttributeDefinitions()[0].GetAttributeId() != "attr_material" {
		t.Fatalf("AttributeId = %q, want attr_material", got.GetAttributeDefinitions()[0].GetAttributeId())
	}
}
