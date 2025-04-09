// backend/cmd/import_dict/main.go
package main

import (
	"context"
	"flag"
	"os"
	"time"

	"word_app/backend/config"
	"word_app/backend/database"
	"word_app/backend/internal/dictimport" // ← 取込ロジック（別途実装）
	"word_app/backend/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	// -------- CLI フラグ --------
	var (
		file      string
		workers   int
		batchSize int
	)
	flag.StringVar(&file, "file", "jmdict.json", "path to JMdict JSON (unzipped)")
	flag.IntVar(&workers, "workers", 4, "concurrent workers")
	flag.IntVar(&batchSize, "batch", 500, "bulk insert chunk size")
	flag.Parse()

	// -------- 共通初期化 --------
	config.LoadEnv()    // .env 読み込み
	logger.InitLogger() // logrus 設定
	database.InitEntClient()
	cli := database.GetEntClient()
	defer cli.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	logrus.Infof("JMdict import start (file=%s)", file)

	// -------- インポート実行 --------
	opts := dictimport.Options{
		Workers:   workers,
		BatchSize: batchSize,
	}
	if err := dictimport.ImportJMdict(ctx, file, cli, opts); err != nil {
		logrus.Fatalf("import failed: %v", err)
	}

	logrus.Infof("JMdict import completed in %s", time.Since(start))
	os.Exit(0)
}
