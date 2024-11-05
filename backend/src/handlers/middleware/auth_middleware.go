package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

// AuthMiddleware はJWTを使った認証ミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c)
		// Authorizationヘッダーからトークンを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		log.Println(authHeader)
		// トークン文字列の先頭にある "Bearer " を取り除く
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// トークンを解析
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			// トークンの署名アルゴリズムを確認
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// トークンが有効であり、クレームを取得できた場合
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			// クレームからemailを取得
			userId := claims.Subject
			// gin.Context にuserIDを保存
			c.Set("userId", userId)

			// 次のハンドラに処理を渡す
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}
