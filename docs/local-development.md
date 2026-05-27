# Local Development Foundation

## 1. Purpose

This document explains the initial local development foundation for **bfstore**.

The goal is to give contributors, reviewers, and potential clients a clear, repeatable way to run local dependencies and work with the contract-first development workflow.

---

## 2. Files Added

```text
Makefile
.env.example
docker-compose.yml
docs/local-development.md
```

These files establish the first local developer workflow.

---

## 3. Developer Workflow

Recommended setup:

```sh
cp .env.example .env
make help
make proto-lint
make proto-generate
make up
```

Stop the local stack:

```sh
make down
```

Remove local volumes:

```sh
make down-volumes
```

---

## 4. Makefile

The `Makefile` provides a consistent command interface.

Key commands:

```text
make help              show available commands
make proto-lint        lint Protobuf contracts with Buf
make proto-breaking    run Buf breaking-change checks
make proto-generate    generate Go code from Protobuf contracts
make proto             lint and generate Protobuf contracts
make up                start local dependencies
make down              stop local dependencies
make logs              tail container logs
make test              run Go tests
make check             run local quality checks
```

This keeps local development predictable and reviewer-friendly.

---

## 5. Environment Configuration

`.env.example` documents the expected local environment variables.

Developers should copy it to `.env`:

```sh
cp .env.example .env
```

The `.env` file should not be committed.

The example includes:

```text
project settings
MySQL settings
service-owned database names
Kafka bootstrap settings
Kafka topic names
gRPC service ports
OpenTelemetry settings
```

---

## 6. Docker Compose

The initial `docker-compose.yml` starts:

```text
MySQL 8.4
Kafka in KRaft mode
Kafka topic initialisation
```

Application services are intentionally left commented as future implementation examples.

This keeps the first local stack focused on dependencies while the service skeletons are built.

---

## 7. MySQL

bfstore uses service-owned MySQL databases.

Each service should own its own database/schema and database user.

Examples:

```text
bfstore_catalog
bfstore_basket
bfstore_inventory
bfstore_order
bfstore_payment
bfstore_shipping
bfstore_notification
```

No service should directly access another service's database.

---

## 8. Kafka

Kafka is used for domain events.

bfstore uses Protocol Buffers for Kafka event payloads.

Initial topics:

```text
bfstore.order.orders.v1
bfstore.inventory.stock.v1
bfstore.payment.payments.v1
bfstore.shipping.shipments.v1
bfstore.notification.notifications.v1
bfstore.catalog.products.v1
```

Kafka auto topic creation is disabled to make topic ownership explicit.

---

## 9. Protobuf and Buf

bfstore uses Protobuf for:

```text
gRPC service APIs
Kafka event payloads
shared common types
```

Buf is used for:

```text
linting
breaking-change checks
code generation
```

Recommended commands:

```sh
make proto-lint
make proto-breaking
make proto-generate
```

---

## 10. Client-Facing Rationale

This foundation demonstrates:

```text
repeatable local development
contract-first engineering
explicit infrastructure dependencies
service-owned data thinking
event-driven architecture preparation
professional developer experience
```

This is deliberately more than “it runs on my machine”.

It provides evidence of how the project will be built, operated, and reviewed.

---

## 11. Next Steps

Recommended next implementation steps:

```text
1. Add MySQL initialisation scripts for service-owned databases.
2. Add catalog-service Go skeleton.
3. Add generated Protobuf code.
4. Implement CatalogService ListProducts and GetProduct.
5. Add catalogue database migrations.
6. Add integration tests against MySQL.
7. Add first Kafka producer for ProductUpdated later.
```
