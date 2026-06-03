package catalog

import "errors"

// Domain errors.
var (
	ErrInvalidProductID = errors.New("invalid product id")
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidPageToken = errors.New("invalid page token")
	ErrInvalidDisplayOrder = errors.New("invalid display order")	
)

