package src

import (
	"github.com/gin-gonic/gin"
)

// Router interface to abstract the gin.Engine
type Router interface {
	Use(middleware ...gin.HandlerFunc)
	GET(path string, handler gin.HandlerFunc)
	POST(path string, handler gin.HandlerFunc)
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}
