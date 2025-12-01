package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"word_app/backend/config"
	"word_app/backend/database"
	"word_app/backend/ent/migrate"
	"word_app/backend/ent/rootconfig"
	"word_app/backend/ent/user"
	"word_app/backend/ent/word"
	"word_app/backend/internal/di"
	"word_app/backend/internal/middleware"
	"word_app/backend/logger"
	routerConfig "word_app/backend/router"
	"word_app/backend/seeder"
	"word_app/backend/src/infrastructure"
	"word_app/backend/src/interfaces"
	"word_app/backend/src/interfaces/sqlexec"
	"word_app/backend/src/validators"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	once      sync.Once
	ginLambda *ginadapter.GinLambda // Lambda 用
)

func fastHealthResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"ok":true,"mode":"health-only","fast":true}`,
	}, nil
}

func main() {
	if isLambda() {
		// aws-lambda-go の Start に渡す"外側"で、即返すルートを作る
		handler := func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			// ① 非常停止パス：環境変数が HEALTH_ONLY なら /health は即返す
			if req.Path == "/health" || req.Resource == "/health" {
				return fastHealthResponse()
			}

			// ② （初回だけ）重い初期化（HEALTH_ONLY なら healthOnlyRouter に限定）
			once.Do(func() {
				if os.Getenv("APP_BOOTSTRAP_MODE") == "HEALTH_ONLY" {
					ginLambda = ginadapter.New(healthOnlyRouter())
				} else {
					router, _, _, _ := mustInitServer(false) // ← ここが重い
					ginLambda = ginadapter.New(router)
				}
			})

			// ③ 通常処理
			return ginLambda.ProxyWithContext(ctx, req)
		}

		lambda.Start(handler)
		return
	}

	// Local
	router, port, env, cleanup := mustInitServer(true)
	defer cleanup()
	startServer(router, port, env)
}

// ヘルスだけの薄いルータ（デプロイ確認用）
func healthOnlyRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	return r
}

func isLambda() bool {
	return os.Getenv("AWS_LAMBDA_RUNTIME_API") != ""
}

// サーバーの初期化関数
func mustInitServer(needCleanup bool) (*gin.Engine, string, string, func()) {
	defer func() {
		if p := recover(); p != nil {
			buf := make([]byte, 1<<16)
			n := runtime.Stack(buf, false)
			logrus.Fatalf("PANIC caught in main: %v\n--- STACK ---\n%s", p, string(buf[:n]))
		}
	}()

	if !isLambda() {
		config.LoadEnv()
	}

	config.ConfigureGinMode()
	logger.InitLogger()

	appEnv, appPort, corsOrigin := config.LoadAppConfig()
	err := database.InitEntClient()
	if err != nil {
		logrus.Fatalf("failed to init ent client: %v", err)
	}
	entClient := database.GetEntClient()
	sqlDB := database.GetSQLDB()

	// cleanup は呼ぶタイミングを呼び出し側に委譲
	cleanup := func() {}
	if needCleanup {
		cleanup = func() {
			if err := entClient.Close(); err != nil {
				logrus.Errorf("failed to close ent client: %v", err)
			}
		}
	}

	client := infrastructure.NewAppClient(entClient)
	runner := sqlexec.NewStdSQLRunner(sqlDB)

	runMig := shouldRun("RUN_MIGRATION")
	runSeed := shouldRun("RUN_SEEDER")

	if !isLambda() && os.Getenv("RUN_MIGRATION") == "" {
		// ローカルのデフォルト：Migrationは On、Seeder は Off
		runMig = true
	}

	if runMig {
		if err := runMigration(client); err != nil {
			logrus.Fatalf("migration failed: %v", err)
		}
	} else {
		logrus.Info("Skip migration on boot")
	}

	if runSeed {
		if err := runSeederIfNeeded(client); err != nil {
			logrus.Fatalf("seeder failed: %v", err)
		}
	} else {
		logrus.Info("Skip seeder on boot")
	}

	router := setupRouter(client, runner, corsOrigin, appEnv)
	return router, appPort, appEnv, cleanup
}

func shouldRun(key string) bool {
	v := strings.ToLower(os.Getenv(key))
	return v == "1" || v == "true" || v == "yes"
}

func runMigration(client interfaces.ClientInterface) error {
	ctx := context.Background()
	entClient := client.EntClient()
	if entClient == nil {
		return fmt.Errorf("ent.Client is nil")
	}

	// 既存テーブルに対してupdated_atカラムのマイグレーションを実行
	// テーブルが存在しない場合はスキップし、Entのマイグレーションに任せる
	if err := migrateRootConfigUpdatedAt(ctx); err != nil {
		return fmt.Errorf("failed to migrate root_configs.updated_at: %w", err)
	}

	// Entスキーマに基づくテーブル・インデックスを作成
	// 外部キーも明確に生成するようにする。
	if err := entClient.Schema.Create(ctx, migrate.WithForeignKeys(true)); err != nil {
		return err
	}

	// PostgreSQL固有のカスタムインデックスを作成
	if err := createCustomIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create custom indexes: %w", err)
	}

	return nil
}

// createCustomIndexes はPostgreSQL固有のインデックスを作成します。
// PostgreSQL以外のDB（例: sqlite）の場合は安全にスキップします。
func createCustomIndexes(ctx context.Context) error {
	sqlDB := database.GetSQLDB()
	if sqlDB == nil {
		logrus.Info("[migrate] SQLDB is nil, skip custom indexes")
		return nil
	}

	db := sqlDB

	// pg_trgm拡張を有効化（trigram検索用）
	if _, err := db.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS pg_trgm;
	`); err != nil {
		return fmt.Errorf("failed to create pg_trgm extension: %w", err)
	}

	// Ent管理のB-treeインデックス: registration_count
	// EntのSchema.Create()は既存テーブルに対して新しいインデックスを自動追加しないため、
	// 手動で作成する必要がある
	if _, err := db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS word_registration_count
		ON words(registration_count);
	`); err != nil {
		return fmt.Errorf("failed to create word_registration_count index: %w", err)
	}

	// words.name の trigram GIN インデックス
	// NameContains()によるLIKE検索を高速化するため
	if _, err := db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS words_name_trgm_idx
		ON words
		USING gin (name gin_trgm_ops);
	`); err != nil {
		return fmt.Errorf("failed to create words_name_trgm_idx: %w", err)
	}

	logrus.Info("[migrate] custom indexes created successfully")
	return nil
}

// migrateRootConfigUpdatedAt は既存のroot_configsテーブルに対してupdated_atカラムを追加・調整します。
// テーブルが存在しない場合は何もせずにスキップします（Entのマイグレーションで作成されるため）。
// この関数は既存のデータベースへの後方互換性のためのマイグレーション処理です。
func migrateRootConfigUpdatedAt(ctx context.Context) error {
	sqlDB := database.GetSQLDB()
	if sqlDB == nil {
		return nil // SQLDBが取得できない場合はスキップ
	}

	// テーブルの存在確認
	tableExists, err := checkTableExists(ctx, sqlDB, "root_configs")
	if err != nil {
		return fmt.Errorf("failed to check if root_configs table exists: %w", err)
	}

	if !tableExists {
		// テーブルが存在しない場合は、Entのマイグレーションに任せる
		logrus.Info("root_configs table does not exist, skipping updated_at migration (will be handled by Ent migration)")
		return nil
	}

	// カラムの存在確認
	columnExists, err := checkColumnExists(ctx, sqlDB, "root_configs", "updated_at")
	if err != nil {
		return fmt.Errorf("failed to check if updated_at column exists: %w", err)
	}

	if !columnExists {
		// テーブル存在を再確認（安全性のため）
		if exists, err := checkTableExists(ctx, sqlDB, "root_configs"); err != nil {
			return fmt.Errorf("failed to re-verify table existence before adding column: %w", err)
		} else if !exists {
			logrus.Warn("root_configs table disappeared during migration, skipping column addition")
			return nil
		}

		// カラムが存在しない場合は、NULLを許可して追加（後でEntがNOT NULLに変更する）
		_, err = sqlDB.ExecContext(ctx, `
			ALTER TABLE root_configs 
			ADD COLUMN updated_at TIMESTAMP
		`)
		if err != nil {
			return fmt.Errorf("failed to add updated_at column to root_configs table: %w", err)
		}
		logrus.Info("Added updated_at column to root_configs table")
	}

	// 既存レコードに対してupdated_atを設定（NULLの場合は現在時刻）
	// テーブル存在を再確認（安全性のため）
	if exists, err := checkTableExists(ctx, sqlDB, "root_configs"); err != nil {
		return fmt.Errorf("failed to verify table existence before updating records: %w", err)
	} else if !exists {
		logrus.Warn("root_configs table disappeared during migration, skipping record update")
		return nil
	}

	_, err = sqlDB.ExecContext(ctx, `
		UPDATE root_configs 
		SET updated_at = COALESCE(updated_at, NOW())
		WHERE updated_at IS NULL
	`)
	if err != nil {
		// テーブルが削除された場合など、他の理由でエラーになる可能性がある
		return fmt.Errorf("failed to update existing root_configs.updated_at values: %w", err)
	}
	logrus.Info("Updated existing root_configs.updated_at values")

	// カラムが存在するが、NOT NULL制約がない場合は、NOT NULL制約を追加
	// （EntのマイグレーションがNOT NULL制約を追加しようとする前に）
	// テーブルとカラムの存在を再確認（安全性のため）
	if exists, err := checkTableExists(ctx, sqlDB, "root_configs"); err != nil {
		return fmt.Errorf("failed to verify table existence before checking nullability: %w", err)
	} else if !exists {
		logrus.Warn("root_configs table disappeared during migration, skipping nullability check")
		return nil
	}

	if exists, err := checkColumnExists(ctx, sqlDB, "root_configs", "updated_at"); err != nil {
		return fmt.Errorf("failed to verify column existence before checking nullability: %w", err)
	} else if !exists {
		logrus.Warn("updated_at column disappeared during migration, skipping nullability check")
		return nil
	}

	var isNullable string
	err = sqlDB.QueryRowContext(ctx, `
		SELECT is_nullable
		FROM information_schema.columns 
		WHERE table_schema = 'public'
		AND table_name = 'root_configs' 
		AND column_name = 'updated_at'
	`).Scan(&isNullable)
	if err != nil {
		return fmt.Errorf("failed to check updated_at column nullability in root_configs table: %w", err)
	}

	if isNullable == "YES" {
		// デフォルト値を設定してからNOT NULL制約を追加
		_, err = sqlDB.ExecContext(ctx, `
			ALTER TABLE root_configs 
			ALTER COLUMN updated_at SET DEFAULT NOW(),
			ALTER COLUMN updated_at SET NOT NULL
		`)
		if err != nil {
			return fmt.Errorf("failed to set NOT NULL constraint on root_configs.updated_at: %w", err)
		}
		logrus.Info("Set NOT NULL constraint on root_configs.updated_at")
	}

	return nil
}

// checkTableExists は指定されたテーブルが存在するかどうかを確認します。
func checkTableExists(ctx context.Context, sqlDB *sql.DB, tableName string) (bool, error) {
	var exists bool
	err := sqlDB.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// checkColumnExists は指定されたテーブルのカラムが存在するかどうかを確認します。
func checkColumnExists(ctx context.Context, sqlDB *sql.DB, tableName, columnName string) (bool, error) {
	var exists bool
	err := sqlDB.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_schema = 'public'
			AND table_name = $1
			AND column_name = $2
		)
	`, tableName, columnName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func runSeederIfNeeded(client interfaces.ClientInterface) error {
	ctx := context.Background()
	entClient := client.EntClient()
	runSeedForWords := shouldRun("RUN_SEEDER_FOR_WORDS")
	if entClient == nil {
		return fmt.Errorf("ent.Client is nil")
	}
	adminExists, err := entClient.User.Query().
		Where(user.Email("root@example.com")).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("admin check failed by user: %w", err)
	}

	seedWordExists, err := entClient.Word.Query().
		Where(word.ID(1)).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("admin check failed by words: %w", err)
	}

	if !adminExists {
		logrus.Info("Running initial seeder...")
		seeder.RunSeeder(ctx, client)
		logrus.Info("Seeder completed.")
	} else {
		logrus.Info("Seed data already exists, skipping.")
		// ENABLE_TEST_USER_MODEがtrueの場合は、RootConfigを更新する
		if shouldRun("ENABLE_TEST_USER_MODE") {
			logrus.Info("ENABLE_TEST_USER_MODE is true, updating RootConfig...")
			seeder.SeedRootConfig(ctx, client)
		}
	}

	// RootConfigが存在しない場合は個別にシード
	rootConfigExists, err := entClient.RootConfig.Query().
		Where(rootconfig.ID(1)).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("root config check failed: %w", err)
	}
	if !rootConfigExists {
		logrus.Info("RootConfig not found, seeding...")
		seeder.SeedRootConfig(ctx, client)
		logrus.Info("RootConfig seeded.")
	}

	if runSeedForWords && !seedWordExists {
		logrus.Info("Running initial seeder for words...")
		seeder.SeedWords(ctx, client)
		logrus.Info("Seeder for words completed.")
	} else {
		logrus.Info("Seed words data already exists, skipping.")
	}
	return nil
}

func setupRouter(client interfaces.ClientInterface, runner sqlexec.Runner, corsOrigin string, appEnv string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// アクセスログミドルウェア（環境変数から設定を読み取り）
	opts := parseAccessLogOpts()
	router.Use(middleware.AccessLog(logrus.StandardLogger(), opts))

	// 1) CORS: 環境変数（カンマ区切り）→ スライス
	allowed := []string{}
	for _, o := range strings.Split(corsOrigin, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			allowed = append(allowed, o)
		}
	}

	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	// 2) 環境変数で指定されたオリジンと *.vercel.app の両方を許可
	// AllowOriginFunc が設定されている場合、AllowOrigins は無視されるため、
	// AllowOriginFunc 内で両方をチェックする
	cfg.AllowOriginFunc = func(origin string) bool {
		// 環境変数で指定されたオリジンをチェック
		for _, o := range allowed {
			if origin == o {
				return true
			}
		}
		// *.vercel.app で終わるオリジンも許可
		return strings.HasSuffix(origin, ".vercel.app")
	}

	router.Use(cors.New(cfg))

	// ルータのセットアップ
	cfgObj := config.NewConfig() // ← env 読み取りなど 1 箇所に集約
	repos := di.NewRepositories(client, runner)
	ucs, err := di.NewUseCases(cfgObj, repos)
	if err != nil {
		logrus.Fatal(err)
	}

	services := di.NewServices(cfgObj, ucs, client, repos)

	middlewares := di.NewMiddlewares(ucs)
	handlers := di.NewHandlers(cfgObj, ucs, client, services)

	routerImpl := routerConfig.NewRouter(
		middlewares.Auth, handlers.Auth, handlers.Bulk, handlers.User,
		handlers.Setting, handlers.Word, handlers.Quiz, handlers.Result)
	routerImpl.MountRoutes(router)

	// テスト用エンドポイント（開発環境のみ動作）
	if appEnv == "development" {
		setupTestEndpoints(router, runner)
	}

	// Trusted Proxies設定（環境変数から読み取り）
	trustedProxies := config.ParseTrustedProxies(config.Getenv("TRUSTED_PROXIES", ""))
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}

	validators.Init()
	binding.Validator = &validators.GinValidator{Validate: validators.V}
	return router
}

// サーバーを起動
func startServer(router *gin.Engine, port, env string) {
	logrus.Infof("Starting server on port %s in %s environment", port, env)
	if err := router.Run(":" + port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

// parseAccessLogOpts reads environment variables and returns AccessLogOpts.
func parseAccessLogOpts() middleware.AccessLogOpts {
	healthPath := config.Getenv("LOG_HEALTH_PATH", "/health")
	excludeHealth := config.GetenvBool("LOG_EXCLUDE_HEALTH", true)

	return middleware.AccessLogOpts{
		HealthPath:    healthPath,
		ExcludeHealth: excludeHealth,
	}
}
