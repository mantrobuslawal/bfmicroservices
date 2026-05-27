package grpc

import (
	"log/slog"

	"github.com/acme-ltd/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc"
)

// NewServer creates the Catalogue Service gRPC server.
//
// Protobuf-generated service registration should be added once generated code
// is present in the repository.
func NewServer(catalogService *catalog.Service, logger *slog.Logger) *grpc.Server {
	server := grpc.NewServer()

	// TODO:
	// Register generated Catalogue Service handler after running buf generate.
	//
	// Example:
	// catalogv1.RegisterCatalogServiceServer(
	//   server,
	//   NewCatalogHandler(catalogService, logger),
	// )

	_ = catalogService
	_ = logger

	return server
}
