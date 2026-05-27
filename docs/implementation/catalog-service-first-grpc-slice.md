# Catalogue Service First gRPC Slice

## 1. Purpose

This document explains the next implementation slice for bfstore Catalogue Service.

The target is to move from a service skeleton to a working gRPC-backed read service.

---

## 2. Goal

Implement:

```text
ListProducts
GetProduct
ListCategories
```

The end-to-end flow should be:

```text
gRPC client
→ Catalogue Service handler
→ Catalogue Service business layer
→ MySQL repository
→ bfstore_catalog database
→ Protobuf response
```

---

## 3. Files Added

```text
services/catalog-service/internal/grpc/mappers.go
services/catalog-service/internal/grpc/catalog_handler.go
services/catalog-service/internal/grpc/server.go
services/catalog-service/internal/catalog/service_test.go
services/catalog-service/internal/catalog/repository_test.go
services/catalog-service/test/integration/catalog_repository_integration_test.go
services/catalog-service/README.testing.md
docker-compose.catalog-service.patch.yml
Makefile.catalog-service.patch
```

---

## 4. Why the Mapper Layer Matters

The mapper layer converts internal domain models into Protobuf response models.

This prevents generated transport code from leaking across the service.

Good boundary:

```text
repository model
→ service model
→ mapper
→ protobuf response
```

Avoid:

```text
database rows
→ protobuf response everywhere
```

---

## 5. Generated Protobuf Dependency

The concrete gRPC handler requires generated Protobuf code.

Expected flow:

```sh
buf generate
```

Then wire:

```go
catalogv1.RegisterCatalogServiceServer(server, handler)
```

The files currently include TODO sections and target implementation examples so that the design is clear before generated code is committed.

---

## 6. Testing Approach

Initial tests include:

```text
service unit tests
repository helper tests
integration test scaffolding
```

Run unit tests:

```sh
cd services/catalog-service
go test ./...
```

Run integration tests:

```sh
cd services/catalog-service
BFSTORE_RUN_INTEGRATION_TESTS=true go test ./test/integration/...
```

---

## 7. Docker Compose Integration

A patch file is included:

```text
docker-compose.catalog-service.patch.yml
```

Add the service block to the main `docker-compose.yml` under `services:` once the service is ready to run in the local stack.

---

## 8. Makefile Integration

A patch file is included:

```text
Makefile.catalog-service.patch
```

Add the targets to the root `Makefile` to support:

```text
catalog-test
catalog-integration-test
catalog-build
catalog-docker-build
```

---

## 9. Next Engineering Steps

Recommended next steps:

```text
1. Run buf generate.
2. Commit generated Protobuf code.
3. Replace placeholder gRPC TODOs with concrete generated method signatures.
4. Register CatalogServiceServer in server.go.
5. Implement product/category mappers.
6. Run unit tests.
7. Run integration tests.
8. Add catalog-service to docker-compose.yml.
```

---

## 10. Client-Facing Value

This slice demonstrates:

```text
contract-first implementation
clean Go service layering
transport/domain separation
database-backed service behaviour
professional testing direction
local runtime integration path
```

It is a strong next portfolio milestone because it turns the Catalogue Service from structure into behaviour.
