// user/handler.go
package user

import (
	"context"
	"net/http"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) MyPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userID の取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// userIDの型チェック
		id, ok := userID.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID type"})
			return
		}

		// ユーザー情報の取得
		signInUser, err := h.userClient.FindByID(context.Background(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		c.JSON(http.StatusOK, models.MyPageResponse{
			User: models.User{
				Name:    signInUser.Name,
				IsAdmin: signInUser.IsAdmin,
				IsRoot:  signInUser.IsRoot,
			},
			IsLogin: true,
		})
	}
}
