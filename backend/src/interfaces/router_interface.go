package interfaces

import (
	"github.com/gin-gonic/gin"
)

type RouterInterface interface {
	SetupRouter(router *gin.Engine)
}
