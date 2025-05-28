package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

type AuthProvider interface {
	AuthURL(state, nonce string) string
	Exchange(ctx context.Context, code string) (*Identity, error)
	ValidateNonce(idToken, nonce string) error
}

type Identity struct {
	Provider string // "line"
	Subject  string // sub
	Email    string
	Name     string
	jwt.RegisteredClaims
}
