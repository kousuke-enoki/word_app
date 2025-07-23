package setting

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetUserConfigHandler() gin.HandlerFunc
	SaveUserConfigHandler() gin.HandlerFunc
	GetRootConfigHandler() gin.HandlerFunc
	SaveRootConfigHandler() gin.HandlerFunc
	GetAuthConfigHandler() gin.HandlerFunc
}
