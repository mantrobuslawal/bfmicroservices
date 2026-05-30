# Admin CLI

This document defines guidance for the bfstore admin command-line interface.

## Purpose

This document explains:

```text
bfstore-admin command structure
safe command design
dry-run and execute flags
local usage examples
```

## Core Rule

```text
Admin commands should be explicit, safe, and boring.
```

## Suggested Command Structure

```text
bfstore-admin migrate up --service=order
bfstore-admin migrate status --service=order
bfstore-admin seed catalog --dataset=borough-demo
bfstore-admin outbox list --service=order --status=FAILED
bfstore-admin outbox replay --service=order --since=2026-05-30T00:00:00Z
bfstore-admin notifications retry --status=FAILED --limit=100
bfstore-admin inventory reconcile-reservations --dry-run
bfstore-admin search reindex --scope=products --mode=full
```

## Dry-run and Execute

Risky commands should support dry-run first.

```bash
bfstore-admin inventory reconcile-reservations --dry-run
bfstore-admin inventory reconcile-reservations --execute
```

Avoid commands where the default behaviour is destructive.

## Safe Flags

Useful safety flags:

```text
--dry-run
--execute
--limit
--since
--until
--service
--dataset
--reason
--request-id
```

High-risk commands should require explicit scope and intent.

## Output

Admin commands should emit structured summaries.

```json
{
  "event": "admin_task_completed",
  "task": "catalog_seed",
  "dataset": "borough-demo",
  "inserted": 6,
  "updated": 0,
  "skipped": 0,
  "duration_ms": 842
}
```

## Local Usage

```bash
docker compose run --rm catalog-service bfstore-admin seed catalog --dataset=borough-demo
docker compose run --rm order-service bfstore-admin outbox list --status=FAILED
```

## Practical Rules

```text
Use explicit subcommands.
Use dry-run for risky tasks.
Require scope for bulk changes.
Log structured summaries.
Return clear exit codes.
Avoid hidden defaults.
Keep commands documented.
```

## Final Rule

```text
A good admin CLI makes dangerous work harder to do accidentally.
```
