package src

import (
	"log"
	"net/http"
	"os"
	"word_app/backend/ent"
	"word_app/backend/src/handlers/middleware"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/handlers/word"
	user_service "word_app/backend/src/service/user"
	word_service "word_app/backend/src/service/word"
	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, client *ent.Client) {
	entClient := user_service.NewEntUserClient(client)
	wordClient := word_service.NewWordService(client)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	jwtGenerator := utils.NewMyJWTGenerator(jwtSecret)

	// jwtGeneratorをUserHandlerに渡す
	userHandler := user.NewUserHandler(entClient, jwtGenerator)

	wordHandler := word.NewWordHandler(wordClient)

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

	// 認証が必要なルート
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.GET("/users/my_page", userHandler.MyPageHandler())

		protectedRoutes.GET("/words/all_list", func(c *gin.Context) {
			wordHandler.AllWordListHandler(c)
		})
		protectedRoutes.GET("/words/:id", func(c *gin.Context) {
			wordHandler.WordShowHandler(c)
		})
	}

	// リクエストログ用ミドルウェア
	router.Use(requestLoggerMiddleware())
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

// リクエストの詳細をログに出力するミドルウェア
func requestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		log.Printf("Request: %s %s, Query: %s, Status: %d", method, path, query, status)
	}
}
