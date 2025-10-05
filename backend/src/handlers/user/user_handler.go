// handlers/user_handler.go
package user

import (
	"word_app/backend/src/interfaces/http/auth"
	"word_app/backend/src/usecase/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase  user.Usecase
	jwtGenerator auth.JWTGenerator
}

func NewHandler(
	usecase user.Usecase,
	jwtGen auth.JWTGenerator,
) *UserHandler {
	return &UserHandler{
		userUsecase:  usecase,
		jwtGenerator: jwtGen,
	}
}

type Handler interface {
	SignUpHandler() gin.HandlerFunc
	SignInHandler() gin.HandlerFunc
	MyPageHandler() gin.HandlerFunc
	ListHandler() gin.HandlerFunc
	EditHandler() gin.HandlerFunc
	DeleteHandler() gin.HandlerFunc
	MeHandler() gin.HandlerFunc
	ShowHandler() gin.HandlerFunc
}
