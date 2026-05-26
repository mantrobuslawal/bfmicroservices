# `db/notification/migrations`

## 1. Purpose

This directory contains database migrations for the Notification Service schema.

Notification Service owns notification records, notification attempts, processed event IDs, delivery state, and notification templates where implemented.

---

## 2. Owning Service

```text
notification-service
```

Owned schema:

```text
bfstore_notification
```

Only Notification Service migrations should modify this schema.

---

## 3. Expected Migration Files

Recommended initial migrations:

```text
db/notification/migrations/
├── README.md
├── 000001_create_notifications.up.sql
├── 000001_create_notifications.down.sql
├── 000002_create_notification_attempts.up.sql
├── 000002_create_notification_attempts.down.sql
├── 000003_create_processed_events.up.sql
├── 000003_create_processed_events.down.sql
├── 000004_create_notification_templates.up.sql
└── 000005_create_outbox_events.up.sql
```

Templates and outbox may be deferred depending on implementation phase.

---

## 4. Candidate Tables

Initial priority:

```text
notifications
notification_attempts
processed_events
```

Later:

```text
notification_templates
outbox_events
notification_preferences
```

---

## 5. Core Data Ownership

Notification Service owns:

```text
notification_id
notification status
notification type
notification channel
delivery attempts
processed event IDs
notification provider response summaries
template references
```

Notification Service does not own:

```text
order lifecycle
customer profile truth
payment state
shipment state
```

Notification may store references such as `order_id`, `customer_id`, and `event_id`.

---

## 6. Initial Table Design Notes

### `notifications`

Recommended fields:

```text
notification_id
event_id
order_id
customer_id
notification_type
channel
status
template_id
recipient_reference
created_at
updated_at
sent_at
failed_at
```

### `notification_attempts`

Recommended fields:

```text
notification_attempt_id
notification_id
provider
status
failure_reason
attempt_number
created_at
completed_at
```

### `processed_events`

Recommended fields:

```text
event_id
event_type
producer
consumer
processed_at
processing_status
failure_reason
```

### `notification_templates`

Recommended fields:

```text
template_id
template_key
channel
subject_template
body_template
status
version
created_at
updated_at
```

---

## 7. Indexing Guidance

Recommended indexes:

```text
idx_notifications_event_id
idx_notifications_order_id
idx_notifications_customer_id
idx_notifications_status
idx_notification_attempts_notification_id
uq_processed_events_event_id_consumer
idx_processed_events_event_type
idx_templates_template_key_version
```

---

## 8. Constraints and Invariants

Notification migrations should support:

```text
notification_id is unique
event_id is tracked for deduplication
attempt_number is positive
notification status is present
channel is present
processed event uniqueness prevents duplicate side effects
```

---

## 9. Idempotency Requirements

Notification consumers must not send duplicate customer notifications for the same event.

Database support should include:

```text
processed_events table
unique event/consumer constraint
notification event_id reference
attempt tracking
```

Expected behaviour:

```text
same OrderCreated event processed twice does not send duplicate confirmation
manual replay can be controlled safely
```

---

## 10. Event and Outbox Considerations

Notification Service should publish:

```text
NotificationSent
NotificationFailed
NotificationSuppressed
```

If reliable publication of notification events becomes important, add:

```text
outbox_events
```

---

## 11. Privacy and Security Rules

Notification data may include contact references.

Rules:

```text
minimise PII
do not store full email bodies with sensitive content unless justified
do not log full recipient details unnecessarily
do not store provider secrets
do not store tokens in notification records
```

Where possible, store recipient references rather than full personal details.

---

## 12. Migration Safety Rules

```text
do not create foreign keys to order or customer schemas
do not store payment data
do not remove processed event deduplication
do not store provider secrets
do not use real customer contact data in seed data
do not edit applied migrations
```

---

## 13. Testing Expectations

Notification migrations should be validated by tests for:

```text
migrations apply cleanly
notification can be created
notification attempt can be recorded
processed event uniqueness works
duplicate event is detected
notification status is updated
sensitive values are not required by schema
```

---

## 14. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/events/ordering-and-idempotency.md
docs/events/retry-and-dlq-strategy.md
proto/acme/notification/v1/README.md
```

---

## 15. Summary

Notification migrations define asynchronous customer communication persistence.

They must support idempotent event processing, attempt tracking, safe retry behaviour, and careful handling of personal data.
