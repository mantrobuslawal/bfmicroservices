# Product Vision

## 1. Document Purpose

This document defines the product vision for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains:

- what the product is
- why it exists
- who it serves
- what business value it should provide
- what capabilities it should support
- what success looks like
- how the product fits into the wider ACME platform estate

This document should guide requirements, architecture, service boundaries, API design, event design, data ownership, testing, and operational planning.

---

## 2. Product Name

**bfstore**

bfstore is the backend platform for ACME Ltd’s online furniture store.

---

## 3. Business Context

ACME Ltd is a fictional furniture retailer that sells household furniture online.

The company wants to provide customers with a reliable digital shopping experience where they can browse furniture, check availability, manage a basket, place orders, pay securely, arrange delivery, receive notifications, leave reviews, and receive product recommendations.

The backend system must support the core commerce workflows required to run an online furniture store while also demonstrating modern engineering practices such as microservices, event-driven architecture, gRPC APIs, Kafka messaging, service-owned databases, observability, secure delivery, and Kubernetes-ready deployment.

---

## 4. Product Vision Statement

bfstore will provide ACME Ltd with a scalable, secure, observable, and maintainable backend platform for online furniture commerce.

The platform will allow customers to browse products, manage baskets, place orders, reserve stock, authorise payments, arrange shipping, and receive notifications, while giving engineering teams a clear service-based architecture that can evolve over time.

The system should be designed as a realistic cloud-native backend that demonstrates production-style software engineering, platform engineering, DevSecOps, and operational readiness.

---

## 5. Product Goals

The main goals of bfstore are:

| Goal | Description |
|---|---|
| Support online furniture shopping | Enable customers to browse, search, review, and purchase furniture online |
| Provide reliable checkout | Ensure basket, stock reservation, payment, order creation, shipping, and notification workflows operate correctly |
| Use clear service ownership | Split the platform into business-aligned services with clear responsibilities |
| Support event-driven workflows | Use Kafka events for asynchronous communication and downstream processing |
| Use contract-first communication | Use protobuf and gRPC for typed service APIs |
| Maintain service-owned data | Ensure each service owns its own MySQL database/schema |
| Improve operability | Provide logs, metrics, traces, health checks, runbooks, and SLOs |
| Support secure delivery | Include secure coding, dependency scanning, container scanning, SBOMs, image signing, and supply-chain controls |
| Support Kubernetes deployment | Provide deployment artefacts suitable for Kubernetes and GitOps workflows |
| Demonstrate senior platform engineering capability | Show realistic design decisions, trade-offs, testing, resilience, and operational readiness |

---

## 6. Problem Statement

ACME Ltd needs a backend platform capable of supporting online furniture sales.

A simple monolithic application may be easier to start with, but it does not demonstrate the architecture, deployment, observability, security, and operational practices ACME wants to evaluate.

The challenge is to design a backend that:

- separates business capabilities clearly
- allows services to evolve independently
- supports synchronous and asynchronous workflows
- avoids shared database coupling
- can be deployed and observed in a Kubernetes environment
- supports secure software delivery practices
- provides enough operational documentation to be credible as a production-style system

---

## 7. Target Users

bfstore serves several user groups.

### 7.1 Customers

Customers use the store to browse products, manage baskets, place orders, track deliveries, receive notifications, and submit reviews.

Customer goals:

- find suitable furniture
- view accurate product information
- understand price, dimensions, material, colour, and availability
- complete checkout reliably
- receive confirmation and delivery updates
- review purchased products

### 7.2 ACME Operations Team

The operations team needs visibility into orders, payments, stock, fulfilment, and service health.

Operations goals:

- monitor order and checkout flows
- identify failed payments, stock reservation failures, or shipping issues
- diagnose service incidents
- restore services when failures occur
- understand system health through dashboards and alerts

### 7.3 ACME Engineering Team

The engineering team builds, tests, deploys, and operates the services.

Engineering goals:

- work with clear service boundaries
- use reliable API and event contracts
- deploy services independently
- test services at multiple levels
- observe behaviour in local and Kubernetes environments
- follow secure software delivery practices

### 7.4 ACME Platform Team

The platform team provides the Kubernetes, CI/CD, GitOps, observability, security, and developer platform capabilities that support bfstore.

Platform goals:

- provide standard deployment patterns
- support self-service service creation
- enforce production readiness standards
- provide observability and incident response tooling
- secure the software supply chain
- manage environment promotion and infrastructure standards

---

## 8. Customer Value

bfstore should provide customers with:

- a reliable product browsing experience
- accurate product and availability information
- a clear basket and checkout process
- reliable order creation
- secure payment handling
- delivery and order notifications
- search and filtering
- product reviews
- useful recommendations

The customer experience should feel consistent even though the backend is implemented using multiple services.

---

## 9. Business Value

bfstore should provide ACME Ltd with:

- a maintainable commerce backend
- clearly separated business capabilities
- improved ability to evolve services independently
- better observability into customer journeys
- stronger reliability around checkout and fulfilment
- a foundation for future capabilities such as search, recommendations, reviews, promotions, and analytics
- a realistic reference architecture for platform engineering and DevSecOps practices

---

## 10. Core Business Capabilities

bfstore should support the following business capabilities.

| Capability | Description | Primary Service |
|---|---|---|
| Authentication and authorisation | Allow users to register, sign in, and access protected actions | `auth-service` |
| Customer management | Manage customer profile, addresses, and preferences | `customer-service` |
| Product catalogue | Manage furniture products, categories, attributes, images, and prices | `catalog-service` |
| Inventory management | Track stock levels, warehouses, and stock reservations | `inventory-service` |
| Basket management | Allow customers to add, update, remove, and view basket items | `basket-service` |
| Order management | Create and manage customer orders | `order-service` |
| Payment processing | Authorise, capture, and refund payments | `payment-service` |
| Shipping and fulfilment | Create shipments, manage delivery options, and track fulfilment status | `shipping-service` |
| Notifications | Send customer-facing order, payment, and delivery notifications | `notification-service` |
| Reviews | Allow customers to submit and view product reviews | `review-service` |
| Search | Support product search, filtering, and faceted search | `search-service` |
| Recommendations | Suggest related, popular, or personalised products | `recommendation-service` |

---

## 11. Product Scope Summary

### 11.1 In Scope

The product should support:

- customer registration and authentication
- customer profile and address management
- browsing furniture products
- viewing product details
- searching and filtering products
- managing a basket
- checking stock availability
- reserving stock during checkout
- creating orders
- authorising payments
- creating shipments
- sending order notifications
- submitting and viewing product reviews
- generating product recommendations
- service-level observability
- service-owned MySQL databases
- gRPC APIs between services
- Kafka events for asynchronous workflows
- local development with Docker Compose
- Kubernetes-ready deployment artefacts
- documentation, tests, runbooks, and production readiness evidence

### 11.2 Out of Scope for Initial Version

The initial version does not need to support:

- real payment provider integration
- real email or SMS provider integration
- real warehouse management system integration
- multi-currency pricing
- international tax calculation
- complex promotions or discount rules
- returns management
- fraud detection
- advanced machine learning recommendations
- a full frontend application
- production cloud deployment on day one

These may be added in later phases.

---

## 12. Initial Vertical Slice

The first working version should focus on a complete checkout path.

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

---

## 13. Product Principles

bfstore should follow these principles.

### 13.1 Business Capability First

Services should be organised around business capabilities, not technical layers.

Good:

```text
catalog-service
inventory-service
order-service
payment-service
```

Avoid:

```text
database-service
validation-service
common-service
```

### 13.2 Contract First

Service contracts should be defined clearly using protobuf before implementation becomes tightly coupled.

The protobuf contracts should describe service behaviour, not database tables.

### 13.3 Service-Owned Data

Each service should own its own MySQL database/schema.

Services must not directly query another service’s database.

### 13.4 gRPC for Commands

Use gRPC when a service needs an immediate response.

Examples:

```text
GetProduct
AddBasketItem
CreateOrder
ReserveStock
AuthorisePayment
CreateShipment
```

### 13.6 Observable by Default

Each service should emit logs, metrics, traces, health checks, and readiness checks.

### 13.7 Secure by Default

The system should avoid insecure defaults, protect customer data, use least privilege access, and support secure delivery practices.

### 13.8 Operable by Design

Services should include runbooks, failure-mode documentation, dashboards, alerts, and production readiness evidence.

---

## 14. Success Criteria

The product is successful when it demonstrates both business functionality and engineering maturity.

### 14.1 Functional Success

bfstore should successfully support:

- browsing products
- adding products to a basket
- creating an order
- reserving stock
- authorising payment
- creating a shipment
- publishing order events
- sending customer notifications

### 14.2 Architecture Success

The architecture should demonstrate:

- clear service boundaries
- service-owned databases
- well-defined gRPC APIs
- well-defined Kafka events
- protobuf contract versioning
- separation of synchronous and asynchronous workflows
- no cross-service database access

### 14.3 Engineering Success

The implementation should demonstrate:

- local development environment
- automated tests
- contract tests
- integration tests
- end-to-end tests
- performance tests
- resilience tests
- CI/CD pipeline
- Kubernetes deployment artefacts

### 14.4 Operational Success

The system should provide:

- structured logs
- metrics
- distributed traces
- dashboards
- alerts
- SLOs
- runbooks
- rollback strategy
- backup and restore approach
- incident response documentation

### 14.5 Security Success

The system should demonstrate:

- authentication and authorisation
- least privilege database users
- secret handling
- no sensitive data in logs
- dependency scanning
- container scanning
- SBOM generation
- image signing
- policy checks

---

## 15. Non-Goals

bfstore is <ins>NOT<ins> intended to be:

- a complete commercial furniture store
- a real payment processing system
- a real fulfilment or logistics platform
- a frontend design project
- a machine learning project
- a replacement for a real enterprise commerce platform
- a production SaaS business

> [!IMPORTANT]
> **The goal is to create a realistic backend and platform engineering portfolio project.**

---

## 16. Key Assumptions

The product vision assumes:

- ACME Ltd is a fictional company
- customer-facing frontend clients will communicate through the API Gateway
- internal services will communicate through gRPC and Kafka
- protobuf will be used for service and event contracts
- MySQL will be used as the primary relational database
- each service will own its own database/schema
- Kafka will be used for asynchronous business events
- the application will be containerised
- Kubernetes will be the target runtime platform
- local development will use Docker Compose
- real third-party integrations may be mocked or simulated initially

---

## 17. Key Constraints

The product is constrained by:

- portfolio and learning project scope
- solo development effort
- need to keep the first vertical slice manageable
- cost limits for cloud infrastructure
- use of open-source or affordable tooling where possible
- preference for technologies relevant to Senior Platform Engineer, DevSecOps Engineer, and Kubernetes Platform Engineer roles

---

## 18. High-Level Risks

| Risk                                              | Impact                                   | Mitigation                                                                  |
| ------------------------------------------------- | ---------------------------------------- | --------------------------------------------------------------------------- |
| Scope becomes too large                           | Project may become difficult to complete | Build in stages and focus on the checkout vertical slice first              |
| Too many services too early                       | Complexity may slow implementation       | Start with core services, then expand                                       |
| Documentation becomes disconnected from code      | Repo becomes misleading                  | Update docs as part of each service change                                  |
| Event-driven flows become hard to reason about    | Debugging and testing become difficult   | Use clear event contracts, correlation IDs, tracing, and idempotency        |
| Local development becomes too heavy               | Harder to run the project                | Keep local Compose profile lightweight                                      |
| Security controls become theoretical only         | Weak DevSecOps evidence                  | Implement scanning, SBOMs, signing, and policy checks in CI                 |
| Kubernetes work overshadows application behaviour | App may remain incomplete                | Complete working application flows before over-investing in platform layers |

---

## 19. Product Roadmap

### Phase 1: Product and Architecture Foundation

- Product vision
- Scope
- User journeys
- Functional requirements
- Non-functional requirements
- Domain model
- Service boundaries
- API and event strategy
- Data ownership model

### Phase 2: Core Backend Vertical Slice

- Catalogue browsing
- Basket management
- Checkout
- Stock reservation
- Payment authorisation
- Order creation
- Shipment creation
- Notification event

### Phase 3: Testing and Observability

- Unit tests
- Integration tests
- Contract tests
- End-to-end tests
- Structured logging
- Metrics
- Distributed tracing
- Dashboards and alerts

### Phase 4: Kubernetes and GitOps

- Kubernetes manifests
- Helm or Kustomize
- Argo CD GitOps deployment
- Environment overlays
- Rollback strategy
- Health and readiness checks

### Phase 5: DevSecOps and Supply Chain

- Dependency scanning
- Container scanning
- SBOM generation
- Image signing
- Policy-as-code
- Secrets strategy
- Threat modelling

### Phase 6: Platform Engineering Estate

- Cloud infrastructure repo
- Terraform/OpenTofu modules
- Zero-trust VPC design
- Developer platform
- Backstage service catalogue
- Golden paths
- Production readiness scorecards

### Phase 7: Advanced Capabilities

- Search
- Reviews
- Recommendations
- Performance testing
- Resilience testing
- Disaster recovery
- FinOps
- Operational maturity

---

##  20. Relationship to Wider ACME BfStore Platform

bfstore is the application backend repo.

It fits into the wider ACME BfStore platform estate:

```text
mantrobuslawal/
├── bfstore-microservices
├── bfstore-platform-infra
├── bfstore-platform-gitops
├── bfstore-terraform-modules
├── bfstore-security-governance
└── bfstore-developer-platform
```

| Repository                 | Relationship to bfstore                                                                                             |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `bfstore-platform-infra`      | Provides cloud infrastructure, Kubernetes, networking, Kafka, MySQL, observability, and CI/CD runner infrastructure |
| `bfstore-platform-gitops`     | Deploys bfstore services and platform add-ons into Kubernetes environments                                          |
| `bfstore-terraform-modules`   | Provides reusable infrastructure modules used by platform infrastructure                                            |
| `bfstore-security-governance` | Defines security standards, zero trust, identity, secrets, policy-as-code, and supply-chain requirements            |
| `bfstore-developer-platform`  | Provides Backstage, golden paths, templates, and production readiness scorecards                                    |


bfstore will focuse on application code, contracts, application-level documentation, tests, and deployment artefacts.

---

## 21. Design Trade-Offs

### Microservices vs Monolith

Microservices are used to demonstrate service boundaries, independent deployability, gRPC APIs, Kafka events, and platform engineering practices.

Trade-off:

- increased complexity
- more infrastructure requirements
- harder local development
- more testing requirements

This is acceptable because the project goal is to demonstrate cloud-native and platform engineering skills.

### gRPC vs REST Internally

gRPC is used for internal service communication because it provides strong contracts, code generation, and efficient service-to-service communication.

The API Gateway may expose REST or another client-friendly interface externally.

### Kafka vs Direct Service Calls

Kafka is used for asynchronous workflows where services need to react to business facts without tightly coupling to the originating service.

This improves decoupling but requires careful event design, idempotency, retry handling, and observability.

### MySQL vs PostgreSQL

MySQL is used as the primary relational database.

The key architectural principle is not the specific database engine, but service-owned data and clear database boundaries.

---

## 22. Open Questions

| Question                                                                             | Status    |
| ------------------------------------------------------------------------------------ | --------- |
| Will the API Gateway expose REST, GraphQL, or gRPC-Web to frontend clients?          | To decide |
| Will payments be fully mocked or integrated with a sandbox provider later?           | To decide |
| Will search use MySQL full-text search initially or a dedicated search engine later? | To decide |
| Will recommendations be rules-based initially or event-driven with behavioural data? | To decide |
| Will Kubernetes deployment use Helm, Kustomize, or both?                             | To decide |
| Will the cloud target be AWS first, or multi-cloud later?                            | To decide |
| Will service mesh be included in the first Kubernetes version?                       | Deferred  |


---

## 23. Related Documents

This product vision should be read alongside:

```text
docs/requirements/scope.md
docs/requirements/user-journeys.md
docs/requirements/functional-requirements.md
docs/requirements/non-functional-requirements.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/events/event-catalog.md
docs/data/data-ownership.md
```

---

## 24. Summary

bfstore is a realistic cloud-native backend for ACME Ltd’s online furniture store.

It is designed to support core commerce workflows while demonstrating senior-level backend, platform, Kubernetes, and DevSecOps engineering practices.

The product should be built incrementally, starting with a complete checkout vertical slice, then expanding into search, reviews, recommendations, Kubernetes deployment, GitOps, supply-chain security, observability, resilience, and production readiness.













