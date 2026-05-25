# Scope

## 1. Document Purpose

This document defines the scope of **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains:

- what is included in the project
- what is excluded from the project
- what will be delivered in the first version
- what may be delivered in later phases
- what assumptions and constraints shape the scope
- how scope decisions support the wider platform engineering portfolio

This document should be read alongside:

```text
docs/requirements/product-vision.md
docs/requirements/user-journeys.md
docs/requirements/functional-requirements.md
docs/requirements/non-functional-requirements.md
docs/architecture/service-boundaries.md
docs/architecture/domain-model.md
```

---

## 2. Product Scope Summary

bfstore is a backend platform for an online furniture store.

The system will support core commerce workflows including product browsing, basket management, checkout, stock reservation, payment authorisation, order creation, shipment creation, notifications, product reviews, search, and recommendations.

The first implementation will focus on a complete checkout vertical slice. Later phases will expand the platform with search, reviews, recommendations, advanced operations, Kubernetes deployment, GitOps, and DevSecOps controls.

---

## 3. Scope Statement

bfstore will provide ACME Ltd with a cloud-native microservice backend for online furniture commerce.

The project includes:

- application-level microservices
- gRPC APIs
- protobuf contracts
- Kafka event contracts
- service-owned MySQL databases
- local development environment
- application tests
- application observability
- application security documentation
- Kubernetes-ready deployment artefacts
- production readiness documentation

bfstore will not contain the full ACME cloud platform, landing zone, reusable Terraform modules, zero-trust governance, or developer platform implementation. Those concerns belong to the wider ACME platform estate.


---

## 4. In Scope

### 4.1 Application Backend
The following backend capabilities are in scope:

| Capability                       | Description                                                                                   |
| -------------------------------- | --------------------------------------------------------------------------------------------- |
| Authentication and authorisation | Customer registration, sign-in, token/session handling, protected operations                  |
| Customer management              | Customer profile, delivery addresses, preferences                                             |
| Product catalogue                | Furniture products, categories, descriptions, images, dimensions, materials, colours, pricing |
| Inventory management             | Stock levels, warehouse availability, stock reservation, reservation expiry                   |
| Basket management                | Add, update, remove, and view basket items                                                    |
| Order management                 | Create orders, view orders, update order lifecycle state                                      |
| Payment processing               | Authorise payments, record payment attempts, handle payment failures and refunds conceptually |
| Shipping and fulfilment          | Delivery options, shipment creation, fulfilment state, tracking references                    |
| Notifications                    | Order confirmation, payment status, shipment updates                                          |
| Reviews                          | Product ratings, review submission, review visibility/moderation status                       |
| Search                           | Product search, filtering, faceted search, search index update events                         |
| Recommendations                  | Related products, popular products, basic recommendation rules                                |

---

### 4.2 Core Microservices

The target service landscape is in scope.

| Service                  | Scope                                          |
| ------------------------ | ---------------------------------------------- |
| `api-gateway`            | Public entry point for frontend clients        |
| `auth-service`           | Authentication, authorisation, sessions/tokens |
| `customer-service`       | Customer profiles, addresses, preferences      |
| `catalog-service`        | Product catalogue and product metadata         |
| `inventory-service`      | Stock levels and reservations                  |
| `basket-service`         | Customer basket management                     |
| `order-service`          | Order creation and order lifecycle             |
| `payment-service`        | Payment authorisation and payment state        |
| `shipping-service`       | Shipment creation and fulfilment status        |
| `notification-service`   | Event-driven customer notifications            |
| `review-service`         | Product reviews and ratings                    |
| `search-service`         | Product search and indexing                    |
| `recommendation-service` | Product recommendations                        |


The first implementation may include only a subset of these services, but the overall architecture should allow all target services to be added.

---

### 4.3 Communication Model

The following communication patterns are in scope:

| Pattern     | Scope                                                 |
| ----------- | ----------------------------------------------------- |
| gRPC        | Internal synchronous service-to-service communication |
| Kafka       | Asynchronous event-driven workflows                   |
| Protobuf    | API and event payload contracts                       |
| API Gateway | Client-facing entry point to backend services         |

The system should clearly separate:

```text commands that require an immediate result```

from:

```text events that describe facts that have already happened```

Examples:

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

## 4.4 Data and Persistence

The following data architecture is in scope:

- MySQL as the primary relational database
- one logical database/schema per service
- one least-privilege database user per service
- versioned database migrations
- seed data for local development
- service-owned data boundaries
- no direct cross-service database access
- documented data ownership
- basic data classification and PII handling
- data retention rules

Example database ownership:

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

---

## 4.5 Local Development

The local development environment is in scope.

It should support:

- Docker Compose
- MySQL
- Kafka
- service containers where practical
- local configuration via ```.env```
- database bootstrap scripts
- database migrations
- local test data
- protobuf generation
- Makefile-based developer commands

Expected local commands:

```bash
</> Bash

make dev-up
make proto
make migrate-up
make test
make run
make dev-down

```

---

## 4.5 Local Development

The local development environment is in scope.

It should support:

- Docker Compose
- MySQL
- Kafka
- service containers where practical
- local configuration via .env
- database bootstrap scripts
- database migrations
- local test data
- protobuf generation
- Makefile-based developer commands

Expected local commands:

```bash
</> Bash

make dev-up
make proto
make migrate-up
make test
make run
make dev-down

```

---

## 4.6 Testing

The following testing areas are in scope:

| Test Type         | Scope                                                                                 |
| ----------------- | ------------------------------------------------------------------------------------- |
| Unit tests        | Validate service logic and package-level behaviour                                    |
| Integration tests | Validate services with MySQL, Kafka, and local dependencies                           |
| Contract tests    | Validate gRPC and Kafka event compatibility                                           |
| End-to-end tests  | Validate complete business journeys                                                   |
| Performance tests | Validate latency, throughput, and error rates                                         |
| Resilience tests  | Validate behaviour under dependency failures, retries, duplicate events, and restarts |

Important flows to test:

```text
browse product
add to basket
checkout
reserve stock
authorise payment
create order
create shipment
publish event
send notification

```

---

## 4.7 Observability

Application-level observability is in scope.

Each service should eventually provide:

- structured logs
- correlation IDs
- request IDs
- metrics
- distributed traces
- health checks
- readiness checks
- liveness checks
- request latency metrics
- error counts
- Kafka consumer lag metrics where relevant
- service dashboards
- service alerts

---

## 4.8 Security

Application-level security is in scope.

This includes:

- customer authentication
- authorisation for protected actions
- service-to-service security design
- least privilege database access
- secure configuration handling
- no secrets committed to source control
- PII-aware logging
- audit events for sensitive operations
- dependency scanning
- container scanning
- SBOM generation
- image signing
- application threat modelling

Wider platform security standards are documented in ```bfstore-security-governance```.

---

## 4.10 Operations Documentation

Application operations documentation is in scope.

This includes:

- runbooks
- deployment strategy
- release strategy
- rollback strategy
- database migration strategy
- incident response guidance
- backup and restore notes
- disaster recovery considerations
- cost-control considerations
- production readiness checklist

---

## 5. Out of Scope

The following are out of scope for the bfstore application repo.

### 5.1 Full Frontend Application

A production-quality frontend is out of scope.

The API Gateway may expose endpoints suitable for frontend clients, but a full customer-facing web application is not required.

Possible later options:

- minimal frontend for demos
- API client collection
- simple admin UI
- simple test harness

---

### 5.2 Real Payment Provider Integration

Real payment provider integration is out of scope for the initial implementation.

The ```payment-service``` may simulate payment authorisation, capture, failure, and refund flows.

Out of scope initially:

- Stripe live integration
- PayPal live integration
- card storage
- PCI DSS implementation
- fraud detection
- chargeback management

---

### 5.3 Real Email or SMS Provider Integration

Real notification provider integration is out of scope initially.

The ```notification-service``` may log, simulate, or write notifications to a local store.

Out of scope initially:

- real email sending
- real SMS sending
- provider webhooks
- marketing campaigns
- notification preference centre

---

### 5.4 Real Warehouse or Carrier Integration

Real fulfilment provider integration is out of scope initially.

The ```shipping-service``` may simulate shipment creation and tracking.

Out of scope initially:

- real warehouse management system integration
- real carrier APIs
- label printing
- route optimisation
- returns logistics

---

###  5.5 Advanced Commerce Features

The following are out of scope for the initial implementation:

- promotions
- vouchers
- loyalty points
- complex tax calculation
- multi-currency pricing
- internationalisation
- subscriptions
- bundles
- wishlists
- gift cards
- marketplaces
- seller onboarding

Some of these may be added later.

---

### 5.6 Advanced Search and Recommendation Engines

Basic search and recommendations are in scope.

Advanced implementations are out of scope initially.

Out of scope initially:

- Elasticsearch/OpenSearch production cluster
- vector search
- machine learning recommendation models
- real-time personalisation
- A/B testing platform
- behavioural data lake

The first implementation may use simple rules, local indexes, or event-fed projections.

---

### 5.7 Full Production Cloud Infrastructure

The bfstore repo does not own full production infrastructure.

Out of scope for this repo:

- cloud landing zone
- account/subscription/project structure
- VPC implementation
- Kubernetes cluster provisioning
- Kafka platform provisioning
- managed MySQL provisioning
- observability platform provisioning
- CI/CD runner infrastructure
- DNS and certificate platform
- cloud firewall and egress inspection

These belong in:
```bfstore-platform-infra```

---

### 5.8 Reusable Infrastructure Modules

Reusable Terraform/OpenTofu modules are out of scope for this repo.

They belong in:
```bfstore-terraform-modules```

---

### 5.9 GitOps Environment State

The bfstore repo may contain application deployment manifests, charts, and examples.

The desired live environment state belongs in:
```bfstore-platform-gitops```

This includes:

- environment-specific Argo CD applications
- cluster overlays
- production image versions
- platform add-ons
- environment policies

---

### 5.10 Platform Security Governance

Application security documentation is in scope.

ACME-wide security governance is out of scope for this repo and belongs in:
```bfstore-security-governance```

Examples:

- zero-trust strategy
- organisation-wide identity model
- secrets governance
- policy-as-code standards
- supply-chain standards
- SLSA target state
- security exceptions process
- threat model templates
- production security standards

---

### 5.11 Developer Platform Implementation

Backstage and the internal developer platform are out of scope for this repo.

They belong in:
```bfstore-developer-platform```

This includes:

- Backstage app configuration
- service catalogue
- TechDocs configuration
- golden path templates
- scorecards
- platform onboarding docs
- service creation templates

---

## 6. Initial Version Scope

The first version should focus on proving the main architecture with a working checkout vertical slice.

### 6.1 Initial Services

The first implementation should include:

- api-gateway
- catalog-service
- inventory-service
- basket-service
- order-service
- payment-service
- shipping-service
- notification-service

### 6.2 Initial Business Flow

The first version should support:

```text
Browse product
    -> Add to basket
    -> Checkout
    -> Reserve stock
    -> Authorise payment
    -> Create order
    -> Create shipment
    -> Publish OrderCreated event
    -> Send notification
```

### 6.3 Initial Technical Capabilities

The first version should include:

- protobuf contracts for the initial services
- gRPC APIs for the initial services
- Kafka event definitions for checkout-related events
- MySQL schemas for the initial services
- Docker Compose local environment
- service-level unit tests
- integration tests for key services
- basic end-to-end checkout test
- structured logs with correlation IDs
- health and readiness endpoints
- basic deployment manifests

### 6.4 Initial Documentation

The first version should include:

```text
docs/requirements/product-vision.md
docs/requirements/scope.md
docs/requirements/user-journeys.md
docs/requirements/functional-requirements.md
docs/requirements/non-functional-requirements.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/api/grpc-overview.md
docs/events/event-catalog.md
docs/data/data-ownership.md
docs/testing/testing-strategy.md
```

---

## 7. Later Phase Scope

Later phases may add the following.

### 7.1 Extended Business Capabilities

- customer profile management
- review service
- search service
- recommendation service
- order cancellation
- refunds
- shipment tracking updates
- product availability notifications
- admin workflows
- product moderation
- review moderation

### 7.2 Advanced Testing

- performance testing
- soak testing
- spike testing
- chaos testing
- fault injection
- duplicate event testing
- consumer lag testing
- database failover testing

### 7.3 Platform Integration

- Kubernetes deployment
- Helm or Kustomize
- Argo CD GitOps deployment
- environment promotion
- image version promotion
- deployment rollback
- blue/green or canary strategy

### 7.4 DevSecOps Controls

- dependency scanning
- SAST
- container scanning
- SBOM generation
- image signing
- provenance
- policy-as-code
- admission policy examples
- secrets integration
- vulnerability reporting

### 7.5 Operational Maturity

- service dashboards
- alerts
- SLOs
- runbooks
- backup and restore tests
- disaster recovery plans
- production readiness scorecards
- incident response simulations
- FinOps documentation

---

## 8. Scope by Repository

The ACME platform estate is split across several repositories.

| Repository                 | In Scope                                                                                                      |
| -------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `bfstore-microservices`                  | Application backend, services, protobuf, events, MySQL schemas, app tests, app docs, app deployment artefacts |
| `bfstore-platform-infra`      | Cloud environments, networking, Kubernetes, Kafka, MySQL platform, observability, CI/CD runners, DR, FinOps   |
| `bfstore-platform-gitops`     | Desired Kubernetes state for apps, platform add-ons, policies, and environment promotion                      |
| `bfstore-terraform-modules`   | Reusable infrastructure modules, module tests, examples, release process                                      |
| `bfstore-security-governance` | Zero trust, identity, secrets, policy-as-code, threat models, supply-chain standards                          |
| `bfstore-developer-platform`  | Backstage, golden paths, scorecards, service templates, onboarding                                            |


This repository split is intentional. It keeps bfstore focused on the application while allowing the wider ACME estate to demonstrate senior platform engineering concerns.

---

## 9. Scope by Service

### 9.1 API Gateway

In scope:

- client-facing entry point
- request routing to backend services
- authentication enforcement
- request correlation IDs
- response mapping
- basic rate limiting design
- error response standardisation

Out of scope initially:

- full API monetisation
- external developer portal
- complex API product management

---

### 9.2 Auth Service

In scope:

- customer registration
- customer login
- token issuing
- token validation
- password hashing
- basic role/permission model

Out of scope initially:

- social login
- enterprise SSO
- MFA
- passkeys
- account recovery workflows
- fraud detection

---

### 9.3 Customer Service

In scope:

- customer profile
- delivery addresses
- customer preferences

Out of scope initially:

- marketing preferences centre
- GDPR subject access automation
- loyalty account management

---

### 9.4 Catalogue Service

In scope:

- products
- categories
- furniture dimensions
- materials
- colours
- product images metadata
- product pricing
- active/inactive product status

Out of scope initially:

- full PIM integration
- supplier catalogue import
- complex pricing engine

---

### 9.5 Inventory Service

In scope:

- stock levels
- stock reservations
- reservation expiry
- stock adjustment events
- warehouse availability model

Out of scope initially:

- real-time warehouse integration
- purchase ordering
- stock forecasting

---

### 9.6 Basket Service

In scope:

- create basket
- add item
- update quantity
- remove item
- view basket
- clear basket after checkout

Out of scope initially:

- guest basket merge
- wishlist conversion
- saved baskets

---

### 9.7 Order Service

In scope:

- create order
- validate checkout
- coordinate stock reservation and payment authorisation
- maintain order lifecycle state
- publish order events
- view order details
- list customer orders

Out of scope initially:

- complex returns
- exchanges
- manual order amendments
- split shipments

---

### 9.8 Payment Service

In scope:

- simulate payment authorisation
- record payment attempts
- handle payment success/failure
- publish payment events
- conceptual refund support

Out of scope initially:

- live payment provider
- card vaulting
- PCI DSS scope
- 3D Secure
- chargebacks

---

### 9.9 Shipping Service

In scope:

- delivery option selection
- shipment creation
- shipment status
- tracking reference
- shipment events

Out of scope initially:

- live carrier integration
- label generation
- route optimisation
- returns logistics

---

### 9.10 Notification Service

In scope:

- consume notification-related events
- send simulated email/SMS notifications
- record notification status
- retry failed notifications conceptually

Out of scope initially:

- real provider integration
- marketing campaigns
- preference centre

---

### 9.11 Review Service

In scope:

- submit review
- view product reviews
- rating summary
- review moderation status

Out of scope initially:

- fraud detection
- media uploads
- sentiment analysis
- abuse detection

---

### 9.12 Search Service

In scope:

- product search
- filtering
- faceted search
- search index update events

Out of scope initially:

- production OpenSearch cluster
- vector search
- semantic search
- personalised ranking

---

### 9.13 Recommendation Service

In scope:

- related products
- popular products
- basic rules-based recommendations
- recommendation events

Out of scope initially:

- machine learning model training
- feature store
- real-time personalisation
-  A/B testing

---

## 10. Non-Functional Scope

The following quality attributes are in scope.

| Category        | Scope                                                                               |
| --------------- | ----------------------------------------------------------------------------------- |
| Performance     | Define latency, throughput, and error-rate targets for key flows                    |
| Reliability     | Design for retries, idempotency, and graceful degradation                           |
| Availability    | Provide health/readiness checks and deployment readiness                            |
| Scalability     | Design services so they can scale independently                                     |
| Security        | Apply authentication, authorisation, least privilege, and secure delivery practices |
| Observability   | Provide logs, metrics, traces, dashboards, and alerts                               |
| Maintainability | Use clear service boundaries, docs, tests, and ADRs                                 |
| Operability     | Provide runbooks, deployment docs, and rollback strategy                            |
| Resilience      | Test behaviour under common dependency failures                                     |
| Cost awareness  | Include local and cloud cost-control considerations                                 |


---

## 11. Scope Control Principles

Scope decisions should follow these principles:

### 11.1 Build a Working Vertical Slice First

Prefer one complete business flow over many incomplete services.

The first major target is:

```browse -> basket -> checkout -> stock -> payment -> order -> shipment -> notification```

### 11.2 Avoid Premature Platform Complexity

Do not block application progress by implementing every platform capability first.

Build enough platform support to run and test the app, then mature the platform iteratively.

### 11.3 Keep Service Boundaries Clear

Avoid adding “shared” services that become distributed utility layers.

Each service should own a business capability.

### 11.4 Prefer Mocked External Integrations Initially

Mock payments, notifications, shipping, and external systems first.

Real integrations can be added later if they strengthen the portfolio.

### 11.5 Document Deferred Work

If a feature is intentionally deferred, record it rather than silently ignoring it.

---

## 12. Assumptions

This scope assumes:

- ACME Ltd is fictional
- bfstore is primarily a backend/platform portfolio project
- Go will be used for backend services
- MySQL will be used for relational persistence
- Kafka will be used for asynchronous messaging
- gRPC and protobuf will be used for internal service contracts
- Docker Compose will support local development
- Kubernetes will be the target deployment platform
- cloud infrastructure will be handled in a separate repo
- external third-party services may be mocked initially
- the first version should prioritise a working checkout flow

---

## 13. Constraints

The project is constrained by:

- solo developer effort
- need to keep local development affordable
- desire to use mostly open-source tooling
- portfolio value for Senior Platform Engineer, DevSecOps Engineer, and Kubernetes Platform Engineer roles
- need to keep first delivery achievable
- need to avoid unnecessary enterprise complexity before the core system works

---

## 14. Dependencies

bfstore depends on:

| Dependency     | Purpose                                                          |
| -------------- | ---------------------------------------------------------------- |
| Go             | Primary service implementation language                          |
| MySQL          | Service-owned relational databases                               |
| Kafka          | Asynchronous event messaging                                     |
| Protobuf       | API and event contract definition                                |
| gRPC           | Internal service communication                                   |
| Docker Compose | Local development environment                                    |
| Kubernetes     | Target runtime platform                                          |
| Buf            | Protobuf linting, breaking change detection, and code generation |
| OpenTelemetry  | Distributed tracing and telemetry instrumentation                |

Wider platform dependencies may include:

| Dependency                    | Purpose                |
| ----------------------------- | ---------------------- |
| Argo CD                       | GitOps deployment      |
| Terraform/OpenTofu            | Infrastructure as code |
| Backstage                     | Developer portal       |
| Kyverno/OPA                   | Policy-as-code         |
| Cosign                        | Image signing          |
| Syft/Trivy/Grype              | SBOM and scanning      |
| Prometheus/Grafana/Loki/Tempo | Observability          |

---

## 15. Deliverables

### 15.1 Application Deliverables

- service source code
- protobuf contracts
- Kafka event definitions
- MySQL migrations
- Dockerfiles
- Docker Compose setup
- Kubernetes manifests/charts
- tests
- documentation
- ADRs
  
### 15.2 Documentation Deliverables

- requirements documentation
- architecture documentation
- API documentation
- event documentation
- data ownership documentation
- testing strategy
- security documentation
- observability documentation
- operations documentation
- runbooks
- production readiness checklist

### 15.3 Platform-Related Deliverables

Within this repo:

- app deployment artefacts
- service security notes
- application supply-chain pipeline notes
- Kubernetes workload requirements

Outside this repo:

- cloud infra
- GitOps state
- reusable modules
- platform security governance
- developer platform

---













