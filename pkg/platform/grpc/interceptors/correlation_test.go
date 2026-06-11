package interceptors

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUnaryCorrelationIDInterceptorPreservesIncomingCorrelationID(t *testing.T) {
	t.Parallel()

	interceptor := UnaryCorrelationIDInterceptor()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(correlationIDMetadataKey, "local-dev-123"),
	)

	_, err := interceptor(
		ctx,
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListProducts",
		},
		func(ctx context.Context, req any) (any, error) {
			got, ok := CorrelationIDFromContext(ctx)
			if !ok {
				t.Fatal("CorrelationIDFromContext() ok = false, want true")
			}

			if got != "local-dev-123" {
				t.Fatalf("correlation ID = %q, want %q", got, "local-dev-123")
			}

			return "response", nil
		},
	)

	if err != nil {
		t.Fatalf("interceptor returned error = %v, want nil", err)
	}
}

func TestUnaryCorrelationIDInterceptorGeneratesCorrelationIDWhenMissing(t *testing.T) {
	t.Parallel()

	interceptor := UnaryCorrelationIDInterceptor()

	_, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListProducts",
		},
		func(ctx context.Context, req any) (any, error) {
			got, ok := CorrelationIDFromContext(ctx)
			if !ok {
				t.Fatal("CorrelationIDFromContext() ok = false, want true")
			}

			if got == "" {
				t.Fatal("correlation ID = empty, want non-empty")
			}

			return "response", nil
		},
	)

	if err != nil {
		t.Fatalf("interceptor returned error = %v, want nil", err)
	}
}

func TestUnaryCorrelationIDInterceptorGeneratesCorrelationIDWhenIncomingValueIsEmpty(t *testing.T) {
	t.Parallel()

	interceptor := UnaryCorrelationIDInterceptor()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(correlationIDMetadataKey, "   "),
	)

	_, err := interceptor(
		ctx,
		"request",
		&grpc.UnaryServerInfo{
			FullMethod: "/bfstore.catalog.v1.CatalogService/ListCategories",
		},
		func(ctx context.Context, req any) (any, error) {
			got, ok := CorrelationIDFromContext(ctx)
			if !ok {
				t.Fatal("CorrelationIDFromContext() ok = false, want true")
			}

			if got == "" {
				t.Fatal("correlation ID = empty, want non-empty")
			}

			if got == "   " {
				t.Fatal("correlation ID preserved whitespace-only value, want generated value")
			}

			return "response", nil
		},
	)

	if err != nil {
		t.Fatalf("interceptor returned error = %v, want nil", err)
	}
}

func TestUnaryCorrelationIDInterceptorReturnsHandlerResponse(t *testing.T) {
	t.Parallel()

	interceptor := UnaryCorrelationIDInterceptor()

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

func TestUnaryCorrelationIDInterceptorReturnsHandlerError(t *testing.T) {
	t.Parallel()

	interceptor := UnaryCorrelationIDInterceptor()

	wantErr := status.Error(codes.InvalidArgument, "invalid request")

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

func TestUnaryCorrelationIDInterceptorCallsHandlerOnce(t *testing.T) {
	t.Parallel()

	interceptor := UnaryCorrelationIDInterceptor()

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

func TestContextWithCorrelationIDAndCorrelationIDFromContext(t *testing.T) {
	t.Parallel()

	ctx := ContextWithCorrelationID(context.Background(), "local-dev-123")

	got, ok := CorrelationIDFromContext(ctx)
	if !ok {
		t.Fatal("CorrelationIDFromContext() ok = false, want true")
	}

	if got != "local-dev-123" {
		t.Fatalf("correlation ID = %q, want %q", got, "local-dev-123")
	}
}

func TestCorrelationIDFromContextReturnsFalseWhenMissing(t *testing.T) {
	t.Parallel()

	got, ok := CorrelationIDFromContext(context.Background())

	if ok {
		t.Fatal("CorrelationIDFromContext() ok = true, want false")
	}

	if got != "" {
		t.Fatalf("correlation ID = %q, want empty string", got)
	}
}

func TestCorrelationIDFromContextReturnsFalseWhenEmpty(t *testing.T) {
	t.Parallel()

	ctx := ContextWithCorrelationID(context.Background(), "   ")

	got, ok := CorrelationIDFromContext(ctx)

	if ok {
		t.Fatal("CorrelationIDFromContext() ok = true, want false")
	}

	if got != "" {
		t.Fatalf("correlation ID = %q, want empty string", got)
	}
}
