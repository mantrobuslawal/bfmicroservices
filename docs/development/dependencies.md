# Dependencies

This document defines how **bfstore** declares and isolates application, tooling, and runtime dependencies.

It is based on the 12-Factor principle:

```text
Explicitly declare and isolate dependencies.
```

---

## Purpose

This document defines bfstore policy for:

```text
runtime dependencies
build dependencies
test dependencies
tool dependencies
system tool policy
Docker dependency policy
CI dependency policy
dependency isolation
```

---

## Core Rule

```text
If bfstore needs it to build, test, lint, generate, migrate, package, deploy, or run,
the dependency must be declared somewhere.
```

A dependency that only exists on one developer’s laptop is not acceptable.

---

## Go Dependencies

Go dependencies are declared in:

```text
go.mod
go.sum
```

Rules:

```text
commit go.mod
commit go.sum
run go mod tidy before commits
use go mod verify in CI where useful
avoid undocumented local replacements
```

Recommended starting approach:

```text
use a root Go module for the bfstore monorepo
revisit per-service modules later if needed
```

---

## Protobuf and Buf Dependencies

bfstore uses Protobuf for:

```text
gRPC APIs
Kafka event payloads
shared contracts
```

The Protobuf toolchain includes:

```text
buf
protoc-gen-go
protoc-gen-go-grpc
protobuf runtime libraries
gRPC runtime libraries
```

Recommended files:

```text
buf.yaml
buf.gen.yaml
go.mod
go.sum
tools.go
Makefile
CI workflow
```

Tool versions should be pinned where practical.

---

## Tool Dependencies

Tools may include:

```text
buf
golangci-lint
gosec
govulncheck
mockgen
migrate
grpcurl
kubectl
helm
terraform
trivy
syft
```

Tools should be declared in at least one of:

```text
tools.go
Makefile
mise.toml
.tool-versions
.devcontainer/devcontainer.json
docs/development/tooling.md
CI workflow
```

---

## Script Dependencies

Scripts must declare or check their required tools.

Example scripts:

```text
scripts/wait-for-mysql.sh
scripts/create-kafka-topics.sh
scripts/run-migrations.sh
scripts/smoke-test.sh
```

Possible requirements:

```text
bash
curl
jq
mysql
grpcurl
docker
kubectl
helm
```

A `make doctor` target should check important local tools.

---

## Docker Dependencies

Dockerfiles must declare what they need.

Rules:

```text
pin base images for production-style builds
avoid latest tags
use minimal runtime images where practical
do not rely on host-installed tools
document OS packages if required
```

Preferred pattern:

```text
multi-stage build
Go builder image
minimal runtime image
explicit binary copy
```

---

## CI Dependencies

CI must install or pin required tools.

CI should not depend on whatever happens to be pre-installed on the runner.

CI should declare:

```text
Go version
Buf version
Protobuf plugin versions
lint tool versions
security scanner versions
Docker build tooling
deployment tooling where applicable
```

Rule:

```text
CI should prove the repo is repeatable, not prove the runner got lucky.
```

---

## Runtime Dependencies vs Backing Services

Runtime dependencies are libraries/tools the app needs:

```text
gRPC library
MySQL driver
Kafka client
OpenTelemetry SDK
```

Backing services are external services the app connects to:

```text
MySQL
Kafka
OpenTelemetry Collector
payment provider simulator
email provider
```

Rule:

```text
The driver is a dependency.
The database is a backing service.
```

---

## Isolation

Dependency isolation may come from:

```text
Go modules
Docker containers
dev containers
pinned CI images
local tools directory
Makefile targets
mise/asdf tool version files
```

Do not rely on random globally installed tools unless explicitly documented and checked.

---

## Local Development

Recommended local flow:

```text
make doctor
make bootstrap
make test
make up
```

`make doctor` should check required tools.

`make bootstrap` may install or prepare local dependencies.

---

## Practical Rules

```text
Declare application dependencies.
Declare tool dependencies.
Declare script dependencies.
Pin versions where practical.
Avoid latest tags in production-style workflows.
Use Go modules properly.
Keep go.sum committed.
Use Docker/dev containers for isolation where useful.
Make CI setup explicit.
Use make doctor to catch missing tools.
```

---

## Final Rule

```text
bfstore should work because the repo declares what it needs, not because a laptop happens to be prepared.
```
