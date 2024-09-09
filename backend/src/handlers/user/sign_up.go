package user

import (
	"github.com/gin-gonic/gin"
	"eng_app/ent"
	"context"
)

func SignUpHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		type SignUpRequest struct {
			Email    string `json:"email" binding:"required"`
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		var req SignUpRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		newUser, err := client.User.
			Create().
			SetEmail(req.Email).
			SetName(req.Name).
			SetPassword(req.Password).
			Save(context.Background())

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "message": "sign up missing"})
			return
		}
		c.JSON(201, gin.H{"user": newUser})
	}
}
