package src

import (
	"log"
	"net/http"
	"word_app/backend/ent"
	"word_app/backend/src/handlers"
	"word_app/backend/src/handlers/middleware"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/handlers/word"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, client *ent.Client) {
	router.Use(CORSMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.GET("/", handlers.RootHandler)
	router.POST("/users/sign_up", user.SignUpHandler(client))
	router.POST("/users/sign_in", user.SignInHandler(client))

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/users/my_page", user.MyPageHandler(client))
	protected.GET("/words/all_list", func(c *gin.Context) {
		word.AllWordListHandler(c, client)
	})
	protected.GET("/words/:id", func(c *gin.Context) {
		word.WordShowHandler(c, client)
	})
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

		// OPTIONSリクエスト（プリフライトリクエスト）の処理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
