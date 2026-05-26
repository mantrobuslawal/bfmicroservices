# `db/payment/migrations`

## 1. Purpose

This directory contains database migrations for the Payment Service schema.

Payment Service owns payment state, payment attempts, provider references, authorisation outcomes, and payment-related events.

---

## 2. Owning Service

```text
payment-service
```

Owned schema:

```text
bfstore_payment
```

Only Payment Service migrations should modify this schema.

---

## 3. Expected Migration Files

Recommended initial migrations:

```text
db/payment/migrations/
├── README.md
├── 000001_create_payments.up.sql
├── 000001_create_payments.down.sql
├── 000002_create_payment_attempts.up.sql
├── 000002_create_payment_attempts.down.sql
├── 000003_create_payment_status_history.up.sql
├── 000003_create_payment_status_history.down.sql
├── 000004_create_refunds.up.sql
└── 000005_create_outbox_events.up.sql
```

Refunds and outbox may be deferred depending on implementation phase.

---

## 4. Candidate Tables

Initial priority:

```text
payments
payment_attempts
```

Later:

```text
payment_status_history
refunds
outbox_events
```

---

## 5. Core Data Ownership

Payment Service owns:

```text
payment_id
payment status
payment amount
payment attempt records
provider name
provider reference
authorisation result
refund state where implemented
```

Payment Service does not own:

```text
order lifecycle
stock reservation
shipment state
customer profile
raw card data
```

Payment stores `order_id` as a reference only.

---

## 6. Initial Table Design Notes

### `payments`

Recommended fields:

```text
payment_id
order_id
customer_id
status
amount_minor
currency_code
provider
provider_reference
idempotency_key
request_hash
created_at
updated_at
authorised_at
captured_at
```

### `payment_attempts`

Recommended fields:

```text
payment_attempt_id
payment_id
order_id
attempt_type
status
failure_reason
provider_reference
created_at
completed_at
```

### `refunds`

Recommended fields:

```text
refund_id
payment_id
amount_minor
currency_code
status
idempotency_key
provider_reference
created_at
updated_at
```

---

## 7. Indexing Guidance

Recommended indexes:

```text
idx_payments_order_id
idx_payments_customer_id
idx_payments_status
uq_payments_idempotency_key
idx_payment_attempts_payment_id
idx_payment_attempts_order_id
idx_refunds_payment_id
uq_refunds_idempotency_key
```

---

## 8. Constraints and Invariants

Payment migrations should support:

```text
payment_id is unique
amount_minor >= 0
currency_code is present
idempotency_key is unique where required
payment status is present
attempt status is present
raw card data is not stored
```

---

## 9. Security Requirements

Payment migrations must not introduce columns for:

```text
raw_card_number
cvv
card_security_code
plaintext_token
secret
password
```

Acceptable fields include:

```text
provider
provider_reference
safe failure code
last four digits only if explicitly justified
```

Even provider references should be treated as sensitive operational data.

---

## 10. Event and Outbox Considerations

Payment Service should publish:

```text
PaymentAuthorised
PaymentFailed
PaymentCaptured
PaymentRefunded
```

Outbox is recommended when payment events become critical to downstream workflows.

---

## 11. Migration Safety Rules

```text
do not store raw card data
do not create foreign keys to order schema
do not remove idempotency columns without migration plan
do not use FLOAT or DOUBLE for money
do not expose provider secrets in seed data
do not edit applied migrations
```

---

## 12. Testing Expectations

Payment migrations should be validated by tests for:

```text
migrations apply cleanly
payment can be inserted
payment attempt can be inserted
idempotency uniqueness works
amount constraints work
raw card data columns do not exist
payment declined scenario can be recorded
```

---

## 13. Related Documents

```text
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/requirements/business-rules.md
docs/api/error-model.md
proto/acme/payment/v1/README.md
```

---

## 14. Summary

Payment migrations define a sensitive and business-critical data model.

They must support idempotency, auditability, safe provider references, and strict avoidance of raw payment data.
