// interfaces/handlers.go
package interfaces

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	SignUpHandler() gin.HandlerFunc
	SignInHandler() gin.HandlerFunc
	MyPageHandler() gin.HandlerFunc
}

type WordHandler interface {
	AllWordListHandler() gin.HandlerFunc
	WordShowHandler() gin.HandlerFunc
}
