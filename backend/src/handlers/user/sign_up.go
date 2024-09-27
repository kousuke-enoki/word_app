package user

import (
	"bytes"
	"context"
	"eng_app/ent"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

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

		// JWTトークンの生成
		token, err := generateJWT(req.Email)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		// 成功時のレスポンス
		c.JSON(201, gin.H{"user": newUser, "token": token})
	}
}
