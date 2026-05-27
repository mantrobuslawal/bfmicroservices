-- Migration: 000001_create_catalog_schema
-- Direction: down
--
-- Purpose:
-- Drop the initial Catalogue Service schema.
--
-- Warning:
-- This is destructive and intended for local/test environments.

USE bfstore_catalog;

DROP TABLE IF EXISTS catalogue_outbox_events;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS product_attribute_values;
DROP TABLE IF EXISTS product_attribute_options;
DROP TABLE IF EXISTS product_attribute_definitions;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
