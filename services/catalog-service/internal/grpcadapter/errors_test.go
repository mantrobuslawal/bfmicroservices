package grpcadapter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapServiceErrorMapsInvalidArgumentErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{name: "invalid product id", err: catalog.ErrInvalidProductID},
		{name: "invalid category id", err: catalog.ErrInvalidCategoryID},
		{name: "invalid page size", err: catalog.ErrInvalidPageSize},
		{name: "invalid page token", err: catalog.ErrInvalidPageToken},
		{name: "invalid display order", err: catalog.ErrInvalidDisplayOrder},
		{name: "wrapped invalid product id", err: fmt.Errorf("wrap: %w", catalog.ErrInvalidProductID)},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mapServiceError(tt.err)

			if status.Code(got) != codes.InvalidArgument {
				t.Fatalf("status.Code() = %v, want %v", status.Code(got), codes.InvalidArgument)
			}
		})
	}
}

func TestMapServiceErrorMapsNotFound(t *testing.T) {
	t.Parallel()

	got := mapServiceError(catalog.ErrProductNotFound)

	if status.Code(got) != codes.NotFound {
		t.Fatalf("status.Code() = %v, want %v", status.Code(got), codes.NotFound)
	}
}

func TestMapServiceErrorMapsUnknownErrorsToInternal(t *testing.T) {
	t.Parallel()

	got := mapServiceError(errors.New("database exploded"))

	if status.Code(got) != codes.Internal {
		t.Fatalf("status.Code() = %v, want %v", status.Code(got), codes.Internal)
	}

	if status.Convert(got).Message() != "internal server error" {
		t.Fatalf("message = %q, want internal server error", status.Convert(got).Message())
	}
}
