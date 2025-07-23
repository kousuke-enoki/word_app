package auth

import (
	"word_app/backend/src/interfaces/http/auth"
)

type Handler struct {
	AuthUsecase  auth.Usecase
	jwtGenerator auth.JWTGenerator
}

func NewHandler(
	authUsecase auth.Usecase,
	jwtGen auth.JWTGenerator,
) *Handler {
	return &Handler{
		AuthUsecase:  authUsecase,
		jwtGenerator: jwtGen,
	}
}
