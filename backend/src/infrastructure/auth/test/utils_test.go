package utils_test

import (
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	// // テスト用の JWT_SECRET を設定
	// testSecret := "test_secret_key"
	// err := os.Setenv("JWT_SECRET", testSecret)
	// assert.NoError(t, err, "Setting JWT_SECRET should not produce an error")

	// userID := "12345"
	// jwtGen := &utils.DefaultJWTGenerator{}

	// // JWT トークンを生成
	// tokenString, err := jwtGen.GenerateJWT(userID)
	// assert.NoError(t, err, "Token generation should not produce an error")

	// // トークンをパースしてクレームを検証
	// token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
	// 	// JWT_SECRET を取得して返す
	// 	secret := os.Getenv("JWT_SECRET")
	// 	if secret == "" {
	// 		return nil, jwt.ErrSignatureInvalid
	// 	}
	// 	return []byte(secret), nil
	// })
	// assert.NoError(t, err, "Token parsing should not produce an error")
	// assert.NotNil(t, token, "Token should not be nil")
	// assert.True(t, token.Valid, "Token should be valid")

	// // クレームの検証
	// claims, ok := token.Claims.(*utils.Claims)
	// assert.True(t, ok, "Claims should be of type *Claims")
	// assert.Equal(t, userID, claims.UserID, "UserID in claims should match the generated userID")
	// assert.WithinDuration(t, time.Now().Add(1*time.Hour), claims.ExpiresAt.Time, 2*time.Second, "Expiration time should be within 2 seconds of the expected time")
}
