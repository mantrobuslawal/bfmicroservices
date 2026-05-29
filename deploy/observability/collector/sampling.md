# Collector Sampling

This document describes how the OpenTelemetry Collector may apply sampling for **bfstore**.

It complements:

```text
docs/architecture/sampling.md
pkg/platform/telemetry/sampling.md
deploy/observability/collector/README.md
```

---

## Purpose

Collector sampling provides centralised sampling policy for bfstore telemetry.

It is especially useful for:

```text
tail sampling
error-based retention
latency-based retention
attribute-based policies
environment-specific policy
reducing backend trace volume
```

---

## When To Use Collector Sampling

Do not start here on day one.

Use Collector sampling when:

```text
trace volume becomes high
backend storage/query cost becomes a concern
load tests generate too many traces
you need to keep errors and slow traces reliably
you want central policy rather than service-by-service configuration
```

---

## Recommended Progression

```text
Phase 1:
  no Collector sampling
  debug/export everything locally

Phase 2:
  sample all traces in dev
  validate gRPC/Kafka/MySQL spans

Phase 3:
  use simple SDK head sampling during load tests if needed

Phase 4:
  introduce Collector tail sampling for production-style environments
```

---

## Tail Sampling Policies

Future policies may include:

```text
keep all error traces
keep slow checkout traces
keep payment timeout traces
sample successful catalogue traffic
sample health checks aggressively
```

Conceptual config:

```yaml
processors:
  tail_sampling:
    decision_wait: 10s
    num_traces: 50000
    policies:
      - name: keep-errors
        type: status_code
        status_code:
          status_codes: [ERROR]

      - name: keep-slow-checkout
        type: latency
        latency:
          threshold_ms: 3000

      - name: sample-successful-catalogue
        type: probabilistic
        probabilistic:
          sampling_percentage: 5
```

This example is a starting concept, not final production config.

---

## Operational Considerations

Tail sampling needs resources.

Monitor:

```text
Collector CPU
Collector memory
dropped spans
queue length
exporter failures
sampling processor errors
trace decision latency
```

Tail sampling requires enough memory to hold traces while decisions are pending.

---

## Data Safety

Sampling does not make unsafe telemetry safe.

Services should not emit:

```text
raw JWTs
passwords
API keys
payment card numbers
CVV
full shipping addresses
customer email
full basket JSON
raw Kafka payloads
full SQL with sensitive values
```

Collector filtering can be a safety net, not the main defence.

---

## Testing Guidance

Test that Collector sampling:

```text
keeps error traces
keeps slow checkout traces
keeps payment timeout traces
samples successful catalogue traces
does not break trace continuity
does not overload Collector memory
exports expected traces to backend
```

---

## Practical Rules

```text
Document sampling policy before writing complex YAML.
Start with no sampling locally.
Use Collector sampling for centralised tail policies.
Monitor Collector resource usage.
Keep failed checkout and payment timeout traces where possible.
Sample boring successful traffic more aggressively.
Do not rely on sampling to fix unsafe instrumentation.
```

---

## Final Rule

```text
Collector sampling is a later-stage control plane for trace volume.
```
