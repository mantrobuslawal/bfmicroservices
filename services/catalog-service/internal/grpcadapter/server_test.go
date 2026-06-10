package grpcadapter

import (
	"log/slog"
	"testing"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/catalog"
)

func TestNewServerReturnsGRPCServer(t *testing.T) {
	t.Parallel()

	server := NewServer(catalog.NewService(fakeCatalogRepository{}), slog.Default())
	if server == nil {
		t.Fatal("NewServer() = nil, want server")
	}

	server.Stop()
}
