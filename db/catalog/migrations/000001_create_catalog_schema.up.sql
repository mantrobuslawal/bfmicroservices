-- Migration: 000001_create_catalog_schema
-- Direction: up
--
-- Purpose:
-- Create the initial Catalogue Service schema.
--
-- Catalogue Service owns:
-- - product identity
-- - categories
-- - product status
-- - product variants
-- - category-scoped attribute definitions
-- - product and variant attribute values
-- - product imagery
--
-- Catalogue Service does not own:
-- - stock levels
-- - basket state
-- - orders
-- - payments
-- - shipping
-- - search indexes

USE bfstore_catalog;

CREATE TABLE categories (
  category_id CHAR(36) NOT NULL,
  parent_category_id CHAR(36) NULL,
  name VARCHAR(160) NOT NULL,
  slug VARCHAR(180) NOT NULL,
  description TEXT NULL,
  status ENUM('draft', 'active', 'inactive', 'archived') NOT NULL DEFAULT 'draft',
  display_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  PRIMARY KEY (category_id),
  UNIQUE KEY uq_categories_slug (slug),
  KEY idx_categories_parent_category_id (parent_category_id),
  KEY idx_categories_status (status),

  CONSTRAINT fk_categories_parent
    FOREIGN KEY (parent_category_id)
    REFERENCES categories (category_id)
    ON DELETE SET NULL
);

CREATE TABLE products (
  product_id CHAR(36) NOT NULL,
  category_id CHAR(36) NOT NULL,
  name VARCHAR(220) NOT NULL,
  slug VARCHAR(240) NOT NULL,
  description TEXT NULL,
  brand VARCHAR(160) NULL,
  status ENUM('draft', 'active', 'inactive', 'archived') NOT NULL DEFAULT 'draft',
  base_price_minor BIGINT NOT NULL,
  currency_code CHAR(3) NOT NULL DEFAULT 'GBP',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  PRIMARY KEY (product_id),
  UNIQUE KEY uq_products_slug (slug),
  KEY idx_products_category_id (category_id),
  KEY idx_products_status (status),
  KEY idx_products_brand (brand),

  CONSTRAINT fk_products_category
    FOREIGN KEY (category_id)
    REFERENCES categories (category_id)
);

CREATE TABLE product_variants (
  variant_id CHAR(36) NOT NULL,
  product_id CHAR(36) NOT NULL,
  sku VARCHAR(120) NOT NULL,
  variant_name VARCHAR(220) NOT NULL,
  status ENUM('draft', 'active', 'inactive', 'archived') NOT NULL DEFAULT 'draft',
  price_minor BIGINT NOT NULL,
  currency_code CHAR(3) NOT NULL DEFAULT 'GBP',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  PRIMARY KEY (variant_id),
  UNIQUE KEY uq_product_variants_sku (sku),
  KEY idx_product_variants_product_id (product_id),
  KEY idx_product_variants_status (status),

  CONSTRAINT fk_product_variants_product
    FOREIGN KEY (product_id)
    REFERENCES products (product_id)
    ON DELETE CASCADE
);

CREATE TABLE product_attribute_definitions (
  attribute_id CHAR(36) NOT NULL,
  category_id CHAR(36) NOT NULL,
  code VARCHAR(120) NOT NULL,
  display_name VARCHAR(160) NOT NULL,
  description TEXT NULL,
  data_type ENUM('string', 'number', 'boolean', 'option', 'multi_option', 'json') NOT NULL,
  unit VARCHAR(40) NULL,
  is_required BOOLEAN NOT NULL DEFAULT FALSE,
  is_filterable BOOLEAN NOT NULL DEFAULT FALSE,
  is_variant_defining BOOLEAN NOT NULL DEFAULT FALSE,
  status ENUM('draft', 'active', 'inactive', 'archived') NOT NULL DEFAULT 'draft',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  PRIMARY KEY (attribute_id),
  UNIQUE KEY uq_attribute_definitions_category_code (category_id, code),
  KEY idx_attribute_definitions_category_id (category_id),
  KEY idx_attribute_definitions_status (status),
  KEY idx_attribute_definitions_filterable (is_filterable),

  CONSTRAINT fk_attribute_definitions_category
    FOREIGN KEY (category_id)
    REFERENCES categories (category_id)
    ON DELETE CASCADE
);

CREATE TABLE product_attribute_options (
  option_id CHAR(36) NOT NULL,
  attribute_id CHAR(36) NOT NULL,
  value VARCHAR(160) NOT NULL,
  display_name VARCHAR(160) NOT NULL,
  display_order INT NOT NULL DEFAULT 0,
  status ENUM('draft', 'active', 'inactive', 'archived') NOT NULL DEFAULT 'active',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  PRIMARY KEY (option_id),
  UNIQUE KEY uq_attribute_options_attribute_value (attribute_id, value),
  KEY idx_attribute_options_attribute_id (attribute_id),
  KEY idx_attribute_options_status (status),

  CONSTRAINT fk_attribute_options_attribute
    FOREIGN KEY (attribute_id)
    REFERENCES product_attribute_definitions (attribute_id)
    ON DELETE CASCADE
);

CREATE TABLE product_attribute_values (
  product_attribute_value_id CHAR(36) NOT NULL,
  product_id CHAR(36) NOT NULL,
  variant_id CHAR(36) NULL,
  attribute_id CHAR(36) NOT NULL,
  value_string VARCHAR(500) NULL,
  value_number DECIMAL(18, 4) NULL,
  value_boolean BOOLEAN NULL,
  value_json JSON NULL,
  unit VARCHAR(40) NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  PRIMARY KEY (product_attribute_value_id),
  UNIQUE KEY uq_product_attribute_values_product_variant_attribute (product_id, variant_id, attribute_id),
  KEY idx_product_attribute_values_product_id (product_id),
  KEY idx_product_attribute_values_variant_id (variant_id),
  KEY idx_product_attribute_values_attribute_id (attribute_id),

  CONSTRAINT fk_product_attribute_values_product
    FOREIGN KEY (product_id)
    REFERENCES products (product_id)
    ON DELETE CASCADE,

  CONSTRAINT fk_product_attribute_values_variant
    FOREIGN KEY (variant_id)
    REFERENCES product_variants (variant_id)
    ON DELETE CASCADE,

  CONSTRAINT fk_product_attribute_values_attribute
    FOREIGN KEY (attribute_id)
    REFERENCES product_attribute_definitions (attribute_id)
);

CREATE TABLE product_images (
  image_id CHAR(36) NOT NULL,
  product_id CHAR(36) NOT NULL,
  url VARCHAR(1000) NOT NULL,
  alt_text VARCHAR(500) NOT NULL,
  display_order INT NOT NULL DEFAULT 0,
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  PRIMARY KEY (image_id),
  KEY idx_product_images_product_id (product_id),
  KEY idx_product_images_display_order (display_order),

  CONSTRAINT fk_product_images_product
    FOREIGN KEY (product_id)
    REFERENCES products (product_id)
    ON DELETE CASCADE
);

CREATE TABLE catalogue_outbox_events (
  outbox_event_id CHAR(36) NOT NULL,
  aggregate_type VARCHAR(120) NOT NULL,
  aggregate_id CHAR(36) NOT NULL,
  event_type VARCHAR(160) NOT NULL,
  event_version VARCHAR(20) NOT NULL,
  payload BLOB NOT NULL,
  content_type VARCHAR(120) NOT NULL DEFAULT 'application/x-protobuf',
  headers JSON NULL,
  status ENUM('pending', 'published', 'failed') NOT NULL DEFAULT 'pending',
  publish_attempts INT NOT NULL DEFAULT 0,
  last_error TEXT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  published_at TIMESTAMP(6) NULL,

  PRIMARY KEY (outbox_event_id),
  KEY idx_catalogue_outbox_events_status_created_at (status, created_at),
  KEY idx_catalogue_outbox_events_aggregate (aggregate_type, aggregate_id),
  KEY idx_catalogue_outbox_events_event_type (event_type)
);
