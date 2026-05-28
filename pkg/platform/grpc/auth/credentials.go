package auth

import "context"

// StaticBearerToken is a simple per-RPC credential that attaches:
//
//	authorization: Bearer <token>
//
// It is suitable for local development or simple internal wiring.
// For production, prefer short-lived tokens, TLS/mTLS, and managed identity where available.
type StaticBearerToken struct {
	Token string
}

// GetRequestMetadata returns authorization metadata for the RPC.
func (c StaticBearerToken) GetRequestMetadata(
	ctx context.Context,
	uri ...string,
) (map[string]string, error) {
	return map[string]string{
		MetadataAuthorisation: "Bearer " + c.Token,
	}, nil
}

// RequireTransportSecurity prevents this token from being sent over insecure transport.
func (c StaticBearerToken) RequireTransportSecurity() bool {
	return true
}
