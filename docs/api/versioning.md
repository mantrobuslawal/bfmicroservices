# API Versioning

## 1. Purpose

This document defines the API versioning strategy for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how gRPC APIs, protobuf packages, Kafka event contracts, and client-facing API Gateway routes should evolve safely over time.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s contract design and long-term maintainability.

---

## 2. Scope

This document covers versioning for:

```text
protobuf packages
gRPC services
request and response messages
Kafka event payloads
Kafka event topics
API Gateway routes
application error codes
generated clients
contract tests
deprecation and migration
```

It does not define release versioning for container images or infrastructure modules. Those should be covered in deployment and release documentation.

---

## 3. Versioning Goals

bfstore versioning should support:

| Goal | Description |
|---|---|
| Safe evolution | APIs can change without breaking consumers unexpectedly |
| Clear compatibility | Teams know which changes are safe or breaking |
| Independent services | Services can evolve at different speeds |
| Contract confidence | Protobuf and event contracts are checked in CI |
| Migration support | Old and new clients can coexist during transition |
| Client clarity | External APIs expose stable, documented versions |
| Operational safety | Releases and rollbacks remain manageable |

---

## 4. Versioning Principles

## 4.1 Contracts Are Public Within the System

Even internal service APIs are contracts.

Once another service depends on a protobuf message, gRPC method, event type, or error code, it must be evolved carefully.

---

## 4.2 Prefer Additive Changes

Prefer backward-compatible additions over breaking changes.

Good:

```text
add optional field
add new RPC
add new event type
add new enum value with safe handling
```

Risky:

```text
remove field
rename field
change field meaning
change field type
change RPC semantics
```

---

## 4.3 Version at the Contract Boundary

Version the contract, not just the implementation.

Examples:

```proto
package acme.order.v1;
```

```text
bfstore.order.order-events.v1
```

A service implementation may change many times without changing contract version.

---

## 4.4 Breaking Changes Require a Migration Plan

Breaking changes must not be introduced casually.

They require:

```text
new version
consumer impact analysis
migration plan
deprecation period where relevant
contract tests
release notes
rollback consideration
```

---

## 5. Protobuf Package Versioning

Protobuf packages should include a major version.

Format:

```proto
package acme.<domain>.v1;
```

Examples:

```proto
package acme.catalog.v1;
package acme.inventory.v1;
package acme.order.v1;
package acme.payment.v1;
```

Event packages:

```proto
package acme.order.events.v1;
package acme.payment.events.v1;
```

## 5.1 Major Versions

Use a new major package version for breaking changes.

Example:

```proto
package acme.order.v1;
```

evolves to:

```proto
package acme.order.v2;
```

when incompatible changes are required.

## 5.2 Minor Versions

Minor versions do not need to appear in protobuf package names.

Use source control, release tags, changelogs, and generated client versions for minor compatible changes.

---

## 6. gRPC Service Versioning

gRPC services are versioned through their protobuf package.

Example:

```proto
package acme.order.v1;

service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
}
```

A breaking version would become:

```proto
package acme.order.v2;

service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
}
```

## 6.1 Running Multiple Versions

During migration, a service may expose both:

```text
acme.order.v1.OrderService
acme.order.v2.OrderService
```

This should be temporary and documented.

---

## 7. Message Versioning

Messages are versioned through package versioning.

Do not include version numbers in message names unless there is a strong reason.

Good:

```proto
package acme.order.v1;

message Order {}
```

Avoid:

```proto
message OrderV1 {}
```

Use package versioning instead.

---

## 8. Compatible Protobuf Changes

Generally compatible:

```text
adding a new optional field
adding a new repeated field
adding a new message type
adding a new RPC method
adding comments
adding new enum values where consumers handle unknown values
relaxing validation where safe
```

Example:

```proto
message Product {
  string product_id = 1;
  string name = 2;
  optional string material = 3;
}
```

Adding:

```proto
optional string care_instructions = 4;
```

is generally compatible.

---

## 9. Breaking Protobuf Changes

Breaking or high-risk changes:

```text
renaming a package
renaming a service
renaming an RPC
removing a field
renumbering a field
reusing a field number
changing a field type
changing field meaning
changing required business semantics
removing enum values
changing enum numeric values
changing streaming behaviour
changing idempotency behaviour
```

Example breaking change:

```proto
string order_id = 1;
```

changed to:

```proto
int64 order_id = 1;
```

This is breaking and should require a new version.

---

## 10. Reserved Fields

When removing fields, reserve field numbers and names.

Example:

```proto
message Product {
  string product_id = 1;
  string name = 2;

  reserved 3;
  reserved "legacy_sku";
}
```

This prevents accidental reuse and unsafe compatibility issues.

---

## 11. Enum Evolution

Enums should evolve carefully.

Example:

```proto
enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED = 0;
  PAYMENT_STATUS_PENDING = 1;
  PAYMENT_STATUS_AUTHORISED = 2;
  PAYMENT_STATUS_FAILED = 3;
}
```

Adding a value may be compatible:

```proto
PAYMENT_STATUS_CANCELLED = 4;
```

only if consumers handle unknown or unexpected values safely.

Rules:

- Do not rename enum values without a migration plan.
- Do not change numeric values.
- Keep zero value as `_UNSPECIFIED`.
- Consumers should not assume only known values exist forever.

---

## 12. Event Versioning

Kafka events are contracts and must be versioned.

Event versioning may occur through:

```text
protobuf package version
event_version field in envelope
topic version
```

Example package:

```proto
package acme.order.events.v1;
```

Example topic:

```text
bfstore.order.order-events.v1
```

Example envelope:

```json
{
  "event_type": "OrderCreated",
  "event_version": "1.0"
}
```

## 12.1 Event Compatibility

Compatible event changes:

```text
add optional field
add new event type
add new enum value if consumers handle unknowns
add metadata field
```

Breaking event changes:

```text
remove field
rename event type
change field type
change field meaning
change event ordering assumptions
change topic semantics
```

---

## 13. Kafka Topic Versioning

Topic names should include a major version.

Format:

```text
bfstore.<domain>.<stream>.v<major>
```

Examples:

```text
bfstore.order.order-events.v1
bfstore.payment.payment-events.v1
bfstore.catalog.product-events.v1
```

Use a new topic version for incompatible event stream changes.

Example:

```text
bfstore.order.order-events.v2
```

## 13.1 When to Create a New Topic Version

Create a new topic version when:

```text
event payloads are incompatible
event semantics change
partition key changes incompatibly
ordering assumptions change
consumers cannot safely process both versions
```

---

## 14. API Gateway Versioning

If the API Gateway exposes REST/JSON externally, routes should be versioned.

Example:

```text
/api/v1/products
/api/v1/baskets
/api/v1/orders
```

A breaking external API change should use:

```text
/api/v2/orders
```

## 14.1 External API Compatibility

Breaking external changes include:

```text
removing response fields used by clients
renaming fields
changing field type
changing error response shape
changing authentication behaviour
changing route semantics
```

External API versioning should be especially conservative because external clients may be harder to coordinate than internal services.

---

## 15. Application Error Code Versioning

Application error codes are also part of the API contract.

Stable examples:

```text
PRODUCT_NOT_FOUND
BASKET_EMPTY
INSUFFICIENT_STOCK
PAYMENT_DECLINED
ORDER_NOT_FOUND
```

Avoid renaming error codes casually.

If an error code is replaced:

```text
mark old code deprecated
support old behaviour where possible
document migration
update contract tests
```

---

## 16. Generated Client Versioning

Generated clients should be versioned through source control and release tags.

Recommended approach:

```text
generate protobuf clients in CI
publish or commit generated clients according to repo strategy
tag releases that include contract changes
document compatibility in changelog
```

If generated clients are committed to the repo, they should be regenerated consistently using the same tooling.

---

## 17. Deprecation Policy

Deprecation should be explicit.

Deprecation notice should include:

```text
what is deprecated
why it is deprecated
replacement API or field
first deprecated version
planned removal version or date where applicable
consumer impact
migration steps
```

## 17.1 Deprecating Protobuf Fields

A field can be marked deprecated:

```proto
string legacy_sku = 3 [deprecated = true];
```

Do not immediately remove it unless all consumers have migrated.

## 17.2 Removing Deprecated Fields

When removed:

```proto
reserved 3;
reserved "legacy_sku";
```

---

## 18. Version Migration Process

For breaking changes, use a controlled migration process.

Example:

```text
1. Add new v2 contract.
2. Implement v1 and v2 side by side.
3. Add tests for both versions.
4. Update internal consumers to v2.
5. Monitor usage of v1.
6. Announce deprecation of v1.
7. Remove v1 after agreed period.
```

For events:

```text
1. Publish v1 and v2 events where required.
2. Update consumers to handle v2.
3. Monitor consumer lag and DLQ.
4. Stop v1 production after migration.
5. Retire v1 topic after retention period.
```

---

## 19. CI Compatibility Checks

CI should enforce contract quality.

Recommended checks:

```text
buf lint
buf breaking
protobuf generation
contract tests
consumer compatibility tests
event schema tests
API Gateway response shape tests
```

## 19.1 Buf Breaking Checks

Buf should compare protobuf changes against a configured baseline.

Breaking changes should fail CI unless intentionally approved through the versioning process.

## 19.2 Contract Tests

Contract tests should verify:

```text
expected fields
expected error codes
expected event payloads
idempotency behaviour
backward-compatible responses
```

---

## 20. Release Notes

Any contract change should be documented in release notes or changelog.

Release notes should include:

```text
new RPCs
new fields
deprecated fields
new event types
event schema changes
error code changes
breaking changes
migration actions
```

---

## 21. Versioning Examples

## 21.1 Safe Additive Change

Adding a field:

```proto
message Order {
  string order_id = 1;
  string customer_id = 2;
  optional string delivery_note = 3;
}
```

This is generally safe if clients can ignore the new field.

---

## 21.2 Breaking Field Type Change

Changing:

```proto
string order_id = 1;
```

to:

```proto
int64 order_id = 1;
```

is breaking.

Use a new field or new version instead.

---

## 21.3 Safe New RPC

Adding:

```proto
rpc GetOrderHistory(GetOrderHistoryRequest) returns (GetOrderHistoryResponse);
```

is generally safe.

---

## 21.4 Risky Behaviour Change

Changing `CreateOrder` from idempotent to non-idempotent is breaking even if the protobuf shape is unchanged.

Behaviour is part of the contract.

---

## 22. Rollback Considerations

Versioning should support rollback.

Rollback risks include:

```text
new service emits events old consumers cannot parse
new API response removes fields old clients need
new enum values crash old consumers
new database schema incompatible with old service version
```

Mitigation:

```text
use additive changes first
avoid destructive migrations
support old and new contracts temporarily
monitor consumer compatibility
use feature flags where appropriate
```

---

## 23. Initial Versioning Scope

Initial implementation should establish:

```text
protobuf packages using v1
Kafka topics using v1
API Gateway external route versioning using v1 if REST/JSON is used
Buf linting
Buf breaking checks
contract test structure
reserved field rules
deprecation policy
```

This creates a professional foundation even before v2 APIs are needed.

---

## 24. Anti-Patterns to Avoid

Avoid:

```text
unversioned protobuf packages
renaming fields casually
reusing field numbers
removing fields without reservation
changing message meaning without versioning
using event_version but ignoring compatibility
publishing incompatible events to the same topic
breaking external clients without migration
treating internal APIs as disposable
```

---

## 25. Versioning Review Checklist

Before approving a contract change, check:

```text
Is the package versioned?
Is the change backward-compatible?
Are any fields removed or renumbered?
Are removed fields reserved?
Are enum values safe for old consumers?
Are event topics still compatible?
Are application error codes stable?
Do generated clients still compile?
Do buf checks pass?
Are contract tests updated?
Is a migration plan required?
Are release notes updated?
```

---

## 26. Open Questions

| Question | Status |
|---|---|
| Will generated clients be committed or generated in CI only? | To decide |
| What baseline will Buf use for breaking-change detection? | To decide |
| What deprecation period should apply to external API versions? | To decide |
| Will event schema compatibility use Buf alone or a schema registry as well? | To decide |
| Will API Gateway expose REST/JSON, GraphQL, or gRPC-Web? | To decide |
| How will version usage be monitored in production-style environments? | To decide |

---

## 27. Related Documents

This document should be read alongside:

```text
docs/api/grpc-overview.md
docs/api/protobuf-style-guide.md
docs/api/error-model.md
docs/events/event-catalog.md
docs/events/event-envelope.md
docs/events/event-versioning.md
docs/architecture/communication-patterns.md
docs/architecture/tradeoffs.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0002-use-grpc-for-service-communication.md
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
adr/0008-use-contract-first-service-design.md
```

---

## 28. Summary

bfstore versioning is designed to make API and event evolution safe, deliberate, and testable.

The most important rules are:

```text
version protobuf packages
prefer additive changes
reserve removed fields
never reuse field numbers
treat behaviour as part of the contract
version Kafka topics for breaking changes
use Buf and contract tests in CI
document deprecations and migrations
```

This approach helps bfstore demonstrate professional API governance and long-term maintainability.
