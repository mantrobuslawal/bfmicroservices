package grpc

import (
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/sliceutil"
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


// Maps a catalog.Product to a *catalogv1.Product
 func productToProto(product catalog.Product) *catalogv1.Product {
 	return &catalogv1.Product{
 		ProductId:   product.ProductID,
 		CategoryId:  product.CategoryID,
 		Name:        product.Name,
 		Description: product.Description,
 		Brand:       product.Brand,
 		Status:      catalogStatusToProto(product.Status),
 		BasePrice: &commonv1.Money{
 			AmountMinor:  product.BasePrice.AmountMinor,
 			CurrencyCode: product.BasePrice.CurrencyCode,
 		},
		Slug: product.Slug,
		Variants: sliceUtil.Map(product.ProductVariants, mapProductVarinatsToProto),
		Attributes:   sliceUtil.Map(product.ProductAttributeValue, ProductAttributeValueToProto),
		Images:   sliceUtil.Map(product.Images, ProductImageToProto),
 		CreatedAt: timestamppb.New(product.CreatedAt),
 		UpdatedAt: timestamppb.New(product.UpdatedAt),
 	}
 }



// Maps a catalog.Category to a *catalogv1.Category
 func categoryToProto(category catalog.Category) *catalogv1.Category {
 	response := &catalogv1.Category{
 		CategoryId:  category.CategoryID,
 		Name:        category.Name,
 		Slug:        category.Slug,
 		Description: category.Description,
 		Status:      catalogStatusToProto(category.Status),
 		CreatedAt:   timestamppb.New(category.CreatedAt),
 		UpdatedAt:   timestamppb.New(category.UpdatedAt),
 	}

 	if category.ParentCategoryID != nil {
 		response.ParentCategoryId = *category.ParentCategoryID
 	}

 	return response
 }


 
// statusToProto maps anything type implementing the Catalog.Status interface
// to catalogv1.ProductStatus
func catalogStatusToProto(status catalog.LifecycleStatus) catalogv1.ProductStatus {
	switch status {
	case "draft":
		return catalogv1.ProductStatus(1)
	case "active":
		return catalogv1.ProductStatus(2)
	case "inactive":
		return catalogv1.ProductStatus(3)
	case "archived":
		return catalogv1.ProductStatus(4)
	}
}


func productVariantToProto(pv catalog.ProductVariant) *catalogv1.ProductVariant{
	return &catalogv1.ProductVariant{
		VariantId: pv.VariantID,
		ProductId: pv.ProductID,
		Sku: pv.Sku,
		Status: catalogStatusToProto(pv.Status),
		Price: &commonv1.Money{
 			AmountMinor:  pv.Price.AmountMinor,
 			CurrencyCode: pv.Price.CurrencyCode,
 		},
		Attributes: sliceUtil.Map(pv.ProductAttributeValue, productAttributeValueToProto),
 		CreatedAt:   timestamppb.New(pv.CreatedAt),
 		UpdatedAt:   timestamppb.New(pv.UpdatedAt),
	}
}


func productAttributeValueToProto(pav catalog.ProductAttributeValue) *catalogv1.ProductAttributeValue {
	var valueOptions []string
	for _, valueOption := range valueOptions {
		valueOptions = append(valueOptions, valueOption)
	}
	return &catalogv1.ProductAttributeValue{
		AttributeId: pav.AttributeID,
		Code: pav.Code,
		DisplayName: pv.DisplayName,
		DataType: productAttributeDataTypeToProto(pv.DataType) ,
		ValueOptions: valueOptions,
		Unit: pav.Unit,
	}
}


func productImageToProto(img catalog.ProductImage) *catalogv1.ProductImage {
	return &catalogv1.ProductImage{
		ImageId: img.ImageID,
		ProductId: img.ProductID,
		Url: img.Url,
		AltText: img.AltText,
		DisplayOrder: img.DisplayOrder,
	}
}


func productAttributeDefinitionToProto(pad catalog.ProductAttributeDefinition) *catalogv1.ProductAttributeDefinition {
	return &catalogv1.ProductAttributeDefinition{
		AttributeId: pad.AttributeID,
		CategoryId: pad.CategoryID,
		Code: pad.Code,
		DisplayName: pad.DisplayName,
		Description: pad.Description,
		DataType: productAttributeDataTypeToProto(pad.DataType),
		Unit: pad.Unit,
		IsRequired: pad.IsRequired,
		IsFilterable: pad.IsFilterable,
		IsVarianDefiniting: pad.IsVariantDefining,
		Options: sliceUtil.Map(pad.ProductAttributeOption, productAttributeOptionToProto)
		Status: CatalogStatusToProto(pad.Status),
	}
}


func productAttributeDataTypeToProto(datatype catalog.ProductAttributeDataType) catalogv1.ProductAttributeDataType {
	switch datatype {
	case "string":
		return catalogv1.ProductAttributeDataType(1)
	case "number":
		return catalogv1.ProductAttributeDataType(2)
	case "boolean":
		return catalogv1.ProductAttributeDataType(3)
	case "option":
		return catalogv1.ProductAttributeDataType(4)
	case "multi_option":
		return catalogv1.ProductAttributeDataType(5)
	case "json":
		return catalogv1.ProductAttributeDataType(6)
	}
}

func productAttributeOptionToProto(p catalog.ProductAttributeOption) *catalogv1.ProducdAttributeOption {
	return &catalogv1.ProductAttributeOption{
		OptionId: p.OptionID,
		Value: p.Value,
		DisplayName: p.DisplayName,
		DisplayOrder: p.DisplayOrder,
		Status: catalogStatusToProto(p.Status)
	}
}


