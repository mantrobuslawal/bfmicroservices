package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Repository defines Catalogue Service persistence behaviour.
type Repository interface {
	ListProducts(ctx context.Context, query ListProductsFilter) ([]Product, error)

	GetProduct(ctx context.Context, productID string) (Product, error)

	ListCategories(ctx context.Context, query ListCategoriesFilter) ([]Category, error)

	ListProductAttributeDefinitions(ctx context.Context, query ListProductAttributeDefinitionsFilter) ([]ProductAttributeDefinition, error)
}

// Repository stores and retrieves catalogue data from MySQL.
type Repository struct{ db *sql.DB }

// ErrProductNotFound is returned when a product does not exist.
var ErrProductNotFound = errors.New("product not found")

// NewMRepository creates a MySQL-backed catalogue repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ListProducts returns products from the catalogue database.
//
// Returns ErrProductNotFound when produccts not found.
// Low-level database errors are wrapped to provide extra context.
func (r *Repository) ListProducts(ctx context.Context, query ListProductsFilter) ([]Product, error) {
	args := []any{query.ID, query.ID}

	sql := `
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
  	`

	if query.Cursor != nil {
		sql += `
  			AND (? = TRUE OR status = 'active') 
  			AND (
				created_at < ?
				OR (created_at = ? AND id < ?)
	  		)
		`

		args = append(
			args,
			query.FilterOptions[0],
			query.Cursor.CreatedAt,
			query.Cursor.CreatedAt,
			query.Cursor.ID,
		)
	}

	sql += `
		ORDER BY created_at DESC, id DESC
		LIMIT ?
	`

	args = append(args, query.Limit)

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errof("list products: %w", err)
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
			return nil, fmt.Errorf("scan product row: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product rows: %w", err)
	}

	return products, nil
}

// GetProduct returns one product by ID, along with its variants, attributes and images.
func (r *Repository) GetProduct(ctx context.Context, productID string) (Product, error) {
	product, err := r.getProductRow(ctx, productID)
	if err != nil {
		return Product{}, err
	}

	variants, err := r.listProductVariants(ctx, productID)
	if err != nil {
		return Product{}, fmt.Error("list product variants for product %q: %w", productID, err)
	}

	attributes, err := r.listProductAttributeValues(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("list product attribute values for %q: %w", productID, err)
	}

	images, err := r.listProductImages(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("list product images for product %q: %w", productID, err)
	}

	product.Variants = variants
	product.Attributes = attributes
	product.Images = images

	return product, nil
}

// ListCategories returns categories from the catalogue database.
func (r *Repository) ListCategories(ctx context.Context, query ListCategoriesFilter) ([]Category, error) {
	args := []any{query.ID}

	sql := `
	       SELECT *
	       FROM categories
	       WHERE (? = '' OR parent_category_id = ?)
	`

	if query.Cursor != nil {
		sql += `
  			AND (? = TRUE OR status = 'active') 
  			AND (
				created_at < ?
				OR (created_at = ? AND id < ?)
	  		)
		`

		args = append(
			args,
			query.FilterOptions[0],
			query.Cursor.CreatedAt,
			query.Cursor.CreatedAt,
			query.Cursor.ID,
		)
	}

	sql += `
		ORDER BY created_at DESC, id DESC
		LIMIT ?
	`

	args = append(args, query.Limit)

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if parentCategoryID.Valid {
			category.ParentCategoryID = &parentCategoryID.String
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// ListProductAttributeDefinitons returns product attribute definitions.
func (r *Repository) ListProductAttributeDefinitions(ctx context.Context, query ListProductAttributeDefinitionsFilter) ([]ProductAttributeDefinition, error) {

	args := []any{query.ID}

	sql := `
		SELECT *
		FROM product_attribute_definitions
		WHERE (category_id = ?)
	`

	if query.Cursor != nil {
		sql += `
   		AND (? = TRUE OR is_filterable = FALSE)
  		AND (? = TRUE OR status = 'active') 
  		AND (
			created_at < ?
			OR (created_at = ? AND id < ?)
	  	)`

		args = append(
			args,
			query.FilterOptions[1],
			query.FilterOptions[0],
			query.Cursor.CreatedAt,
			query.Cursor.CreatedAt,
			query.Cursor.ID,
		)
	}

	sql += `
		ORDER BY created_at DESC, id DESC
		LIMIT ?
	`

	args = append(args, query.Limit)

	productAttributeDefinitions := make([]ProductAttributeDefinition, 0)

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productAttributeDefinition ProductAttributeDefinition

		if err := rows.Scan(
			&productAttributeDefinition.AttributeID,
			&productAttributeDefinition.CategoryID,
			&productAttributeDefinition.Code,
			&productAttributeDefinition.DisplayName,
			&productAttributeDefinition.Description,
			&productAttributeDefinition.DataType,
			&productAttributeDefinition.Unit,
			&productAttributeDefinition.IsRequired,
			&productAttributeDefinition.IsFilterable,
			&productAttributeDefinition.IsVariantDefining,
			&productAttributeDefinition.Status,
			&productAttributeDefinition.CreatedAt,
			&productAttributeDefinition.UpdatedAt,
		); err != nil {
			return nil, err
		}

		productAttributeDefinitons = append(productAttributeDefinitons, productAttributeDefinition)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productAttributeDefinitions, nil
}

func (r *Respository) getProductRow(ctx context.Context, productID string) (Product, error) {
	const query = `
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
          upddated_at,
       FROM products
       WHERE product_id = ?    
  `

	var product Product
	var rawStatus string
	var description sql.NullString
	var brand sql.NullString
	var amountMinor int64
	var currencyCode string

	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ProductID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&description,
		&brand,
		&rawStatus,
		&amountMinor,
		&currencyCode,
		&product.CreatedAt,
		&productUpdatedAt,
	)

	if err != nil {
		if error.Is(err, sql.ErrNoRows) {
			return Product{}, ErrProductNotFound
		}

		return Product{}, fmt.Errorf("query product %q: %w", productID, err)
	}

	status, err := ParseLifecycleStatus(rawStatus)
	if err != nil {
		return Product{}, fmt.Errorf("parse product status for product %q: %w", productID, err)
	}

	product.Status = status
	product.Description = description.String
	product.Brand = brand.String
	product.BasePrice = Money{
		AmountMinor:  amountMinor,
		CurrencyCode: currencyCode,
	}

	return product, nil
}

func (r *Respository) listProductVariants(ctx context.Context, productID string) ([]ProductVariant, error) {
	const query = `
	SELECT
          variant_id,
          product_id,
          sku,
          variant_name,
          status,
          price_minor,
          currency_code,
          created_at,
          updated_at
        FROM product_variants
	WHERE product_id = ?
        ORDER BY variant_name ASC, variant_id ASC
   `

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product variants: %w", err)
	}
	defer rows.Close()

	var variants []ProductVariant

	for rows.Next() {
		var variant ProductVariant
		var rawStatus string
		var priceMinor int64
		var currenyCode string

		if err := rows.Scan(
			&variant.VariantID,
			&variant.ProductID,
			&variant.Sku,
			&variant.VariantName,
			&rawStatus,
			&priceMinor,
			&currencyCode,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product variant row: %w", err)
		}

		status, err := ParseLifecycleStatus(rawStatus)
		if err != nil {
			return nil, fmt.Errorf("parse product variant status for variant %q: %w", variant.VariantID, err)
		}

		variant.Status = status
		variant.Price = Money{
			AmountMinor:  priceMinor,
			CurrencyCode: currencyCode,
		}

		variants = append(variants, variant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product variant rows: %w", err)
	}

	return variants, nil
}

func (r *Respository) listProductAttributeValues(ctx context.Context, productID string) ([]ProductAttributeValue, error) {
	const query = `
     SELECT
        product_attribute_value_id,
        product_id,
        variant_id,
        attribute_id,
        value_string,
        CAST(value_number AS CHAR),
        value_boolean,
        value_json,
        unit,
        created_at,
        update_at
     FROM product_attribute_values
     WHERE product_id = ?
     ORDER BY attribute_id ASC, variant_id ASC, product_attribute_value_id ASC 	
  `

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product attribute values: %w", err)
	}
	defer rows.Close()

	var values []ProductAttributeValue

	for rows.Next() {
		var value ProductAttributeValue

		var variantID sql.NullString
		var valueString sql.NullString
		var valueNumber sql.NullString
		var valueBoolean sql.NullBool
		var valueJSON []byte
		var unit sql.NullString

		if err := rows.Scan(
			&value.ProductAttributeID,
			&value.ProductID,
			&variantID,
			&value.AttributeID,
			&valueString,
			&valueNumber,
			&valueBoolean,
			&valueJSON,
			&unit,
			&value.CreatedAt,
			&value.UpdatedAt,
		); err != nil {
			return nil, fmt.Errof("scan product attribute value row: %w", err)
		}

		if variantID.Valid {
			value.VariantID = &variantID.String
		}

		if valueString.Valid {
			value.ValueString = &valueString.String
		}

		if valueNumber.Valid {
			value.Number = &valueNumber.String
		}

		if valueBoolean.Valid {
			value.ValueBoolean = &valueBoolean.Bool
		}

		if len(valueJSON) > 0 {
			value.ValueJSON = valueJSON
		}

		if unit.Valid {
			value.Unit = &unit.String
		}

		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product attribute value rows: %w", err)
	}

	return values, nil
}

func (r *Respository) listProductImages(ctx context.Context, productID string) ([]ProductImage, error) {
	const query = `
	   SELECT
	     image_id,
             product_id,
             url,
             alt_text,
             display_order,
             is_primary,
             created_at,
             updated_at
           FROM product_images
           WHERE product_id = ?
           ORDER BY display_order ASC, image_id ASC   
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product images: %w", err)
	}
	defer rows.Close()

	var images []ProductImage

	for rows.Next() {
		var image ProductImage

		if err := rows.Scan(
			&image.ImageID,
			&image.ProductID,
			&image.URL,
			&image.AlText,
			&image.DisplayOrder,
			&image.IsPrimary,
			&image.CreatedAt,
			&image.UpdateAt,
		); err != nil {
			return nil, fmt.Errorf("scan product image row: %w", err)
		}

		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product image rows: %w", err)
	}

	return images, nil
}
