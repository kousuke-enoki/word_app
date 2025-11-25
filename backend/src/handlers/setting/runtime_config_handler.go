package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthSettingHandler) GetRuntimeConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		config, err := h.settingUsecase.GetRuntimeConfig(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Cache-Control ヘッダを設定
		c.Header("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
		c.JSON(http.StatusOK, config)
	}
}
