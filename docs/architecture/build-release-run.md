# Build, Release, Run

This document defines how **bfstore** separates build, release, and run stages.

The principle is:

```text
build
  -> release
      -> run
```

---

## Purpose

This document explains:

```text
build/release/run stages
immutable release model
build once deploy many
service-specific lifecycle
config separation
runtime responsibilities
```

---

## Core Rule

```text
Build creates an artefact.
Release combines artefact with config.
Run starts the selected release.
```

---

## Build Stage

Build turns source code into executable artefacts.

For bfstore, build may include:

```text
checkout source
download dependencies
run tests
run linting
run buf lint
run buf breaking checks
generate code if needed
compile Go binaries
build Docker images
tag images with commit SHA
push images to registry
```

Build must not bake environment-specific config into images.

---

## Release Stage

Release combines a build artefact with environment config.

A release includes:

```text
service image tag
environment config
ConfigMap references
Secret references
service addresses
replica/resource configuration
release metadata
```

Example:

```text
service: order-service
environment: staging
image: ghcr.io/mantrobuslawal/bfstore/order-service:abc123
config version: staging-config-v8
release id: order-service-staging-20260530-2215
```

Rule:

```text
A release is image plus config.
```

---

## Run Stage

Run starts the selected release.

Runtime should:

```text
start process
load env vars
validate config
connect to backing services
serve gRPC/HTTP
consume Kafka events
emit telemetry
respond to health checks
shutdown gracefully
```

Runtime should not:

```text
compile code
download dependencies
generate protobuf code
install packages
patch source files
mutate release config
```

Rule:

```text
Runtime should be boring.
```

---

## Build Once, Deploy Many

Recommended flow:

```text
commit abc123
  -> build image once
  -> deploy same image to dev
  -> promote same image to staging
  -> promote same image to production-style environment
```

Environment differences should come from config, not new builds.

---

## Immutable Releases

Each release should have a unique ID and should not be changed after creation.

If image or config changes, create a new release.

Record:

```text
service name
environment
image tag
git commit
config version
release timestamp
migration version if applicable
```

---

## Service-specific Lifecycle

Each deployable service should have its own traceable lifecycle.

Examples:

```text
catalog-service:
  build image
  release with CATALOG_MYSQL_DSN and OTEL config
  run pods

order-service:
  build image
  release with ORDER_MYSQL_DSN, Kafka, dependency addresses, OTEL config
  run pods

notification-service:
  build image
  release with Kafka and SMTP config
  run consumers
```

---

## Database Migrations

Migrations should be handled deliberately.

Avoid:

```text
every app pod automatically running migrations at startup
```

Prefer:

```text
dedicated migration job
CI/CD controlled migration step
backward-compatible migrations
expand/contract migration pattern later
```

Rule:

```text
App runtime should not surprise-migrate production databases.
```

---

## Kafka and Protobuf Compatibility

bfstore releases can affect other services through Kafka and Protobuf contracts.

Release safety should include:

```text
Buf breaking-change checks
event compatibility rules
schema evolution discipline
consumer backward compatibility
topic versioning where needed
```

---

## OpenTelemetry Release Metadata

Each service should emit release metadata.

Example attributes:

```text
service.name=bfstore-order-service
service.version=abc123
deployment.environment.name=staging
```

This helps identify which release introduced latency, errors, retries, or failed calls.

---

## Practical Rules

```text
Separate build, release, and run.
Build from a specific commit.
Do not bake deploy config into images.
Tag images with commit SHA.
Treat releases as image plus config.
Give each release an ID.
Promote artefacts across environments.
Keep runtime simple.
Handle migrations deliberately.
Make running services traceable to commits and releases.
```

---

## Final Rule

```text
bfstore deploys should be repeatable, traceable, and boring at runtime.
```
