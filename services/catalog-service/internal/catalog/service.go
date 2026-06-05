package catalog

import (
	"context"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"encoding/base64"
)


const (
	defaultPageSize = 20
	maxPageSize = 100
)


// Service contains Catalogue Service business logic.
type Service struct {
	repository Repository
}

// NewService creates a catalogue service.
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

type CatalogQueryResult interface {
	Product |
	Category |
	ProductAttributeDefinition
}

type ListCatalogQueryResult struct {
	Result []CatalogQueryResult
	NextPageToken string
}

type ListQuery struct {
	ID CatalogID
	FilterOptions []bool
	Limit int
	Cursor catalogCursor
}

type catalogCursor struct {
	CreatedAt time.Time  `json:"created_at"`
	ID CatalogID	`json:"id"`
}

// ListProducts returns customer-visible catalogue products.
func (s *Service) ListProducts(ctx context.Context, input ListProductsFilter) (ListCatalogQueryResult, error) {
	pageSize, err := normalisePageSize(input.PageSize)
	if err != nil {
		return ListCatalogQueryResult{}, ErrInvalidPageSize
	}
	
	var cursor catalogCursor
	if strings.Trim(input.PageToken) != "" {
		cursor, err = decodedCursor(input.PageToken)
		if err != nil {
			return ListCatalogQueryResult{}, fmt.Errorf("invalid page token: %w", err)

		}
	}	

	products, err := s.repository.ListProducts(ctx, ListQuery{
		ID: input.CatalogID,
		FilterOptions: []bool{input.IncludeInactive},
		Limit: pageSize + 1,
		Cursor: cursor,
	})
	if err != nil {
		return ListCatalogQueryResult{}, nil
	}

	nextToken := ""
	hasMore := len(products) > pageSize

	if hasMore {
		products = products[:pageSize]
		last := products[len(products) - 1]
		nexToken, err = encodeCatalogCursor(last)
		if err != nil {
			return ListCatalogQueryResult{}, nil
		}
	}

	return ListCatalogQueryResult{
		Result: products,
		NextPageToken: nextToken,
		}, nil
}

// GetProduct returns a single product.
func (s *Service) GetProduct(ctx context.Context, productID string) (Product, error) {
	if strings.Trim(productID) == "" {
		return Product{}, ErrInvalidProductID
	}

	product, err := s.repository.GetProduct(ctx, productID)
	if err != nil {
		return Product{}, ErrProductNotFound
	}

	return product, nil
}

// ListCategories returns customer-visible catalogue categories.
func (s *Service) ListCategories(ctx context.Context, input ListCategoriesFilter) (ListCatalogQueryResult, error) {

	pageSize, err := normalisePageSize(input.PageSize)
	if err != nil {
		return ListCatalogQueryResult{}, ErrInvalidPageSize
	}

	var cursor catalogCursor
	if strings.Trim(input.PageToken) != "" {
		cursor, err = decodedCursor(input.PageToken)
		if err != nil {
			return ListCatalogQueryResult{},  fmt.Errorf("invalid page token: %w", err)
		}
	}

	categories, err := s.repository.ListCategories(ctx, ListQuery{
		ID: input.ParentCategoryID,
		FilterOptions: []bool{input.IncludeInactive},
		Limit: pageSize + 1,
		Cursor: cursor,
	})
	if err != nil {
		return ListCatalogQueryResult{}, nil
	}

	nextToken := ""
	hasMore := len(categories) > pageSize

	if hasMore {
		categories = categories[:pageSize]
		last := categories[len(categories) - 1]
		nexToken, err = encodeCatalogCursor(last)
		if err != nil {
			return ListCatalogQueryResult{}, nil
		}
	}

	return ListCatalogQueryResult{
		Result: categories,
		NextPageToken: nextToken,
		}, nil
}

// ListProductAttributeDefinitions returns catalogue product attribute definitions.
func (s *Service) ListProductAttributeDefinitions(ctx context.Context, 
		  input ListProductAttributeDefinitionsFilter) 
		(ListCatalogQueryResult, error) {
	
	pageSize, err := normalisePageSize(input.PageSize)
	if err != nil {
		return ListCatalogQueryResult{}, ErrInvalidPageSize
	}

	var cursor catalogCursor
	if strings.Trim(input.PageToken) != "" {
		cursor, err = decodedCursor(input.PageToken)
		if err != nil {
			return ListCatalogQueryResult{}, fmt.Errorf("invalid page token: %w", err)
		}
	}

	attributeDefinitions, err := s.repository.ListAttributeDefinitions(ctx, ListQuery{
		ID: input.CategoryID,
		FilterOptions: []bool{input.IncludeInactive, input.IsFilterable},
		Limit: pageSize + 1,
		Cursor: cursor,
	})
	if err != nil {
		return ListCatalogQueryResult{}, nil
	}

	nextToken := ""
	hasMore := len(attributeDefinitions) > pageSize

	if hasMore {
		attributeDefinitions = attributeDefinitions[:pageSize]
		last := attributeDefinitions[len(attributeDefinitions) - 1]
		nexToken, err = encodeCatalogCursor(last)
		if err != nil {
			return ListCatalogQueryResult{}, nil
		}
	}

	return ListCatalogQueryResult{
		Result: attributeDefinitions,
		NextPageToken: nextToken,
		}, nil

}


// Helper functions for PageSize and PageToken

func normalisePageSize(size int) (int, err) {
	if size <= 0 {
		return defaultPageSize, nil
	}

	if size > maxPageSize {
		return 0, ErrInvalidPageSize
	}

	return size, nil
}

func decodeCursor(pageToken string) (catalogCursor, error) {
	data, err := base64.RawURLEncoding.DecodeString(pageToken)
	if err != nil {
		return catalogCursor{}, err
	}

	var token catalogCursor
	if err := json.Unmarshal(data, &token); err != nil {
		return catalogCursor{}, err
	}
	
	if token.ID == "" | token.CreatedAt.IsZero() {
		return catalogCursor{}, errors.New("invalid field in catalog cursor")
	}

	return token, nil
}

func encodeProductCursor(c catalogQueryResult) (string, error) {
	token := catalogCursor{
		CreatedAt: c.CreatedAt,		
	}

	switch catalogQueryType := c.(type) {
	case Product:
		token.ID = catalogQueryType.ProductID
	case Category, ProductAttributeDefiniton:
		token.ID =  catalogQueryType.CategoryID
	default:
		"", fmt.Errorf("unknown catalog domain type: %v", c)
	}	

	data, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}
