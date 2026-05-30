# Documentation

This directory contains the documentation for **bfstore**, the cloud-native microservice backend for ACME Ltd’s fictional online furniture store.

The documentation is intended to explain:

- what the system must do
- why the system is designed this way
- how services communicate
- how data is owned and stored
- how the system is tested
- how the system is secured
- how the system is observed and operated

The aim is to keep design, implementation, testing, and operations aligned as the platform evolves.

---

## Documentation Structure

```text
docs/
├── README.md
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

## Directory Guide

| Directory        | Purpose                                                                                                  |
| ---------------- | -------------------------------------------------------------------------------------------------------- |
| `requirements/`  | Defines what bfstore must do and what quality standards it must meet                                     |
| `architecture/`  | Explains the system design, service boundaries, communication patterns, and trade-offs                   |
| `api/`           | Documents gRPC APIs, protobuf conventions, API gateway behaviour, errors, and versioning                 |
| `events/`        | Documents Kafka topics, event contracts, event versioning, ordering, retries, and replay                 |
| `data/`          | Defines service-owned data, MySQL design standards, migrations, consistency, PII, and retention          |
| `testing/`       | Defines unit, integration, contract, end-to-end, performance, and resilience testing                     |
| `security/`      | Covers application security, threat modelling, authentication, authorisation, secrets, and secure coding |
| `observability/` | Covers logging, metrics, tracing, dashboards, alerts, SLOs, and Kafka consumer lag                       |
| `operations/`    | Covers runbooks, deployments, rollbacks, incidents, DR, backup/restore, and production readiness         |


---

## Recommended Reading Order

Start with the documents that explain the business and system shape before moving into technical details.

1. requirements/product-vision.md
2. requirements/scope.md
3. requirements/user-journeys.md
4. requirements/functional-requirements.md
5. requirements/non-functional-requirements.md
6. architecture/domain-model.md
7. architecture/service-boundaries.md
8. architecture/communication-patterns.md
9. api/grpc-overview.md
10. events/event-catalog.md
11. data/data-ownership.md
12. testing/testing-strategy.md
13. security/threat-model.md
14. observability/logging.md
15. operations/production-readiness.md

This order reflects the intended design flow:
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

---

## Development Standards

- [Go Code Quality Tooling](./development/go-code-quality-tooling.md)





