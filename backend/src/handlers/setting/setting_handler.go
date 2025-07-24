// handlers/setting_handler.go
package setting

import (
	settingUc "word_app/backend/src/usecase/setting"

	"github.com/gin-gonic/gin"
)

type AuthSettingHandler struct {
	settingUsecase settingUc.SettingFacade
}

func NewHandler(client settingUc.SettingFacade) *AuthSettingHandler {
	return &AuthSettingHandler{
		settingUsecase: client,
	}
}

type Handler interface {
	GetUserConfigHandler() gin.HandlerFunc
	SaveUserConfigHandler() gin.HandlerFunc
	GetRootConfigHandler() gin.HandlerFunc
	SaveRootConfigHandler() gin.HandlerFunc
	GetAuthConfigHandler() gin.HandlerFunc
}
