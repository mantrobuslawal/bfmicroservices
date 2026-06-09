package grpcadapter

import (
	"errors"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapServiceError(err error) error {
	switch {
	case errors.Is(err, catalog.ErrInvalidProductID),
		errors.Is(err, catalog.ErrInvalidPageSize),
		errors.Is(err, catalog.ErrInvalidPageToken),
		errors.Is(err, catalog.ErrInvalidDisplayOrder):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, catalog.ErrProductNotFound):
		return status.Error(codes.NotFound, err.Error())

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
