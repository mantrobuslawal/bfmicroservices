package healthcheck


import (
	"sort"
	"sync"

	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

// Manager registers and manages the standard gRPC health service for a bfstore
// service process.
//
// Manager owns platform-level health status plumbing only. It does not decide
// whether a service is ready. Service-specific readiness checks should run in
// service-local packages, then call MarkServing or MarkNotServing based on the
// result.
type Manager struct {
	mu       sync.RWMutex
	server   *grpchealth.Server
	services map[string]struct{}
}

// NewManager registers a standard gRPC health server on grpcServer and returns a
// Manager for updating whole-server and service-specific health status.
//
// The whole-server health status is represented by the empty service name, as
// defined by the standard gRPC health-checking protocol. New managers start as
// NOT_SERVING so services must explicitly mark themselves serving after startup
// readiness checks pass.
func NewManager(grpcServer *grpc.Server) *Manager {
	healthServer := grpchealth.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthServer)

	manager := &Manager{
		server:   healthServer,
		services: make(map[string]struct{}),
	}

	manager.MarkNotServing()

	return manager
}

// RegisterService registers a service name with the health manager and marks it
// NOT_SERVING by default.
//
// RegisterService is safe to call multiple times with the same service name.
// Empty service names are ignored because the empty service name is reserved for
// whole-server health status.
func (m *Manager) RegisterService(serviceName string) {
	if m == nil || serviceName == "" {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.services[serviceName] = struct{}{}
	m.server.SetServingStatus(serviceName, healthv1.HealthCheckResponse_NOT_SERVING)
}

// MarkServing marks the whole server and all registered services as SERVING.
func (m *Manager) MarkServing() {
	if m == nil {
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	m.server.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)

	for serviceName := range m.services {
		m.server.SetServingStatus(serviceName, healthv1.HealthCheckResponse_SERVING)
	}
}

// MarkNotServing marks the whole server and all registered services as
// NOT_SERVING.
func (m *Manager) MarkNotServing() {
	if m == nil {
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	m.server.SetServingStatus("", healthv1.HealthCheckResponse_NOT_SERVING)

	for serviceName := range m.services {
		m.server.SetServingStatus(serviceName, healthv1.HealthCheckResponse_NOT_SERVING)
	}
}

// Shutdown marks every service known to the health server as NOT_SERVING and
// prevents future status updates from changing serving state.
//
// It should be called during graceful shutdown before the gRPC server stops
// accepting traffic.
func (m *Manager) Shutdown() {
	if m == nil {
		return
	}

	m.server.Shutdown()
}

// Services returns the registered service names in sorted order.
//
// It is mainly useful for tests, diagnostics, and startup logging.
func (m *Manager) Services() []string {
	if m == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	services := make([]string, 0, len(m.services))
	for serviceName := range m.services {
		services = append(services, serviceName)
	}

	sort.Strings(services)

	return services
}

// Server returns the underlying standard gRPC health server.
//
// Most service code should not need this. It is exposed for advanced use cases
// and tests that need to call the standard health API directly.
func (m *Manager) Server() healthv1.HealthServer {
	if m == nil {
		return nil
	}

	return m.server
}
