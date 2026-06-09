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

type attributeData struct {
	attributes  []*catalog.ProductAttributeValue
	definitions []catalog.ProductAttributeDefinition
}

type hydratedProductAttributeValue struct {
	AttributeID  catalog.AttributeID
	Code         string
	DisplayName  string
	DataType     catalog.ProductAttributeDataType
	ValueString  string
	ValueNumber  string
	ValueBoolean *bool
	ValueJSON    []byte
	Unit         string
	ValueOptions []*catalog.ProductAttributeOption
}

type variantData struct {
	variants    []*catalog.ProductVariant
	definitions []catalog.ProductAttributeDefinition
}

type hydratedProductVariant struct {
	VariantID   catalog.VariantID
	ProductID   catalog.ProductID
	Sku         catalog.Sku
	VariantName string
	Status      catalog.ProductVariantStatus
	Price       catalog.Money
	Attributes  []*hydratedProductAttributeValue
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// statusToProto maps anything type implementing the Catalog.Status interface
// to catalogv1.ProductStatus
func catalogStatusToProto[T catalog.LifecycleStatus, E catalogv1.ProductStatus, err error](t T) (E, error) {
	status, ok := any(t).(string)
	if ok {
		switch status {
		case "draft":
			return E(catalogv1.ProductStatus_PRODUCT_STATUS_DRAFT), nil
		case "active":
			return E(catalogv1.ProductStatus_PRODUCT_STATUS_ACTIVE), nil
		case "inactive":
			return E(catalogv1.ProductStatus_PRODUCT_STATUS_INACTIVE), nil
		case "archived":
			return E(catalogv1.ProductStatus_PRODUCT_STATUS_ARCHIVED), nil
		}
	}
	return E(catalogv1.ProductStatus_PRODUCT_STATUS_UNSPECIFIED), fmt.Errorf("unknown catalog lifecycle status: %q", status)

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

func productImageToProto(img *catalog.ProductImage) (*catalogv1.ProductImage, error) {
	if img == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to ProductImage proto")
	}

	return &catalogv1.ProductImage{
		ImageId:      string(img.ImageID),
		ProductId:    string(img.ProductID),
		Url:          img.Url,
		AltText:      img.AltText,
		DisplayOrder: int32(img.DisplayOrder),
	}, nil
}

// Maps a catalog.Category to a *catalogv1.Category
func categoryToProto(category *catalog.Category) (*catalogv1.Category, error) {
	if category == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to Category proto")
	}

	status, err := catalogStatusToProto[catalog.CategoryStatus, catalogv1.ProductStatus, error](category.Status)
	if err != nil {
		return nil, err
	}

	response := &catalogv1.Category{
		CategoryId:  string(category.CategoryID),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Status:      status,
		CreatedAt:   timestamppb.New(category.CreatedAt),
		UpdatedAt:   timestamppb.New(category.UpdatedAt),
	}

	if category.ParentCategoryID != nil {
		response.ParentCategoryId = string(category.CategoryID)
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

	status, err := catalogStatusToProto[catalog.ProductStatus, catalogv1.ProductStatus, error](product.Status)
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

func createProductAttributes(ad *attributeData) ([]*hydratedProductAttributeValue, error) {
	if ad == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to []*hydratedProductAttributeValue")
	}

	out := make([]*hydratedProductAttributeValue, 0, len(ad.attributes))

	for _, attribute := range ad.attributes {
		id := attribute.AttributeID
		var hpav hydratedProductAttributeValue
		hpav.AttributeID = id
		hpav.Unit = attribute.Unit
		hpav.ValueString = attribute.ValueString
		hpav.ValueNumber = attribute.ValueNumber
		hpav.ValueBoolean = attribute.ValueBoolean
		hpav.ValueJSON = attribute.ValueJSON

		for _, definition := range ad.definitions {
			if id == definition.AttributeID {
				hpav.Code = definition.Code
				hpav.DisplayName = definition.DisplayName
				hpav.DataType = definition.DataType
				hpav.ValueOptions = definition.Options
			}
			break
		}
		out = append(out, &hpav)
	}

	return out, nil
}

func createProductVariant(data *variantData) ([]*hydratedProductVariant, error) {

	if data == nil {
		return nil, fmt.Errorf("unable to map nil pointer to []*hydratedProductVariant")
	}

	out := make([]*hydratedProductVariant, 0, len(data.variants))

	for _, variant := range data.variants {
		var hpv hydratedProductVariant

		hydratedAttributes, err := createProductAttributes(&attributeData{variant.Attributes, data.definitions})
		if err != nil {
			return nil, fmt.Errorf("create *[]hydratedProductAttribute: %w", err)
		}

		hpv.VariantID = variant.VariantID
		hpv.ProductID = variant.ProductID
		hpv.Sku = variant.Sku
		hpv.VariantName = variant.VariantName
		hpv.Status = variant.Status
		hpv.Price = variant.Price
		hpv.CreatedAt = variant.CreatedAt
		hpv.UpdatedAt = variant.UpdatedAt
		hpv.Attributes = hydratedAttributes

		out = append(out, &hpv)
	}

	return out, nil
}

// Maps a catalog.Product to a *catalogv1.Product
func productToProto(product *catalog.Product, definitions []catalog.ProductAttributeDefinition) (*catalogv1.Product, error) {
	if product == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to Product proto")
	}

	status, err := catalogStatusToProto[catalog.ProductStatus, catalogv1.ProductStatus, error](product.Status)
	if err != nil {
		return nil, err
	}

	data := attributeData{
		product.Attributes,
		definitions,
	}
	attributeSlice, err := createProductAttributes(&data)
	if err != nil {
		return nil, fmt.Errorf("adding attribute definition data to product attribute values :%w", err)
	}
	attributes, err := sliceutil.Map[*hydratedProductAttributeValue, *catalogv1.ProductAttributeValue, error](attributeSlice, productAttributeValueToProto)
	if err != nil {
		return nil, fmt.Errorf("product to protobuf mapping, attributes: %w", err)
	}

	varData := variantData{
		product.Variants,
		definitions,
	}
	variantSlice, err := createProductVariant(&varData)
	if err != nil {
		return nil, fmt.Errorf("adding attribute definition data to product attribute values :%w", err)
	}
	variants, err := sliceutil.Map[*hydratedProductVariant, *catalogv1.ProductVariant, error](variantSlice, productVariantToProto)
	if err != nil {
		return nil, fmt.Errorf("product to protobuf mapping, variants: %w", err)
	}

	images, err := sliceutil.Map[*catalog.ProductImage, *catalogv1.ProductImage, error](product.Images, productImageToProto)
	if err != nil {
		return nil, fmt.Errorf("product to protobuf mapping, images: %w", err)
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
		Slug:       product.Slug, // TODO: regen product.pb.go - updated version includes slug
		Variants:   variants,
		Attributes: attributes,
		Images:     images,
		CreatedAt:  timestamppb.New(product.CreatedAt),
		UpdatedAt:  timestamppb.New(product.UpdatedAt),
	}, nil
}

func productAttributeOptionToProto(p *catalog.ProductAttributeOption) (*catalogv1.ProductAttributeOption, error) {
	if p == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to ProductAttributeOption proto")
	}

	status, err := catalogStatusToProto[catalog.ProductAttributeOptionStatus, catalogv1.ProductStatus, error](p.Status)
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

func productVariantToProto(pv *hydratedProductVariant) (*catalogv1.ProductVariant, error) {
	if pv == nil {
		return nil, fmt.Errorf("unable to convert nil pointer to ProductVariant proto")
	}

	status, err := catalogStatusToProto[catalog.ProductVariantStatus, catalogv1.ProductStatus, error](pv.Status)
	if err != nil {
		return nil, err
	}

	attributes, err := sliceutil.Map[*catalog.ProductAttributeValue, *catalogv1.ProductAttributeValue, error](pv.Attributes, productAttributeValueToProto)
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
			return nil, fmt.Errorf("parsing ValueNumber to float64: %w")
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

	status, err := catalogStatusToProto[catalog.ProductAttributeDefinitionStatus, catalogv1.ProductStatus, error](pad.Status)
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
