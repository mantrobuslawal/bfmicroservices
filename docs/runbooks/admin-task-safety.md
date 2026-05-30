# Admin Task Safety

This runbook defines safety expectations for bfstore admin tasks.

## Purpose

This document explains:

```text
risk levels
approval expectations
dry-run policy
rollback considerations
audit/logging requirements
```

## Core Rule

```text
The more business state a task can change, the more ceremony it deserves.
```

## Risk Levels

### Low risk

Examples:

```text
read-only inspection
list failed outbox records
show migration status
validate config
```

Expected controls:

```text
structured logs
clear command
safe read-only access
```

### Medium risk

Examples:

```text
seed demo data
retry notifications
replay outbox events
refresh search projection
```

Expected controls:

```text
dry-run where practical
bounded scope
idempotency
summary output
documented command
```

### High risk

Examples:

```text
payment repair
inventory correction
bulk order status changes
customer data modification
destructive migrations
```

Expected controls:

```text
dry-run
approval
limit/scope
audit log
backup/rollback thought
extra review
clear reason
```

## Dry-run Policy

Risky commands should support dry-run.

Dry-run output should include:

```text
records that would be changed
side effects that would be triggered
estimated counts
warnings
required execute command
```

## Rollback Considerations

Before running high-risk tasks, document:

```text
what will change
how to verify success
how to detect partial failure
how to recover or rollback
who approved it
where logs will be stored
```

## Audit and Logging

Admin tasks should log:

```text
task name
release version
environment
operator/request identifier where available
parameters
dry-run or execute mode
records scanned/changed
duration
success/failure
```

Never log secrets or sensitive raw payloads.

## Practical Rules

```text
Prefer read-only inspection first.
Use dry-run before execute.
Limit scope.
Avoid broad destructive commands.
Use known release images.
Keep logs.
Document runbooks.
Treat payments, inventory, and customer data as high-risk.
```

## Final Rule

```text
Admin work should leave evidence, not mysteries.
```
