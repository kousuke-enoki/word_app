package main

import (
	"context"
	"fmt"
	"os"

	"word_app/backend/config"
	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/logger"
	routerConfig "word_app/backend/router"
	"word_app/backend/seeder"
	auth "word_app/backend/src/handlers/middleware"
	userHandler "word_app/backend/src/handlers/user"
	"word_app/backend/src/handlers/word"
	"word_app/backend/src/infrastructure"
	"word_app/backend/src/interfaces"
	userService "word_app/backend/src/service/user"
	wordService "word_app/backend/src/service/word"
	"word_app/backend/src/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// サーバーの初期化
	initializeServer()
}

// サーバーの初期化関数
func initializeServer() {
	config.LoadEnv()

	config.ConfigureGinMode()
	logger.InitLogger()

	appEnv, appPort, corsOrigin := config.LoadAppConfig()
	entClient := connectToDatabase()
	defer entClient.Close()

	client := infrastructure.NewAppClient(entClient)
	setupDatabase(client)

	router := setupRouter(client, corsOrigin)

	startServer(router, appPort, appEnv)
}

// データベースに接続
func connectToDatabase() *ent.Client {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	// entDebug := os.Getenv("ENT_DEBUG")

	dsn := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable", dbHost, dbUser, dbName, dbPassword)
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// if entDebug == "1" {
	// 	client = client.Debug()
	// 	logrus.Info("Ent debug mode enabled")
	// }

	logrus.Info("Database connection established")
	return client
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
	jwtGenerator := utils.NewMyJWTGenerator(jwtSecret)
	entClient := userService.NewEntUserClient(client)
	wordClient := wordService.NewWordService(client)
	userHandler := userHandler.NewUserHandler(entClient, jwtGenerator)

	wordHandler := word.NewWordHandler(wordClient)
	authHandler := auth.NewAuthHandler()

	routerImpl := routerConfig.NewRouter(authHandler, userHandler, wordHandler)
	routerImpl.SetupRouter(router)
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}
	logrus.Info("Router setup completed")

	return router
}

// サーバーを起動
func startServer(router *gin.Engine, port, env string) {
	logrus.Infof("Starting server on port %s in %s environment", port, env)
	if err := router.Run(":" + port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
