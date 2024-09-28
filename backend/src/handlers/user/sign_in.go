package user

import (
	"context"
	"time"
	"word_app/ent"
	"word_app/ent/user"

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
			Subject:   email,
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

func SignInHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		type SignInRequest struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var req SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// ユーザーの検索
		signInUser, err := client.User.Query().
			Where(user.EmailEQ(req.Email)).
			First(context.Background())

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// パスワードの検証
		err = bcrypt.CompareHashAndPassword([]byte(signInUser.Password), []byte(req.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// JWTトークンの生成
		token, err := generateJWT(req.Email)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, gin.H{
			"message": "Sign in successful",
			"token":   token,
		})
	}
}
