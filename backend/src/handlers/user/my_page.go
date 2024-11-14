// user/handler.go
package user

import (
	"context"
	"net/http"
	"strconv"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) MyPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userID の取得とチェック
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// userID を int に変換
		userIDInt, err := strconv.Atoi(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		// ユーザー情報の取得
		signInUser, err := h.userClient.FindUserByID(context.Background(), userIDInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		c.JSON(http.StatusOK, models.MyPageResponse{
			User: models.User{
				Name:  signInUser.Name,
				Admin: signInUser.Admin,
			},
		})
	}
}
