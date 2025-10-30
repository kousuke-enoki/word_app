// src/handlers/auth/test_logout_handler.go
package auth

import (
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) TestLogoutHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		if err := h.AuthUsecase.TestLogout(c.Request.Context(), userID); err != nil {
			httperr.Write(c, err)
			return
		}

		// 204で十分（本文なし）。クライアント側はトークン破棄＆トップ遷移。
		c.AbortWithStatus(http.StatusNoContent)
	})
}
