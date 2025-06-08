package auth

import (
	"word_app/backend/src/interfaces/http/auth"
)

type AuthHandler struct {
	AuthUsecase  auth.AuthUsecase
	jwtGenerator auth.JWTGenerator
}

func NewAuthHandler(
	authUsecase auth.AuthUsecase,
	jwtGen auth.JWTGenerator,
) *AuthHandler {
	return &AuthHandler{
		AuthUsecase:  authUsecase,
		jwtGenerator: jwtGen,
	}
}
