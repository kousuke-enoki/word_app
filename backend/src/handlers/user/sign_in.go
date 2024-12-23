package user

import (
	"context"
	"fmt"
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid request"})
			return
		}

		validationErrors := user.ValidateSignIn(&req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid request"})
			return
		}

		// ユーザーの検索
		signInUser, err := h.userClient.FindUserByEmail(context.Background(), req.Email)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(signInUser.Password), []byte(req.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid request"})
			return
		}

		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", signInUser.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful", "token": token})
	}
}
