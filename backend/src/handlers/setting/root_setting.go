package setting

import (
	"errors"
	"net/http"
	settingUc "word_app/backend/src/usecase/setting"
	"word_app/backend/src/utils/contextutil"
	"word_app/backend/src/validators/setting"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *AuthSettingHandler) GetRootSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
			return
		}
		userRoles, err := contextutil.GetUserRoles(c)
		if err != nil || userRoles == nil || !userRoles.IsRoot {
			if err == nil {
				err = errors.New("unauthorized: root access required")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var req settingUc.InputGetRootConfig
		req.UserID = userID.(int)
		rootConfig, err := h.settingUsecase.GetRoot(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}

func (h *AuthSettingHandler) SaveRootSettingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
			return
		}
		userRoles, err := contextutil.GetUserRoles(c)
		if err != nil || userRoles == nil || !userRoles.IsRoot {
			if err == nil {
				err = errors.New("unauthorized: root access required")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var req settingUc.InputUpdateRootConfig
		req.UserID = userID.(int)
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErrors := setting.ValidateRootConfig(&req)
		if len(validationErrors) > 0 {
			for _, err := range validationErrors {
				logrus.Error(err)
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		rootConfig, err := h.settingUsecase.UpdateRoot(c, req)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}
