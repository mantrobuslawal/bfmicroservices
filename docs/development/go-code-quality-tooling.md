# Go Code Quality Tooling

## Purpose

This document lists the Go tools recommended for checking, formatting, testing, securing, and maintaining Go code in the **bfstore** project.

The goal is to build a boring, repeatable quality gate that can run locally and later in CI.

> Keep it boring where production matters.

---

## Recommended Tooling Summary

| Area | Tool | Purpose |
|---|---|---|
| Formatting | `gofmt` | Formats Go code using the standard Go formatter |
| Imports | `goimports` | Formats Go code and organises imports |
| Suspicious code checks | `go vet` | Finds likely bugs and suspicious constructs |
| Tests | `go test` | Runs unit and integration tests |
| Race detection | `go test -race` | Detects data races at runtime |
| Coverage | `go test -cover` | Reports test coverage |
| Lint aggregation | `golangci-lint` | Runs multiple linters consistently |
| Static analysis | `staticcheck` | Finds bugs, simplifications, and performance issues |
| Vulnerability scanning | `govulncheck` | Checks for known vulnerabilities affecting Go code |
| Module hygiene | `go mod tidy` | Cleans and verifies module dependencies |
| Dependency graph | `go mod graph` | Shows module dependency graph |
| Build verification | `go build ./...` | Confirms packages compile |
| Generated code check | `buf generate` | Generates Protobuf/gRPC code where relevant |
| Protobuf linting | `buf lint` | Lints Protobuf contracts |
| Protobuf compatibility | `buf breaking` | Checks for breaking Protobuf contract changes |

---

## Core Go Tools

## `gofmt`

Use `gofmt` to format all Go files.

```sh
gofmt -w .
```

Why it matters:

```text
consistent formatting
less style debate
cleaner diffs
standard Go appearance
```

Practical rule:

```text
All Go code should be gofmt-formatted before commit.
```

---

## `goimports`

Use `goimports` when you want formatting plus import management.

```sh
goimports -w .
```

Why it matters:

```text
adds missing imports
removes unused imports
groups imports consistently
runs gofmt behaviour too
```

Recommended local installation:

```sh
go install golang.org/x/tools/cmd/goimports@latest
```

Practical rule:

```text
Use goimports in editors and pre-commit workflows where possible.
```

---

## `go vet`

Use `go vet` to report suspicious constructs that compile but may be wrong.

```sh
go vet ./...
```

Why it matters:

```text
finds likely bugs
catches suspicious patterns
supports basic static analysis
should run in CI
```

Examples of issues it can detect:

```text
bad format strings
unreachable code patterns
misused struct tags
copying lock values
suspicious tests
```

Practical rule:

```text
go vet is not optional for production-facing Go services.
```

---

## Testing Tools

## `go test`

Run all tests:

```sh
go test ./...
```

Run tests with verbose output:

```sh
go test -v ./...
```

Run a single package:

```sh
go test ./internal/catalog
```

Run a single test:

```sh
go test ./internal/catalog -run TestServiceListProducts
```

Practical rule:

```text
Every service should have fast unit tests that run without Docker or cloud dependencies.
```

---

## Coverage

Run tests with coverage:

```sh
go test -cover ./...
```

Generate a coverage profile:

```sh
go test -coverprofile=coverage.out ./...
```

View coverage in the browser:

```sh
go tool cover -html=coverage.out
```

Practical rule:

```text
Coverage is useful, but do not worship the percentage.
High coverage with weak assertions is theatre.
```

---

## Race Detection

Run tests with the race detector:

```sh
go test -race ./...
```

Use this especially for code involving:

```text
goroutines
channels
shared maps
caches
background workers
Kafka consumers
HTTP/gRPC servers
connection pools
```

Practical rule:

```text
Run race detection before merging concurrency-heavy changes.
```

---

## Linting and Static Analysis

## `golangci-lint`

`golangci-lint` is a linter runner. It lets the project run many linters through one command and one configuration file.

Install:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

Run:

```sh
golangci-lint run ./...
```

Recommended config file:

```text
.golangci.yml
```

Suggested starting configuration:

```yaml
version: "2"

run:
  timeout: 5m

linters:
  enable:
    - govet
    - staticcheck
    - ineffassign
    - unused
    - errcheck
    - misspell
    - revive
    - gosec
    - bodyclose
    - noctx
    - nilerr
    - unconvert
    - unparam
    - prealloc

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
```

Practical rule:

```text
Start with a useful but not ridiculous linter set.
A noisy linter setup gets ignored.
```

---

## `staticcheck`

`staticcheck` is a strong static analysis tool for Go.

Install:

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
```

Run:

```sh
staticcheck ./...
```

If you use `golangci-lint`, Staticcheck can also run through that.

Practical rule:

```text
Use staticcheck either directly or through golangci-lint.
Do not ignore warnings without understanding them.
```

---

## Security Checks

## `govulncheck`

`govulncheck` reports known vulnerabilities that affect Go code.

Install:

```sh
go install golang.org/x/vuln/cmd/govulncheck@latest
```

Run:

```sh
govulncheck ./...
```

Practical rule:

```text
Run govulncheck locally before releases and in CI for service repositories.
```

---

## `gosec`

`gosec` checks Go code for common security issues.

It is often run through `golangci-lint`.

Direct installation:

```sh
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

Run:

```sh
gosec ./...
```

Practical rule:

```text
Use gosec findings as review prompts.
Some warnings need judgement, not blind suppression.
```

---

## Dependency and Module Hygiene

## `go mod tidy`

Run:

```sh
go mod tidy
```

Why it matters:

```text
removes unused dependencies
adds missing dependencies
keeps go.mod and go.sum clean
reduces dependency drift
```

Practical rule:

```text
Run go mod tidy before committing module changes.
```

---

## `go list`

List packages:

```sh
go list ./...
```

List module dependencies:

```sh
go list -m all
```

Practical rule:

```text
Use go list when debugging package/module structure.
```

---

## `go mod graph`

Inspect dependency graph:

```sh
go mod graph
```

Useful when investigating:

```text
unexpected dependency versions
transitive dependency risk
bloated dependency trees
security review questions
```

---

## Build Verification

## `go build`

Run:

```sh
go build ./...
```

For a service binary:

```sh
go build -o bin/catalog-service ./cmd/catalog-service
```

Practical rule:

```text
CI should prove the code builds, not just that tests pass.
```

---

## Recommended Local Command Sequence

For day-to-day development:

```sh
goimports -w .
go mod tidy
go test ./...
go vet ./...
golangci-lint run ./...
govulncheck ./...
```

For concurrency-heavy changes:

```sh
go test -race ./...
```

For release candidates:

```sh
goimports -w .
go mod tidy
go test -race -cover ./...
go vet ./...
golangci-lint run ./...
govulncheck ./...
go build ./...
```

---

## Suggested Makefile Targets

Add these targets to service-level or repo-level Makefiles where appropriate:

```makefile
.PHONY: fmt
fmt: ## Format Go code and organise imports
	goimports -w .

.PHONY: tidy
tidy: ## Tidy Go modules
	go mod tidy

.PHONY: test
test: ## Run Go tests
	go test ./...

.PHONY: test-race
test-race: ## Run Go tests with race detector
	go test -race ./...

.PHONY: coverage
coverage: ## Run tests and generate coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run ./...

.PHONY: vuln
vuln: ## Run govulncheck
	govulncheck ./...

.PHONY: build
build: ## Build all packages
	go build ./...

.PHONY: check
check: fmt tidy vet lint test vuln build ## Run local quality gate
```

---

## Suggested CI Quality Gate

A sensible first CI quality gate:

```text
1. checkout code
2. setup Go
3. install goimports
4. install golangci-lint
5. install govulncheck
6. run go mod tidy check
7. run gofmt/goimports check
8. run go vet
9. run golangci-lint
10. run go test
11. run govulncheck
12. run go build
```

Later, add:

```text
race detection job
integration test job
container build job
SBOM generation
image vulnerability scanning
Protobuf breaking-change checks
```

---

## Suggested GitHub Actions Outline

```yaml
name: Go Quality

on:
  pull_request:
  push:
    branches:
      - main

jobs:
  go-quality:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: services/catalog-service/go.mod
          cache-dependency-path: services/catalog-service/go.sum

      - name: Install tools
        run: |
          go install golang.org/x/tools/cmd/goimports@latest
          go install golang.org/x/vuln/cmd/govulncheck@latest
          go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

      - name: Check formatting
        working-directory: services/catalog-service
        run: |
          goimports -w .
          git diff --exit-code

      - name: Tidy modules
        working-directory: services/catalog-service
        run: |
          go mod tidy
          git diff --exit-code

      - name: Vet
        working-directory: services/catalog-service
        run: go vet ./...

      - name: Lint
        working-directory: services/catalog-service
        run: golangci-lint run ./...

      - name: Test
        working-directory: services/catalog-service
        run: go test ./...

      - name: Vulnerability check
        working-directory: services/catalog-service
        run: govulncheck ./...

      - name: Build
        working-directory: services/catalog-service
        run: go build ./...
```

---

## Editor Recommendations

Configure your editor to run:

```text
goimports on save
gopls language server
go test shortcuts
lint feedback where available
```

Recommended editor tooling:

```text
gopls
goimports
golangci-lint integration
test explorer support
```

---

## bfstore Practical Standard

For bfstore services, aim for this before opening a PR:

```sh
make check
```

Where `make check` should cover:

```text
formatting
module tidying
vetting
linting
unit tests
vulnerability checks
build verification
```

---

## Kuti Judgement

Use layers, not vibes:

```text
gofmt/goimports  -> code shape
go vet           -> suspicious constructs
go test          -> behaviour
go test -race    -> concurrency safety
golangci-lint    -> broad quality gate
staticcheck      -> deeper static analysis
govulncheck      -> known vulnerability exposure
go build         -> compile verification
```

A tool cannot make bad design good, but it can stop avoidable nonsense reaching main.

That is the boring magic.
