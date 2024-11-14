package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

// JWTGenerator インターフェースを定義
type JWTGenerator interface {
	GenerateJWT(userID string) (string, error)
}

// DefaultJWTGenerator はデフォルトの JWTGenerator 実装
type DefaultJWTGenerator struct{}

func (d *DefaultJWTGenerator) GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) // 本番はDayを想定
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
