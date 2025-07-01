package setting

import (
	"github.com/gin-gonic/gin"
)

type SettingHandler interface {
	GetUserSettingHandler() gin.HandlerFunc
	SaveUserSettingHandler() gin.HandlerFunc
	GetRootSettingHandler() gin.HandlerFunc
	SaveRootSettingHandler() gin.HandlerFunc
	GetAuthSettingHandler() gin.HandlerFunc
}
