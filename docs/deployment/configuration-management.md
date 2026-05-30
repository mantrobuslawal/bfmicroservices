# Configuration Management

This document defines how bfstore manages configuration across deployment environments.

It complements:

```text
docs/architecture/configuration.md
docs/development/local-configuration.md
docs/security/secrets-management.md
```

---

## Purpose

This document explains:

```text
ConfigMaps
Secrets
Helm/Kustomize overlays
environment-specific values
build once deploy many
promotion between environments
```

---

## Core Rule

```text
Same artefact, different config.
```

The same Docker image should be deployable to multiple environments by changing configuration, not rebuilding source code.

---

## Environment Targets

bfstore may use:

```text
local
dev
test
staging
production-style
```

Each environment is a deploy target, not a separate codebase or branch.

---

## Kubernetes Config

Use:

```text
ConfigMap:
  non-secret config

Secret:
  sensitive config

Deployment:
  injects both as environment variables
```

Example:

```yaml
env:
  - name: BFSTORE_ENV
    value: staging

  - name: OTEL_EXPORTER_OTLP_ENDPOINT
    valueFrom:
      configMapKeyRef:
        name: observability-config
        key: otlp_endpoint

  - name: CATALOG_MYSQL_DSN
    valueFrom:
      secretKeyRef:
        name: catalog-db
        key: dsn
```

---

## Overlays

Possible Kubernetes layout:

```text
deploy/kubernetes/
├── base/
└── overlays/
    ├── local/
    ├── dev/
    ├── staging/
    └── prod/
```

Environment overlays should change config, replicas, resources, and routing — not application source code.

---

## Build Once, Deploy Many

Recommended flow:

```text
commit
  -> CI checks
  -> build Docker image
  -> tag image with commit SHA
  -> deploy to dev
  -> promote same image to staging
  -> promote same image to production-style environment
```

Avoid:

```text
rebuilding different images for each environment from different branches
manual production patches
environment-specific source code
```

---

## Promotion

Promotion should move a known artefact between environments.

Track:

```text
service name
image tag
git commit
environment
deployment time
configuration version
```

---

## Config Validation

Deployment pipelines should validate:

```text
required ConfigMaps exist
required Secrets exist
manifests render
no forbidden placeholder values remain
resource values are set
OTLP endpoints are environment-appropriate
```

---

## Practical Rules

```text
Use ConfigMaps for non-secret config.
Use Secrets or secret managers for sensitive config.
Use overlays for environment differences.
Do not create environment branches.
Do not rebuild source code per environment.
Promote artefacts between environments.
Track config and artefact versions.
Validate config before rollout.
```

---

## Final Rule

```text
Deployment config changes where bfstore runs, not what bfstore is.
```
