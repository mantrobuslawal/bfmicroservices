# Release Policy

This document defines bfstore release and artefact versioning policy.

## Core Rule

```text
A release tag should be a receipt, not a moving target.
```

## Git tags

```text
v0.8.0
v1.0.0-rc.1
v1.0.0
```

## Docker image tags

```text
ghcr.io/mantrobuslawal/bfstore/order-service:0.6.0
ghcr.io/mantrobuslawal/bfstore/order-service:sha-9f3a21c
ghcr.io/mantrobuslawal/bfstore/order-service:0.6.0-sha.9f3a21c
```

Avoid `latest` in controlled environments.

## Pre-release flow

```text
1.0.0-alpha.1:
  early testing

1.0.0-beta.1:
  feature complete

1.0.0-rc.1:
  release candidate

1.0.0:
  stable release
```

## Promotion flow

```text
build once
tag image immutably
deploy to dev
promote same artefact to staging
promote same artefact to production-style environment later
```

## Rollback expectations

Rollback requires:

```text
immutable release tags
compatible migrations
documented release notes
known previous good version
observable deployment state
```

## Practical rules

```text
Do not overwrite release tags.
Do not deploy latest to controlled environments.
Use SemVer and SHA tags.
Record release notes.
Document breaking changes.
Keep artefacts immutable.
```

## Final rule

```text
A release should be traceable, repeatable, and rollback-friendly.
```
