// src/interface/http/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"word_app/backend/src/interfaces"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(v interfaces.TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: Bearer xxx
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "authorization header required"})
			return
		}

		roles, err := v.Validate(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": err.Error()})
			return
		}

		// 後続ハンドラで使えるようにセット
		c.Set("userID", roles.UserID)
		c.Set("isAdmin", roles.IsAdmin)
		c.Set("isRoot", roles.IsRoot)

		c.Next()
	}
}
