package catalog

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	listProductsFunc                     func(context.Context, ListQuery) ([]Product, error)
	getProductFunc                       func(context.Context, ProductID) (Product, error)
	listCategoriesFunc                   func(context.Context, ListQuery) ([]Category, error)
	listProductAttributeDefinitionsFunc  func(context.Context, ListQuery) ([]ProductAttributeDefinition, error)
}

func (f fakeRepository) ListProducts(ctx context.Context, query ListQuery) ([]Product, error) {
	if f.listProductsFunc == nil {
		return nil, errors.New("listProductsFunc not configured")
	}

	return f.listProductsFunc(ctx, query)
}

func (f fakeRepository) GetProduct(ctx context.Context, productID ProductID) (Product, error) {
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

func TestService_GetProductReturnsHydratedProductDetails(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 6, 8, 11, 0, 0, 0, time.UTC)

	productLevelValueID := ProductAttributeValueID("pav_material")
	variantValueID := ProductAttributeValueID("pav_colour")
	variantID := VariantID("var_blue")

	product := Product{
		ProductID:    ProductID("prod_gopher_lamp"),
		CategoryID:   CategoryID("cat_lighting"),
		Name:         "Gopher Desk Lamp",
		Slug:         "gopher-desk-lamp",
		Description:  "A cheerful lamp for late-night debugging.",
		Brand:        "Borough",
		Status:       ProductStatusActive,
		BasePrice:    Money{AmountMinor: 4999, CurrencyCode: "GBP"},
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Images: []*ProductImage{
			{
				ImageID:      ImageID("img_primary"),
				ProductID:    ProductID("prod_gopher_lamp"),
				Url:          "https://example.test/gopher-lamp.png",
				AltText:      "Gopher Desk Lamp on a desk",
				DisplayOrder: 10,
			},
		},
		Variants: []*ProductVariant{
			{
				VariantID:   variantID,
				ProductID:   ProductID("prod_gopher_lamp"),
				Sku:         Sku("BFS-GO-LAMP-BLUE"),
				VariantName: "Blue shade",
				Status:      ProductVariantStatusActive,
				Price:       Money{AmountMinor: 5499, CurrencyCode: "GBP"},
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			},
		},
		Attributes: []*ProductAttributeValue{
			{
				ProductAttributeValueID: productLevelValueID,
				ProductID:              ProductID("prod_gopher_lamp"),
				AttributeID:            AttributeID("attr_material"),
				ValueString:            "steel",
				CreatedAt:              createdAt,
				UpdatedAt:              updatedAt,
			},
			{
				ProductAttributeValueID: variantValueID,
				ProductID:              ProductID("prod_gopher_lamp"),
				VariantID:              &variantID,
				AttributeID:            AttributeID("attr_colour"),
				ValueString:            "blue",
				CreatedAt:              createdAt,
				UpdatedAt:              updatedAt,
			},
		},
	}

	definitions := []ProductAttributeDefinition{
		{
			AttributeID: AttributeID("attr_material"),
			CategoryID:  CategoryID("cat_lighting"),
			Code:        "material",
			DisplayName: "Material",
			DataType:    ProductAttributeDataTypeString,
			Status:      ProductAttributeDefinitionStatusActive,
			CreatedAt:   createdAt,
		},
		{
			AttributeID: AttributeID("attr_colour"),
			CategoryID:  CategoryID("cat_lighting"),
			Code:        "colour",
			DisplayName: "Colour",
			DataType:    ProductAttributeDataTypeString,
			Options: []*ProductAttributeOption{
				{
					OptionID:     OptionID("opt_blue"),
					Value:        "blue",
					DisplayName:  "Blue",
					DisplayOrder: 10,
					Status:       ProductAttributeOptionStatusActive,
				},
			},
			Status:    ProductAttributeDefinitionStatusActive,
			CreatedAt: createdAt,
		},
	}

	service := NewService(fakeRepository{
		getProductFunc: func(ctx context.Context, productID ProductID) (Product, error) {
			if productID != product.ProductID {
				t.Fatalf("productID = %q, want %q", productID, product.ProductID)
			}

			return product, nil
		},
		listProductAttributeDefinitionsFunc: func(ctx context.Context, query ListQuery) ([]ProductAttributeDefinition, error) {
			if query.ID != string(product.CategoryID) {
				t.Fatalf("query.ID = %q, want %q", query.ID, product.CategoryID)
			}

			if got, want := query.Limit, 500; got != want {
				t.Fatalf("query.Limit = %d, want %d", got, want)
			}

			if len(query.FilterOptions) != 2 {
				t.Fatalf("len(query.FilterOptions) = %d, want 2", len(query.FilterOptions))
			}

			if !query.FilterOptions[0] {
				t.Fatal("includeInactive filter = false, want true for hydration")
			}

			if query.FilterOptions[1] {
				t.Fatal("isFilterable filter = true, want false for hydration")
			}

			if query.Cursor != nil {
				t.Fatalf("query.Cursor = %+v, want nil", query.Cursor)
			}

			return definitions, nil
		},
	})

	got, err := service.GetProduct(context.Background(), product.ProductID)
	if err != nil {
		t.Fatalf("GetProduct() error = %v, want nil", err)
	}

	if got.ProductID != product.ProductID {
		t.Fatalf("ProductID = %q, want %q", got.ProductID, product.ProductID)
	}

	if got.CategoryID != product.CategoryID {
		t.Fatalf("CategoryID = %q, want %q", got.CategoryID, product.CategoryID)
	}

	if got.Name != product.Name {
		t.Fatalf("Name = %q, want %q", got.Name, product.Name)
	}

	if got.BasePrice != product.BasePrice {
		t.Fatalf("BasePrice = %+v, want %+v", got.BasePrice, product.BasePrice)
	}

	if got.Status != product.Status {
		t.Fatalf("Status = %q, want %q", got.Status, product.Status)
	}

	if len(got.Images) != 1 {
		t.Fatalf("len(Images) = %d, want 1", len(got.Images))
	}

	if len(got.Attributes) != 1 {
		t.Fatalf("len(Attributes) = %d, want 1 product-level attribute", len(got.Attributes))
	}

	productAttribute := got.Attributes[0]

	if productAttribute.ProductAttributeValueID != productLevelValueID {
		t.Fatalf("product attribute ID = %q, want %q", productAttribute.ProductAttributeValueID, productLevelValueID)
	}

	if productAttribute.Code != "material" {
		t.Fatalf("product attribute Code = %q, want material", productAttribute.Code)
	}

	if productAttribute.DisplayName != "Material" {
		t.Fatalf("product attribute DisplayName = %q, want Material", productAttribute.DisplayName)
	}

	if productAttribute.DataType != ProductAttributeDataTypeString {
		t.Fatalf("product attribute DataType = %q, want %q", productAttribute.DataType, ProductAttributeDataTypeString)
	}

	if productAttribute.ValueString != "steel" {
		t.Fatalf("product attribute ValueString = %q, want steel", productAttribute.ValueString)
	}

	if len(got.Variants) != 1 {
		t.Fatalf("len(Variants) = %d, want 1", len(got.Variants))
	}

	variant := got.Variants[0]

	if variant.VariantID != variantID {
		t.Fatalf("variant ID = %q, want %q", variant.VariantID, variantID)
	}

	if len(variant.Attributes) != 1 {
		t.Fatalf("len(variant.Attributes) = %d, want 1", len(variant.Attributes))
	}

	variantAttribute := variant.Attributes[0]

	if variantAttribute.ProductAttributeValueID != variantValueID {
		t.Fatalf("variant attribute ID = %q, want %q", variantAttribute.ProductAttributeValueID, variantValueID)
	}

	if variantAttribute.VariantID == nil || *variantAttribute.VariantID != variantID {
		t.Fatalf("variant attribute VariantID = %v, want %q", variantAttribute.VariantID, variantID)
	}

	if variantAttribute.Code != "colour" {
		t.Fatalf("variant attribute Code = %q, want colour", variantAttribute.Code)
	}

	if variantAttribute.ValueString != "blue" {
		t.Fatalf("variant attribute ValueString = %q, want blue", variantAttribute.ValueString)
	}

	if got.Variants[0].Attributes[0].Options[0].Value != "blue" {
		t.Fatalf("variant attribute option value = %q, want blue", got.Variants[0].Attributes[0].Options[0].Value)
	}
}

func TestService_GetProductReturnsNotFound(t *testing.T) {
	t.Parallel()

	service := NewService(fakeRepository{
		getProductFunc: func(ctx context.Context, productID ProductID) (Product, error) {
			return Product{}, ErrProductNotFound
		},
	})

	_, err := service.GetProduct(context.Background(), ProductID("prod_missing"))
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("GetProduct() error = %v, want ErrProductNotFound", err)
	}
}

func TestService_GetProductReturnsDefinitionLoadError(t *testing.T) {
	t.Parallel()

	definitionErr := errors.New("definition query failed")

	service := NewService(fakeRepository{
		getProductFunc: func(ctx context.Context, productID ProductID) (Product, error) {
			return Product{
				ProductID:  ProductID("prod_gopher_lamp"),
				CategoryID: CategoryID("cat_lighting"),
			}, nil
		},
		listProductAttributeDefinitionsFunc: func(ctx context.Context, query ListQuery) ([]ProductAttributeDefinition, error) {
			return nil, definitionErr
		},
	})

	_, err := service.GetProduct(context.Background(), ProductID("prod_gopher_lamp"))
	if err == nil {
		t.Fatal("GetProduct() error = nil, want definition load error")
	}

	if !errors.Is(err, definitionErr) {
		t.Fatalf("GetProduct() error = %v, want wrapped definition load error", err)
	}
}

func TestService_GetProductReturnsMissingDefinitionError(t *testing.T) {
	t.Parallel()

	service := NewService(fakeRepository{
		getProductFunc: func(ctx context.Context, productID ProductID) (Product, error) {
			return Product{
				ProductID:  ProductID("prod_gopher_lamp"),
				CategoryID: CategoryID("cat_lighting"),
				Attributes: []*ProductAttributeValue{
					{
						ProductAttributeValueID: ProductAttributeValueID("pav_missing"),
						ProductID:              ProductID("prod_gopher_lamp"),
						AttributeID:            AttributeID("attr_missing"),
						ValueString:            "unknown",
					},
				},
			}, nil
		},
		listProductAttributeDefinitionsFunc: func(ctx context.Context, query ListQuery) ([]ProductAttributeDefinition, error) {
			return []ProductAttributeDefinition{}, nil
		},
	})

	_, err := service.GetProduct(context.Background(), ProductID("prod_gopher_lamp"))
	if err == nil {
		t.Fatal("GetProduct() error = nil, want missing definition error")
	}
}

func TestHydrateProductAttributeValueRejectsNilValue(t *testing.T) {
	t.Parallel()

	_, err := hydrateProductAttributeValue(nil, map[AttributeID]ProductAttributeDefinition{})
	if err == nil {
		t.Fatal("hydrateProductAttributeValue() error = nil, want error")
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
