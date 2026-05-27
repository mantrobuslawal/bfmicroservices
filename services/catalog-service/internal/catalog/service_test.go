package catalog

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	products   []Product
	categories []Category
	product    Product
	err        error
}

func (f *fakeRepository) ListProducts(ctx context.Context, filter ListProductsFilter) ([]Product, error) {
	return f.products, f.err
}

func (f *fakeRepository) GetProduct(ctx context.Context, productID string) (Product, error) {
	if f.err != nil {
		return Product{}, f.err
	}

	return f.product, nil
}

func (f *fakeRepository) ListCategories(ctx context.Context, filter ListCategoriesFilter) ([]Category, error) {
	return f.categories, f.err
}

func TestServiceListProducts(t *testing.T) {
	t.Parallel()

	repository := &fakeRepository{
		products: []Product{
			{
				ProductID:  "product-1",
				Name:       "Gopher Desk Lamp",
				Status:     "active",
				BasePrice:  Money{AmountMinor: 4599, CurrencyCode: "GBP"},
				CategoryID: "lighting",
			},
		},
	}

	service := NewService(repository)

	products, err := service.ListProducts(context.Background(), ListProductsFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	if products[0].Name != "Gopher Desk Lamp" {
		t.Fatalf("expected Gopher Desk Lamp, got %s", products[0].Name)
	}
}

func TestServiceGetProductRequiresProductID(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})

	_, err := service.GetProduct(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty product ID")
	}
}

func TestServiceGetProductNotFound(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{err: ErrNotFound})

	_, err := service.GetProduct(context.Background(), "missing-product")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceListCategories(t *testing.T) {
	t.Parallel()

	repository := &fakeRepository{
		categories: []Category{
			{
				CategoryID: "category-1",
				Name:       "Developer Homeware",
				Status:     "active",
			},
		},
	}

	service := NewService(repository)

	categories, err := service.ListCategories(context.Background(), ListCategoriesFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(categories))
	}
}
