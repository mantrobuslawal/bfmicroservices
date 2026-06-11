package interceptors

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = slog.Default()
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := status.Code(err)

		attrs := []any{
			"grpc.method", info.FullMethod,
			"grpc.code", code.String(),
			"duration_ms", duration.Milliseconds(),
		}

		if correlationID, ok := CorrelationIDFromContext(ctx); ok {
			attrs = append(attrs, "correlation_id", correlationID)
		}

		if err != nil {
			logger.ErrorContext(ctx, "grpc request failed", append(attrs, "error", err)...)
			return resp, err
		}

		logger.InfoContext(ctx, "grpc request completed", attrs...)

		return resp, nil
	}
}
