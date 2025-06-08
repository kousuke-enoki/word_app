package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *SettingHandler) GetAuthSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userID, ok := c.Get("userID")
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		// 	return
		// }
		// userRoles, err := contextutil.GetUserRoles(c)
		// if err != nil || userRoles == nil || !userRoles.IsRoot {
		// 	if err == nil {
		// 		err = errors.New("unauthorized: root access required")
		// 	}
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// 	return
		// }
		rootConfig, err := h.settingService.GetAuthConfig(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}
