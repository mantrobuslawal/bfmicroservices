# Catalogue Service

## 1. Purpose

The **Catalogue Service** owns product catalogue truth for bfstore.

It is responsible for:

```text
products
categories
product variants
category-scoped product attributes
product images
catalogue lifecycle state
```

It is not responsible for:

```text
stock levels
basket state
orders
payments
shipping
search ranking
recommendations
```

---

## 2. Service Role

Catalogue Service provides customer-facing product and category data over gRPC.

Initial API focus:

```text
ListProducts
GetProduct
ListCategories
```

Later API expansion:

```text
ListProductAttributeDefinitions
admin product management
catalogue event publishing
search projection updates
```

---

## 3. Architecture

The service follows a small layered structure:

```text
cmd/catalog-service/
└── main.go

internal/
├── config/
│   └── config.go
├── database/
│   └── mysql.go
├── health/
│   └── health.go
├── catalog/
│   ├── repository.go
│   └── service.go
└── grpc/
    ├── server.go
    └── catalog_handler.go
```

Layer responsibilities:

```text
cmd              application entry point
config           environment configuration
database         MySQL connection setup
catalog          domain service and repository logic
grpc             transport handlers
health           service readiness/liveness checks
```

---

## 4. Data Ownership

Catalogue Service owns the `bfstore_catalog` database.

It must not directly access databases owned by other services.

Other services should use Catalogue Service APIs or consume Catalogue events where appropriate.

---

## 5. Contracts

Catalogue Service contracts live under:

```text
proto/acme/catalog/v1/
```

Expected generated Go package:

```text
gen/go/acme/catalog/v1
```

This skeleton assumes generated Protobuf code will be added by running:

```sh
buf generate
```

---

## 6. Local Development

From the repository root:

```sh
make up
make proto-generate
```

Then, from this service directory or root-level service commands later:

```sh
go test ./...
```

The service expects MySQL configuration through environment variables.

Example:

```text
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=bfstore_catalog
MYSQL_USER=bfstore_catalog_user
MYSQL_PASSWORD=bfstore_catalog_password
CATALOG_SERVICE_GRPC_PORT=50051
```

---

## 7. Implementation Status

Current status:

```text
service skeleton
configuration loader
MySQL connection helper
repository interfaces
service layer placeholders
gRPC server wiring
health checks
Dockerfile
```

Next implementation target:

```text
ListProducts
GetProduct
ListCategories
```

---

## 8. Client-Facing Engineering Evidence

This service demonstrates:

```text
contract-first service design
clean service boundaries
service-owned database access
layered Go service structure
container-ready implementation
clear local development path
future-ready observability and health hooks
```
