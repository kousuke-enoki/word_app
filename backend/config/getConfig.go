// Package config centralizes application configuration loaded from
// environment variables and provides a single struct (Config) that is
// easy to pass through dependency injection.
package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// AppCfg holds runtime application settings.
// Env is the current environment name (e.g. "development", "production").
// Port is the TCP port the HTTP server listens on.
type AppCfg struct {
	Env  string
	Port string
}

// JWTCfg defines secrets and default expiration for JWTs issued by the app.
// Secret is the primary signing key.
// TempSecret is intended for short-lived/temporary flows (e.g. email link).
// ExpireHour and ExpireMinute together define the default token TTL.
type JWTCfg struct {
	Secret       string
	TempSecret   string
	ExpireHour   int // 例: 1
	ExpireMinute int // 例: 30
}

// DBCfg contains database connectivity settings.
// DSN is a full connection string consumable by the Ent client / driver.
type DBCfg struct {
	DSN string
}

// LineOAuthCfg holds LINE Login (OAuth) client configuration.
// ClientID and ClientSecret are the app credentials issued by LINE.
// RedirectURI must match the one registered with the provider.
type LineOAuthCfg struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// Config aggregates all sub-config sections used across the application.
// It is designed to be constructed once in main and passed into modules.
type Config struct {
	App  AppCfg
	JWT  JWTCfg
	DB   DBCfg
	Line LineOAuthCfg
}

// NewConfig reads environment variables, applies sane defaults, and returns
// a Config. It terminates the process (logrus.Fatalf) if required variables
// are missing. If you need testability without process exit, consider a
// constructor that returns (Config, error) instead.
// NOTE: TempSecret currently reads JWT_SECRET as well. If you intend a
// separate key for temporary tokens, change it to read TEMP_JWT_SECRET.
func NewConfig() *Config {
	// ♦ 1. 必須値を取得
	jwtSecret := must("JWT_SECRET")
	tempSecret := must("JWT_SECRET")
	// tempJwt := tempjwt.New(os.Getenv("TEMP_JWT_SECRET"))
	lineClientID := must("LINE_CLIENT_ID")
	lineClientSec := must("LINE_CLIENT_SECRET")
	lineRedirect := must("LINE_REDIRECT_URI")

	// ♦ 2. オプション値を取得（デフォルトあり）
	appEnv := getenv("APP_ENV", "development")
	appPort := getenv("APP_PORT", "8080")
	dbDSN := getenv("DB_DSN",
		"postgres://postgres:password@db:5432/db?sslmode=disable")

	// ♦ 3. 構造体に詰めて返す
	return &Config{
		App: AppCfg{
			Env:  appEnv,
			Port: appPort,
		},
		JWT: JWTCfg{
			Secret:       jwtSecret,
			TempSecret:   tempSecret,
			ExpireHour:   1, // 固定値でも OK。env にしても良い
			ExpireMinute: 30,
		},
		DB: DBCfg{
			DSN: dbDSN,
		},
		Line: LineOAuthCfg{
			ClientID:     lineClientID,
			ClientSecret: lineClientSec,
			RedirectURI:  lineRedirect,
		},
	}
}

/*──────────────── helpers ────────────────*/

// must fetches an environment variable and fatally exits if it is empty.
func must(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logrus.Fatalf("environment variable %s is required", key)
	}
	return val
}

// getenv returns the value of an environment variable, or the provided
// default if the variable is unset or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
