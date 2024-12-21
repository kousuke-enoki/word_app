package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// JWT トークン生成のモック関数
type JwtGenerator struct {
	mock.Mock
}

func (m *JwtGenerator) GenerateJWT(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", 1) // 常に userID を 1 に設定
		c.Next()
	}
}
