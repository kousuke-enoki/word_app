package setting

import (
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/usecase/apperror"
	settingUc "word_app/backend/src/usecase/setting"
	"word_app/backend/src/validators/setting"

	"github.com/gin-gonic/gin"
)

func (h *AuthSettingHandler) GetRootConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		principal, ok := jwt.GetPrincipal(c)
		if !ok || !principal.IsRoot {
			httperr.Write(c, apperror.Unauthorizedf("unauthorized", nil))
			return
		}
		var req settingUc.InputGetRootConfig
		req.UserID = principal.UserID
		rootConfig, err := h.settingUsecase.GetRoot(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}

func (h *AuthSettingHandler) SaveRootConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		principal, ok := jwt.GetPrincipal(c)
		if !ok || !principal.IsRoot {
			httperr.Write(c, apperror.Unauthorizedf("unauthorized", nil))
			return
		}

		var req settingUc.InputUpdateRootConfig
		req.UserID = principal.UserID
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
			return
		}

		validationErrors := setting.ValidateRootConfig(&req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		rootConfig, err := h.settingUsecase.UpdateRoot(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rootConfig)
	}
}
