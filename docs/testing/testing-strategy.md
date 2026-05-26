# Testing Strategy

## 1. Purpose

This document defines the testing strategy for **bfstore**, ACME Ltd’s fictional online furniture store backend.

The purpose of this document is to describe how the system will be tested across:

- service logic
- gRPC APIs
- Kafka events
- MySQL persistence
- service integrations
- end-to-end business journeys
- performance behaviour
- resilience and failure handling
- security and supply-chain checks
- Kubernetes deployment readiness

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s engineering quality.

---

## 2. Testing Goals

The testing strategy should prove that bfstore is:

| Goal | Description |
|---|---|
| Correct | Services implement required business behaviour |
| Reliable | Important flows work consistently under expected conditions |
| Safe | Failures are handled without corrupting business state |
| Observable | Failures can be diagnosed through logs, metrics, and traces |
| Secure | Common security and supply-chain issues are checked |
| Maintainable | Changes can be made with confidence |
| Deployable | Services can be safely built, released, and operated |
| Client-ready | The project demonstrates professional engineering discipline |

---

## 3. Scope

This strategy covers tests for:

```text
Go service code
gRPC APIs
protobuf contracts
Kafka event contracts
MySQL repositories and migrations
Docker Compose local environment
Kubernetes deployment artefacts
CI/CD quality gates
security and supply-chain checks
performance and resilience behaviour
```

It does not replace detailed test plans for each service. Service-specific requirements should be documented in:

```text
docs/requirements/service-requirements/
```

---

## 4. Testing Principles

## 4.1 Test Business Behaviour, Not Just Code

Tests should prove meaningful business behaviour.

Good:

```text
Given a basket contains available products
When checkout is submitted
Then stock is reserved, payment is authorised, an order is created, and OrderCreated is published
```

Weak:

```text
test function returns true
```

---

## 4.2 Push Fast Tests Earlier

Fast tests should run frequently.

Recommended order:

```text
unit tests
static checks
contract tests
integration tests
end-to-end tests
performance tests
resilience tests
```

The faster the test, the earlier it should run in CI.

---

## 4.3 Test Service Boundaries

Because bfstore is a microservice system, testing must validate service boundaries.

Tests should confirm:

- services do not depend on another service’s database
- gRPC APIs match protobuf contracts
- Kafka events match documented schemas
- consumers handle duplicate events
- each service owns its own data and behaviour

---

## 4.4 Test Failure Paths

A professional system must test failure, not only happy paths.

Important failures:

```text
insufficient stock
payment declined
payment service unavailable
inventory service unavailable
Kafka unavailable
duplicate checkout request
duplicate OrderCreated event
MySQL unavailable
slow downstream service
```

---

## 4.5 Tests Should Produce Useful Evidence

The testing strategy should support a senior contractor portfolio by producing evidence such as:

- CI test reports
- coverage summaries
- contract test output
- performance test results
- resilience test results
- security scan results
- deployment validation results

---

## 5. Test Levels

## 5.1 Unit Tests

### Purpose

Unit tests validate small pieces of service logic without external dependencies.

### Scope

Unit tests should cover:

```text
domain logic
validation rules
state transitions
error mapping
idempotency helpers
price calculations
status transitions
retry decision logic
```

### Examples

| Service | Unit Test Examples |
|---|---|
| `basket-service` | quantity validation, basket item updates, remove item behaviour |
| `inventory-service` | stock reservation rules, reservation expiry rules |
| `order-service` | order state transitions, checkout validation, duplicate request handling |
| `payment-service` | payment state transitions, payment failure mapping |
| `shipping-service` | shipment status transitions |
| `notification-service` | duplicate notification suppression |

### Expected Characteristics

Unit tests should be:

- fast
- deterministic
- isolated
- easy to run locally
- part of every pull request

### Example Command

```sh
go test ./...
```

---

## 5.2 Repository and Persistence Tests

### Purpose

Repository tests validate database access logic against MySQL.

### Scope

Repository tests should cover:

```text
SQL queries
transactions
constraints
migrations
data mapping
duplicate key behaviour
pagination
locking behaviour where relevant
```

### Recommended Approach

Use test containers or Docker Compose-backed MySQL for integration-style repository tests.

### Important Services

Repository tests are especially important for:

```text
inventory-service
order-service
payment-service
basket-service
shipping-service
```

### Example Test Cases

```text
stock reservation cannot reduce available stock below zero
order insert creates order and order items in one transaction
payment attempt is recorded with failure reason
basket item quantity update is persisted correctly
```

---

## 5.3 gRPC Contract Tests

### Purpose

gRPC contract tests validate that service APIs conform to protobuf contracts and expected behaviour.

### Scope

Contract tests should cover:

```text
request validation
response shape
error model
status codes
backward compatibility
required business fields
timeout behaviour
```

### Services

All gRPC services should have contract tests.

Initial focus:

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
```

### Example Contract Tests

```text
GetProduct returns active product details
AddBasketItem rejects inactive product
ReserveStock returns reservation ID when stock is available
AuthorisePayment returns authorised result for valid payment request
CreateShipment returns shipment ID and tracking reference
CreateOrder returns existing order for duplicate idempotency key
```

### Tooling

Potential tooling:

```text
Buf breaking-change checks
Go gRPC test clients
generated protobuf clients
custom contract test suites
```

---

## 5.4 Kafka Event Contract Tests

### Purpose

Kafka event contract tests validate that producers and consumers agree on event payloads.

### Scope

Event contract tests should cover:

```text
event envelope
required metadata
payload schema
event version
producer correctness
consumer compatibility
unknown field handling
duplicate event handling
```

### Critical Events

Initial critical events:

```text
StockReserved
StockReservationFailed
PaymentAuthorised
PaymentFailed
ShipmentCreated
ShipmentFailed
OrderCreated
OrderFailed
NotificationSent
NotificationFailed
```

### Example Tests

```text
Order Service publishes valid OrderCreated event
Notification Service consumes OrderCreated event
Duplicate OrderCreated does not send duplicate notification
Payment Service publishes PaymentFailed with safe failure reason
Search Service handles ProductUpdated events idempotently
```

---

## 5.5 Integration Tests

### Purpose

Integration tests validate a service working with real dependencies such as MySQL, Kafka, or dependent service test doubles.

### Scope

Integration tests should cover:

```text
service startup
database connectivity
migrations
Kafka producer and consumer behaviour
gRPC client/server behaviour
transaction boundaries
configuration loading
health and readiness checks
```

### Example Integration Tests

| Service | Integration Test |
|---|---|
| `inventory-service` | reserve stock using MySQL transaction |
| `order-service` | create order with repository and mocked downstream services |
| `payment-service` | record payment attempt and publish event |
| `notification-service` | consume `OrderCreated` and record notification |
| `search-service` | consume `ProductUpdated` and update projection |

### Local Dependencies

Integration tests may use:

```text
Docker Compose
Testcontainers
MySQL
Kafka
service mocks
```

---

## 5.6 End-to-End Tests

### Purpose

End-to-end tests validate complete user journeys across multiple services.

### Initial Critical Journey

```text
Browse product
    -> Add to basket
    -> Checkout
    -> Reserve stock
    -> Authorise payment
    -> Create order
    -> Create shipment
    -> Publish OrderCreated
    -> Send notification
```

### Priority E2E Tests

| Test | Priority |
|---|---|
| Browse catalogue | Must |
| View product details | Must |
| Add product to basket | Must |
| Update basket item quantity | Must |
| Successful checkout | Must |
| Insufficient stock checkout failure | Must |
| Payment failure checkout failure | Must |
| OrderCreated triggers notification | Must |
| Duplicate checkout request does not create duplicate order | Must |
| View order history | Should |
| Shipment tracking | Should |
| Submit product review | Could |
| Search products | Should |
| Recommendations | Could |

### E2E Test Rules

- E2E tests should be few but valuable.
- They should validate real business flows.
- They should run against a realistic local environment.
- They should produce logs and traces that can be inspected.
- They should not replace lower-level tests.

---

## 5.7 Performance Tests

### Purpose

Performance tests validate latency, throughput, resource usage, and stability under load.

### Candidate Journeys

| Journey | Reason |
|---|---|
| Browse catalogue | High read traffic |
| View product details | Read latency and optional service fan-out |
| Add to basket | Write latency |
| Checkout | Multi-service latency and reliability |
| OrderCreated notification flow | Kafka throughput and consumer lag |
| Search products | Query latency and projection performance |

### Metrics

Performance tests should capture:

```text
request rate
response time
p50 latency
p95 latency
p99 latency
error rate
throughput
CPU usage
memory usage
database connection usage
Kafka consumer lag
```

### Initial Targets

Initial performance targets should be defined in:

```text
docs/requirements/non-functional-requirements.md
docs/testing/performance-testing.md
```

Example placeholder targets:

| Flow | Example Target |
|---|---|
| Catalogue browse | p95 under agreed threshold |
| Product details | p95 under agreed threshold |
| Add to basket | p95 under agreed threshold |
| Checkout | p95 under agreed threshold excluding external provider latency |
| Notification processing | consumer lag remains within agreed threshold |

### Tooling Options

Potential tools:

```text
k6
Vegeta
Locust
Go benchmark tests
Prometheus
Grafana
```

---

## 5.8 Resilience Tests

### Purpose

Resilience tests validate system behaviour under failure.

### Critical Scenarios

| Scenario | Expected Behaviour |
|---|---|
| Catalogue Service unavailable | Product browsing fails safely with clear error |
| Inventory Service unavailable during checkout | Checkout fails safely or enters defined state |
| Insufficient stock | Payment is not attempted |
| Payment Service unavailable | Checkout fails safely or enters defined pending state |
| Payment declined | Order is not confirmed and reservation is released or expires |
| Shipping Service unavailable | Order flow follows documented failure or pending fulfilment behaviour |
| Kafka unavailable | Event publishing is retried or outbox stores pending event |
| Notification Service unavailable | Order creation still succeeds |
| Duplicate `OrderCreated` event | Duplicate notification is avoided where possible |
| MySQL unavailable | Affected service readiness fails and requests fail safely |
| Slow downstream service | Timeout prevents indefinite request hanging |

### Resilience Patterns to Validate

```text
timeouts
retries
idempotency
dead-letter queues
outbox publishing
duplicate event handling
graceful degradation
readiness failure
compensation actions
```

### Tooling Options

```text
Docker Compose fault injection
Toxiproxy
k6
Chaos Mesh
LitmusChaos
custom failure scripts
```

---

## 5.9 Security Tests

### Purpose

Security tests reduce the risk of vulnerable code, dependencies, containers, and configuration.

### Scope

Security testing should include:

```text
dependency scanning
container image scanning
secret scanning
SAST
IaC scanning
Kubernetes manifest scanning
protobuf/API input validation tests
authentication tests
authorisation tests
sensitive logging tests
```

### Example Tools

```text
Trivy
Grype
Syft
Gitleaks
Semgrep
govulncheck
Checkov
tfsec
Conftest
Kyverno
OpenSSF Scorecard
```

### Critical Security Test Cases

```text
tokens are not logged
raw payment details are not stored
inactive users cannot access protected actions
customers cannot view another customer's orders
service database users cannot access other schemas
containers do not run as root where avoidable
Kubernetes manifests define resource limits
```

---

## 5.10 Deployment and Smoke Tests

### Purpose

Deployment tests confirm that services can be deployed and started correctly.

### Scope

Smoke tests should validate:

```text
service starts
health endpoint responds
readiness endpoint responds
database connection works
Kafka connection works
gRPC server responds
basic request succeeds
basic event can be published or consumed
```

### Kubernetes Readiness Checks

Kubernetes deployment validation should include:

```text
liveness probes
readiness probes
resource requests
resource limits
service account usage
config injection
secret references
network policy compatibility
startup behaviour
```

---

## 6. Test Pyramid

Recommended test distribution:

```text
many unit tests
many contract tests
some integration tests
few but valuable end-to-end tests
targeted performance tests
targeted resilience tests
```

Conceptual pyramid:

```text
        resilience / chaos tests
          performance tests
        end-to-end tests
      integration tests
   contract tests
unit tests and static checks
```

The project should avoid relying on a large number of slow end-to-end tests for confidence.

---

## 7. Test Environments

## 7.1 Local Development

Used for:

```text
unit tests
service integration tests
manual testing
local E2E checkout flow
```

Expected tooling:

```text
Docker Compose
MySQL
Kafka
service containers
Makefile commands
```

Example commands:

```sh
make dev-up
make migrate-up
make test
make smoke-test
make dev-down
```

---

## 7.2 CI Environment

Used for:

```text
static checks
unit tests
contract tests
integration tests
security checks
container builds
protobuf checks
```

CI should run on every pull request.

---

## 7.3 Kubernetes Test Environment

Used for:

```text
deployment validation
smoke tests
observability validation
resilience tests
performance tests
GitOps validation
```

This may be added after the local vertical slice is working.

---

## 8. CI Quality Gates

A professional CI pipeline should include:

| Stage | Checks |
|---|---|
| Format | Go formatting, Markdown linting where practical |
| Static checks | Go vet, staticcheck, linting |
| Unit tests | `go test ./...` |
| Protobuf checks | Buf lint and breaking-change checks |
| Contract tests | gRPC and event contract tests |
| Integration tests | MySQL and Kafka-backed tests |
| Security checks | secret scan, dependency scan, container scan |
| Build | container image build |
| SBOM | generate software bill of materials |
| Sign | sign image where implemented |
| Smoke | start service and validate health/readiness |

Example quality gate:

```text
A pull request cannot merge if unit tests, protobuf linting, contract tests, or critical security checks fail.
```

---

## 9. Test Data Strategy

## 9.1 Test Data Principles

Test data should be:

- realistic enough to reflect the business domain
- deterministic
- safe to store in Git
- free from real personal data
- resettable between tests
- documented

## 9.2 Example Test Data

Products:

```text
sofa
dining table
wardrobe
office chair
coffee table
bed frame
```

Checkout test customer:

```text
customer_id: test-customer-001
address: synthetic UK delivery address
```

Payment simulation:

```text
payment token: test-payment-success
payment token: test-payment-declined
payment token: test-payment-timeout
```

## 9.3 Data Isolation

Tests should avoid depending on shared mutable data.

Where possible:

- create data per test
- use unique IDs
- reset database state between tests
- isolate Kafka topics or consumer groups in tests

---

## 10. Traceability

Tests should trace back to requirements.

Example:

```text
FR-003: The system shall reserve stock before confirming an order.
    -> inventory-service unit tests
    -> ReserveStock gRPC contract test
    -> order-service integration test
    -> successful checkout E2E test
    -> insufficient stock E2E test
```

Traceability links may appear in:

```text
test names
test descriptions
pull request descriptions
requirements documents
CI reports
```

---

## 11. Service-Specific Testing Focus

## 11.1 Catalog Service

Focus:

```text
product visibility
active/inactive product behaviour
category filtering
product detail retrieval
product event publishing
```

## 11.2 Inventory Service

Focus:

```text
stock reservation
reservation expiry
stock release
stock commit
concurrent reservation attempts
no overselling
```

## 11.3 Basket Service

Focus:

```text
add item
update quantity
remove item
empty basket validation
inactive product rejection
basket checkout transition
```

## 11.4 Order Service

Focus:

```text
checkout orchestration
order creation
idempotency
order state transitions
stock failure handling
payment failure handling
shipment failure handling
event publishing
```

## 11.5 Payment Service

Focus:

```text
authorisation success
authorisation failure
payment attempt audit
idempotent payment requests
sensitive data exclusion
```

## 11.6 Shipping Service

Focus:

```text
delivery options
shipment creation
tracking status
shipment failure
idempotent shipment requests
```

## 11.7 Notification Service

Focus:

```text
event consumption
notification request creation
duplicate event handling
send simulation
retry behaviour
failure recording
```

## 11.8 Search Service

Focus:

```text
product projection updates
search query behaviour
inactive product exclusion
index rebuild
event replay
```

## 11.9 Recommendation Service

Focus:

```text
rules-based recommendations
inactive product exclusion
signal processing
fallback behaviour
```

---

## 12. Observability in Tests

Tests should validate observability where practical.

## 12.1 Logs

Validate that:

- correlation IDs are present
- sensitive data is not logged
- important business failures are logged
- duplicate event handling is visible

## 12.2 Metrics

Validate or inspect:

```text
request count
request latency
error count
checkout success count
checkout failure count
Kafka consumer lag
notification failure count
```

## 12.3 Traces

Important E2E tests should generate traces across:

```text
api-gateway
order-service
basket-service
inventory-service
payment-service
shipping-service
notification-service
```

This demonstrates operational maturity as well as functional correctness.

---

## 13. Definition of Done for a Service

A service is considered test-ready when it has:

```text
unit tests for core logic
repository/integration tests where it uses MySQL
gRPC contract tests for exposed APIs
event contract tests for produced or consumed events
health and readiness tests
basic security checks
test data fixtures
CI integration
documentation of important failure cases
```

A service is considered production-readiness-test-ready when it also has:

```text
performance test coverage for critical endpoints
resilience test coverage for key dependency failures
dashboard and alert validation
deployment smoke tests
runbook-linked failure scenarios
```

---

## 14. Initial Implementation Testing Plan

The first implementation should focus on the checkout vertical slice.

## 14.1 Initial Test Scope

Services:

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

## 14.2 Initial Must-Have Tests

```text
catalog-service lists active products
basket-service adds active product to basket
basket-service rejects inactive product
inventory-service reserves available stock
inventory-service rejects insufficient stock
payment-service authorises simulated successful payment
payment-service rejects simulated declined payment
order-service creates order successfully
order-service does not confirm order when stock reservation fails
order-service does not confirm order when payment fails
shipping-service creates shipment
notification-service consumes OrderCreated
duplicate OrderCreated does not create duplicate notification where idempotency is implemented
successful checkout E2E passes
insufficient stock checkout E2E passes
payment failure checkout E2E passes
```

## 14.3 Initial CI Gates

```text
go test ./...
buf lint
buf breaking
unit tests
contract tests
integration tests for core services
secret scan
dependency scan
container scan
```

---

## 15. Later Testing Enhancements

Later phases should add:

```text
full auth and authorisation tests
customer profile tests
review service tests
search projection replay tests
recommendation signal tests
performance test suite
soak tests
chaos tests
Kubernetes deployment tests
GitOps sync validation
policy-as-code tests
SBOM and image signing verification
production readiness scorecard checks
```

---

## 16. Test Reporting

Test outputs should be easy to review.

Recommended artefacts:

```text
unit test report
coverage report
contract test report
integration test report
E2E test report
performance test summary
resilience test summary
security scan report
SBOM artefact
container scan report
```

These artefacts strengthen the portfolio because they show evidence, not just claims.

---

## 17. Risks and Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| Too many E2E tests | Slow and fragile CI | Keep E2E tests focused on core journeys |
| Tests require too much local setup | Poor developer experience | Use Makefile commands and Docker Compose profiles |
| Contract tests are skipped | Service compatibility issues | Add contract tests to CI gates |
| Event consumers are not idempotent | Duplicate side effects | Test duplicate event handling |
| Performance testing is deferred too long | Unknown scalability | Add small k6 tests after vertical slice |
| Security scanning is superficial | Weak DevSecOps evidence | Add scans to CI and store reports |
| Docs drift from tests | Misleading portfolio | Link tests to requirements and update docs with changes |

---

## 18. Open Questions

| Question | Status |
|---|---|
| Which test framework will be standard for Go services? | To decide |
| Will Testcontainers or Docker Compose be preferred for integration tests? | To decide |
| Will k6 be used for performance tests? | Proposed |
| Will Toxiproxy be used for local resilience tests? | Proposed |
| Should contract tests run before integration tests in CI? | Proposed |
| Should coverage thresholds be enforced from the start? | To decide |
| Should end-to-end tests run on every pull request or nightly? | To decide |
| Should Kubernetes smoke tests run on every merge to main? | To decide |

---

## 19. Related Documents

This document should be read alongside:

```text
docs/requirements/user-journeys.md
docs/requirements/functional-requirements.md
docs/requirements/non-functional-requirements.md
docs/architecture/service-boundaries.md
docs/architecture/resilience-patterns.md
docs/api/grpc-overview.md
docs/events/event-catalog.md
docs/data/data-ownership.md
docs/security/supply-chain-security.md
docs/observability/logging.md
docs/observability/tracing.md
docs/operations/production-readiness.md
```

Relevant ADRs:

```text
adr/0002-use-grpc-for-service-communication.md
adr/0003-use-kafka-for-events.md
adr/0004-use-service-owned-databases.md
adr/0006-use-buf-for-protobuf.md
adr/0007-use-opentelemetry.md
adr/0008-use-contract-first-service-design.md
```

---

## 20. Summary

bfstore’s testing strategy is designed to prove more than basic code correctness.

It validates:

```text
business behaviour
service boundaries
API contracts
event contracts
data ownership
checkout reliability
failure handling
security posture
deployment readiness
operational visibility
```

The first testing priority is the checkout vertical slice:

```text
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

This gives the project a strong professional testing foundation and demonstrates the level of engineering discipline expected from a senior platform, DevSecOps, or Kubernetes-focused contractor.
