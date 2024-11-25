package user

import (
	"context"
	"fmt"
	"net/http"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// ユーザーの検索
		signInUser, err := h.userClient.FindUserByEmail(context.Background(), req.Email)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(signInUser.Password), []byte(req.Password)) != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", signInUser.ID))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful", "token": token})
	}
}
