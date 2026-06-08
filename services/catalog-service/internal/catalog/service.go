package catalog

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Service contains Catalog Service business logic.
type Service struct{ repository Repository }

// NewService creates a catalog service.
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// CatalogQueryResult represents domain object returned by a catalog service query.
//
// Represents products, product categories and product attribute definiton objects.
type CatalogQueryResult interface {
	Product |
		Category |
		ProductAttributeDefinition
}

// ListCatalogQueryResult represents the combination of a collection of catalog objects
// and next page token.
type ListResult[T CatalogQueryResult] struct {
	Result        []T
	NextPageToken string
}

// ListQuery represents the collection of catalog object id (i.e. product id, category id etc),
// search filter options, max page size and cursor for pagination of catalog results.
type ListQuery struct {
	ID            string
	FilterOptions []bool
	Limit         int
	Cursor        *catalogCursor
}

// catalogCursor represents object encoded to page token string for result pagination.
type catalogCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

// ListProducts returns customer-visible catalogue products.
func (s *Service) ListProducts(ctx context.Context, input ListProductsFilter) (ListResult[Product], error) {
	pageSize, err := normalisePageSize(input.PageSize)
	if err != nil {
		return ListResult[Product]{}, ErrInvalidPageSize
	}

	var cursor *catalogCursor
	if strings.TrimSpace(input.PageToken) != "" {
		cursor, err = decodeCursor(input.PageToken)
		if err != nil {
			return ListResult[Product]{}, fmt.Errorf("invalid page token: %w", err)

		}
	}

	var id string

	id, err = isValidCatalogID(input.CategoryID)
	if err != nil {
		return ListResult[Product]{}, ErrInvalidCategoryID
	}

	products, err := s.repository.ListProducts(ctx, ListQuery{
		ID:            id,
		FilterOptions: []bool{input.IncludeInactive},
		Limit:         pageSize + 1,
		Cursor:        cursor,
	})
	if err != nil {
		return ListResult[Product]{}, fmt.Errorf("list products category id:%q :%w", input.CategoryID, err)
	}

	nextToken := ""
	hasMore := len(products) > pageSize

	if hasMore {
		products = products[:pageSize]
		last := products[len(products)-1]
		nextToken, err = encodeCursor(last)
		if err != nil {
			return ListResult[Product]{}, fmt.Errorf("encode page token: %w", err)
		}
	}

	return ListResult[Product]{
		Result:        products,
		NextPageToken: nextToken,
	}, nil
}

// GetProduct returns a single product, it's variants and attributes.
func (s *Service) GetProduct(ctx context.Context, productID ProductID) (Product, error) {
	id, err := isValidCatalogID(productID)
	if err != nil {
		return Product{}, ErrInvalidProductID
	}

	product, err := s.repository.GetProduct(ctx, ProductID(id))
	if err != nil {
		return Product{}, fmt.Errorf("get product %q: %w", id, err)
	}

	return product, nil
}

// ListCategories returns customer-visible catalogue categories.
func (s *Service) ListCategories(ctx context.Context, input ListCategoriesFilter) (ListResult[Category], error) {

	pageSize, err := normalisePageSize(input.PageSize)
	if err != nil {
		return ListResult[Category]{}, ErrInvalidPageSize
	}

	var cursor *catalogCursor
	if strings.TrimSpace(input.PageToken) != "" {
		cursor, err = decodeCursor(input.PageToken)
		if err != nil {
			return ListResult[Category]{}, fmt.Errorf("invalid page token: %w", err)
		}
	}

	var id string
	id, err = isValidCatalogID(input.ParentCategoryID)
	if err != nil {
		return ListResult[Category]{}, ErrInvalidCategoryID
	}

	categories, err := s.repository.ListCategories(ctx, ListQuery{
		ID:            id,
		FilterOptions: []bool{input.IncludeInactive},
		Limit:         pageSize + 1,
		Cursor:        cursor,
	})
	if err != nil {
		return ListResult[Category]{}, fmt.Errorf("list categories: %w", err)
	}

	nextToken := ""
	hasMore := len(categories) > pageSize

	if hasMore {
		categories = categories[:pageSize]
		last := categories[len(categories)-1]
		nextToken, err = encodeCursor(last)
		if err != nil {
			return ListResult[Category]{}, fmt.Errorf("create next page token: %w", err)
		}
	}

	return ListResult[Category]{
		Result:        categories,
		NextPageToken: nextToken,
	}, nil
}

// ListProductAttributeDefinitions returns catalogue product attribute definitions.
func (s *Service) ListProductAttributeDefinitions(ctx context.Context, input ListProductAttributeDefinitionsFilter) (ListResult[ProductAttributeDefinition], error) {

	pageSize, err := normalisePageSize(input.PageSize)
	if err != nil {
		return ListResult[ProductAttributeDefinition]{}, ErrInvalidPageSize
	}

	var cursor *catalogCursor
	if strings.TrimSpace(input.PageToken) != "" {
		cursor, err = decodeCursor(input.PageToken)
		if err != nil {
			return ListResult[ProductAttributeDefinition]{}, fmt.Errorf("invalid page token: %w", err)
		}
	}

	var id string
	id, err = isValidCatalogID(input.CategoryID)
	if err != nil {
		return ListResult[ProductAttributeDefinition]{}, ErrInvalidCategoryID
	}
	attributeDefinitions, err := s.repository.ListProductAttributeDefinitions(ctx, ListQuery{
		ID:            id,
		FilterOptions: []bool{input.IncludeInactive, input.IsFilterable},
		Limit:         pageSize + 1,
		Cursor:        cursor,
	})
	if err != nil {
		return ListResult[ProductAttributeDefinition]{}, fmt.Errorf("list product attribute definitions: %w", err)
	}

	nextToken := ""
	hasMore := len(attributeDefinitions) > pageSize

	if hasMore {
		attributeDefinitions = attributeDefinitions[:pageSize]
		last := attributeDefinitions[len(attributeDefinitions)-1]
		nextToken, err = encodeCursor(last)
		if err != nil {
			return ListResult[ProductAttributeDefinition]{}, fmt.Errorf("create next page token: %w", err)
		}
	}

	return ListResult[ProductAttributeDefinition]{
		Result:        attributeDefinitions,
		NextPageToken: nextToken,
	}, nil

}

// Helper functions for PageSize and PageToken

func normalisePageSize(size int) (int, error) {
	if size <= 0 {
		return defaultPageSize, nil
	}

	if size > maxPageSize {
		return 0, ErrInvalidPageSize
	}

	return size, nil
}

func decodeCursor(pageToken string) (*catalogCursor, error) {
	data, err := base64.RawURLEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, err
	}

	var token catalogCursor
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	if token.ID == "" {
		return nil, errors.New("invalid field in catalog cursor")
	}

	if token.CreatedAt.IsZero() {
		return nil, errors.New("invalid field in catalog cursor")
	}

	return &token, nil
}

func encodeCursor[T CatalogQueryResult](c T) (string, error) {
	var token catalogCursor
	product, ok := any(c).(Product)
	if ok {
		token.CreatedAt = product.CreatedAt
		token.ID = string(product.ProductID)
	}

	category, ok := any(c).(Category)
	if ok {
		token.CreatedAt = category.CreatedAt
		token.ID = string(category.CategoryID)
	}

	pad, ok := any(c).(ProductAttributeDefinition)
	if ok {
		token.CreatedAt = pad.CreatedAt
		token.ID = string(pad.AttributeID)
	}

	data, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("marshal struct to json: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

// TODO: Implement proper validation
func isValidCatalogID[T CatalogID](id T) (string, error) {
	return string(id), nil
}
