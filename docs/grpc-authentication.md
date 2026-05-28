# gRPC Authentication

This document defines the recommended authentication approach for **bfstore** gRPC services.

bfstore uses gRPC for synchronous service-to-service commands and queries. Authentication ensures callers are identified, connections are protected, and sensitive operations are not exposed to untrusted callers.

---

## Purpose

This document explains:

```text
local development authentication expectations
TLS and mTLS direction
token-based authentication
authentication vs authorisation
metadata rules
interceptor usage
service identity rules
recommended rollout phases
```

---

## Authentication vs Authorisation

```text
Authentication = Who are you?
Authorisation  = Are you allowed to do this?
```

Authentication proves identity. Authorisation applies access rules to that identity.

bfstore must keep these concerns separate.

---

## Local Development

Local Docker Compose may use insecure gRPC transport:

```go
grpc.WithTransportCredentials(insecure.NewCredentials())
```

This is acceptable only for local Docker Compose, local development, local integration tests, and temporary scaffolding.

Local development must not use real production secrets or real payment credentials.

Auth interceptors may initially use a local development validator such as `Bearer local-dev-token`.

---

## Production Direction

For production-like environments, bfstore should move towards:

```text
TLS for encrypted service traffic
mTLS for internal service identity
server-side authentication interceptors
server-side authorisation interceptors
short-lived tokens where tokens are used
no raw token logging
```

For sensitive flows such as payment authorisation, prefer mTLS/service identity over trusting caller-supplied metadata alone.

---

## Channel Credentials and Call Credentials

Use channel credentials for transport security:

```text
TLS
mTLS
```

Use call credentials for per-request authentication data:

```text
Bearer tokens
service tokens
request-scoped credentials
```

Practical rule:

```text
Use channel credentials to secure the connection.
Use call credentials or metadata to describe the caller/request.
```

---

## Metadata Rules

Authentication tokens usually travel as gRPC metadata:

```text
authorization: Bearer <token>
```

Rules:

```text
Use the standard metadata key: authorization.
Do not log raw authorization metadata.
Do not send tokens over insecure connections outside local development.
Do not put business payloads in metadata.
Do not trust diagnostic metadata such as x-bfstore-client for security decisions.
```

The key is intentionally spelled `authorization`, even though project documentation uses UK English generally.

---

## Interceptor Rules

Server-side authentication should live in interceptors.

Recommended server chain:

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

Handlers should not repeatedly parse and validate tokens.

---

## Status Code Rules

| Condition | Code |
|---|---|
| Missing token | `Unauthenticated` |
| Invalid token | `Unauthenticated` |
| Expired token | `Unauthenticated` |
| Valid identity but disallowed method | `PermissionDenied` |
| Auth validator unavailable due to infrastructure failure | `Unavailable` or `Internal`, depending on cause |
| Authorisation engine fails unexpectedly | `Internal` |

Do not return `Internal` for ordinary authentication failures.

---

## Service Identity

bfstore should define stable service identity names:

```text
api-gateway
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-service
```

These identities may later map to mTLS certificate subjects, SPIFFE IDs, service mesh identities, or internal service JWT claims.

Do not use `x-bfstore-client` alone for security decisions. It is useful for diagnostics, not proof of identity.

---

## Recommended Rollout Phases

### Phase 1: Local development shape

```text
Use insecure credentials inside Docker Compose only.
Create auth package and interfaces.
Use local dev token validator where needed.
Do not send real secrets.
```

### Phase 2: Token validation

```text
Use authorization metadata.
Validate tokens in server interceptors.
Use TLS where real tokens are transmitted.
Do not log tokens.
```

### Phase 3: TLS/mTLS

```text
Use TLS for encrypted service traffic.
Use mTLS for internal service identity.
Use service identity for authorisation decisions.
```

### Phase 4: Platform identity

```text
mTLS may be handled by the mesh
service identity may come from SPIFFE/SPIRE or mesh identity
application auth still validates user/request context where required
```

---

## Sensitive Flow: Checkout

Checkout may involve:

```text
api-gateway -> order-service
order-service -> inventory-service
order-service -> payment-service
order-service -> shipping-service
```

Recommended direction:

```text
api-gateway validates user token
order-service receives authenticated request context
order-service calls downstream services using internal service identity
payment-service only accepts AuthorisePayment from trusted order-service identity
```

Payment calls should not rely on arbitrary caller-supplied metadata.

---

## What Not To Do

Avoid:

```text
sending tokens over insecure gRPC outside local development
logging authorization headers
trusting x-bfstore-client as proof of identity
using Google OAuth tokens for bfstore internal services
copying auth parsing into every handler
returning Internal for missing/invalid tokens
putting payment card data or personal data in metadata
```

---

## Testing Guidance

Recommended tests:

```text
missing authorization returns Unauthenticated
invalid authorization returns Unauthenticated
valid token attaches Principal to context
valid principal but disallowed method returns PermissionDenied
allowed principal reaches handler
auth interceptor does not call handler when auth fails
authorisation interceptor does not call handler when permission fails
tokens are not logged
```

---

## Practical Rules

```text
Use insecure credentials only for local development.
Use TLS/mTLS for real service-to-service traffic.
Use authorization metadata for bearer tokens.
Do not log raw tokens.
Validate tokens in server interceptors.
Separate authentication from authorisation.
Use Unauthenticated for missing/invalid identity.
Use PermissionDenied for valid identity without permission.
Use service identity for internal authorisation where possible.
Keep auth boring, explicit, and testable.
```

---

## Final Rule

```text
Authentication is identity.
Authorisation is permission.
Transport security is protection.
```

bfstore should start simple but keep clean seams so stronger authentication can be added without rewriting business handlers.
