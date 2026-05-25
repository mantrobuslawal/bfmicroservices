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

### Example:

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
---
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

---

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

---

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

---

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

---

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

---

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

---

## Communication Model

bfstore uses two communication styles.

### Synchronous gRPC

Used when the caller needs an immediate response.

### Examples:

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

### Examples:

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

### Example Checkout Flow

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

---

## Documentation

The ```text docs/ directory``` is the main source of project documentation.
```text
docs/
├── requirements/
├── architecture/
├── api/
├── events/
├── data/
├── testing/
├── security/
├── observability/
└── operations/
```

---

## Requirements

```text docs/requirements/``` explains what the system must do.

It includes:

- Product vision
- Scope
- User journeys
- Functional requirements
- Non-functional requirements
- Business rules
- Acceptance criteria
- Service-level requirements

---

## Architecture

```text docs/architecture/``` explains how the system is designed.

It includes:

- System context
- Container view
- Service boundaries
- Domain model
- Communication patterns
- Event-driven design
- Deployment view
- Resilience patterns
- Architecture diagrams
- Trade-offs

---

## API Documentation

```text docs/api/``` explains the synchronous API model.

It includes:

- gRPC design
- Protobuf style guide
- API Gateway behaviour
- Error model
- Authentication
- Versioning strategy

---

## Event Documentation

```text docs/events/``` explains the asynchronous messaging model.

It includes:

- Kafka topic design
- Event envelope
- Event catalogue
- Event versioning
- Ordering and idempotency
- Retry and dead-letter queue strategy
- Event replay strategy
- Consumer contracts
- Data Documentation

---

## Data Documentation

```text docs/data/``` explains service data ownership and persistence design.

It includes:

- Data ownership
- MySQL standards
- Service database design
- Migration strategy
- Consistency model
- Data classification
- Retention
- PII handling

---

## Testing Documentation

```text docs/testing/``` defines the testing strategy.

It covers:

- Unit testing
- Contract testing
- Integration testing
- End-to-end testing
- Performance testing
- Resilience testing
- Chaos testing
- Test environments

---

## Security Documentation

```text docs/security/``` explains application-level security.

It covers:

- Threat modelling
- Authentication
- Authorisation
- Service-to-service security
- Secrets management
- Privacy and PII
- Audit events
- Secure coding
- Supply-chain security

Platform-wide zero trust, identity governance, CI/CD hardening, policy-as-code, and software supply-chain standards are defined in ```text bfstore-security-governance```.

---

## Observability Documentation

```text docs/observability/``` explains how the system is monitored and diagnosed.

It covers:

- Structured logging
- Metrics
- Distributed tracing
- Dashboards
- Alerts
- SLOs
- Kafka consumer lag

---

## Operations Documentation

```text docs/operations/``` explains how the system is released, operated, restored, and supported.

It covers:

- Runbooks
- Deployment strategy
- Release strategy
- Rollback strategy
- Database migration strategy
- Incident response
- Disaster recovery
- Backup and restore
- Cost controls
- Production readiness

---

## Architecture Decision Records

The ```text adr/``` directory contains Architecture Decision Records.

ADRs document significant decisions and their trade-offs.

### Examples:
```text
0001-use-microservices.md
0002-use-grpc-for-service-communication.md
0003-use-kafka-for-events.md
0004-use-service-owned-databases.md
0005-use-mysql.md
0006-use-buf-for-protobuf.md
0007-use-opentelemetry.md
0008-use-contract-first-service-design.md
```
ADRs explain not only what was chosen, but why it was chosen.

---

## Protobuf Contracts

The ```text proto/``` directory contains shared protobuf definitions for gRPC APIs and Kafka event payloads.

### Example layout:

```text
proto/acme/
├── common/v1/
├── auth/v1/
├── customer/v1/
├── catalog/v1/
├── inventory/v1/
├── basket/v1/
├── order/v1/
├── payment/v1/
├── shipping/v1/
├── notification/v1/
├── review/v1/
├── search/v1/
└── recommendation/v1/
```

Protobuf is used because it provides strongly typed contracts, supports code generation, and works well with gRPC and event-driven systems.

Generated code should not be edited manually.

---

## Service Design Principles

Each service should:

1. Own a clear business capability
2. Expose a documented gRPC API
3. Publish and consume documented events
4. Own its own MySQL database/schema
5. Avoid direct database access across services
6. Emit structured logs, metrics, and traces
7. Include health, readiness, and liveness checks
8. Include unit, integration, and contract tests
9. Be independently buildable and deployable
10. Document operational behaviour and failure modes
11. Define timeouts, retries, and idempotency rules where needed
12. Provide a runbook and production readiness evidence

---

## Local Development

The local development environment is intended to run through Docker Compose.

Expected local dependencies include:

- MySQL
- Kafka
- Protobuf tooling
- Optional observability components
- Application services

Typical commands:
```bash
</> Bash

make dev-up
make proto
make migrate-up
make test
make run
make dev-down
```

The exact commands may evolve as the project implementation matures.

---

## Make Targets

The root ```text Makefile``` exposes the following common development commands.

```bash
</> Bash

make help          # Show available commands
make dev-up        # Start local dependencies with Docker Compose
make dev-down      # Stop local dependencies
make proto         # Generate protobuf/gRPC code
make lint          # Run linters
make test          # Run unit tests
make test-int      # Run integration tests
make test-e2e      # Run end-to-end tests
make test-perf     # Run performance tests
make test-res      # Run resilience tests
make migrate-up    # Apply database migrations
make migrate-down  # Roll back database migrations
make run           # Run services locally
make build         # Build all services

```

---

## Database Migrations

Each service owns its own migrations under:

```text db/<service-name>/migrations/```
### Example
```text
db/order/migrations/
db/catalog/migrations/
db/inventory/migrations/
db/payment/migrations/
```

Local MySQL bootstrap scripts live under:

```text db/mysql-init/```

The local setup creates one logical database/schema per service and one least-privilege database user per service.

---

## Testing Strategy

bfstore uses several layers of testing.

```text
Unit tests
    Validate individual functions and packages.

Service integration tests
    Validate a service with its database, Kafka, and dependencies.

Contract tests
    Validate protobuf, gRPC, and event compatibility.

End-to-end tests
    Validate complete user journeys across services.

Performance tests
    Validate latency, throughput, error rate, and behaviour under load.

Resilience tests
    Validate behaviour during dependency failures, retries, timeouts, duplicate events, and restarts.
```

Important resilience scenarios include:

- Kafka unavailable
- MySQL unavailable
- Payment service slow
- Inventory reservation failure
- Duplicate Kafka events
- Kafka consumer lag
- Service restart during checkout
- Network latency between services

---

## Observability

Each service should emit:
- Structured logs
- Metrics
- Distributed traces
- Health check status
- Readiness status
- Error counts
- Request latency
- Kafka consumer lag where applicable

The observability goal is to make it possible to answer:
```text
Is the service healthy?
Is it fast enough?
Is it producing errors?
Where is latency being introduced?
Are Kafka consumers keeping up?
Are downstream dependencies causing failures?
Can we diagnose a failed checkout flow end to end?
```

---

## Security Principles

bfstore should follow secure-by-default application practices:
- Authentication at the edge
- Authorisation for protected operations
- Service-to-service identity
- Least privilege database access
- No secrets in source control
- Secure configuration handling
- Dependency scanning
- Container scanning
- SBOM generation
- Signed container images
- Structured audit events for sensitive operations
- PII-aware logging and data handling

Application-level security lives in this repo.

Platform-wide controls live in ```text bfstore-security-governance```, including:

- Zero trust principles
- Identity and access model
- Secrets strategy
- Policy-as-code
- Threat models
- Kubernetes security standards
- Container security standards
- Software supply-chain standards
- SLSA target state
- Production readiness standard

---

### Deployment

This repo includes application deployment configuration for local and Kubernetes-based environments.

```text
deploy/
├── docker/
├── kubernetes/
├── helm/
└── kustomize/
```

The application should be deployable locally for development and to Kubernetes environments such as:

```text
dev
test
staging
prod
```

The production cloud infrastructure itself is managed outside this repo in ```text bfstore-platform-infra```.

The desired Kubernetes state for applications, platform add-ons, and policies is managed in ```text bfstore-platform-gitops```.

---

## Production Readiness

bfstore provides production readiness evidence for each service.

Each service should eventually have:

- Service owner
- README
- API documentation
- Event documentation
- Runbook
- Health checks
- Readiness checks
- Metrics
- Logs
- Traces
- Dashboard
- Alerts
- SLOs
- Threat model
- SBOM
- Signed image
- Resource requests and limits
- Deployment strategy
- Rollback strategy
- Database migration strategy
- Performance test evidence
- Resilience test evidence

Production readiness scorecards are managed in the wider developer platform repository:

```text  bfstore-developer-platform/scorecards/```

---

## Engineering Workflow

A typical change will update all affected parts of the system.

For example, changing order creation behaviour may require updates to:

```text
services/order-service/
proto/acme/order/v1/order_service.proto
proto/acme/order/v1/order_events.proto
docs/requirements/service-requirements/order-service.md
docs/events/event-catalog.md
docs/testing/contract-testing.md
docs/architecture/resilience-patterns.md
adr/
tests/contract/
tests/integration/
tests/resilience/
```
This keeps requirements, contracts, implementation, tests, and operations aligned.

---

## Design Philosophy

bfstore follows a documentation-first and contract-first approach.

The preferred design order is:

```text
Requirements
    -> Domain model
    -> Service boundaries
    -> API and event contracts
    -> Data ownership
    -> Database design
    -> Implementation
    -> Tests
    -> Deployment
    -> Operations
```

This prevents the system from being shaped too early by database tables or framework choices.

The API is the contract.

The database is an implementation detail owned by the service.

---

### Senior Platform Engineering Scope

Although bfstore is the application repo, it is designed to support a wider senior-platform portfolio.

The complete ACME estate demonstrates:

- Cloud-native application design
- Kubernetes platform design
- GitOps deployment
- Infrastructure as code
- Reusable Terraform/OpenTofu modules
- Zero-trust networking
- Identity and access governance
- Secrets management
- Policy-as-code
- Secure CI/CD
- Software supply-chain security
- SBOMs, signing, and provenance
- Observability and SLOs
- Incident response
- Disaster recovery
- Backup and restore
- FinOps and cost controls
- Developer platform and golden paths
- Production readiness scorecards

This makes the project suitable for demonstrating skills relevant to:

- Senior Platform Engineer
- DevSecOps Engineer
- Kubernetes Platform Engineer
- Cloud Platform Engineer
- Site Reliability Engineer

---

## Status

This project is under active design and development.

Current focus areas:

- Requirements analysis
- Domain modelling
- Service boundary definition
- Protobuf contract design
- MySQL data ownership model
- Local development environment
- Core service implementation
- Testing strategy
- Observability baseline
- Resilience and production readiness planning

---

## Licence

This project is intended for learning, demonstration, and portfolio purposes.

Licence to be decided.
