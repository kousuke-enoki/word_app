// Package config centralizes application configuration loaded from
// environment variables and provides a single struct (Config) that is
// easy to pass through dependency injection.
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"

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
	App    AppCfg
	JWT    JWTCfg
	DB     DBCfg
	Line   LineOAuthCfg
	Lambda LambdaCfg
}

type dbSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type appSecret struct {
	JWTSecret        string `json:"JWT_SECRET"`
	LineClientID     string `json:"LINE_CLIENT_ID"`
	LineClientSecret string `json:"LINE_CLIENT_SECRET"`
	LineRedirectURI  string `json:"LINE_REDIRECT_URI"`
}

type LambdaCfg struct {
	LambdaRuntime string `json:"LambdaRuntime"`
}

// NewConfig reads environment variables, applies sane defaults, and returns
// a Config. It terminates the process (logrus.Fatalf) if required variables
// are missing. If you need testability without process exit, consider a
// constructor that returns (Config, error) instead.
// NOTE: TempSecret currently reads JWT_SECRET as well. If you intend a
// separate key for temporary tokens, change it to read TEMP_JWT_SECRET.
func NewConfig() *Config {
	// 1) 非秘匿は環境変数
	appEnv := getenv("APP_ENV", "production")
	appPort := getenv("APP_PORT", "8080") // Lambda では未使用でもOK

	// // 2) DB ホスト/ポート/DB名（必須）
	// dbHost := must("DB_HOST")
	// dbPort := getenv("DB_PORT", "5432")
	// dbName := must("DB_NAME")
	lambdaRuntime := getenv("AWS_LAMBDA_RUNTIME_API", "")

	// 3) Secrets Manager から読み出し（存在すれば）
	var jwtSecret, lineID, lineSec, lineRedirect string
	if arn := os.Getenv("APP_SECRET_ARN"); arn != "" {
		if s, err := fetchSecretJSON[appSecret](context.Background(), arn); err != nil {
			logrus.Fatalf("read APP_SECRET_ARN: %v", err)
		} else {
			jwtSecret = s.JWTSecret
			lineID = s.LineClientID
			lineSec = s.LineClientSecret
			lineRedirect = s.LineRedirectURI
		}
	} else {
		// Secrets を使わない運用なら従来どおり env から
		jwtSecret = must("JWT_SECRET")
		lineID = must("LINE_CLIENT_ID")
		lineSec = must("LINE_CLIENT_SECRET")
		lineRedirect = must("LINE_REDIRECT_URI")
	}

	// // 4) DB 認証（ユーザー/パス）は Secrets Manager
	// // var dbUser, dbPass string
	// if arn := os.Getenv("DB_SECRET_ARN"); arn != "" {
	// 	s, err := fetchSecretJSON[dbSecret](context.Background(), arn)
	// 	if err != nil {
	// 		logrus.Fatalf("read DB_SECRET_ARN: %v", err)
	// 	}
	// 	dbUser, dbPass = s.Username, s.Password
	// } else {
	// 	// フォールバック（テスト用）
	// 	dbUser = getenv("DB_USER", "postgres")
	// 	dbPass = getenv("DB_PASSWORD", "password")
	// }

	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	// 	url.QueryEscape(dbUser), url.QueryEscape(dbPass),
	// 	dbHost, dbPort, dbName,
	// )

	// ♦ 3. 構造体に詰めて返す
	return &Config{
		App: AppCfg{Env: appEnv, Port: appPort},
		JWT: JWTCfg{Secret: jwtSecret, TempSecret: jwtSecret, ExpireHour: 1, ExpireMinute: 30},
		// DB:  DBCfg{DSN: dsn},
		Line: LineOAuthCfg{
			ClientID: lineID, ClientSecret: lineSec, RedirectURI: lineRedirect,
		},
		Lambda: LambdaCfg{LambdaRuntime: lambdaRuntime},
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

func fetchSecretJSON[T any](ctx context.Context, arn string) (T, error) {
	var zero T

	cfg, err := awscfg.LoadDefaultConfig(ctx)
	// リージョンを強制したいなら:
	// cfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion("ap-northeast-1"))
	if err != nil {
		return zero, fmt.Errorf("load aws config: %w", err)
	}

	client := sm.NewFromConfig(cfg)

	out, err := client.GetSecretValue(ctx, &sm.GetSecretValueInput{
		SecretId: aws.String(arn),
	})
	if err != nil {
		return zero, fmt.Errorf("get secret value: %w", err)
	}
	if out.SecretString == nil {
		return zero, fmt.Errorf("secret %s has no SecretString (maybe binary)", arn)
	}

	var v T
	if err := json.Unmarshal([]byte(*out.SecretString), &v); err != nil {
		return zero, fmt.Errorf("unmarshal secret json: %w", err)
	}
	return v, nil
}
