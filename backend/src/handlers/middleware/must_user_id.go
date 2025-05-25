package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	// あなたのDB接続用パッケージ
)

func MustUserID(c *gin.Context) (int, error) {
	r, err := GetUserRoles(c)
	if err != nil {
		return 0, errors.New("unauthorized: userID not found in context")
	}
	return r.UserID, nil
}
