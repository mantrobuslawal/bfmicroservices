# Borough Online Furniture Store (bfstore)

**bfstore** is a cloud-native microservice backend for a fictional online furniture store operated by **ACME Ltd**.

The project is designed to demonstrate production-style backend engineering and platform-aware application design. It forms the application layer of a wider ACME platform engineering portfolio covering Kubernetes, GitOps, cloud infrastructure, DevSecOps, software supply-chain security, observability, developer experience, and production operations.

---

## Project Goals

bfstore demonstrates:

- Microservice architecture design
- Domain-led service boundaries
- gRPC-based service-to-service communication
- Protobuf-based API and event contracts
- Kafka-based asynchronous messaging
- MySQL service-owned databases
- Structured logging, metrics, and tracing
- Contract, integration, end-to-end, performance, and resilience testing
- Containerised local development
- Kubernetes-ready deployment configuration
- Secure software delivery practices
- Production readiness documentation
- Operational runbooks and failure-mode analysis

---

## Business Context

ACME Ltd operates an online furniture store where customers can browse products, manage a basket, place orders, reserve stock, make payments, arrange delivery, receive notifications, submit reviews, search the catalogue, and receive product recommendations.

The backend is designed as a microservice system because the store contains distinct business capabilities:

- Product catalogue management
- Inventory and stock reservation
- Basket management
- Customer management
- Order management
- Payment processing
- Shipping and fulfilment
- Notification delivery
- Review management
- Search
- Recommendations
- Authentication and authorisation

Each capability is owned by a dedicated service with its own API, data model, operational behaviour, and database boundary.

---

## Architecture Overview

bfstore uses a hybrid communication model:

- **gRPC** for synchronous service-to-service requests
- **Kafka** for asynchronous business events
- **Protobuf** for service contracts and event payloads
- **MySQL** for service-owned relational data stores

The design follows this principle:

> Commands that require an immediate response use gRPC. Facts that have already happened are published as Kafka events.

Example:

```text
CreateOrder           -> gRPC command
ReserveStock          -> gRPC command
AuthorisePayment      -> gRPC command

OrderCreated          -> Kafka event
StockReserved         -> Kafka event
PaymentAuthorised     -> Kafka event
ShipmentCreated       -> Kafka event
NotificationRequested -> Kafka event
```

## High-Level System Flow

The API Gateway exposes client-facing APIs and routes internal requests to backend services.

Services communicate through gRPC when a direct response is required. Services publish Kafka events when downstream services need to react asynchronously.

```text
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

## Core Services

| Service                  | Responsibility                                                           |
| ------------------------ | ------------------------------------------------------------------------ |
| `api-gateway`            | Public entry point for frontend clients                                  |
| `auth-service`           | Authentication, authorisation, tokens, user sessions                     |
| `customer-service`       | Customer profiles, addresses, preferences                                |
| `catalog-service`        | Products, categories, furniture details, pricing                         |
| `inventory-service`      | Stock levels, stock reservations, warehouse availability                 |
| `basket-service`         | Customer basket and basket items                                         |
| `order-service`          | Order creation, order lifecycle, order history                           |
| `payment-service`        | Payment authorisation, capture, refunds                                  |
| `shipping-service`       | Delivery options, shipment creation, fulfilment status, tracking updates |
| `notification-service`   | Email/SMS/event-driven customer notifications                            |
| `review-service`         | Product reviews, ratings, moderation status                              |
| `search-service`         | Product search, filtering, faceted search, search index updates          |
| `recommendation-service` | Product recommendations, related items, personalised suggestions         |

The initial implementation may focus on a smaller vertical slice first:

```text
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
api-gateway
```

## Repository Layout

```text
bfstore/
├── README.md
├── Makefile
├── docker-compose.yml
├── buf.yaml
├── buf.gen.yaml
├── .env.example
├── .gitignore
├── .editorconfig
│
├── docs/
│   ├── README.md
│   ├── requirements/
│   ├── architecture/
│   ├── api/
│   ├── events/
│   ├── data/
│   ├── testing/
│   ├── security/
│   ├── observability/
│   └── operations/
│
├── adr/
│   ├── README.md
│   ├── 0001-use-microservices.md
│   ├── 0002-use-grpc-for-service-communication.md
│   ├── 0003-use-kafka-for-events.md
│   ├── 0004-use-service-owned-databases.md
│   ├── 0005-use-mysql.md
│   ├── 0006-use-buf-for-protobuf.md
│   ├── 0007-use-opentelemetry.md
│   └── 0008-use-contract-first-service-design.md
│
├── proto/
│   └── acme/
│       ├── common/
│       ├── auth/
│       ├── customer/
│       ├── catalog/
│       ├── inventory/
│       ├── basket/
│       ├── order/
│       ├── payment/
│       ├── shipping/
│       ├── notification/
│       ├── review/
│       ├── search/
│       └── recommendation/
│
├── services/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── customer-service/
│   ├── catalog-service/
│   ├── inventory-service/
│   ├── basket-service/
│   ├── order-service/
│   ├── payment-service/
│   ├── shipping-service/
│   ├── notification-service/
│   ├── review-service/
│   ├── search-service/
│   └── recommendation-service/
│
├── packages/
│   └── go/
│       ├── logger/
│       ├── config/
│       ├── grpc/
│       ├── kafka/
│       ├── telemetry/
│       ├── auth/
│       ├── errors/
│       ├── health/
│       ├── middleware/
│       └── testkit/
│
├── db/
│   ├── README.md
│   ├── mysql-init/
│   ├── catalog/
│   ├── inventory/
│   ├── basket/
│   ├── order/
│   ├── payment/
│   ├── customer/
│   ├── shipping/
│   ├── notification/
│   ├── review/
│   ├── search/
│   └── recommendation/
│
├── deploy/
│   ├── docker/
│   ├── kubernetes/
│   ├── helm/
│   └── kustomize/
│
├── tests/
│   ├── contract/
│   ├── integration/
│   ├── e2e/
│   ├── performance/
│   ├── resilience/
│   └── testdata/
│
├── tools/
├── scripts/
└── .github/
    ├── workflows/
    ├── CODEOWNERS
    ├── dependabot.yml
    └── pull_request_template.md
```

## Repository Rationale

This repository is structured as an application monorepo.

The microservices, protobuf contracts, shared Go packages, service documentation, local development tooling, database migrations, deployment manifests, and tests live together because they are tightly related.

This makes it easier to:

- Review the full backend system in one place
- Evolve service contracts alongside implementations
- Run local integration environments
- Keep documentation close to code
- Coordinate cross-service tests
- Share internal Go packages safely
- Demonstrate the complete application architecture clearly

The wider ACME platform estate is intentionally split into separate repositories so that application, infrastructure, security governance, GitOps, reusable modules, and developer platform concerns remain cleanly separated.

## Wider ACME Platform Estate

bfstore is the application backend repository. It is part of a wider ACME Ltd platform engineering estate:

```text
mantrobuslawal/
├── bfstore-microservices
├── bfstore-platform-infra
├── bfstore-platform-gitops
├── bfstore-terraform-modules
├── bfstore-security-governance
└── bfstore-developer-platform
```

| Repository                 | Purpose                                                                              |
| -------------------------- | ------------------------------------------------------------------------------------ |
| `bfstore-mircorservices`                  | Application backend, services, protobuf, events, MySQL schemas, tests, app docs      |
| `bfstore-platform-infra`      | Cloud infrastructure, VPCs, Kubernetes, Kafka, MySQL, observability, CI/CD runners   |
| `bfstore-platform-gitops`     | Desired Kubernetes state for apps, platform add-ons, policies, and environments      |
| `bfstore-terraform-modules`   | Versioned reusable infrastructure modules                                            |
| `bfstore-security-governance` | Zero trust, identity, secrets, policy-as-code, threat models, supply-chain standards |
| `bfstore-developer-platform`  | Backstage, golden paths, production readiness scorecards, service templates          |


The relationship between the repositories is:

```text
bfstore-microservices
    Application source, service contracts, app tests, app docs
        |
        | container images deployed by
        v
bfstore-platform-gitops
    Desired Kubernetes state for apps, platform add-ons, and policies
        |
        | runs on infrastructure built by
        v
bfstore-platform-infra
    Cloud environments, VPCs, Kubernetes, Kafka, MySQL, observability, runners
        |
        | consumes reusable modules from
        v
bfstore-terraform-modules
    Versioned, tested, reusable infrastructure modules

bfstore-security-governance
    Defines standards, policies, threat models, supply-chain controls
        |
        | governs
        v
bfstore-microservices + bfstore-platform-infra + bfstore-platform-gitops + bfstore-terraform-modules

bfstore-developer-platform
    Provides Backstage, golden paths, scorecards, and service templates
        |
        | supports developers building and operating
        v
bfstore-microservices and the wider ACME platform
```

This separation reflects a realistic platform engineering model.

## Data Ownership

Each service owns its own MySQL database/schema.

```text
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

Services must not directly access another service’s database.

Services communicate through:

- gRPC APIs
- Kafka events
- Documented protobuf contracts

The API is the contract. The database is a private implementation detail owned by the service.

## Communication Model

bfstore uses two communication styles.

### Synchronous gRPC

Used when the caller needs an immediate response.

Examples:

```text
api-gateway -> catalog-service
api-gateway -> basket-service
api-gateway -> order-service
order-service -> inventory-service
order-service -> payment-service
order-service -> shipping-service
Asynchronous Kafka Events
```
### Asynchronous Kafka Events

Used when one service needs to announce that something has happened.

Examples:

```text
ProductCreated
ProductUpdated
StockReserved
StockReservationFailed
BasketCheckedOut
OrderCreated
OrderConfirmed
PaymentAuthorised
PaymentFailed
ShipmentCreated
NotificationRequested
ReviewCreated
SearchIndexUpdated
RecommendationGenerated
```
Events allow downstream services to react without tightly coupling the originating service to every consumer.

## Example Checkout Flow

1. Customer submits checkout request
2. API Gateway calls Order Service using gRPC
3. Order Service validates the request
4. Order Service requests stock reservation from Inventory Service
5. Order Service requests payment authorisation from Payment Service
6. Order Service creates the order
7. Order Service requests shipment creation from Shipping Service
8. Order Service publishes OrderCreated event
9. Notification Service consumes event and sends confirmation
10. Search and Recommendation services update downstream data where required

This flow combines synchronous and asynchronous communication.

Stock reservation, payment authorisation, and shipment creation may be synchronous because they affect whether the order can progress.

Notifications, search updates, analytics, and recommendations are asynchronous because the checkout request should not block on every downstream consumer.



