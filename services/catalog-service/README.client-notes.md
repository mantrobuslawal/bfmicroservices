# Catalogue Service Client Notes

## What this skeleton demonstrates

This service skeleton shows how bfstore will move from architecture documents into implementation.

It establishes:

```text
a real Go service boundary
a service-owned MySQL dependency
clean configuration loading
database connection management
repository and service layers
gRPC server wiring
container build path
health check extension point
```

## Why Catalogue Service first?

Catalogue Service is the best first implementation candidate because it already has:

```text
clear product requirements
Protobuf API contracts
MySQL schema
seed data
bounded service ownership
low-risk read APIs
```

This allows the first implementation slice to be useful without immediately dealing with checkout orchestration, payment safety, or distributed transaction concerns.

## First implementation target

```text
ListProducts
GetProduct
ListCategories
```

This creates a complete vertical slice:

```text
MySQL seed data
repository query
domain service
gRPC handler
Protobuf response
containerised runtime
```

## Next engineering steps

```text
1. Run buf generate from the repo root.
2. Commit generated Protobuf code.
3. Wire the generated CatalogServiceServer into internal/grpc/server.go.
4. Implement ListProducts.
5. Implement GetProduct.
6. Implement ListCategories.
7. Add integration tests against local MySQL.
8. Add the service to docker-compose.yml.
```
