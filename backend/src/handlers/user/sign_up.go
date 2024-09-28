package user

import (
	"bytes"
	"context"
	"io"
	"log"
	"word_app/ent"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		type SignUpRequest struct {
			Email    string `json:"email" binding:"required"`
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

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
		var req SignUpRequest
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

		// レスポンスを返す
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "message": "sign up missing"})
			return
		}
		c.JSON(201, gin.H{"user": newUser})
	}
}
