package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) AuthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ミドルウェアで設定された userID を取得
		userID, exists := c.Get("userID")
		if !exists {
			// 通常ここには到達しないが念のためチェック
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// ログイン状態を返す
		c.JSON(http.StatusOK, gin.H{
			"message": "Authenticated",
			"userID":  userID, // 必要であればユーザーIDを含める
		})
	}
}
