package jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerWithUser func(c *gin.Context, userID int)

func WithUser(h HandlerWithUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := GetPrincipal(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		h(c, p.UserID)
	}
}
