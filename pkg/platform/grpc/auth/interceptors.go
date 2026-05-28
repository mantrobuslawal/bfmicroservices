package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const MetadataAuthorisation = "authorization"

// Authoriser checks whether an authenticated principal can call a gRPC method.
type Authoriser interface {
	CanCall(ctx context.Context, principal Principal, method string) (bool, error)
}

// UnaryAuthentication validates authorization metadata and attaches a Principal to context.
func UnaryAuthentication(validator TokenValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		values := metadata.ValueFromIncomingContext(ctx, MetadataAuthorisation)
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization is required")
		}

		principal, err := validator.Validate(ctx, values[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization")
		}

		ctx = ContextWithPrincipal(ctx, principal)

		return handler(ctx, req)
	}
}

// UnaryAuthorisation checks whether the authenticated principal may call the gRPC method.
func UnaryAuthorisation(authoriser Authoriser) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		principal, ok := PrincipalFromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "principal is required")
		}

		allowed, err := authoriser.CanCall(ctx, principal, info.FullMethod)
		if err != nil {
			return nil, status.Error(codes.Internal, "authorisation failed")
		}

		if !allowed {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}

		return handler(ctx, req)
	}
}
