package integration_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/acme-ltd/bfstore/services/catalog-service/internal/catalog"
)

// These tests are intended to run against the local MySQL database started by:
//
//   make up
//
// They expect the catalogue migration and seed data to have been applied.
//
// Run:
//
//   BFSTORE_RUN_INTEGRATION_TESTS=true go test ./test/integration/...
func TestCatalogRepositoryListProducts(t *testing.T) {
	if os.Getenv("BFSTORE_RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("set BFSTORE_RUN_INTEGRATION_TESTS=true to run integration tests")
	}

	db := openTestDB(t)
	defer db.Close()

	repository := catalog.NewMySQLRepository(db)

	products, err := repository.ListProducts(context.Background(), catalog.ListProductsFilter{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) == 0 {
		t.Fatal("expected seeded products, got none")
	}
}

func TestCatalogRepositoryGetProduct(t *testing.T) {
	if os.Getenv("BFSTORE_RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("set BFSTORE_RUN_INTEGRATION_TESTS=true to run integration tests")
	}

	db := openTestDB(t)
	defer db.Close()

	repository := catalog.NewMySQLRepository(db)

	product, err := repository.GetProduct(context.Background(), "cccccccc-cccc-cccc-cccc-cccccccc0001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Name != "Gopher Desk Lamp" {
		t.Fatalf("expected Gopher Desk Lamp, got %s", product.Name)
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := getenv("CATALOG_TEST_MYSQL_DSN", "bfstore_catalog_user:bfstore_catalog_password@tcp(localhost:3306)/bfstore_catalog?parseTime=true&charset=utf8mb4,utf8")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	return db
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
