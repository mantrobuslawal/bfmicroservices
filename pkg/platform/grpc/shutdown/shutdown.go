package shutdown

import (
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

// HealthManager is the minimal health-check integration needed during shutdown.
type HealthManager interface {
	MarkNotServing()
	Shutdown()
}

// Graceful attempts to gracefully stop a gRPC server, falling back to Stop after timeout.
//
// Shutdown order:
//   - mark health as NOT_SERVING, if a health manager is provided
//   - notify health watchers, if supported by the health manager
//   - call GracefulStop in a goroutine
//   - call Stop if the timeout expires
func Graceful(
	logger *slog.Logger,
	grpcServer *grpc.Server,
	healthManager HealthManager,
	timeout time.Duration,
) {
	if logger == nil {
		logger = slog.Default()
	}

	logger.Info("starting grpc graceful shutdown", "timeout", timeout.String())

	if healthManager != nil {
		logger.Info("marking grpc health not serving")
		healthManager.MarkNotServing()
		healthManager.Shutdown()
	}

	done := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("grpc server stopped gracefully")

	case <-time.After(timeout):
		logger.Warn("grpc graceful shutdown timed out; forcing stop")
		grpcServer.Stop()
	}
}
