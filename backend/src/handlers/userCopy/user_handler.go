// handlers/user_handler.go
package user

import (
	"word_app/backend/src/interfaces/http/auth"
	"word_app/backend/src/interfaces/http/user"
)

type Handler struct {
	userClient   user.Client
	jwtGenerator auth.JWTGenerator
}

func NewHandler(client user.Client, jwtGen auth.JWTGenerator) *Handler {
	return &Handler{
		userClient:   client,
		jwtGenerator: jwtGen,
	}
}
