# ADR-0006: Use Buf for Protobuf Tooling

## Status

Accepted

## Date

2026-05-26

## Context

bfstore uses Protobuf for gRPC service contracts and Kafka event payload contracts.

To keep contracts consistent, generated correctly, and safe to evolve, the project needs standard tooling for:

```text
protobuf linting
breaking-change detection
code generation
consistent style enforcement
CI quality gates
```

## Decision

bfstore will use Buf for Protobuf linting, breaking-change detection, and code generation.

## Drivers

This decision supports:

```text
contract-first development
consistent protobuf style
CI enforcement
safe API evolution
generated Go clients and servers
client-facing professionalism
```

## Alternatives Considered

### Option 1: Raw protoc Scripts

Benefits:

```text
simple in small projects
direct use of standard compiler
```

Costs:

```text
more custom scripting
less built-in linting
harder breaking-change governance
more inconsistent developer experience
```

### Option 2: Buf

Benefits:

```text
standardised protobuf tooling
linting
breaking-change detection
generation config
good CI integration
professional contract governance
```

Costs:

```text
additional tool to learn
requires buf.yaml and buf.gen.yaml management
```

## Consequences

### Positive

```text
protobuf contracts can be validated automatically
breaking changes can fail CI
generated code paths can be standardised
API style can remain consistent across services
```

### Negative

```text
developers must install or use Buf through containers/scripts
CI must be configured for Buf
generated output must be managed carefully
```

## Implementation Notes

Expected files:

```text
buf.yaml
buf.gen.yaml
```

Expected commands:

```sh
buf lint
buf breaking
buf generate
```

Expected Makefile targets:

```sh
make proto-lint
make proto-breaking
make proto-generate
make proto
```

## CI Expectations

CI should fail when:

```text
protobuf lint rules fail
breaking changes are introduced without a new version
code generation fails
generated files are out of date if committed
```

## Risks

| Risk | Mitigation |
|---|---|
| Developers bypass Buf | Make Buf part of Makefile and CI |
| Breaking-change baseline unclear | Document baseline strategy |
| Generated code becomes inconsistent | Use single buf.gen.yaml |
| Tooling slows early development | Start with simple rules and mature over time |

## Review Triggers

Revisit this decision if:

```text
Buf does not support required workflow
another schema registry/tooling approach becomes mandatory
contract generation becomes too complex for repo structure
```

## Related Documents

```text
docs/api/protobuf-style-guide.md
docs/api/versioning.md
docs/api/grpc-overview.md
docs/events/event-envelope.md
```
