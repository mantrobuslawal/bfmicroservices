# Deployment View

## 1. Purpose

This document defines the deployment view for **bfstore**, ACME Ltd’s fictional online furniture store backend.

It explains how the application is expected to run locally, how it should be packaged, how it may be deployed to Kubernetes, and how deployment responsibilities are split across the wider ACME platform repository estate.

This document is intended for engineers, reviewers, technical leads, and potential clients evaluating bfstore’s deployment and platform engineering maturity.

---

## 2. Deployment Context

bfstore is the application backend repo.

It contains:

```text
service source code
protobuf contracts
Kafka event definitions
database migrations
application tests
Dockerfiles
local Docker Compose setup
Kubernetes-ready deployment artefacts
application-level documentation
```

The wider ACME platform estate is split across multiple repositories:

```text
acme-ltd/
├── bfstore
├── acme-platform-infra
├── acme-platform-gitops
├── acme-terraform-modules
├── acme-security-governance
└── acme-developer-platform
```

## 2.1 Repository Responsibilities

| Repository | Responsibility |
|---|---|
| `bfstore` | Application services, contracts, migrations, tests, app deployment artefacts |
| `acme-platform-infra` | Cloud infrastructure, VPCs, Kubernetes, Kafka, MySQL, observability platform |
| `acme-platform-gitops` | Desired Kubernetes state for environments, applications, and platform add-ons |
| `acme-terraform-modules` | Reusable infrastructure modules |
| `acme-security-governance` | Security standards, policy-as-code, threat models, supply-chain controls |
| `acme-developer-platform` | Backstage, golden paths, templates, scorecards |

This document focuses on the application deployment view for the `bfstore` repo.

---

## 3. Deployment Goals

bfstore deployments should support:

| Goal | Description |
|---|---|
| Local development | Engineers can run the application locally |
| Repeatable builds | Services are built consistently from source |
| Containerisation | Each service is packaged as a container image |
| Kubernetes readiness | Services can run on Kubernetes with health checks and resource controls |
| Environment separation | Dev, test, staging, and production can have separate configuration |
| Observability | Deployments expose logs, metrics, traces, and health status |
| Security | Services run with least privilege and safe defaults |
| GitOps compatibility | Desired runtime state can be managed by GitOps |
| Rollback | Failed releases can be safely rolled back |
| Production readiness | Deployments are testable, supportable, and auditable |

---

## 4. Deployment Units

Each backend service should be independently buildable and deployable.

Target services:

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

Initial deployment focus:

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

Each service should eventually provide:

```text
Dockerfile
configuration example
health endpoint
readiness endpoint
metrics endpoint
database migrations where required
Kubernetes deployment artefact
service README
runbook
```

---

## 5. Local Development Deployment

Local development should be supported through Docker Compose.

## 5.1 Local Components

Initial local environment:

```text
MySQL
Kafka
api-gateway
catalog-service
inventory-service
basket-service
order-service
payment-service
shipping-service
notification-service
```

Optional local observability:

```text
Prometheus
Grafana
Loki
Tempo
OpenTelemetry Collector
```

## 5.2 Local Development Commands

Expected commands:

```sh
make dev-up
make proto
make migrate-up
make test
make smoke-test
make dev-down
```

## 5.3 Local Database Model

For local development, a single MySQL container may host multiple logical schemas:

```text
bfstore_catalog
bfstore_inventory
bfstore_basket
bfstore_order
bfstore_payment
bfstore_shipping
bfstore_notification
```

This keeps local development manageable while preserving service-owned database boundaries.

## 5.4 Local Kafka Model

Local Kafka should support the initial checkout event flow:

```text
bfstore.inventory.stock-events.v1
bfstore.payment.payment-events.v1
bfstore.shipping.shipment-events.v1
bfstore.order.order-events.v1
bfstore.notification.notification-events.v1
```

---

## 6. Container Image Strategy

Each service should produce its own container image.

Example image names:

```text
bfstore/api-gateway
bfstore/catalog-service
bfstore/inventory-service
bfstore/basket-service
bfstore/order-service
bfstore/payment-service
bfstore/shipping-service
bfstore/notification-service
```

## 6.1 Image Tagging

Recommended tags:

```text
git SHA
semantic version
environment promotion tag where appropriate
```

Examples:

```text
catalog-service:sha-abc1234
catalog-service:v0.1.0
catalog-service:dev
```

## 6.2 Image Build Requirements

Images should:

- use minimal base images where practical
- run as non-root where possible
- avoid embedding secrets
- include only required runtime artefacts
- expose documented ports
- support health and readiness checks
- include build metadata labels where useful

## 6.3 Supply Chain Controls

Container builds should eventually include:

```text
dependency scanning
container scanning
SBOM generation
image signing
provenance metadata
policy checks
```

---

## 7. Configuration Strategy

Services should be configured through environment variables or mounted configuration, not hard-coded values.

Configuration examples:

```text
SERVICE_NAME
ENVIRONMENT
GRPC_PORT
HTTP_PORT
DATABASE_DSN
KAFKA_BROKERS
LOG_LEVEL
OTEL_EXPORTER_ENDPOINT
JWT_ISSUER
```

## 7.1 Configuration Rules

- No secrets in Git.
- Use `.env.example` for local development examples.
- Use Kubernetes ConfigMaps for non-sensitive configuration.
- Use Kubernetes Secrets or External Secrets for sensitive values.
- Keep configuration names consistent across services.
- Validate configuration at startup.

---

## 8. Secrets Strategy

Secrets may include:

```text
database passwords
JWT signing keys
provider credentials
Kafka credentials
TLS certificates
```

## 8.1 Local Development

Local development may use safe dummy secrets stored in `.env.example`.

Real secrets must not be committed.

## 8.2 Kubernetes

Kubernetes deployments should reference secrets through:

```text
Kubernetes Secrets
External Secrets Operator
Vault or cloud secret manager
```

The broader secrets governance model belongs in:

```text
acme-security-governance
```

---

## 9. Kubernetes Deployment View

In Kubernetes, each service should run as an independent workload.

Conceptual deployment:

```text
Kubernetes Cluster
    Namespace: bfstore-dev
        api-gateway Deployment
        catalog-service Deployment
        inventory-service Deployment
        basket-service Deployment
        order-service Deployment
        payment-service Deployment
        shipping-service Deployment
        notification-service Deployment
        Services
        ConfigMaps
        Secret references
        NetworkPolicies
```

## 9.1 Kubernetes Objects

Each service may require:

```text
Deployment
Service
ServiceAccount
ConfigMap
Secret reference
HorizontalPodAutoscaler
PodDisruptionBudget
NetworkPolicy
ServiceMonitor or PodMonitor
```

## 9.2 Namespaces

Recommended namespace model:

```text
bfstore-dev
bfstore-test
bfstore-staging
bfstore-prod
```

Alternative:

```text
acme-dev
acme-test
acme-staging
acme-prod
```

Final environment naming should align with the wider platform standards.

---

## 10. Kubernetes Health Checks

Each service should expose:

```text
liveness check
readiness check
startup check where needed
```

## 10.1 Liveness

Liveness indicates whether the process should be restarted.

It should not fail just because a temporary downstream dependency is unavailable.

## 10.2 Readiness

Readiness indicates whether the service can receive traffic.

Readiness may fail when:

```text
database unavailable
required configuration invalid
migration incomplete
critical dependency unavailable
```

## 10.3 Startup

Startup probes may be useful for services with slower bootstrapping or migration checks.

---

## 11. Resource Management

Each service should define resource requests and limits.

Example categories:

```text
CPU request
CPU limit
memory request
memory limit
```

## 11.1 Rules

- Avoid deploying services with no resource requests.
- Avoid unrealistic memory limits that cause frequent OOM kills.
- Tune resources through testing.
- Use separate profiles for local development and Kubernetes deployment.

---

## 12. Scaling Strategy

Services should be independently scalable.

Initial scaling model:

| Service | Scaling Consideration |
|---|---|
| `api-gateway` | Scale by request volume |
| `catalog-service` | Scale by product read traffic |
| `basket-service` | Scale by basket write traffic |
| `order-service` | Scale carefully due to checkout orchestration |
| `inventory-service` | Scale carefully due to stock consistency |
| `payment-service` | Scale based on payment requests and provider limits |
| `shipping-service` | Scale based on fulfilment requests |
| `notification-service` | Scale by Kafka consumer lag |
| `search-service` | Scale by query load |
| `recommendation-service` | Scale by recommendation traffic |

## 12.1 Horizontal Autoscaling

Horizontal Pod Autoscaling may use:

```text
CPU
memory
request rate
Kafka consumer lag
custom business metrics
```

Use custom metrics only after basic service behaviour is stable.

---

## 13. Network and Access Model

Kubernetes network access should be restricted.

## 13.1 Network Policy Goals

- API Gateway can call public-facing backend services.
- Order Service can call Basket, Inventory, Payment, and Shipping.
- Notification Service can consume Kafka and call provider integrations.
- Services can access only their own database.
- Unnecessary east-west traffic should be blocked.

## 13.2 Example Desired Access

```text
api-gateway -> catalog-service
api-gateway -> basket-service
api-gateway -> order-service
order-service -> basket-service
order-service -> inventory-service
order-service -> payment-service
order-service -> shipping-service
notification-service -> Kafka
```

## 13.3 Forbidden Access

```text
order-service -> catalog database
basket-service -> inventory database
notification-service -> order database
search-service -> catalog database
```

---

## 14. Environment Strategy

bfstore should support multiple environments.

| Environment | Purpose |
|---|---|
| `dev` | Developer integration and early testing |
| `test` | Automated integration and E2E testing |
| `staging` | Production-like validation |
| `prod` | Production target state |

## 14.1 Environment Differences

Configuration may vary by:

```text
replica count
resource limits
database endpoints
Kafka brokers
log level
feature flags
observability endpoints
secrets provider
external provider mode
```

## 14.2 Promotion Model

Recommended promotion flow:

```text
dev
    -> test
        -> staging
            -> prod
```

Container images should be promoted between environments rather than rebuilt differently for each environment.

---

## 15. GitOps Relationship

The bfstore repo may contain reusable deployment artefacts.

The desired live state should live in:

```text
acme-platform-gitops
```

## 15.1 bfstore Repo Owns

```text
Dockerfiles
Helm charts or Kustomize bases
example Kubernetes manifests
service deployment requirements
application configuration examples
```

## 15.2 GitOps Repo Owns

```text
environment-specific overlays
image versions deployed to each environment
Argo CD applications
cluster-specific configuration
platform add-on desired state
policy bindings
```

This separation keeps application code and live environment state cleanly separated.

---

## 16. Deployment Strategies

## 16.1 Rolling Deployment

Default strategy for most services.

Suitable for:

```text
catalog-service
basket-service
notification-service
search-service
recommendation-service
```

## 16.2 Blue/Green Deployment

Useful when release risk is higher or fast rollback is required.

Potentially useful for:

```text
api-gateway
order-service
payment-service
```

## 16.3 Canary Deployment

Useful when traffic can be gradually shifted.

Potentially useful for:

```text
api-gateway
catalog-service
order-service
```

## 16.4 Initial Recommendation

Start with rolling deployments, then introduce blue/green or canary for critical services after observability and traffic control are mature.

---

## 17. Database Migration Deployment

Each service owns its own migrations.

Migration options:

```text
run migrations before service rollout
run migrations as Kubernetes Jobs
run migrations in CI/CD before deployment
run migrations through service startup only if carefully controlled
```

## 17.1 Migration Rules

- Migrations must be versioned.
- Destructive migrations require explicit review.
- Backward-compatible migrations should be preferred.
- Application rollout and migration order must be documented.
- Rollback must consider database compatibility.

## 17.2 Expand and Contract Pattern

For safer changes:

```text
1. Expand schema with backward-compatible change.
2. Deploy application that writes both old and new fields where needed.
3. Backfill data.
4. Switch reads to new field.
5. Remove old field in later release.
```

---

## 18. Rollback Strategy

Rollback should be safe and documented.

## 18.1 Application Rollback

Application rollback may involve:

```text
revert image tag
sync GitOps state
verify readiness
run smoke tests
monitor error rates
```

## 18.2 Database Rollback

Database rollback is more sensitive.

Rules:

- avoid destructive migrations during the same release as application changes
- prefer backward-compatible schema changes
- document rollback limitations
- test rollback where practical

## 18.3 Event Contract Rollback

Event contract rollback should consider:

```text
consumer compatibility
producer version
topic version
schema version
replay behaviour
```

---

## 19. Observability Deployment

Services should emit telemetry to the platform observability stack.

Expected signals:

```text
structured logs
metrics
distributed traces
health status
readiness status
Kafka consumer lag
business metrics
```

Possible stack:

```text
OpenTelemetry Collector
Prometheus
Grafana
Loki
Tempo
Alertmanager
```

## 19.1 Minimum Service Metrics

```text
request count
request latency
error count
dependency latency
database errors
Kafka publish failures
Kafka consumer lag where applicable
```

---

## 20. Security Deployment Requirements

Kubernetes workloads should follow secure defaults.

Recommended controls:

```text
run as non-root where possible
read-only root filesystem where practical
drop unnecessary Linux capabilities
set resource requests and limits
use service accounts per service
avoid default service account
use network policies
avoid secrets in environment where stronger options exist
scan images before deployment
sign images where implemented
```

Policy enforcement may be managed by:

```text
Kyverno
OPA Gatekeeper
Conftest
```

The broader policy model belongs in:

```text
acme-security-governance
```

---

## 21. Smoke Tests

After deployment, smoke tests should validate:

```text
API Gateway reachable
Catalog Service returns active products
Basket Service can add item
Order Service can start checkout
Inventory Service can reserve test stock
Payment Service can simulate authorisation
Shipping Service can create test shipment
Notification Service can consume test event
health and readiness endpoints pass
```

Smoke tests should run after deployments to dev/test/staging and after production releases where appropriate.

---

## 22. Initial Deployment Scope

Initial deployment should focus on local and Kubernetes-ready application artefacts for:

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

Initial infrastructure dependencies:

```text
MySQL
Kafka
OpenTelemetry-compatible logging/tracing
```

Initial deployment artefacts:

```text
Dockerfiles
docker-compose.yml
Kubernetes base manifests or Helm charts
environment overlays
health/readiness endpoints
Makefile commands
```

---

## 23. Target Deployment Maturity

A mature service deployment should include:

```text
container image
SBOM
image scan report
image signature
Kubernetes Deployment
Kubernetes Service
ServiceAccount
ConfigMap
Secret reference
NetworkPolicy
resource requests and limits
readiness and liveness probes
metrics scraping configuration
dashboard
alerts
runbook
rollback instructions
```

---

## 24. Deployment Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Local environment becomes too heavy | Poor developer experience | Use Compose profiles and phased startup |
| Database migrations break rollback | Production instability | Use expand/contract migration strategy |
| Services deploy without resource limits | Cluster instability | Enforce resource policies |
| Missing readiness checks | Traffic sent to unhealthy pods | Add readiness probes for all services |
| Secrets committed to Git | Security incident | Use secret scanning and External Secrets |
| Event contract changes break consumers | Runtime failures | Use versioning and compatibility checks |
| Too much platform work before app works | Delayed progress | Prove local checkout vertical slice first |

---

## 25. Open Questions

| Question | Status |
|---|---|
| Will application deployment use Helm, Kustomize, or both? | To decide |
| Will local development use Docker Compose profiles per service group? | Proposed |
| Will migrations run as Kubernetes Jobs or as CI/CD steps? | To decide |
| Should production-style deployments use canary for critical services? | Later |
| Which cloud provider will host the first Kubernetes environment? | To decide |
| Should service mesh be introduced for mTLS and traffic management? | Deferred |
| Should image signing be required before GitOps promotion? | Proposed for mature stage |

---

## 26. Related Documents

This document should be read alongside:

```text
docs/architecture/service-boundaries.md
docs/architecture/communication-patterns.md
docs/architecture/resilience-patterns.md
docs/data/data-ownership.md
docs/data/migrations.md
docs/testing/testing-strategy.md
docs/security/supply-chain-security.md
docs/observability/logging.md
docs/observability/metrics.md
docs/observability/tracing.md
docs/operations/deployment-strategy.md
docs/operations/rollback-strategy.md
docs/operations/production-readiness.md
```

Relevant repositories:

```text
acme-platform-infra
acme-platform-gitops
acme-terraform-modules
acme-security-governance
acme-developer-platform
```

---

## 27. Summary

bfstore’s deployment model is designed to support a professional application lifecycle:

```text
local development
container build
test automation
Kubernetes deployment
GitOps promotion
observability
security controls
rollback
production readiness
```

The first deployment goal is to run the checkout vertical slice locally and produce Kubernetes-ready artefacts.

The mature deployment goal is to demonstrate that bfstore can be built, tested, scanned, deployed, observed, secured, rolled back, and operated as part of a realistic platform engineering estate.
