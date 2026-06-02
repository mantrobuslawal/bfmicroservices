package catalog

import (
	"context"
	"fmt"
)

// Service contains Catalogue Service business logic.
type Service struct {
	repository Repository
}

// NewService creates a catalogue service.
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// ListProducts returns customer-visible catalogue products.
func (s *Service) ListProducts(ctx context.Context, filter ListProductsFilter) ([]Product, error) {
	products, err := s.repository.ListProducts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return products, nil
}

// GetProduct returns a single product.
func (s *Service) GetProduct(ctx context.Context, productID string) (Product, error) {
	if productID == "" {
		return Product{}, fmt.Errorf("product_id is required")
	}

	product, err := s.repository.GetProduct(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("get product: %w", err)
	}

	return product, nil
}

// ListCategories returns customer-visible catalogue categories.
func (s *Service) ListCategories(ctx context.Context, filter ListCategoriesFilter) ([]Category, error) {
	categories, err := s.repository.ListCategories(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return categories, nil
}

// ListProductAttributeDefinitions returns catalogue product attribute definitions.
func (s *Service) ListProductAttributeDefinitions(ctx context.Context, 
						  filter ListProductAttributeDefinitionsFilter) 
						  ([]AttributeDefinition, error) {
	attributeDefinitions, err := s.repository.ListProductProductAttributeDefinitions(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list product attribute definitions: %w", err)
	}

	return attributeDefinitions, nil
}
