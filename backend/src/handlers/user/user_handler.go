// handlers/user_handler.go
package user

import (
	"word_app/backend/src/interfaces"
	"word_app/backend/src/utils"
)

type UserHandler struct {
	userClient   interfaces.UserClient
	jwtGenerator utils.JWTGenerator
}

func NewUserHandler(client interfaces.UserClient, jwtGen utils.JWTGenerator) *UserHandler {
	return &UserHandler{
		userClient:   client,
		jwtGenerator: jwtGen,
	}
}
