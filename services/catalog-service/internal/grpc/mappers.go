package grpc

import (
	"time"

	"github.com/acme-ltd/bfstore/services/catalog-service/internal/catalog"

	// TODO:
	// Uncomment after generated Protobuf code is committed.
	//
	// catalogv1 "github.com/acme-ltd/bfstore/gen/go/acme/catalog/v1"
	// commonv1 "github.com/acme-ltd/bfstore/gen/go/acme/common/v1"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

// This file intentionally documents the mapper boundary between internal
// domain models and generated Protobuf response messages.
//
// The mapper layer is important because it prevents transport-specific
// generated code from leaking throughout the service internals.
//
// Once generated Protobuf code exists, replace the placeholder examples below
// with concrete mapping functions.

// Example target implementation:
//
// func productToProto(product catalog.Product) *catalogv1.Product {
// 	return &catalogv1.Product{
// 		ProductId:   product.ProductID,
// 		CategoryId:  product.CategoryID,
// 		Name:        product.Name,
// 		Description: product.Description,
// 		Brand:       product.Brand,
// 		Status:      productStatusToProto(product.Status),
// 		BasePrice: &commonv1.Money{
// 			AmountMinor:  product.BasePrice.AmountMinor,
// 			CurrencyCode: product.BasePrice.CurrencyCode,
// 		},
// 		CreatedAt: timestamppb.New(product.CreatedAt),
// 		UpdatedAt: timestamppb.New(product.UpdatedAt),
// 	}
// }
//
// func categoryToProto(category catalog.Category) *catalogv1.Category {
// 	response := &catalogv1.Category{
// 		CategoryId:  category.CategoryID,
// 		Name:        category.Name,
// 		Slug:        category.Slug,
// 		Description: category.Description,
// 		Status:      productStatusToProto(category.Status),
// 		CreatedAt:   timestamppb.New(category.CreatedAt),
// 		UpdatedAt:   timestamppb.New(category.UpdatedAt),
// 	}
//
// 	if category.ParentCategoryID != nil {
// 		response.ParentCategoryId = *category.ParentCategoryID
// 	}
//
// 	return response
// }

// mapProductForDocumentation keeps this file compile-safe before generated
// Protobuf code is committed.
//
// Remove this once concrete Protobuf mapper functions are implemented.
func mapProductForDocumentation(product catalog.Product) map[string]any {
	return map[string]any{
		"product_id":    product.ProductID,
		"category_id":   product.CategoryID,
		"name":          product.Name,
		"slug":          product.Slug,
		"description":   product.Description,
		"brand":         product.Brand,
		"status":        product.Status,
		"amount_minor":  product.BasePrice.AmountMinor,
		"currency_code": product.BasePrice.CurrencyCode,
		"created_at":    formatTime(product.CreatedAt),
		"updated_at":    formatTime(product.UpdatedAt),
	}
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.UTC().Format(time.RFC3339Nano)
}
