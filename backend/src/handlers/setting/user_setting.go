package setting

import (
	"errors"
	"net/http"

	"word_app/backend/src/middleware/jwt"
	settingUc "word_app/backend/src/usecase/setting"

	"github.com/gin-gonic/gin"
)

func (h *AuthSettingHandler) GetUserConfigHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		var req settingUc.InputGetUserConfig
		req.UserID = userID

		setting, err := h.settingUsecase.GetUser(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, setting)
	})
}

func (h *AuthSettingHandler) SaveUserConfigHandler() gin.HandlerFunc {
	return jwt.WithUser(func(c *gin.Context, userID int) {
		ctx := c.Request.Context()
		var req settingUc.InputUpdateUserConfig
		req.UserID = userID
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		setting, err := h.settingUsecase.UpdateUser(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, setting)
	})
}

var ErrUserNotFound = errors.New("user not found")
