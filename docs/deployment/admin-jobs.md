# Admin Jobs

This document defines deployment patterns for bfstore admin jobs.

## Purpose

This document explains:

```text
Kubernetes Job pattern
Docker Compose local command pattern
CI/CD migration flow
logging and exit code expectations
```

## Core Rule

```text
One-off operational tasks should run as controlled jobs.
```

## Local Docker Compose Pattern

```bash
docker compose run --rm catalog-service bfstore-admin migrate up --service=catalog
docker compose run --rm catalog-service bfstore-admin seed catalog --dataset=borough-demo
docker compose run --rm order-service bfstore-admin outbox list --status=FAILED
```

Use the same service image and config style as normal local services.

## Kubernetes Job Pattern

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: order-migration
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: migrate
          image: ghcr.io/mantrobuslawal/bfstore/order-service:abc123
          command: ["bfstore-admin", "migrate", "up", "--service=order"]
          envFrom:
            - secretRef:
                name: order-service-secrets
            - configMapRef:
                name: order-service-config
```

## CI/CD Migration Flow

```text
build image
run tests
run buf checks
deploy migration job
wait for migration success
deploy service
run smoke tests
```

Migrations should not run automatically inside every service replica.

## Logging

Admin jobs should log to stdout/stderr using structured JSON.

Useful fields:

```text
task
service
environment
release
started_at
duration_ms
records_scanned
records_changed
success
failure_reason
```

## Exit Codes

```text
0:
  success

non-zero:
  task failed or partially failed
```

Jobs should fail loudly when they cannot safely complete.

## Practical Rules

```text
Use Jobs for staging/prod-style admin work.
Use docker compose run for local admin work.
Tie jobs to known release images.
Use the same secrets/config model.
Wait for migration jobs before rollout.
Keep logs and summaries.
```

## Final Rule

```text
Admin jobs should be repeatable enough for automation and clear enough for humans.
```
