package user

import (
	"context"
	"log"
	"net/http"
	"word_app/ent"
	"word_app/ent/user"

	"github.com/gin-gonic/gin"
)

func MyPageHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWTトークンからユーザーのEmailを取得
		email, exists := c.Get("email") // ここは user_id ではなく email に
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Emailでユーザー情報をデータベースから取得
		signInUser, err := client.User.Query().
			Where(user.EmailEQ(email.(string))). // <- クエリ部分
			First(context.Background())

		if err != nil {
			log.Println("Error retrieving user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		// ユーザー情報を返す
		c.JSON(http.StatusOK, gin.H{
			"name":  signInUser.Name,
			"email": email,
		})
	}
}
