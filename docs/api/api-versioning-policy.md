# API Versioning Policy

This document defines versioning rules for bfstore gRPC and Protobuf APIs.

## Core Rule

```text
Protobuf contracts are public API.
```

## Package versions

Use versioned Protobuf packages.

```proto
package bfstore.catalog.v1;
```

Breaking API changes should move to a new major package version.

```proto
package bfstore.catalog.v2;
```

## PATCH changes

```text
fix comments
fix generated documentation
fix server implementation without changing contract
```

## MINOR changes

```text
add optional field
add new message
add new RPC method
add new enum value carefully
mark field as deprecated
```

## MAJOR changes

```text
remove field
rename field in a way that breaks generated APIs
change field type
reuse field number
remove RPC
change request/response shape incompatibly
change documented meaning incompatibly
```

## Field numbering

```text
never reuse field numbers
reserve removed field numbers
reserve removed field names where appropriate
prefer optional for scalar presence where needed
```

## Deprecation path

```text
minor release:
  mark field/RPC deprecated
  document replacement

major release:
  remove deprecated field/RPC
```

## Buf checks

Use Buf for:

```text
linting
generation
breaking-change detection
```

## Final rule

```text
A Protobuf change is not safe just because the compiler accepts it.
```
