# Non-Functional Requirements

## 1. Purpose

This document defines the non-functional requirements for **bfstore**, ACME Ltd’s fictional online furniture store backend.

Non-functional requirements describe the qualities the system must have, including reliability, performance, scalability, security, observability, maintainability, operability, data integrity, and deployment readiness.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s engineering quality.

---

## 2. Scope

This document covers non-functional requirements for:

```text
performance
availability
reliability
resilience
scalability
security
privacy
observability
maintainability
testability
deployability
data integrity
compatibility
operability
cost awareness
developer experience
```

These requirements apply across the bfstore service estate unless stated otherwise.

---

## 3. Requirement Priority

Priorities use the following scale:

| Priority | Meaning |
|---|---|
| Must | Required for the initial credible implementation |
| Should | Important for a mature professional version |
| Could | Valuable improvement after core maturity |
| Won't Yet | Explicitly deferred |

---

## 4. Performance Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-PERF-001 | The system shall define latency targets for core user journeys. | Must | Targets may be refined after baseline testing |
| NFR-PERF-002 | Catalogue browse shall be optimised for low-latency read access. | Should | High traffic path |
| NFR-PERF-003 | Checkout shall avoid unnecessary synchronous fan-out. | Must | Protect critical path |
| NFR-PERF-004 | Services shall use timeouts for remote calls. | Must | Prevent hanging requests |
| NFR-PERF-005 | Database queries for critical paths shall be indexed according to access patterns. | Must | Orders, baskets, stock, payments |
| NFR-PERF-006 | Kafka consumers shall process events fast enough to keep lag within defined thresholds. | Should | Thresholds to be refined |
| NFR-PERF-007 | Performance tests shall cover browse, basket, checkout, and event processing flows. | Should | k6 or equivalent |
| NFR-PERF-008 | Services shall expose latency metrics for inbound and outbound calls. | Should | Observability requirement |

## 4.1 Initial Performance Targets

Initial targets are placeholders and should be refined through measurement.

| Flow | Initial Target |
|---|---|
| Product list | p95 under 500 ms in local/test environment |
| Product details | p95 under 300 ms in local/test environment |
| Add basket item | p95 under 500 ms in local/test environment |
| Checkout | p95 under 2 seconds excluding simulated provider delays |
| Kafka notification processing | event processed within 5 seconds under normal local/test load |

These are not production SLOs. They are starting points for portfolio validation.

---

## 5. Availability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-AVL-001 | Each service shall expose a liveness endpoint. | Must | Kubernetes compatibility |
| NFR-AVL-002 | Each service shall expose a readiness endpoint. | Must | Traffic safety |
| NFR-AVL-003 | Services shall fail readiness when critical dependencies are unavailable. | Must | Avoid routing to unhealthy pods |
| NFR-AVL-004 | Notification failure shall not make the order service unavailable. | Must | Async decoupling |
| NFR-AVL-005 | Search or recommendation failure shall not block checkout. | Must | Non-critical downstream services |
| NFR-AVL-006 | The system shall support rolling deployment for most services. | Should | Deployment maturity |
| NFR-AVL-007 | Critical services should later support multiple replicas. | Should | Kubernetes maturity |

---

## 6. Reliability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-REL-001 | Checkout shall not create duplicate confirmed orders for duplicate submissions. | Must | Idempotency |
| NFR-REL-002 | Payment authorisation shall be idempotent where retries are possible. | Must | Avoid duplicate charges |
| NFR-REL-003 | Stock reservation shall be idempotent. | Must | Avoid double reservation |
| NFR-REL-004 | Shipment creation shall be idempotent. | Must | Avoid duplicate shipments |
| NFR-REL-005 | Event consumers shall handle duplicate events safely. | Must | Kafka reliability |
| NFR-REL-006 | Poison events shall not block consumer progress indefinitely. | Should | DLQ strategy |
| NFR-REL-007 | Important domain events should use an outbox pattern where practical. | Should | Prevent lost events |
| NFR-REL-008 | Order, payment, shipment, and notification state changes shall be auditable. | Should | Support and recovery |

---

## 7. Resilience Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-RES-001 | Every outbound gRPC call shall have a timeout. | Must | Prevent cascading failure |
| NFR-RES-002 | Retries shall only be used where operations are safe or idempotent. | Must | Avoid duplicate side effects |
| NFR-RES-003 | Kafka consumers shall use bounded retries. | Should | Retry/DLQ strategy |
| NFR-RES-004 | Kafka consumers shall send unprocessable messages to a DLQ after retry exhaustion. | Should | Operational maturity |
| NFR-RES-005 | Services shall log dependency failures with correlation context. | Must | Troubleshooting |
| NFR-RES-006 | Checkout shall fail safely when stock reservation fails. | Must | Payment must not be attempted |
| NFR-RES-007 | Checkout shall fail safely when payment authorisation fails. | Must | Order must not be confirmed |
| NFR-RES-008 | Notification failure shall not roll back order creation. | Must | Async side effect |
| NFR-RES-009 | Resilience tests shall cover stock failure, payment failure, duplicate events, and notification failure. | Should | Testing evidence |

---

## 8. Scalability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-SCL-001 | Services shall be independently deployable in principle. | Must | Microservice boundary |
| NFR-SCL-002 | Services shall be independently scalable in Kubernetes. | Should | Platform maturity |
| NFR-SCL-003 | Kafka consumers shall support horizontal scaling through consumer groups. | Should | Event processing |
| NFR-SCL-004 | Read-heavy paths such as catalogue and search shall be scalable independently of checkout. | Should | Service separation |
| NFR-SCL-005 | Services shall define CPU and memory requests and limits for Kubernetes deployment. | Should | Cluster stability |
| NFR-SCL-006 | Scaling decisions shall be based on metrics such as latency, request rate, CPU, memory, or consumer lag. | Could | Mature autoscaling |

---

## 9. Security Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-SEC-001 | Services shall not log secrets, tokens, raw payment data, or passwords. | Must | Secure logging |
| NFR-SEC-002 | Payment Service shall not store raw payment card data. | Must | Payment safety |
| NFR-SEC-003 | Service database users shall follow least privilege. | Must | One service, one schema |
| NFR-SEC-004 | Client-facing APIs shall return safe error messages. | Must | No internals leaked |
| NFR-SEC-005 | Protected customer operations shall require authentication once auth is implemented. | Should | Customer account security |
| NFR-SEC-006 | Customers shall not be able to access other customers’ orders. | Should | Authorisation |
| NFR-SEC-007 | CI shall include secret scanning. | Should | DevSecOps evidence |
| NFR-SEC-008 | CI shall include dependency vulnerability scanning. | Should | Supply-chain security |
| NFR-SEC-009 | Container images shall be scanned before deployment promotion. | Should | Container security |
| NFR-SEC-010 | Kubernetes workloads should run as non-root where practical. | Should | Runtime security |
| NFR-SEC-011 | Network policies should restrict unnecessary service-to-service traffic. | Could | Kubernetes security maturity |

---

## 10. Privacy and Data Protection Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-PRI-001 | Services shall minimise copying of personal data between service boundaries. | Must | PII minimisation |
| NFR-PRI-002 | Events shall avoid unnecessary customer PII. | Must | Event safety |
| NFR-PRI-003 | Order and shipment address snapshots shall be justified by historical fulfilment needs. | Should | Data ownership |
| NFR-PRI-004 | Logs shall avoid full customer addresses, emails, phone numbers, and payment details. | Must | Secure observability |
| NFR-PRI-005 | Data retention rules shall be defined for baskets, notifications, search logs, and recommendation signals. | Should | Lifecycle governance |
| NFR-PRI-006 | Customer deletion or anonymisation requirements shall be documented before production use. | Could | Mature privacy support |

---

## 11. Observability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-OBS-001 | All services shall emit structured logs. | Must | JSON logs preferred |
| NFR-OBS-002 | Logs shall include service name, operation, severity, and correlation ID. | Must | Troubleshooting |
| NFR-OBS-003 | Services shall expose request count, latency, and error metrics. | Should | Prometheus compatible |
| NFR-OBS-004 | gRPC client calls shall record dependency latency and failure metrics. | Should | Dependency visibility |
| NFR-OBS-005 | Kafka producers and consumers shall expose publish, consume, failure, retry, and lag metrics. | Should | Event operations |
| NFR-OBS-006 | Distributed tracing shall be supported through OpenTelemetry. | Should | End-to-end checkout trace |
| NFR-OBS-007 | The checkout journey shall be traceable across synchronous and asynchronous boundaries. | Should | Senior-level evidence |
| NFR-OBS-008 | Dashboards shall be created for service health, checkout, Kafka, and database behaviour. | Could | Mature operations |
| NFR-OBS-009 | Alerts shall be defined for critical failure indicators. | Could | Production readiness |

---

## 12. Maintainability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-MNT-001 | Services shall have clear ownership and boundaries. | Must | Architecture quality |
| NFR-MNT-002 | Shared packages shall not contain business ownership logic. | Must | Avoid shared domain coupling |
| NFR-MNT-003 | Protobuf contracts shall follow the project style guide. | Must | API consistency |
| NFR-MNT-004 | Documentation shall be kept aligned with implementation. | Must | Portfolio credibility |
| NFR-MNT-005 | ADRs shall capture significant architectural decisions. | Should | Decision history |
| NFR-MNT-006 | Services shall have README files explaining local run, config, tests, and operations. | Should | Developer experience |
| NFR-MNT-007 | Code should follow consistent formatting and linting standards. | Must | Review quality |

---

## 13. Testability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-TST-001 | Core business logic shall have unit tests. | Must | Fast feedback |
| NFR-TST-002 | gRPC contracts shall have contract tests. | Must | API confidence |
| NFR-TST-003 | Kafka producers and consumers shall have event contract tests. | Must | Event compatibility |
| NFR-TST-004 | MySQL repositories shall have integration tests where practical. | Should | Persistence confidence |
| NFR-TST-005 | The checkout vertical slice shall have end-to-end tests. | Must | Business proof |
| NFR-TST-006 | CI shall run unit, contract, and integration tests. | Must | Quality gate |
| NFR-TST-007 | Performance tests shall cover selected critical flows. | Should | Evidence |
| NFR-TST-008 | Resilience tests shall cover selected failure scenarios. | Should | Failure proof |
| NFR-TST-009 | Security checks shall run in CI. | Should | DevSecOps evidence |

---

## 14. Deployability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-DEP-001 | Each service shall have a container image. | Must | Containerisation |
| NFR-DEP-002 | Local development shall be supported through Docker Compose. | Must | Developer workflow |
| NFR-DEP-003 | Services shall be configurable through environment variables or mounted config. | Must | Twelve-factor style |
| NFR-DEP-004 | Secrets shall not be committed to Git. | Must | Security |
| NFR-DEP-005 | Kubernetes manifests, Helm charts, or Kustomize bases shall be provided. | Should | Platform evidence |
| NFR-DEP-006 | Deployments shall include liveness and readiness probes. | Should | Kubernetes readiness |
| NFR-DEP-007 | Deployments shall define resource requests and limits. | Should | Cluster stability |
| NFR-DEP-008 | Database migrations shall be deployment-aware. | Must | Safe releases |
| NFR-DEP-009 | GitOps integration shall be supported through environment overlays or deployable artefacts. | Could | Mature platform workflow |

---

## 15. Data Integrity Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-DATA-001 | Each service shall own its own schema. | Must | Service-owned data |
| NFR-DATA-002 | Services shall not query other services’ schemas. | Must | Boundary enforcement |
| NFR-DATA-003 | Cross-service references shall use IDs, not database foreign keys. | Must | Independent schemas |
| NFR-DATA-004 | Monetary values shall be stored in minor units. | Must | Avoid floating point money |
| NFR-DATA-005 | Order items shall preserve historical product and price snapshots. | Must | Historical accuracy |
| NFR-DATA-006 | Stock levels shall not become negative. | Must | Inventory invariant |
| NFR-DATA-007 | Migrations shall be versioned and tested. | Must | Database quality |
| NFR-DATA-008 | Important lifecycle entities should record status history. | Should | Auditability |

---

## 16. Compatibility and Versioning Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-COMP-001 | Protobuf packages shall be versioned. | Must | Example: `acme.order.v1` |
| NFR-COMP-002 | Kafka topics shall include a major version. | Must | Example: `.v1` |
| NFR-COMP-003 | Breaking protobuf changes shall fail CI where possible. | Should | Buf breaking checks |
| NFR-COMP-004 | Event consumers shall handle compatible event additions safely. | Must | Forward evolution |
| NFR-COMP-005 | API Gateway external routes shall be versioned if exposing REST/JSON. | Should | External compatibility |
| NFR-COMP-006 | Deprecated fields or APIs shall have a documented migration path. | Should | Governance |

---

## 17. Operability Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-OPS-001 | Services shall provide health and readiness endpoints. | Must | Operability |
| NFR-OPS-002 | Runbooks shall exist for critical operational scenarios. | Should | Incident response |
| NFR-OPS-003 | The system shall expose enough telemetry to diagnose checkout failures. | Must | Client-facing maturity |
| NFR-OPS-004 | DLQ growth shall be visible and actionable. | Should | Event operations |
| NFR-OPS-005 | Outbox backlog shall be visible if outbox is implemented. | Should | Event reliability |
| NFR-OPS-006 | Deployment rollback procedures shall be documented. | Should | Release safety |
| NFR-OPS-007 | Database migration failure handling shall be documented. | Should | Operational safety |

---

## 18. Cost and Resource Awareness

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-COST-001 | Local development should run on reasonable workstation resources. | Must | Avoid excessive dependencies |
| NFR-COST-002 | Optional services should be startable through Compose profiles. | Should | Developer efficiency |
| NFR-COST-003 | Kubernetes resources should define requests and limits. | Should | Cost control |
| NFR-COST-004 | Observability retention should be controlled. | Could | Avoid excessive storage |
| NFR-COST-005 | Cloud deployment should include basic cost controls and tagging where implemented. | Could | FinOps evidence |

---

## 19. Developer Experience Requirements

| ID | Requirement | Priority | Notes |
|---|---|---:|---|
| NFR-DEV-001 | The repo shall provide a clear root README. | Must | Onboarding |
| NFR-DEV-002 | The repo shall provide Makefile or script commands for common tasks. | Must | Ease of use |
| NFR-DEV-003 | Developers shall be able to start the local environment with one command. | Should | `make dev-up` |
| NFR-DEV-004 | Developers shall be able to regenerate protobuf code consistently. | Must | Contract-first workflow |
| NFR-DEV-005 | Developers shall be able to run all tests locally. | Should | Confidence |
| NFR-DEV-006 | `.env.example` shall document required local variables. | Must | Configuration clarity |

---

## 20. Initial Non-Functional Requirements for Vertical Slice

The first implementation should prioritise:

```text
NFR-REL-001
NFR-REL-002
NFR-REL-003
NFR-REL-005
NFR-RES-001
NFR-RES-002
NFR-SEC-001
NFR-SEC-002
NFR-OBS-001
NFR-OBS-002
NFR-TST-001
NFR-TST-002
NFR-TST-003
NFR-TST-005
NFR-DEP-001
NFR-DEP-002
NFR-DATA-001
NFR-DATA-002
NFR-DATA-004
NFR-DATA-006
NFR-COMP-001
NFR-COMP-002
NFR-DEV-001
```

---

## 21. Measurement and Evidence

Where possible, non-functional requirements should produce evidence.

Examples:

```text
test reports
coverage output
performance test summaries
container scan results
SBOM artefacts
CI pipeline logs
OpenTelemetry traces
Grafana dashboards
Kubernetes deployment screenshots
runbooks
ADR links
```

A strong client-facing portfolio should show evidence, not only claims.

---

## 22. Open Questions

| Question | Status |
|---|---|
| What production-like SLOs should be defined later? | To decide |
| Which performance tool will be standard? | Proposed: k6 |
| Which tracing backend will be used locally? | To decide |
| Which vulnerability scanners will be included first? | Proposed: Gitleaks, govulncheck, Trivy |
| Should the first Kubernetes deployment use Helm or Kustomize? | To decide |
| What level of auth is required for the initial slice? | To decide |
| What observability stack will be used locally? | To decide |

---

## 23. Related Documents

This document should be read alongside:

```text
docs/requirements/functional-requirements.md
docs/requirements/business-rules.md
docs/requirements/acceptance-criteria.md
docs/architecture/resilience-patterns.md
docs/architecture/deployment-view.md
docs/api/error-model.md
docs/events/retry-and-dlq-strategy.md
docs/data/mysql-standards.md
docs/testing/testing-strategy.md
docs/security/supply-chain-security.md
docs/observability/logging.md
docs/operations/production-readiness.md
```

---

## 24. Summary

bfstore’s non-functional requirements define the engineering quality expected of the system.

The most important early qualities are:

```text
safe checkout behaviour
idempotency
service-owned data
secure handling of sensitive data
observability
testability
containerisation
versioned contracts
local developer experience
```

These requirements help bfstore demonstrate senior-level engineering judgement beyond basic feature implementation.
