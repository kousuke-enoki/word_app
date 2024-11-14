package utils_test

import (
	"testing"
	"time"
	"word_app/backend/src/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	userID := "12345"
	jwtGen := &utils.DefaultJWTGenerator{}
	var jwtKey = []byte("your_secret_key")

	tokenString, err := jwtGen.GenerateJWT(userID)
	assert.NoError(t, err, "Token generation should not produce an error")

	// トークンをパースしてクレームを検証
	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	assert.NoError(t, err, "Token parsing should not produce an error")
	assert.NotNil(t, token, "Token should not be nil")
	assert.True(t, token.Valid, "Token should be valid")

	// クレームの検証
	claims, ok := token.Claims.(*utils.Claims)
	assert.True(t, ok, "Claims should be of type *Claims")
	assert.Equal(t, userID, claims.UserID, "UserID in claims should match the generated userID")
	assert.WithinDuration(t, time.Now().Add(1*time.Minute), claims.ExpiresAt.Time, 2*time.Second, "Expiration time should be within 2 seconds of the expected time")
}
