package grpcadapter

import (
	"log/slog"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
	platforminterceptors "github.com/mantrobuslawal/bfstore/pkg/platform/grpc/interceptors"
)

// NewServer creates the Catalogue Service gRPC server.
func NewServer(catalogService *catalog.Service, logger *slog.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			platforminterceptors.UnaryRecoveryInterceptor(logger),
			platforminterceptors.UnaryLoggingInterceptor(logger),
		),
	)

	handler := NewCatalogHandler(catalogService, logger)

	catalogv1.RegisterCatalogServiceServer(server, handler)

	return server
}
