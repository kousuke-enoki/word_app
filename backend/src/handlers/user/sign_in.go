package user

import (
	"context"
	"fmt"
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("SignInHandler")
		var req models.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		validationErrors := user.ValidateSignIn(&req)
		if len(validationErrors) > 0 {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// ユーザーの検索
		signInUser, err := h.userClient.FindUserByEmail(context.Background(), req.Email)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(signInUser.Password), []byte(req.Password)) != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		logrus.Info("signInUser")
		logrus.Info(signInUser)
		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", signInUser.ID))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		logrus.Info(token)
		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful", "token": token})
	}
}
