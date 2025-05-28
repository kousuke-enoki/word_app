package auth

import (
	"context"
	"word_app/backend/src/infrastructure/auth"

	"github.com/golang-jwt/jwt/v4"
)

type AuthProvider interface {
	AuthURL(state, nonce string) string
	Exchange(ctx context.Context, code string) (*auth.Identity, error)
	ValidateNonce(idToken, nonce string) error
}

type Identity struct {
	Provider string // "line"
	Subject  string // sub
	Email    string
	Name     string
	jwt.RegisteredClaims
}
