# OpenTelemetry Distributions

This document defines the bfstore policy for OpenTelemetry distributions.

A distribution is a customised, packaged version of an upstream OpenTelemetry component. For bfstore, the most relevant component is the OpenTelemetry Collector.

---

## Purpose

This document explains:

```text
what OpenTelemetry distributions are
distribution vs fork
Pure, Plus, and Minus distribution types
bfstore maturity model
vendor lock-in policy
when not to build a custom distribution
when to consider one
supply-chain and security expectations
```

---

## Distribution vs Fork

A distribution packages or wraps upstream OpenTelemetry with customisation.

A fork copies the source code and creates a separate maintenance path.

bfstore policy:

```text
Prefer distributions.
Avoid forks.
```

Forking OpenTelemetry would create unnecessary maintenance, upgrade, and security patching burden unless there is a very strong reason.

---

## Distribution Types

### Pure

A Pure distribution keeps the same functionality as upstream and focuses on packaging, defaults, compatibility, or ease of use.

bfstore example:

```text
standard Collector image
bfstore Collector config
bfstore Docker Compose/Kubernetes packaging
```

### Plus

A Plus distribution adds functionality beyond the standard upstream component.

bfstore should avoid Plus distributions unless there is a clear operational need.

### Minus

A Minus distribution removes components and provides a smaller approved subset.

bfstore may consider this later for:

```text
security hardening
smaller image size
reduced attack surface
supportability
approved component control
```

---

## Current bfstore Policy

Current recommendation:

```text
Use the upstream OpenTelemetry Collector image.
Store Collector config in repo.
Use OTLP from Go services.
Use debug exporter locally.
Add Prometheus/Tempo/Loki routing later.
Do not build a custom Collector binary yet.
```

Practical rule:

```text
Configuration first.
Custom distribution later.
```

---

## Maturity Model

### Phase 1: Standard upstream Collector

```text
official Collector image
OTLP receiver
batch processor
debug exporter
local Docker Compose
```

### Phase 2: Version-pinned Collector

```text
pin Collector image version
pin config in repo
validate config in CI
document chosen components
```

### Phase 3: Hardened Collector config

```text
memory_limiter processor
health_check extension
restricted receivers/exporters
security review
```

### Phase 4: Custom Collector distribution, if justified

```text
build with Collector Builder
include only approved components
scan image
sign image
generate SBOM
version release
```

### Phase 5: Kubernetes production-style operation

```text
deploy with Helm/Kustomize/GitOps
monitor Collector internal telemetry
document upgrade process
document rollback plan
```

---

## When To Consider a Custom Distribution

Consider a custom Collector distribution when there is a real operational reason.

Examples:

```text
security:
  remove unused receivers/exporters/extensions

supply chain:
  pin exact component versions
  scan and sign the image
  produce an SBOM

operability:
  standardise one Collector binary for all environments

performance:
  remove unused components
  reduce memory footprint

compliance:
  restrict outbound exporters
  ensure only approved telemetry paths exist
```

Do not build a custom distribution just because it sounds advanced.

Rule:

```text
Custom distributions are for control, not decoration.
```

---

## Vendor Lock-in Policy

bfstore should remain portable.

Recommended:

```text
use OTLP from services
send telemetry to OpenTelemetry Collector
avoid vendor-only SDK APIs
keep exporters replaceable
document exporter choices
```

Avoid:

```text
service code tied directly to one vendor backend
vendor-specific SDKs where OpenTelemetry APIs would work
unexplained vendor Collector distributions
```

---

## Collector Component Policy

Start with standard components:

```text
receivers:
  otlp

processors:
  batch

exporters:
  debug
```

Later hardening may add:

```text
processors:
  memory_limiter
  tail_sampling

extensions:
  health_check
  pprof
  zpages
```

Only add components when they serve a documented purpose.

---

## Testing Expectations

For the standard Collector phase:

```text
Collector config validates
Collector container starts
OTLP gRPC receiver works
OTLP HTTP receiver works
debug exporter shows traces
service.name appears correctly
metrics are exported correctly
logs pipeline works if enabled
```

For a future custom distribution:

```text
binary starts
only approved components are present
config validates against the custom binary
image scan passes
SBOM generated
image signed
upgrade process documented
rollback tested
```

---

## ADR Guidance

Consider these ADRs later:

```text
docs/adr/00xx-use-standard-opentelemetry-collector-before-custom-distribution.md
docs/adr/00xx-build-custom-opentelemetry-collector-distribution.md
```

The first ADR may be useful now.

The second ADR should only be written if bfstore actually needs a custom Collector binary.

---

## Practical Rules

```text
Start with upstream Collector.
Do not fork OpenTelemetry.
Prefer Pure distribution behaviour first.
Consider Minus distribution later for security/supportability.
Avoid Plus distribution unless there is a clear need.
Keep OTLP as the service export path.
Keep backends replaceable.
Pin versions before production-style use.
Test Collector config in CI.
Scan, sign, and generate SBOMs for custom images if built.
Document vendor lock-in risks.
```

---

## Final Rule

```text
Use the standard Collector until the standard Collector is the problem.
```
