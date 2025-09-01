package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"word_app/backend/config"
	"word_app/backend/database"
	"word_app/backend/ent/user"
	"word_app/backend/internal/di"
	"word_app/backend/logger"
	routerConfig "word_app/backend/router"
	"word_app/backend/seeder"
	"word_app/backend/src/infrastructure"
	"word_app/backend/src/interfaces"
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
	// bootstrapMode := os.Getenv("APP_BOOTSTRAP_MODE") // "HEALTH_ONLY" / "FULL" など

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
// func initializeServer() {
func mustInitServer(needCleanup bool) (*gin.Engine, string, string, func()) {
	defer func() {
		if p := recover(); p != nil {
			logrus.Fatalf("PANIC caught in main: %v\n", p)
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
		logrus.Error(err)
	}
	entClient := database.GetEntClient()
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
	// setupDatabase(client)

	runMig := shouldRun("RUN_MIGRATION")
	runSeed := shouldRun("RUN_SEEDER")
	if !isLambda() && os.Getenv("RUN_MIGRATION") == "" {
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

	router := setupRouter(client, corsOrigin)

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
	return entClient.Schema.Create(ctx)
}

func runSeederIfNeeded(client interfaces.ClientInterface) error {
	ctx := context.Background()
	entClient := client.EntClient()
	if entClient == nil {
		return fmt.Errorf("ent.Client is nil")
	}
	adminExists, err := entClient.User.Query().Where(user.Email("root@example.com")).Exist(ctx)
	if err != nil {
		return fmt.Errorf("admin check failed: %w", err)
	}
	if !adminExists {
		logrus.Info("Running initial seeder...")
		seeder.RunSeeder(ctx, client)
		logrus.Info("Seeder completed.")
	} else {
		logrus.Info("Seed data already exists, skipping.")
	}
	return nil
}

// データベースのセットアップ
func setupDatabase(client interfaces.ClientInterface) {
	ctx := context.Background()
	entClient := client.EntClient()
	if entClient == nil {
		logrus.Fatalf("ent.Client is nil")
	}
	// Schema を作成
	if err := entClient.Schema.Create(ctx); err != nil {
		logrus.Fatalf("Failed to create schema: %v", err)
	}
	// Admin の存在を確認
	adminExists, err := entClient.User.Query().Where(user.Email("root@example.com")).Exist(ctx)
	if err != nil {
		logrus.Fatalf("Failed to check admin existence: %v", err)
	}
	// Seeder の実行
	if !adminExists {
		logrus.Info("Running initial seeder...")
		seeder.RunSeeder(ctx, client)
		logrus.Info("Seeder completed.")
	} else {
		logrus.Info("Seed data already exists, skipping.")
	}
}

// ルートを構築する関数
func setupRouter(client interfaces.ClientInterface, corsOrigin string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

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
	repos := di.NewRepositories(client)
	ucs, err := di.NewUseCases(cfgObj, repos)
	if err != nil {
		logrus.Fatal(err)
	}

	handlers := di.NewHandlers(cfgObj, ucs, client)

	routerImpl := routerConfig.NewRouter(
		handlers.JWTMiD, handlers.Auth, handlers.User,
		handlers.Setting, handlers.Word, handlers.Quiz, handlers.Result)
	routerImpl.MountRoutes(router)
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
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
