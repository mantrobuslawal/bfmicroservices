# OpenTelemetry Collector Distributions

This document describes how bfstore should approach OpenTelemetry Collector distributions.

It complements:

```text
docs/architecture/opentelemetry-distributions.md
deploy/observability/collector/README.md
deploy/observability/collector/sampling.md
```

---

## Purpose

The Collector distribution policy exists so bfstore can start simple and only introduce custom Collector builds when they solve a real operational, security, performance, or supply-chain problem.

---

## Current Approach

bfstore should initially use an upstream OpenTelemetry Collector image.

Recommended starting point:

```text
official upstream Collector image
OTLP receiver
batch processor
debug exporter
version-controlled Collector config
```

Do not build a custom Collector binary in the early stages.

---

## Recommended Directory Structure

Current or near-future structure:

```text
deploy/observability/collector/
├── README.md
├── distributions.md
├── sampling.md
├── otel-collector.local.yaml
└── otel-collector.k8s.yaml
```

Future custom distribution structure, only if justified:

```text
deploy/observability/collector/
├── builder-config.yaml
├── Dockerfile
├── README.md
└── distributions.md
```

---

## Approved Initial Components

Start with:

```text
receivers:
  otlp

processors:
  batch

exporters:
  debug
```

Later consider:

```text
processors:
  memory_limiter
  tail_sampling

exporters:
  prometheus
  otlp

extensions:
  health_check
  pprof
  zpages
```

Each component should have a documented reason.

---

## When To Build a Custom Collector Distribution

Only consider a custom distribution when one or more of these applies:

```text
unused components need to be removed
security team wants a smaller approved binary
image size or attack surface matters
supply-chain control is required
SBOM/signing/scanning is part of the release process
a single standardised Collector binary is needed across environments
```

Do not create one just because custom sounds impressive.

---

## Collector Builder

A future custom distribution may use the OpenTelemetry Collector Builder.

Conceptual builder config:

```yaml
dist:
  name: bfstore-otel-collector
  description: bfstore OpenTelemetry Collector distribution
  output_path: ./dist/bfstore-otel-collector

receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver vX.Y.Z

processors:
  - gomod: go.opentelemetry.io/collector/processor/batchprocessor vX.Y.Z
  - gomod: go.opentelemetry.io/collector/processor/memorylimiterprocessor vX.Y.Z

exporters:
  - gomod: go.opentelemetry.io/collector/exporter/debugexporter vX.Y.Z
```

This is intentionally future-facing and should not be treated as immediate implementation work.

---

## Version Pinning

Before production-style use:

```text
pin Collector image version
pin component versions if building custom distribution
record upgrade notes
test configuration before rollout
document rollback path
```

Avoid floating tags such as:

```text
latest
```

for production-style environments.

---

## Security and Supply Chain Expectations

If bfstore later builds a custom Collector image:

```text
scan the image
generate an SBOM
sign the image
pin dependencies
store build config in Git
document upgrade procedure
document rollback procedure
```

A custom image without this discipline is not a platform asset.

---

## Testing Expectations

Standard Collector phase:

```text
Collector config validates
Collector container starts
OTLP gRPC receiver works
OTLP HTTP receiver works
debug exporter receives traces
service.name appears correctly
metrics pipeline works if enabled
logs pipeline works if enabled
```

Custom distribution phase:

```text
custom binary starts
only approved components are present
config validates against custom binary
image scan passes
SBOM generated
image signed
upgrade tested
rollback tested
```

---

## Vendor Distribution Policy

Vendor distributions may be useful, but must be evaluated carefully.

Check:

```text
does it preserve OTLP compatibility?
does it create vendor lock-in?
does it require vendor-specific SDK APIs?
who supports it?
can bfstore swap backends later?
```

bfstore should prefer portability.

---

## Practical Rules

```text
Start with upstream Collector.
Keep Collector config version-controlled.
Pin versions before production-style use.
Only add components with a documented purpose.
Do not fork OpenTelemetry.
Build a custom distribution only for real control needs.
Scan/sign/SBOM custom images.
Keep service telemetry exported via OTLP.
Keep observability backends replaceable.
```

---

## Final Rule

```text
Collector distributions are a later-stage hardening option, not a day-one requirement.
```
