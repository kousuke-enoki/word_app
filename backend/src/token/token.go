package token

import "github.com/golang-jwt/jwt/v4"

// JWTGeneratorは、JWTの生成を行うインターフェース
type JWTGenerator interface {
	GenerateJWT(userID string) (string, error)
}

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
		// 必要に応じてさらにクレームを追加
	})
	return token.SignedString([]byte(j.secretKey))
}
