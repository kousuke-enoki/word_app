package router

import (
	"net/http"
	"os"
	"time"

	"word_app/backend/src/handlers/middleware"
	"word_app/backend/src/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouterImplementation struct {
	UserHandler interfaces.UserHandler
	WordHandler interfaces.WordHandler
	JWTSecret   string
}

func NewRouter(userHandler interfaces.UserHandler, wordHandler interfaces.WordHandler) *RouterImplementation {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET environment variable is required")
	}

	return &RouterImplementation{
		UserHandler: userHandler,
		WordHandler: wordHandler,
		JWTSecret:   jwtSecret,
	}
}

func (r *RouterImplementation) SetupRouter(router *gin.Engine) {
	router.Use(requestLoggerMiddleware())
	router.Use(CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/sign_up", r.UserHandler.SignUpHandler())
		userRoutes.POST("/sign_in", r.UserHandler.SignInHandler())
	}

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.GET("/users/my_page", r.UserHandler.MyPageHandler())
		protectedRoutes.GET("/words/all_list", r.WordHandler.AllWordListHandler())
		protectedRoutes.GET("/words/:id", r.WordHandler.WordShowHandler())
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
