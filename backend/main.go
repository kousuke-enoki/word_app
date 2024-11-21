package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/seeder"
	"word_app/backend/src"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	log.Println("Starting server...")

	// 環境変数をロード
	loadEnv()

	// 環境変数の読み取り
	appEnv := os.Getenv("APP_ENV")
	appPort := os.Getenv("APP_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	corsOrigin := os.Getenv("CORS_ORIGIN")

	log.Printf("Environment: %s, Port: %s\n", appEnv, appPort)

	// Ginのモードを設定
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// PostgreSQLに接続
	dsn := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbUser, dbName, dbPassword)
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// コンテキストの作成
	ctx := context.Background()

	// マイグレーションを実行
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed creating schema resources: %v", err)
	}

	// 初回のみシードを実行
	seedAdminExists, err := client.User.Query().Where(user.Email("admin@example.com")).Exist(ctx)
	if err != nil {
		log.Fatalf("Failed checking for admin existence: %v", err)
	}

	if !seedAdminExists {
		log.Println("Running initial seeder...")
		seeder.RunSeeder(ctx, client)
		log.Println("Seeder completed.")
	} else {
		log.Println("Seed data already exists, skipping.")
	}

	// ロガーを設定
	configureLogger()
	setupFileLogger()

	// Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.New()
	setupFileLogger()

	// リクエストログの有効化/無効化
	configureRouterLogging(router)
	router.Use(gin.Logger(), gin.Recovery())

	// CORSの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ルータのセットアップ
	src.SetupRouter(router, client)
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// サーバー起動
	if err := router.Run(":" + appPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	log.Printf("Server successfully started on port %s, environment: %s\n", appPort, appEnv)
}

func loadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	log.Printf("APP_ENV is set to: %s", env)

	envFile := ".env." + env
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("No %s file found, falling back to system environment variables", envFile)
	} else {
		log.Printf("Loaded environment file: %s", envFile)
	}
}

func configureLogger() {
	// 環境変数からログレベルを取得
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // デフォルトは info
	}

	// ログレベルを設定
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Fatalf("Invalid LOG_LEVEL: %s", logLevel)
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Infof("Log level set to %s", logLevel)
}

func configureRouterLogging(router *gin.Engine) {
	enableLogging := os.Getenv("ENABLE_REQUEST_LOGGING")
	if enableLogging == "true" {
		router.Use(gin.Logger()) // ログミドルウェアを有効化
		logrus.Info("Request logging is enabled.")
	} else {
		logrus.Info("Request logging is disabled.")
	}
}

func setupFileLogger() {
	logPath := "log/app.log"

	// ディレクトリがなければ作成
	logDir := filepath.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// ログのローテーションを設定
	logFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     30, // 日
		Compress:   true,
	}

	// ログ出力をファイルに変更
	logrus.SetOutput(logFile)
	logrus.Infof("Logging to file: %s", logPath)
}
