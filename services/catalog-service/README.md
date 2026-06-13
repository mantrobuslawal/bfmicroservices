# Catalog Service

The `catalog-service` owns product catalog data for bfstore.

## Observability

Current local telemetry flow:

```text
catalog-service
  -> OpenTelemetry Collector
  -> Jaeger for traces
  -> Prometheus for metrics
  -> Grafana for dashboards
```

The service currently emits:

```text
gRPC request spans
database/sql spans
database connection pool metrics
structured request logs
correlation IDs
```

## Grafana dashboard

Dashboard:

```text
Catalog DB Pool Overview
```

Dashboard file:

```text
deployments/local/grafana/dashboards/catalog-db-pool-overview.json
```

Dashboard provider:

```text
deployments/local/grafana/provisioning/dashboards/dashboards.yml
```

Prometheus data source:

```text
deployments/local/grafana/provisioning/datasources/prometheus.yml
```

## Running with telemetry

```bash
make observability-up
```

```bash
cd services/catalog-service

TELEMETRY_ENABLED=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_EXPORTER_OTLP_INSECURE=true \
GRPC_REFLECTION_ENABLED=true \
go run ./cmd/catalog-service
```

Send a request:

```bash
grpcurl -plaintext \
  -H 'x-correlation-id: local-dev-grafana-123' \
  -d '{"page":{"page_size":5}}' \
  localhost:50051 \
  bfstore.catalog.v1.CatalogService/ListProducts
```

## UIs

```text
Jaeger:     http://localhost:16686
Prometheus: http://localhost:9090
Grafana:    http://localhost:3000
```

## Current runtime foundation

The service currently demonstrates:

```text
gRPC API serving
standard gRPC health checks
gRPC reflection for local development
structured request logging
correlation ID propagation
panic recovery
graceful shutdown
database readiness checks
OpenTelemetry bootstrap
gRPC server tracing instrumentation
database/sql tracing through otelsql
database pool metrics through dbmetrics
OpenTelemetry Collector integration
Jaeger trace visualisation
Prometheus metric querying
Grafana dashboard provisioning
```

## Practical rule

```text
Traces explain request paths.
Metrics explain resource health over time.
Dashboards make the signal easy to inspect.
```

Keep it boring where production matters.
