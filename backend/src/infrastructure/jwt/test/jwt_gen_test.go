package jwt_test

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	jwt_infra "word_app/backend/src/infrastructure/jwt"
)

func TestGenerateJWT(t *testing.T) {
	// 1. 共通キーを準備
	testSecret := "test_secret_key"
	_ = os.Setenv("JWT_SECRET", testSecret)

	userID := "12345"

	// 2. 正しく初期化
	jwtGen := jwt_infra.NewMyJWTGenerator(testSecret)

	// 3. トークン生成
	tokenString, err := jwtGen.GenerateJWT(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString, "token string must not be empty")

	// 4. パース & 検証
	tok, err := jwt.ParseWithClaims(tokenString, &jwt_infra.Claims{}, func(tk *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	assert.NoError(t, err, "parsing/verification failed")
	assert.True(t, tok.Valid, "token should be valid")

	// 5. クレームの検証
	claims := tok.Claims.(*jwt_infra.Claims)
	assert.Equal(t, userID, claims.UserID)
	assert.WithinDuration(t,
		time.Now().Add(1*time.Hour),
		claims.ExpiresAt.Time,
		2*time.Second,
	)
}
