package interceptors

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const correlationIDMetadataKey = "x-correlation-id"

type correlationIDContextKey struct{}

// ContextWithCorrelationID returns a child context containing the supplied
// correlation ID.
//
// Empty correlation IDs are still stored as given, but callers should usually
// avoid storing empty IDs. Use CorrelationIDFromContext to safely retrieve only
// non-empty values.
func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDContextKey{}, correlationID)
}

// CorrelationIDFromContext returns the correlation ID stored in ctx.
//
// It returns false when the context has no correlation ID, when the stored value
// is not a string, or when the stored string is empty after trimming whitespace.
func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	correlationID, ok := ctx.Value(correlationIDContextKey{}).(string)
	if !ok {
		return "", false
	}

	correlationID = strings.TrimSpace(correlationID)
	if correlationID == "" {
		return "", false
	}

	return correlationID, true
}

// UnaryCorrelationIDInterceptor returns a gRPC unary server interceptor that
// ensures every request has a correlation ID.
//
// The interceptor reads x-correlation-id from incoming gRPC metadata. If the
// caller supplied a non-empty value, that value is reused. If no usable value is
// present, a new correlation ID is generated.
//
// The correlation ID is stored in the request context and sent back as response
// header metadata so callers can see which ID the service used.
func UnaryCorrelationIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		correlationID := correlationIDFromIncomingMetadata(ctx)
		if correlationID == "" {
			correlationID = newCorrelationID()
		}

		ctx = ContextWithCorrelationID(ctx, correlationID)

		_ = grpc.SetHeader(ctx, metadata.Pairs(
			correlationIDMetadataKey,
			correlationID,
		))

		return handler(ctx, req)
	}
}

func correlationIDFromIncomingMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(correlationIDMetadataKey)
	if len(values) == 0 {
		return ""
	}

	return strings.TrimSpace(values[0])
}

func newCorrelationID() string {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(b[:])
}
