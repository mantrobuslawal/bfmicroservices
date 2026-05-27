# Protobuf Field Presence Convention

## 1. Purpose

This document defines how bfstore uses proto3 field presence.

bfstore uses Protocol Buffers for:

```text
gRPC service contracts
Kafka event payloads
shared service messages
```

The goal is to avoid ambiguity around scalar default values.

---

## 2. The proto3 Problem

In proto3, plain scalar fields do not track explicit presence.

For example:

```proto
message Example {
  string name = 1;
  int64 amount_minor = 2;
  bool active = 3;
}
```

The receiver cannot reliably distinguish:

```text
name omitted        vs name set to ""
amount omitted      vs amount set to 0
active omitted      vs active set to false
```

This can be dangerous for update and patch APIs.

---

## 3. Decision

Use `optional` for scalar fields when bfstore must distinguish:

```text
omitted
explicitly set to default value
explicitly set to non-default value
```

Do not use `optional` blindly.

---

## 4. Use Optional For

Use `optional` for:

```text
patch request fields
partial update request fields
nullable scalar business values
optional failure details
optional event context
selected list/search filters where absence has meaning
```

Examples:

```proto
message UpdateProductRequest {
  string product_id = 1;

  optional string name = 2;
  optional string description = 3;
  optional int64 base_price_minor = 4;
  optional bool active = 5;
}
```

```proto
message OrderFailed {
  string order_id = 1;
  OrderFailureReason reason = 2;
  optional string failure_message = 3;
}
```

---

## 5. Do Not Usually Use Optional For

Avoid unnecessary `optional` fields for:

```text
normal response fields
core event facts
required-by-validation command identifiers
fields inside value objects that are always validated as a unit
repeated fields
message fields
```

Examples:

```proto
message Product {
  string product_id = 1;
  string name = 2;
  string category_id = 3;
}
```

```proto
message Money {
  int64 amount_minor = 1;
  string currency_code = 2;
}
```

---

## 6. Message Fields

Message fields already have presence in proto3.

Example:

```proto
message Product {
  acme.common.v1.Money base_price = 1;
}
```

The service can tell whether `base_price` was supplied.

Adding `optional` to message fields is usually unnecessary.

---

## 7. Events

Events are facts.

Most event fields should be plain scalar fields because they should always be populated by the producer.

Use `optional` only for genuinely conditional values.

Example:

```proto
message OrderFailed {
  string order_id = 1;
  OrderFailureReason reason = 2;
  optional string failure_message = 3;
}
```

---

## 8. Filters

Filters require judgement.

This is acceptable:

```proto
message ListProductsRequest {
  optional string category_id = 1;
  optional string search_query = 2;
  optional bool include_inactive = 3;
}
```

This lets the service distinguish:

```text
filter omitted
filter explicitly supplied as empty/default
```

If omitted and default values mean the same thing, plain scalars are also acceptable.

Document the behaviour.

---

## 9. Go Implementation Guidance

Generated Go code represents optional scalar fields differently from plain scalar fields.

Handlers should convert Protobuf requests into domain commands.

Recommended boundary:

```text
Protobuf request
→ gRPC handler
→ domain command/query
→ service layer
→ repository
```

Avoid passing generated Protobuf request messages into repositories.

---

## 10. Validation

`optional` does not replace validation.

Services must still validate:

```text
required identifiers
valid statuses
allowed enum values
valid money values
allowed state transitions
required event metadata
business invariants
```

Examples:

```text
product_id must not be empty
currency_code must be valid
event_id must not be empty
occurred_at must be present
OrderFailed reason must not be UNSPECIFIED
```

---

## 11. Review Checklist

When reviewing a `.proto` change, ask:

```text
Does this field need to distinguish omitted from default?
Is this an update/patch request?
Is this a nullable business value?
Is this a core event fact that should always be present?
Would optional make the handler clearer or noisier?
Does validation still enforce required-by-convention fields?
```

---

## 12. Summary

Use `optional` deliberately.

The bfstore rule is:

```text
optional when presence matters
plain scalar when the value is required-by-convention or always populated
message fields already have presence
validation still belongs in service logic
```
