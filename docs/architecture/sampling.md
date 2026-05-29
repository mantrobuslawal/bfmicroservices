# Sampling

This document defines the trace sampling strategy for **bfstore**.

Sampling decides which traces are kept, exported, and stored when trace volume becomes too high to keep everything.

---

## Purpose

bfstore uses sampling to control trace volume, cost, and backend noise while preserving useful debugging visibility.

Sampling is mainly a trace concern. It does not replace metrics, logs, SLIs, or SLOs.

---

## Core Rule

Early bfstore environments should keep all traces until instrumentation behaviour is understood.

```text
local:
  keep 100% of traces

early dev/test:
  keep 100% of traces

load testing and production-like environments:
  introduce sampling deliberately
```

Practical rule:

```text
Do not sample before you know what normal telemetry looks like.
```

---

## When Sampling Is Useful

Sampling becomes useful when:

```text
catalogue browsing creates high trace volume
health checks create noisy traces
successful requests dominate storage
load tests generate too much telemetry
observability backend becomes expensive or noisy
```

---

## Head Sampling

Head sampling makes the decision near the start of a trace.

Use cases:

```text
simple percentage-based sampling
early volume protection
high-volume successful traffic reduction
```

Trade-off:

```text
simple and efficient
but may drop traces that later contain errors
```

Recommended bfstore use:

```text
sample successful high-volume catalogue traffic during load testing
```

---

## Tail Sampling

Tail sampling waits until most or all of a trace has arrived before deciding whether to keep it.

Use cases:

```text
keep all error traces
keep all slow traces
keep all payment timeout traces
keep important checkout traces
sample successful low-value traces
```

Trade-off:

```text
more powerful
but more operationally complex
```

Tail sampling should be introduced later through the OpenTelemetry Collector when trace volume justifies it.

---

## Environment Policy

### Local

```text
sample rate = 100%
```

Reason:

```text
learning
debugging
proving gRPC/Kafka/MySQL instrumentation
proving context propagation
```

### Dev/Test

```text
sample rate = 100% initially
```

Later during load tests:

```text
sample successful high-volume traces
keep error traces
```

### Staging

Keep:

```text
all error traces
slow checkout traces
payment timeout traces
traces from newly deployed versions at higher rate
```

Sample:

```text
successful catalogue browsing
health checks
metrics scrape routes
```

### Production-style Later

Use a combination of:

```text
head sampling for simple volume protection
tail sampling for errors, latency, and business-flow rules
metrics for complete aggregate visibility
```

---

## Recommended Retention Rules

Always keep where possible:

```text
grpc.status_code != OK
error.type exists
PaymentService/AuthorisePayment timeout
OrderService/Checkout duration > 3s
Kafka publish OrderCreated failed
failed checkout traces
```

Keep at higher rate:

```text
new service versions
staging environment
order-service traces
payment-service traces
checkout traces
```

Sample lower:

```text
successful CatalogService/GetProduct
successful ListProducts
health checks
metrics scrape routes
```

Example policy:

```text
CatalogService/GetProduct successful:
  keep 5%

OrderService/Checkout successful:
  keep 50%

OrderService/Checkout failed:
  keep 100%

PaymentService/AuthorisePayment DeadlineExceeded:
  keep 100%
```

---

## Sampling and Metrics

Sampling must not be used as the only source of system health.

Metrics should still preserve aggregate visibility:

```text
checkout.failed_total
payment.timeout_total
kafka.publish.failed_total
notification.delivery.failed_total
```

Rule:

```text
Metrics tell you how often something happened.
Traces show representative examples.
```

---

## Sampling and Propagation

Sampling decisions should propagate through trace context.

bfstore propagation boundaries:

```text
gRPC metadata
Kafka headers
traceparent
```

Goal:

```text
kept traces should remain connected across services
dropped traces should not create confusing partial visibility
```

---

## Collector Sampling

Later, bfstore may add Collector-based sampling.

Possible future file:

```text
deploy/observability/collector/otel-collector.sampling.yaml
```

Conceptual policy:

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

---

## Risks

Sampling can hide important traces if configured poorly.

Risks:

```text
rare checkout bug not sampled
payment timeout trace dropped
new deployment regression under-sampled
Kafka consumer failure not kept
partial traces confuse investigations
```

Rule:

```text
A sampling policy is an engineering decision, not just a percentage.
```

---

## Testing Expectations

Sampling tests should verify:

```text
local environment samples all traces
configured probability is applied
errors are retained where policy says so
slow checkout traces are retained where policy says so
health checks are sampled or dropped as expected
sampling decisions propagate across gRPC
sampling decisions propagate through Kafka headers
```

---

## Practical Rules

```text
Sample all traces locally.
Introduce sampling only after instrumentation is understood.
Use head sampling for simple early volume control.
Use tail sampling later for error, latency, and business-flow policies.
Keep failed checkout and payment timeout traces where possible.
Sample high-volume successful catalogue traces more aggressively.
Do not rely on sampled traces for aggregate health.
Make sampling environment-specific.
Monitor the Collector when using tail sampling.
Document sampling decisions.
```

---

## Final Rule

```text
Sampling keeps the interesting traces and reduces the boring ones.
```
