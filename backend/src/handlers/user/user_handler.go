// handlers/user_handler.go
package user

import (
	"word_app/backend/src/interfaces/http/auth"
	"word_app/backend/src/interfaces/http/user"
)

type Handler struct {
	userUsecase  user.Usecase
	jwtGenerator auth.JWTGenerator
}

func NewHandler(
	usecase user.Usecase,
	jwtGen auth.JWTGenerator,
) *Handler {
	return &Handler{
		userUsecase:  usecase,
		jwtGenerator: jwtGen,
	}
}
