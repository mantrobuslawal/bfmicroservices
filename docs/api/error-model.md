# API Error Model

## 1. Purpose

This document defines the API error model for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how services should represent, classify, log, propagate, and test errors across gRPC APIs and client-facing API Gateway responses.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s API design, reliability, and operational maturity.

---

## 2. Scope

This document covers:

```text
gRPC error categories
business error mapping
validation errors
authentication and authorisation errors
dependency failures
timeouts
idempotency conflicts
API Gateway error mapping
safe error responses
logging and observability
testing expectations
```

It does not define service-specific business rules in full. Those should be documented in:

```text
docs/requirements/service-requirements/
docs/requirements/business-rules.md
```

---

## 3. Error Model Goals

bfstore errors should be:

| Goal | Description |
|---|---|
| Consistent | Similar failures produce similar error categories |
| Safe | Sensitive implementation details are not exposed |
| Actionable | Clients and operators can understand what happened |
| Observable | Errors can be traced, logged, counted, and alerted |
| Business-aware | Expected business failures are distinct from system failures |
| Testable | Error mappings are covered by contract and integration tests |
| Client-friendly | External responses are clear without leaking internals |

---

## 4. Error Handling Principles

## 4.1 Errors Are Part of the Contract

Error behaviour is part of the API contract.

A service should document:

```text
which errors can be returned
what they mean
whether retry is safe
whether the error is caused by the caller
whether the error is a business rule failure
```

---

## 4.2 Expected Business Failures Are Not System Failures

Examples of expected business failures:

```text
product not found
product inactive
basket empty
insufficient stock
payment declined
invalid delivery option
order already cancelled
```

These should not be logged as severe platform incidents unless volume or pattern suggests a problem.

---

## 4.3 Internal Failures Must Not Leak Internals

Do not expose:

```text
SQL errors
stack traces
secret names
database hostnames
provider credentials
raw payment failure details
internal table names
```

External clients should receive safe messages.

Operators should get detailed diagnostics through logs, metrics, and traces.

---

## 4.4 Error Responses Must Be Correlatable

Every error should be traceable through:

```text
correlation_id
request_id
trace_id
service name
operation name
business entity ID where safe
```

A customer-facing error should include a reference that support or operators can use to find the internal trace.

---

## 5. gRPC Error Categories

bfstore services should use standard gRPC status codes consistently.

| gRPC Code | Use When |
|---|---|
| `OK` | Request completed successfully |
| `INVALID_ARGUMENT` | Request shape or field value is invalid |
| `NOT_FOUND` | Requested entity does not exist or is not visible to caller |
| `ALREADY_EXISTS` | Entity already exists where uniqueness is required |
| `FAILED_PRECONDITION` | Request is valid but current business state prevents action |
| `ABORTED` | Operation was aborted due to concurrency or transaction conflict |
| `OUT_OF_RANGE` | Value is outside an allowed range |
| `UNAUTHENTICATED` | Caller is not authenticated |
| `PERMISSION_DENIED` | Caller is authenticated but not authorised |
| `RESOURCE_EXHAUSTED` | Quota, rate limit, or capacity limit reached |
| `CANCELLED` | Request was cancelled by caller |
| `DEADLINE_EXCEEDED` | Deadline or timeout exceeded |
| `UNAVAILABLE` | Dependency or service temporarily unavailable |
| `INTERNAL` | Unexpected internal failure |
| `UNKNOWN` | Avoid where possible; only when error cannot be classified |

---

## 6. Business Error Mapping

## 6.1 Catalogue Errors

| Scenario | gRPC Code | Retry? |
|---|---|---|
| Product ID missing | `INVALID_ARGUMENT` | No |
| Product not found | `NOT_FOUND` | No |
| Product inactive | `FAILED_PRECONDITION` | No |
| Category not found | `NOT_FOUND` | No |
| Catalogue database unavailable | `UNAVAILABLE` | Maybe |
| Unexpected catalogue failure | `INTERNAL` | Maybe |

---

## 6.2 Basket Errors

| Scenario | gRPC Code | Retry? |
|---|---|---|
| Basket not found | `NOT_FOUND` | No |
| Basket item quantity invalid | `INVALID_ARGUMENT` | No |
| Product cannot be added to basket | `FAILED_PRECONDITION` | No |
| Basket already checked out | `FAILED_PRECONDITION` | No |
| Basket update conflict | `ABORTED` | Maybe |
| Basket database unavailable | `UNAVAILABLE` | Maybe |

---

## 6.3 Inventory Errors

| Scenario | gRPC Code | Retry? |
|---|---|---|
| Product or variant missing | `INVALID_ARGUMENT` or `NOT_FOUND` | No |
| Quantity invalid | `INVALID_ARGUMENT` | No |
| Insufficient stock | `FAILED_PRECONDITION` | No |
| Reservation expired | `FAILED_PRECONDITION` | No |
| Duplicate reservation request with same idempotency key | `OK` with previous result or `ALREADY_EXISTS` depending on API design | Safe if idempotent |
| Reservation conflict | `ABORTED` | Maybe |
| Inventory database unavailable | `UNAVAILABLE` | Maybe |

---

## 6.4 Order Errors

| Scenario | gRPC Code | Retry? |
|---|---|---|
| Checkout request invalid | `INVALID_ARGUMENT` | No |
| Basket empty | `FAILED_PRECONDITION` | No |
| Basket already checked out | `FAILED_PRECONDITION` | No |
| Insufficient stock | `FAILED_PRECONDITION` | No |
| Payment declined | `FAILED_PRECONDITION` | No |
| Duplicate checkout idempotency key | `OK` with previous result | Safe |
| Order not found | `NOT_FOUND` | No |
| Order cancellation not allowed | `FAILED_PRECONDITION` | No |
| Downstream service timeout | `DEADLINE_EXCEEDED` | Maybe |
| Downstream service unavailable | `UNAVAILABLE` | Maybe |
| Unexpected order failure | `INTERNAL` | Maybe |

---

## 6.5 Payment Errors

| Scenario | gRPC Code | Retry? |
|---|---|---|
| Payment request invalid | `INVALID_ARGUMENT` | No |
| Payment method invalid | `FAILED_PRECONDITION` | No |
| Payment declined | `FAILED_PRECONDITION` | No |
| Duplicate payment idempotency key | `OK` with previous result | Safe |
| Payment provider timeout | `DEADLINE_EXCEEDED` | Maybe, only with idempotency |
| Payment provider unavailable | `UNAVAILABLE` | Maybe, only with idempotency |
| Payment provider returned unexpected response | `INTERNAL` | Maybe |
| Raw payment data supplied where not allowed | `INVALID_ARGUMENT` | No |

---

## 6.6 Shipping Errors

| Scenario | gRPC Code | Retry? |
|---|---|---|
| Delivery option invalid | `INVALID_ARGUMENT` | No |
| Address invalid | `INVALID_ARGUMENT` | No |
| Shipment not found | `NOT_FOUND` | No |
| Shipment already exists for idempotency key | `OK` with previous result | Safe |
| Carrier unavailable | `UNAVAILABLE` | Maybe |
| Carrier timeout | `DEADLINE_EXCEEDED` | Maybe |
| Shipment cannot be cancelled | `FAILED_PRECONDITION` | No |

---

## 6.7 Notification Errors

| Scenario | gRPC Code | Retry? |
|---|---|---|
| Notification request invalid | `INVALID_ARGUMENT` | No |
| Template not found | `NOT_FOUND` | No |
| Unsupported notification channel | `FAILED_PRECONDITION` | No |
| Duplicate notification event | `OK` or ignored idempotently | Safe |
| Provider unavailable | `UNAVAILABLE` | Maybe |
| Provider timeout | `DEADLINE_EXCEEDED` | Maybe |
| Notification permanently rejected | `FAILED_PRECONDITION` | No |

---

## 7. Validation Errors

Invalid input should return `INVALID_ARGUMENT`.

Examples:

```text
missing required product_id
quantity less than 1
invalid currency code
invalid page size
malformed customer_id
unsupported enum value
```

Validation errors should be specific enough for developers to fix the request, but should not expose internal implementation details.

Example safe message:

```text
quantity must be greater than zero
```

Poor message:

```text
SQL constraint basket_items_quantity_check failed
```

---

## 8. Business Rule Errors

Business rule failures usually map to `FAILED_PRECONDITION`.

Examples:

```text
product is inactive
basket has already been checked out
stock is insufficient
payment was declined
order cannot be cancelled after dispatch
review cannot be submitted before delivery
```

These are not system outages. They are valid business outcomes.

---

## 9. Authentication and Authorisation Errors

## 9.1 Unauthenticated

Use `UNAUTHENTICATED` when the caller has not provided valid authentication.

Examples:

```text
missing token
expired token
invalid token signature
unsupported issuer
```

## 9.2 Permission Denied

Use `PERMISSION_DENIED` when the caller is authenticated but not allowed to perform the action.

Examples:

```text
customer attempts to view another customer's order
user lacks admin role
service account lacks permission for operation
```

Do not reveal whether protected resources exist when doing so would leak information.

---

## 10. Timeouts

Use `DEADLINE_EXCEEDED` when an operation exceeds its deadline.

Examples:

```text
payment provider timeout
inventory reservation timeout
shipping provider timeout
database query timeout
downstream gRPC deadline exceeded
```

Timeout logs should include:

```text
dependency
operation
configured timeout
elapsed time
correlation_id
trace_id
```

---

## 11. Unavailable Dependencies

Use `UNAVAILABLE` when a dependency is temporarily unavailable.

Examples:

```text
MySQL unavailable
Kafka unavailable
downstream service unavailable
provider unavailable
DNS failure
connection refused
```

`UNAVAILABLE` may be retryable, but only where the operation is safe to retry.

---

## 12. Idempotency Error Behaviour

Idempotent operations should behave consistently for duplicate requests.

Preferred behaviour:

```text
same idempotency key + same request body -> return original result
same idempotency key + different request body -> reject as conflict
```

Possible mapping:

| Scenario | gRPC Code |
|---|---|
| Same key, same request | `OK` |
| Same key, different request | `ALREADY_EXISTS` or `FAILED_PRECONDITION` |
| Idempotency key missing for required operation | `INVALID_ARGUMENT` |

Operations requiring idempotency:

```text
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
RefundPayment
SendNotification
```

---

## 13. API Gateway Error Mapping

The API Gateway should translate internal gRPC errors to safe client-facing HTTP responses if the external API is REST/JSON.

Example mapping:

| gRPC Code | HTTP Status |
|---|---:|
| `OK` | `200`, `201`, or `204` |
| `INVALID_ARGUMENT` | `400` |
| `UNAUTHENTICATED` | `401` |
| `PERMISSION_DENIED` | `403` |
| `NOT_FOUND` | `404` |
| `ALREADY_EXISTS` | `409` |
| `FAILED_PRECONDITION` | `409` or `422` |
| `ABORTED` | `409` |
| `RESOURCE_EXHAUSTED` | `429` |
| `CANCELLED` | `499` where supported, otherwise `400` or `499` equivalent |
| `DEADLINE_EXCEEDED` | `504` |
| `UNAVAILABLE` | `503` |
| `INTERNAL` | `500` |
| `UNKNOWN` | `500` |

The chosen external status mapping should be documented and used consistently.

---

## 14. Client-Facing Error Shape

If the API Gateway exposes REST/JSON, use a consistent error shape.

Example:

```json
{
  "error": {
    "code": "INSUFFICIENT_STOCK",
    "message": "One or more items are no longer available in the requested quantity.",
    "correlation_id": "corr_01HX...",
    "details": [
      {
        "field": "items[0].quantity",
        "reason": "requested quantity exceeds available stock"
      }
    ]
  }
}
```

## 14.1 External Error Fields

| Field | Description |
|---|---|
| `code` | Stable application error code |
| `message` | Safe human-readable message |
| `correlation_id` | Support and troubleshooting reference |
| `details` | Optional structured field or business details |

## 14.2 Message Rules

Client-facing messages should be:

```text
clear
safe
stable enough for clients
free of internal stack traces
free of database details
free of secret values
```

---

## 15. Application Error Codes

Use stable application-level error codes where useful.

Examples:

```text
PRODUCT_NOT_FOUND
PRODUCT_INACTIVE
BASKET_EMPTY
BASKET_ALREADY_CHECKED_OUT
INSUFFICIENT_STOCK
PAYMENT_DECLINED
PAYMENT_PROVIDER_UNAVAILABLE
DELIVERY_OPTION_INVALID
ORDER_NOT_FOUND
ORDER_CANNOT_BE_CANCELLED
NOTIFICATION_TEMPLATE_NOT_FOUND
```

Application error codes should not replace gRPC status codes. They add business-specific detail.

---

## 16. Logging Errors

## 16.1 Log Levels

Recommended log levels:

| Level | Use |
|---|---|
| `DEBUG` | Developer diagnostics, disabled or reduced in production |
| `INFO` | Important normal business events |
| `WARN` | Expected failures or recoverable issues needing attention |
| `ERROR` | Unexpected failures or failed operations |
| `FATAL` | Service cannot continue |

## 16.2 Business Errors

Expected business failures may be `INFO` or `WARN` depending on context.

Examples:

```text
insufficient stock
payment declined
basket already checked out
```

These are not automatically platform incidents.

## 16.3 System Errors

Unexpected system failures should be `ERROR`.

Examples:

```text
database unavailable
Kafka publish failure
unexpected nil reference
provider integration failure
migration failure
```

## 16.4 Sensitive Logging Rules

Do not log:

```text
passwords
tokens
raw card data
secret values
full customer addresses
provider credentials
```

Prefer:

```text
order_id
payment_id
customer_id
correlation_id
trace_id
event_id
```

---

## 17. Metrics

Services should expose error metrics.

Examples:

```text
grpc_server_requests_total
grpc_server_errors_total
grpc_client_errors_total
checkout_failures_total
payment_declined_total
stock_reservation_failures_total
kafka_publish_failures_total
kafka_consumer_failures_total
database_errors_total
```

Metrics should include useful labels such as:

```text
service
operation
error_code
grpc_code
dependency
```

Avoid high-cardinality labels such as raw IDs.

---

## 18. Tracing

Errors should be attached to distributed traces.

Trace spans should include:

```text
service name
operation name
status
error category
dependency
correlation_id
business entity ID where safe
```

A failed checkout should be traceable across:

```text
api-gateway
order-service
basket-service
inventory-service
payment-service
shipping-service
notification-service
```

---

## 19. Kafka Error Handling

Kafka consumers should classify failures.

| Failure | Handling |
|---|---|
| Temporary database error | Retry |
| Temporary downstream error | Retry |
| Invalid event payload | DLQ |
| Unsupported event version | DLQ |
| Duplicate event | Ignore or idempotently handle |
| Poison message | DLQ and alert |

Kafka consumer error handling should include:

```text
retry count
DLQ count
consumer lag
event_id
event_type
correlation_id
failure reason
```

---

## 20. Security Considerations

Error handling must not weaken security.

Avoid:

```text
revealing whether protected resources exist
returning detailed auth failure reasons to attackers
leaking token validation internals
exposing provider credentials
including stack traces in external responses
```

Example safe authentication error:

```text
authentication required
```

Poor authentication error:

```text
JWT signature validation failed using key kid=prod-key-001
```

---

## 21. Testing Requirements

## 21.1 Unit Tests

Unit tests should cover:

```text
business error classification
validation error generation
idempotency conflict handling
error-to-status-code mapping
safe message generation
```

## 21.2 Contract Tests

Contract tests should verify:

```text
expected gRPC status codes
expected application error codes
safe error details
backward-compatible error behaviour
```

## 21.3 Integration Tests

Integration tests should cover:

```text
database unavailable
downstream service unavailable
timeout behaviour
Kafka publish failure
duplicate request behaviour
```

## 21.4 End-to-End Tests

Critical E2E error cases:

```text
checkout fails for insufficient stock
checkout fails for payment declined
duplicate checkout request returns original result
notification failure does not roll back order
product not found returns safe response
```

---

## 22. Initial Error Handling Scope

The first implementation should define and test errors for:

```text
product not found
product inactive
basket empty
invalid quantity
insufficient stock
payment declined
payment timeout
shipment creation failure
duplicate checkout request
OrderCreated notification duplicate
database unavailable
Kafka publish failure
```

---

## 23. Anti-Patterns to Avoid

Avoid:

```text
returning OK with success=false for failed RPCs
exposing SQL errors to clients
using UNKNOWN for most failures
logging secrets or raw payment data
treating payment declined as system outage
retrying non-idempotent operations blindly
mapping every business failure to INTERNAL
using inconsistent error codes across services
```

---

## 24. Error Review Checklist

Before approving an API change, check:

```text
Are expected errors documented?
Are gRPC status codes appropriate?
Are business errors distinct from system errors?
Are client-facing messages safe?
Are application error codes stable?
Are errors logged with correlation IDs?
Are sensitive values excluded from logs?
Are timeout and retry behaviours clear?
Are idempotency conflicts handled?
Are contract tests updated?
```

---

## 25. Related Documents

This document should be read alongside:

```text
docs/api/grpc-overview.md
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/architecture/communication-patterns.md
docs/architecture/resilience-patterns.md
docs/events/retry-and-dlq-strategy.md
docs/testing/testing-strategy.md
docs/security/secure-coding.md
docs/observability/logging.md
docs/observability/tracing.md
```

Relevant ADRs:

```text
adr/0002-use-grpc-for-service-communication.md
adr/0006-use-buf-for-protobuf.md
adr/0008-use-contract-first-service-design.md
```

---

## 26. Summary

bfstore’s error model is designed to make failures predictable, safe, observable, and testable.

The most important rules are:

```text
use standard gRPC status codes
distinguish business failures from system failures
do not leak internal implementation details
include correlation IDs
support idempotent retry where needed
map API Gateway errors consistently
test error behaviour as part of the contract
```

A strong error model improves developer experience, operational support, security posture, and client confidence.
