# Log Redaction

This document defines sensitive data and redaction rules for bfstore logs.

---

## Purpose

This document explains:

```text
sensitive data rules
redaction
payment/customer data handling
token/secret prevention
safe error logging
```

---

## Core Rule

```text
Logs are operational evidence, not a dumping ground.
```

Logs are copied, indexed, searched, retained, and often accessed by more people than production databases.

---

## Never Log

Do not log:

```text
passwords
API keys
JWTs
session tokens
authorisation headers
card numbers
CVV
raw payment provider payloads
full addresses unless genuinely needed
full email contents
secret DSNs
private keys
raw credentials
```

---

## Prefer Safe Identifiers

Prefer:

```text
customer_id instead of full customer details
order_id instead of full order payload
payment_attempt_id instead of payment details
masked email if needed
safe error_code instead of raw provider body
```

---

## Payment Logging

Allowed examples:

```text
payment_attempt_id
provider name
safe provider reference
status
error_code
duration_ms
idempotency_key hash/reference
```

Do not log:

```text
card number
CVV
full provider request/response
authorisation headers
secret API keys
raw tokens
```

---

## Email / Notification Logging

Allowed examples:

```text
notification_type
order_id
event_id
delivery_status
provider_message_id
masked email if genuinely needed
```

Avoid:

```text
full email address unless required
email body
customer private data
```

---

## Redaction Strategy

Use:

```text
structured logging helpers
safe error types
allowlisted fields
middleware/interceptors to prevent unsafe headers
tests for known sensitive field names
```

Prefer allowlisting safe fields over trying to remove unsafe fields after the fact.

---

## Practical Rules

```text
Never log secrets.
Never log raw payment details.
Never log full sensitive payloads.
Prefer IDs and safe codes.
Mask values only when necessary.
Use allowlisted fields.
Review logs as part of security hygiene.
```

---

## Final Rule

```text
A useful log should not become a data breach.
```
