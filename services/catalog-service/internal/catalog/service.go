package catalog

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
type ListCatalogQueryResult struct {
	Result        []CatalogQueryResult
	NextPageToken string
}

// ListQuery represents the collection of catalog object id (i.e. product id, category id etc),
// search filter options, max page size and cursor for pagination of catalog results.
type ListQuery struct {
	ID            string
	FilterOptions []bool
	Limit         int
	Cursor        catalogCursor
}

// catalogCursor represents object encoded to page token string for result pagination.
type catalogCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
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

	var id string

	id, err = isValidCatalogID(input.CatalogID)
	if err != nil {
		return ListCatalogQueryResult{}, ErrInvalidCategoryID
	}

	products, err := s.repository.ListProducts(ctx, ListQuery{
		ID:            id,
		FilterOptions: []bool{input.IncludeInactive},
		Limit:         pageSize + 1,
		Cursor:        cursor,
	})
	if err != nil {
		return ListCatalogQueryResult{}, fmt.Errorf("list products category id:%q :%w", input.CatalogID, err)
	}

	nextToken := ""
	hasMore := len(products) > pageSize

	if hasMore {
		products = products[:pageSize]
		last := products[len(products)-1]
		nexToken, err = encodeCatalogCursor(last)
		if err != nil {
			return ListCatalogQueryResult{}, fmt.Errorf("encode page token: %w", err)
		}
	}

	return ListCatalogQueryResult{
		Result:        products,
		NextPageToken: nextToken,
	}, nil
}

// GetProduct returns a single product, it's variants and attributes.
func (s *Service) GetProduct(ctx context.Context, productID CatalogID) (Product, error) {
	id, err := isValidCatalogID(productID)
	if err != nil {
		return Product{}, ErrInvalidProductID
	}

	product, err := s.repository.GetProduct(ctx, id)
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
			return ListCatalogQueryResult{}, fmt.Errorf("invalid page token: %w", err)
		}
	}

	var id string
	id, err = isValidCatalogID(input.ParentCategoryID)
	if err != nil {
		return ListCatalogQueryResult{}, ErrInvalidCategoryID
	}

	categories, err := s.repository.ListCategories(ctx, ListQuery{
		ID:            id,
		FilterOptions: []bool{input.IncludeInactive},
		Limit:         pageSize + 1,
		Cursor:        cursor,
	})
	if err != nil {
		return ListCatalogQueryResult{}, nil
	}

	nextToken := ""
	hasMore := len(categories) > pageSize

	if hasMore {
		categories = categories[:pageSize]
		last := categories[len(categories)-1]
		nexToken, err = encodeCatalogCursor(last)
		if err != nil {
			return ListCatalogQueryResult{}, nil
		}
	}

	return ListCatalogQueryResult{
		Result:        categories,
		NextPageToken: nextToken,
	}, nil
}

// ListProductAttributeDefinitions returns catalogue product attribute definitions.
func (s *Service) ListProductAttributeDefinitions(ctx context.Context, input ListProductAttributeDefinitionsFilter) (ListCatalogQueryResult, error) {

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

	var id string
	id, err = isValidCatalogID(input.CategoryID)
	if err != nil {
		return ListCatalogQueryResult{}, ErrInvalidCategoryID
	}
	attributeDefinitions, err := s.repository.ListAttributeDefinitions(ctx, ListQuery{
		ID:            id,
		FilterOptions: []bool{input.IncludeInactive, input.IsFilterable},
		Limit:         pageSize + 1,
		Cursor:        cursor,
	})
	if err != nil {
		return ListCatalogQueryResult{}, nil
	}

	nextToken := ""
	hasMore := len(attributeDefinitions) > pageSize

	if hasMore {
		attributeDefinitions = attributeDefinitions[:pageSize]
		last := attributeDefinitions[len(attributeDefinitions)-1]
		nexToken, err = encodeCatalogCursor(last)
		if err != nil {
			return ListCatalogQueryResult{}, nil
		}
	}

	return ListCatalogQueryResult{
		Result:        attributeDefinitions,
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

	if token.ID == ""|token.CreatedAt.IsZero() {
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
		token.ID = catalogQueryType.CategoryID
	default:
		return "", fmt.Errorf("unknown catalog domain type: %v", c)
	}

	data, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

// TODO: Implement proper validation
func isValidCatalogID(id CatalogID) (string, error) {
	return string(id), nil
}
