package logger

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// ロガーを初期化
func InitLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log.Println(logLevel)
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Fatalf("Invalid LOG_LEVEL: %s", logLevel)
	}
	log.Println(level)

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	setupFileLogger()
	logrus.Info("Logger initialized successfully ", logLevel)
}

// ログファイルの設定
func setupFileLogger() {
	logPath := "log/app.log"
	logDir := filepath.Dir(logPath)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Fatalf("Failed to create log directory: %v", err)
		}
	}

	logFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	logrus.SetOutput(logFile)
	logrus.Infof("Logging to file: %s", logPath)
}
