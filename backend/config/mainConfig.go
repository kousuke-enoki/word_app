package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadEnv() {
	if inLambda() {
		// Lambda では .env を読まない（環境変数だけ使う）
		return
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No %s file found, using system environment variables", envFile)
	} else {
		log.Printf("Loaded environment file: %s", envFile)
	}
}

func inLambda() bool { return os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" }

func ConfigureGinMode() {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug" // デフォルトは debug
	}
	gin.SetMode(ginMode)
	log.Printf("Gin mode set to %s", ginMode)
}

// アプリケーションの設定をロード
func LoadAppConfig() (string, string, string) {
	appEnv := os.Getenv("APP_ENV")
	appPort := os.Getenv("APP_PORT")
	corsOrigin := os.Getenv("CORS_ORIGIN")

	if appPort == "" {
		logrus.Fatal("APP_PORT is not set")
	}

	logrus.Infof("Environment: %s, Port: %s", appEnv, appPort)
	return appEnv, appPort, corsOrigin
}
