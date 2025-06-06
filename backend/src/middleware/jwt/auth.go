package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (m *JwtMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			logrus.Warn("AuthMiddleware: empty token for", c.FullPath())
			c.AbortWithStatus(403)
			return
		}
		logrus.Info(token)
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
