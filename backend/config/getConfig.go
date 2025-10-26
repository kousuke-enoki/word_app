// Package config centralizes application configuration loaded from
// environment variables and provides a single struct (Config) that is
// easy to pass through dependency injection.
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

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

// リミット設定
// 登録単語数の上限など
type LimitsCfg struct {
	RegisteredWordsPerUser int // 200
	QuizMaxPerDay          int // 20
	QuizMaxQuestions       int // 100
	BulkMaxPerDay          int // 5
	BulkMaxBytes           int // 51200 (=50KB)
	BulkTokenizeMaxTokens  int // 200
	BulkRegisterMaxItems   int // 200
}

// Config aggregates all sub-config sections used across the application.
// It is designed to be constructed once in main and passed into modules.
type Config struct {
	App    AppCfg
	JWT    JWTCfg
	DB     DBCfg
	Line   LineOAuthCfg
	Lambda LambdaCfg
	Limits LimitsCfg
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
	jwtSecret := getenv("JWT_SECRET", "")
	lineID := getenv("LINE_CLIENT_ID", "")
	lineSec := getenv("LINE_CLIENT_SECRET", "")
	lineRedirect := getenv("LINE_REDIRECT_URI", "")

	// 3) Secrets Manager から読み出し（存在すれば）
	// 未設定の時だけ Secrets Manager を使いたい場合はフォールバック
	if (jwtSecret == "" || lineID == "" || lineSec == "" || lineRedirect == "") && os.Getenv("APP_SECRET_ARN") != "" {
		// 必須の一部がない時だけ取りに行く
		s, err := fetchSecretJSON[appSecret](context.Background(), os.Getenv("APP_SECRET_ARN"))
		if err != nil {
			logrus.Fatalf("read APP_SECRET_ARN: %v", err)
		}
		if jwtSecret == "" {
			jwtSecret = s.JWTSecret
		}
		if lineID == "" {
			lineID = s.LineClientID
		}
		if lineSec == "" {
			lineSec = s.LineClientSecret
		}
		if lineRedirect == "" {
			lineRedirect = s.LineRedirectURI
		}
	}

	// registeredWord の登録可能数/ユーザー
	registeredWordsPerUser := getenvInt("LIMIT_REGISTERED_WORDS_PER_USER", 200)
	// quiz_create を使用できる回数の上限/日
	quizMaxPerDay := getenvInt("LIMIT_QUIZ_MAX_PER_DAY", 20)
	// quiz_create 作成時に一度に作成できる質問数(quiz_questionの上限)
	quizMaxQuestions := getenvInt("LIMIT_QUIZ_MAX_QUESTIONS", 100)
	// bulk_tokenize を使用できる回数の上限/日
	bulkMaxPerDay := getenvInt("LIMIT_BULK_MAX_PER_DAY", 5)
	// bulk_tokenize で一度に処理できるデータサイズの上限(byte)
	bulkMaxBytes := getenvInt("LIMIT_BULK_MAX_BYTES", 50*1024) // 51200
	// bulk_tokenize で一度に処理できるトークン数の上限
	bulkTokenizeMaxTokens := getenvInt("LIMIT_BULK_TOKENIZE_MAX_TOKENS", 50*1024) // 51200
	// bulk_register で一度に登録できる単語数の上限
	bulkRegisterMaxItems := getenvInt("LIMIT_BULK_REGISTER_MAX_ITEMS", 50*1024) // 51200

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
		Limits: LimitsCfg{
			RegisteredWordsPerUser: registeredWordsPerUser,
			QuizMaxPerDay:          quizMaxPerDay,
			QuizMaxQuestions:       quizMaxQuestions,
			BulkMaxPerDay:          bulkMaxPerDay,
			BulkMaxBytes:           bulkMaxBytes, // 51200
			BulkTokenizeMaxTokens:  bulkTokenizeMaxTokens,
			BulkRegisterMaxItems:   bulkRegisterMaxItems,
		},
	}
}

/*──────────────── helpers ────────────────*/

// getenv returns the value of an environment variable, or the provided
// default if the variable is unset or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		logrus.Warnf("invalid int for %s=%q, using default=%d", key, v, def)
		return def
	}
	return i
}

// // "51200" / "50KB" / "1MB" を許容（大文字小文字OK）
// func getenvSizeBytes(key string, def int) int {
// 	v := strings.TrimSpace(os.Getenv(key))
// 	if v == "" {
// 		return def
// 	}
// 	n, ok := parseSizeToBytes(v)
// 	if !ok {
// 		logrus.Warnf("invalid size for %s=%q, using default=%d", key, v, def)
// 		return def
// 	}
// 	return n
// }

// var sizeRe = regexp.MustCompile(`(?i)^\s*(\d+)\s*([km]?b)?\s*$`)

// // 返り値: (bytes, ok)
// func parseSizeToBytes(s string) (int, bool) {
// 	m := sizeRe.FindStringSubmatch(s)
// 	if m == nil {
// 		// 数値だけの可能性にも対応
// 		if i, err := strconv.Atoi(s); err == nil {
// 			return i, true
// 		}
// 		return 0, false
// 	}
// 	numStr := m[1]
// 	unit := strings.ToUpper(strings.TrimSpace(m[2])) // "", "KB", "MB", "B"
// 	n, err := strconv.Atoi(numStr)
// 	if err != nil || n < 0 {
// 		return 0, false
// 	}
// 	switch unit {
// 	case "", "B":
// 		return n, true
// 	case "KB":
// 		return n * 1024, true
// 	case "MB":
// 		return n * 1024 * 1024, true
// 	default:
// 		return 0, false
// 	}
// }

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
