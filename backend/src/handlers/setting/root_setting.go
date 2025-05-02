package setting

import (
	"errors"
	"net/http"
	"word_app/backend/src/handlers/middleware"
	"word_app/backend/src/models"
	"word_app/backend/src/validators/setting"

	"github.com/gin-gonic/gin"
)

func (h *SettingHandler) GetRootSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
			return
		}
		userRoles, err := middleware.GetUserRoles(c)
		if err != nil || userRoles == nil || !userRoles.IsRoot {
			if err == nil {
				err = errors.New("unauthorized: root access required")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		rootConfig, err := h.settingService.GetRootConfig(c, userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}

func (h *SettingHandler) SaveRootSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
			return
		}
		userRoles, err := middleware.GetUserRoles(c)
		if err != nil || userRoles == nil || !userRoles.IsRoot {
			if err == nil {
				err = errors.New("unauthorized: root access required")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var req models.RootConfig
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErrors := setting.ValidateRootConfig(&req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		rootConfig, err := h.settingService.UpdateRootConfig(
			c, userID.(int), req.EditingPermission, req.IsTestUserMode, req.IsEmailAuthCheck, req.IsLineAuth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}
