// user/handler.go
package user

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) MyPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("mypage")
		log.Println(c)
		// userId の取得とチェック
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		log.Println(userId)
		// userId を int に変換
		userIDInt, err := strconv.Atoi(userId.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		// ユーザー情報の取得
		signInUser, err := h.userClient.FindUserByID(context.Background(), userIDInt)
		if err != nil {
			log.Println("Error retrieving user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		// ユーザー情報を返す
		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"name": signInUser.Name,
			},
		})
	}
}
