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

func (h *Handler) SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("1")
		var req models.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		logrus.Info("1")
		validationErrors := user.ValidateSignIn(&req)
		if len(validationErrors) > 0 {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		logrus.Info("2")
		// ユーザーの検索
		signInUser, err := h.userClient.FindByEmail(context.Background(), req.Email)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(signInUser.Password), []byte(req.Password)) != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		logrus.Info("3")
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
