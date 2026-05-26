# Architecture Trade-Offs

## 1. Purpose

This document records the major architecture trade-offs for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains why key design choices were made, which alternatives were considered, what benefits were gained, and what costs or risks were accepted.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating the maturity of bfstore’s architecture.

---

## 2. Why Trade-Offs Matter

Good architecture is not about choosing perfect technologies. It is about making deliberate decisions based on goals, constraints, risks, and context.

This document makes those decisions explicit so that the project can show:

- clear engineering judgement
- awareness of alternatives
- honest discussion of complexity
- ability to explain why decisions were made
- ability to revisit decisions when context changes

---

## 3. Project Context

bfstore is designed as a serious portfolio project for demonstrating backend, platform, Kubernetes, and DevSecOps engineering capability.

The system is intentionally more sophisticated than a basic CRUD application.

Primary goals:

```text
demonstrate service boundaries
use contract-first APIs
use event-driven patterns
show service-owned data
support Kubernetes deployment
include observability and resilience thinking
show secure delivery practices
produce client-facing documentation
```

Important constraints:

```text
solo development effort
portfolio project scope
need to remain implementable
local development should stay manageable
costs should stay controlled
external integrations may be mocked initially
```

---

## 4. Summary of Major Trade-Offs

| Decision | Benefit | Cost |
|---|---|---|
| Microservices over monolith | Demonstrates service boundaries and platform skills | More complexity |
| gRPC for internal APIs | Strong contracts and efficient communication | More tooling and learning curve |
| Kafka for events | Decoupled asynchronous workflows | Eventual consistency and operational complexity |
| MySQL per service | Clear data ownership | More schemas and migrations |
| API Gateway | Hides internal topology | Risk of gateway becoming monolith |
| Order Service orchestrates checkout | Clear ownership of order workflow | Risk of growing orchestration complexity |
| Simulated external providers first | Keeps scope manageable | Less realistic integrations initially |
| Docker Compose for local dev | Easy local environment | May diverge from Kubernetes behaviour |
| Kubernetes-ready artefacts | Platform relevance | More deployment complexity |
| Documentation-first foundation | Clear design and client-facing quality | Risk of delaying implementation |

---

## 5. Microservices vs Monolith

## Decision

bfstore uses a microservice architecture.

## Alternatives Considered

| Option | Description |
|---|---|
| Modular monolith | Single deployable application with internal modules |
| Microservices | Multiple independently deployable services |
| Distributed modular monolith | Multiple services but tightly coupled data and releases |

## Rationale

Microservices are chosen because the project is intended to demonstrate:

```text
service boundaries
independent ownership
gRPC communication
Kafka events
service-owned databases
Kubernetes deployment
observability across services
resilience and failure handling
platform engineering maturity
```

## Benefits

- Clear business capability boundaries
- Independent service deployment potential
- Better fit for Kubernetes and platform engineering demonstration
- Supports event-driven architecture
- Allows realistic DevSecOps and operational patterns

## Costs

- More local development complexity
- More testing complexity
- More deployment complexity
- Need for observability and tracing
- Risk of distributed monolith if boundaries are poor

## Mitigation

- Start with a focused checkout vertical slice.
- Keep service boundaries explicit.
- Avoid shared databases.
- Use contract-first APIs.
- Keep documentation aligned with implementation.
- Avoid building all services fully at once.

## Review Trigger

Revisit if:

```text
service boundaries become artificial
local development becomes unmanageable
most changes require coordinated releases
implementation progress stalls due to distributed complexity
```

---

## 6. gRPC vs REST for Internal Communication

## Decision

bfstore uses gRPC for internal service-to-service communication.

## Alternatives Considered

| Option | Description |
|---|---|
| REST/JSON | Simple HTTP APIs |
| gRPC | Strongly typed RPC using protobuf |
| Async-only messaging | Services communicate mainly through Kafka |

## Rationale

gRPC is chosen because it provides:

```text
typed contracts
code generation
efficient binary protocol
clear service APIs
good fit for internal service communication
protobuf reuse for events
```

## Benefits

- Strong service contracts
- Better compatibility checking with Buf
- Reduced ambiguity in API payloads
- Good fit for Go services
- Useful senior-level engineering signal

## Costs

- Less browser-friendly than REST
- Requires protobuf tooling
- Requires client generation
- Error model must be designed carefully
- Debugging can be less straightforward than plain JSON

## Mitigation

- Use API Gateway for client-facing protocol adaptation.
- Document error model clearly.
- Use Buf for linting and breaking-change detection.
- Provide local tooling and examples.

## Review Trigger

Revisit if:

```text
external clients need direct browser access
API Gateway complexity increases significantly
team familiarity with gRPC becomes a blocker
```

---

## 7. Kafka vs Synchronous Calls for Downstream Workflows

## Decision

bfstore uses Kafka for asynchronous business events.

## Alternatives Considered

| Option | Description |
|---|---|
| Direct synchronous calls | Producer calls every downstream service |
| Kafka events | Producer publishes facts; consumers react |
| Job queue | Work items are processed asynchronously |
| Database polling | Consumers poll source tables |

## Rationale

Kafka is chosen to decouple services that react to business events.

Examples:

```text
OrderCreated -> notification-service
ProductUpdated -> search-service
ReviewApproved -> recommendation-service
```

## Benefits

- Decouples producers from consumers
- Supports multiple independent consumers
- Good fit for notifications, search, recommendations, and analytics
- Enables event replay and projections
- Demonstrates event-driven architecture

## Costs

- Eventual consistency
- Consumer idempotency required
- Kafka operational complexity
- Schema/versioning complexity
- Harder debugging without strong observability

## Mitigation

- Use clear event catalog.
- Use consistent event envelope.
- Make consumers idempotent.
- Add consumer lag and DLQ monitoring.
- Start with a small event set.
- Consider outbox pattern for reliable publishing.

## Review Trigger

Revisit if:

```text
events are used as hidden RPC
event flows become hard to reason about
Kafka operation outweighs benefit for early stages
```

---

## 8. Service-Owned Databases vs Shared Database

## Decision

Each service owns its own logical MySQL database or schema.

## Alternatives Considered

| Option | Description |
|---|---|
| Shared database | All services use one database |
| Schema per service | One MySQL instance, separate schemas |
| Database per service | Physically separate databases |
| Event-sourced storage | State primarily reconstructed from events |

## Rationale

Service-owned databases preserve microservice boundaries and prevent hidden coupling.

## Benefits

- Clear data ownership
- Services can evolve internal schemas independently
- Avoids cross-service joins
- Supports least privilege access
- More realistic microservice architecture

## Costs

- More migrations
- More schemas
- More integration complexity
- Reporting becomes harder
- Cross-service queries require APIs, events, or projections

## Mitigation

- Use one MySQL container with separate schemas locally.
- Use clear migration standards.
- Use APIs/events for cross-service data.
- Use projections for search and reporting.
- Document ownership explicitly.

## Review Trigger

Revisit if:

```text
services frequently need direct joins
data ownership is unclear
local database management becomes too heavy
```

---

## 9. MySQL vs PostgreSQL

## Decision

bfstore uses MySQL as the primary relational database.

## Alternatives Considered

| Option | Description |
|---|---|
| MySQL | Widely used relational database |
| PostgreSQL | Feature-rich relational database |
| NoSQL database | Document or key-value storage |
| Mixed persistence | Different database per service need |

## Rationale

MySQL is suitable for transactional commerce-style data such as:

```text
products
baskets
orders
payments
shipments
customers
reviews
```

The most important architectural choice is not the database engine but the service-owned data model.

## Benefits

- Familiar relational model
- Good fit for transactional data
- Widely understood
- Easy to run locally
- Suitable for portfolio implementation

## Costs

- Some PostgreSQL-specific features are unavailable
- Search may need separate tooling later
- Advanced JSON/query features may be less attractive than PostgreSQL in some cases

## Mitigation

- Keep data access behind service APIs.
- Avoid database-specific coupling where practical.
- Use migrations consistently.
- Use Search Service projection if search requirements outgrow MySQL.

## Review Trigger

Revisit if:

```text
advanced relational features become necessary
search requirements exceed MySQL capabilities
cloud provider constraints favour another engine
```

---

## 10. API Gateway vs Direct Service Exposure

## Decision

bfstore uses an API Gateway as the client-facing entry point.

## Alternatives Considered

| Option | Description |
|---|---|
| Direct service exposure | Clients call services directly |
| API Gateway | Single external entry point |
| Backend-for-Frontend | Separate gateway per client type |

## Rationale

The API Gateway hides internal service topology and provides a stable external API.

## Benefits

- Central entry point
- Authentication enforcement at the edge
- Client-friendly response shaping
- Internal service topology remains private
- Supports protocol translation

## Costs

- Gateway can become a bottleneck
- Gateway can become a monolith if business logic leaks in
- Additional component to deploy and observe

## Mitigation

- Keep business logic in domain services.
- Use gateway for routing, authentication, request shaping, and error mapping.
- Monitor gateway latency and error rates.
- Document gateway responsibilities clearly.

## Review Trigger

Revisit if:

```text
gateway accumulates domain logic
different clients require very different APIs
gateway becomes a performance bottleneck
```

---

## 11. Order Service as Checkout Orchestrator

## Decision

The initial architecture places checkout orchestration in `order-service`.

## Alternatives Considered

| Option | Description |
|---|---|
| Order Service orchestration | Order Service coordinates checkout |
| Dedicated workflow service | Separate orchestration component |
| Choreography via events | Services react to events without central orchestrator |
| Saga/workflow engine | Use Temporal, Cadence, or similar |

## Rationale

Order creation is the business process being completed. It is reasonable for Order Service to coordinate the first checkout implementation.

## Benefits

- Clear ownership of checkout outcome
- Easier to reason about initial flow
- Simpler than introducing a workflow engine
- Keeps first vertical slice achievable

## Costs

- Order Service may grow too much orchestration logic
- Risk of becoming a process manager for many domains
- Failure handling must be carefully designed
- Long-running workflows may become complex

## Mitigation

- Keep service ownership clear.
- Payment, stock, and shipping remain owned by their services.
- Use idempotency and explicit state transitions.
- Consider workflow tooling later if orchestration grows.

## Review Trigger

Revisit if:

```text
checkout becomes long-running and complex
many compensation paths are added
order-service becomes hard to maintain
workflow visibility becomes poor
```

---

## 12. Simulated External Providers First

## Decision

Payment, notification, and shipping providers may be simulated initially.

## Alternatives Considered

| Option | Description |
|---|---|
| Simulated providers | Local deterministic provider behaviour |
| Sandbox integrations | Use real provider test environments |
| Live integrations | Real payment/email/shipping services |

## Rationale

The first goal is to prove architecture, service boundaries, and checkout behaviour without being blocked by provider complexity.

## Benefits

- Keeps scope manageable
- Supports deterministic tests
- Avoids cost and compliance issues
- Allows failure scenarios to be simulated easily

## Costs

- Less realistic provider behaviour
- Less integration evidence initially
- Some operational concerns deferred

## Mitigation

- Design provider interfaces cleanly.
- Simulate success, failure, timeout, and retry scenarios.
- Add sandbox provider integration later if valuable.
- Document limitations clearly.

## Review Trigger

Revisit if:

```text
portfolio needs stronger real-world integration evidence
payment or notification behaviour becomes central to the project
```

---

## 13. Docker Compose vs Kubernetes for Initial Development

## Decision

Start with Docker Compose for local development and add Kubernetes-ready artefacts.

## Alternatives Considered

| Option | Description |
|---|---|
| Docker Compose first | Simple local environment |
| Kubernetes only | Develop directly on a cluster |
| Hybrid | Compose locally, Kubernetes for deployment validation |

## Rationale

Docker Compose keeps the first implementation accessible while Kubernetes artefacts support platform engineering goals.

## Benefits

- Easier local setup
- Faster development loop
- Lower cost
- Good for integration tests
- Avoids platform complexity before the app works

## Costs

- Compose differs from Kubernetes
- Some Kubernetes behaviours are not tested locally
- Risk of local-only assumptions

## Mitigation

- Add Kubernetes manifests or Helm charts.
- Run smoke tests in Kubernetes later.
- Keep configuration twelve-factor friendly.
- Avoid Compose-specific application logic.

## Review Trigger

Revisit if:

```text
Kubernetes behaviour becomes central
local Compose diverges from real deployment
platform validation becomes the main project focus
```

---

## 14. Documentation-First vs Code-First

## Decision

bfstore begins with strong requirements and architecture documentation before heavy implementation.

## Alternatives Considered

| Option | Description |
|---|---|
| Code-first | Build services quickly, document later |
| Documentation-first | Define architecture and requirements first |
| Iterative | Write enough docs, build, update docs |

## Rationale

The project is intended to demonstrate senior contractor judgement, not just coding ability.

## Benefits

- Clear design intent
- Better service boundaries
- Better client-facing evidence
- Easier to explain decisions
- Stronger foundation for tests and implementation

## Costs

- Risk of over-documentation
- Implementation may lag
- Some assumptions may change once coding starts

## Mitigation

- Treat docs as living documents.
- Build the checkout vertical slice soon.
- Update docs when implementation reveals better choices.
- Add personal design notes and ADRs.

## Review Trigger

Revisit if:

```text
documentation delays implementation too much
docs no longer match code
documents become generic instead of decision-focused
```

---

## 15. Search Service Now vs Later

## Decision

Search Service is part of the target architecture but may be deferred from the first implementation.

## Alternatives Considered

| Option | Description |
|---|---|
| Catalogue search first | Use Catalog Service for simple listing/filtering |
| Dedicated Search Service | Build search projection early |
| External search engine | Use OpenSearch/Elasticsearch |

## Rationale

The first priority is the checkout vertical slice. Search is valuable but not essential to prove checkout.

## Benefits of Deferring

- Reduces early service count
- Keeps implementation focused
- Avoids premature projection complexity
- Allows search requirements to mature

## Costs

- Less complete customer browsing experience initially
- Search architecture remains theoretical until implemented

## Mitigation

- Document Search Service boundary.
- Add ProductUpdated events early if useful.
- Implement basic catalogue filtering first.
- Add Search Service after core flow is working.

---

## 16. Recommendation Service Now vs Later

## Decision

Recommendation Service is part of the target architecture but deferred from the initial checkout slice.

## Rationale

Recommendations are useful for a commerce domain but not required to prove core platform behaviour.

## Benefits of Deferring

- Keeps first implementation realistic
- Avoids building weak recommendation logic too early
- Allows event signals to emerge naturally

## Costs

- Recommendation architecture remains theoretical initially
- Fewer advanced product capabilities in early demo

## Mitigation

- Start with simple rules later.
- Consume real events from basket and order flows.
- Avoid pretending to have machine learning before there is data.

---

## 17. Helm vs Kustomize

## Decision

To be decided.

## Options

| Option | Benefit | Cost |
|---|---|---|
| Helm | Templating, packaging, widely used | Can become complex |
| Kustomize | Overlay model, simple patching | Less powerful templating |
| Both | Common enterprise pattern | More moving parts |

## Initial Recommendation

Use simple Kubernetes manifests or Kustomize first, then introduce Helm if service packaging requires it.

For a platform portfolio, understanding both is useful, but using both from day one may be unnecessary.

## Review Trigger

Revisit when:

```text
deployment artefacts become repetitive
environment overlays grow
GitOps repo needs reusable app definitions
```

---

## 18. Service Mesh Now vs Later

## Decision

Service mesh is deferred.

## Alternatives Considered

| Option | Description |
|---|---|
| No mesh initially | Simpler Kubernetes deployment |
| Linkerd | Lightweight service mesh |
| Istio | Feature-rich service mesh |
| Cilium service mesh | Network/security integrated option |

## Rationale

Service mesh can add value for mTLS, traffic policy, retries, observability, and zero trust. However, it adds complexity and should not block the first working application.

## Benefits of Deferring

- Simpler deployment
- Easier debugging
- Lower learning overhead
- Faster vertical slice delivery

## Costs

- mTLS and advanced traffic policy deferred
- Some platform security maturity delayed

## Mitigation

- Design services to be mesh-compatible.
- Use clear service identities.
- Add mesh later as a platform maturity phase.

---

## 19. Outbox Pattern Now vs Later

## Decision

Outbox pattern is proposed for serious implementation but may be phased in.

## Rationale

The outbox pattern improves reliability between database writes and Kafka publishing.

## Benefits

- Prevents lost events after database commits
- Improves recovery
- Stronger event-driven reliability story
- Useful professional architecture signal

## Costs

- Additional tables
- Publisher process complexity
- Operational monitoring needed
- More tests required

## Initial Recommendation

Consider implementing outbox first in `order-service`, because `OrderCreated` is critical.

Then expand to:

```text
payment-service
inventory-service
shipping-service
catalog-service
review-service
```

---

## 20. Trade-Off Register

| ID | Decision | Status |
|---|---|---|
| `TD-001` | Use microservices for bfstore | Accepted |
| `TD-002` | Use gRPC for internal service communication | Accepted |
| `TD-003` | Use Kafka for asynchronous business events | Accepted |
| `TD-004` | Use service-owned MySQL schemas | Accepted |
| `TD-005` | Use API Gateway for client-facing access | Accepted |
| `TD-006` | Use Order Service as initial checkout orchestrator | Accepted for initial version |
| `TD-007` | Simulate external providers initially | Accepted |
| `TD-008` | Use Docker Compose for local development | Accepted |
| `TD-009` | Add Kubernetes-ready artefacts | Accepted |
| `TD-010` | Defer Search Service implementation | Proposed |
| `TD-011` | Defer Recommendation Service implementation | Proposed |
| `TD-012` | Defer service mesh | Proposed |
| `TD-013` | Consider outbox pattern for critical events | Proposed |
| `TD-014` | Choose Helm, Kustomize, or both | To decide |

---

## 21. Decision Review Process

Architecture trade-offs should be revisited when:

```text
implementation reveals unexpected complexity
service boundaries become unclear
testing becomes too fragile
deployment becomes too complex
performance goals are not met
operational risk increases
new requirements change the context
```

Major decisions should be recorded as ADRs in:

```text
adr/
```

---

## 22. Related Documents

This document should be read alongside:

```text
docs/requirements/product-vision.md
docs/requirements/scope.md
docs/architecture/domain-model.md
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/architecture/event-driven-design.md
docs/architecture/resilience-patterns.md
docs/architecture/deployment-view.md
docs/data/data-ownership.md
docs/events/event-catalog.md
docs/testing/testing-strategy.md
```

Relevant ADRs:

```text
adr/0001-use-microservices.md
adr/0002-use-grpc-for-service-communication.md
adr/0003-use-kafka-for-events.md
adr/0004-use-service-owned-databases.md
adr/0005-use-mysql.md
adr/0006-use-buf-for-protobuf.md
adr/0007-use-opentelemetry.md
adr/0008-use-contract-first-service-design.md
```

---

## 23. Summary

bfstore’s architecture deliberately favours professional cloud-native patterns:

```text
microservices
gRPC
Kafka
service-owned databases
API Gateway
Kubernetes-ready deployment
observability
resilience
security-aware delivery
```

These choices create more complexity than a simple monolith, but that complexity is intentional because the project is designed to demonstrate senior-level platform, DevSecOps, Kubernetes, and backend architecture capability.

The key to keeping the project credible is to deliver the design in phases, starting with the checkout vertical slice and then expanding into search, recommendations, GitOps, security controls, resilience testing, and production readiness.
