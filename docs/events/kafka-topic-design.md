# Kafka Topic Design

## 1. Purpose

This document defines the Kafka topic design standards for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how topics are named, owned, versioned, partitioned, retained, observed, and governed.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s event-driven architecture and operational maturity.

---

## 2. Scope

This document applies to Kafka topics used by bfstore services, including:

```text
domain event topics
dead-letter topics
retry topics where used
projection topics
operational event topics
```

It covers:

```text
topic naming
topic ownership
domain alignment
versioning
partition keys
retention
consumer groups
DLQ naming
environment naming
security considerations
observability
testing expectations
```

---

## 3. Topic Design Goals

Kafka topics should be:

| Goal | Description |
|---|---|
| Predictable | Engineers can infer purpose from topic name |
| Domain-aligned | Topics reflect service and business ownership |
| Versioned | Incompatible changes can be introduced safely |
| Observable | Throughput, lag, and failures can be monitored |
| Secure | Access can be restricted by service and domain |
| Scalable | Partitioning supports expected throughput and ordering |
| Maintainable | Topic sprawl and event soup are avoided |
| Testable | Topic behaviour is covered in integration and contract tests |

---

## 4. Topic Naming Standard

Recommended topic format:

```text
bfstore.<domain>.<stream>.v<major>
```

Example:

```text
bfstore.order.order-events.v1
```

## 4.1 Naming Components

| Component | Example | Description |
|---|---|---|
| `bfstore` | `bfstore` | Application or product namespace |
| `<domain>` | `order` | Business or service domain |
| `<stream>` | `order-events` | Logical stream purpose |
| `v<major>` | `v1` | Major contract version |

---

## 5. Standard Topic Examples

| Topic | Owner | Purpose |
|---|---|---|
| `bfstore.catalog.product-events.v1` | `catalog-service` | Product catalogue events |
| `bfstore.inventory.stock-events.v1` | `inventory-service` | Stock and reservation events |
| `bfstore.basket.basket-events.v1` | `basket-service` | Basket lifecycle events |
| `bfstore.order.order-events.v1` | `order-service` | Order lifecycle events |
| `bfstore.payment.payment-events.v1` | `payment-service` | Payment lifecycle events |
| `bfstore.shipping.shipment-events.v1` | `shipping-service` | Shipment lifecycle events |
| `bfstore.notification.notification-events.v1` | `notification-service` | Notification delivery events |
| `bfstore.review.review-events.v1` | `review-service` | Review and moderation events |
| `bfstore.search.search-events.v1` | `search-service` | Search projection events |
| `bfstore.recommendation.recommendation-events.v1` | `recommendation-service` | Recommendation events |

---

## 6. Domain Ownership

Each topic must have a clear owning service.

Rules:

```text
the owner defines event meaning
the owner publishes events to the topic
the owner manages schema evolution
the owner documents event semantics
consumers must not publish to another service's owned event stream
```

Example:

```text
order-service owns bfstore.order.order-events.v1
payment-service owns bfstore.payment.payment-events.v1
```

A topic without an owner becomes difficult to govern and should be avoided.

---

## 7. Topic Granularity

## 7.1 Avoid One Topic for Everything

Avoid:

```text
bfstore.events.v1
```

Problems:

```text
unclear ownership
difficult access control
difficult retention tuning
noisy consumers
harder observability
```

## 7.2 Avoid One Topic per Tiny Event by Default

Avoid excessive topic sprawl such as:

```text
bfstore.order.order-created.v1
bfstore.order.order-cancelled.v1
bfstore.order.order-failed.v1
```

unless there is a strong operational reason.

## 7.3 Recommended Balance

Use one topic per domain event stream.

Example:

```text
bfstore.order.order-events.v1
```

containing:

```text
OrderCreated
OrderConfirmed
OrderCancelled
OrderFailed
```

This balances ownership, observability, and operational simplicity.

---

## 8. Topic Versioning

Topic major versions should be used for incompatible stream changes.

Example:

```text
bfstore.order.order-events.v1
bfstore.order.order-events.v2
```

Create a new topic version when:

```text
event payloads become incompatible
partition key changes incompatibly
ordering guarantees change
topic semantics change
consumers cannot safely process both versions
```

Do not create a new topic version for every additive field.

Compatible event changes should usually remain on the same topic version.

---

## 9. Event Version vs Topic Version

Topic version and event version are related but not identical.

| Concept | Example | Purpose |
|---|---|---|
| Topic version | `bfstore.order.order-events.v1` | Stream-level compatibility |
| Event version | `event_version: 1.1` | Event schema or payload compatibility |
| Protobuf package version | `acme.order.events.v1` | Contract package version |

A minor event schema addition may update `event_version` without changing topic version.

A breaking stream-level change should use a new topic version.

---

## 10. Partitioning Strategy

Partitioning affects ordering, throughput, and scalability.

## 10.1 Recommended Partition Keys

| Topic | Suggested Partition Key |
|---|---|
| `bfstore.catalog.product-events.v1` | `product_id` |
| `bfstore.inventory.stock-events.v1` | `product_id` or `reservation_id` |
| `bfstore.basket.basket-events.v1` | `basket_id` |
| `bfstore.order.order-events.v1` | `order_id` |
| `bfstore.payment.payment-events.v1` | `payment_id` or `order_id` |
| `bfstore.shipping.shipment-events.v1` | `shipment_id` or `order_id` |
| `bfstore.notification.notification-events.v1` | `notification_id` |
| `bfstore.review.review-events.v1` | `review_id` or `product_id` |
| `bfstore.search.search-events.v1` | `projection_id` or `product_id` |
| `bfstore.recommendation.recommendation-events.v1` | `customer_id` or `product_id` |

## 10.2 Partitioning Rules

```text
choose keys based on ordering needs
keep related aggregate events on the same key where possible
avoid keys with extreme skew
document the selected key per topic
do not assume ordering across different keys
```

---

## 11. Ordering Expectations

Kafka ordering is guaranteed only within a partition.

This means:

```text
events with the same partition key can be processed in order
events with different partition keys may be processed independently
global ordering should not be assumed
```

For bfstore:

```text
OrderCreated and OrderCancelled for the same order should use order_id
ProductUpdated events for the same product should use product_id
Payment events for the same payment should use payment_id
```

---

## 12. Consumer Groups

Each logical consumer should use its own consumer group.

Examples:

```text
notification-service.order-events-consumer
search-service.product-events-consumer
recommendation-service.order-events-consumer
analytics.order-events-consumer
```

Rules:

```text
consumer group names should identify service and purpose
do not share consumer groups between unrelated services
consumer lag should be monitored per group
consumer group names should be stable across deployments
```

---

## 13. Dead-Letter Topic Naming

DLQ topic format:

```text
<source-topic>.dlq
```

Examples:

```text
bfstore.order.order-events.v1.dlq
bfstore.payment.payment-events.v1.dlq
bfstore.notification.notification-events.v1.dlq
```

Alternative format:

```text
bfstore.<domain>.<stream>.dlq.v<major>
```

Example:

```text
bfstore.order.order-events.dlq.v1
```

## 13.1 Initial Recommendation

Use:

```text
bfstore.<domain>.<stream>.dlq.v<major>
```

This keeps DLQ topics clearly named and easy to manage.

---

## 14. Retry Topics

Retry topics may be introduced if the consumer retry model requires delayed retries.

Example format:

```text
bfstore.order.order-events.retry.5m.v1
bfstore.order.order-events.retry.30m.v1
```

Initial recommendation:

```text
start with consumer-level retry and DLQ
introduce retry topics only when needed
```

Retry topic complexity should be justified by operational requirements.

---

## 15. Environment Naming

There are two common approaches.

## 15.1 Environment in Topic Name

Example:

```text
dev.bfstore.order.order-events.v1
test.bfstore.order.order-events.v1
prod.bfstore.order.order-events.v1
```

## 15.2 Environment in Kafka Cluster or Namespace

Example:

```text
bfstore.order.order-events.v1
```

with separate Kafka clusters or namespaces per environment.

## 15.3 Recommendation

Prefer environment separation through Kafka cluster, namespace, or platform configuration where practical.

If shared Kafka infrastructure is used, include environment in the topic name or access policy.

The chosen approach should align with platform standards.

---

## 16. Retention Strategy

Retention should reflect event purpose.

| Topic Type | Example Retention Consideration |
|---|---|
| Domain events | Retain long enough for replay, audit, and recovery needs |
| Notification events | Retain for operational troubleshooting |
| Search projection events | Retain long enough for index rebuild if needed |
| DLQ topics | Retain long enough for investigation and replay |
| Retry topics | Retain according to retry window |

Initial retention values should be documented per environment.

Production retention should consider:

```text
storage cost
audit needs
replay needs
privacy requirements
operational support
```

---

## 17. Compaction

Kafka log compaction may be useful for projection-style topics where the latest state per key matters.

Potential candidates:

```text
product projection events
inventory availability summary
search projection state
```

Use compaction carefully.

Compaction is not a replacement for audit event history where the sequence of events matters.

---

## 18. Access Control

Topic access should follow least privilege.

Examples:

| Service | Allowed Actions |
|---|---|
| `order-service` | Produce to order events |
| `notification-service` | Consume order events, produce notification events |
| `search-service` | Consume catalogue and inventory events |
| `payment-service` | Produce payment events |
| `inventory-service` | Produce inventory events |

Services should not have broad produce/consume access to all topics.

---

## 19. Topic Creation and Governance

Topic creation should be deliberate.

Each topic should have:

```text
owner
purpose
event types
partition key
retention policy
consumer groups
DLQ strategy
schema references
observability expectations
```

Topic definitions may be managed through:

```text
Terraform
GitOps
Kafka operator
platform automation
```

The final implementation belongs in the platform infrastructure and GitOps repositories.

---

## 20. Observability Requirements

Topic dashboards should show:

```text
messages produced
messages consumed
consumer lag
publish failures
processing failures
DLQ count
topic throughput
partition skew
oldest unprocessed message age
```

Important alerts:

```text
consumer lag above threshold
DLQ message count above threshold
publish failure rate above threshold
no messages published when expected
partition skew causing processing delay
```

---

## 21. Testing Requirements

Kafka topic design should be tested through:

```text
producer integration tests
consumer integration tests
contract tests
DLQ tests
consumer group tests
partition key tests
event replay tests
```

Test cases:

```text
OrderCreated is published to bfstore.order.order-events.v1
Notification Service consumes from correct topic and group
Invalid OrderCreated is sent to order events DLQ
ProductUpdated uses product_id as partition key
Duplicate events are handled safely by consumer
```

---

## 22. Initial Topic Set

The first implementation should focus on checkout topics.

Required initial topics:

```text
bfstore.inventory.stock-events.v1
bfstore.payment.payment-events.v1
bfstore.shipping.shipment-events.v1
bfstore.order.order-events.v1
bfstore.notification.notification-events.v1
```

Optional early topics:

```text
bfstore.catalog.product-events.v1
```

Deferred topics:

```text
bfstore.review.review-events.v1
bfstore.search.search-events.v1
bfstore.recommendation.recommendation-events.v1
```

---

## 23. Anti-Patterns to Avoid

Avoid:

```text
one global topic for all events
unclear topic ownership
unversioned topics
consumer groups shared by unrelated services
no DLQ strategy
no retention policy
no partition key documentation
topic names tied to implementation details
putting environment names inconsistently in topics
giving every service access to every topic
```

---

## 24. Topic Design Checklist

Before creating a topic, confirm:

```text
Who owns this topic?
What business domain does it represent?
What events will it contain?
What is the topic name?
What is the major version?
What is the partition key?
What ordering is expected?
What consumers are expected?
What is the retention policy?
Does it need a DLQ?
Does it require compaction?
What access controls apply?
What dashboards and alerts are needed?
What tests cover it?
```

---

## 25. Open Questions

| Question | Status |
|---|---|
| Will topics be created manually, through Terraform, or through a Kafka operator? | To decide |
| Will environments use separate Kafka clusters or environment-prefixed topics? | To decide |
| What retention period should be used for production domain events? | To decide |
| Which topics, if any, should use log compaction? | To decide |
| Should retry topics be implemented initially or deferred? | Proposed: defer |
| What ACL model will be used in local and Kubernetes environments? | To decide |

---

## 26. Related Documents

This document should be read alongside:

```text
docs/events/event-catalog.md
docs/events/event-envelope.md
docs/events/ordering-and-idempotency.md
docs/events/retry-and-dlq-strategy.md
docs/architecture/event-driven-design.md
docs/architecture/communication-patterns.md
docs/testing/testing-strategy.md
docs/operations/runbooks.md
```

---

## 27. Summary

bfstore Kafka topics should be domain-aligned, versioned, owned, observable, and governed.

The recommended naming pattern is:

```text
bfstore.<domain>.<stream>.v<major>
```

The first implementation should focus on checkout-related topics, then expand into catalogue, search, review, recommendation, analytics, and operational streams as the platform matures.
