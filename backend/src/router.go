package src

import (
	"eng_app/ent"
	"eng_app/src/handlers"
	"eng_app/src/handlers/user"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, client *ent.Client) {
	router.Use(CORSMiddleware())

	router.GET("/", handlers.RootHandler)
	router.POST("users/sign_up", user.SignUpHandler(client))
	router.POST("users/sign_in", user.SignInHandler(client))

	// リクエストの詳細をログに出力
	router.Use(func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		body := c.Request.Body
		log.Printf("Request: %s %s, Query: %s, Body: %v, Status: %d", method, path, query, body, status)
	})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
