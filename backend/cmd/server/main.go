package main

import (
	"context"
	"log"
	"os"

	"word_app/backend/config"
	"word_app/backend/database"
	"word_app/backend/ent/user"
	"word_app/backend/logger"
	routerConfig "word_app/backend/router"
	"word_app/backend/seeder"
	AuthHandler "word_app/backend/src/handlers/auth"
	"word_app/backend/src/handlers/quiz"
	"word_app/backend/src/handlers/result"
	settingHandler "word_app/backend/src/handlers/setting"
	userHandler "word_app/backend/src/handlers/user"
	"word_app/backend/src/handlers/word"
	"word_app/backend/src/infrastructure"
	"word_app/backend/src/infrastructure/auth/line"
	"word_app/backend/src/infrastructure/jwt"
	authRepository "word_app/backend/src/infrastructure/repository/auth"
	userRepository "word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/interfaces"
	JwtMiddlewarePackage "word_app/backend/src/middleware/jwt"
	quizService "word_app/backend/src/service/quiz"
	resultService "word_app/backend/src/service/result"
	settingService "word_app/backend/src/service/setting"
	userService "word_app/backend/src/service/user"
	wordService "word_app/backend/src/service/word"
	"word_app/backend/src/usecase/auth"
	"word_app/backend/src/utils/tempjwt"
	"word_app/backend/src/validators"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// サーバーの初期化
	initializeServer()
}

// サーバーの初期化関数
func initializeServer() {
	defer func() {
		if p := recover(); p != nil {
			logrus.Fatalf("PANIC caught in main: %v\n", p)
		}
	}()
	config.LoadEnv()

	config.ConfigureGinMode()
	logger.InitLogger()

	appEnv, appPort, corsOrigin := config.LoadAppConfig()
	database.InitEntClient()
	entClient := database.GetEntClient()
	// entClient := connectToDatabase()
	defer entClient.Close()

	client := infrastructure.NewAppClient(entClient)
	setupDatabase(client)

	router := setupRouter(client, corsOrigin)

	startServer(router, appPort, appEnv)
}

// データベースのセットアップ
func setupDatabase(client interfaces.ClientInterface) {
	ctx := context.Background()
	entClient := client.EntClient()
	if entClient == nil {
		logrus.Fatalf("ent.Client is nil")
	}
	// Schema を作成
	if err := entClient.Schema.Create(ctx); err != nil {
		logrus.Fatalf("Failed to create schema: %v", err)
	}

	// Admin の存在を確認
	adminExists, err := entClient.User.Query().Where(user.Email("root@example.com")).Exist(ctx)
	if err != nil {
		logrus.Fatalf("Failed to check admin existence: %v", err)
	}

	// Seeder の実行
	if !adminExists {
		logrus.Info("Running initial seeder...")
		seeder.RunSeeder(ctx, client)
		logrus.Info("Seeder completed.")
	} else {
		logrus.Info("Seed data already exists, skipping.")
	}
}

// ルータのセットアップ
func setupRouter(client interfaces.ClientInterface, corsOrigin string) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Handler の初期化
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET environment variable is required")
	}
	jwtGen := jwt.NewMyJWTGenerator(jwtSecret)
	entUserClient := userService.NewEntUserClient(client)
	entSettingClient := settingService.NewEntSettingClient(client)
	wordClient := wordService.NewWordService(client)
	quizClient := quizService.NewQuizService(client)
	resultClient := resultService.NewResultService(client)
	authClient := jwt.NewJWTValidator(jwtSecret, client)

	userHandler := userHandler.NewUserHandler(entUserClient, jwtGen)
	settingHandler := settingHandler.NewSettingHandler(entSettingClient)

	wordHandler := word.NewWordHandler(wordClient)
	quizHandler := quiz.NewQuizHandler(quizClient)
	resultHandler := result.NewResultHandler(resultClient)
	JwtMiddleware := JwtMiddlewarePackage.NewJwtMiddleware(authClient)
	lineCfg := config.LoadLineConfig()
	lineProvider, err := line.NewProvider(lineCfg)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := userRepository.NewEntUserRepo(client)
	extAuthRepo := authRepository.NewEntExtAuthRepo(client)

	tempJwt := tempjwt.TempJWTNew(os.Getenv("TEMP_JWT_SECRET"))
	authUC := auth.NewAuthUsecase(
		lineProvider,
		userRepo,
		extAuthRepo,
		jwtGen,
		tempJwt,
	)
	authHandler := AuthHandler.NewAuthHandler(authUC, jwtGen)

	routerImpl := routerConfig.NewRouter(JwtMiddleware, authHandler, userHandler, settingHandler, wordHandler, quizHandler, resultHandler)

	routerImpl.SetupRouter(router)
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}

	validators.Init()
	binding.Validator = &validators.GinValidator{Validate: validators.V}

	return router
}

// サーバーを起動
func startServer(router *gin.Engine, port, env string) {
	logrus.Infof("Starting server on port %s in %s environment", port, env)
	if err := router.Run(":" + port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
