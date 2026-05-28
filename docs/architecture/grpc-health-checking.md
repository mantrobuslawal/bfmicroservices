# gRPC Health Checking

This document defines how **bfstore** services expose and manage gRPC health checking.

bfstore uses the standard gRPC health checking service rather than custom `Ping`, `IsAlive`, or service-specific health RPCs.

---

## Purpose

The purpose of gRPC health checking is to provide a standard traffic safety signal for:

```text
Kubernetes readiness probes
Kubernetes liveness probes
local development checks
CI smoke tests
load balancers
future client-side health checking
graceful shutdown
```

A service reporting healthy should be able to safely receive the type of traffic being checked.

---

## Standard Health Service

bfstore services should register the standard gRPC health service:

```go
import (
    "google.golang.org/grpc/health"
    healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

healthServer := health.NewServer()
healthpb.RegisterHealthServer(grpcServer, healthServer)
```

Services should not define custom health RPCs such as:

```proto
rpc Ping(PingRequest) returns (PingResponse);
rpc IsAlive(IsAliveRequest) returns (IsAliveResponse);
```

---

## Health Statuses

bfstore services use the standard health statuses:

| Status | Meaning |
|---|---|
| `SERVING` | The service is ready to accept traffic. |
| `NOT_SERVING` | The service should not receive traffic right now. |
| `UNKNOWN` | Health state is unknown. Avoid as normal runtime state. |
| `SERVICE_UNKNOWN` | Requested service name is not known to the health server. |

---

## Whole-server and Service-specific Health

bfstore services should expose both:

```text
""                                  = whole gRPC server health
"bfstore.<service>.v1.<Service>"    = service-specific health
```

Example for `catalog-service`:

```text
""                                  = whole server
bfstore.catalog.v1.CatalogService   = catalogue API readiness
```

The whole-server health is useful for simple liveness checks.

The service-specific health is useful for readiness checks and client-side health checking.

---

## Startup Behaviour

Services must start as:

```text
NOT_SERVING
```

They should only switch to:

```text
SERVING
```

after critical startup checks pass.

Example critical checks for `catalog-service`:

```text
configuration loaded
gRPC handlers registered
MySQL reachable
catalogue schema compatible
service is not shutting down
```

Example:

```go
healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
healthServer.SetServingStatus(catalogServiceName, healthpb.HealthCheckResponse_NOT_SERVING)

if err := dependencies.Ready(ctx); err == nil {
    healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
    healthServer.SetServingStatus(catalogServiceName, healthpb.HealthCheckResponse_SERVING)
}
```

---

## Runtime Dependency Behaviour

Services should update health when critical dependencies fail or recover.

For example, if `catalog-service` loses MySQL:

```text
health becomes NOT_SERVING
Kubernetes removes the pod from ready endpoints
clients stop routing new catalogue traffic
service continues attempting recovery
health returns to SERVING when MySQL recovers
```

Dependency checks should be cheap and deliberate.

Do not make every health check perform expensive business queries.

Preferred pattern:

```text
background dependency monitor updates health status
health RPC returns latest known status
```

---

## Dependency Guidance by Service

### `catalog-service`

`SERVING` requires:

```text
MySQL reachable
catalogue schema available
configuration valid
server not shutting down
```

### `inventory-service`

`SERVING` requires:

```text
MySQL reachable
inventory schema available
reservation logic enabled
server not shutting down
```

### `order-service`

`SERVING` requires:

```text
order database reachable
configuration valid
core orchestration dependencies configured
server not shutting down
```

`order-service` should not necessarily become `NOT_SERVING` just because one downstream service has a brief failure. Checkout should handle downstream errors using domain-specific failure paths.

### `payment-service`

`SERVING` requires:

```text
configuration valid
idempotency store reachable
payment provider adapter available, if applicable
server not shutting down
```

### `notification-service`

`SERVING` requires:

```text
Kafka dependency available, where required
notification templates/config loaded
server not shutting down
```

---

## Health Check Cost Rules

Health checks must be:

```text
fast
safe
cheap
side-effect free
```

Good:

```text
read internal readiness state
ping database with timeout
verify dependency monitor state
verify configuration loaded
```

Bad:

```text
query entire product catalogue
publish Kafka events
perform checkout simulation
call five downstream services per probe
run database migrations
```

Health checks should observe readiness, not create load.

---

## Kubernetes Probe Strategy

bfstore should use native gRPC probes where available.

### Readiness

Use service-specific health for readiness.

Example:

```yaml
readinessProbe:
  grpc:
    port: 50051
    service: bfstore.catalog.v1.CatalogService
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 2
  failureThreshold: 3
```

Readiness answers:

```text
Should this pod receive traffic?
```

### Liveness

Use whole-server health for liveness.

Example:

```yaml
livenessProbe:
  grpc:
    port: 50051
    service: ""
  initialDelaySeconds: 10
  periodSeconds: 20
  timeoutSeconds: 2
  failureThreshold: 3
```

Liveness answers:

```text
Is this process alive enough to keep running?
```

Do not make liveness too dependency-sensitive. A short MySQL outage should usually make a pod not-ready, not immediately killed.

### Startup

Use startup probes later if a service may take longer to initialise.

This can protect slow-starting services from premature liveness failures.

---

## Shutdown Behaviour

Before graceful shutdown, services must mark health as:

```text
NOT_SERVING
```

Then drain existing calls.

Example:

```go
healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
healthServer.Shutdown()

grpcServer.GracefulStop()
```

Recommended shutdown flow:

```text
receive termination signal
mark whole-server and service-specific health NOT_SERVING
notify health watchers
stop accepting new traffic
allow in-flight RPCs to complete
force stop if graceful timeout expires
exit
```

This supports safe Kubernetes rollouts and graceful deployments.

---

## Client-side Health Checking

Client-side health checking may be added later using a gRPC service config.

Example:

```json
{
  "healthCheckConfig": {
    "serviceName": "bfstore.catalog.v1.CatalogService"
  }
}
```

This causes clients to use the health `Watch` RPC and avoid sending calls to unhealthy backends.

Use this later when bfstore has:

```text
multiple service instances
client-side load balancing
service discovery
more realistic Kubernetes networking
```

Initial priority:

```text
server-side health
Kubernetes readiness/liveness
local smoke tests
```

---

## Local Development Checks

Example with `grpcurl`:

```bash
grpcurl -plaintext \
  -d '{"service":"bfstore.catalog.v1.CatalogService"}' \
  localhost:50051 \
  grpc.health.v1.Health/Check
```

Expected response:

```json
{
  "status": "SERVING"
}
```

---

## Testing Guidance

Each service should have tests for health behaviour.

Recommended tests:

```text
starts as NOT_SERVING before dependencies are ready
moves to SERVING after dependency checks pass
moves to NOT_SERVING when critical dependency fails
sets NOT_SERVING during shutdown
returns SERVICE_UNKNOWN for unknown service names where applicable
```

For package-level tests, use the standard health client against an in-memory/bufconn gRPC server where practical.

---

## Practical Rules

```text
Use the standard gRPC health service.
Do not invent custom Ping RPCs.
Start as NOT_SERVING.
Become SERVING only after critical dependencies are ready.
Expose whole-server health and service-specific health.
Use service-specific health for readiness.
Use whole-server health for liveness.
Keep checks cheap and side-effect free.
Update health when dependencies fail or recover.
Mark NOT_SERVING before graceful shutdown.
Use Watch/client-side health later when topology needs it.
```

---

## Final Rule

```text
Health checks are traffic safety signals, not decorative endpoints.
```

For bfstore, health checking should make deployments, dependency failures, and service shutdowns safer and easier to operate.
