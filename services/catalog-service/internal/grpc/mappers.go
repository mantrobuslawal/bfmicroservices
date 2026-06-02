package grpc

import (
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	commonv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/common/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// This file intentionally documents the mapper boundary between internal
// domain models and generated Protobuf response messages.
//
// The mapper layer is important because it prevents transport-specific
// generated code from leaking throughout the service internals.
//

 func productStatusToProto(status string) catalogv1.ProductStatus {
        var productStatus int32
	switch status {
 	case "draft" :
             productStatus = 1
        case "active":
              productStatus = 2
        case "inactive":
              productStatus = 3
        case "archived":
              productStatus = 4
	default:
             productStatus = 0
	   
      }
      return catalogv1.ProductStatus(productStatus)
 }


 func productToProto(product catalog.Product) *catalogv1.Product {
 	return &catalogv1.Product{
 		ProductId:   product.ProductID,
 		CategoryId:  product.CategoryID,
 		Name:        product.Name,
 		Description: product.Description,
 		Brand:       product.Brand,
 		Status:      productStatusToProto(product.Status),
 		BasePrice: &commonv1.Money{
 			AmountMinor:  product.BasePrice.AmountMinor,
 			CurrencyCode: product.BasePrice.CurrencyCode,
 		},
 		CreatedAt: timestamppb.New(product.CreatedAt),
 		UpdatedAt: timestamppb.New(product.UpdatedAt),
 	}
 }


 func categoryToProto(category catalog.Category) *catalogv1.Category {
 	response := &catalogv1.Category{
 		CategoryId:  category.CategoryID,
 		Name:        category.Name,
 		Slug:        category.Slug,
 		Description: category.Description,
 		Status:      productStatusToProto(category.Status),
 		CreatedAt:   timestamppb.New(category.CreatedAt),
 		UpdatedAt:   timestamppb.New(category.UpdatedAt),
 	}

 	if category.ParentCategoryID != nil {
 		response.ParentCategoryId = *category.ParentCategoryID
 	}

 	return response
 }


 func productAttributeDefToProto(prodAttrDef catalog.productAttributeDefinition) *catalogv1.ProductAttributeDefinition {
      opts []*catalogv1.ProdAttributeOption
      for _, opt := range proAttrDef.Options {
		opts = append(opts, prodOptionToProto(opt))
	}
      return &catalogv1.ProductAttributeDefinition{
		AttributeId: prodAttrDef.AttributeID,
		CategoryId:  prodAttrDef.CategoryID,
                Code:        prodAttrDef.Code,
                DisplayName: prodAttrDef.DisplayName,
                Description: prodAttrDef.Description,
                DataType:    ,
                Unit: ,
                IsRequired: ,
                IsFilterable: ,
                IsVarianrDefining: ,
                Options:      opts,
                Status: ,   
	}
 }
