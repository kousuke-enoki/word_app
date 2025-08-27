package router

import (
	"net/http"
	"time"

	"word_app/backend/src/interfaces/http/auth"
	middleware_interface "word_app/backend/src/interfaces/http/middleware"
	"word_app/backend/src/interfaces/http/quiz"
	"word_app/backend/src/interfaces/http/result"
	"word_app/backend/src/interfaces/http/setting"
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/interfaces/http/word"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Implementation struct {
	JwtMiddleware  middleware_interface.Middleware
	AuthHandler    auth.Handler
	UserHandler    user.Handler
	SettingHandler setting.Handler
	WordHandler    word.Handler
	QuizHandler    quiz.Handler
	ResultHandler  result.Handler
	// JWTSecret      string
}

func NewRouter(
	jwtMiddleware middleware_interface.Middleware,
	authHandler auth.Handler,
	userHandler user.Handler,
	settingHandler setting.Handler,
	wordHandler word.Handler,
	quizHandler quiz.Handler,
	resultHandler result.Handler,
) *Implementation {
	// jwtSecret := os.Getenv("JWT_SECRET")
	// if jwtSecret == "" {
	// 	logrus.Fatal("JWT_SECRET environment variable is required")
	// }

	return &Implementation{
		JwtMiddleware:  jwtMiddleware,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		SettingHandler: settingHandler,
		WordHandler:    wordHandler,
		QuizHandler:    quizHandler,
		ResultHandler:  resultHandler,
		// JWTSecret:      jwtSecret,
	}
}

// ルートを取り付ける関数
func (r *Implementation) MountRoutes(router *gin.Engine) {
	router.Use(requestLoggerMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/sign_up", r.UserHandler.SignUpHandler())
		userRoutes.POST("/sign_in", r.UserHandler.SignInHandler())
		userRoutes.GET("/auth/line/login", r.AuthHandler.LineLogin())
		userRoutes.GET("/auth/line/callback", r.AuthHandler.LineCallback())
		userRoutes.POST("/auth/line/complete", r.AuthHandler.LineComplete())
	}

	SettingRoutes := router.Group("/setting")
	{
		SettingRoutes.GET("/auth", r.SettingHandler.GetAuthConfigHandler())
	}

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(r.JwtMiddleware.AuthMiddleware())
	{
		protectedRoutes.GET("/auth/check", r.JwtMiddleware.JwtCheckMiddleware())

		protectedRoutes.GET("/users/my_page", r.UserHandler.MyPageHandler())
		protectedRoutes.GET("/setting/user_config", r.SettingHandler.GetUserConfigHandler())
		protectedRoutes.POST("/setting/user_config", r.SettingHandler.SaveUserConfigHandler())
		protectedRoutes.GET("/setting/root_config", r.SettingHandler.GetRootConfigHandler())
		protectedRoutes.POST("/setting/root_config", r.SettingHandler.SaveRootConfigHandler())
		protectedRoutes.GET("/words", r.WordHandler.ListHandler())
		protectedRoutes.GET("/words/:id", r.WordHandler.ShowHandler())
		protectedRoutes.POST("/words/register", r.WordHandler.RegisterHandler())
		protectedRoutes.POST("/words/memo", r.WordHandler.SaveMemoHandler())

		protectedRoutes.POST("/words/new", r.WordHandler.CreateHandler())
		protectedRoutes.PUT("/words/:id", r.WordHandler.UpdateHandler())
		protectedRoutes.DELETE("/words/:id", r.WordHandler.DeleteHandler())
		protectedRoutes.POST("/words/bulk_tokenize", r.WordHandler.BulkTokenizeHandler())
		protectedRoutes.POST("/words/bulk_register", r.WordHandler.BulkRegisterHandler())

		protectedRoutes.POST("/quizzes/new", r.QuizHandler.CreateHandler())
		protectedRoutes.POST("/quizzes/answers/:id", r.QuizHandler.PostAnswerAndRouteHandler())
		protectedRoutes.GET("/quizzes", r.QuizHandler.GetHandler())

		protectedRoutes.GET("/results", r.ResultHandler.GetIndexHandler())
		protectedRoutes.GET("/results/:quizNo", r.ResultHandler.GetHandler())
	}
}

// カスタム用cors
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

// 		// OPTIONSリクエスト（プリフライトリクエスト）の処理
// 		if c.Request.Method == "OPTIONS" {
// 			logrus.WithFields(logrus.Fields{
// 				"method": c.Request.Method,
// 				"path":   c.Request.URL.Path,
// 			}).Info("Handling CORS preflight request")
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

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
