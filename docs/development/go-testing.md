# Go Testing

This document defines bfstore's Go testing conventions.

## Core Rule

```text
go test ./... should be boring enough to run constantly.
```

## Test Files

Go test files end with:

```text
_test.go
```

Test functions use:

```go
func TestXxx(t *testing.T) {
}
```

## Running Tests

```bash
go test ./...
go test ./services/basket/internal/basket
go test -v ./...
go test -race ./...
```

## Table-driven Tests

Use table-driven tests for validation, mapping, and edge cases.

```go
func TestValidateQuantity(t *testing.T) {
    tests := []struct {
        name    string
        qty     int
        wantErr bool
    }{
        {name: "zero", qty: 0, wantErr: true},
        {name: "one", qty: 1, wantErr: false},
        {name: "too high", qty: 100, wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateQuantity(tt.qty)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ValidateQuantity(%d) error = %v, wantErr %v", tt.qty, err, tt.wantErr)
            }
        })
    }
}
```

## Unit Tests

Unit tests should cover:

```text
domain validation
state transitions
mapping
error classification
small pure functions
```

## Integration Tests

Integration tests may cover:

```text
MySQL repositories
Kafka producers/consumers
gRPC handlers
outbox behaviour
```

## Practical Rules

```text
Keep unit tests fast.
Use table-driven tests for edge cases.
Separate unit and integration test assumptions.
Avoid hidden local machine dependencies.
Run tests in CI.
Use race tests for concurrency-sensitive code.
```

## Final Rule

```text
Tests should make Go changes safer, not make developers afraid to touch the repo.
```
