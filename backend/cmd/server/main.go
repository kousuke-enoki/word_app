package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"word_app/backend/config"
	"word_app/backend/database"
	"word_app/backend/ent/migrate"
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
	// 外部キーも明確に生成するようにする。
	return entClient.Schema.Create(ctx, migrate.WithForeignKeys(true))
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
		AllowOrigins:     allowed,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	// 2) もし *.vercel.app を許可したい場合（合わせ技OK）
	cfg.AllowOriginFunc = func(origin string) bool {
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
