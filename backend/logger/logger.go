// Package logger centralizes application-wide logging setup.
// It configures logrus with a level derived from the LOG_LEVEL environment
// variable and a rotating file sink powered by lumberjack.
// Typical usage is to call InitLogger() once during application startup.
package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	Level        string // "debug"|"info"|...
	Format       string // "json"|"text"
	Stdout       bool
	FilePath     string // emptyならファイル出力なし
	ReportCaller bool
	MaxSizeMB    int
	MaxBackups   int
	MaxAgeDays   int
	Compress     bool
}

// InitLogger initializes logrus for the application.
// It reads LOG_LEVEL (e.g. "debug", "info", "warn", "error") and falls back
// to "info" when unset. The function configures a human-readable TextFormatter
// with full timestamps and routes output to a size-rotated file sink.
// It terminates the process if LOG_LEVEL is invalid.
//
// Note: By default this replaces logrus' output with a file writer.
// If you want to log to both stdout and a file, consider using io.MultiWriter.
func InitLogger() {
	opt := readOptionsFromEnv()

	// level
	level, err := logrus.ParseLevel(strings.ToLower(opt.Level))
	if err != nil {
		log.Fatalf("Invalid LOG_LEVEL: %s", opt.Level)
	}
	logrus.SetLevel(level)

	// format
	switch strings.ToLower(opt.Format) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// outputs
	var writers []io.Writer
	if opt.Stdout {
		writers = append(writers, os.Stdout)
	}
	if opt.FilePath != "" {
		ensureDir(filepath.Dir(opt.FilePath))
		writers = append(writers, &lumberjack.Logger{
			Filename:   opt.FilePath,
			MaxSize:    max1(opt.MaxSizeMB, 3),
			MaxBackups: max1(opt.MaxBackups, 3),
			MaxAge:     max1(opt.MaxAgeDays, 30),
			Compress:   opt.Compress,
		})
	}
	if len(writers) == 0 {
		// デフォルトはstdout
		writers = []io.Writer{os.Stdout}
	}
	logrus.SetOutput(io.MultiWriter(writers...))
	logrus.SetReportCaller(opt.ReportCaller)

	logrus.WithFields(logrus.Fields{
		"level":  level.String(),
		"format": opt.Format,
		"stdout": opt.Stdout,
		"file":   opt.FilePath,
	}).Info("logger initialized")
}

// lambdaでの動作かどうか
func inLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" || os.Getenv("LAMBDA_TASK_ROOT") != ""
}

func readOptionsFromEnv() Options {
	// Lambdaならstdout+json固定（CloudWatchへ）
	if inLambda() {
		return Options{
			Level:        envOr("LOG_LEVEL", "info"),
			Format:       "json",
			Stdout:       true,
			FilePath:     "",
			ReportCaller: envBool("LOG_REPORT_CALLER", false),
		}
	}
	return Options{
		Level:        envOr("LOG_LEVEL", "info"),
		Format:       envOr("LOG_FORMAT", "json"),
		Stdout:       envBool("LOG_STDOUT", true),      // ローカルでも見えるように
		FilePath:     envOr("LOG_FILE", "log/app.log"), // 空でファイル出力OFF
		ReportCaller: envBool("LOG_REPORT_CALLER", false),
		MaxSizeMB:    envInt("LOG_ROTATE_SIZE_MB", 1),
		MaxBackups:   envInt("LOG_ROTATE_BACKUPS", 5),
		MaxAgeDays:   envInt("LOG_ROTATE_MAX_DAYS", 30),
		Compress:     envBool("LOG_ROTATE_COMPRESS", true),
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envBool(k string, def bool) bool {
	if v := os.Getenv(k); v != "" {
		b, _ := strconv.ParseBool(v)
		return b
	}
	return def
}

func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func ensureDir(dir string) {
	if dir == "" {
		return
	}
	_ = os.MkdirAll(dir, 0o755)
}

func max1(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}
