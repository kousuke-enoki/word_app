package user

import (
	"context"
	"fmt"
	"word_app/ent"
	"word_app/ent/user"
	"word_app/src/utils"

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

		if err != nil || bcrypt.CompareHashAndPassword([]byte(signInUser.Password), []byte(req.Password)) != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := utils.GenerateJWT(fmt.Sprintf("%d", signInUser.ID))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		utils.SendTokenResponse(c, token)
	}
}
