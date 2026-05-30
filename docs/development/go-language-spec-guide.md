# Go Language Spec Guide

This document explains how bfstore uses important Go language rules from the Go specification.

## Purpose

```text
types
packages
interfaces
method sets
slices
maps
channels
package initialisation
```

## Core Rule

```text
When production correctness depends on language behaviour, do not guess.
Check the language rule.
```

## Packages

Use packages to create clear boundaries.

```text
services/catalog/
  cmd/catalog-service/
  internal/catalog/
  internal/repository/
  internal/transport/grpc/

pkg/platform/
  logging/
  telemetry/
  config/
  grpc/
```

Use `internal` packages to protect service internals.

## Declarations

Use clear declarations for domain concepts.

```go
type ProductID string
type BasketID string
type OrderID string
type PaymentAttemptID string
```

## Zero Values

Avoid letting zero values accidentally become business meaning.

Bad:

```go
type PaymentAttempt struct {
    Authorised bool
}
```

Better:

```go
type PaymentStatus string
```

with explicit statuses.

## Types

Use types to model business concepts.

```go
type Money struct {
    AmountMinor int64
    Currency    string
}
```

## Interfaces

Define interfaces at consumer boundaries.

```go
type ProductRepository interface {
    GetProduct(ctx context.Context, id ProductID) (Product, error)
}
```

## Slices and Maps

Slices can share underlying arrays. Maps are reference-like and not safe for concurrent writes.

Use copies or locks where needed.

## Channels

Use channels for in-process coordination.

Use Kafka for durable cross-service asynchronous work.

## Package Initialisation

Avoid magical startup in `init()` functions.

Prefer explicit construction in `main`.

## Practical Rules

```text
Use packages to enforce boundaries.
Use internal packages for service internals.
Use domain-specific types.
Avoid magical init functions.
Be careful with slices and maps.
Use channels for in-process coordination.
Use Kafka for durable cross-service messaging.
```

## Final Rule

```text
The Go spec tells us what the language guarantees; bfstore should use those guarantees deliberately.
```
