package user

import (
	"github.com/gin-gonic/gin"
	"eng_app/ent"
	"eng_app/ent/user"
	"context"
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

		sign_in_user, err := client.User.Query().
			Where(user.EmailEQ(req.Email), user.PasswordEQ(req.Password)).
			First(context.Background())

		if err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		c.JSON(200, gin.H{"user": sign_in_user})
	}
}
