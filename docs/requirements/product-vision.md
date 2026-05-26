# Product Vision

## 1. Product Name

**bfstore** — **Borough Furniture Store**

---

## 2. Vision Statement

Borough Furniture Store is a fictional online furniture and homeware retailer built to demonstrate a modern, cloud-native ecommerce backend.

The vision is to create a technically credible, client-reviewable ecommerce platform that showcases strong backend engineering, platform engineering, DevSecOps, Kubernetes, observability, testing, and operational maturity.

bfstore should feel like a real business system, not a toy demo.

---

## 3. Company Backstory

Borough Furniture Store began as a small speciality shop selling Golang-themed homeware to Go enthusiasts.

The name **Borough** is inspired by the home of the Go gopher mascot.

The company’s early products included playful developer-themed items such as:

```text
Gopher cushions
Gopher desk lamps
Gopher mugs
Gopher wall art
```

Over time, Borough Furniture Store expanded beyond Go-themed products into broader programming-language mascot-themed and computer-science-inspired homeware.

Example product lines include:

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

This gives bfstore a distinctive identity while still supporting realistic ecommerce requirements.

---

## 4. Why This Product Exists

bfstore exists to demonstrate how a modern ecommerce backend can be:

```text
designed
documented
implemented
tested
secured
deployed
observed
operated
```

It is intended as a portfolio-quality system for someone targeting roles such as:

```text
Senior Platform Engineer
DevSecOps Engineer
Kubernetes Platform Engineer
Cloud Engineer
Backend Platform Engineer
```

---

## 5. Target Users

## 5.1 Customer

A customer wants to browse, filter, and buy developer-themed furniture and homeware.

Customer goals:

```text
browse interesting products
view product details
choose product variants
add items to basket
checkout securely
receive order confirmation
track fulfilment
```

## 5.2 Store Operator

A store operator wants to manage products, stock, fulfilment, and customer communication.

Operator goals:

```text
publish products
manage product categories
manage category-specific attributes
monitor stock
process orders
handle failed payments
track shipments
monitor notifications
```

## 5.3 Engineering/Platform Reviewer

A reviewer wants to understand whether the system demonstrates professional engineering maturity.

Reviewer goals:

```text
understand architecture decisions
review service boundaries
see clear API and event contracts
validate data ownership
assess deployment and testing strategy
evaluate operational readiness
```

---

## 6. Product Differentiator

bfstore is memorable because it combines:

```text
a playful developer-themed product catalogue
serious microservice architecture
professional documentation
cloud-native deployment thinking
platform engineering evidence
```

The catalogue story gives the project personality.

The architecture gives it credibility.

---

## 7. Business Capabilities

bfstore should support these business capabilities:

```text
Product catalogue management
Category and product attribute management
Product browsing
Basket management
Checkout
Stock reservation
Payment authorisation
Shipment creation
Order lifecycle management
Customer notification
Product search
Product recommendations
Reviews and ratings
```

---

## 8. Product Catalogue Vision

The catalogue must support varied product types.

Examples:

```text
curtains
bed frames
mattresses
sofas
rugs
lamps
wardrobes
wall art
lockboxes
desk accessories
developer-themed homeware
```

These products share common data, but have different product-specific attributes.

Examples:

```text
curtains: drop, width, lining, heading type
bed frames: bed size, material, storage type
rugs: shape, pile height, material
lamps: bulb type, wattage, fitting type
lockboxes: lock type, security rating, dimensions
wall tapestries: fabric, width, height, mounting style
```

The vision is to support this using:

```text
governed product data in Catalogue Service
category-scoped attribute definitions
product attribute values
Search Service projections for browse and filtering
```

This keeps product data flexible without losing quality or governance.

---

## 9. Initial Product Experience

The first customer experience should be:

```text
A customer opens the store.
The customer browses developer-themed homeware.
The customer selects a product such as a Gopher desk lamp or Rob Pike wall tapestry.
The customer adds it to their basket.
The customer checks out.
The system reserves stock.
The system authorises payment.
The system creates a shipment.
The system creates an order.
The system sends an order confirmation notification.
```

---

## 10. Initial Technical Vision

The first implementation should prove:

```text
service boundaries
contract-first gRPC APIs
Kafka event flow
service-owned MySQL schemas
checkout orchestration
idempotency
basic observability
database migrations
testable vertical slice
```

The initial version should be small but architecturally meaningful.

---

## 11. Out of Scope for First Version

The first version does not need to include:

```text
real payment provider integration
real shipping provider integration
full customer account management
full admin UI
advanced recommendation engine
full-text search engine
complex promotions
returns and refunds
multi-region deployment
full production security hardening
```

These can be added later as portfolio stages.

---

## 12. Success Criteria

The product vision is successful when bfstore can demonstrate:

```text
a clear business story
a complete checkout vertical slice
well-documented architecture
strong service boundaries
working gRPC contracts
working Kafka events
service-owned MySQL schemas
migration discipline
tests proving core behaviour
local development workflow
deployment path to Kubernetes
```

---

## 13. Portfolio Value

bfstore should help demonstrate that the engineer can:

```text
turn business requirements into architecture
design service boundaries
model complex product data
use MySQL professionally
design APIs and events
reason about trade-offs
document decisions
build observable services
think about platform operations
communicate with clients and reviewers
```

---

## 14. Summary

Borough Furniture Store gives bfstore a memorable identity.

It is a developer-themed furniture and homeware ecommerce platform that started with Go-inspired products and expanded into broader computer-science-themed homeware.

The project combines a charming product story with serious backend and platform engineering practices.
