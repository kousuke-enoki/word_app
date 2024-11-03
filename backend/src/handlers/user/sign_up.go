package user

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"word_app/backend/ent"
	"word_app/backend/src/models"
	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// リクエストボディの内容をログに出力
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to read request body"})
			return
		}
		log.Println("Request Body:", string(body))

		// リクエストボディを再度読み取れるようにする
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// リクエストのバインディングと検証
		var req models.SignUpRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("Binding Error:", err)
			c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		// パスワードのハッシュ化
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		// 新しいユーザーの作成
		newUser, err := client.User.
			Create().
			SetEmail(req.Email).
			SetName(req.Name).
			SetPassword(string(hashedPassword)).
			Save(context.Background())

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "message": "sign up failed"})
			return
		}
		log.Println("user", newUser)
		log.Println("user", newUser.ID)
		// サインアップ後にJWTトークンを生成
		token, err := utils.GenerateJWT(fmt.Sprintf("%d", newUser.ID))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}
		log.Println("token", c)
		log.Println("token", token)

		utils.SendTokenResponse(c, token)
	}
}
