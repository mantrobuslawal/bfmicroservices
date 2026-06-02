# Effective Go Guide

This document defines how bfstore applies Effective Go principles across Go services.

## Purpose

```text
how bfstore applies Effective Go principles
what is modernised or supplemented
where Effective Go is intentionally not enough
```

## Core Rule

```text
Write Go that is boring, clear, readable, and unsurprising.
```

## Applies To

```text
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-worker
bfstore-admin
pkg/platform/*
```

## bfstore Effective Go Principles

```text
use gofmt
use short package names
use MixedCaps
use explicit errors
avoid unnecessary else after return
keep main.go thin
define interfaces at consumer boundaries
prefer behaviour methods over generic setters
use defer for cleanup
use make for slices, maps, and channels
use goroutines carefully
use channels only for in-process coordination
use Kafka for durable cross-service events
```

## Modern Supplements

Effective Go should be combined with:

```text
Go Modules Reference
How to Write Go Code
Go Code Review Comments
modern context usage
structured logging
OpenTelemetry
govulncheck
staticcheck
gRPC/Kafka package docs
```

## Practical Rules

```text
Prefer simple packages.
Prefer explicit construction.
Prefer clear errors.
Prefer domain types.
Do not let transport types swallow the domain.
Avoid cleverness that makes incidents harder.
```

## Final Rule

```text
Effective Go should make bfstore easier to read, test, review, and operate.
```
