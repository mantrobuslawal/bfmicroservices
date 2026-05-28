# gRPC Auth Package

This package contains shared authentication and authorisation helpers for **bfstore** gRPC services.

Recommended location:

```text
pkg/platform/grpc/auth/
```

Suggested files:

```text
pkg/platform/grpc/auth/
├── README.md
├── principal.go
├── token_validator.go
├── interceptors.go
└── credentials.go
```

---

## Purpose

This package provides reusable building blocks for:

```text
representing authenticated principals
validating bearer tokens
attaching principals to context
server-side authentication interceptors
server-side authorisation interceptors
client-side bearer token credentials
local development token validation
```

It should keep auth concerns out of business handlers.

---

## Practical Rule

```text
Handlers should receive authenticated context.
They should not repeatedly parse authorization metadata.
```

---

## Principal

A `Principal` represents the authenticated caller.

Example fields:

```go
type Principal struct {
    Subject string
    Roles   []string
    Service string
}
```

Possible meanings:

```text
Subject = user ID, service ID, or token subject
Roles   = admin/customer/service roles
Service = internal service identity, where applicable
```

---

## Token Validator

A token validator verifies raw authorisation metadata.

```go
type TokenValidator interface {
    Validate(ctx context.Context, rawAuthorisation string) (Principal, error)
}
```

The raw value will usually look like:

```text
Bearer <token>
```

For local development, a simple validator can accept:

```text
Bearer local-dev-token
```

Production validators may later validate JWT signatures, token expiry, issuer/audience claims, and service identity claims.

---

## Authentication Interceptor

The authentication interceptor should:

```text
read authorization metadata
validate the token
attach Principal to context
return Unauthenticated if identity is missing/invalid
call the handler only if authentication succeeds
```

---

## Authorisation Interceptor

The authorisation interceptor should:

```text
read Principal from context
check whether the principal can call the method
return PermissionDenied when identity is valid but not allowed
call the handler only if allowed
```

Example rules:

```text
CatalogService/ListProducts -> public/customer/admin
CatalogAdminService/CreateProduct -> admin only
PaymentService/AuthorisePayment -> order-service only
```

---

## Client Credentials

For client-side per-RPC bearer credentials, use a type that implements:

```go
credentials.PerRPCCredentials
```

This package includes a simple `StaticBearerToken` for non-production or internal development scenarios.

Tokens should require transport security:

```go
func (c StaticBearerToken) RequireTransportSecurity() bool {
    return true
}
```

Do not send real tokens over insecure gRPC.

---

## Example Server Usage

```go
grpcServer := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        interceptors.UnaryRecovery(logger),
        interceptors.UnaryCorrelationID(logger),
        interceptors.UnaryLogging(logger),
        auth.UnaryAuthentication(tokenValidator),
        auth.UnaryAuthorisation(authoriser),
    ),
)
```

---

## Example Client Usage

```go
tokenCreds := auth.StaticBearerToken{
    Token: cfg.InternalServiceToken,
}

conn, err := grpc.NewClient(
    "inventory-service:50051",
    grpc.WithTransportCredentials(tlsCreds),
    grpc.WithPerRPCCredentials(tokenCreds),
)
```

---

## Testing Guidance

Recommended tests:

```text
missing authorization metadata returns Unauthenticated
invalid token returns Unauthenticated
valid token attaches Principal to context
principal can be retrieved from context
disallowed principal returns PermissionDenied
allowed principal reaches handler
StaticBearerToken returns authorization metadata
StaticBearerToken requires transport security
```

---

## What This Package Should Not Do

Do not put service-specific business rules directly in this package.

Bad:

```text
auth package knows checkout orchestration rules
auth package knows catalogue admin workflows
auth package knows payment business logic
```

Good:

```text
auth package provides interfaces and interceptors
services provide their own validators/authorisers where needed
```

The shared package provides auth plumbing. Service/application layers own business-specific policy decisions.

---

## Practical Rules

```text
Keep auth explicit.
Keep auth testable.
Do not log raw tokens.
Do not send real tokens over insecure transport.
Use Unauthenticated for missing/invalid identity.
Use PermissionDenied for valid identity without permission.
Prefer mTLS/service identity for internal service-to-service security later.
```

---

## Final Rule

```text
The auth package should prove identity and expose it cleanly.
Business services decide what that identity may do.
```
