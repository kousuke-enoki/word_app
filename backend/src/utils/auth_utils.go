package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

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
	// JWT_SECRET の取得とチェック
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// 環境変数から有効期限を取得
	expirationHoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
	expirationMinutesStr := os.Getenv("JWT_EXPIRATION_MINUTES")

	// デフォルト値を設定 (1時間)
	expirationTime := time.Now().Add(1 * time.Hour)

	// 有効期限の時間を設定
	if expirationHoursStr != "" {
		expirationHours, err := strconv.Atoi(expirationHoursStr)
		if err != nil {
			log.Fatalf("Invalid JWT_EXPIRATION_HOURS: %v", err)
		}
		expirationTime = time.Now().Add(time.Duration(expirationHours) * time.Hour)
	}

	// 有効期限の分を追加設定（オプション）
	if expirationMinutesStr != "" {
		expirationMinutes, err := strconv.Atoi(expirationMinutesStr)
		if err != nil {
			log.Fatalf("Invalid JWT_EXPIRATION_MINUTES: %v", err)
		}
		expirationTime = expirationTime.Add(time.Duration(expirationMinutes) * time.Minute)
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
	return token.SignedString([]byte(jwtSecret))
}
