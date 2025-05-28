package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type TempJWT struct {
	secret []byte
}

func New(secret string) *TempJWT { return &TempJWT{secret: []byte(secret)} }

type Identity struct {
	Provider string `json:"provider"`
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func (t *TempJWT) GenerateTemp(id *Identity, ttl time.Duration) (string, error) {
	id.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, id)
	return token.SignedString(t.secret)
}

func (t *TempJWT) ParseTemp(tok string) (*Identity, error) {
	var id Identity
	_, err := jwt.ParseWithClaims(tok, &id, func(token *jwt.Token) (interface{}, error) {
		return t.secret, nil
	})
	return &id, err
}
