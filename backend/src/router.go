package src

import (
	"log"
	"net/http"
	"word_app/backend/ent"
	"word_app/backend/src/adapters"
	"word_app/backend/src/handlers"
	"word_app/backend/src/handlers/middleware"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, client *ent.Client) {
	entClient := adapters.NewEntUserClient(client)
	// JWTGeneratorを初期化
	jwtGenerator := utils.NewMyJWTGenerator("your_secret_key")

	// jwtGeneratorをUserHandlerに渡す
	userHandler := user.NewUserHandler(entClient, jwtGenerator)

	wordHandler := handlers.NewWordHandler(client)

	router.Use(CORSMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.GET("/", func(c *gin.Context) { /* RootHandlerの処理 */ })
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/sign_up", userHandler.SignUpHandler())
		userRoutes.POST("/sign_in", userHandler.SignInHandler())
	}

	// protected := router.Group("/protected")
	router.Use(middleware.AuthMiddleware())
	router.GET("/users/my_page", userHandler.MyPageHandler())
	router.GET("/words/all_list", wordHandler.AllWordListHandler())
	router.GET("/words/:id", wordHandler.WordShowHandler())

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
