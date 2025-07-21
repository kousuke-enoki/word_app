package auth

import (
	"word_app/backend/src/interfaces/http/auth"
)

type Handler struct {
	AuthUsecase  auth.AuthUsecase
	jwtGenerator auth.JWTGenerator
}

func NewHandler(
	authUsecase auth.AuthUsecase,
	jwtGen auth.JWTGenerator,
) *Handler {
	return &Handler{
		AuthUsecase:  authUsecase,
		jwtGenerator: jwtGen,
	}
}
