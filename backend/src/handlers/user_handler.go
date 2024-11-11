// handlers/user_handler.go
package handlers

import (
	"word_app/backend/ent"
	"word_app/backend/src/handlers/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	client *ent.Client
}

func NewUserHandler(client *ent.Client) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) SignUpHandler() gin.HandlerFunc {
	return user.SignUpHandler(h.client)
}

func (h *UserHandler) SignInHandler() gin.HandlerFunc {
	return user.SignInHandler(h.client)
}

func (h *UserHandler) MyPageHandler() gin.HandlerFunc {
	return user.MyPageHandler(h.client)
}
