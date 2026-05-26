# Architecture Documentation

This directory contains the architecture documentation for **bfstore**, ACME Ltd’s fictional online furniture store backend.

The purpose of this section is to explain **how the system is designed**, **why major design decisions were made**, **how services are separated**, **how services communicate**, **how data ownership works**, and **how the system can be deployed, observed, secured, tested, and operated**.

Architecture documentation should connect business requirements to implementation decisions.

---

## Directory Structure

```text
docs/architecture/
├── README.md
├── system-context.md
├── container-view.md
├── service-boundaries.md
├── domain-model.md
├── communication-patterns.md
├── event-driven-design.md
├── deployment-view.md
├── resilience-patterns.md
├── diagrams/
│   ├── c4/
│   ├── mermaid/
│   ├── plantuml/
│   └── drawio/
└── tradeoffs.md
```

---

## Purpose

The architecture documentation should answer:
```
What is bfstore?
What are the major parts of the system?
Why is the system designed as microservices?
Which service owns which business capability?
How do services communicate?
Where is gRPC used?
Where is Kafka used?
How is data owned?
How is the system deployed?
How does the system handle failure?
What trade-offs were accepted?
```

This documentation should guide implementation, testing, operations, and future design decisions.

---

## Recommended Reading Order

Read the architecture documents in this order:

```
1. system-context.md
2. container-view.md
3. domain-model.md
4. service-boundaries.md
5. communication-patterns.md
6. event-driven-design.md
7. deployment-view.md
8. resilience-patterns.md
9. tradeoffs.md
10. diagrams/
```

This follows the intended design flow:

```
Business context
    -> System context
    -> Domain model
    -> Service boundaries
    -> Communication model
    -> Event-driven design
    -> Deployment view
    -> Resilience and trade-offs
```

---

## Architecture Principles

bfstore follows these architecture principles.

1. Business Capability First

Services are organised around business capabilities rather than technical layers.

Good examples:
```
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

Avoid technical service boundaries such as:
```
database-service
validation-service
utility-service
common-service
```
A service should exist because it owns a meaningful business responsibility.

---

## 2. Service-Owned Data

Each service owns its own MySQL database/schema.

```
catalog-service          -> bfstore_catalog
inventory-service        -> bfstore_inventory
basket-service           -> bfstore_basket
order-service            -> bfstore_order
payment-service          -> bfstore_payment
customer-service         -> bfstore_customer
shipping-service         -> bfstore_shipping
notification-service     -> bfstore_notification
review-service           -> bfstore_review
search-service           -> bfstore_search
recommendation-service   -> bfstore_recommendation
```

Services must not directly query or modify another service’s database.

Cross-service communication must happen through:

```
gRPC APIs
Kafka events
documented protobuf contracts
```

The API is the contract.
The database is a private implementation detail.

---

## 3. Contract-First Design

Service APIs and events should be designed before implementation details become fixed.

bfstore uses:

```
Protobuf for service and event contracts
gRPC for synchronous service APIs
Kafka events for asynchronous business facts
Buf for protobuf linting, generation, and breaking-change checks
```

Contracts should describe service behaviour, not internal database tables.

---

## 4. gRPC for Immediate Decisions

Use gRPC when the caller needs an immediate response.

Examples:
```
GetProduct
ListProducts
AddBasketItem
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
```
gRPC is suitable for request/response operations where the caller needs a clear success or failure result.

---

## 5. Kafka for Business Facts

Use Kafka when a service needs to publish that something has happened.

Examples:
```
ProductCreated
StockReserved
PaymentAuthorised
ShipmentCreated
OrderCreated
NotificationRequested
ReviewCreated
SearchIndexUpdated
```
Events should represent facts that have already occurred.

Good event names:

```
OrderCreated
PaymentFailed
StockReservationExpired
```

Avoid command-like event names:

```
CreateOrder
SendNotification
ReserveStock
```

---

## 6. Prefer Clear Boundaries Over Shared Convenience

Avoid centralising business logic in shared packages or shared services.

Shared packages may contain technical utilities such as:

```
logging
configuration
gRPC interceptors
Kafka clients
OpenTelemetry helpers
error handling
test helpers
```

Shared packages should not contain business rules that belong to a specific service.

---

## 7. Observable by Default

Each service should emit:

```
structured logs
correlation IDs
request IDs
metrics
distributed traces
health checks
readiness checks
liveness checks
service-specific business metrics
```

The architecture should allow an engineer to trace a checkout request across API Gateway, Order Service, Inventory Service, Payment Service, Shipping Service, Kafka, and Notification Service.

---

## 8. Secure by Default

The architecture should support:

```
authentication at the edge
authorisation for protected actions
service-to-service identity
least privilege database access
no secrets in source control
PII-aware logging
secure configuration
dependency scanning
container scanning
SBOM generation
image signing
policy checks
```

Application-level security is documented in docs/security/.

ACME-wide security governance lives in acme-security-governance.

---

## 9. Operable by Design

Each service should be designed as something that can be deployed, monitored, debugged, rolled back, and restored.

Architecture decisions should consider:

```
deployment strategy
rollback strategy
database migration strategy
health/readiness behaviour
failure modes
runbooks
alerts
SLOs
disaster recovery
backup and restore
```

---

## Document Guide

| Document                    | Purpose                                                                                                     |
| --------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `system-context.md`         | Explains bfstore in relation to users, external systems, and the wider ACME platform estate                 |
| `container-view.md`         | Shows the major deployable/runtime parts of the system, such as services, databases, Kafka, and API Gateway |
| `domain-model.md`           | Describes the main business concepts such as Product, Basket, Order, Payment, Shipment, and Customer        |
| `service-boundaries.md`     | Defines what each service owns, what is out of scope, and how services should not overlap                   |
| `communication-patterns.md` | Explains where gRPC, Kafka, and API Gateway communication are used                                          |
| `event-driven-design.md`    | Defines event-driven architecture rules, topic design, event ownership, idempotency, retries, and replay    |
| `deployment-view.md`        | Explains how bfstore is expected to run locally and on Kubernetes                                           |
| `resilience-patterns.md`    | Documents timeouts, retries, circuit breaking, idempotency, fallback behaviour, DLQs, and failure handling  |
| `tradeoffs.md`              | Records major design trade-offs and their consequences                                                      |
| `diagrams/`                 | Stores architecture diagrams as source files and exported images where needed                               |

---

## Architecture Views

The architecture should be documented using several views.

1. System Context View

The system context view shows bfstore as a whole.

It should answer:
```
Who uses bfstore?
What external systems does it interact with?
How does it relate to the wider ACME platform?
What is inside and outside the system boundary?
```

Example context:

```
Customer
    -> bfstore API Gateway
        -> bfstore backend services

ACME Platform
    -> Kubernetes
    -> Kafka
    -> MySQL
    -> Observability
    -> GitOps
    -> Security controls
```

Expected file:

```docs/architecture/system-context.md```

---

## 3. Domain Model View

The domain model explains core business concepts.

Important domain concepts include:
```
Customer
Address
Product
Category
FurnitureVariant
InventoryItem
StockReservation
Basket
BasketItem
Order
OrderItem
Payment
PaymentAttempt
Shipment
Notification
Review
SearchIndexEntry
Recommendation
```
The domain model should explain relationships without becoming a physical database design.

Expected file:

```docs/architecture/domain-model.md```

---

## 4. Service Boundary View

The service boundary view explains ownership.

| Service                | Owns                                        | Does Not Own                                    |
| ---------------------- | ------------------------------------------- | ----------------------------------------------- |
| `catalog-service`      | Products, categories, product attributes    | Stock reservations, basket state, orders        |
| `inventory-service`    | Stock levels, reservations, warehouses      | Product descriptions, customer baskets          |
| `basket-service`       | Basket and basket items                     | Stock reservation, payment, order history       |
| `order-service`        | Orders, order items, order status           | Product catalogue, payment provider integration |
| `payment-service`      | Payment attempts, payment state             | Order lifecycle, stock                          |
| `shipping-service`     | Shipments, delivery options, tracking state | Order creation, payment                         |
| `notification-service` | Notification requests and delivery status   | Order state, customer profile ownership         |


Expected file:

```docs/architecture/service-boundaries.md```

---

## 5. Communication View

The communication view explains synchronous and asynchronous communication.

Example:

```
Frontend
    -> API Gateway

API Gateway
    -> Catalog Service
    -> Basket Service
    -> Order Service

Order Service
    -> Basket Service
    -> Inventory Service
    -> Payment Service
    -> Shipping Service

Order Service
    -> Kafka: OrderCreated

Notification Service
    <- Kafka: OrderCreated / NotificationRequested
```

Expected file:

```docs/architecture/communication-patterns.md```

---

## 6. Event-Driven View

The event-driven view explains Kafka usage.

It should cover:

```
topic naming
event ownership
event envelope
event versioning
producer responsibilities
consumer responsibilities
idempotency
ordering
retries
dead-letter queues
event replay
schema evolution
```

Expected file:

```docs/architecture/event-driven-design.md```

Detailed event contracts live in:

```docs/events/```

---

## 7. Deployment View

The deployment view explains how bfstore runs.

It should include:

```
local Docker Compose deployment
Kubernetes deployment
namespaces
service accounts
config
secrets
resource requests and limits
health checks
readiness checks
horizontal scaling
environment overlays
GitOps relationship
```

Expected file:

```docs/architecture/deployment-view.md```

The application repo may contain deployment artefacts in:

```deploy/```

The desired live Kubernetes state belongs in:

```bfstore-platform-gitops/```

The cloud infrastructure belongs in:

```bfstore-platform-infra/```

---

## 8. Resilience View

The resilience view explains how the system behaves during failures.

It should cover:

```
timeouts
retries
idempotency
circuit breakers
bulkheads
fallbacks
dead-letter queues
compensation
event replay
duplicate event handling
partial failure
dependency failure
```

Important scenarios:

```
Inventory Service unavailable during checkout
Payment Service returns failure
Shipping Service fails after payment authorisation
Kafka unavailable
Notification Service unavailable
MySQL unavailable
duplicate OrderCreated event
slow downstream service
customer retries checkout
```

Expected file:

```docs/architecture/resilience-patterns.md```

---

## High-Level Architecture

The target architecture is:

```
Customer / Frontend
        |
        v
API Gateway
        |
        +--> Auth Service
        +--> Customer Service
        +--> Catalog Service
        +--> Basket Service
        +--> Order Service
                  |
                  +--> Basket Service
                  +--> Inventory Service
                  +--> Payment Service
                  +--> Shipping Service
                  |
                  v
                Kafka
                  |
                  +--> Notification Service
                  +--> Search Service
                  +--> Recommendation Service
```

This structure supports the first major vertical slice:
```
Browse product
    -> Add to basket
    -> Checkout
    -> Reserve stock
    -> Authorise payment
    -> Create order
    -> Create shipment
    -> Publish event
    -> Send notification
```

---

## Initial Architecture Scope

The initial architecture should focus on the core checkout flow.

Initial services:
```
api-gateway
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

Initial infrastructure components:
```
MySQL
Kafka
Docker Compose
Protobuf tooling
basic observability
basic Kubernetes manifests
```

Initial quality requirements:
```
service-owned data
clear gRPC contracts
clear Kafka events
local development environment
structured logging
correlation IDs
health and readiness checks
integration tests
end-to-end checkout test
basic resilience behaviour
```

---

## Target Architecture Scope

The target architecture expands to include:
```
auth-service
customer-service
review-service
search-service
recommendation-service
advanced observability
performance testing
resilience testing
Kubernetes deployment
GitOps
supply-chain security
policy-as-code
production readiness scorecards
```

The broader platform architecture lives across the wider ACME estate:

```
mantrobuslawal/
├── bfstore
├── acme-platform-infra
├── acme-platform-gitops
├── acme-terraform-modules
├── acme-security-governance
└── acme-developer-platform
```

---

## Relationship to Wider ACME Platform

bfstore is the application backend repo.

It depends on the wider ACME platform estate for cloud runtime, deployment, security governance, and developer experience.

| Repository                 | Architecture Responsibility                                                                 |
| -------------------------- | ------------------------------------------------------------------------------------------- |
| `bfstore`                  | Application architecture, service boundaries, protobuf contracts, app deployment artefacts  |
| `bfstore-platform-infra`      | Cloud infrastructure, VPCs, Kubernetes, Kafka, MySQL platform, observability, CI/CD runners |
| `bfstore-platform-gitops`     | Desired Kubernetes state for apps, environments, policies, and platform add-ons             |
| `bfstore-terraform-modules`   | Reusable infrastructure modules                                                             |
| `bfstore-security-governance` | Zero trust, identity, secrets, policy-as-code, supply-chain standards                       |
| `bfstore-developer-platform`  | Backstage, golden paths, service templates, production readiness scorecards                 |

This separation keeps the application architecture focused while still showing how it fits into a realistic platform engineering model.

---

### Diagram Guidelines

Architecture diagrams should be stored as source files where possible.

Recommended tools:

```
Structurizr DSL    C4 architecture diagrams
Mermaid            README-friendly diagrams and simple flows
PlantUML           Sequence diagrams
draw.io            Polished cloud/network/platform diagrams
```

Suggested directory:

```
docs/architecture/diagrams/
├── c4/
├── mermaid/
├── plantuml/
└── drawio/
```

Recommended diagrams:

| Diagram                    | Purpose                                    | Suggested Tool         |
| -------------------------- | ------------------------------------------ | ---------------------- |
| System context             | Shows users, bfstore, and external systems | Structurizr            |
| Container view             | Shows services, Kafka, MySQL, API Gateway  | Structurizr            |
| Service dependency diagram | Shows service-to-service relationships     | Mermaid or Structurizr |
| Checkout sequence diagram  | Shows checkout flow across services        | PlantUML or Mermaid    |
| Event flow diagram         | Shows Kafka event publishing and consuming | Mermaid                |
| Deployment diagram         | Shows local/Kubernetes runtime layout      | Structurizr or draw.io |
| Observability diagram      | Shows logs, metrics, traces, dashboards    | draw.io or Mermaid     |
| CI/CD supply-chain diagram | Shows build, scan, sign, deploy flow       | draw.io or Mermaid     |

---

## Suggested Diagram Set

Create these diagrams first:
```
01-system-context
02-container-view
03-service-boundaries
04-checkout-sequence
05-kafka-event-flow
06-deployment-view
07-observability-view
08-secure-delivery-flow
```

Example files:
```
docs/architecture/diagrams/mermaid/checkout-sequence.mmd
docs/architecture/diagrams/mermaid/kafka-event-flow.mmd
docs/architecture/diagrams/c4/workspace.dsl
docs/architecture/diagrams/drawio/platform-overview.drawio
```

---

## Architecture Decision Records

Major architecture decisions should be recorded as ADRs in:

```adr/```

Examples:

```
adr/0001-use-microservices.md
adr/0002-use-grpc-for-service-communication.md
adr/0003-use-kafka-for-events.md
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
adr/0006-use-buf-for-protobuf.md
adr/0007-use-opentelemetry.md
adr/0008-use-contract-first-service-design.md
```

Use ADRs for decisions that affect:
```
service boundaries
communication style
database ownership
event design
deployment model
security model
observability approach
testing strategy
technology selection
```

---

## Architecture Documentation Template

Use this structure for major architecture documents.

# <Architecture Topic>

## 1. Purpose

Explain what this document covers.

## 2. Context

Explain the business or technical context.

## 3. Goals

List what the design should achieve.

## 4. Non-Goals

List what this design does not attempt to solve.

## 5. Current or Proposed Design

Describe the design.

## 6. Components Involved

List relevant services, databases, topics, tools, or platform components.

## 7. Data Flow or Request Flow

Explain how information moves through the system.

## 8. Failure Behaviour

Explain what happens when dependencies fail.

## 9. Security Considerations

Explain authentication, authorisation, secrets, data protection, and access boundaries.

## 10. Observability Considerations

Explain logs, metrics, traces, dashboards, alerts, and SLOs.

## 11. Operational Considerations

Explain deployment, rollback, scaling, maintenance, and runbooks.

## 12. Trade-Offs

Explain the advantages, disadvantages, and alternatives.

## 13. Related ADRs

Link relevant architecture decision records.

## 14. Related Documents

Link requirements, API docs, event docs, data docs, tests, and operations docs.

---

## Service Architecture Template

Use this structure when documenting an individual service architecture.

# <Service Name> Architecture

## 1. Purpose

Why does this service exist?

## 2. Business Capability

What business capability does it own?

## 3. Responsibilities

What does this service do?

## 4. Out of Scope

What does this service deliberately not do?

## 5. Owned Data

What data does this service own?

## 6. Inbound APIs

Which gRPC APIs does it expose?

## 7. Outbound Dependencies

Which services does it call?

## 8. Events Published

Which Kafka events does it publish?

## 9. Events Consumed

Which Kafka events does it consume?

## 10. Failure Modes

How can it fail, and what should happen?

## 11. Security Model

What authentication, authorisation, or data protection applies?

## 12. Observability

What logs, metrics, traces, dashboards, and alerts are required?

## 13. Deployment Notes

How is it configured, deployed, scaled, and rolled back?

## 14. Tests

Which tests validate this service?

## 15. Related Documents

Link requirements, protobuf contracts, database docs, and runbooks.

---

## Service Boundary Checklist

Use this checklist before creating or changing a service boundary.

```
Does the service own a clear business capability?
Does it have a clear reason to change independently?
Does it own its own data?
Can its API be described clearly?
Can its events be described clearly?
Does it avoid owning another service’s business rules?
Does it avoid direct database access to another service?
Does it have clear failure behaviour?
Can it be tested independently?
Can it be deployed independently?
Can it be observed independently?
```

If the answer to several of these is “no”, the boundary may need rethinking.

---

## Communication Design Checklist

Use this checklist when choosing gRPC or Kafka.

Use gRPC when:
```
the caller needs an immediate answer
the request is a command or query
the caller needs validation result
the caller must know success or failure before continuing
```

Use Kafka when:
```
something has already happened
multiple consumers may care
the producer should not block on downstream work
eventual consistency is acceptable
the event is useful for audit, projections, search, notifications, or recommendations
```

Avoid using Kafka as a hidden remote procedure call mechanism.

Avoid using synchronous gRPC for every downstream reaction.

---

## Data Ownership Checklist

Before adding data to a service, confirm:

```
Which service owns this data?
Who is allowed to write it?
Who is allowed to read it directly?
Which other services need to know about changes?
Should changes be exposed through gRPC, Kafka events, or both?
Is this source-of-truth data or a local projection?
What is the retention requirement?
Does it contain PII?
Does it need audit logging?
```

---

## Resilience Checklist

For each important flow, document:
```
timeouts
retries
idempotency keys
duplicate event handling
dead-letter queues
compensation actions
fallback behaviour
partial failure behaviour
manual recovery steps
alerts
runbooks
```
Important flows:
```
checkout
stock reservation
payment authorisation
shipment creation
notification delivery
search indexing
order cancellation
refunds
```

---

## Architecture Quality Bar

Architecture documentation is useful when it is:
```
clear
specific
decision-oriented
linked to requirements
linked to ADRs
linked to service contracts
linked to tests
honest about trade-offs
updated when design changes
```

Avoid documents that only contain generic claims.

Weak:
```
The system is scalable and secure.
```

Better:
```
The order-service can scale independently from the catalog-service because it owns its own deployment, API, database schema, metrics, and autoscaling policy. It communicates with other services through gRPC and Kafka rather than direct database access.
```

Weak:
```
Kafka improves performance.
```

Better:
```
Kafka is used for asynchronous downstream workflows such as notifications, search indexing, and recommendations so that order creation does not block on every consumer. This introduces eventual consistency and requires idempotent consumers.
```

---

## Documentation and Implementation Relationship

Architecture documentation should stay aligned with code.

When changing service design, also update related files.

Example: changing checkout orchestration may affect:
```
docs/architecture/communication-patterns.md
docs/architecture/resilience-patterns.md
docs/requirements/user-journeys.md
docs/requirements/service-requirements/order-service.md
docs/events/event-catalog.md
docs/data/data-ownership.md
proto/acme/order/v1/
proto/acme/inventory/v1/
proto/acme/payment/v1/
services/order-service/
tests/e2e/
tests/resilience/
adr/
```

Architecture drift should be treated as technical debt.

---

## Open Architecture Questions

Track unresolved design questions in the relevant document or ADR.

Initial questions:

| Question                                                                                           | Status    |
| -------------------------------------------------------------------------------------------------- | --------- |
| Will the API Gateway expose REST, GraphQL, or gRPC-Web externally?                                 | To decide |
| Should checkout orchestration live entirely in `order-service`?                                    | Proposed  |
| Should shipment creation block order confirmation?                                                 | To decide |
| Should notifications consume `OrderCreated` directly or a dedicated `NotificationRequested` event? | To decide |
| Should search initially use MySQL projections or a dedicated search engine?                        | To decide |
| Should recommendation logic be rules-based first?                                                  | Proposed  |
| Should service mesh be introduced in the first Kubernetes version?                                 | Deferred  |
| Should payment authorisation and capture be separate flows?                                        | To decide |
| Should guest checkout be supported in the first version?                                           | To decide |

---

## Relationship to Other Documentation

Architecture documentation connects to the rest of the project.

```
requirements/
    defines what the system must do

architecture/
    explains how the system is shaped

api/
    defines synchronous gRPC contracts

events/
    defines asynchronous event contracts

data/
    defines service data ownership and persistence design

testing/
    defines how the architecture is validated

security/
    defines how the application is protected

observability/
    defines how the system is monitored and diagnosed

operations/
    defines how the system is deployed, restored, and supported

adr/
    records major decisions and trade-offs
```

---

## Definition of Done for Architecture Documents

An architecture document is ready when it has:
```
a clear purpose
context
goals and non-goals
components involved
request or data flow
service ownership impact
security considerations
observability considerations
failure behaviour
operational considerations
trade-offs
related ADRs
related documents
```

A service architecture document is ready when it defines:
```
business capability
responsibilities
out of scope
owned data
inbound APIs
outbound calls
events published
events consumed
failure modes
security model
observability model
deployment notes
tests
```

---

## Getting Started

Start by creating these files:
```
system-context.md
container-view.md
domain-model.md
service-boundaries.md
communication-patterns.md
event-driven-design.md
```
Then add:
```
deployment-view.md
resilience-patterns.md
tradeoffs.md
diagrams/
```

The first diagrams to create should be:
```
system context diagram
container diagram
checkout sequence diagram
Kafka event flow diagram
deployment view
```

The first service boundaries to document should be:
```
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
api-gateway
```

These support the initial checkout vertical slice and give the project a strong architecture foundation.
