// handlers/user_handler.go
package user

import (
	"word_app/backend/src/interfaces"
	"word_app/backend/src/interfaces/http/auth"
)

type UserHandler struct {
	userClient   interfaces.UserClient
	jwtGenerator auth.JWTGenerator
}

func NewUserHandler(client interfaces.UserClient, jwtGen auth.JWTGenerator) *UserHandler {
	return &UserHandler{
		userClient:   client,
		jwtGenerator: jwtGen,
	}
}
