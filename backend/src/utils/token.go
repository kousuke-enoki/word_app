package utils

import jwt "github.com/golang-jwt/jwt/v4"

type MyJWTGenerator struct {
	secretKey string
}

func NewMyJWTGenerator(secretKey string) *MyJWTGenerator {
	return &MyJWTGenerator{secretKey: secretKey}
}

func (j *MyJWTGenerator) GenerateJWT(userID string) (string, error) {
	// JWTトークンの生成ロジック
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		// "Subject": userId,
		// 必要に応じてさらにクレームを追加
	})
	return token.SignedString([]byte(j.secretKey))
}
