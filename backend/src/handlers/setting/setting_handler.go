// handlers/setting_handler.go
package setting

import (
	settingUc "word_app/backend/src/usecase/setting"

	"github.com/gin-gonic/gin"
)

type AuthSettingHandler struct {
	settingUsecase settingUc.SettingFacade
}

func NewAuthSettingHandler(client settingUc.SettingFacade) *AuthSettingHandler {
	return &AuthSettingHandler{
		settingUsecase: client,
	}
}

type SettingHandler interface {
	GetUserSettingHandler() gin.HandlerFunc
	SaveUserSettingHandler() gin.HandlerFunc
	GetRootSettingHandler() gin.HandlerFunc
	SaveRootSettingHandler() gin.HandlerFunc
	GetAuthSettingHandler() gin.HandlerFunc
}
