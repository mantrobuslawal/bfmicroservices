# `notification/v1`

## 1. Purpose

This document describes the `notification/v1` Protobuf package for bfstore.

The package defines notification records, delivery attempts, notification status, templates, and notification event contracts.

These contracts are client-visible engineering artefacts and should be treated as stable integration boundaries.

---

## 2. Package Name

```proto
package acme.notification.v1;
```

Recommended Go package option:

```proto
option go_package = "github.com/acme-ltd/bfstore/gen/go/acme/notification/v1;notificationv1";
```

---

## 3. Ownership

| Area | Owner |
|---|---|
| Package owner | `notification-service` |
| Contract review | API governance / service owner |
| Backward compatibility | Service owner |
| Generated clients | Build/tooling pipeline |

Ownership rule:

> The service that owns the business capability owns the contract. Other services integrate through APIs and events, not database tables.

---

## 4. Expected Files

```text
proto/acme/notification/v1/
├── README.md
├── notification_service.proto
├── notification.proto
├── notification_events.proto
```

Event contracts may later move to a dedicated package such as:

```text
proto/acme/<domain>/events/v1/
```

if that separation improves clarity.

---

## 5. Primary Service Contract

```proto
service NotificationService {
  rpc GetNotification(... ) returns (...);
  rpc ListNotifications(... ) returns (...);
  rpc RetryNotification(... ) returns (...);
}
```

---

## 6. Core Message Types

Recommended message types:

- `Notification`
- `NotificationAttempt`
- `NotificationStatus`
- `NotificationChannel`
- `NotificationType`

These messages should describe business concepts and API contracts, not database rows.

---

## Recommended Enums

```proto
enum NotificationStatus {
  NOTIFICATION_STATUS_UNSPECIFIED = 0;
  NOTIFICATION_STATUS_PENDING = 1;
  NOTIFICATION_STATUS_SENT = 2;
  NOTIFICATION_STATUS_FAILED = 3;
  NOTIFICATION_STATUS_SUPPRESSED = 4;
}

enum NotificationChannel {
  NOTIFICATION_CHANNEL_UNSPECIFIED = 0;
  NOTIFICATION_CHANNEL_EMAIL = 1;
  NOTIFICATION_CHANNEL_SMS = 2;
  NOTIFICATION_CHANNEL_PUSH = 3;
}
```

---

## 7. Contract Rules

- Notification Service is primarily event-driven.
- OrderCreated must be consumed idempotently.
- Duplicate events must not create duplicate customer messages.
- Notification failure must not roll back order creation.
- Events should minimise customer PII.

---

## 8. Event Contracts

Expected or potential events:

- `NotificationRequested`
- `NotificationSent`
- `NotificationFailed`
- `NotificationSuppressed`

Event rules:

```text
events describe facts that have already happened
events must include the standard event envelope
events must be versioned
events must carry correlation context
consumers must be idempotent
```

---

## 9. Error Behaviour

Expected error behaviour:

- Duplicate event is ignored or handled idempotently
- Invalid event payload is retried or sent to DLQ
- Unsupported event version is safely rejected or sent to DLQ
- Provider unavailable is retried
- Notification not found returns NOT_FOUND for gRPC lookup

All gRPC errors should follow:

```text
docs/api/error-model.md
```

---

## 10. Versioning

This package is versioned as `v1`.

Compatible changes include:

```text
adding optional fields
adding new messages
adding new RPCs
adding new event types
adding comments
```

Breaking changes include:

```text
renaming services or RPCs
renaming packages
removing fields
renumbering fields
changing field types
changing field meaning
changing idempotency behaviour
changing error semantics
```

Breaking changes require a new package version, such as:

```proto
package acme.notification.v2;
```

Removed fields must be reserved.

---

## 11. Security and Privacy

Contracts must avoid unnecessary sensitive data.

Do not expose:

```text
passwords
tokens
raw payment card data
secret values
internal stack traces
internal database IDs
unnecessary customer PII
```

Prefer opaque IDs and service-owned lookups where sensitive details are required.

---

## 12. Testing Expectations

This package should be covered by:

```text
buf lint
buf breaking
protobuf generation checks
gRPC contract tests
event contract tests where events are defined
integration tests for critical service behaviours
```

---

## 13. Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/error-model.md
docs/api/versioning.md
docs/architecture/communication-patterns.md
docs/architecture/service-boundaries.md
docs/events/event-catalog.md
docs/testing/testing-strategy.md
```

---

## 14. Summary

The `acme.notification.v1` package is part of bfstore's contract-first service design.

It should remain business-focused, versioned, testable, and aligned with the owning service boundary.
