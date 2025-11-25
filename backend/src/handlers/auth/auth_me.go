// handlers/auth/me.go
package auth

import (
	"net/http"

	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) AuthMeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := jwt.GetPrincipal(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, models.MeResponse{
			User: models.User{
				ID:      p.UserID,
				IsAdmin: p.IsAdmin,
				IsRoot:  p.IsRoot,
				IsTest:  p.IsTest,
			},
		})
	}
}
