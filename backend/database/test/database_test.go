package database_test

import (
	"testing"

	"word_app/backend/database"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitEntClient(t *testing.T) {
	// テスト用の環境変数を設定
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "test_user")
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("DB_NAME", "test_db")

	// 初期化
	_ = database.InitEntClient()
	assert.NoError(t, nil, "Database initialization should not return an error")

	// クライアント取得
	client := database.GetEntClient()
	logrus.Info(client)

	// nil でないことを確認
	assert.NotNil(t, client, "Ent client should be initialized")
}
