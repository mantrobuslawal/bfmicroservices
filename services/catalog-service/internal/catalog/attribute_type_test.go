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

func TestProductAttributeValue_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value ProductAttributeValue
		want  bool
	}{
		{
			name:  "string is valid",
			value: ProductAttributeValue("string"),
			want:  true,
		},
		{
			name:  "number is valid",
			value: ProductAttributeValue("number"),
			want:  true,
		},
		{
			name:  "boolean is valid",
			value: ProductAttributeValue("boolean"),
			want:  true,
		},
		{
			name:  "option is valid",
			value: ProductAttributeValue("option"),
			want:  true,
		},
		{
			name:  "multi option is valid",
			value: ProductAttributeValue("multi_option"),
			want:  true,
		},
		{
			name:  "json is valid",
			value: ProductAttributeValue("json"),
			want:  true,
		},
		{
			name:  "unspecified is invalid",
			value: ProductAttributeValue("unspecified"),
			want:  false,
		},
		{
			name:  "empty is invalid",
			value: ProductAttributeValue(""),
			want:  false,
		},
		{
			name:  "unknown is invalid",
			value: ProductAttributeValue("banana"),
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

func TestProductAttributeValue_String(t *testing.T) {
	t.Parallel()

	value := ProductAttributeValue("json")

	if got, want := value.String(), "json"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
