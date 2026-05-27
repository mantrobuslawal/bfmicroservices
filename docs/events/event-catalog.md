# Event Catalog

## 1. Purpose

This document defines the event catalog for **bfstore**.

bfstore uses Kafka for asynchronous service communication and **Protocol Buffers** for Kafka event payloads.

Events represent facts that have already happened.

---

## 2. Event Principles

```text
events are past-tense facts
events are owned by the producing service
events are Protobuf messages
events include standard metadata
events are versioned
events are backwards-compatible where possible
events exclude unnecessary sensitive data
```

Commands that need an immediate response use gRPC. Facts that have already happened use Kafka.

---

## 3. Serialisation Standard

```text
Kafka key     = stable business identifier
Kafka value   = Protobuf event message
Kafka headers = event_type, event_version, correlation_id, traceparent, content_type
```

Recommended content type:

```text
application/x-protobuf
```

---

## 4. Event Topic Summary

| Domain | Topic | Producer |
|---|---|---|
| Catalogue | `bfstore.catalog.products.v1` | `catalog-service` |
| Inventory | `bfstore.inventory.stock.v1` | `inventory-service` |
| Basket | `bfstore.basket.baskets.v1` | `basket-service` |
| Order | `bfstore.order.orders.v1` | `order-service` |
| Payment | `bfstore.payment.payments.v1` | `payment-service` |
| Shipping | `bfstore.shipping.shipments.v1` | `shipping-service` |
| Notification | `bfstore.notification.notifications.v1` | `notification-service` |

---

## 5. Catalogue Events

```text
ProductCreated
ProductUpdated
ProductActivated
ProductDeactivated
ProductArchived
CategoryCreated
CategoryUpdated
ProductAttributeDefinitionCreated
ProductAttributeDefinitionUpdated
ProductAttributeDefinitionDeprecated
```

Consumers:

```text
search-service
recommendation-service
analytics later
```

Catalogue events allow Search Service to maintain denormalised browse and filter projections.

---

## 6. Inventory Events

```text
StockReserved
StockReservationFailed
StockReservationReleased
StockReservationExpired
StockCommitted
InventoryAdjusted
```

Initial priority:

```text
StockReserved
StockReservationFailed
StockReservationReleased
```

---

## 7. Order Events

```text
OrderCreated
OrderFailed
OrderConfirmed
OrderCancelled
```

Initial priority:

```text
OrderCreated
OrderFailed
```

`OrderCreated` is a critical event and should use the outbox pattern where practical.

---

## 8. Payment Events

```text
PaymentAuthorised
PaymentFailed
PaymentCaptured
PaymentRefunded
PaymentCancelled
```

Payment events must not contain raw card data, CVV, provider secrets, or sensitive payment credentials.

---

## 9. Shipping Events

```text
ShipmentCreated
ShipmentFailed
ShipmentDispatched
ShipmentDelivered
ShipmentCancelled
```

Shipment events should minimise address data.

---

## 10. Notification Events

```text
NotificationRequested
NotificationSent
NotificationFailed
NotificationSuppressed
```

Notification consumers must be idempotent to avoid duplicate customer messages.

---

## 11. Event Proto Example

```proto
message OrderCreatedEvent {
  acme.events.v1.EventMetadata metadata = 1;
  OrderCreated payload = 2;
}

message OrderCreated {
  string order_id = 1;
  string order_number = 2;
  string customer_id = 3;
  acme.common.v1.Money total_amount = 4;
  repeated OrderItemSnapshot items = 5;
  google.protobuf.Timestamp created_at = 6;
}
```

---

## 12. Compatibility Rules

```text
never reuse field numbers
reserve removed fields
reserve removed field names
add fields instead of changing existing fields
avoid changing field meaning
include UNKNOWN/UNSPECIFIED enum value at 0
run Buf breaking checks in CI
```

---

## 13. Initial Checkout Event Flow

```text
Inventory Service publishes StockReserved or StockReservationFailed.
Payment Service publishes PaymentAuthorised or PaymentFailed.
Shipping Service publishes ShipmentCreated or ShipmentFailed.
Order Service publishes OrderCreated or OrderFailed.
Notification Service consumes OrderCreated.
Notification Service publishes NotificationSent or NotificationFailed.
```

---

## 14. Summary

bfstore events are Kafka messages with Protobuf payloads. The event catalog defines the typed domain facts exchanged between services.


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
