# bfstore — Borough Furniture Store

## 1. Overview

**bfstore** is the backend platform for **Borough Furniture Store**, a fictional cloud-native ecommerce company used to demonstrate senior-level platform engineering, DevSecOps, Kubernetes, microservices, observability, testing, and cloud infrastructure skills.

The system is designed as a professional portfolio project for demonstrating how a modern ecommerce backend can be documented, designed, implemented, tested, deployed, and operated.

---

## 2. Company Backstory

**Borough Furniture Store** began as a small niche shop selling speciality Golang-themed homeware to Go enthusiasts.

The name **Borough** is inspired by the home of the Go gopher mascot.

The company originally sold playful, developer-themed homeware such as:

```text
Gopher cushions
Gopher desk lamps
Gopher mugs
Gopher wall art
```

As demand grew, Borough Furniture Store expanded into broader programming-language mascot-themed and computer-science-inspired homeware.

Example product ideas include:

```text
Rob Pike wall tapestry
Rivest super-secure lockbox
Dijkstra pathfinding rug
Grace Hopper debugging blanket
Turing machine wall clock
Rust crab coat hooks
Python plush beanbag
Kubernetes helm bookends
```

This backstory makes the project more memorable than a generic furniture store while still supporting realistic ecommerce architecture.

---

## 3. Purpose of the Project

The purpose of bfstore is to demonstrate a production-style ecommerce backend using:

```text
Go
gRPC
Protocol Buffers
Kafka
MySQL
Docker
Kubernetes
CI/CD
Infrastructure as Code
Observability
Security controls
Testing strategy
Operational documentation
```

The project is intentionally designed to support client-facing review.

It should show not only that the application works, but that the engineering decisions are documented, defensible, and aligned with real-world platform engineering expectations.

---

## 4. Architecture Summary

bfstore is designed as a microservice-based backend.

Core services include:

```text
api-gateway
auth-service
customer-service
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
review-service
search-service
recommendation-service
```

The first implementation focuses on a smaller checkout vertical slice:

```text
Browse product
→ Add product to basket
→ Checkout
→ Reserve stock
→ Authorise payment
→ Create shipment
→ Create order
→ Publish OrderCreated
→ Send notification
```

---

## 5. Communication Model

bfstore uses a hybrid communication model.

```text
gRPC = commands and queries requiring an immediate response
Kafka = facts that have already happened
```

Examples:

```text
ReserveStock is a gRPC command.
StockReserved is a Kafka event.

AuthorisePayment is a gRPC command.
PaymentAuthorised is a Kafka event.

CreateOrder is a gRPC command.
OrderCreated is a Kafka event.
```

---

## 6. Data Model Summary

Each service owns its own data.

Catalogue Service owns product truth, including:

```text
product identity
product descriptions
category taxonomy
product variants
category-scoped product attributes
catalogue pricing in the initial version
product status
```

Inventory Service owns stock.

Order Service owns orders.

Payment Service owns payments.

Shipping Service owns shipments.

Notification Service owns notification delivery state.

Search Service owns denormalised product search projections.

---

## 7. Product Catalogue Direction

Borough Furniture Store sells varied product types, so the catalogue must support flexible product attributes.

Examples:

```text
curtains need drop, width, lining, heading type
bed frames need bed size, material, storage type
rugs need shape, pile height, weave
lamps need bulb type, wattage, fitting type
```

The project uses:

```text
MySQL as the governed catalogue source of truth
category-scoped product attribute definitions
product attribute values
Search Service as a denormalised browse/search projection
```

This avoids both:

```text
one huge products table with hundreds of nullable fields
uncontrolled schemaless product documents with weak governance
```

---

## 8. Repository Structure

Recommended high-level layout:

```text
bfstore/
├── README.md
├── Makefile
├── docker-compose.yml
├── buf.yaml
├── buf.gen.yaml
├── .env.example
│
├── docs/
├── adr/
├── proto/
├── services/
├── packages/
├── db/
├── deploy/
├── tests/
├── tools/
├── scripts/
└── .github/
```

---

## 9. Documentation

Key documentation areas:

```text
docs/requirements/
docs/architecture/
docs/api/
docs/events/
docs/data/
docs/testing/
docs/security/
docs/observability/
docs/operations/
adr/
```

The documentation is intended to be client-readable and to explain:

```text
what the system does
why it is designed this way
how services communicate
which service owns which data
how APIs and events are versioned
how the platform is tested
how the system is deployed and operated
```

---

## 10. Initial Implementation Target

The first implementation should prove the architecture through one complete customer journey.

Initial target:

```text
List products
Get product
Create basket
Add item to basket
Checkout
Reserve stock
Authorise payment
Create shipment
Create order
Publish OrderCreated
Notification Service consumes OrderCreated
```

---

## 11. Example Seed Products

Initial demo catalogue seed data should reflect the Borough Furniture Store story.

Suggested products:

```text
Gopher desk lamp
Gopher cushion set
Rob Pike wall tapestry
Rivest super-secure lockbox
Dijkstra pathfinding rug
Grace Hopper debugging blanket
Turing machine wall clock
Rust crab coat hooks
Python plush beanbag
Kubernetes helm bookends
```

The seed data should include multiple product categories and varied attributes to prove the flexible catalogue model.

---

## 12. Professional Goals Demonstrated

bfstore is intended to demonstrate:

```text
microservice architecture
contract-first API design
event-driven design
service-owned databases
MySQL schema design
Kubernetes deployment thinking
CI/CD quality gates
observability
resilience patterns
security and privacy awareness
testing strategy
technical documentation
architecture decision records
```

---

## 13. Summary

bfstore is not just a demo shop.

It is a professional platform engineering portfolio project wrapped in a memorable, developer-themed ecommerce story.

**Borough Furniture Store** gives the system a clear identity while still supporting serious engineering conversations about microservices, data ownership, events, testing, deployment, and operations.
