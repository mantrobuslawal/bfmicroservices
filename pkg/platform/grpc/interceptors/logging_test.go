package interceptors

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
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

func TestUnaryLoggingInterceptorReturnsHandlerError(t *testing.T) {
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

func TestUnaryLoggingInterceptorLogsCorrelationIDWhenPresent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	logger := slog.New(slog.NewTextHandler(&buf, nil))
	interceptor := UnaryLoggingInterceptor(logger)

	ctx := ContextWithCorrelationID(context.Background(), "local-dev-123")

	_, err := interceptor(
		ctx,
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListProducts",
		},
		func(ctx context.Context, req any) (any, error) {
			return "response", nil
		},
	)

	if err != nil {
		t.Fatalf("interceptor returned error = %v, want nil", err)
	}

	got := buf.String()

	if !strings.Contains(got, "correlation_id=local-dev-123") {
		t.Fatalf("log output = %q, want correlation_id", got)
	}

	if !strings.Contains(got, "grpc.method=/bfstore.catalog.v1.CatalogService/ListProducts") {
		t.Fatalf("log output = %q, want grpc.method", got)
	}

	if !strings.Contains(got, "grpc.code=OK") {
		t.Fatalf("log output = %q, want grpc.code=OK", got)
	}
}

func TestUnaryLoggingInterceptorDoesNotLogCorrelationIDWhenMissing(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	logger := slog.New(slog.NewTextHandler(&buf, nil))
	interceptor := UnaryLoggingInterceptor(logger)

	_, err := interceptor(
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

	got := buf.String()

	if strings.Contains(got, "correlation_id=") {
		t.Fatalf("log output = %q, did not expect correlation_id", got)
	}

	if !strings.Contains(got, "grpc.method=/bfstore.catalog.v1.CatalogService/ListCategories") {
		t.Fatalf("log output = %q, want grpc.method", got)
	}

	if !strings.Contains(got, "grpc.code=OK") {
		t.Fatalf("log output = %q, want grpc.code=OK", got)
	}
}

func TestUnaryLoggingInterceptorLogsCorrelationIDOnError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	logger := slog.New(slog.NewTextHandler(&buf, nil))
	interceptor := UnaryLoggingInterceptor(logger)

	ctx := ContextWithCorrelationID(context.Background(), "local-dev-error-123")
	wantErr := status.Error(codes.InvalidArgument, "invalid page size")

	response, err := interceptor(
		ctx,
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

	got := buf.String()

	if !strings.Contains(got, "correlation_id=local-dev-error-123") {
		t.Fatalf("log output = %q, want correlation_id", got)
	}

	if !strings.Contains(got, "grpc.code=InvalidArgument") {
		t.Fatalf("log output = %q, want grpc.code=InvalidArgument", got)
	}

	if !strings.Contains(got, "grpc.method=/bfstore.catalog.v1.CatalogService/ListProducts") {
		t.Fatalf("log output = %q, want grpc.method", got)
	}
}
