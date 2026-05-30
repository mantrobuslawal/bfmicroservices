# Go Style and Design

This document defines Go style and design guidance for bfstore services.

## Core Rule

```text
Write Go that is simple to read, easy to test, and boring to operate.
```

## Package Layout

Recommended service layout:

```text
services/order/
  cmd/order-service/
  internal/order/
  internal/repository/
  internal/transport/grpc/
  internal/config/
```

Recommended shared platform layout:

```text
pkg/platform/
  config/
  logging/
  telemetry/
  grpc/
  shutdown/
```

## Naming

Prefer clear domain names:

```go
type OrderID string
type PaymentAttemptID string
type StockReservationID string
```

Avoid vague names like `Data`, `Manager`, and `Helper`.

## Receiver Choices

Use pointer receivers for mutation or large structs.

Use value receivers for small immutable value objects.

## Interface Placement

Define interfaces where they are consumed.

Avoid giant exported interfaces created “just in case”.

## Constructors

Use constructors when invariants matter.

```go
func NewMoney(amountMinor int64, currency string) (Money, error) {
    if amountMinor < 0 {
        return Money{}, errors.New("amount must not be negative")
    }
    if currency == "" {
        return Money{}, errors.New("currency is required")
    }
    return Money{AmountMinor: amountMinor, Currency: currency}, nil
}
```

## Transport vs Domain

Keep generated Protobuf types near the transport layer.

Recommended flow:

```text
gRPC request
  -> transport validation/mapping
  -> domain command
  -> application service
  -> repository/client
  -> response mapping
```

## Practical Rules

```text
Keep packages small and purposeful.
Use clear names.
Use domain types.
Define interfaces where they are consumed.
Keep generated transport types near the transport layer.
Use constructors for important invariants.
Avoid clever abstractions until repetition proves they are needed.
```

## Final Rule

```text
Good Go design should make the boring path obvious.
```
