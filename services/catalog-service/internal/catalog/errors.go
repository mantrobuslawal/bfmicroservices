package catalog

import "errors"

// Domain errors.
var (
	// Common Service layer errors.
	ErrInvalidProductID    = errors.New("invalid product id")
	ErrInvalidCategoryID   = errors.New("invalid category id")
	ErrProductNotFound     = errors.New("product not found")
	ErrInvalidPageSize     = errors.New("invalid page size")
	ErrInvalidPageToken    = errors.New("invalid page token")
	ErrInvalidDisplayOrder = errors.New("invalid display order")

	// Lifecycle status error
	ErrInvalidLifecycleStatus = errors.New("invalid lifecycle status")
)
