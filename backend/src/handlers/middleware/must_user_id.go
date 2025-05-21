package middleware

import (
	"github.com/gin-gonic/gin"
	// あなたのDB接続用パッケージ
)

func MustUserID(c *gin.Context) (int, error) {
	r, err := GetUserRoles(c)
	if err != nil {
		return 0, err
	}
	return r.UserID, nil
}
