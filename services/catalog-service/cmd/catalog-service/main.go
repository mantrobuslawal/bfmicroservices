package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/config"
	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/database"
	cataloggrpc "github.com/mantrobuslawal/bfstore/services/catalog-service/internal/grpcadapter"
	cataloghealth "github.com/mantrobuslawal/bfstore/services/catalog-service/internal/health"

	grpchealth "google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer closeDatabase(logger, db)

	logger.Info("database connection opened")

	repository := catalog.NewMySQLRepository(db)
	service := catalog.NewService(repository)
	grpcServer := cataloggrpc.NewServer(service, logger)

	if cfg.EnableGRPCReflection {
		reflection.Register(grpcServer)
		logger.Info("grpc reflection enabled")
	}

	healthchecker := cataloghealth.NewChecker(db)

	if err := healthchecker.Ready(ctx); err != nil {
		logger.Error("service is not ready", "error", err)
		os.Exit(1)
	}

	logger.Info("database readiness check passed")

	healthServer := grpchealth.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthServer)

	const catalogServiceName = "bfstore.catalog.v1.CatalogService"

	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(catalogServiceName, healthv1.HealthCheckResponse_SERVING)

	logger.Info("grpc health service is registered")

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Error("failed to listen for gRPC", "port", cfg.GRPCPort, "error", err)
		os.Exit(1)
	}

	serverErr := make(chan error, 1)

	go func() {
		logger.Info("catalog-service started", "grpc_port", cfg.GRPCPort)
		serverErr <- grpcServer.Serve(listener)
	}()

	go func() { monitorReadiness(ctx, logger, healthchecker, healthServer, catalogServiceName) }()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")

	case err := <-serverErr:
		if err != nil {
			logger.Error("gRPC server failed", "error", err)
		}
		stop()
	}

	logger.Info("marking service not serving")
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_NOT_SERVING)
	healthServer.SetServingStatus(catalogServiceName, healthv1.HealthCheckResponse_NOT_SERVING)
	healthServer.Shutdown()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("catalog-service stopped gracefully")
	case <-shutdownCtx.Done():
		logger.Warn("graceful shutdown timed out; forcing stop")
		grpcServer.Stop()
	}
}

func closeDatabase(logger *slog.Logger, db *sql.DB) {
	if err := db.Close(); err != nil {
		logger.Error("failed to close database", "error", err)
	}
}

func monitorReadiness(
	ctx context.Context,
	logger *slog.Logger,
	checker *cataloghealth.Checker,
	healthServer *grpchealth.Server,
	serviceName string,
) {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			healthServer.SetServingStatus("", healthv1.HealthCheckResponse_NOT_SERVING)
			healthServer.SetServingStatus(serviceName, healthv1.HealthCheckResponse_NOT_SERVING)
			return

		case <-ticker.C:
			if err := checker.Ready(ctx); err != nil {
				logger.Warn("readiness check failed", "error", err)
				healthServer.SetServingStatus("", healthv1.HealthCheckResponse_NOT_SERVING)
				healthServer.SetServingStatus(serviceName, healthv1.HealthCheckResponse_NOT_SERVING)
				continue
			}

			healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
			healthServer.SetServingStatus(serviceName, healthv1.HealthCheckResponse_SERVING)
		}
	}

}
