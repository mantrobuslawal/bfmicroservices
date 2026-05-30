# Semantic Versioning

This document defines bfstore's Semantic Versioning policy for development.

## Core Rule

```text
SemVer is a compatibility promise.
```

## Version format

```text
MAJOR.MINOR.PATCH
```

## PATCH

Use PATCH for backward-compatible bug fixes.

```text
1.4.2 -> 1.4.3
```

Examples:

```text
fix logging timestamp
fix catalogue sorting
fix retry error handling
```

## MINOR

Use MINOR for backward-compatible additions.

```text
1.4.2 -> 1.5.0
```

Examples:

```text
add optional Protobuf field
add new gRPC method without breaking existing clients
add new admin CLI command
deprecate an existing field/API
```

## MAJOR

Use MAJOR for breaking public API changes.

```text
1.4.2 -> 2.0.0
```

Examples:

```text
remove RPC
remove Protobuf field
change field type
rename required config key
change Kafka event key semantics
```

## 0.y.z

Use `0.y.z` while bfstore contracts are still unstable.

## Pre-release versions

```text
1.0.0-alpha.1
1.0.0-beta.1
1.0.0-rc.1
```

## Deprecation

```text
Deprecate in MINOR.
Remove in MAJOR.
```

## Practical rules

```text
Do not make random version bumps.
Do not hide breaking changes in patch releases.
Document breaking changes.
Use changelogs.
Do not overwrite released tags.
```

## Final rule

```text
Version numbers should tell the truth.
```
