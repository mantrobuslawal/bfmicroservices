# Common Protobuf Package

## Package

```proto
package acme.common.v1;
```

## Purpose

This package contains Protobuf messages shared across bfstore services.

Common messages should be small, stable, and broadly reusable.

Examples:

```text
Money
PageRequest
PageResponse
```

Do not place service-specific business concepts in this package.

---

## Field Presence Convention

bfstore uses proto3.

Use `optional` for scalar fields only when the service must distinguish between:

```text
field omitted
field explicitly set to the scalar default value
field explicitly set to a non-default value
```

This is most useful for:

```text
patch requests
partial update commands
nullable scalar business values
optional event context
selected filter fields where absence has meaning
```

Do not use `optional` blindly for every scalar field.

For normal response models and values required by validation, prefer plain scalar fields.

Message fields already have presence in proto3.

---

## Money

The `Money` message should represent monetary values using minor units.

Recommended shape:

```proto
message Money {
  int64 amount_minor = 1;
  string currency_code = 2;
}
```

Example:

```text
£45.99 = amount_minor: 4599, currency_code: "GBP"
```

### Money field presence

`Money` is a message type, so proto3 already supports presence for the field itself when used inside another message.

For example:

```proto
message Product {
  Money base_price = 1;
}
```

The service can distinguish whether `base_price` was set.

Inside `Money`, keep `amount_minor` and `currency_code` as plain scalar fields because a valid `Money` value should always include both values and be validated by application logic.

---

## Pagination

Pagination messages should use plain scalar fields unless presence has a specific meaning.

Recommended shape:

```proto
message PageRequest {
  int32 page_size = 1;
  string page_token = 2;
}

message PageResponse {
  string next_page_token = 1;
  int32 total_size = 2;
}
```

For many APIs, an omitted or empty `page_token` naturally means "first page".

A zero `page_size` can mean "use service default", but this should be documented and validated by the service.

---

## Validation

Protobuf does not enforce required fields in proto3.

Services must validate required-by-convention fields.

Examples:

```text
currency_code must not be empty
currency_code should be a valid ISO 4217 currency code
amount_minor must satisfy business rules
page_size must be within allowed limits
```

---

## Design Rule

Keep this package boring and stable.

Common types should not import service-specific packages or create hidden coupling between services.
