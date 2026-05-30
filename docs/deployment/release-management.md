# Release Management

This document defines release management expectations for bfstore.

It complements:

```text
docs/architecture/build-release-run.md
docs/deployment/rollback-strategy.md
docs/development/build-artifacts.md
```

---

## Purpose

This document explains:

```text
release IDs
image tags
config versions
environment promotion
GitOps flow
migration handling
release metadata
```

---

## Core Rule

```text
A release is a build artefact plus deploy configuration.
```

---

## Release Metadata

Each release should record:

```text
service name
environment
image tag
git commit
config version
release ID
release timestamp
migration version if applicable
who/what created the release
```

Example:

```text
service: order-service
environment: staging
image: ghcr.io/mantrobuslawal/bfstore/order-service:abc123
git commit: abc123
config version: staging-config-v8
release id: order-service-staging-20260530-2215
```

---

## Image Tags

Preferred image tags:

```text
commit SHA
semantic version plus commit SHA
```

Examples:

```text
ghcr.io/mantrobuslawal/bfstore/order-service:abc123
ghcr.io/mantrobuslawal/bfstore/order-service:v0.3.0-abc123
```

Avoid:

```text
latest
current
final
```

---

## Config Versions

Config should be versioned or traceable.

Sources may include:

```text
GitOps environment folders
Helm values
Kustomize overlays
ConfigMaps
Secret references
cloud secret versions
```

---

## Promotion

Promote the same artefact through environments.

Example:

```text
dev:
  order-service:abc123 + dev config

staging:
  order-service:abc123 + staging config

production-style:
  order-service:abc123 + production-style config
```

Do not rebuild separate images for each environment.

---

## GitOps Flow

With GitOps, release state can live in Git.

Example:

```text
update deploy/kubernetes/overlays/staging image tag to abc123
commit change
Argo CD syncs environment
Kubernetes runs selected release
```

Rule:

```text
GitOps makes release state reviewable and auditable.
```

---

## Migration Handling

Releases may include database schema expectations.

Migration handling should document:

```text
migration version
whether migration is backward-compatible
whether rollback is possible
whether migration must run before app rollout
```

Prefer:

```text
dedicated migration job
controlled pipeline step
expand/contract pattern later
```

---

## Release Validation

Before marking a release healthy, verify:

```text
pods are ready
health checks pass
gRPC endpoints respond
Kafka producers/consumers work
database connectivity works
OpenTelemetry emits service.version
smoke tests pass
error rates are acceptable
```

---

## Practical Rules

```text
Give each release an ID.
Record image tag and commit SHA.
Record config version.
Promote artefacts, not branches.
Avoid latest tags.
Use GitOps to make release state visible where practical.
Validate releases after rollout.
```

---

## Final Rule

```text
A good release can answer: what code, what config, which environment, and when?
```
