package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 後続ハンドラでユーザーロールを使えるようにcontextにセットするミドルウェア
// ユーザーID、isAdmin、isRootをcontextにセットする
func (m *JwtMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "authorization header required"})
			return
		}

		roles, err := m.tokenValidator.Validate(c.Request.Context(), token)
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
