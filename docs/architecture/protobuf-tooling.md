# Protobuf Tooling

This document defines the standard Protocol Buffers tooling workflow for **bfstore**.

bfstore uses **Buf CLI** as the primary Protobuf toolchain for:

```text
building proto schemas
linting proto contracts
checking breaking changes
generating Go and gRPC code
managing proto dependencies
formatting .proto files
supporting CI contract checks
```

---

## Purpose

The purpose of this standard is to keep bfstore Protobuf contracts:

```text
consistent
versioned
linted
safe to evolve
easy to generate from
suitable for gRPC and Kafka payloads
```

Buf provides a professional workflow around Protocol Buffers so the project does not rely on long, fragile `protoc` commands.

---

## Repository Layout

Recommended layout:

```text
bfstore/
├── buf.yaml
├── buf.gen.yaml
├── buf.lock
├── proto/
│   └── bfstore/
│       ├── common/
│       │   └── v1/
│       ├── catalog/
│       │   └── v1/
│       ├── basket/
│       │   └── v1/
│       ├── inventory/
│       │   └── v1/
│       ├── order/
│       │   ├── v1/
│       │   └── events/
│       │       └── v1/
│       ├── payment/
│       │   └── v1/
│       ├── shipping/
│       │   └── v1/
│       └── notification/
│           └── v1/
├── gen/
│   └── go/
└── services/
```

`proto/` contains source contracts.

`gen/go/` contains generated Go code.

Services import generated code rather than hand-writing Protobuf or gRPC types.

---

## `buf.yaml`

`buf.yaml` defines the Protobuf module, lint rules, breaking-change policy, and dependencies.

Recommended starting config:

```yaml
version: v2

modules:
  - path: proto
    name: buf.build/mantrobuslawal/bfstore

lint:
  use:
    - STANDARD

breaking:
  use:
    - FILE
```

Rules:

```text
Buf should operate over proto/, not the whole repo.
Use versioned packages such as bfstore.catalog.v1.
Use STANDARD lint rules unless there is a strong reason not to.
Use FILE breaking-change rules initially.
```

---

## `buf.gen.yaml`

`buf.gen.yaml` defines code generation.

Recommended Go/gRPC config:

```yaml
version: v2

plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen/go
    opt:
      - paths=source_relative

  - remote: buf.build/grpc/go
    out: gen/go
    opt:
      - paths=source_relative
```

This generates:

```text
Protobuf message types
gRPC server interfaces
gRPC clients
```

Example generated package import:

```go
import catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
```

---

## Required Commands

bfstore should support these commands:

```bash
buf build
buf lint
buf format -w
buf generate
buf breaking --against '.git#branch=main'
```

Recommended meanings:

| Command | Purpose |
|---|---|
| `buf build` | Compile proto schema and check imports/types. |
| `buf lint` | Enforce schema style and quality rules. |
| `buf format -w` | Format `.proto` files consistently. |
| `buf generate` | Generate Go/gRPC code. |
| `buf breaking` | Detect incompatible contract changes. |

---

## Makefile Targets

Recommended Makefile targets:

```makefile
.PHONY: proto-build
proto-build:
	buf build

.PHONY: proto-lint
proto-lint:
	buf lint

.PHONY: proto-format
proto-format:
	buf format -w

.PHONY: proto-generate
proto-generate:
	buf generate

.PHONY: proto-breaking
proto-breaking:
	buf breaking --against '.git#branch=main'

.PHONY: proto-check
proto-check: proto-build proto-lint proto-breaking

.PHONY: proto
proto: proto-check proto-generate
```

Developer workflow:

```bash
make proto
```

---

## CI Requirements

Pull requests that change Protobuf contracts must run:

```text
buf build
buf lint
buf breaking --against main
buf generate
git diff --exit-code
```

Example GitHub Actions workflow:

```yaml
name: protobuf

on:
  pull_request:
    paths:
      - "proto/**"
      - "buf.yaml"
      - "buf.gen.yaml"
      - ".github/workflows/protobuf.yml"

jobs:
  protobuf:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: bufbuild/buf-action@v1
        with:
          setup_only: true

      - name: Build proto
        run: buf build

      - name: Lint proto
        run: buf lint

      - name: Check breaking changes
        run: buf breaking --against '.git#branch=main'

      - name: Generate proto code
        run: buf generate

      - name: Check generated code is committed
        run: git diff --exit-code
```

This ensures generated code stays in sync with contract changes.

---

## gRPC Contract Workflow

For gRPC APIs:

```text
edit .proto service contract
run buf format
run buf lint
run buf breaking
run buf generate
implement generated server interface
use generated clients
run Go tests
```

Example service:

```proto
service CatalogService {
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
}
```

Generated Go code provides:

```text
CatalogServiceServer interface
RegisterCatalogServiceServer function
CatalogServiceClient client
```

---

## Kafka Event Workflow

bfstore uses Protobuf for Kafka event payloads.

For Kafka events:

```text
define event message in proto
generate Go event type
marshal event using proto.Marshal
publish bytes to Kafka
consumer unmarshals using proto.Unmarshal
```

Example:

```go
event := &ordereventsv1.OrderCreated{
    OrderId: orderID,
}

payload, err := proto.Marshal(event)
```

Breaking-change checks are especially important for Kafka because old consumers may continue reading messages after producers change.

---

## Breaking Change Rules

Do not make incompatible changes to published messages or services.

Avoid:

```text
renaming fields without care
reusing field numbers
changing field types
removing fields without reserving names/numbers
changing package names casually
changing service/method names casually
```

When removing fields:

```proto
message Product {
  reserved 4;
  reserved "old_field_name";
}
```

Practical rule:

```text
Field numbers are part of the wire contract.
Treat them like public API.
```

---

## Dependencies

When external proto dependencies are needed, add them to `buf.yaml`:

```yaml
deps:
  - buf.build/googleapis/googleapis
```

Then run:

```bash
buf dep update
```

This creates or updates:

```text
buf.lock
```

`buf.lock` should be committed.

Use external dependencies deliberately for shared Google APIs or well-known schema packages.

---

## Remote Plugins

Use Buf remote plugins for repeatable generation.

Benefits:

```text
less local machine setup
consistent plugin versions
cleaner CI
fewer protoc plugin installation issues
```

For Go/gRPC:

```yaml
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen/go
    opt:
      - paths=source_relative

  - remote: buf.build/grpc/go
    out: gen/go
    opt:
      - paths=source_relative
```

---

## Generated Code Policy

Generated code lives under:

```text
gen/go/
```

Rules:

```text
Do not hand-edit generated files.
Regenerate after changing .proto files.
Commit generated code if this is the chosen repo policy.
CI should fail if generated code is stale.
Services should import generated code from gen/go paths.
```

Example import:

```go
import catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
```

---

## Local API Testing

Use `buf curl` or `grpcurl` for manual API testing.

Example:

```bash
buf curl \
  --schema proto \
  --protocol grpc \
  --http2-prior-knowledge \
  --data '{"product_id":"cccccccc-cccc-cccc-cccc-cccccccc0001"}' \
  http://localhost:50051/bfstore.catalog.v1.CatalogService/GetProduct
```

Use generated clients for application code.

Use CLI tools for smoke tests, debugging, and learning.

---

## Buf Schema Registry

Publishing to the Buf Schema Registry is optional in early bfstore phases.

Initial priority:

```text
local Buf workflow
CI linting
CI breaking-change checks
repeatable generation
```

Later, consider publishing:

```text
buf.build/mantrobuslawal/bfstore
```

Benefits later:

```text
hosted schema documentation
module dependency management
schema history
SDK generation options
professional API governance story
```

---

## Practical Rules

```text
Use Buf as the standard Protobuf toolchain.
Keep proto files under proto/.
Use versioned packages.
Use buf.yaml for contract policy.
Use buf.gen.yaml for generation policy.
Run buf lint before committing.
Run buf breaking in CI.
Run buf generate after contract changes.
Do not hand-edit generated files.
Use remote plugins for repeatable generation.
Commit buf.lock when dependencies are used.
Wrap commands in Makefile targets.
```

---

## Final Rule

```text
Buf keeps bfstore Protobuf contracts boring, governed, and safe to evolve.
```

That is exactly what production-grade contract-first development should look like.
