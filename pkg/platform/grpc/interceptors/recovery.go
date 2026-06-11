package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryRecoveryInterceptor returns a gRPC unary server interceptor that recovers
// from panics and converts them into a safe Internal gRPC error.
//
// The interceptor prevents a single panicking request handler from crashing the
// whole service process. Panic details are logged for operators, but the client
// receives a generic Internal error so implementation details are not leaked.
func UnaryRecoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = slog.Default()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.ErrorContext(
					ctx,
					"grpc request panic recovered",
					"grpc.method", info.FullMethod,
					"panic", fmt.Sprint(recovered),
					"stack", string(debug.Stack()),
				)

				resp = nil
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
