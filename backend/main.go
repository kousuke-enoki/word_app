package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/seeder"
	"word_app/backend/src"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	// Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// CORSの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
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
