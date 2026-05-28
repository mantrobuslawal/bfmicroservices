# gRPC Metadata Conventions

This document defines the standard gRPC metadata keys used across **bfstore** services.

Metadata is used for call context such as tracing, correlation, authentication, and diagnostics. Business data belongs in Protobuf request and response messages, not metadata.

---

## Purpose

The purpose of this standard is to keep service-to-service calls consistent, observable, and secure.

bfstore services use gRPC metadata for:

```text
correlation IDs
request IDs
trace propagation
authentication context
caller identification
small diagnostic hints
```

bfstore services must not use metadata for:

```text
basket items
product details
order totals
payment amounts
shipping addresses
large JSON payloads
raw image bytes
debug dumps
```

---

## Practical Rule

```text
Protobuf messages carry business meaning.
gRPC metadata carries call context.
```

If a value is required by the service to perform business logic, it usually belongs in the Protobuf message.

If a value is required by platform middleware to trace, authenticate, authorise, route, or observe the call, it usually belongs in metadata.

---

## Standard Metadata Keys

| Key | Purpose | Direction | Required | Notes |
|---|---|---|---|---|
| `x-correlation-id` | Trace one business request across services | Incoming and outgoing | Generated if absent at the edge | Should remain stable across a full user/business flow |
| `x-request-id` | Identify one RPC/request instance | Incoming and outgoing | Generated per service boundary if absent | Useful for distinguishing retries or repeated calls |
| `authorization` | Bearer token or auth credential | Client/API gateway to service | Required for protected APIs | Do not log raw values |
| `traceparent` | W3C/OpenTelemetry trace propagation | Incoming and outgoing | Required when tracing is enabled | Usually managed by tracing middleware/instrumentation |
| `x-bfstore-client` | Identify the calling service/client | Outgoing | Optional | Useful in local/dev diagnostics and logs |

---

## Key Naming Rules

Use lowercase metadata keys.

Good:

```text
x-correlation-id
x-request-id
authorization
traceparent
x-bfstore-client
```

Avoid:

```text
X-Correlation-ID
Grpc-Correlation-Id
grpc-custom-key
```

The `grpc-` prefix is reserved for gRPC internals and must not be used for custom bfstore metadata keys.

---

## Suggested Go Constants

Define shared constants in a small platform package, for example:

```text
pkg/platform/grpc/metadata
```

Example:

```go
package grpcmetadata

const (
    CorrelationID = "x-correlation-id"
    RequestID     = "x-request-id"
    Authorisation = "authorization"
    Traceparent   = "traceparent"
    BFStoreClient = "x-bfstore-client"
)
```

Although UK English normally uses `authorisation`, the standard metadata key should remain:

```text
authorization
```

This aligns with common HTTP/gRPC conventions.

---

## `x-correlation-id`

### Purpose

`x-correlation-id` ties together all work related to one business flow.

Example checkout flow:

```text
api-gateway
-> order-service
-> inventory-service
-> payment-service
-> shipping-service
```

All logs and metrics for this flow should include the same correlation ID.

### Behaviour

At the edge:

```text
If x-correlation-id is present, preserve it.
If x-correlation-id is absent, generate one.
```

On outbound service calls:

```text
Propagate the current x-correlation-id.
```

### Example

```go
ctx = metadata.AppendToOutgoingContext(
    ctx,
    grpcmetadata.CorrelationID, correlationID,
)
```

---

## `x-request-id`

### Purpose

`x-request-id` identifies a single RPC/request instance.

This differs from `x-correlation-id`.

```text
x-correlation-id = one business flow
x-request-id     = one specific request/RPC attempt
```

A checkout flow may have one correlation ID but many request IDs.

### Behaviour

Each service boundary may generate a request ID if one is not present.

Request IDs are useful for debugging:

```text
retries
duplicate requests
load balancer behaviour
individual RPC attempts
```

---

## `authorization`

### Purpose

`authorization` carries auth credentials, commonly a bearer token.

Example:

```text
authorization: Bearer <token>
```

### Rules

Do not log raw authorisation metadata.

Do not propagate user tokens to downstream services unless the downstream service genuinely needs them.

Prefer service identity or mTLS for internal service-to-service authentication in later phases.

For local development, authentication may be disabled or simplified, but the metadata convention should still be documented.

---

## `traceparent`

### Purpose

`traceparent` carries W3C trace context and is commonly used by OpenTelemetry.

This allows distributed traces to connect work across services.

Example:

```text
api-gateway
-> catalog-service
-> database query
```

### Rules

Prefer allowing OpenTelemetry/gRPC instrumentation to manage trace propagation.

Do not hand-roll tracing semantics unless necessary.

---

## `x-bfstore-client`

### Purpose

`x-bfstore-client` identifies the logical caller.

Examples:

```text
api-gateway
order-service
admin-cli
integration-test
```

This is useful in logs and diagnostics.

Example:

```go
ctx = metadata.AppendToOutgoingContext(
    ctx,
    grpcmetadata.BFStoreClient, "order-service",
)
```

### Rules

This is a diagnostic helper, not an authentication mechanism.

Do not trust `x-bfstore-client` for security decisions.

---

## Metadata Security Rules

Metadata can appear in:

```text
logs
traces
proxies
debug tools
grpcurl output
test fixtures
```

Treat metadata as part of the security surface.

Do not put these in metadata:

```text
raw card numbers
CVV
passwords
large personal data
shipping addresses
full customer profiles
large debug payloads
```

Auth tokens may appear in metadata, but they must not be logged.

---

## Metadata Size Rules

Keep metadata small.

Good:

```text
x-correlation-id
x-request-id
authorization
traceparent
x-bfstore-client
```

Bad:

```text
x-basket-json
x-product-image
x-debug-dump
x-order-snapshot
```

If metadata starts looking like a document, it belongs somewhere else.

---

## Headers and Trailers

Use request metadata for:

```text
correlation IDs
request IDs
auth credentials
trace context
caller identification
```

Use response headers sparingly for small service information:

```text
x-bfstore-service
x-bfstore-service-version
x-correlation-id
```

Use trailers sparingly for diagnostics:

```text
x-query-cost
x-cache-status
x-rate-limit-remaining
x-bfstore-debug-id
```

Business logic must not depend heavily on custom response headers or trailers.

---

## Interceptor Usage

Metadata handling should usually live in interceptors rather than business handlers.

Recommended pattern:

```text
client interceptor:
adds x-correlation-id and x-bfstore-client to outgoing metadata

server interceptor:
reads x-correlation-id from incoming metadata
generates one if absent
adds it to context/logs

logging interceptor:
logs method, status, duration, correlation ID

auth interceptor:
reads authorization metadata
```

This keeps handlers focused on:

```text
request validation
domain orchestration
response mapping
error mapping
```

---

## Example Checkout Flow

```text
api-gateway receives checkout request
api-gateway sets x-correlation-id: checkout-abc-123

api-gateway -> order-service
order-service logs correlation_id=checkout-abc-123

order-service -> inventory-service
inventory-service logs correlation_id=checkout-abc-123

order-service -> payment-service
payment-service logs correlation_id=checkout-abc-123

order-service -> shipping-service
shipping-service logs correlation_id=checkout-abc-123
```

When checkout fails, the same correlation ID allows the full path to be traced across services.

Without metadata, debugging becomes guesswork.

---

## Review Checklist

Before adding a new metadata key, ask:

```text
Is this call context rather than business data?
Is it small?
Is it safe to appear in logs/traces/debug tools?
Does it need to be propagated across services?
Could it be represented better in the Protobuf message?
Is the key name lowercase?
Does it avoid the grpc- prefix?
```

If the answer is unclear, prefer explicit Protobuf fields for business data and keep metadata for platform concerns.

---

## Final Rule

```text
Metadata is the envelope note, not the parcel.
```

For bfstore, metadata should help services route, trace, authenticate, observe, and diagnose calls. The actual business contents belong in typed Protobuf messages.
