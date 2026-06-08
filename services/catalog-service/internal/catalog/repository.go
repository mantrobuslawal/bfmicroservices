package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Repository defines Catalogue Service persistence behaviour.
type Repository interface {
	ListProducts(ctx context.Context, query ListQuery) ([]Product, error)

	GetProduct(ctx context.Context, productID ProductID) (Product, error)

	ListCategories(ctx context.Context, query ListQuery) ([]Category, error)

	ListProductAttributeDefinitions(ctx context.Context, query ListQuery) ([]ProductAttributeDefinition, error)
}

// Repository stores and retrieves catalogue data from MySQL.
type MySQLRepository struct{ db *sql.DB }

var _ Repository = (*MySQLRepository)(nil)

// NewMRepository creates a MySQL-backed catalogue repository.
func MySQLNewRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// ListProducts returns products from the catalogue database.
//
// Returns ErrProductNotFound when produccts not found.
// Low-level database errors are wrapped to provide extra context.
func (r *MySQLRepository) ListProducts(ctx context.Context, query ListQuery) ([]Product, error) {
	args := []any{query.ID, query.ID}

	sqlText := `
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
		AND (? = TRUE OR status = 'active') 
  	`

	if query.Cursor != nil {
		sqlText += `
  			AND (
				created_at < ?
				OR (created_at = ? AND product_id < ?)
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

	sqlText += `
		ORDER BY created_at DESC, product_id DESC
		LIMIT ?
	`

	args = append(args, query.Limit)

	rows, err := r.db.QueryContext(ctx, sqlText, args...)
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
func (r *MySQLRepository) GetProduct(ctx context.Context, productID ProductID) (Product, error) {
	product, err := r.getProductRow(ctx, productID)
	if err != nil {
		return Product{}, err
	}

	variants, err := r.listProductVariants(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("list product variants for product %q: %w", productID, err)
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
func (r *MySQLRepository) ListCategories(ctx context.Context, query ListQuery) ([]Category, error) {
	args := []any{query.ID, query.ID}

	sqlText := `
	       SELECT 
		   category_id,
		   parent_category_id,
		   name,
		   slug,
		   description,
		   status,
		   display_order,
		   created_at,
		   updated_at
	       FROM categories
	       WHERE (? = '' OR parent_category_id = ?)
		   AND (? = TRUE OR status = 'active') 
	`

	if query.Cursor != nil {
		sqlText += `
  			AND (
				created_at < ?
				OR (created_at = ? AND category_id < ?)
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

	sqlText += `
		ORDER BY created_at DESC, category_id DESC
		LIMIT ?
	`

	args = append(args, query.Limit)

	rows, err := r.db.QueryContext(ctx, sqlText, args...)
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
			parentId := CategoryID(parentCategoryID.String)
			category.ParentCategoryID = &parentId
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// ListProductAttributeDefinitons returns product attribute definitions.
func (r *MySQLRepository) ListProductAttributeDefinitions(ctx context.Context, query ListQuery) ([]ProductAttributeDefinition, error) {

	args := []any{query.ID}

	sqlText := `
		SELECT 
		attribute_id,
		category_id,
		code,
		display_name,
		description,
		data_type,
		unit,
		is_required,
		is_filterable,
		is_variant_defining,
		status,
		created_at
		FROM product_attribute_definitions
		WHERE (category_id = ?)
		AND (? = FALSE OR is_filterable = TRUE)
  		AND (? = TRUE OR status = 'active') 
	`

	if query.Cursor != nil {
		sqlText += `
  		AND (
			created_at < ?
			OR (created_at = ? AND attribute_id < ?)
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

	sqlText += `
		ORDER BY created_at DESC, attribute_id DESC
		LIMIT ?
	`

	args = append(args, query.Limit)

	productAttributeDefinitions := make([]ProductAttributeDefinition, 0)

	rows, err := r.db.QueryContext(ctx, sqlText, args...)
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
			&productAttributeDefinition.CreatedAt); err != nil {
			return nil, err
		}

		productAttributeDefinitions = append(productAttributeDefinitions, productAttributeDefinition)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productAttributeDefinitions, nil
}

func (r *MySQLRepository) getProductRow(ctx context.Context, productID ProductID) (Product, error) {
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
       updated_at
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
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, ErrProductNotFound
		}

		return Product{}, fmt.Errorf("query product %q: %w", productID, err)
	}

	status, err := ParseToProductStatus(rawStatus)
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

func (r *MySQLRepository) listProductVariants(ctx context.Context, productID ProductID) ([]*ProductVariant, error) {
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

	var variants []*ProductVariant

	for rows.Next() {
		var variant ProductVariant
		var rawStatus string
		var priceMinor int64
		var currencyCode string

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

		status, err := ParseToProductVariantStatus(rawStatus)
		if err != nil {
			return nil, fmt.Errorf("parse product variant status for variant %q: %w", variant.VariantID, err)
		}

		variant.Status = status
		variant.Price = Money{
			AmountMinor:  priceMinor,
			CurrencyCode: currencyCode,
		}

		variants = append(variants, &variant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product variant rows: %w", err)
	}

	return variants, nil
}

func (r *MySQLRepository) listProductAttributeValues(ctx context.Context, productID ProductID) ([]*ProductAttributeValue, error) {
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
        updated_at
     FROM product_attribute_values
     WHERE product_id = ?
     ORDER BY attribute_id ASC, variant_id ASC, product_attribute_value_id ASC 	
  `

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product attribute values: %w", err)
	}
	defer rows.Close()

	var values []*ProductAttributeValue

	for rows.Next() {
		var value ProductAttributeValue

		var variantID sql.NullString
		var valueString sql.NullString
		var valueNumber sql.NullString
		var valueBoolean sql.NullBool
		var valueJSON []byte
		var unit sql.NullString

		if err := rows.Scan(
			&value.ProductAttributeValueID,
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
			return nil, fmt.Errorf("scan product attribute value row: %w", err)
		}

		if variantID.Valid {
			vid := VariantID(variantID.String)
			value.VariantID = &vid
		}

		if valueString.Valid {
			value.ValueString = &valueString.String
		}

		if valueNumber.Valid {
			value.ValueNumber = &valueNumber.String
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

		values = append(values, &value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product attribute value rows: %w", err)
	}

	return values, nil
}

func (r *MySQLRepository) listProductImages(ctx context.Context, productID ProductID) ([]*ProductImage, error) {
	const query = `
	   SELECT
	     image_id,
             product_id,
             url,
             alt_text,
             display_order
           FROM product_images
           WHERE product_id = ?
           ORDER BY display_order ASC, image_id ASC   
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product images: %w", err)
	}
	defer rows.Close()

	var images []*ProductImage

	for rows.Next() {
		var image ProductImage

		if err := rows.Scan(
			&image.ImageID,
			&image.ProductID,
			&image.Url,
			&image.AltText,
			&image.DisplayOrder,
		); err != nil {
			return nil, fmt.Errorf("scan product image row: %w", err)
		}

		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product image rows: %w", err)
	}

	return images, nil
}
