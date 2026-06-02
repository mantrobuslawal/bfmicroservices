# Go Data and Allocation

This document defines bfstore conventions for Go data structures and allocation.

## Core Rule

```text
Make zero values useful where natural.
Reject them where they hide invalid business state.
```

## new and make

Use `make` for slices, maps, and channels.

```go
items := make([]BasketItem, 0, 10)
attributes := make(map[string]string)
jobs := make(chan ProductID, 100)
```

## Zero Values

Good zero-value design:

```go
var b Basket
b.AddItem(item)
```

Use constructors where zero value would be invalid.

```go
func NewMoney(amountMinor int64, currency string) (Money, error) {
    if currency == "" {
        return Money{}, ErrCurrencyRequired
    }
    if amountMinor < 0 {
        return Money{}, ErrNegativeAmount
    }
    return Money{AmountMinor: amountMinor, Currency: currency}, nil
}
```

## Slices and append

If a function appends to a slice it receives, return the new slice.

```go
func AddItem(items []BasketItem, item BasketItem) []BasketItem {
    return append(items, item)
}
```

Or mutate through an owning struct.

```go
func (b *Basket) AddItem(item BasketItem) {
    b.items = append(b.items, item)
}
```

## Maps

Initialise maps before assignment.

```go
attrs := make(map[string]string)
attrs["colour"] = "blue"
```

Remember:

```text
map iteration order is not stable
maps are reference-like
maps are not safe for concurrent writes without synchronisation
```

## String Methods

Avoid recursion traps.

```go
func (id ProductID) String() string {
    return string(id)
}
```

## Type Assertions

Use comma-ok form.

```go
productID, ok := value.(ProductID)
if !ok {
    return ErrInvalidProductID
}
```

## Blank Identifier

Use `_` deliberately, not to hide ignored errors.

Bad:

```go
_ = err
```

Good:

```go
if err != nil {
    return err
}
```

## Practical Rules

```text
Use make for slices, maps, and channels.
Use constructors for invalid zero values.
Return appended slices.
Protect maps used concurrently.
Avoid recursive String methods.
Use comma-ok assertions.
Do not ignore errors with _.
```

## Final Rule

```text
Data structure choices should make business rules safer, not more mysterious.
```
