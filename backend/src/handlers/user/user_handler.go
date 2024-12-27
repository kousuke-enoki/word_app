// handlers/user_handler.go
package user

import (
	"word_app/backend/src/interfaces"
)

type UserHandler struct {
	userClient   interfaces.UserClient
	jwtGenerator interfaces.JWTGenerator
}

func NewUserHandler(client interfaces.UserClient, jwtGen interfaces.JWTGenerator) *UserHandler {
	return &UserHandler{
		userClient:   client,
		jwtGenerator: jwtGen,
	}
}
