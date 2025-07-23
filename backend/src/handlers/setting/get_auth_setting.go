package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthSettingHandler) GetAuthConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rootConfig, err := h.settingUsecase.GetAuth(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}
