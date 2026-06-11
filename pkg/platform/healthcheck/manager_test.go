package healthcheck

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

const testCatalogServiceName = "bfstore.catalog.v1.CatalogService"

func TestNewManagerRegistersWholeServerAsNotServing(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)

	response, err := manager.Server().Check(context.Background(), &healthv1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check() error = %v, want nil", err)
	}

	if response.GetStatus() != healthv1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("status = %v, want %v", response.GetStatus(), healthv1.HealthCheckResponse_NOT_SERVING)
	}
}

func TestRegisterServiceMarksServiceNotServing(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)
	manager.RegisterService(testCatalogServiceName)

	response, err := manager.Server().Check(context.Background(), &healthv1.HealthCheckRequest{
		Service: testCatalogServiceName,
	})
	if err != nil {
		t.Fatalf("Check() error = %v, want nil", err)
	}

	if response.GetStatus() != healthv1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("status = %v, want %v", response.GetStatus(), healthv1.HealthCheckResponse_NOT_SERVING)
	}
}

func TestMarkServingMarksWholeServerAndRegisteredServicesServing(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)
	manager.RegisterService(testCatalogServiceName)

	manager.MarkServing()

	wholeServerResponse, err := manager.Server().Check(context.Background(), &healthv1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("whole server Check() error = %v, want nil", err)
	}

	if wholeServerResponse.GetStatus() != healthv1.HealthCheckResponse_SERVING {
		t.Fatalf("whole server status = %v, want %v", wholeServerResponse.GetStatus(), healthv1.HealthCheckResponse_SERVING)
	}

	serviceResponse, err := manager.Server().Check(context.Background(), &healthv1.HealthCheckRequest{
		Service: testCatalogServiceName,
	})
	if err != nil {
		t.Fatalf("service Check() error = %v, want nil", err)
	}

	if serviceResponse.GetStatus() != healthv1.HealthCheckResponse_SERVING {
		t.Fatalf("service status = %v, want %v", serviceResponse.GetStatus(), healthv1.HealthCheckResponse_SERVING)
	}
}

func TestMarkNotServingMarksWholeServerAndRegisteredServicesNotServing(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)
	manager.RegisterService(testCatalogServiceName)
	manager.MarkServing()

	manager.MarkNotServing()

	wholeServerResponse, err := manager.Server().Check(context.Background(), &healthv1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("whole server Check() error = %v, want nil", err)
	}

	if wholeServerResponse.GetStatus() != healthv1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("whole server status = %v, want %v", wholeServerResponse.GetStatus(), healthv1.HealthCheckResponse_NOT_SERVING)
	}

	serviceResponse, err := manager.Server().Check(context.Background(), &healthv1.HealthCheckRequest{
		Service: testCatalogServiceName,
	})
	if err != nil {
		t.Fatalf("service Check() error = %v, want nil", err)
	}

	if serviceResponse.GetStatus() != healthv1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("service status = %v, want %v", serviceResponse.GetStatus(), healthv1.HealthCheckResponse_NOT_SERVING)
	}
}

func TestRegisterServiceIsIdempotent(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)

	manager.RegisterService(testCatalogServiceName)
	manager.RegisterService(testCatalogServiceName)
	manager.RegisterService(testCatalogServiceName)

	services := manager.Services()

	if len(services) != 1 {
		t.Fatalf("len(Services()) = %d, want 1", len(services))
	}

	if services[0] != testCatalogServiceName {
		t.Fatalf("Services()[0] = %q, want %q", services[0], testCatalogServiceName)
	}
}

func TestServicesReturnsSortedServiceNames(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)

	manager.RegisterService("bfstore.shipping.v1.ShippingService")
	manager.RegisterService("bfstore.catalog.v1.CatalogService")
	manager.RegisterService("bfstore.basket.v1.BasketService")

	services := manager.Services()

	want := []string{
		"bfstore.basket.v1.BasketService",
		"bfstore.catalog.v1.CatalogService",
		"bfstore.shipping.v1.ShippingService",
	}

	if len(services) != len(want) {
		t.Fatalf("len(Services()) = %d, want %d", len(services), len(want))
	}

	for i := range want {
		if services[i] != want[i] {
			t.Fatalf("Services()[%d] = %q, want %q", i, services[i], want[i])
		}
	}
}

func TestRegisterServiceIgnoresEmptyServiceName(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)

	manager.RegisterService("")

	if len(manager.Services()) != 0 {
		t.Fatalf("len(Services()) = %d, want 0", len(manager.Services()))
	}
}

func TestShutdownMarksRegisteredServicesNotServing(t *testing.T) {
	t.Parallel()

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	manager := NewManager(grpcServer)
	manager.RegisterService(testCatalogServiceName)
	manager.MarkServing()

	manager.Shutdown()

	response, err := manager.Server().Check(context.Background(), &healthv1.HealthCheckRequest{
		Service: testCatalogServiceName,
	})
	if err != nil {
		t.Fatalf("Check() error = %v, want nil", err)
	}

	if response.GetStatus() != healthv1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("status = %v, want %v", response.GetStatus(), healthv1.HealthCheckResponse_NOT_SERVING)
	}
}

func TestNilManagerMethodsDoNotPanic(t *testing.T) {
	t.Parallel()

	var manager *Manager

	manager.RegisterService(testCatalogServiceName)
	manager.MarkServing()
	manager.MarkNotServing()
	manager.Shutdown()

	if manager.Services() != nil {
		t.Fatalf("Services() = %#v, want nil", manager.Services())
	}

	if manager.Server() != nil {
		t.Fatalf("Server() = %#v, want nil", manager.Server())
	}
}
