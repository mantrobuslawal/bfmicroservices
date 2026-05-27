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

	"github.com/acme-ltd/bfstore/services/catalog-service/internal/catalog"
	"github.com/acme-ltd/bfstore/services/catalog-service/internal/config"
	"github.com/acme-ltd/bfstore/services/catalog-service/internal/database"
	cataloggrpc "github.com/acme-ltd/bfstore/services/catalog-service/internal/grpc"
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

	if err := database.Ping(ctx, db); err != nil {
		logger.Error("database ping failed", "error", err)
		os.Exit(1)
	}

	repository := catalog.NewMySQLRepository(db)
	service := catalog.NewService(repository)
	grpcServer := cataloggrpc.NewServer(service, logger)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Error("failed to listen for gRPC", "port", cfg.GRPCPort, "error", err)
		os.Exit(1)
	}

	go func() {
		logger.Info("catalog-service started", "grpc_port", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server stopped unexpectedly", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("shutdown signal received")

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
