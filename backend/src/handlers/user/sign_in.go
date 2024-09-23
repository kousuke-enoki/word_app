package user

import (
	"context"
	"eng_app/ent"
	"eng_app/ent/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		type SignInRequest struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		var req SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// ユーザーの検索
		signInUser, err := client.User.Query().
			Where(user.EmailEQ(req.Email)).
			First(context.Background())

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// パスワードの検証
		err = bcrypt.CompareHashAndPassword([]byte(signInUser.Password), []byte(req.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(200, gin.H{"user": signInUser})
	}
}
