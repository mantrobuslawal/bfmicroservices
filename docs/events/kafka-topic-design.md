# Kafka Topic Design

## 1. Purpose

This document defines Kafka topic design for **bfstore**.

bfstore uses Kafka for domain events and **Protocol Buffers** for event payloads.

---

## 2. Topic Naming

Use:

```text
bfstore.<domain>.<stream>.v<major>
```

Examples:

```text
bfstore.catalog.products.v1
bfstore.inventory.stock.v1
bfstore.order.orders.v1
bfstore.payment.payments.v1
bfstore.shipping.shipments.v1
bfstore.notification.notifications.v1
```

---

## 3. Serialisation Standard

```text
Kafka key     = stable business identifier
Kafka value   = Protobuf event message
Kafka headers = lightweight metadata
```

Recommended header:

```text
content_type = application/x-protobuf
```

JSON is not the production event payload format.

---

## 4. Topic Ownership

| Topic | Owning Producer |
|---|---|
| `bfstore.catalog.products.v1` | `catalog-service` |
| `bfstore.inventory.stock.v1` | `inventory-service` |
| `bfstore.basket.baskets.v1` | `basket-service` |
| `bfstore.order.orders.v1` | `order-service` |
| `bfstore.payment.payments.v1` | `payment-service` |
| `bfstore.shipping.shipments.v1` | `shipping-service` |
| `bfstore.notification.notifications.v1` | `notification-service` |

Only the owning service should publish to its topic.

---

## 5. Protobuf Package Mapping

Topic streams should map clearly to event proto packages.

```text
bfstore.order.orders.v1
  -> proto/acme/order/events/v1/order_events.proto

bfstore.catalog.products.v1
  -> proto/acme/catalog/events/v1/catalog_events.proto
```

Shared metadata:

```text
proto/acme/events/v1/event_metadata.proto
```

---

## 6. Message Key Strategy

| Topic | Recommended Key |
|---|---|
| `bfstore.catalog.products.v1` | `product_id` |
| `bfstore.inventory.stock.v1` | `reservation_id`, `product_id`, or `variant_id` depending on ordering need |
| `bfstore.basket.baskets.v1` | `basket_id` |
| `bfstore.order.orders.v1` | `order_id` |
| `bfstore.payment.payments.v1` | `order_id` or `payment_id` |
| `bfstore.shipping.shipments.v1` | `order_id` or `shipment_id` |
| `bfstore.notification.notifications.v1` | `notification_id` |

Ordering is only guaranteed within a partition.

---

## 7. Headers

Recommended headers:

```text
event_type
event_version
correlation_id
traceparent
producer
content_type
```

Headers support observability and routing. The Protobuf payload remains the contract.

---

## 8. Topic Versioning

Topic major versions change for breaking stream changes.

```text
bfstore.order.orders.v1
bfstore.order.orders.v2
```

Compatible Protobuf additions do not automatically require a new topic.

---

## 9. Retry and DLQ Topics

Retry:

```text
bfstore.<domain>.<stream>.retry.v<major>
```

DLQ:

```text
bfstore.<domain>.<stream>.dlq.v<major>
```

DLQ records should preserve original binary payload, headers, topic, partition, offset, consumer group, and failure reason.

---

## 10. Protobuf Decode Failures

Consumers must handle decode failures safely.

Decode failures may indicate:

```text
wrong topic
wrong message type
unsupported event version
corrupt payload
producer bug
consumer contract drift
```

Handling:

```text
log safe metadata from headers
increment decode failure metric
send to DLQ where appropriate
alert if repeated
avoid crash loops
```

---

## 11. Initial Topic Set

```text
bfstore.order.orders.v1
bfstore.inventory.stock.v1
bfstore.payment.payments.v1
bfstore.shipping.shipments.v1
bfstore.notification.notifications.v1
```

Later:

```text
bfstore.catalog.products.v1
```

---

## 12. Summary

bfstore Kafka topics are domain-owned, versioned streams carrying Protobuf event payloads.


---

## Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/events/event-envelope.md
docs/events/event-catalog.md
docs/events/kafka-topic-design.md
adr/0003-use-kafka-for-events.md
adr/0006-use-buf-for-protobuf.md
```
