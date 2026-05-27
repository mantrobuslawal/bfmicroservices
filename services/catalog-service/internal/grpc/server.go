package grpc

import (
	"log/slog"

	"github.com/acme-ltd/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc"

	// TODO:
	// Uncomment after generated Protobuf code is committed.
	//
	// catalogv1 "github.com/acme-ltd/bfstore/gen/go/acme/catalog/v1"
)

// NewServer creates the Catalogue Service gRPC server.
func NewServer(catalogService *catalog.Service, logger *slog.Logger) *grpc.Server {
	server := grpc.NewServer()

	handler := NewCatalogHandler(catalogService, logger)

	// TODO:
	// Register the generated Catalogue Service handler after running buf generate.
	//
	// catalogv1.RegisterCatalogServiceServer(server, handler)

	_ = handler

	return server
}
