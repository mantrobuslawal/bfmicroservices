package catalog

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	listProductsFunc                     func(context.Context, ListQuery) ([]Product, error)
	getProductFunc                       func(context.Context, string) (Product, error)
	listCategoriesFunc                   func(context.Context, ListQuery) ([]Category, error)
	listProductAttributeDefinitionsFunc  func(context.Context, ListQuery) ([]ProductAttributeDefinition, error)
}

func (f fakeRepository) ListProducts(ctx context.Context, query ListQuery) ([]Product, error) {
	if f.listProductsFunc == nil {
		return nil, errors.New("listProductsFunc not configured")
	}

	return f.listProductsFunc(ctx, query)
}

func (f fakeRepository) GetProduct(ctx context.Context, productID string) (Product, error) {
	if f.getProductFunc == nil {
		return Product{}, errors.New("getProductFunc not configured")
	}

	return f.getProductFunc(ctx, productID)
}

func (f fakeRepository) ListCategories(ctx context.Context, query ListQuery) ([]Category, error) {
	if f.listCategoriesFunc == nil {
		return nil, errors.New("listCategoriesFunc not configured")
	}

	return f.listCategoriesFunc(ctx, query)
}

func (f fakeRepository) ListProductAttributeDefinitions(
	ctx context.Context,
	query ListQuery,
) ([]ProductAttributeDefinition, error) {
	if f.listProductAttributeDefinitionsFunc == nil {
		return nil, errors.New("listProductAttributeDefinitionsFunc not configured")
	}

	return f.listProductAttributeDefinitionsFunc(ctx, query)
}

func TestService_GetProduct(t *testing.T) {
	t.Parallel()

	want := Product{
		ProductID: ProductID("prod_test_001"),
		Name:      "Gopher Desk Lamp",
		Status:    ProductStatusActive,
	}

	service := NewService(fakeRepository{
		getProductFunc: func(ctx context.Context, productID string) (Product, error) {
			if productID != string(want.ProductID) {
				t.Fatalf("productID = %q, want %q", productID, want.ProductID)
			}

			return want, nil
		},
	})

	got, err := service.GetProduct(context.Background(), want.ProductID)
	if err != nil {
		t.Fatalf("GetProduct() error = %v, want nil", err)
	}

	if got.ProductID != want.ProductID {
		t.Fatalf("ProductID = %q, want %q", got.ProductID, want.ProductID)
	}

	if got.Name != want.Name {
		t.Fatalf("Name = %q, want %q", got.Name, want.Name)
	}
}

func TestService_GetProductMapsRepositoryNotFound(t *testing.T) {
	t.Parallel()

	service := NewService(fakeRepository{
		getProductFunc: func(ctx context.Context, productID string) (Product, error) {
			return Product{}, ErrProductNotFound
		},
	})

	_, err := service.GetProduct(context.Background(), ProductID("prod_missing"))

	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("GetProduct() error = %v, want ErrProductNotFound", err)
	}
}

func TestService_ListProductsUsesDefaultPageSize(t *testing.T) {
	t.Parallel()

	service := NewService(fakeRepository{
		listProductsFunc: func(ctx context.Context, query ListQuery) ([]Product, error) {
			if query.Limit != defaultPageSize+1 {
				t.Fatalf("Limit = %d, want %d", query.Limit, defaultPageSize+1)
			}

			return []Product{
				{
					ProductID: ProductID("prod_test_001"),
					Name:      "Gopher Desk Lamp",
					CreatedAt: time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC),
				},
			}, nil
		},
	})

	result, err := service.ListProducts(context.Background(), ListProductsFilter{})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}

	if len(result.Result) != 1 {
		t.Fatalf("len(Result) = %d, want 1", len(result.Result))
	}

	if result.NextPageToken != "" {
		t.Fatalf("NextPageToken = %q, want empty", result.NextPageToken)
	}
}

func TestService_ListProductsRejectsOversizedPage(t *testing.T) {
	t.Parallel()

	service := NewService(fakeRepository{})

	_, err := service.ListProducts(context.Background(), ListProductsFilter{
		PageSize: maxPageSize + 1,
	})

	if !errors.Is(err, ErrInvalidPageSize) {
		t.Fatalf("ListProducts() error = %v, want ErrInvalidPageSize", err)
	}
}

func TestService_ListProductsReturnsNextPageTokenWhenMoreResultsExist(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)

	service := NewService(fakeRepository{
		listProductsFunc: func(ctx context.Context, query ListQuery) ([]Product, error) {
			if query.Limit != 3 {
				t.Fatalf("Limit = %d, want 3", query.Limit)
			}

			return []Product{
				{
					ProductID: ProductID("prod_001"),
					Name:      "Gopher Desk Lamp",
					CreatedAt: createdAt.Add(2 * time.Minute),
				},
				{
					ProductID: ProductID("prod_002"),
					Name:      "Rob Pike Wall Tapestry",
					CreatedAt: createdAt.Add(1 * time.Minute),
				},
				{
					ProductID: ProductID("prod_003"),
					Name:      "Dijkstra Pathfinding Rug",
					CreatedAt: createdAt,
				},
			}, nil
		},
	})

	result, err := service.ListProducts(context.Background(), ListProductsFilter{
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}

	if len(result.Result) != 2 {
		t.Fatalf("len(Result) = %d, want 2", len(result.Result))
	}

	if result.NextPageToken == "" {
		t.Fatal("NextPageToken is empty, want token")
	}

	cursor, err := decodeCursor(result.NextPageToken)
	if err != nil {
		t.Fatalf("decode next page token: %v", err)
	}

	if cursor.ID != "prod_002" {
		t.Fatalf("cursor.ID = %q, want prod_002", cursor.ID)
	}
}

func TestService_ListCategoriesRejectsOversizedPage(t *testing.T) {
	t.Parallel()

	service := NewService(fakeRepository{})

	_, err := service.ListCategories(context.Background(), ListCategoriesFilter{
		PageSize: maxPageSize + 1,
	})

	if !errors.Is(err, ErrInvalidPageSize) {
		t.Fatalf("ListCategories() error = %v, want ErrInvalidPageSize", err)
	}
}

func TestService_ListProductAttributeDefinitionsRejectsOversizedPage(t *testing.T) {
	t.Parallel()

	service := NewService(fakeRepository{})

	_, err := service.ListProductAttributeDefinitions(context.Background(), ListProductAttributeDefinitionsFilter{
		PageSize: maxPageSize + 1,
	})

	if !errors.Is(err, ErrInvalidPageSize) {
		t.Fatalf("ListProductAttributeDefinitions() error = %v, want ErrInvalidPageSize", err)
	}
}
