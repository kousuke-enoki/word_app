package user

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"word_app/ent"
	"word_app/ent/user"

	"github.com/gin-gonic/gin"
)

func MyPageHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// userIdをint型に変換
		userIDInt, err := strconv.Atoi(userId.(string)) // 変換処理を追加
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}
		// userIdでユーザー情報をデータベースから取得
		signInUser, err := client.User.Query().
			Where(user.ID(userIDInt)). // <- クエリ部分
			First(context.Background())

		if err != nil {
			log.Println("Error retrieving user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		// ユーザー情報を返す
		c.JSON(http.StatusOK, gin.H{
			"name": signInUser.Name,
		})
	}
}
