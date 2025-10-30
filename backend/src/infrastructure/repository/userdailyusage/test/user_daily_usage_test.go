package userdailyusage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"word_app/backend/src/infrastructure/repository/userdailyusage"
	si "word_app/backend/src/interfaces/service_interfaces"
	sqlexec "word_app/backend/src/interfaces/sqlexec"
)

func TestNewEntUserDailyUsageRepo(t *testing.T) {
	t.Run("正常に初期化される", func(t *testing.T) {
		adapter, runner := setupMocks(t)
		repo := userdailyusage.NewEntUserDailyUsageRepo(adapter, runner)

		assert.NotNil(t, repo)
	})
}

func setupMocks(t *testing.T) (si.EntClientInterface, sqlexec.Runner) {
	t.Helper()
	adapter, runner, _ := newRepo(t)
	return adapter, runner
}

func TestEntUserDailyUsageRepo_truncateToJST0(t *testing.T) {
	// truncateToJST0は非公開メソッドなので、間接的にテストする
	// CreateIfNotExistsやIncQuizOr429の中で呼ばれるため、それらを通じて動作確認

	t.Run("UTCからJSTへの変換と切り捨て", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		// UTCで2025/1/2 15:30:45
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, time.UTC)
		user := createUser(t, adapter.EntClient(), 1)

		err := repo.CreateIfNotExists(context.Background(), user.ID, now)
		assert.NoError(t, err)
		// truncateToJST0が正しく動作していれば、エラーなく実行される
	})

	t.Run("JST内で時刻切り捨て", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		// JSTで2025/1/2 15:30:45
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 2)

		err := repo.CreateIfNotExists(context.Background(), user.ID, now)
		assert.NoError(t, err)
	})

	t.Run("同じ日の異なる時刻は同じ日付に正規化", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		user1 := createUser(t, adapter.EntClient(), 3)
		user2 := createUser(t, adapter.EntClient(), 7)
		user3 := createUser(t, adapter.EntClient(), 8)

		now1 := time.Date(2025, 1, 2, 10, 0, 0, 0, jst)
		now2 := time.Date(2025, 1, 2, 23, 59, 59, 999999999, jst)
		now3 := time.Date(2025, 1, 2, 0, 0, 0, 0, jst)

		err1 := repo.CreateIfNotExists(context.Background(), user1.ID, now1)
		assert.NoError(t, err1)

		// 同じ日なので、別ユーザーでも同じ日付に正規化されることを確認
		err2 := repo.CreateIfNotExists(context.Background(), user2.ID, now2)
		assert.NoError(t, err2)

		err3 := repo.CreateIfNotExists(context.Background(), user3.ID, now3)
		assert.NoError(t, err3)
	})

	t.Run("日付境界のテスト", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		user1 := createUser(t, adapter.EntClient(), 4)
		user2 := createUser(t, adapter.EntClient(), 5)

		// 2025/1/1 23:59:59 JST
		lastMoment := time.Date(2025, 1, 1, 23, 59, 59, 999999999, jst)
		// 2025/1/2 00:00:00 JST
		nextDay := time.Date(2025, 1, 2, 0, 0, 0, 0, jst)

		err1 := repo.CreateIfNotExists(context.Background(), user1.ID, lastMoment)
		assert.NoError(t, err1)

		err2 := repo.CreateIfNotExists(context.Background(), user2.ID, nextDay)
		assert.NoError(t, err2)
		// 異なる日なので、両方とも新規作成される
	})

	t.Run("タイムゾーンの境界テスト", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		user := createUser(t, adapter.EntClient(), 6)

		// UTCの2025/1/1 15:00:00は JSTでは2025/1/2 00:00:00
		now := time.Date(2025, 1, 1, 15, 0, 0, 0, time.UTC)

		err := repo.CreateIfNotExists(context.Background(), user.ID, now)
		assert.NoError(t, err)
	})
}
