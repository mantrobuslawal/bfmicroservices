# Build Artefacts

This document defines bfstore build artefact expectations.

It complements:

```text
docs/architecture/build-release-run.md
docs/deployment/release-management.md
```

---

## Purpose

This document explains:

```text
Docker images
commit SHA tags
generated Protobuf code
binary naming
SBOM/image scanning later
registry naming
```

---

## Core Rule

```text
Build artefacts must be traceable to source code.
```

---

## Docker Images

Each deployable service should produce a Docker image.

Examples:

```text
ghcr.io/mantrobuslawal/bfstore/catalog-service:abc123
ghcr.io/mantrobuslawal/bfstore/order-service:abc123
ghcr.io/mantrobuslawal/bfstore/payment-service:abc123
```

Avoid:

```text
latest
local-only tags for real environments
manually patched images
```

---

## Commit SHA Tags

Images should be tagged with the Git commit SHA.

Benefits:

```text
traceability
rollback
auditability
release comparison
observability correlation
```

Possible format:

```text
service-name:<short-sha>
service-name:v0.3.0-<short-sha>
```

---

## Go Binaries

Binary names should be predictable.

Examples:

```text
catalog-service
order-service
payment-service
notification-service
```

Build output may live in:

```text
bin/
out/
Docker build stage only
```

---

## Generated Protobuf Code

The source of truth is the `.proto` file.

Generated code may be:

```text
generated during build/CI
committed and verified as up to date
```

Do not manually edit generated files.

Build should run:

```text
buf lint
buf breaking
buf generate where required
```

---

## Image Scanning

Later, bfstore should scan images with tools such as:

```text
Trivy
Grype
```

Scan for:

```text
known vulnerabilities
secrets accidentally included
large/unnecessary packages
outdated base images
```

---

## SBOM Plan

Later, bfstore may generate SBOMs for service images.

Possible tool:

```text
Syft
```

SBOMs help with:

```text
supply-chain visibility
dependency auditing
client-facing DevSecOps evidence
incident response
```

---

## Registry Naming

Preferred registry:

```text
ghcr.io/mantrobuslawal/bfstore/<service-name>:<tag>
```

Examples:

```text
ghcr.io/mantrobuslawal/bfstore/catalog-service:abc123
ghcr.io/mantrobuslawal/bfstore/order-service:abc123
```

---

## Practical Rules

```text
Build images from specific commits.
Tag images with commit SHA.
Keep artefacts immutable.
Avoid latest tags.
Do not bake environment config into images.
Keep generated code traceable to proto source.
Scan images later.
Generate SBOMs later for production-style evidence.
```

---

## Final Rule

```text
An artefact should tell you exactly what source produced it.
```
