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
func (s *Service) GetProduct(ctx context.Context, productID ProductID) (ProductDetails, error) {
	id, err := isValidCatalogID(productID)
	if err != nil {
		return ProductDetails{}, ErrInvalidProductID
	}

	product, err := s.repository.GetProduct(ctx, ProductID(id))
	if err != nil {
		return ProductDetails{}, fmt.Errorf("get product %q: %w", id, err)
	}

	definitions, err := s.repository.ListProductAttributeDefinitions(ctx, ListQuery{
		ID:            string(product.CategoryID),
		FilterOptions: []bool{true, false},
		Limit:         500, // Magic Number - create ListAllProductAttributeDefinitionForCategory(ctx, catgeoryID)([]ProductAttributeDefinition, error)
	})
	if err != nil {
		return ProductDetails{}, fmt.Errorf("list attribute defintions for category %q: %w", product.CategoryID, err)
	}

	details, err := hydrateProduct(product, definitions)
	if err != nil {
		return ProductDetails{}, fmt.Errorf("hydrate product %q: %w", product.ProductID, err)
	}

	return details, nil

}

func hydrateProduct(product Product, definitions []ProductAttributeDefinition) (ProductDetails, error) {
	definitionsByID := make(map[AttributeID]ProductAttributeDefinition, len(definitions))

	for _, definition := range definitions {
		definitionsByID[definition.AttributeID] = definition
	}

	details := ProductDetails{
		ProductID:   product.ProductID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Brand:       product.Brand,
		Status:      product.Status,
		BasePrice:   product.BasePrice,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Images:      product.Images,
	}

	for _, value := range product.Attributes {
		if value.VariantID != nil { // no nil indicates variant product details, ignore to avoid duplication
			continue
		}

		hydrated, err := hydrateProductAttributeValue(value, definitionsByID)
		if err != nil {
			return ProductDetails{}, err
		}

		details.Attributes = append(details.Attributes, hydrated)
	}

	for _, variant := range product.Variants {
		variantDetails := &ProductVariantDetails{
			VariantID:   variant.VariantID,
			ProductID:   variant.ProductID,
			Sku:         variant.Sku,
			VariantName: variant.VariantName,
			Status:      variant.Status,
			Price:       variant.Price,
			CreatedAt:   variant.CreatedAt,
			UpdatedAt:   variant.UpdatedAt,
		}

		for _, value := range product.Attributes {
			if value.VariantID == nil {
				continue
			}

			if variant.VariantID != *value.VariantID {
				continue
			}

			hydrated, err := hydrateProductAttributeValue(value, definitionsByID)
			if err != nil {
				return ProductDetails{}, err
			}

			variantDetails.Attributes = append(variantDetails.Attributes, hydrated)
		}

		details.Variants = append(details.Variants, variantDetails)
	}

	return details, nil
}
func hydrateProductAttributeValue(value *ProductAttributeValue, definitionsByID map[AttributeID]ProductAttributeDefinition) (*ProductAttributeValueDetails, error) {
	if value == nil {
		return nil, fmt.Errorf("nil product attribute value")
	}

	definition, ok := definitionsByID[value.AttributeID]
	if !ok {
		return nil, fmt.Errorf("missing attribute definition %q", value.AttributeID)
	}

	return &ProductAttributeValueDetails{
		ProductAttributeValueID: value.ProductAttributeValueID,
		ProductID:               value.ProductID,
		VariantID:               value.VariantID,
		AttributeID:             value.AttributeID,

		Code:        definition.Code,
		DisplayName: definition.DisplayName,
		DataType:    definition.DataType,
		Options:     definition.Options,

		ValueString:  value.ValueString,
		ValueNumber:  value.ValueNumber,
		ValueBoolean: value.ValueBoolean,
		ValueJSON:    value.ValueJSON,
		Unit:         value.Unit,

		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}, nil
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
