# Changelog Policy

This document defines how bfstore records release changes.

## Core Rule

```text
SemVer without a changelog is only half a communication system.
```

## Changelog location

```text
CHANGELOG.md
```

## Suggested sections

```text
Added
Changed
Deprecated
Removed
Fixed
Security
Migration Notes
```

## Breaking changes

Breaking changes must be clearly marked.

```text
## Breaking Changes

- Changed CheckoutRequest payment field shape.
- Consumers must update generated Protobuf clients.
- See migration notes below.
```

## Deprecations

Deprecations should include:

```text
what is deprecated
what replaces it
when it may be removed
migration guidance
```

## Migration notes

Include migration notes for:

```text
API changes
event schema changes
database migrations
config key changes
Helm/Terraform changes later
```

## Security fixes

Security fixes should be clearly marked without leaking exploit details unnecessarily.

```text
Fixed unsafe logging of payment provider error payloads.
```

## Practical rules

```text
Update changelog for every release.
Group changes by type.
Call out breaking changes loudly.
Document migration steps.
Record deprecations before removals.
Mention security fixes safely.
```

## Final rule

```text
A changelog turns version numbers into useful history.
```
