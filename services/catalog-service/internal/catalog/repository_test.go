package catalog

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockRepository(t *testing.T) (*MySQLRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock database: %v", err)
	}

	cleanup := func() {
		t.Helper()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}

		if err := db.Close(); err != nil {
			t.Fatalf("close mock database: %v", err)
		}
	}

	return &MySQLRepository{db: db}, mock, cleanup
}

func productColumns() []string {
	return []string{
		"product_id",
		"category_id",
		"name",
		"slug",
		"description",
		"brand",
		"status",
		"base_price_minor",
		"currency_code",
		"created_at",
		"updated_at",
	}
}

func productVariantColumns() []string {
	return []string{
		"variant_id",
		"product_id",
		"sku",
		"variant_name",
		"status",
		"price_minor",
		"currency_code",
		"created_at",
		"updated_at",
	}
}

func productAttributeValueColumns() []string {
	return []string{
		"product_attribute_value_id",
		"product_id",
		"variant_id",
		"attribute_id",
		"value_string",
		"value_number",
		"value_boolean",
		"value_json",
		"unit",
		"created_at",
		"updated_at",
	}
}

func productImageColumns() []string {
	return []string{
		"image_id",
		"product_id",
		"url",
		"alt_text",
		"display_order",
	}
}

func categoryColumns() []string {
	return []string{
		"category_id",
		"parent_category_id",
		"name",
		"slug",
		"description",
		"status",
		"display_order",
		"created_at",
		"updated_at",
	}
}

func attributeDefinitionColumns() []string {
	return []string{
		"attribute_id",
		"category_id",
		"code",
		"display_name",
		"description",
		"data_type",
		"unit",
		"is_required",
		"is_filterable",
		"is_variant_defining",
		"status",
		"created_at",
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatalf("parse time %q: %v", value, err)
	}

	return parsed
}

func TestMySQLRepository_ListProductsReturnsProducts(t *testing.T) {
	t.Parallel()

	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	createdAt := mustParseTime(t, "2026-06-08T10:00:00Z")
	updatedAt := mustParseTime(t, "2026-06-08T11:00:00Z")

	rows := sqlmock.NewRows(productColumns()).
		AddRow(
			"prod_001",
			"cat_lighting",
			"Gopher Desk Lamp",
			"gopher-desk-lamp",
			"A cheerful lamp for late-night debugging.",
			"Borough",
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		).
		AddRow(
			"prod_002",
			"cat_wall_decor",
			"Rob Pike Wall Tapestry",
			"rob-pike-wall-tapestry",
			"A tasteful wall tapestry for simple systems.",
			"Borough",
			"active",
			int64(7999),
			"GBP",
			createdAt.Add(-time.Hour),
			updatedAt,
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			product_id,
			category_id,
			name,
			slug,
			description,
			brand,
			status,
			base_price_minor,
			currency_code,
			created_at,
			updated_at
		FROM products
		WHERE (? = '' OR category_id = ?)
		ORDER BY created_at DESC, product_id DESC
		LIMIT ?
	`)).
		WithArgs("", "", 20).
		WillReturnRows(rows)

	products, err := repo.ListProducts(context.Background(), ListQuery{
		ID:    "",
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}

	if len(products) != 2 {
		t.Fatalf("len(products) = %d, want 2", len(products))
	}

	if products[0].ProductID != ProductID("prod_001") {
		t.Fatalf("first ProductID = %q, want prod_001", products[0].ProductID)
	}

	if products[0].CategoryID != CategoryID("cat_lighting") {
		t.Fatalf("first CategoryID = %q, want cat_lighting", products[0].CategoryID)
	}

	if products[0].Status != ProductStatusActive {
		t.Fatalf("first Status = %q, want %q", products[0].Status, ProductStatusActive)
	}

	if products[0].BasePrice.AmountMinor != 4999 {
		t.Fatalf("first AmountMinor = %d, want 4999", products[0].BasePrice.AmountMinor)
	}

	if products[0].BasePrice.CurrencyCode != "GBP" {
		t.Fatalf("first CurrencyCode = %q, want GBP", products[0].BasePrice.CurrencyCode)
	}
}

func TestMySQLRepository_ListProductsWrapsQueryError(t *testing.T) {
	t.Parallel()

	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	dbErr := errors.New("database unavailable")

	mock.ExpectQuery("FROM products").
		WithArgs("", "", 20).
		WillReturnError(dbErr)

	_, err := repo.ListProducts(context.Background(), ListQuery{
		ID:    "",
		Limit: 20,
	})
	if err == nil {
		t.Fatal("ListProducts() error = nil, want error")
	}

	if !errors.Is(err, dbErr) {
		t.Fatalf("ListProducts() error = %v, want wrapped database error", err)
	}
}

func TestMySQLRepository_GetProductReturnsProductWithChildren(t *testing.T) {
	t.Parallel()

	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	createdAt := mustParseTime(t, "2026-06-08T10:00:00Z")
	updatedAt := mustParseTime(t, "2026-06-08T11:00:00Z")

	productRows := sqlmock.NewRows(productColumns()).
		AddRow(
			"prod_001",
			"cat_lighting",
			"Gopher Desk Lamp",
			"gopher-desk-lamp",
			"A cheerful lamp for late-night debugging.",
			"Borough",
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM products").
		WithArgs("prod_001").
		WillReturnRows(productRows)

	variantRows := sqlmock.NewRows(productVariantColumns()).
		AddRow(
			"var_001",
			"prod_001",
			"BFS-GO-LAMP-001",
			"Default",
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM product_variants").
		WithArgs("prod_001").
		WillReturnRows(variantRows)

	attrRows := sqlmock.NewRows(productAttributeValueColumns()).
		AddRow(
			"pav_001",
			"prod_001",
			nil,
			"attr_material",
			"steel",
			nil,
			nil,
			nil,
			nil,
			createdAt,
			updatedAt,
		).
		AddRow(
			"pav_002",
			"prod_001",
			"var_001",
			"attr_colour",
			"blue",
			nil,
			nil,
			nil,
			nil,
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM product_attribute_values").
		WithArgs("prod_001").
		WillReturnRows(attrRows)

	imageRows := sqlmock.NewRows(productImageColumns()).
		AddRow(
			"img_001",
			"prod_001",
			"https://example.test/gopher-lamp.png",
			"Gopher Desk Lamp on a desk",
			10,
		)

	mock.ExpectQuery("FROM product_images").
		WithArgs("prod_001").
		WillReturnRows(imageRows)

	product, err := repo.GetProduct(context.Background(), "prod_001")
	if err != nil {
		t.Fatalf("GetProduct() error = %v, want nil", err)
	}

	if product.ProductID != ProductID("prod_001") {
		t.Fatalf("ProductID = %q, want prod_001", product.ProductID)
	}

	if len(product.Variants) != 1 {
		t.Fatalf("len(Variants) = %d, want 1", len(product.Variants))
	}

	if product.Variants[0].Status != ProductVariantStatusActive {
		t.Fatalf("variant Status = %q, want %q", product.Variants[0].Status, ProductVariantStatusActive)
	}

	if len(product.Attributes) != 2 {
		t.Fatalf("len(Attributes) = %d, want 2", len(product.Attributes))
	}

	if product.Attributes[0].VariantID != nil {
		t.Fatalf("first attribute VariantID = %v, want nil", product.Attributes[0].VariantID)
	}

	if product.Attributes[1].VariantID == nil || *product.Attributes[1].VariantID != VariantID("var_001") {
		t.Fatalf("second attribute VariantID = %v, want var_001", product.Attributes[1].VariantID)
	}

	if product.Attributes[0].ValueString == nil || *product.Attributes[0].ValueString != "steel" {
		t.Fatalf("first attribute ValueString = %v, want steel", product.Attributes[0].ValueString)
	}

	if len(product.Images) != 1 {
		t.Fatalf("len(Images) = %d, want 1", len(product.Images))
	}

	if product.Images[0].DisplayOrder != 10 {
		t.Fatalf("image DisplayOrder = %d, want 10", product.Images[0].DisplayOrder)
	}
}

func TestMySQLRepository_GetProductReturnsNotFound(t *testing.T) {
	t.Parallel()

	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	mock.ExpectQuery("FROM products").
		WithArgs("prod_missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetProduct(context.Background(), "prod_missing")
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("GetProduct() error = %v, want ErrProductNotFound", err)
	}
}

func TestMySQLRepository_GetProductWrapsChildQueryError(t *testing.T) {
	t.Parallel()

	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	createdAt := mustParseTime(t, "2026-06-08T10:00:00Z")
	updatedAt := mustParseTime(t, "2026-06-08T11:00:00Z")

	productRows := sqlmock.NewRows(productColumns()).
		AddRow(
			"prod_001",
			"cat_lighting",
			"Gopher Desk Lamp",
			"gopher-desk-lamp",
			"A cheerful lamp.",
			"Borough",
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM products").
		WithArgs("prod_001").
		WillReturnRows(productRows)

	dbErr := errors.New("variant query failed")

	mock.ExpectQuery("FROM product_variants").
		WithArgs("prod_001").
		WillReturnError(dbErr)

	_, err := repo.GetProduct(context.Background(), "prod_001")
	if err == nil {
		t.Fatal("GetProduct() error = nil, want error")
	}

	if !errors.Is(err, dbErr) {
		t.Fatalf("GetProduct() error = %v, want wrapped variant query error", err)
	}
}

func TestMySQLRepository_ListCategoriesReturnsCategories(t *testing.T) {
	t.Parallel()

	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	createdAt := mustParseTime(t, "2026-06-08T10:00:00Z")
	updatedAt := mustParseTime(t, "2026-06-08T11:00:00Z")

	rows := sqlmock.NewRows(categoryColumns()).
		AddRow(
			"cat_lighting",
			nil,
			"Lighting",
			"lighting",
			"Developer-themed lighting.",
			"active",
			10,
			createdAt,
			updatedAt,
		).
		AddRow(
			"cat_desk_lamps",
			"cat_lighting",
			"Desk Lamps",
			"desk-lamps",
			"Desk lamps and task lighting.",
			"active",
			20,
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM categories").
		WithArgs("", 20).
		WillReturnRows(rows)

	categories, err := repo.ListCategories(context.Background(), ListQuery{
		ID:    "",
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("ListCategories() error = %v, want nil", err)
	}

	if len(categories) != 2 {
		t.Fatalf("len(categories) = %d, want 2", len(categories))
	}

	if categories[0].ParentCategoryID != nil {
		t.Fatalf("first ParentCategoryID = %v, want empty", categories[0].ParentCategoryID)
	}

	if *categories[1].ParentCategoryID != CategoryID("cat_lighting") {
		t.Fatalf("second ParentCategoryID = %v, want cat_lighting", categories[1].ParentCategoryID)
	}

	if categories[0].Status != CategoryStatusActive {
		t.Fatalf("first Status = %q, want %q", categories[0].Status, CategoryStatusActive)
	}
}

func TestMySQLRepository_ListProductAttributeDefinitionsReturnsDefinitions(t *testing.T) {
	t.Parallel()

	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	createdAt := mustParseTime(t, "2026-06-08T10:00:00Z")

	rows := sqlmock.NewRows(attributeDefinitionColumns()).
		AddRow(
			"attr_material",
			"cat_lighting",
			"material",
			"Material",
			"Primary product material.",
			"string",
			nil,
			true,
			true,
			false,
			"active",
			createdAt,
		)

	mock.ExpectQuery("FROM product_attribute_definitions").
		WithArgs("cat_lighting", 20).
		WillReturnRows(rows)

	definitions, err := repo.ListProductAttributeDefinitions(context.Background(), ListQuery{
		ID:    "cat_lighting",
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("ListProductAttributeDefinitions() error = %v, want nil", err)
	}

	if len(definitions) != 1 {
		t.Fatalf("len(definitions) = %d, want 1", len(definitions))
	}

	if definitions[0].AttributeID != AttributeID("attr_material") {
		t.Fatalf("AttributeID = %q, want attr_material", definitions[0].AttributeID)
	}

	if definitions[0].DataType != ProductAttributeTypeString {
		t.Fatalf("DataType = %q, want %q", definitions[0].DataType, ProductAttributeTypeString)
	}

	if definitions[0].Status != ProductAttributeDefinitionStatusActive {
		t.Fatalf("Status = %q, want %q", definitions[0].Status, ProductAttributeDefinitionStatusActive)
	}
}
