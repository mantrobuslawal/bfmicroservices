# Go Error Handling

This document defines Go error handling guidance for bfstore services.

## Core Rule

```text
Errors are part of the contract, not an afterthought.
```

## Explicit Error Returns

```go
product, err := repo.GetProduct(ctx, productID)
if err != nil {
    return Product{}, err
}
```

## Sentinel Errors

```go
var ErrProductNotFound = errors.New("product not found")
var ErrBasketEmpty = errors.New("basket is empty")
var ErrPaymentDeclined = errors.New("payment declined")
```

Check with `errors.Is`.

## Error Wrapping

```go
return fmt.Errorf("get product %s: %w", productID, err)
```

## gRPC Status Mapping

```go
switch {
case errors.Is(err, ErrProductNotFound):
    return nil, status.Error(codes.NotFound, "product not found")
case errors.Is(err, ErrBasketEmpty):
    return nil, status.Error(codes.FailedPrecondition, "basket is empty")
case errors.Is(err, ErrPaymentDeclined):
    return nil, status.Error(codes.FailedPrecondition, "payment declined")
default:
    return nil, status.Error(codes.Internal, "internal error")
}
```

## Retryable vs Non-retryable

Retryable examples:

```text
temporary network failure
database deadlock
Kafka broker temporarily unavailable
payment provider timeout where idempotency exists
```

Non-retryable examples:

```text
invalid product ID
basket is empty
payment declined
unauthorised request
```

## Logging Errors

Log useful context, not sensitive payloads.

Include:

```text
event
service
trace_id
correlation_id
error_code
safe business ID
```

Avoid:

```text
secrets
tokens
payment details
raw provider payloads
full request bodies
```

## Practical Rules

```text
Return errors explicitly.
Use errors.Is/errors.As for classification.
Wrap errors with context.
Map domain errors to gRPC status codes.
Separate validation, retryable, and fatal errors.
Do not leak sensitive details.
Log errors with safe structured fields.
```

## Final Rule

```text
Good error handling makes failure understandable without making it unsafe.
```
