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
	AuthHandler    interfaces.AuthHandler
	UserHandler    interfaces.UserHandler
	SettingHandler interfaces.SettingHandler
	WordHandler    interfaces.WordHandler
	QuizHandler    interfaces.QuizHandler
	JWTSecret      string
}

func NewRouter(
	authHandler interfaces.AuthHandler,
	userHandler interfaces.UserHandler,
	settingHandler interfaces.SettingHandler,
	wordHandler interfaces.WordHandler,
	quizHandler interfaces.QuizHandler,
) *RouterImplementation {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET environment variable is required")
	}

	return &RouterImplementation{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		SettingHandler: settingHandler,
		WordHandler:    wordHandler,
		QuizHandler:    quizHandler,
		JWTSecret:      jwtSecret,
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
		protectedRoutes.GET("/auth/check", r.AuthHandler.AuthCheckHandler())

		protectedRoutes.GET("/users/my_page", r.UserHandler.MyPageHandler())
		protectedRoutes.GET("/setting/user_config", r.SettingHandler.GetUserSettingHandler())
		protectedRoutes.POST("/setting/user_config", r.SettingHandler.SaveUserSettingHandler())
		protectedRoutes.GET("/setting/root_config", r.SettingHandler.GetRootSettingHandler())
		protectedRoutes.POST("/setting/root_config", r.SettingHandler.SaveRootSettingHandler())
		protectedRoutes.GET("/words", r.WordHandler.WordListHandler())
		protectedRoutes.GET("/words/:id", r.WordHandler.WordShowHandler())
		protectedRoutes.POST("/words/register", r.WordHandler.RegisterWordHandler())
		protectedRoutes.POST("/words/memo", r.WordHandler.SaveMemoHandler())

		protectedRoutes.POST("/words/new", r.WordHandler.CreateWordHandler())
		protectedRoutes.PUT("/words/:id", r.WordHandler.UpdateWordHandler())
		protectedRoutes.DELETE("/words/:id", r.WordHandler.DeleteWordHandler())
		protectedRoutes.POST("/words/bulk_tokenize", r.WordHandler.BulkTokenizeHandler())
		protectedRoutes.POST("/words/bulk_register", r.WordHandler.BulkRegisterHandler())

		protectedRoutes.POST("/quiz/new", r.QuizHandler.CreateQuizHandler())
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
