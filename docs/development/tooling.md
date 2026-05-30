# Tooling

This document defines bfstore tooling expectations and version-management guidance.

It complements:

```text
docs/development/dependencies.md
docs/security/dependency-management.md
```

---

## Purpose

This document explains:

```text
required tools
how versions are pinned
how to install/use tools
Makefile targets
doctor/bootstrap flow
local development expectations
CI tooling expectations
```

---

## Required Core Tools

Core tools may include:

```text
go
docker
docker compose
buf
protoc-gen-go
protoc-gen-go-grpc
git
make
```

Later platform tools may include:

```text
kubectl
helm
terraform
trivy
syft
cosign
```

---

## Version Pinning

Tool versions should be pinned where practical.

Possible files:

```text
mise.toml
.tool-versions
tools.go
go.mod
devcontainer.json
Makefile
CI workflow
```

Avoid:

```text
install latest
assume global tool exists
undocumented local setup
```

---

## tools.go Pattern

Go-based tools can be pinned through a `tools.go` file.

Example:

```go
//go:build tools

package tools

import (
    _ "github.com/bufbuild/buf/cmd/buf"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
```

This lets Go modules track tool versions.

---

## Makefile Targets

Recommended targets:

```text
make doctor
make bootstrap
make lint
make test
make proto-lint
make proto-generate
make proto-breaking
make security-scan
make up
make down
```

---

## make doctor

`make doctor` should verify required tools.

Example checks:

```text
go installed
docker installed
docker compose available
buf available or runnable through pinned command
git available
```

It should fail fast with clear messages.

---

## make bootstrap

`make bootstrap` may:

```text
download Go modules
install pinned local tools
prepare local directories
verify Docker availability
prepare Git hooks if used
```

It should not install unpinned tools silently.

---

## Protobuf Tooling

bfstore Protobuf tooling should use:

```text
buf.yaml
buf.gen.yaml
buf lint
buf generate
buf breaking
```

Tooling must be repeatable in:

```text
developer laptop
CI runner
containerised environment
```

---

## CI Tooling

CI should install or run pinned versions.

Examples:

```text
actions/setup-go with go-version-file
go run github.com/bufbuild/buf/cmd/buf@vX.Y.Z
pinned golangci-lint action/version
pinned security scanner versions
```

CI should not rely on mystery tools already present on the runner.

---

## Dev Container Option

A future dev container may provide:

```text
Go
Docker CLI
Buf
Protobuf plugins
linting tools
security tools
kubectl/helm where needed
```

Possible file:

```text
.devcontainer/devcontainer.json
```

This can improve onboarding and reviewer experience.

---

## Practical Rules

```text
Document required tools.
Pin important tool versions.
Prefer repeatable commands.
Use Makefile targets for common workflows.
Use make doctor to detect missing tools.
Avoid global mystery dependencies.
Keep CI and local workflows aligned.
```

---

## Final Rule

```text
Tooling should make the correct workflow obvious and repeatable.
```
