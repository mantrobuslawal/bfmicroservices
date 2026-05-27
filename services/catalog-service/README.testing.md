# Catalogue Service Testing

## 1. Purpose

This document explains the initial testing approach for Catalogue Service.

The goal is to keep the first implementation slice small, professional, and easy to review.

---

## 2. Test Types

The initial test set includes:

```text
unit tests
repository helper tests
integration test scaffolding
```

Later test expansion should include:

```text
gRPC handler tests
contract tests
container-based integration tests
performance smoke tests
```

---

## 3. Unit Tests

Run from the service directory:

```sh
go test ./...
```

Initial unit tests cover:

```text
catalogue service behaviour
required product ID validation
not found propagation
pagination limit normalisation
pagination offset normalisation
```

---

## 4. Integration Tests

Integration tests are opt-in.

They require:

```text
local MySQL running
catalogue schema migrated
seed data applied
```

Start local dependencies from the repository root:

```sh
make up
```

Then run:

```sh
BFSTORE_RUN_INTEGRATION_TESTS=true go test ./test/integration/...
```

Optional DSN override:

```sh
CATALOG_TEST_MYSQL_DSN='bfstore_catalog_user:bfstore_catalog_password@tcp(localhost:3306)/bfstore_catalog?parseTime=true&charset=utf8mb4,utf8'
```

---

## 5. Why Integration Tests Are Opt-In

Integration tests need external dependencies.

Making them opt-in keeps normal unit test runs fast and reliable.

CI can later run them in a dedicated job with MySQL available.

---

## 6. Next Testing Improvements

Recommended next improvements:

```text
add testcontainers-go for database integration tests
add generated gRPC handler tests
add contract tests against Protobuf fixtures
add migration up/down tests
add seed data validation tests
add repository tests for category filters
```

---

## 7. Client-Facing Value

This testing structure demonstrates:

```text
clear test separation
fast local feedback
dependency-aware integration testing
database-backed verification
a path towards CI maturity
```
