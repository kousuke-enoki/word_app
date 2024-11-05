// handlers/user_handler.go
package user

import (
	"word_app/backend/src/interfaces"
)

type UserHandler struct {
	userClient interfaces.UserClient
}

func NewUserHandler(userClient interfaces.UserClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}
