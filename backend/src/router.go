package src

import (
	"log"
	"net/http"
	"word_app/backend/ent"
	"word_app/backend/src/handlers"
	"word_app/backend/src/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, client *ent.Client) {
	userHandler := handlers.NewUserHandler(client)
	wordHandler := handlers.NewWordHandler(client)

	router.Use(CORSMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.GET("/", func(c *gin.Context) { /* RootHandlerの処理 */ })
	router.POST("/users/sign_up", userHandler.SignUpHandler())
	router.POST("/users/sign_in", userHandler.SignInHandler())

	// router.POST("/users/sign_up", userHandler.SignUp)
	// router.POST("/users/sign_in", userHandler.SignIn)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/users/my_page", userHandler.MyPageHandler())
	protected.GET("/words/all_list", wordHandler.AllWordListHandler())
	protected.GET("/words/:id", wordHandler.WordShowHandler())
	// protected.GET("/users/my_page", userHandler.MyPage)
	// protected.GET("/words/all_list", wordHandler.AllWordList)
	// protected.GET("/words/:id", wordHandler.WordShow)

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
