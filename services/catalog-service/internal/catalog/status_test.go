package catalog

import "testing"

func TestLifecycleStatuses_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{
			name: "product status draft is valid",
			got:  ProductStatus("draft").IsValid(),
			want: true,
		},
		{
			name: "product status active is valid",
			got:  ProductStatus("active").IsValid(),
			want: true,
		},
		{
			name: "product status inactive is valid",
			got:  ProductStatus("inactive").IsValid(),
			want: true,
		},
		{
			name: "product status archived is valid",
			got:  ProductStatus("archived").IsValid(),
			want: true,
		},
		{
			name: "product status empty is invalid",
			got:  ProductStatus("").IsValid(),
			want: false,
		},
		{
			name: "product status unknown is invalid",
			got:  ProductStatus("deleted").IsValid(),
			want: false,
		},
		{
			name: "category status active is valid",
			got:  CategoryStatus("active").IsValid(),
			want: true,
		},
		{
			name: "category status unknown is invalid",
			got:  CategoryStatus("hidden").IsValid(),
			want: false,
		},
		{
			name: "variant status active is valid",
			got:  ProductVariantStatus("active").IsValid(),
			want: true,
		},
		{
			name: "variant status unknown is invalid",
			got:  ProductVariantStatus("backorder").IsValid(),
			want: false,
		},
		{
			name: "attribute definition status active is valid",
			got:  ProductAttributeDefinitionStatus("active").IsValid(),
			want: true,
		},
		{
			name: "attribute definition status unknown is invalid",
			got:  ProductAttributeDefinitionStatus("pending").IsValid(),
			want: false,
		},
		{
			name: "attribute option status active is valid",
			got:  ProductAttributeOptionStatus("active").IsValid(),
			want: true,
		},
		{
			name: "attribute option status unknown is invalid",
			got:  ProductAttributeOptionStatus("expired").IsValid(),
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestLifecycleStatuses_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "product status string",
			got:  ProductStatus("active").String(),
			want: "active",
		},
		{
			name: "category status string",
			got:  CategoryStatus("draft").String(),
			want: "draft",
		},
		{
			name: "variant status string",
			got:  ProductVariantStatus("inactive").String(),
			want: "inactive",
		},
		{
			name: "attribute definition status string",
			got:  ProductAttributeDefinitionStatus("archived").String(),
			want: "archived",
		},
		{
			name: "attribute option status string",
			got:  ProductAttributeOptionStatus("active").String(),
			want: "active",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.got != tt.want {
				t.Fatalf("String() = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestIsKnownLifecycleStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{
			name: "known product status",
			got:  isKnownLifecycleStatus(ProductStatus("active")),
			want: true,
		},
		{
			name: "known category status",
			got:  isKnownLifecycleStatus(CategoryStatus("archived")),
			want: true,
		},
		{
			name: "known variant status",
			got:  isKnownLifecycleStatus(ProductVariantStatus("draft")),
			want: true,
		},
		{
			name: "known attribute definition status",
			got:  isKnownLifecycleStatus(ProductAttributeDefinitionStatus("inactive")),
			want: true,
		},
		{
			name: "known attribute option status",
			got:  isKnownLifecycleStatus(ProductAttributeOptionStatus("active")),
			want: true,
		},
		{
			name: "unknown product status",
			got:  isKnownLifecycleStatus(ProductStatus("whatever")),
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.got != tt.want {
				t.Fatalf("isKnownLifecycleStatus() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}
