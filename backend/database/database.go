package database

import (
	"fmt"
	"os"

	"word_app/backend/ent"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

var entClient *ent.Client

func InitEntClient() {
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

	// スキーマのマイグレーションが必要ならここで実行
	// err = client.Schema.Create(context.Background())
	// if err != nil {
	//     log.Fatalf("failed creating schema resources: %v", err)
	// }

	entClient = client
	fmt.Println("Ent client initialized successfully")
}

func GetEntClient() *ent.Client {
	// InitEntClient() で作成したClientを返す
	return entClient
}

func SetEntClient(c *ent.Client) {
	entClient = c
}
