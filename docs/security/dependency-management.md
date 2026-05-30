# Dependency Management

This document defines bfstore dependency management and security expectations.

It complements:

```text
docs/development/dependencies.md
docs/development/tooling.md
```

---

## Purpose

Dependency management is part of bfstore’s DevSecOps posture.

This document defines expectations for:

```text
dependency scanning
vulnerability management
image scanning
SBOM planning
update policy
license review later
version pinning
supply-chain hygiene
```

---

## Core Rule

```text
You cannot secure dependencies you have not declared.
```

---

## Go Dependency Hygiene

Recommended:

```text
go mod tidy
go mod verify
govulncheck
gosec
Dependabot or Renovate
```

Rules:

```text
commit go.sum
review dependency changes
avoid unnecessary dependencies
remove unused dependencies
do not ignore vulnerability alerts without explanation
```

---

## Container Image Hygiene

Recommended:

```text
minimal runtime images
pinned base images
image scanning with Trivy or Grype
avoid latest tags
multi-stage builds
```

Later:

```text
SBOM generation with Syft
image signing with Cosign
provenance/SLSA-style metadata
```

---

## Tool Dependency Hygiene

Tool versions should be intentional.

Important tools:

```text
buf
protoc-gen-go
golangci-lint
gosec
govulncheck
trivy
syft
cosign
terraform
kubectl
helm
```

Avoid unpinned installer scripts where possible.

---

## Vulnerability Management

When a vulnerability is found:

```text
identify affected dependency
check exploitability in bfstore context
upgrade if possible
document risk if not immediately fixed
add compensating controls where needed
verify with tests/scans
```

Do not blindly suppress findings.

---

## Update Policy

Dependency updates should be:

```text
regular
reviewed
tested
scanned
committed with clear messages
```

Possible cadence:

```text
patch/security updates:
  as soon as practical

minor updates:
  regular maintenance window

major updates:
  planned and tested
```

---

## License Review

Later, bfstore may add license checks.

Watch for:

```text
unexpected copyleft dependencies
unknown licences
commercial-use restrictions
abandoned packages
```

---

## SBOM Plan

Later bfstore may generate SBOMs for:

```text
service container images
custom CI/tooling images
custom OpenTelemetry Collector distribution if built
```

Possible tool:

```text
syft
```

---

## CI Expectations

CI should eventually include:

```text
go test
go vet
golangci-lint
govulncheck
gosec
container image scan
dependency review
secret scan
SBOM generation later
```

---

## Practical Rules

```text
Keep dependency manifests committed.
Scan dependencies.
Scan images.
Pin important versions.
Avoid latest tags.
Review updates.
Remove unused dependencies.
Do not suppress vulnerabilities casually.
Generate SBOMs later for production-style artefacts.
Treat dependency management as security work.
```

---

## Final Rule

```text
Dependency management is not admin; it is supply-chain security.
```
