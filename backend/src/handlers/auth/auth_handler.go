package auth

import (
	"word_app/backend/config"
	"word_app/backend/src/infrastructure/jwt"
	"word_app/backend/src/usecase/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthUsecase  auth.Usecase
	jwtGenerator jwt.JWTGenerator
	config       *config.Config
}

func NewHandler(
	authUsecase auth.Usecase,
	jwtGen jwt.JWTGenerator,
	config *config.Config,
) *AuthHandler {
	return &AuthHandler{
		AuthUsecase:  authUsecase,
		jwtGenerator: jwtGen,
		config:       config,
	}
}

type Handler interface {
	LineLogin() gin.HandlerFunc
	LineCallback() gin.HandlerFunc
	LineComplete() gin.HandlerFunc
	TestLoginHandler() gin.HandlerFunc
	TestLogoutHandler() gin.HandlerFunc
	AuthMeHandler() gin.HandlerFunc
}
