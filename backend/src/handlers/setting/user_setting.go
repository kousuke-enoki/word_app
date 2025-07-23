package setting

import (
	"errors"
	"net/http"

	settingUc "word_app/backend/src/usecase/setting"

	"github.com/gin-gonic/gin"
)

func (h *AuthSettingHandler) GetUserConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUserNotFound})
			return
		}
		var req settingUc.InputGetUserConfig
		req.UserID = userID.(int)

		setting, err := h.settingUsecase.GetUser(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, setting)
	}
}

func (h *AuthSettingHandler) SaveUserConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var req settingUc.InputUpdateUserConfig
		req.UserID = userID.(int)
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		setting, err := h.settingUsecase.UpdateUser(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, setting)
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
)
