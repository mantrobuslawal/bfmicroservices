package grpcadapter

import (
	"fmt"
	"strconv"
	"time"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	commonv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/common/v1"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/sliceutil"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// This file intentionally documents the mapper boundary between internal
// domain models and generated Protobuf response messages.
//
// The mapper layer is important because it prevents transport-specific
// generated code from leaking throughout the service internals.
//

// statusToProto maps anything type implementing the Catalog.Status interface
// to catalogv1.ProductStatus
func catalogStatusToProto[T catalog.LifecycleStatus](status T) (catalogv1.ProductStatus, error) {
	switch string(status) {
	case "draft":
		return catalogv1.ProductStatus_PRODUCT_STATUS_DRAFT, nil
	case "active":
		return catalogv1.ProductStatus_PRODUCT_STATUS_ACTIVE, nil
	case "inactive":
		return catalogv1.ProductStatus_PRODUCT_STATUS_INACTIVE, nil
	case "archived":
		return catalogv1.ProductStatus_PRODUCT_STATUS_ARCHIVED, nil
	default:
		return catalogv1.ProductStatus_PRODUCT_STATUS_UNSPECIFIED, fmt.Errorf("unknown catalog lifecycle status: %q", string(status))
	}
}

func productAttributeDataTypeToProto(datatype catalog.ProductAttributeDataType) (catalogv1.ProductAttributeDataType, error) {
	switch datatype {
	case "string":
		return catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_STRING, nil
	case "number":
		return catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_NUMBER, nil
	case "boolean":
		return catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_BOOLEAN, nil
	case "option":
		return catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_OPTION, nil
	case "multi_option":
		return catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_MULTI_OPTION, nil
	case "json":
		return catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_JSON, nil
	default:
		return catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_UNSPECIFIED, fmt.Errorf("unknown product attribute data type: %q", datatype)
	}
}

func productImageToProto(img *catalog.ProductImage) *catalogv1.ProductImage {
	if img == nil {
		return nil
	}

	return &catalogv1.ProductImage{
		ImageId:      string(img.ImageID),
		ProductId:    string(img.ProductID),
		Url:          img.Url,
		AltText:      img.AltText,
		DisplayOrder: int32(img.DisplayOrder),
	}
}

// Maps a catalog.Category to a *catalogv1.Category
func categoryToProto(category *catalog.Category) (*catalogv1.Category, error) {
	if category == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to Category proto")
	}

	status, err := catalogStatusToProto(category.Status)
	if err != nil {
		return nil, err
	}

	response := &catalogv1.Category{
		CategoryId:   string(category.CategoryID),
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		Status:       status,
		DisplayOrder: int32(category.DisplayOrder),
		CreatedAt:    timestamppb.New(category.CreatedAt),
		UpdatedAt:    timestamppb.New(category.UpdatedAt),
	}

	if category.ParentCategoryID != nil {
		response.ParentCategoryId = string(*category.ParentCategoryID)
	}

	return response, nil
}

// Maps a catalog.Product to a *catalogv1.Product.
// Used for ListProducts handler, which doesn't fully hydrate
// product fields to avoid providing unneeded product data (including variants, options, etc)
func listProductToProto(product *catalog.Product) (*catalogv1.Product, error) {
	if product == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to Product proto")
	}

	status, err := catalogStatusToProto(product.Status)
	if err != nil {
		return nil, err
	}

	return &catalogv1.Product{
		ProductId:   string(product.ProductID),
		CategoryId:  string(product.CategoryID),
		Name:        product.Name,
		Description: product.Description,
		Brand:       product.Brand,
		Status:      status,
		BasePrice: &commonv1.Money{
			AmountMinor:  product.BasePrice.AmountMinor,
			CurrencyCode: product.BasePrice.CurrencyCode,
		},
		Slug:      product.Slug,
		CreatedAt: timestamppb.New(product.CreatedAt),
		UpdatedAt: timestamppb.New(product.UpdatedAt),
	}, nil
}

// productDetailsToProto provides a fully hyrdated catalog product.
// Returned when calls are made to GetProduct endpoint.
func productDetailsToProto(product catalog.ProductDetails) (*catalogv1.Product, error) {
	status, err := catalogStatusToProto(product.Status)
	if err != nil {
		return nil, err
	}

	attributes, err := sliceutil.Map[*catalog.ProductAttributeValueDetails, *catalogv1.ProductAttributeValue, error](
		product.Attributes,
		productAttributeValueDetailsToProto,
	)
	if err != nil {
		return nil, err
	}

	variants, err := sliceutil.Map[*catalog.ProductVariantDetails, *catalogv1.ProductVariant, error](
		product.Variants,
		productVariantDetailsToProto,
	)
	if err != nil {
		return nil, err
	}

	images := sliceutil.MapNoError[*catalog.ProductImage, *catalogv1.ProductImage](
		product.Images,
		productImageToProto,
	)

	return &catalogv1.Product{
		ProductId:   string(product.ProductID),
		CategoryId:  string(product.CategoryID),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Brand:       product.Brand,
		Status:      status,
		BasePrice:   moneyToProto(product.BasePrice),
		Attributes:  attributes,
		Variants:    variants,
		Images:      images,
		CreatedAt:   timetoProto(product.CreatedAt),
		UpdatedAt:   timetoProto(product.UpdatedAt),
	}, nil

}

func productAttributeValueDetailsToProto(p *catalog.ProductAttributeValueDetails) (*catalogv1.ProductAttributeValue, error) {
	if p == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to proto")
	}

	var proto catalogv1.ProductAttributeValue

	proto.AttributeId = string(p.AttributeID)
	proto.Code = p.Code
	proto.DisplayName = p.DisplayName
	proto.Unit = p.Unit

	datatype, err := productAttributeDataTypeToProto(p.DataType)
	if err != nil {
		return nil, fmt.Errorf("product attribute datatype map to proto: %w", err)
	}

	proto.DataType = datatype

	Options := p.Options
	var options []string
	for _, Option := range Options {
		options = append(options, Option.Value)
	}

	proto.ValueOptions = options

	switch {
	case p.ValueString != "":
		proto.Value = &catalogv1.ProductAttributeValue_ValueString{
			ValueString: p.ValueString,
		}

	case p.ValueNumber != "":
		number, err := strconv.ParseFloat(p.ValueNumber, 64)
		if err != nil {
			return nil, fmt.Errorf("parse product attribute number value %q: %w", p.ValueNumber, err)
		}
		proto.Value = &catalogv1.ProductAttributeValue_ValueNumber{
			ValueNumber: number,
		}

	case p.ValueBoolean != nil:
		proto.Value = &catalogv1.ProductAttributeValue_ValueBoolean{
			ValueBoolean: *p.ValueBoolean,
		}

	case len(p.ValueJSON) > 0:
		proto.Value = &catalogv1.ProductAttributeValue_ValueJson{
			ValueJson: string(p.ValueJSON),
		}

	default:
		return nil, fmt.Errorf("product attribute value %q has no value set", p.AttributeID)

	}

	return &proto, nil
}

func productVariantDetailsToProto(p *catalog.ProductVariantDetails) (*catalogv1.ProductVariant, error) {
	if p == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to proto")
	}

	status, err := catalogStatusToProto(p.Status)
	if err != nil {
		return nil, err
	}

	attributes, err := sliceutil.Map[*catalog.ProductAttributeValueDetails, *catalogv1.ProductAttributeValue, error](p.Attributes, productAttributeValueDetailsToProto)
	if err != nil {
		return nil, fmt.Errorf("product variant to proto mapping, product attribute values: %w", err)
	}

	return &catalogv1.ProductVariant{
		VariantId:   string(p.VariantID),
		VariantName: p.VariantName,
		ProductId:   string(p.ProductID),
		Sku:         string(p.Sku),
		Status:      status,
		Price:       moneyToProto(p.Price),
		Attributes:  attributes,
		CreatedAt:   timetoProto(p.CreatedAt),
		UpdatedAt:   timetoProto(p.UpdatedAt),
	}, nil
}

func moneyToProto(money catalog.Money) *commonv1.Money {
	return &commonv1.Money{
		AmountMinor:  money.AmountMinor,
		CurrencyCode: money.CurrencyCode,
	}
}
func timetoProto(time time.Time) *timestamppb.Timestamp {
	return timestamppb.New(time)
}

func productAttributeOptionToProto(p *catalog.ProductAttributeOption) (*catalogv1.ProductAttributeOption, error) {
	if p == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to ProductAttributeOption proto")
	}

	status, err := catalogStatusToProto(p.Status)
	if err != nil {
		return nil, err
	}

	return &catalogv1.ProductAttributeOption{
		OptionId:     string(p.OptionID),
		Value:        p.Value,
		DisplayName:  p.DisplayName,
		DisplayOrder: int32(p.DisplayOrder),
		Status:       status,
	}, nil
}

func productAttributeDefinitionToProto(pad *catalog.ProductAttributeDefinition) (*catalogv1.ProductAttributeDefinition, error) {
	if pad == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to ProductAttributeDefinition proto")
	}

	dt, err := productAttributeDataTypeToProto(pad.DataType)
	if err != nil {
		return nil, fmt.Errorf("map product attribute definition data type to proto: %w", err)
	}

	options, err := sliceutil.Map[*catalog.ProductAttributeOption, *catalogv1.ProductAttributeOption, error](pad.Options, productAttributeOptionToProto)
	if err != nil {
		return nil, fmt.Errorf("map product attribute option to to proto: %w", err)
	}

	status, err := catalogStatusToProto(pad.Status)
	if err != nil {
		return nil, fmt.Errorf("map product attribute definition status to product status proto: %w", err)
	}

	return &catalogv1.ProductAttributeDefinition{
		AttributeId:       string(pad.AttributeID),
		CategoryId:        string(pad.CategoryID),
		Code:              pad.Code,
		DisplayName:       pad.DisplayName,
		Description:       pad.Description,
		DataType:          dt,
		Unit:              pad.Unit,
		IsRequired:        pad.IsRequired,
		IsFilterable:      pad.IsFilterable,
		IsVariantDefining: pad.IsVariantDefining,
		Options:           options,
		Status:            status,
	}, nil
}

/* TODO: NO LONGER REQUIRED - DELETE ONCE ALL TESTS PASS
// Maps a catalog.Product to a *catalogv1.Product
func productToProto(hydratedProduct *catalog.HydratedProduct) (*catalogv1.Product, error) {
	if hydratedProduct == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to Product proto")
	}

	status, err := catalogStatusToProto(hydratedProduct.Status)
	if err != nil {
		return nil, err
	}

	attributes, err := sliceutil.Map[*hydratedProductAttributeValue, *catalogv1.ProductAttributeValue, error](hydratedProduct.Attributes, productAttributeValueToProto)
	if err != nil {
		return nil, fmt.Errorf("product to protobuf mapping, attributes: %w", err)
	}

	variants, err := sliceutil.Map[*hydratedProductVariant, *catalogv1.ProductVariant, error](hydratedProduct.Variants, productVariantToProto)
	if err != nil {
		return nil, fmt.Errorf("product to protobuf mapping, variants: %w", err)
	}

	images, err := sliceutil.Map[*catalog.ProductImage, *catalogv1.ProductImage, error](hydratedProduct.Images, productImageToProto)
	if err != nil {
		return nil, fmt.Errorf("product to protobuf mapping, images: %w", err)
	}

	return &catalogv1.Product{
		ProductId:   string(hydratedProduct.ProductID),
		CategoryId:  string(hydratedProduct.CategoryID),
		Name:        hydratedProduct.Name,
		Description: hydratedProduct.Description,
		Brand:       hydratedProduct.Brand,
		Status:      status,
		BasePrice: &commonv1.Money{
			AmountMinor:  hydratedProduct.BasePrice.AmountMinor,
			CurrencyCode: hydratedProduct.BasePrice.CurrencyCode,
		},
		Slug:       hydratedProduct.Slug, // TODO: regen product.pb.go - updated version includes slug
		Variants:   variants,
		Attributes: attributes,
		Images:     images,
		CreatedAt:  timestamppb.New(hydratedProduct.CreatedAt),
		UpdatedAt:  timestamppb.New(hydratedProduct.UpdatedAt),
	}, nil
}



func productVariantToProto(pv *hydratedProductVariant) (*catalogv1.ProductVariant, error) {
	if pv == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to ProductVariant proto")
	}

	status, err := catalogStatusToProto(pv.Status)
	if err != nil {
		return nil, err
	}

	attributes, err := sliceutil.Map[*hydratedProductAttributeValue, *catalogv1.ProductAttributeValue, error](pv.Attributes, productAttributeValueToProto)
	if err != nil {
		return nil, fmt.Errorf("product variant to proto mapping, product attribute values: %w", err)
	}

	return &catalogv1.ProductVariant{
		VariantId: string(pv.VariantID),
		ProductId: string(pv.ProductID),
		Sku:       string(pv.Sku),
		Status:    status,
		Price: &commonv1.Money{
			AmountMinor:  pv.Price.AmountMinor,
			CurrencyCode: pv.Price.CurrencyCode,
		},
		Attributes: attributes,
		CreatedAt:  timestamppb.New(pv.CreatedAt),
		UpdatedAt:  timestamppb.New(pv.UpdatedAt),
	}, nil
}

func productAttributeValueToProto(hpav *hydratedProductAttributeValue) (*catalogv1.ProductAttributeValue, error) {
	if hpav == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to hydrated ProductAttributeValue proto")
	}

	dt, err := productAttributeDataTypeToProto(hpav.DataType)
	if err != nil {
		return nil, fmt.Errorf("product attribute datatype map to proto: %w", err)
	}

	valueOptions := hpav.ValueOptions
	var options []string

	for _, valueOption := range valueOptions {
		options = append(options, valueOption.Value)
	}

	protoValue := &catalogv1.ProductAttributeValue{
		AttributeId:  string(hpav.AttributeID),
		Code:         hpav.Code,
		DisplayName:  hpav.DisplayName,
		DataType:     dt,
		Unit:         hpav.Unit,
		ValueOptions: options,
	}

	switch {
	case hpav.ValueString != "":
		protoValue.Value = &catalogv1.ProductAttributeValue_ValueString{
			ValueString: hpav.ValueString,
		}

	case hpav.ValueNumber != "":
		number, err := strconv.ParseFloat(hpav.ValueNumber, 64)
		if err != nil {
			return nil, fmt.Errorf("parse product attribute number value %q: %w", hpav.ValueNumber, err)
		}
		protoValue.Value = &catalogv1.ProductAttributeValue_ValueNumber{
			ValueNumber: number,
		}

	case hpav.ValueBoolean != nil:
		protoValue.Value = &catalogv1.ProductAttributeValue_ValueBoolean{
			ValueBoolean: *hpav.ValueBoolean,
		}

	case len(hpav.ValueJSON) > 0:
		protoValue.Value = &catalogv1.ProductAttributeValue_ValueJson{
			ValueJson: string(hpav.ValueJSON),
		}

	default:
		return nil, fmt.Errorf("product attribute value %q has no value set", hpav.AttributeID)

	}

	return protoValue, nil
}
*/
