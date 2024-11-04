package interfaces

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	MyPage(c *gin.Context)
}

type WordHandler interface {
	AllWordList(c *gin.Context)
	WordShow(c *gin.Context)
}
