package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type UserRoles struct {
	UserID  int
	IsAdmin bool
	IsRoot  bool
}

// ユーザーのロールでハンドラーごとのアクセス制限をする
func GetUserRoles(c *gin.Context) (*UserRoles, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("user ID not found in context")
	}

	isAdmin, _ := c.Get("isAdmin")
	isRoot, _ := c.Get("isRoot")

	return &UserRoles{
		UserID:  userID.(int),
		IsAdmin: isAdmin.(bool),
		IsRoot:  isRoot.(bool),
	}, nil
}
