package catalog

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockMySQLRepository(t *testing.T) (*MySQLRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("create sqlmock database: %v", err)
	}

	cleanup := func() {
		t.Helper()

		mock.ExpectClose()

		if err := db.Close(); err != nil {
			t.Fatalf("close mock database: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}
	}

	return NewMySQLRepository(db), mock, cleanup
}

func testTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatalf("parse time %q: %v", value, err)
	}

	return parsed
}

func boolPtr(value bool) *bool {
	return &value
}

func productRowColumns() []string {
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

func productVariantRowColumns() []string {
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

func productAttributeValueRowColumns() []string {
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

func productImageRowColumns() []string {
	return []string{
		"image_id",
		"product_id",
		"url",
		"alt_text",
		"display_order",
	}
}

func categoryRowColumns() []string {
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

func productAttributeDefinitionRowColumns() []string {
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

func TestMySQLRepository_ListProductsReturnsProducts(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	rows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_gopher_lamp",
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
			"prod_rob_pike_tapestry",
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

	mock.ExpectQuery("FROM products").
		WithArgs("", "", false, 20).
		WillReturnRows(rows)

	products, err := repository.ListProducts(context.Background(), ListQuery{
		ID:            "",
		FilterOptions: []bool{false},
		Limit:         20,
		Cursor:        nil,
	})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}

	if got, want := len(products), 2; got != want {
		t.Fatalf("len(products) = %d, want %d", got, want)
	}

	first := products[0]

	if got, want := first.ProductID, ProductID("prod_gopher_lamp"); got != want {
		t.Fatalf("first.ProductID = %q, want %q", got, want)
	}

	if got, want := first.CategoryID, CategoryID("cat_lighting"); got != want {
		t.Fatalf("first.CategoryID = %q, want %q", got, want)
	}

	if got, want := first.Description, "A cheerful lamp for late-night debugging."; got != want {
		t.Fatalf("first.Description = %q, want %q", got, want)
	}

	if got, want := first.Brand, "Borough"; got != want {
		t.Fatalf("first.Brand = %q, want %q", got, want)
	}

	if got, want := first.Status, ProductStatusActive; got != want {
		t.Fatalf("first.Status = %q, want %q", got, want)
	}

	if got, want := first.BasePrice.AmountMinor, int64(4999); got != want {
		t.Fatalf("first.BasePrice.AmountMinor = %d, want %d", got, want)
	}

	if got, want := first.BasePrice.CurrencyCode, "GBP"; got != want {
		t.Fatalf("first.BasePrice.CurrencyCode = %q, want %q", got, want)
	}
}

func TestMySQLRepository_ListProductsHandlesNullableDescriptionAndBrand(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	rows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_plain_mug",
			"cat_homeware",
			"Plain Debug Mug",
			"plain-debug-mug",
			nil,
			nil,
			"active",
			int64(1299),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM products").
		WithArgs("cat_homeware", "cat_homeware", false, 10).
		WillReturnRows(rows)

	products, err := repository.ListProducts(context.Background(), ListQuery{
		ID:            "cat_homeware",
		FilterOptions: []bool{false},
		Limit:         10,
	})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}

	if got, want := len(products), 1; got != want {
		t.Fatalf("len(products) = %d, want %d", got, want)
	}

	if products[0].Description != "" {
		t.Fatalf("Description = %q, want empty string", products[0].Description)
	}

	if products[0].Brand != "" {
		t.Fatalf("Brand = %q, want empty string", products[0].Brand)
	}
}

func TestMySQLRepository_ListProductsWithCursor(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	cursorTime := testTime(t, "2026-06-08T10:00:00Z")

	mock.ExpectQuery("FROM products").
		WithArgs(
			"cat_lighting",
			"cat_lighting",
			true,
			cursorTime,
			cursorTime,
			"prod_cursor",
			5,
		).
		WillReturnRows(sqlmock.NewRows(productRowColumns()))

	products, err := repository.ListProducts(context.Background(), ListQuery{
		ID:            "cat_lighting",
		FilterOptions: []bool{true},
		Limit:         5,
		Cursor: &catalogCursor{
			CreatedAt: cursorTime,
			ID:        "prod_cursor",
		},
	})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}

	if got, want := len(products), 0; got != want {
		t.Fatalf("len(products) = %d, want %d", got, want)
	}
}

func TestMySQLRepository_ListProductsDefaultsMissingFilterOptionsToActiveOnly(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	mock.ExpectQuery("FROM products").
		WithArgs("", "", false, 10).
		WillReturnRows(sqlmock.NewRows(productRowColumns()))

	_, err := repository.ListProducts(context.Background(), ListQuery{
		ID:    "",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListProducts() error = %v, want nil", err)
	}
}

func TestMySQLRepository_ListProductsReturnsInvalidStatusError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	rows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_invalid",
			"cat_lighting",
			"Invalid Status Product",
			"invalid-status-product",
			nil,
			nil,
			"nonsense",
			int64(1000),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM products").
		WithArgs("", "", false, 10).
		WillReturnRows(rows)

	_, err := repository.ListProducts(context.Background(), ListQuery{
		FilterOptions: []bool{false},
		Limit:         10,
	})
	if err == nil {
		t.Fatal("ListProducts() error = nil, want invalid status error")
	}
}

func TestMySQLRepository_ListProductsWrapsQueryError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	databaseErr := errors.New("database unavailable")

	mock.ExpectQuery("FROM products").
		WithArgs("", "", false, 20).
		WillReturnError(databaseErr)

	_, err := repository.ListProducts(context.Background(), ListQuery{
		ID:            "",
		FilterOptions: []bool{false},
		Limit:         20,
	})
	if err == nil {
		t.Fatal("ListProducts() error = nil, want error")
	}

	if !errors.Is(err, databaseErr) {
		t.Fatalf("ListProducts() error = %v, want wrapped database error", err)
	}
}

func TestMySQLRepository_ListProductsWrapsRowsError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")
	rowsErr := errors.New("rows iteration failed")

	rows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_row_error",
			"cat_lighting",
			"Gopher Desk Lamp",
			"gopher-desk-lamp",
			"A cheerful lamp for late-night debugging",
			"Borough",
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		).
		RowError(0, rowsErr)

	mock.ExpectQuery("FROM products").
		WithArgs("", "", false, 20).
		WillReturnRows(rows)

	_, err := repository.ListProducts(context.Background(), ListQuery{
		ID:            "",
		FilterOptions: []bool{false},
		Limit:         20,
	})
	if err == nil {
		t.Fatal("ListProducts() error = nil, want rows error")
	}

	if !errors.Is(err, rowsErr) {
		t.Fatalf("ListProducts() error = %v, want wrapped rows error", err)
	}
}

func TestMySQLRepository_GetProductReturnsProductWithChildren(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	productRows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_gopher_lamp",
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
		WithArgs("prod_gopher_lamp").
		WillReturnRows(productRows)

	variantRows := sqlmock.NewRows(productVariantRowColumns()).
		AddRow(
			"var_blue",
			"prod_gopher_lamp",
			"BFS-GO-LAMP-BLUE",
			"Blue shade",
			"active",
			int64(5499),
			"GBP",
			createdAt,
			updatedAt,
		).
		AddRow(
			"var_green",
			"prod_gopher_lamp",
			"BFS-GO-LAMP-GREEN",
			"Green shade",
			"inactive",
			int64(5499),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM product_variants").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(variantRows)

	jsonPayload := []byte(`{"lumens":800,"power":"USB-C"}`)

	attributeRows := sqlmock.NewRows(productAttributeValueRowColumns()).
		AddRow(
			"pav_material",
			"prod_gopher_lamp",
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
			"pav_dimmable",
			"prod_gopher_lamp",
			nil,
			"attr_dimmable",
			nil,
			nil,
			true,
			nil,
			nil,
			createdAt,
			updatedAt,
		).
		AddRow(
			"pav_weight",
			"prod_gopher_lamp",
			nil,
			"attr_weight",
			nil,
			"1.2500",
			nil,
			nil,
			"kg",
			createdAt,
			updatedAt,
		).
		AddRow(
			"pav_specs",
			"prod_gopher_lamp",
			"var_blue",
			"attr_specs",
			nil,
			nil,
			nil,
			jsonPayload,
			nil,
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM product_attribute_values").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(attributeRows)

	imageRows := sqlmock.NewRows(productImageRowColumns()).
		AddRow(
			"img_primary",
			"prod_gopher_lamp",
			"https://example.test/gopher-lamp-primary.png",
			"Gopher Desk Lamp on a desk",
			10,
		).
		AddRow(
			"img_detail",
			"prod_gopher_lamp",
			"https://example.test/gopher-lamp-detail.png",
			"Close-up of the Gopher Desk Lamp shade",
			20,
		)

	mock.ExpectQuery("FROM product_images").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(imageRows)

	product, err := repository.GetProduct(context.Background(), ProductID("prod_gopher_lamp"))
	if err != nil {
		t.Fatalf("GetProduct() error = %v, want nil", err)
	}

	if got, want := product.ProductID, ProductID("prod_gopher_lamp"); got != want {
		t.Fatalf("ProductID = %q, want %q", got, want)
	}

	if got, want := len(product.Variants), 2; got != want {
		t.Fatalf("len(Variants) = %d, want %d", got, want)
	}

	if got, want := product.Variants[0].VariantID, VariantID("var_blue"); got != want {
		t.Fatalf("first variant ID = %q, want %q", got, want)
	}

	if got, want := product.Variants[0].Status, ProductVariantStatusActive; got != want {
		t.Fatalf("first variant status = %q, want %q", got, want)
	}

	if got, want := product.Variants[0].Price.AmountMinor, int64(5499); got != want {
		t.Fatalf("first variant price amount = %d, want %d", got, want)
	}

	if got, want := len(product.Attributes), 4; got != want {
		t.Fatalf("len(Attributes) = %d, want %d", got, want)
	}

	if product.Attributes[0].VariantID != nil {
		t.Fatalf("first attribute VariantID = %v, want nil", product.Attributes[0].VariantID)
	}

	if got, want := product.Attributes[0].ValueString, "steel"; got != want {
		t.Fatalf("first attribute ValueString = %q, want %q", got, want)
	}

	if product.Attributes[1].ValueBoolean == nil || *product.Attributes[1].ValueBoolean != true {
		t.Fatalf("second attribute ValueBoolean = %v, want true", product.Attributes[1].ValueBoolean)
	}

	if got, want := product.Attributes[2].ValueNumber, "1.2500"; got != want {
		t.Fatalf("third attribute ValueNumber = %q, want %q", got, want)
	}

	if got, want := product.Attributes[2].Unit, "kg"; got != want {
		t.Fatalf("third attribute Unit = %q, want %q", got, want)
	}

	if product.Attributes[3].VariantID == nil || *product.Attributes[3].VariantID != VariantID("var_blue") {
		t.Fatalf("fourth attribute VariantID = %v, want var_blue", product.Attributes[3].VariantID)
	}

	if string(product.Attributes[3].ValueJSON) != string(jsonPayload) {
		t.Fatalf("fourth attribute ValueJSON = %s, want %s", product.Attributes[3].ValueJSON, jsonPayload)
	}

	if got, want := len(product.Images), 2; got != want {
		t.Fatalf("len(Images) = %d, want %d", got, want)
	}

	if got, want := product.Images[0].ProductID, ProductID("prod_gopher_lamp"); got != want {
		t.Fatalf("first image ProductID = %q, want %q", got, want)
	}

	if got, want := product.Images[0].AltText, "Gopher Desk Lamp on a desk"; got != want {
		t.Fatalf("first image AltText = %q, want %q", got, want)
	}

	if got, want := product.Images[1].AltText, "Close-up of the Gopher Desk Lamp shade"; got != want {
		t.Fatalf("second image AltText = %q, want %q", got, want)
	}
}

func TestMySQLRepository_GetProductReturnsNotFound(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	mock.ExpectQuery("FROM products").
		WithArgs("prod_missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repository.GetProduct(context.Background(), ProductID("prod_missing"))
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("GetProduct() error = %v, want ErrProductNotFound", err)
	}
}

func TestMySQLRepository_GetProductWrapsProductQueryError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	databaseErr := errors.New("connection lost")

	mock.ExpectQuery("FROM products").
		WithArgs("prod_gopher_lamp").
		WillReturnError(databaseErr)

	_, err := repository.GetProduct(context.Background(), ProductID("prod_gopher_lamp"))
	if err == nil {
		t.Fatal("GetProduct() error = nil, want error")
	}

	if !errors.Is(err, databaseErr) {
		t.Fatalf("GetProduct() error = %v, want wrapped database error", err)
	}
}

func TestMySQLRepository_GetProductWrapsVariantQueryError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	productRows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_gopher_lamp",
			"cat_lighting",
			"Gopher Desk Lamp",
			"gopher-desk-lamp",
			nil,
			nil,
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM products").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(productRows)

	databaseErr := errors.New("variant query failed")

	mock.ExpectQuery("FROM product_variants").
		WithArgs("prod_gopher_lamp").
		WillReturnError(databaseErr)

	_, err := repository.GetProduct(context.Background(), ProductID("prod_gopher_lamp"))
	if err == nil {
		t.Fatal("GetProduct() error = nil, want error")
	}

	if !errors.Is(err, databaseErr) {
		t.Fatalf("GetProduct() error = %v, want wrapped variant query error", err)
	}
}

func TestMySQLRepository_GetProductWrapsAttributeValueQueryError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	productRows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_gopher_lamp",
			"cat_lighting",
			"Gopher Desk Lamp",
			"gopher-desk-lamp",
			nil,
			nil,
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM products").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(productRows)

	mock.ExpectQuery("FROM product_variants").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(sqlmock.NewRows(productVariantRowColumns()))

	databaseErr := errors.New("attribute value query failed")

	mock.ExpectQuery("FROM product_attribute_values").
		WithArgs("prod_gopher_lamp").
		WillReturnError(databaseErr)

	_, err := repository.GetProduct(context.Background(), ProductID("prod_gopher_lamp"))
	if err == nil {
		t.Fatal("GetProduct() error = nil, want error")
	}

	if !errors.Is(err, databaseErr) {
		t.Fatalf("GetProduct() error = %v, want wrapped attribute value query error", err)
	}
}

func TestMySQLRepository_GetProductWrapsImageQueryError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	productRows := sqlmock.NewRows(productRowColumns()).
		AddRow(
			"prod_gopher_lamp",
			"cat_lighting",
			"Gopher Desk Lamp",
			"gopher-desk-lamp",
			nil,
			nil,
			"active",
			int64(4999),
			"GBP",
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM products").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(productRows)

	mock.ExpectQuery("FROM product_variants").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(sqlmock.NewRows(productVariantRowColumns()))

	mock.ExpectQuery("FROM product_attribute_values").
		WithArgs("prod_gopher_lamp").
		WillReturnRows(sqlmock.NewRows(productAttributeValueRowColumns()))

	databaseErr := errors.New("image query failed")

	mock.ExpectQuery("FROM product_images").
		WithArgs("prod_gopher_lamp").
		WillReturnError(databaseErr)

	_, err := repository.GetProduct(context.Background(), ProductID("prod_gopher_lamp"))
	if err == nil {
		t.Fatal("GetProduct() error = nil, want error")
	}

	if !errors.Is(err, databaseErr) {
		t.Fatalf("GetProduct() error = %v, want wrapped image query error", err)
	}
}

func TestMySQLRepository_ListCategoriesReturnsCategories(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")
	updatedAt := testTime(t, "2026-06-08T11:00:00Z")

	rows := sqlmock.NewRows(categoryRowColumns()).
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
			nil,
			"active",
			20,
			createdAt,
			updatedAt,
		)

	mock.ExpectQuery("FROM categories").
		WithArgs("", "", false, 20).
		WillReturnRows(rows)

	categories, err := repository.ListCategories(context.Background(), ListQuery{
		ID:            "",
		FilterOptions: []bool{false},
		Limit:         20,
	})
	if err != nil {
		t.Fatalf("ListCategories() error = %v, want nil", err)
	}

	if got, want := len(categories), 2; got != want {
		t.Fatalf("len(categories) = %d, want %d", got, want)
	}

	root := categories[0]

	if got, want := root.CategoryID, CategoryID("cat_lighting"); got != want {
		t.Fatalf("root.CategoryID = %q, want %q", got, want)
	}

	if root.ParentCategoryID != nil {
		t.Fatalf("root.ParentCategoryID = %v, want nil", root.ParentCategoryID)
	}

	if got, want := root.Description, "Developer-themed lighting."; got != want {
		t.Fatalf("root.Description = %q, want %q", got, want)
	}

	if got, want := root.Status, CategoryStatusActive; got != want {
		t.Fatalf("root.Status = %q, want %q", got, want)
	}

	child := categories[1]

	if child.ParentCategoryID == nil || *child.ParentCategoryID != CategoryID("cat_lighting") {
		t.Fatalf("child.ParentCategoryID = %v, want cat_lighting", child.ParentCategoryID)
	}

	if child.Description != "" {
		t.Fatalf("child.Description = %q, want empty string", child.Description)
	}
}

func TestMySQLRepository_ListCategoriesWithCursor(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	cursorTime := testTime(t, "2026-06-08T10:00:00Z")

	mock.ExpectQuery("FROM categories").
		WithArgs(
			"cat_lighting",
			"cat_lighting",
			true,
			cursorTime,
			cursorTime,
			"cat_desk_lamps",
			10,
		).
		WillReturnRows(sqlmock.NewRows(categoryRowColumns()))

	categories, err := repository.ListCategories(context.Background(), ListQuery{
		ID:            "cat_lighting",
		FilterOptions: []bool{true},
		Limit:         10,
		Cursor: &catalogCursor{
			CreatedAt: cursorTime,
			ID:        "cat_desk_lamps",
		},
	})
	if err != nil {
		t.Fatalf("ListCategories() error = %v, want nil", err)
	}

	if got, want := len(categories), 0; got != want {
		t.Fatalf("len(categories) = %d, want %d", got, want)
	}
}

func TestMySQLRepository_ListCategoriesDefaultsMissingFilterOptionsToActiveOnly(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	mock.ExpectQuery("FROM categories").
		WithArgs("", "", false, 10).
		WillReturnRows(sqlmock.NewRows(categoryRowColumns()))

	_, err := repository.ListCategories(context.Background(), ListQuery{
		ID:    "",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListCategories() error = %v, want nil", err)
	}
}

func TestMySQLRepository_ListCategoriesWrapsQueryError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	databaseErr := errors.New("database unavailable")

	mock.ExpectQuery("FROM categories").
		WithArgs("", "", false, 20).
		WillReturnError(databaseErr)

	_, err := repository.ListCategories(context.Background(), ListQuery{
		ID:            "",
		FilterOptions: []bool{false},
		Limit:         20,
	})
	if err == nil {
		t.Fatal("ListCategories() error = nil, want error")
	}

	if !errors.Is(err, databaseErr) {
		t.Fatalf("ListCategories() error = %v, want wrapped database error", err)
	}
}

func TestMySQLRepository_ListProductAttributeDefinitionsReturnsDefinitions(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")

	rows := sqlmock.NewRows(productAttributeDefinitionRowColumns()).
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
		).
		AddRow(
			"attr_weight",
			"cat_lighting",
			"weight",
			"Weight",
			nil,
			"number",
			"kg",
			false,
			true,
			false,
			"active",
			createdAt.Add(-time.Minute),
		)

	mock.ExpectQuery("FROM product_attribute_definitions").
		WithArgs("cat_lighting", false, false, 20).
		WillReturnRows(rows)

	definitions, err := repository.ListProductAttributeDefinitions(context.Background(), ListQuery{
		ID:            "cat_lighting",
		FilterOptions: []bool{false, false},
		Limit:         20,
	})
	if err != nil {
		t.Fatalf("ListProductAttributeDefinitions() error = %v, want nil", err)
	}

	if got, want := len(definitions), 2; got != want {
		t.Fatalf("len(definitions) = %d, want %d", got, want)
	}

	first := definitions[0]

	if got, want := first.AttributeID, AttributeID("attr_material"); got != want {
		t.Fatalf("first.AttributeID = %q, want %q", got, want)
	}

	if got, want := first.Description, "Primary product material."; got != want {
		t.Fatalf("first.Description = %q, want %q", got, want)
	}

	if got, want := first.DataType, ProductAttributeDataTypeString; got != want {
		t.Fatalf("first.DataType = %q, want %q", got, want)
	}

	if first.Unit != "" {
		t.Fatalf("first.Unit = %q, want empty string", first.Unit)
	}

	second := definitions[1]

	if second.Description != "" {
		t.Fatalf("second.Description = %q, want empty string", second.Description)
	}

	if got, want := second.Unit, "kg"; got != want {
		t.Fatalf("second.Unit = %q, want %q", got, want)
	}

	if got, want := second.DataType, ProductAttributeDataTypeNumber; got != want {
		t.Fatalf("second.DataType = %q, want %q", got, want)
	}
}

func TestMySQLRepository_ListProductAttributeDefinitionsWithFiltersAndCursor(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	cursorTime := testTime(t, "2026-06-08T10:00:00Z")

	mock.ExpectQuery("FROM product_attribute_definitions").
		WithArgs(
			"cat_lighting",
			true,
			true,
			cursorTime,
			cursorTime,
			"attr_material",
			10,
		).
		WillReturnRows(sqlmock.NewRows(productAttributeDefinitionRowColumns()))

	definitions, err := repository.ListProductAttributeDefinitions(context.Background(), ListQuery{
		ID:            "cat_lighting",
		FilterOptions: []bool{true, true},
		Limit:         10,
		Cursor: &catalogCursor{
			CreatedAt: cursorTime,
			ID:        "attr_material",
		},
	})
	if err != nil {
		t.Fatalf("ListProductAttributeDefinitions() error = %v, want nil", err)
	}

	if got, want := len(definitions), 0; got != want {
		t.Fatalf("len(definitions) = %d, want %d", got, want)
	}
}

func TestMySQLRepository_ListProductAttributeDefinitionsDefaultsMissingFilterOptions(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	mock.ExpectQuery("FROM product_attribute_definitions").
		WithArgs("cat_lighting", false, false, 10).
		WillReturnRows(sqlmock.NewRows(productAttributeDefinitionRowColumns()))

	_, err := repository.ListProductAttributeDefinitions(context.Background(), ListQuery{
		ID:    "cat_lighting",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListProductAttributeDefinitions() error = %v, want nil", err)
	}
}

func TestMySQLRepository_ListProductAttributeDefinitionsReturnsInvalidDataTypeError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	createdAt := testTime(t, "2026-06-08T10:00:00Z")

	rows := sqlmock.NewRows(productAttributeDefinitionRowColumns()).
		AddRow(
			"attr_invalid",
			"cat_lighting",
			"invalid",
			"Invalid",
			nil,
			"mystery",
			nil,
			false,
			false,
			false,
			"active",
			createdAt,
		)

	mock.ExpectQuery("FROM product_attribute_definitions").
		WithArgs("cat_lighting", false, false, 10).
		WillReturnRows(rows)

	_, err := repository.ListProductAttributeDefinitions(context.Background(), ListQuery{
		ID:            "cat_lighting",
		FilterOptions: []bool{false, false},
		Limit:         10,
	})
	if err == nil {
		t.Fatal("ListProductAttributeDefinitions() error = nil, want invalid data type error")
	}
}

func TestMySQLRepository_ListProductAttributeDefinitionsWrapsQueryError(t *testing.T) {
	t.Parallel()

	repository, mock, cleanup := newMockMySQLRepository(t)
	defer cleanup()

	databaseErr := errors.New("database unavailable")

	mock.ExpectQuery("FROM product_attribute_definitions").
		WithArgs("cat_lighting", false, false, 20).
		WillReturnError(databaseErr)

	_, err := repository.ListProductAttributeDefinitions(context.Background(), ListQuery{
		ID:            "cat_lighting",
		FilterOptions: []bool{false, false},
		Limit:         20,
	})
	if err == nil {
		t.Fatal("ListProductAttributeDefinitions() error = nil, want error")
	}

	if !errors.Is(err, databaseErr) {
		t.Fatalf("ListProductAttributeDefinitions() error = %v, want wrapped database error", err)
	}
}
