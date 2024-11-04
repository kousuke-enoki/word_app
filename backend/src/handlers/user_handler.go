package handlers

import (
	"word_app/backend/ent"
	"word_app/backend/src/handlers/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	client *ent.Client
}

func NewUserHandler(client *ent.Client) *userHandler {
	return &userHandler{client: client}
}

func (h *UserHandler) SignUpHandler() gin.HandlerFunc {
	return user.SignUpHandler(h.client)
}

func (h *userHandler) SignIn(c *gin.Context) {
	// サインイン処理
}

func (h *userHandler) MyPage(c *gin.Context) {
	// マイページ取得処理
}
