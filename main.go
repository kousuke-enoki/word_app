package main

import (
	"github.com/gin-gonic/gin"
	"context"
  	"log"
	"eng_app/src"
	"eng_app/ent"
	_ "github.com/lib/pq"
)

func main() {

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
	src.SetupRouter(router, client)

	// サーバー起動
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
