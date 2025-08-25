package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"word_app/backend/ent"

	entsql "entgo.io/ent/dialect/sql"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/sirupsen/logrus"
)

var entClient *ent.Client

// rds用オブジェクト
type rdsSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Dbname   string `json:"dbname"`
	Engine   string `json:"engine"`
}

func InitEntClient() error {
	return initEntClientWithContext(context.Background())
}

func initEntClientWithContext(ctx context.Context) error {
	cfg, err := loadDbConfig(ctx)
	if err != nil {
		logrus.WithError(err).Error("loadDbConfig failed")
		entClient = nil
		return err
	}

	sslMode := getenv("DB_SSLMODE", "")
	if sslMode == "" {
		if isLambda() {
			sslMode = "require" // 本番(Lambda)のデフォルト
		} else {
			sslMode = "disable" // ローカルのデフォルト
		}
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, sslMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.WithError(err).Error("sql.Open failed")
		entClient = nil
		return err
	}

	// Lambda は同時接続を絞るのが吉
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ent のドライバに載せ替えて Client を作成
	drv := entsql.OpenDB(dialect.Postgres, db)
	c := ent.NewClient(ent.Driver(drv))
	entClient = c
	logrus.WithFields(logrus.Fields{
		"host": cfg.Host, "port": cfg.Port, "db": cfg.Name, "user": cfg.User,
	}).Info("Ent client initialized")
	return nil
}

type dbCfg struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

func isLambda() bool {
	return os.Getenv("AWS_LAMBDA_RUNTIME_API") != ""
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func loadDbConfig(ctx context.Context) (*dbCfg, error) {
	// 1) RDS のシークレット優先（CDK で DB_SECRET_ARN を渡している想定）
	if arn := os.Getenv("DB_SECRET_ARN"); arn != "" {
		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("load aws cfg: %w", err)
		}
		sm := secretsmanager.NewFromConfig(awsCfg)
		out, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &arn})
		if err != nil {
			return nil, fmt.Errorf("get secret: %w", err)
		}
		if out.SecretString == nil {
			return nil, fmt.Errorf("secret has no SecretString")
		}
		var s rdsSecret
		if err := json.Unmarshal([]byte(*out.SecretString), &s); err != nil {
			return nil, fmt.Errorf("unmarshal secret: %w", err)
		}
		// env で上書きも許可（CDK から渡しているならそれを優先したい時に使える）
		name := s.Dbname
		if v := os.Getenv("DB_NAME"); v != "" {
			name = v
		}
		port := s.Port
		if v := os.Getenv("DB_PORT"); v != "" {
			if p, err := strconv.Atoi(v); err == nil {
				port = p
			}
		}
		return &dbCfg{
			Host: firstNonEmpty(os.Getenv("DB_HOST"), s.Host),
			Port: port,
			User: firstNonEmpty(os.Getenv("DB_USER"), s.Username),
			Pass: firstNonEmpty(os.Getenv("DB_PASSWORD"), s.Password),
			Name: firstNonEmpty(name, "postgres"),
		}, nil
	}

	// 2) フォールバック：環境変数から全部読む
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	if host == "" || user == "" || pass == "" || name == "" {
		return nil, fmt.Errorf("DB envs missing (need DB_HOST/DB_USER/DB_PASSWORD/DB_NAME or DB_SECRET_ARN)")
	}
	port := 5432
	if v := os.Getenv("DB_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}
	return &dbCfg{Host: host, Port: port, User: user, Pass: pass, Name: name}, nil
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func GetEntClient() *ent.Client {
	// InitEntClient() で作成したClientを返す
	return entClient
}

func SetEntClient(c *ent.Client) {
	entClient = c
}
