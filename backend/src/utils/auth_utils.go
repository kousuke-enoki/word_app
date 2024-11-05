package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTGenerator インターフェースを定義
type JWTGenerator interface {
	GenerateJWT(userID string) (string, error)
}

// DefaultJWTGenerator はデフォルトの JWTGenerator 実装
type DefaultJWTGenerator struct{}

func (d *DefaultJWTGenerator) GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("token: ", token)
	return token.SignedString(jwtKey)
	// SendTokenResponseは生成したJWTトークンを含むレスポンスを返します
	//
	//	func SendTokenResponse(c *gin.Context, token string) {
	//		c.JSON(200, gin.H{
	//			"message": "Authentication successful",
	//			"token":   token,
	//		})
}
