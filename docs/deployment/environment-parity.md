# Environment Parity

This document defines parity expectations across bfstore development, staging, and production-style environments.

---

## Purpose

This document explains:

```text
dev/staging/prod-style environment expectations
Kubernetes parity
config/secrets parity
promotion flow
```

---

## Core Rule

```text
Environments may differ in scale, but not in architectural truth.
```

---

## Environment Expectations

### Local

```text
Docker Compose
single replica services
local MySQL
local Kafka
local OpenTelemetry Collector
fake/simulated providers
```

### Dev

```text
deployed from main regularly
real container images
same config shape
same service contracts
basic smoke tests
```

### Staging

```text
Kubernetes deployment
multiple replicas where useful
realistic ConfigMaps/Secrets
realistic rollout/rollback
observability enabled
integration/smoke testing
```

### Production-style portfolio

```text
Kubernetes
service-owned backing services
observability dashboards
release management
rollback strategy
network exposure controls
documented operational behaviour
```

---

## Kubernetes Parity

Kubernetes-specific behaviour should be tested in Kubernetes environments.

Examples:

```text
Services
ConfigMaps
Secrets
probes
rolling updates
rollbacks
resource requests/limits
NetworkPolicy later
```

Compose proves service architecture. Kubernetes proves platform behaviour.

---

## Config and Secrets Parity

Use the same config interface across environments.

```text
environment variables
secret references
service addresses
resource handles
```

Avoid hard-coded local-only behaviour.

---

## Promotion Flow

Recommended flow:

```text
local validation
CI validation
dev deploy
staging promotion
production-style promotion later
```

Promote the same build artefact where practical. Change config per environment.

---

## Practical Rules

```text
Keep config shape consistent.
Keep service contracts consistent.
Use realistic backing services.
Deploy frequently to dev.
Use staging to prove Kubernetes behaviour.
Document intentional differences.
Avoid environment-specific code paths.
```

---

## Final Rule

```text
A deployed environment should feel like the same bfstore system at a different scale.
```
