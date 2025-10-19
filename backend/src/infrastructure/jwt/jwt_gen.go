package jwt

import (
	"log"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type MyJWTGenerator struct {
	secretKey string
}
type Claims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

type JWTGenerator interface {
	GenerateJWT(userID string) (string, error)
}

func NewMyJWTGenerator(secretKey string) *MyJWTGenerator {
	if secretKey == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	return &MyJWTGenerator{secretKey: secretKey}
}

func (j *MyJWTGenerator) GenerateJWT(userID string) (string, error) {
	// 有効期限のデフォルト: 1時間
	expirationTime := time.Now().Add(1 * time.Minute)

	// 環境変数から有効期限を取得（オプション）
	if hours := os.Getenv("JWT_EXPIRATION_HOURS"); hours != "" {
		if h, err := strconv.Atoi(hours); err == nil {
			expirationTime = time.Now().Add(time.Duration(h) * time.Hour)
		}
	}
	if minutes := os.Getenv("JWT_EXPIRATION_MINUTES"); minutes != "" {
		if m, err := strconv.Atoi(minutes); err == nil {
			expirationTime = expirationTime.Add(time.Duration(m) * time.Minute)
		}
	}

	// クレームを作成
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}
