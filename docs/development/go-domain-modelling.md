# Go Domain Modelling

This document defines domain modelling guidance for bfstore Go services.

## Core Rule

```text
If a concept matters to the business, give it a type.
```

## Domain IDs

Prefer defined types:

```go
type ProductID string
type BasketID string
type OrderID string
type PaymentAttemptID string
type ReservationID string
```

## Money

Do not use `float64` for money.

```go
type Money struct {
    AmountMinor int64
    Currency    string
}
```

## Status Types

Avoid booleans when business state has more than two meanings.

```go
type PaymentStatus string

const (
    PaymentStatusPending    PaymentStatus = "PENDING"
    PaymentStatusAuthorised PaymentStatus = "AUTHORISED"
    PaymentStatusDeclined   PaymentStatus = "DECLINED"
    PaymentStatusFailed     PaymentStatus = "FAILED"
)
```

## Basket

```go
type BasketItem struct {
    ProductID ProductID
    Quantity  int
    UnitPrice Money
}

type Basket struct {
    ID     BasketID
    Items  []BasketItem
    Status BasketStatus
}
```

## Order

Use methods to protect state transitions.

```go
func (o *Order) MarkPaymentAuthorised() error {
    if o.Status != OrderStatusPendingPayment {
        return errors.New("order is not pending payment")
    }
    o.Status = OrderStatusPaymentAuthorised
    return nil
}
```

## Zero-value Pitfalls

Check:

```text
empty string ID
false boolean state
zero amount
nil slice
nil map
```

## Practical Rules

```text
Use defined types for IDs.
Use integer minor units for money.
Use explicit status types.
Avoid business meaning hidden in zero values.
Put behaviour near data when it protects invariants.
Keep domain models separate from Protobuf transport models.
```

## Final Rule

```text
Domain modelling is where Go types start protecting bfstore from silly mistakes.
```
