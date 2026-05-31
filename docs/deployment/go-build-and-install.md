# Go Build and Install

This document defines bfstore's Go build, install, and run conventions.

## Core Rule

```text
Docker should build the same Go package CI builds.
```

## Build Everything

```bash
go build ./...
```

This should pass before code is considered healthy.

## Build a Service

```bash
go build -o bin/catalog-service ./services/catalog/cmd/catalog-service
```

## Run a Service Locally

```bash
go run ./services/catalog/cmd/catalog-service
```

## Install Developer Tools

Use `go install` for developer tools and CLIs.

```bash
go install ./cmd/bfstore-admin
```

Then:

```bash
bfstore-admin migrate status --service=catalog
```

## Docker Build Pattern

```dockerfile
FROM golang:1.23 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/catalog-service ./services/catalog/cmd/catalog-service

FROM gcr.io/distroless/base-debian12

COPY --from=build /out/catalog-service /catalog-service

ENTRYPOINT ["/catalog-service"]
```

## CI Build

CI should run:

```bash
go build ./...
```

## Make Targets

```makefile
.PHONY: go-build
go-build:
	go build ./...

.PHONY: catalog-build
catalog-build:
	go build -o bin/catalog-service ./services/catalog/cmd/catalog-service

.PHONY: catalog-run
catalog-run:
	go run ./services/catalog/cmd/catalog-service
```

## Practical Rules

```text
Use go build ./... as a baseline health check.
Use explicit command package paths for service binaries.
Use go install for local developer CLIs.
Use go run for local development.
Use Docker builds that target the same command packages as CI.
Avoid hidden build steps.
```

## Final Rule

```text
Builds should be repeatable locally, in CI, and in Docker.
```
