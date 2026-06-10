package grpcadapter

import (
	"testing"
	"time"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
)

func grpcadapterTestTime() time.Time {
	return time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)
}

func grpcadapterBoolPtr(value bool) *bool {
	return &value
}

func TestCatalogStatusToProtoMapsLifecycleStatuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status catalog.ProductStatus
		want   catalogv1.ProductStatus
	}{
		{name: "draft", status: catalog.ProductStatusDraft, want: catalogv1.ProductStatus_PRODUCT_STATUS_DRAFT},
		{name: "active", status: catalog.ProductStatusActive, want: catalogv1.ProductStatus_PRODUCT_STATUS_ACTIVE},
		{name: "inactive", status: catalog.ProductStatusInactive, want: catalogv1.ProductStatus_PRODUCT_STATUS_INACTIVE},
		{name: "archived", status: catalog.ProductStatusArchived, want: catalogv1.ProductStatus_PRODUCT_STATUS_ARCHIVED},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := catalogStatusToProto(tt.status)
			if err != nil {
				t.Fatalf("catalogStatusToProto() error = %v, want nil", err)
			}

			if got != tt.want {
				t.Fatalf("catalogStatusToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCatalogStatusToProtoRejectsUnknownStatus(t *testing.T) {
	t.Parallel()

	_, err := catalogStatusToProto(catalog.ProductStatus("retired"))
	if err == nil {
		t.Fatal("catalogStatusToProto() error = nil, want error")
	}
}

func TestProductAttributeDataTypeToProtoMapsDataTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dataType catalog.ProductAttributeDataType
		want     catalogv1.ProductAttributeDataType
	}{
		{name: "string", dataType: catalog.ProductAttributeDataTypeString, want: catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_STRING},
		{name: "number", dataType: catalog.ProductAttributeDataTypeNumber, want: catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_NUMBER},
		{name: "boolean", dataType: catalog.ProductAttributeDataTypeBool, want: catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_BOOLEAN},
		{name: "option", dataType: catalog.ProductAttributeDataTypeOption, want: catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_OPTION},
		{name: "multi option", dataType: catalog.ProductAttributeDataTypeMultiOption, want: catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_MULTI_OPTION},
		{name: "json", dataType: catalog.ProductAttributeDataTypeJson, want: catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_JSON},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := productAttributeDataTypeToProto(tt.dataType)
			if err != nil {
				t.Fatalf("productAttributeDataTypeToProto() error = %v, want nil", err)
			}

			if got != tt.want {
				t.Fatalf("productAttributeDataTypeToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductAttributeDataTypeToProtoRejectsUnknownDataType(t *testing.T) {
	t.Parallel()

	_, err := productAttributeDataTypeToProto(catalog.ProductAttributeDataType("mystery"))
	if err == nil {
		t.Fatal("productAttributeDataTypeToProto() error = nil, want error")
	}
}

func TestCategoryToProtoMapsParentCategoryID(t *testing.T) {
	t.Parallel()

	parentID := catalog.CategoryID("cat_lighting")
	now := grpcadapterTestTime()

	got, err := categoryToProto(&catalog.Category{
		CategoryID:       catalog.CategoryID("cat_desk_lamps"),
		ParentCategoryID: &parentID,
		Name:             "Desk Lamps",
		Slug:             "desk-lamps",
		Description:      "Desk lamps for dev caves.",
		Status:           catalog.CategoryStatusActive,
		DisplayOrder:     20,
		CreatedAt:        now,
		UpdatedAt:        now.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("categoryToProto() error = %v, want nil", err)
	}

	if got.GetCategoryId() != "cat_desk_lamps" {
		t.Fatalf("CategoryId = %q, want cat_desk_lamps", got.GetCategoryId())
	}

	if got.GetParentCategoryId() != "cat_lighting" {
		t.Fatalf("ParentCategoryId = %q, want cat_lighting", got.GetParentCategoryId())
	}

	if got.GetStatus() != catalogv1.ProductStatus_PRODUCT_STATUS_ACTIVE {
		t.Fatalf("Status = %v, want ACTIVE", got.GetStatus())
	}

	if got.GetDisplayOrder() != 20 {
		t.Fatalf("DisplayOrder = %d, want 20", got.GetDisplayOrder())
	}
}

func TestCategoryToProtoRejectsNilCategory(t *testing.T) {
	t.Parallel()

	_, err := categoryToProto(nil)
	if err == nil {
		t.Fatal("categoryToProto(nil) error = nil, want error")
	}
}

func TestListProductToProtoMapsSummaryProduct(t *testing.T) {
	t.Parallel()

	now := grpcadapterTestTime()

	got, err := listProductToProto(&catalog.Product{
		ProductID:   catalog.ProductID("prod_gopher_lamp"),
		CategoryID:  catalog.CategoryID("cat_lighting"),
		Name:        "Gopher Desk Lamp",
		Slug:        "gopher-desk-lamp",
		Description: "A cheerful lamp.",
		Brand:       "Borough",
		Status:      catalog.ProductStatusActive,
		BasePrice:   catalog.Money{AmountMinor: 4999, CurrencyCode: "GBP"},
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("listProductToProto() error = %v, want nil", err)
	}

	if got.GetProductId() != "prod_gopher_lamp" {
		t.Fatalf("ProductId = %q, want prod_gopher_lamp", got.GetProductId())
	}

	if got.GetSlug() != "gopher-desk-lamp" {
		t.Fatalf("Slug = %q, want gopher-desk-lamp", got.GetSlug())
	}

	if got.GetBasePrice().GetAmountMinor() != 4999 {
		t.Fatalf("AmountMinor = %d, want 4999", got.GetBasePrice().GetAmountMinor())
	}

	if got.GetBasePrice().GetCurrencyCode() != "GBP" {
		t.Fatalf("CurrencyCode = %q, want GBP", got.GetBasePrice().GetCurrencyCode())
	}

	if len(got.GetVariants()) != 0 {
		t.Fatalf("summary product variants len = %d, want 0", len(got.GetVariants()))
	}

	if len(got.GetAttributes()) != 0 {
		t.Fatalf("summary product attributes len = %d, want 0", len(got.GetAttributes()))
	}

	if len(got.GetImages()) != 0 {
		t.Fatalf("summary product images len = %d, want 0", len(got.GetImages()))
	}
}

func TestProductAttributeValueDetailsToProtoMapsStringValue(t *testing.T) {
	t.Parallel()

	got, err := productAttributeValueDetailsToProto(&catalog.ProductAttributeValueDetails{
		AttributeID: catalog.AttributeID("attr_material"),
		Code:        "material",
		DisplayName: "Material",
		DataType:    catalog.ProductAttributeDataTypeString,
		ValueString: "steel",
	})
	if err != nil {
		t.Fatalf("productAttributeValueDetailsToProto() error = %v, want nil", err)
	}

	value, ok := got.GetValue().(*catalogv1.ProductAttributeValue_ValueString)
	if !ok {
		t.Fatalf("Value type = %T, want ValueString", got.GetValue())
	}

	if value.ValueString != "steel" {
		t.Fatalf("ValueString = %q, want steel", value.ValueString)
	}
}

func TestProductAttributeValueDetailsToProtoMapsNumberValue(t *testing.T) {
	t.Parallel()

	got, err := productAttributeValueDetailsToProto(&catalog.ProductAttributeValueDetails{
		AttributeID: catalog.AttributeID("attr_weight"),
		Code:        "weight",
		DisplayName: "Weight",
		DataType:    catalog.ProductAttributeDataTypeNumber,
		ValueNumber: "1.25",
		Unit:        "kg",
	})
	if err != nil {
		t.Fatalf("productAttributeValueDetailsToProto() error = %v, want nil", err)
	}

	value, ok := got.GetValue().(*catalogv1.ProductAttributeValue_ValueNumber)
	if !ok {
		t.Fatalf("Value type = %T, want ValueNumber", got.GetValue())
	}

	if value.ValueNumber != 1.25 {
		t.Fatalf("ValueNumber = %v, want 1.25", value.ValueNumber)
	}

	if got.GetUnit() != "kg" {
		t.Fatalf("Unit = %q, want kg", got.GetUnit())
	}
}

func TestProductAttributeValueDetailsToProtoReturnsParseErrorForInvalidNumber(t *testing.T) {
	t.Parallel()

	_, err := productAttributeValueDetailsToProto(&catalog.ProductAttributeValueDetails{
		AttributeID: catalog.AttributeID("attr_weight"),
		Code:        "weight",
		DisplayName: "Weight",
		DataType:    catalog.ProductAttributeDataTypeNumber,
		ValueNumber: "heavy",
	})
	if err == nil {
		t.Fatal("productAttributeValueDetailsToProto() error = nil, want parse error")
	}
}

func TestProductAttributeValueDetailsToProtoMapsBooleanFalseValue(t *testing.T) {
	t.Parallel()

	got, err := productAttributeValueDetailsToProto(&catalog.ProductAttributeValueDetails{
		AttributeID:  catalog.AttributeID("attr_dimmable"),
		Code:         "dimmable",
		DisplayName:  "Dimmable",
		DataType:     catalog.ProductAttributeDataTypeBool,
		ValueBoolean: grpcadapterBoolPtr(false),
	})
	if err != nil {
		t.Fatalf("productAttributeValueDetailsToProto() error = %v, want nil", err)
	}

	value, ok := got.GetValue().(*catalogv1.ProductAttributeValue_ValueBoolean)
	if !ok {
		t.Fatalf("Value type = %T, want ValueBoolean", got.GetValue())
	}

	if value.ValueBoolean {
		t.Fatal("ValueBoolean = true, want false")
	}
}

func TestProductAttributeValueDetailsToProtoMapsJSONValue(t *testing.T) {
	t.Parallel()

	got, err := productAttributeValueDetailsToProto(&catalog.ProductAttributeValueDetails{
		AttributeID: catalog.AttributeID("attr_specs"),
		Code:        "specs",
		DisplayName: "Specs",
		DataType:    catalog.ProductAttributeDataTypeJson,
		ValueJSON:   []byte(`{"lumens":800}`),
	})
	if err != nil {
		t.Fatalf("productAttributeValueDetailsToProto() error = %v, want nil", err)
	}

	value, ok := got.GetValue().(*catalogv1.ProductAttributeValue_ValueJson)
	if !ok {
		t.Fatalf("Value type = %T, want ValueJson", got.GetValue())
	}

	if value.ValueJson != `{"lumens":800}` {
		t.Fatalf("ValueJson = %q, want JSON string", value.ValueJson)
	}
}

func TestProductAttributeValueDetailsToProtoMapsValueOptions(t *testing.T) {
	t.Parallel()

	got, err := productAttributeValueDetailsToProto(&catalog.ProductAttributeValueDetails{
		AttributeID: catalog.AttributeID("attr_colour"),
		Code:        "colour",
		DisplayName: "Colour",
		DataType:    catalog.ProductAttributeDataTypeOption,
		ValueString: "blue",
		Options: []*catalog.ProductAttributeOption{
			{OptionID: catalog.OptionID("opt_blue"), Value: "blue", DisplayName: "Blue", Status: catalog.ProductAttributeOptionsStatusActive},
			{OptionID: catalog.OptionID("opt_green"), Value: "green", DisplayName: "Green", Status: catalog.ProductAttributeOptionsStatusActive},
		},
	})
	if err != nil {
		t.Fatalf("productAttributeValueDetailsToProto() error = %v, want nil", err)
	}

	if len(got.GetValueOptions()) != 2 {
		t.Fatalf("len(ValueOptions) = %d, want 2", len(got.GetValueOptions()))
	}

	if got.GetValueOptions()[0] != "blue" {
		t.Fatalf("ValueOptions[0] = %q, want blue", got.GetValueOptions()[0])
	}
}

func TestProductAttributeValueDetailsToProtoRejectsMissingValue(t *testing.T) {
	t.Parallel()

	_, err := productAttributeValueDetailsToProto(&catalog.ProductAttributeValueDetails{
		AttributeID: catalog.AttributeID("attr_empty"),
		Code:        "empty",
		DisplayName: "Empty",
		DataType:    catalog.ProductAttributeDataTypeString,
	})
	if err == nil {
		t.Fatal("productAttributeValueDetailsToProto() error = nil, want error")
	}
}

func TestProductVariantDetailsToProtoMapsVariant(t *testing.T) {
	t.Parallel()

	now := grpcadapterTestTime()

	got, err := productVariantDetailsToProto(&catalog.ProductVariantDetails{
		VariantID:   catalog.VariantID("var_blue"),
		ProductID:   catalog.ProductID("prod_gopher_lamp"),
		Sku:         catalog.Sku("BFS-GO-LAMP-BLUE"),
		VariantName: "Blue shade",
		Status:      catalog.ProductVariantStatusActive,
		Price:       catalog.Money{AmountMinor: 5499, CurrencyCode: "GBP"},
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Hour),
		Attributes: []*catalog.ProductAttributeValueDetails{
			{
				AttributeID: catalog.AttributeID("attr_colour"),
				Code:        "colour",
				DisplayName: "Colour",
				DataType:    catalog.ProductAttributeDataTypeString,
				ValueString: "blue",
			},
		},
	})
	if err != nil {
		t.Fatalf("productVariantDetailsToProto() error = %v, want nil", err)
	}

	if got.GetVariantId() != "var_blue" {
		t.Fatalf("VariantId = %q, want var_blue", got.GetVariantId())
	}

	if got.GetSku() != "BFS-GO-LAMP-BLUE" {
		t.Fatalf("Sku = %q, want BFS-GO-LAMP-BLUE", got.GetSku())
	}

	if got.GetVariantName() != "Blue shade" {
		t.Fatalf("VariantName = %q, want Blue shade", got.GetVariantName())
	}

	if len(got.GetAttributes()) != 1 {
		t.Fatalf("len(Attributes) = %d, want 1", len(got.GetAttributes()))
	}
}

func TestProductDetailsToProtoMapsFullyHydratedProduct(t *testing.T) {
	t.Parallel()

	now := grpcadapterTestTime()

	got, err := productDetailsToProto(catalog.ProductDetails{
		ProductID:   catalog.ProductID("prod_gopher_lamp"),
		CategoryID:  catalog.CategoryID("cat_lighting"),
		Name:        "Gopher Desk Lamp",
		Slug:        "gopher-desk-lamp",
		Description: "A cheerful lamp.",
		Brand:       "Borough",
		Status:      catalog.ProductStatusActive,
		BasePrice:   catalog.Money{AmountMinor: 4999, CurrencyCode: "GBP"},
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Hour),
		Attributes: []*catalog.ProductAttributeValueDetails{
			{
				AttributeID: catalog.AttributeID("attr_material"),
				Code:        "material",
				DisplayName: "Material",
				DataType:    catalog.ProductAttributeDataTypeString,
				ValueString: "steel",
			},
		},
		Variants: []*catalog.ProductVariantDetails{
			{
				VariantID:   catalog.VariantID("var_blue"),
				ProductID:   catalog.ProductID("prod_gopher_lamp"),
				Sku:         catalog.Sku("BFS-GO-LAMP-BLUE"),
				VariantName: "Blue shade",
				Status:      catalog.ProductVariantStatusActive,
				Price:       catalog.Money{AmountMinor: 5499, CurrencyCode: "GBP"},
				CreatedAt:   now,
				UpdatedAt:   now.Add(time.Hour),
			},
		},
		Images: []*catalog.ProductImage{
			{
				ImageID:      catalog.ImageID("img_primary"),
				ProductID:    catalog.ProductID("prod_gopher_lamp"),
				Url:          "https://example.test/gopher-lamp.png",
				AltText:      "Gopher Desk Lamp on a desk",
				DisplayOrder: 10,
			},
		},
	})
	if err != nil {
		t.Fatalf("productDetailsToProto() error = %v, want nil", err)
	}

	if got.GetProductId() != "prod_gopher_lamp" {
		t.Fatalf("ProductId = %q, want prod_gopher_lamp", got.GetProductId())
	}

	if got.GetSlug() != "gopher-desk-lamp" {
		t.Fatalf("Slug = %q, want gopher-desk-lamp", got.GetSlug())
	}

	if len(got.GetAttributes()) != 1 {
		t.Fatalf("len(Attributes) = %d, want 1", len(got.GetAttributes()))
	}

	if len(got.GetVariants()) != 1 {
		t.Fatalf("len(Variants) = %d, want 1", len(got.GetVariants()))
	}

	if len(got.GetImages()) != 1 {
		t.Fatalf("len(Images) = %d, want 1", len(got.GetImages()))
	}
}

func TestProductAttributeOptionToProtoMapsOption(t *testing.T) {
	t.Parallel()

	got, err := productAttributeOptionToProto(&catalog.ProductAttributeOption{
		OptionID:     catalog.OptionID("opt_blue"),
		Value:        "blue",
		DisplayName:  "Blue",
		DisplayOrder: 10,
		Status:       catalog.ProductAttributeOptionsStatusActive,
	})
	if err != nil {
		t.Fatalf("productAttributeOptionToProto() error = %v, want nil", err)
	}

	if got.GetOptionId() != "opt_blue" {
		t.Fatalf("OptionId = %q, want opt_blue", got.GetOptionId())
	}

	if got.GetStatus() != catalogv1.ProductStatus_PRODUCT_STATUS_ACTIVE {
		t.Fatalf("Status = %v, want ACTIVE", got.GetStatus())
	}
}

func TestProductAttributeDefinitionToProtoMapsDefinitionWithOptions(t *testing.T) {
	t.Parallel()

	got, err := productAttributeDefinitionToProto(&catalog.ProductAttributeDefinition{
		AttributeID:       catalog.AttributeID("attr_colour"),
		CategoryID:        catalog.CategoryID("cat_lighting"),
		Code:              "colour",
		DisplayName:       "Colour",
		Description:       "Lamp shade colour.",
		DataType:          catalog.ProductAttributeDataTypeOption,
		IsRequired:        true,
		IsFilterable:      true,
		IsVariantDefining: true,
		Status:            catalog.ProductAttributeDefinitionStatusActive,
		Options: []*catalog.ProductAttributeOption{
			{
				OptionID:     catalog.OptionID("opt_blue"),
				Value:        "blue",
				DisplayName:  "Blue",
				DisplayOrder: 10,
				Status:       catalog.ProductAttributeOptionsStatusActive,
			},
		},
	})
	if err != nil {
		t.Fatalf("productAttributeDefinitionToProto() error = %v, want nil", err)
	}

	if got.GetAttributeId() != "attr_colour" {
		t.Fatalf("AttributeId = %q, want attr_colour", got.GetAttributeId())
	}

	if got.GetDataType() != catalogv1.ProductAttributeDataType_PRODUCT_ATTRIBUTE_DATA_TYPE_OPTION {
		t.Fatalf("DataType = %v, want OPTION", got.GetDataType())
	}

	if !got.GetIsRequired() {
		t.Fatal("IsRequired = false, want true")
	}

	if !got.GetIsFilterable() {
		t.Fatal("IsFilterable = false, want true")
	}

	if !got.GetIsVariantDefining() {
		t.Fatal("IsVariantDefining = false, want true")
	}

	if len(got.GetOptions()) != 1 {
		t.Fatalf("len(Options) = %d, want 1", len(got.GetOptions()))
	}
}
