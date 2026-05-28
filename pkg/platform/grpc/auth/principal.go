package auth

import "context"

type principalContextKey struct{}

// Principal represents an authenticated caller.
//
// Subject may be a user ID, service ID, or token subject.
// Roles may contain coarse-grained roles such as customer, admin, or service.
// Service may contain an internal service identity such as order-service.
type Principal struct {
	Subject string
	Roles   []string
	Service string
}

// ContextWithPrincipal returns a child context containing the authenticated principal.
func ContextWithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

// PrincipalFromContext returns the authenticated principal from context, if present.
func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(Principal)
	return principal, ok
}
