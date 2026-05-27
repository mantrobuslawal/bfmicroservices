# ADR 0012: Use explicit proto3 field presence for patch and nullable scalar fields

## Status

Accepted

## Context

bfstore uses Protocol Buffers version 3 for gRPC APIs and Kafka event payloads.

In proto3, scalar fields such as `string`, `int64`, `bool`, and enums use implicit presence by default. This means a service cannot always distinguish between:

```text
field omitted
field explicitly set to the default value
```

Examples:

```text
string omitted      vs string explicitly set to ""
int64 omitted       vs int64 explicitly set to 0
bool omitted        vs bool explicitly set to false
enum omitted        vs enum explicitly set to UNKNOWN / UNSPECIFIED
```

This distinction matters for:

```text
patch requests
partial update commands
nullable business values
optional event context
failure detail fields
filter requests where omitted means "do not filter"
```

It matters less for:

```text
normal response messages
event facts that should always be populated
metadata fields required by validation
command identifiers such as product_id or order_id
```

## Decision

Use `optional` for proto3 scalar fields when bfstore must distinguish between:

```text
not provided
provided with default value
provided with non-default value
```

Use `optional` especially for:

```text
update request fields
patch request fields
nullable scalar business values
optional event context fields
optional failure detail fields
selected list/search filters where absence has meaning
```

Do not use `optional` blindly for every scalar field.

Keep normal scalar fields for values that are conceptually required and validated by application code.

Message fields already have presence in proto3, so `optional` is usually unnecessary for message fields.

## Examples

### Patch request

```proto
message UpdateProductRequest {
  string product_id = 1;

  optional string name = 2;
  optional string description = 3;
  optional int64 base_price_minor = 4;
  optional bool active = 5;
}
```

This allows the service to distinguish between:

```text
description omitted
description intentionally cleared to ""
```

and:

```text
active omitted
active intentionally set to false
```

### Event context

```proto
message EventMetadata {
  string event_id = 1;
  string event_type = 2;
  string event_version = 3;
  google.protobuf.Timestamp occurred_at = 4;
  string producer = 5;
  string subject = 6;

  optional string correlation_id = 7;
  optional string causation_id = 8;
  optional string trace_id = 9;
  optional string idempotency_key = 10;
}
```

The required-by-convention fields remain plain scalars and are validated by service code.

Optional context fields use explicit presence because some producers may not always have them.

### Failure detail

```proto
message OrderFailed {
  string order_id = 1;
  OrderFailureReason reason = 2;
  optional string failure_message = 3;
}
```

The event still records the required failure fact, while allowing a human-readable message to be absent or intentionally empty.

## Consequences

### Positive

This decision improves correctness for partial updates and optional business fields.

It avoids accidental behaviour where omitted scalar values are interpreted as intentional default values.

It makes Go handlers clearer because generated code can check field presence for optional scalar fields.

It documents a repeatable Protobuf convention across bfstore services.

### Negative

It adds some complexity to generated code because optional scalar fields are represented differently from plain scalar fields.

Engineers must decide deliberately whether field presence matters.

Reviewers must check that `optional` is used consistently and not sprayed everywhere.

## Validation

Because proto3 does not enforce required fields, services must still validate required-by-convention values.

Examples:

```text
product_id must be present and non-empty
event_id must be present and non-empty
occurred_at must be present
currency_code must be valid ISO 4217
amount_minor must be valid for the business rule
```

Validation belongs in service logic, interceptors, or dedicated validation helpers.

## Related Decisions

```text
ADR 0003: Use Kafka for events
ADR 0006: Use Buf for Protobuf
ADR 0011: Use outbox pattern for critical events
```
