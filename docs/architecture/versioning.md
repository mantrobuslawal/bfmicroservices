# Versioning

This document defines how **bfstore** uses versioning across services, APIs, events, images, modules, and release artefacts.

## Purpose

```text
what SemVer means across bfstore
what counts as public API
how service, API, event, image, and module versions relate
```

## Core Rule

```text
A version number should communicate compatibility.
```

## Public API in bfstore

Public API can include:

```text
gRPC service methods
Protobuf messages and fields
Kafka event schemas
Kafka topic/key semantics
Docker image commands and environment variables
Helm chart values later
Terraform module inputs/outputs later
Go package exported APIs if reused
```

## Version types

```text
Project version:
  overall bfstore platform milestone

Service version:
  deployable service artefact

API version:
  gRPC/Protobuf compatibility

Event version:
  Kafka event schema/topic compatibility

Image version:
  immutable container artefact
```

## Practical rules

```text
Use SemVer for release communication.
Use explicit API/event package versions for wire contracts.
Use immutable Git and Docker tags.
Use commit SHAs for traceability.
Document breaking changes.
Keep a changelog.
```

## Final rule

```text
Versioning should make compatibility visible before deployment.
```
