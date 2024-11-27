package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	// 環境変数から JWT_SECRET を取得
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return func(c *gin.Context) {
		// Authorization ヘッダーからトークンを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// トークン文字列の先頭にある "Bearer " を取り除く
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// トークンを解析
		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// トークンが有効であり、クレームから userID を取得する
		if claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
			userID := claims.UserID
			if userID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: userID not found"})
				c.Abort()
				return
			}

			// userID を int に変換
			userIDInt, err := strconv.Atoi(fmt.Sprintf("%v", userID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
				return
			}

			// gin.Context に userID を保存
			c.Set("userID", userIDInt)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}
