package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"word_app/backend/src/models"
	"word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// ユーザー検索
		signInUser, err := h.userClient.FindByEmail(context.Background(), req.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// パスワード未設定（外部認証のみ）のユーザーはパスワードサインイン不可
		if err := comparePasswordPtr(signInUser.Password, req.Password); err != nil {
			// エラーメッセージは統一して情報リークを防ぐ
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", signInUser.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful", "token": token})
	}
}

// hashed が nil/空ならエラー。平文とbcryptで比較。
func comparePasswordPtr(hashed *string, plain string) error {
	if hashed == nil || *hashed == "" {
		return errors.New("password not set")
	}
	return bcrypt.CompareHashAndPassword([]byte(*hashed), []byte(plain))
}
