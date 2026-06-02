# Go Control Flow and Errors

This document defines bfstore control-flow and error-handling conventions for Go code.

## Core Rule

```text
Handle failure early.
Let the happy path keep walking down the page.
```

## Early Returns

Avoid unnecessary `else` after return.

```go
product, err := s.products.Get(ctx, id)
if err != nil {
    return Product{}, err
}

return product, nil
```

## Validation

Validate commands early.

```go
func (s *Service) Checkout(ctx context.Context, cmd CheckoutCommand) (Order, error) {
    if err := cmd.Validate(); err != nil {
        return Order{}, err
    }

    // continue happy path
}
```

## switch

Use switch for readable business states.

```go
func (o *Order) CanCancel() bool {
    switch o.status {
    case OrderStatusCreated, OrderStatusPendingPayment:
        return true
    case OrderStatusShipped, OrderStatusCancelled:
        return false
    default:
        return false
    }
}
```

## Explicit Errors

If something can fail, make failure explicit.

```go
func (r *Repository) GetProduct(ctx context.Context, id ProductID) (Product, error) {
    // ...
}
```

## Wrapping Errors

Wrap errors when adding useful context.

```go
return Product{}, fmt.Errorf("get product %s: %w", productID, err)
```

## gRPC Mapping

Map domain errors at the transport boundary.

```go
switch {
case errors.Is(err, catalog.ErrInvalidProductID):
    return nil, status.Error(codes.InvalidArgument, "invalid product_id")
case errors.Is(err, catalog.ErrProductNotFound):
    return nil, status.Error(codes.NotFound, "product not found")
default:
    return nil, status.Error(codes.Internal, "failed to get product")
}
```

## defer

Use `defer` immediately after successful acquisition.

```go
rows, err := db.QueryContext(ctx, query)
if err != nil {
    return err
}
defer rows.Close()
```

## Practical Rules

```text
Prefer early returns.
Avoid unnecessary else.
Use switch for state classification.
Return explicit errors.
Wrap errors with context.
Map errors at boundaries.
Use defer for cleanup.
Do not panic for normal business errors.
```

## Final Rule

```text
Clear control flow makes failure easier to reason about.
```
