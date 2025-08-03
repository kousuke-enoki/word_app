// backend/cmd/import_dict/main.go
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"word_app/backend/config"
	"word_app/backend/database"
	"word_app/backend/internal/dictimport"
	"word_app/backend/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	log.Println("JMdict import start")
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
	defer func() {
		if err := cli.Close(); err != nil {
			logrus.Fatalf("failed to close ent client: %v", err)
		}
	}()

	// スキーマを作成（存在すれば no‑op）
	if err := cli.Schema.Create(context.Background()); err != nil {
		log.Fatalf("schema create failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	log.Printf("JMdict import start (file=%s)", file)

	// -------- インポート実行 --------
	opts := dictimport.Options{
		Workers:   workers,
		BatchSize: batchSize,
	}

	errs, fatal := dictimport.ImportJMdict(ctx, file, cli, opts)
	if fatal != nil {
		log.Fatalf("import failed: %v", fatal)
	}
	log.Printf("import finished. failures=%d\n", len(errs))
	for _, e := range errs {
		log.Println(e.ID, e.Message)
	}

	log.Printf("JMdict import completed in %s", time.Since(start))
	os.Exit(0)
}
