// Package logger centralizes application-wide logging setup.
// It configures logrus with a level derived from the LOG_LEVEL environment
// variable and a rotating file sink powered by lumberjack.
// Typical usage is to call InitLogger() once during application startup.
package logger

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger initializes logrus for the application.
// It reads LOG_LEVEL (e.g. "debug", "info", "warn", "error") and falls back
// to "info" when unset. The function configures a human-readable TextFormatter
// with full timestamps and routes output to a size-rotated file sink.
// It terminates the process if LOG_LEVEL is invalid.
//
// Note: By default this replaces logrus' output with a file writer.
// If you want to log to both stdout and a file, consider using io.MultiWriter.
func InitLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log.Println("logLevel:", logLevel)
	logrus.Info("logLevel:", logLevel)
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Fatalf("Invalid LOG_LEVEL: %s", logLevel)
	}

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

	// setupFileLogger configures a lumberjack rotating file writer and directs
	// logrus output to it. Logs are written to "log/app.log". The directory is
	// created when it does not exist. Rotation rules are:
	//   - MaxSize:    3 MB per file
	//   - MaxBackups: keep at most 3 old files
	//   - MaxAge:     30 days
	//   - Compress:   gzip-compress old files
	//
	// Consider making the path and rotation policy configurable via environment
	// variables if you need different behavior per environment.
	logFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    3,  // megabytes
		MaxBackups: 3,  // number of files
		MaxAge:     30, // days
		Compress:   true,
	}

	logrus.SetOutput(logFile)
	logrus.Infof("Logging to file: %s", logPath)
}
