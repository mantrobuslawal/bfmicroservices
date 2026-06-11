package interceptors

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryRecoveryInterceptorReturnsHandlerResponse(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryRecoveryInterceptor(logger)

	wantResponse := "catalog response"

	response, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListProducts",
		},
		func(ctx context.Context, req any) (any, error) {
			return wantResponse, nil
		},
	)

	if err != nil {
		t.Fatalf("interceptor returned error = %v, want nil", err)
	}

	if response != wantResponse {
		t.Fatalf("response = %v, want %v", response, wantResponse)
	}
}

func TestUnaryRecoveryInterceptorReturnsHandlerError(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryRecoveryInterceptor(logger)

	wantErr := status.Error(codes.InvalidArgument, "invalid page size")

	response, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListProducts",
		},
		func(ctx context.Context, req any) (any, error) {
			return nil, wantErr
		},
	)

	if response != nil {
		t.Fatalf("response = %v, want nil", response)
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestUnaryRecoveryInterceptorRecoversPanic(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryRecoveryInterceptor(logger)

	response, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/GetProduct",
		},
		func(ctx context.Context, req any) (any, error) {
			panic("catalog mapper exploded")
		},
	)

	if response != nil {
		t.Fatalf("response = %v, want nil", response)
	}

	if err == nil {
		t.Fatal("error = nil, want non-nil")
	}

	if gotCode := status.Code(err); gotCode != codes.Internal {
		t.Fatalf("status.Code(error) = %v, want %v", gotCode, codes.Internal)
	}
}

func TestUnaryRecoveryInterceptorRecoversNonStringPanic(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryRecoveryInterceptor(logger)

	_, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/GetProduct",
		},
		func(ctx context.Context, req any) (any, error) {
			panic(struct {
				Message string
			}{
				Message: "structured panic value",
			})
		},
	)

	if err == nil {
		t.Fatal("error = nil, want non-nil")
	}

	if gotCode := status.Code(err); gotCode != codes.Internal {
		t.Fatalf("status.Code(error) = %v, want %v", gotCode, codes.Internal)
	}
}

func TestUnaryRecoveryInterceptorAcceptsNilLogger(t *testing.T) {
	t.Parallel()

	interceptor := UnaryRecoveryInterceptor(nil)

	response, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListCategories",
		},
		func(ctx context.Context, req any) (any, error) {
			return "response", nil
		},
	)

	if err != nil {
		t.Fatalf("interceptor returned error = %v, want nil", err)
	}

	if response != "response" {
		t.Fatalf("response = %v, want %v", response, "response")
	}
}

func TestUnaryRecoveryInterceptorWithNilLoggerRecoversPanic(t *testing.T) {
	t.Parallel()

	interceptor := UnaryRecoveryInterceptor(nil)

	response, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListCategories",
		},
		func(ctx context.Context, req any) (any, error) {
			panic("nil logger recovery path")
		},
	)

	if response != nil {
		t.Fatalf("response = %v, want nil", response)
	}

	if err == nil {
		t.Fatal("error = nil, want non-nil")
	}

	if gotCode := status.Code(err); gotCode != codes.Internal {
		t.Fatalf("status.Code(error) = %v, want %v", gotCode, codes.Internal)
	}
}

func TestUnaryRecoveryInterceptorCallsHandlerOnce(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryRecoveryInterceptor(logger)

	var calls int

	_, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/GetProduct",
		},
		func(ctx context.Context, req any) (any, error) {
			calls++
			return "response", nil
		},
	)

	if err != nil {
		t.Fatalf("interceptor returned error = %v, want nil", err)
	}

	if calls != 1 {
		t.Fatalf("handler calls = %d, want 1", calls)
	}
}

func TestUnaryRecoveryInterceptorCallsPanickingHandlerOnce(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryRecoveryInterceptor(logger)

	var calls int

	_, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/GetProduct",
		},
		func(ctx context.Context, req any) (any, error) {
			calls++
			panic("only once")
		},
	)

	if err == nil {
		t.Fatal("error = nil, want non-nil")
	}

	if calls != 1 {
		t.Fatalf("handler calls = %d, want 1", calls)
	}
}
