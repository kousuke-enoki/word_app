package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// JWT トークン生成のモック関数
type MockJwtGenerator struct {
	mock.Mock
}

func (m *MockJwtGenerator) GenerateJWT(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", 1) // 常に userID を 1 に設定
		c.Next()
	}
}

func MockAuthMiddlewareWithoutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userIDを設定しない
		c.Next()
	}
}

func MockAuthMiddlewareWithInvalidUserType() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userIDに不正な型を設定
		c.Set("userID", "invalid_user_type")
		c.Next()
	}
}
