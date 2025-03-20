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
		userRoles, err := GetUserRoles(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// ログイン状態を返す
		c.JSON(http.StatusOK, gin.H{
			"message": "Authenticated",
			"userID":  userRoles.UserID,
			"isLogin": true,
			"isAdmin": userRoles.IsAdmin,
			"isRoot":  userRoles.IsRoot,
		})
	}
}
