package test

import (
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

// userID とロールをコンテキストに埋め込む
func InjectUser(c *gin.Context, id int, isRoot bool) {
	p := models.Principal{
		UserID:  id,
		IsAdmin: false,
		IsRoot:  isRoot,
		IsTest:  false,
	}
	jwt.SetPrincipal(c, p)
}
