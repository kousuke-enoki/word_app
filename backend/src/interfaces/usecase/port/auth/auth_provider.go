package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
)

type AuthProvider interface {
	AuthURL(state, nonce string) string
	Exchange(ctx context.Context, code string) (*Identity, error)
	ValidateNonce(idTok *oidc.IDToken, expected string) error
}

type Identity struct {
	Provider string // "line"
	Subject  string // sub
	Email    string
	Name     string
	jwt.RegisteredClaims
}
