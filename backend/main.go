package main

import (
	"context"
	"eng_app/ent"
	"eng_app/src"
	"log"

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

	// マイグレーションを実行
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed creating schema resources: %v", err)
	}

	// Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	// 現在はデバッグモードなので本番では下記
	// gin.SetMode(gin.ReleaseMode)

	// CORSの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	src.SetupRouter(router, client)

	// サーバー起動
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	} else {
		log.Println("Server successfully started on port 8080")
	}
}
