package main

import (
	"context"
	"log"
	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/seeder"
	"word_app/backend/src"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting server...")

	// PostgreSQLに接続
	client, err := ent.Open("postgres", "host=db port=5432 user=postgres dbname=db password=password sslmode=disable")
	if err != nil {
		log.Fatalf("Failed opening connection to postgres: %v", err)
	}
	defer client.Close() // データベース接続を閉じる

	// コンテキストの作成
	ctx := context.Background()

	// マイグレーションを実行
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed creating schema resources: %v", err)
	}

	// 初回のみシードを実行
	seedAdminExists, err := client.User.Query().Where(user.Email("admin@example.com")).Exist(ctx)
	if err != nil {
		log.Fatalf("Failed checking for admin existence: %v", err)
	}

	if !seedAdminExists {
		log.Println("Running initial seeder...")
		seeder.RunSeeder(ctx, client) // Seederを呼び出す
		log.Println("Seeder completed.")
	} else {
		log.Println("Seed data already exists, skipping.")
	}

	// Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	// 現在はデバッグモードなので本番では下記
	// gin.SetMode(gin.ReleaseMode)

	// CORSの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Authorization ヘッダーを許可
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ルータのセットアップ
	r := gin.Default()
	userHandler := user.SignUpHandler(client)
	wordHandler := handlers.NewWordHandler(client)

	src.SetupRouter(r, client, userHandler, wordHandler)

	// サーバー起動
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	} else {
		log.Println("Server successfully started on port 8080")
	}
}
