package user

import (
	"context"
	"eng_app/ent"
	"eng_app/ent/user"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func MyPageHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("MyPageHandler called")
		tokenString := c.GetHeader("Authorization")
		log.Println("c", c)
		log.Println("dc", client)
		log.Println("asdf", tokenString)
		if tokenString == "" {
			log.Println("Authorization header is required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		log.Println("Authorization header:", tokenString)

		// "Bearer "プレフィックスを削除
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		user, err := client.User.
			Query().
			Where(user.EmailEQ(claims.Email)).
			Only(context.Background())

		if err != nil {
			log.Println("User not found:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		log.Println("User found:", user.Name)
		c.JSON(http.StatusOK, gin.H{"name": user.Name, "date": time.Now().Format("2006-01-02")})
	}
}
