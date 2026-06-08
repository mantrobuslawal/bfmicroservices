package catalog

import "testing"

func TestProductAttributeDataType_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dataType ProductAttributeDataType
		want     bool
	}{
		{
			name:     "string is valid",
			dataType: ProductAttributeDataType("string"),
			want:     true,
		},
		{
			name:     "number is valid",
			dataType: ProductAttributeDataType("number"),
			want:     true,
		},
		{
			name:     "boolean is valid",
			dataType: ProductAttributeDataType("boolean"),
			want:     true,
		},
		{
			name:     "option is valid",
			dataType: ProductAttributeDataType("option"),
			want:     true,
		},
		{
			name:     "multi option is valid",
			dataType: ProductAttributeDataType("multi_option"),
			want:     true,
		},
		{
			name:     "json is valid",
			dataType: ProductAttributeDataType("json"),
			want:     true,
		},
		{
			name:     "empty is invalid",
			dataType: ProductAttributeDataType(""),
			want:     false,
		},
		{
			name:     "unspecified is invalid",
			dataType: ProductAttributeDataType("unspecified"),
			want:     false,
		},
		{
			name:     "unknown is invalid",
			dataType: ProductAttributeDataType("colourful-chaos"),
			want:     false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.dataType.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductAttributeDataType_String(t *testing.T) {
	t.Parallel()

	dataType := ProductAttributeDataType("multi_option")

	if got, want := dataType.String(), "multi_option"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestProductAttributeValueKind_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value ProductAttributeValueKind
		want  bool
	}{
		{
			name:  "string is valid",
			value: ProductAttributeValueKind("string"),
			want:  true,
		},
		{
			name:  "number is valid",
			value: ProductAttributeValueKind("number"),
			want:  true,
		},
		{
			name:  "boolean is valid",
			value: ProductAttributeValueKind("boolean"),
			want:  true,
		},
		{
			name:  "option is valid",
			value: ProductAttributeValueKind("option"),
			want:  true,
		},
		{
			name:  "multi option is valid",
			value: ProductAttributeValueKind("multi_option"),
			want:  true,
		},
		{
			name:  "json is valid",
			value: ProductAttributeValueKind("json"),
			want:  true,
		},
		{
			name:  "unspecified is invalid",
			value: ProductAttributeValueKind("unspecified"),
			want:  false,
		},
		{
			name:  "empty is invalid",
			value: ProductAttributeValueKind(""),
			want:  false,
		},
		{
			name:  "unknown is invalid",
			value: ProductAttributeValueKind("banana"),
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.value.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductAttributeValueKind_String(t *testing.T) {
	t.Parallel()

	value := ProductAttributeValueKind("json")

	if got, want := value.String(), "json"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
