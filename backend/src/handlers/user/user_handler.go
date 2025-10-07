// handlers/user_handler.go
package user

import (
	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/usecase/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase  user.Usecase
	jwtGenerator jwt.JWTGenerator
}

func NewHandler(
	usecase user.Usecase,
	jwtGen jwt.JWTGenerator,
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
