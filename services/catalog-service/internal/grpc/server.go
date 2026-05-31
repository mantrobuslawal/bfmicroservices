package grpc

import (
	"log/slog"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
	"google.golang.org/grpc"

	catalogv1 "github.com/mantrobuslawal/bfstore/gen/go/bfstore/catalog/v1"
)

// NewServer creates the Catalogue Service gRPC server.
func NewServer(catalogService *catalog.Service, logger *slog.Logger) *grpc.Server {
	server := grpc.NewServer()

	handler := NewCatalogHandler(catalogService, logger)

	catalogv1.RegisterCatalogServiceServer(server, handler)

	return server
}
