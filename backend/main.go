package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"context"
  	"log"
	"eng_app/src"
	"eng_app/ent"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("start server...")
	// PostgreSQLに接続
	client, err := ent.Open("postgres", "host=db port=5432 user=postgres dbname=db password=password sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close() // データベース接続を閉じる

	// マイグレーションを実行
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.Default()

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
		log.Fatalf("failed to run server: %v", err)
	}
}
