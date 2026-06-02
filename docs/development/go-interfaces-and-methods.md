# Go Interfaces and Methods

This document defines bfstore interface and method conventions.

## Core Rule

```text
Use interfaces to describe behaviour at the point of consumption.
```

## Consumer-side Interfaces

Define interfaces where they are used.

```go
type ProductRepository interface {
    Get(ctx context.Context, id ProductID) (Product, error)
    List(ctx context.Context, filter Filter) ([]Product, error)
}
```

## Small Interfaces

Small interfaces are useful, but avoid creating noise.

Good:

```go
type Publisher interface {
    Publish(ctx context.Context, event Event) error
}
```

Avoid unnecessary interface fragmentation.

## One-method Interface Names

Use `-er` names where natural.

```go
type Publisher interface {
    Publish(ctx context.Context, event Event) error
}
```

## Behaviour Methods

Prefer behaviour methods when invariants matter.

```go
func (o *Order) MarkPaymentAuthorised() error {
    if o.status != OrderStatusPendingPayment {
        return ErrInvalidOrderTransition
    }

    o.status = OrderStatusPaymentAuthorised
    return nil
}
```

Avoid generic setters that bypass rules.

## Receiver Choices

Use pointer receivers when:

```text
method mutates the receiver
struct is large
identity matters
```

Use value receivers when:

```text
type is small
type is immutable-like
copying is cheap
```

## Constructors

Return concrete types when the type has useful behaviour.

```go
func NewService(products ProductRepository) *Service {
    return &Service{products: products}
}
```

Return interfaces when the abstraction itself is the product.

## Practical Rules

```text
Define interfaces at consumer boundaries.
Keep interfaces small and meaningful.
Avoid interface soup.
Prefer behaviour methods over generic setters.
Use pointer receivers for mutation.
Use concrete constructors by default.
Return interfaces only when it clarifies the abstraction.
```

## Final Rule

```text
Interfaces should reduce coupling, not increase ceremony.
```
