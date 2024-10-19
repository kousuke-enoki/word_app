package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

// AuthMiddleware はJWTトークンの認証をする
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil {
			// トークンが無効な場合にエラーメッセージを詳細にする
			var errorMsg string
			if err == jwt.ErrSignatureInvalid {
				errorMsg = "Invalid token signature"
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					errorMsg = "TokenExpired"
				} else {
					errorMsg = "InvalidToken"
				}
			} else {
				errorMsg = "InvalidToken"
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			userId := claims.Subject
			c.Set("userId", userId)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "InvalidToken"})
			c.Abort()
			return
		}
	}
}
