# Admin Processes

This document defines how **bfstore** runs one-off administrative and maintenance tasks.

## Purpose

This document explains:

```text
what admin processes are
how they differ from long-running processes
same release/config/dependencies rule
examples for migrations, seed data, replay, and repair
```

## Core Rule

```text
Admin work is production work.
Treat it with production discipline.
```

Admin tasks should be explicit, repeatable, logged, and tied to a known release.

## Admin Processes vs Long-running Processes

Long-running processes include:

```text
api-gateway
catalog-service
basket-service
inventory-service
order-service
payment-service
shipping-service
notification-worker
```

Admin processes include:

```text
migration jobs
seed jobs
outbox replay jobs
notification retry jobs
inventory reconciliation jobs
search reindex jobs later
data repair jobs
```

Admin processes usually run once, do a defined task, then exit.

## Same Release / Config / Dependencies

Admin processes should use:

```text
same repo
same release image
same config/secrets model
same dependency versions
same logging conventions
same telemetry conventions where useful
```

## Example Admin Tasks

```text
bfstore-admin migrate up --service=order
bfstore-admin migrate status --service=order
bfstore-admin seed catalog --dataset=borough-demo
bfstore-admin outbox replay --service=order
bfstore-admin notifications retry --status=FAILED
bfstore-admin inventory reconcile-reservations --dry-run
bfstore-admin search reindex --scope=products --mode=full
```

## Safety Principles

Admin tasks should support:

```text
dry-run for risky changes
bounded scope
clear flags
structured logs
clear exit codes
idempotency where practical
summary output
documented runbooks
```

## Practical Rules

```text
Keep admin code in the repo.
Run tasks from known images/releases.
Use env vars/secrets for config.
Avoid manual SQL for repeatable tasks.
Use dedicated jobs for migrations.
Use dry-run for data repairs.
Log structured summaries to stdout.
```

## Final Rule

```text
Admin processes maintain the business without turning production into a manual terminal session.
```
