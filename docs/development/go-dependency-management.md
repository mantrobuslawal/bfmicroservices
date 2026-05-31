# Go Dependency Management

This document defines bfstore's Go dependency management conventions.

## Core Rule

```text
go.mod declares dependency intent.
go.sum protects dependency integrity.
Commit both.
```

## Adding Dependencies

After adding imports, run:

```bash
go mod tidy
```

Then check:

```bash
git diff -- go.mod go.sum
```

## Removing Dependencies

After removing imports, run:

```bash
go mod tidy
```

## CI Check

CI should check that dependency files are tidy:

```bash
go mod tidy
git diff --exit-code go.mod go.sum
```

## Module Cache

Useful command:

```bash
go clean -modcache
```

Use this when dependency cache issues need a clean fetch.

## Dependency Review

Review dependencies for:

```text
license
maintenance
security advisories
transitive dependency weight
necessity
```

## Tooling Dependencies

For tools such as linters or generators, prefer pinned versions and documented installation.

Possible tools:

```text
buf
staticcheck
govulncheck
protoc plugins
```

## Practical Rules

```text
Run go mod tidy after dependency changes.
Commit go.mod and go.sum.
Do not add dependencies casually.
Review dependency purpose and risk.
Use caches for speed, not correctness.
Keep tool versions documented.
```

## Final Rule

```text
Dependency management is part of supply-chain hygiene, not admin paperwork.
```
