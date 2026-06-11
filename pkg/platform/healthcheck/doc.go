// Package healthcheck provides shared gRPC health-check wiring for bfstore
// services.
//
// The package wraps the standard gRPC health service and gives services a small,
// consistent way to register health status, mark services as serving, mark
// services as not serving, and coordinate health state during startup and
// graceful shutdown.
//
// This package owns platform-level health plumbing. It should not contain
// service-specific readiness checks, database pings, Kafka checks, repository
// calls, or domain logic. Those responsibilities belong inside each service's
// internal health package.
//
// Typical responsibilities include:
//
//   - registering the standard grpc.health.v1.Health service;
//   - tracking whole-server and service-specific health status;
//   - marking services as SERVING after startup readiness checks pass;
//   - marking services as NOT_SERVING during graceful shutdown;
//   - providing a consistent health-check integration pattern across bfstore
//     services.
//
// Service-specific packages decide what readiness means. For example, the
// catalog service may check whether its MySQL database is reachable, then use
// this package to expose the resulting status through the standard gRPC health
// API.
//
// In summary this package offers: simple status management,
// clear service names, predictable shutdown behaviour, and no hidden service
// orchestration.
package healthcheck
