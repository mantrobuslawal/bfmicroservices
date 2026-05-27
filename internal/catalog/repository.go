package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("catalogue record not found")

// Repository defines Catalogue Service persistence behaviour.
type Repository interface {
	ListProducts(ctx context.Context, filter ListProductsFilter) ([]Product, error)
	GetProduct(ctx context.Context, productID string) (Product, error)
	ListCategories(ctx context.Context, filter ListCategoriesFilter) ([]Category, error)
}

// MySQLRepository implements Repository using MySQL.
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a MySQL-backed catalogue repository.
func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// ListProducts returns products from the catalogue database.
func (r *MySQLRepository) ListProducts(ctx context.Context, filter ListProductsFilter) ([]Product, error) {
	limit := normaliseLimit(filter.Limit)
	offset := normaliseOffset(filter.Offset)

	query := `
SELECT
  product_id,
  category_id,
  name,
  slug,
  COALESCE(description, ''),
  COALESCE(brand, ''),
  status,
  base_price_minor,
  currency_code,
  created_at,
  updated_at
FROM products
WHERE (? = '' OR category_id = ?)
  AND (? = TRUE OR status = 'active')
ORDER BY created_at DESC
LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		filter.CategoryID,
		filter.CategoryID,
		filter.IncludeInactive,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var product Product

		if err := rows.Scan(
			&product.ProductID,
			&product.CategoryID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Brand,
			&product.Status,
			&product.BasePrice.AmountMinor,
			&product.BasePrice.CurrencyCode,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	return products, nil
}

// GetProduct returns one product by ID.
func (r *MySQLRepository) GetProduct(ctx context.Context, productID string) (Product, error) {
	query := `
SELECT
  product_id,
  category_id,
  name,
  slug,
  COALESCE(description, ''),
  COALESCE(brand, ''),
  status,
  base_price_minor,
  currency_code,
  created_at,
  updated_at
FROM products
WHERE product_id = ?
LIMIT 1`

	var product Product

	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ProductID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Brand,
		&product.Status,
		&product.BasePrice.AmountMinor,
		&product.BasePrice.CurrencyCode,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, fmt.Errorf("get product: %w", err)
	}

	return product, nil
}

// ListCategories returns categories from the catalogue database.
func (r *MySQLRepository) ListCategories(ctx context.Context, filter ListCategoriesFilter) ([]Category, error) {
	limit := normaliseLimit(filter.Limit)
	offset := normaliseOffset(filter.Offset)

	query := `
SELECT
  category_id,
  parent_category_id,
  name,
  slug,
  COALESCE(description, ''),
  status,
  display_order,
  created_at,
  updated_at
FROM categories
WHERE (? = '' OR parent_category_id = ?)
  AND (? = TRUE OR status = 'active')
ORDER BY display_order ASC, name ASC
LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		filter.ParentCategoryID,
		filter.ParentCategoryID,
		filter.IncludeInactive,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	categories := make([]Category, 0)

	for rows.Next() {
		var category Category
		var parentCategoryID sql.NullString

		if err := rows.Scan(
			&category.CategoryID,
			&parentCategoryID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Status,
			&category.DisplayOrder,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		if parentCategoryID.Valid {
			category.ParentCategoryID = &parentCategoryID.String
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}

	return categories, nil
}

func normaliseLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normaliseOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
