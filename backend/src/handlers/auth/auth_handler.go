package auth

import (
	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/usecase/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthUsecase  auth.Usecase
	jwtGenerator jwt.JWTGenerator
}

func NewHandler(
	authUsecase auth.Usecase,
	jwtGen jwt.JWTGenerator,
) *AuthHandler {
	return &AuthHandler{
		AuthUsecase:  authUsecase,
		jwtGenerator: jwtGen,
	}
}

type Handler interface {
	LineLogin() gin.HandlerFunc
	LineCallback() gin.HandlerFunc
	LineComplete() gin.HandlerFunc
	TestLoginHandler() gin.HandlerFunc
}
