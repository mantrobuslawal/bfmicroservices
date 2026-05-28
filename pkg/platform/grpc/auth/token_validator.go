package auth

import (
	"context"
	"errors"
)

var ErrInvalidToken = errors.New("invalid token")

// TokenValidator validates raw authorisation metadata, usually in the form:
//
//	Authorization: Bearer <token>
type TokenValidator interface {
	Validate(ctx context.Context, rawAuthorisation string) (Principal, error)
}

// LocalDevValidator is a simple validator intended for local development only.
type LocalDevValidator struct {
	Token   string
	Subject string
	Roles   []string
	Service string
}

// Validate accepts exactly "Bearer <Token>" and rejects everything else.
func (v LocalDevValidator) Validate(ctx context.Context, rawAuthorisation string) (Principal, error) {
	expected := "Bearer " + v.Token
	if v.Token == "" || rawAuthorisation != expected {
		return Principal{}, ErrInvalidToken
	}

	return Principal{
		Subject: v.Subject,
		Roles:   append([]string(nil), v.Roles...),
		Service: v.Service,
	}, nil
}
