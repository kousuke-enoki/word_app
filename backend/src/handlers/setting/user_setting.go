package setting

import (
	"errors"
	"net/http"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *SettingHandler) GetUserSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUserNotFound})
			return
		}
		setting, err := h.settingService.GetUserConfig(c, userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, setting)
	}
}

func (h *SettingHandler) SaveUserSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var req models.UserConfig
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		setting, err := h.settingService.UpdateUserConfig(c, userID.(int), req.IsDarkMode)
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
