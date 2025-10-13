// src/handlers/auth/test_logout_handler.go
package auth

import (
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/usecase/apperror"
	"word_app/backend/src/utils/contextutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *AuthHandler) TestLogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("TestLogoutHandler")
		userID, err := contextutil.MustUserID(c)
		if err != nil || userID <= 0 {
			httperr.Write(c, apperror.Unauthorizedf("invalid token", nil))
			return
		}
		logrus.Info(userID)
		if err := h.AuthUsecase.TestLogout(c.Request.Context(), userID); err != nil {
			httperr.Write(c, err)
			return
		}
		logrus.Info("204")
		// 204で十分（本文なし）。クライアント側はトークン破棄＆トップ遷移。
		c.Status(http.StatusNoContent)
	}
}
