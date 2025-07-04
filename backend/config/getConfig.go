package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// ─── サブ構造体 ──────────────────────────────
type AppCfg struct {
	Env  string
	Port string
}

type TempJWT struct {
	secret []byte
}

type JWTCfg struct {
	Secret       string
	TempSecret   string
	ExpireHour   int // 例: 1
	ExpireMinute int // 例: 30
}

type DBCfg struct {
	DSN string
}

type LineOAuthCfg struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// ─── 集約構造体 ──────────────────────────────
type Config struct {
	App  AppCfg
	JWT  JWTCfg
	DB   DBCfg
	Line LineOAuthCfg
}

// ─── Public constructor ─────────────────────
func NewConfig() *Config {
	// ♦ 1. 必須値を取得
	jwtSecret := must("JWT_SECRET")
	tempSecret := must("JWT_SECRET")
	// tempJwt := tempjwt.TempJWTNew(os.Getenv("TEMP_JWT_SECRET"))
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

/*──────────────── ヘルパ ────────────────*/

// 必須変数: 空なら Fatal
func must(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logrus.Fatalf("environment variable %s is required", key)
	}
	return val
}

// 任意変数: 空なら default
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
