# Go Modules and Packages

This document defines bfstore's Go module and package conventions.

## Core Rule

```text
Use one Go module at the repo root until there is a real release-boundary reason to split.
```

## Module Path

Use the real repository path:

```go
module github.com/mantrobuslawal/bfstore
```

Example import:

```go
import "github.com/mantrobuslawal/bfstore/pkg/platform/logging"
```

## go.mod and go.sum

Commit both:

```text
go.mod
go.sum
```

Rules:

```text
go.mod declares dependency intent
go.sum protects dependency integrity
run go mod tidy after dependency changes
commit dependency changes deliberately
```

## Packages

A package is a directory of Go files compiled together.

Better service package layout:

```text
internal/catalog
internal/repository/mysql
internal/transport/grpc
internal/config
```

## internal Packages

Use `internal` to protect service internals.

Example:

```text
services/catalog/internal/repository
```

should not be imported by other services.

## pkg/platform

Use `pkg/platform` for genuinely shared platform code:

```text
logging
telemetry
config
grpc
shutdown
```

## Practical Rules

```text
Use the real module path from the start.
Avoid multiple modules early.
Keep packages purposeful.
Use internal packages for service-private code.
Use pkg/platform sparingly.
Run go mod tidy.
Commit go.mod and go.sum.
```

## Final Rule

```text
Go modules and packages should make dependency boundaries obvious.
```
