# User Journeys

## 1. Document Purpose

This document describes the key user journeys for **bfstore**, ACME Ltd’s fictional online furniture store backend.

User journeys explain how different users interact with the system to achieve their goals. They help guide:

- functional requirements
- service boundaries
- gRPC API design
- Kafka event design
- data ownership
- acceptance criteria
- end-to-end testing
- observability requirements
- operational runbooks

This document focuses on the business flows that the backend must support.

---

## 2. Product Context

bfstore is the backend platform for ACME Ltd’s online furniture store.

Customers should be able to browse furniture, search the catalogue, manage a basket, place orders, reserve stock, authorise payments, arrange delivery, receive notifications, submit reviews, and receive product recommendations.

The backend is implemented as a microservice architecture using:

- gRPC for synchronous service-to-service communication
- Kafka for asynchronous business events
- Protobuf for API and event contracts
- MySQL for service-owned databases

---

## 3. Journey Design Principles

The user journeys should follow these principles:

1. **Start with user goals**  
   Journeys should describe what the user is trying to achieve, not which service is called first.

2. **Map business flow to service behaviour**  
   Each journey should identify the services involved and the expected backend responsibilities.

3. **Separate synchronous and asynchronous work**  
   Immediate decisions use gRPC. Downstream reactions use Kafka events.

4. **Include success and failure paths**  
   Real systems must handle unavailable stock, failed payments, duplicate events, and service failures.

5. **Link journeys to tests**  
   Each journey should be testable through end-to-end, integration, contract, performance, or resilience tests.

6. **Make observability explicit**  
   Important journeys should define the logs, metrics, traces, and alerts needed to diagnose them.

---

## 4. Primary Actors

| Actor | Description |
|---|---|
| Customer | A person browsing and purchasing furniture from ACME Ltd |
| Registered Customer | A signed-in customer with saved profile, addresses, orders, and preferences |
| Guest Customer | A customer browsing before signing in or registering |
| ACME Support User | An internal user who may need to inspect orders, payments, shipments, or customer issues |
| ACME Operations Team | Team responsible for monitoring and operating the platform |
| ACME Engineering Team | Team responsible for building, testing, deploying, and maintaining services |
| External Payment Provider | Simulated or future third-party payment provider |
| External Notification Provider | Simulated or future email/SMS provider |
| External Shipping Provider | Simulated or future carrier/fulfilment provider |

---

## 5. Core User Journeys Summary

| ID | Journey | Primary Actor | Priority | Initial Version |
|---|---|---|---|---|
| `UJ-001` | Browse furniture catalogue | Customer | Must | Yes |
| `UJ-002` | View product details | Customer | Must | Yes |
| `UJ-003` | Search and filter products | Customer | Should | Later |
| `UJ-004` | Register account | Customer | Should | Later |
| `UJ-005` | Sign in | Registered Customer | Should | Later |
| `UJ-006` | Manage customer profile and addresses | Registered Customer | Should | Later |
| `UJ-007` | Add product to basket | Customer | Must | Yes |
| `UJ-008` | Update basket | Customer | Must | Yes |
| `UJ-009` | Checkout and create order | Customer | Must | Yes |
| `UJ-010` | Reserve stock during checkout | System | Must | Yes |
| `UJ-011` | Authorise payment during checkout | System | Must | Yes |
| `UJ-012` | Create shipment | System | Must | Yes |
| `UJ-013` | Send order confirmation notification | System | Must | Yes |
| `UJ-014` | View order history | Registered Customer | Should | Later |
| `UJ-015` | Track shipment | Registered Customer | Should | Later |
| `UJ-016` | Submit product review | Registered Customer | Could | Later |
| `UJ-017` | Receive product recommendations | Customer | Could | Later |
| `UJ-018` | Handle payment failure | Customer/System | Must | Yes |
| `UJ-019` | Handle insufficient stock | Customer/System | Must | Yes |
| `UJ-020` | Cancel order | Registered Customer | Should | Later |

---

# 6. Initial Vertical Slice

The first implementation should focus on a complete checkout journey.

```text
Browse product
    -> View product details
    -> Add to basket
    -> Checkout
    -> Reserve stock
    -> Authorise payment
    -> Create order
    -> Create shipment
    -> Publish OrderCreated event
    -> Send notification
```

Initial services involved:

```text
api-gateway
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

The purpose of this vertical slice is to prove the core architecture before expanding into search, reviews, recommendations, advanced account features, and production platform maturity.

---

## 7. Journey Details

### UJ-001: Browse Furniture Catalogue

**Goal**

A customer wants to browse available furniture products by category.

**Primary Actor**

Customer

**Priority**


Must

**Services Involved**

| Service             | Responsibility                                                |
| ------------------- | ------------------------------------------------------------- |
| `api-gateway`       | Receives client request and routes to catalogue               |
| `catalog-service`   | Returns product catalogue data                                |
| `inventory-service` | May provide stock availability summary                        |
| `search-service`    | Later phase: supports search/filter-backed catalogue browsing |


**Preconditions**

- Catalogue products exist.
- Products have active/inactive status.
- Product categories exist.
- API Gateway is available.
- Catalogue Service is available.

**Main Success Flow**

1. Customer opens the furniture catalogue.
2. API Gateway receives a request to list products.
3. API Gateway calls Catalogue Service using gRPC.
4. Catalogue Service retrieves active products.
5. Catalogue Service returns product summaries.
6. API Gateway returns the product list to the client.

**Example Product Summary**

```text
product_id
name
category
thumbnail_url
price
currency
material
colour
dimensions_summary
availability_summary
```

**Alternative Flows**

A1: Category Filter Applied

1. Customer selects a category, such as sofas or dining tables.
2. API Gateway sends category filter to Catalogue Service.
3. Catalogue Service returns active products in that category.

A2: No Products Found

1. Customer selects a category with no active products.
2. Catalogue Service returns an empty result set.
3. API Gateway returns an empty catalogue response.

   
**Failure Flows**

F1: Catalogue Service Unavailable
API Gateway calls Catalogue Service.
Catalogue Service is unavailable.
API Gateway returns an appropriate service unavailable error.
Error is logged with correlation ID.
Events

No event is required for simple browsing in the initial version.

Later phases may emit analytics events, for example:

```text
ProductListViewed
CategoryViewed
```

**Observability Requirements**

- Log catalogue list requests.
- Record request latency.
- Record error count.
- Trace API Gateway to Catalogue Service call.
- Track empty result responses.

**Acceptance Criteria**

```text
Given active products exist
When a customer browses the catalogue
Then the system returns a list of active product summaries

Given a category has no active products
When a customer browses that category
Then the system returns an empty product list without error
```

---




