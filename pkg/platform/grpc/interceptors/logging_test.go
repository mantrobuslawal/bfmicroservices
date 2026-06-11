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

func TestUnaryLoggingInterceptorReturnsHandlerResponse(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryLoggingInterceptor(logger)

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

func TestUnaryLoggingInterceptorReturnHandlerError(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryLoggingInterceptor(logger)

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

func TestUnaryLoggingInterceptorCallsHandlerOnce(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryLoggingInterceptor(logger)

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

func TestUnaryLoggingInterceptorAcceptsNilLogger(t *testing.T) {
	t.Parallel()

	interceptor := UnaryLoggingInterceptor(nil)

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

func TestUnaryLoggingInterceptorReturnsNonGRPCError(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	interceptor := UnaryLoggingInterceptor(logger)

	wantErr := errors.New("database unavailable")

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
