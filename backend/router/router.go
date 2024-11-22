package router

import (
	"log"
	"net/http"
	"os"
	"time"
	"word_app/backend/ent"
	"word_app/backend/src/handlers/middleware"
	"word_app/backend/src/handlers/user"
	"word_app/backend/src/handlers/word"
	user_service "word_app/backend/src/service/user"
	word_service "word_app/backend/src/service/word"
	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRouter(router *gin.Engine, client *ent.Client) {
	// リクエストログ用ミドルウェア
	router.Use(requestLoggerMiddleware())

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
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// OPTIONSリクエスト（プリフライトリクエスト）の処理
		if c.Request.Method == "OPTIONS" {
			logrus.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			}).Info("Handling CORS preflight request")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// リクエストの詳細をログに出力するミドルウェア
func requestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// リクエスト処理前
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"ip":     c.ClientIP(),
		}).Info("Incoming request")

		c.Next()

		// リクエスト処理後
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": duration,
		}).Info("Request completed")
	}
}
