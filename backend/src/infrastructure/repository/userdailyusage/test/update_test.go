package userdailyusage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	si "word_app/backend/src/interfaces/service_interfaces"
)

// ヘルパー関数: 既存のusageレコードを作成
// 注意: ent.Clientで作成しても、sql.Runnerの*sql.DBとは別のDBになる可能性がある
func createUsage(t *testing.T, adapter si.EntClientInterface, userID int, year, month, day, quizCount, bulkCount int) {
	t.Helper()
	jst, _ := time.LoadLocation("Asia/Tokyo")
	lastResetDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, jst)

	// ent.Clientで直接作成
	_, err := adapter.EntClient().UserDailyUsage.Create().
		SetUserID(userID).
		SetLastResetDate(lastResetDate).
		SetQuizCount(quizCount).
		SetBulkCount(bulkCount).
		Save(context.Background())
	require.NoError(t, err)
}

func TestEntUserDailyUsageRepo_IncQuizOr429(t *testing.T) {
	ctx := context.Background()

	t.Run("success - first quiz today", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 1)
		// 初期化
		err := repo.CreateIfNotExists(ctx, user.ID, now)
		require.NoError(t, err)

		result, err := repo.IncQuizOr429(ctx, user.ID, now, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.QuizCount)
		assert.Equal(t, 0, result.BulkCount)
	})

	t.Run("success - second quiz today", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 2)
		// 1回目のIncQuizOr429を実行してquiz_countを1にする
		err := repo.CreateIfNotExists(ctx, user.ID, now)
		require.NoError(t, err)
		_, err = repo.IncQuizOr429(ctx, user.ID, now, 10)
		require.NoError(t, err)

		// 2回目のIncQuizOr429を実行
		result, err := repo.IncQuizOr429(ctx, user.ID, now, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.QuizCount)
		assert.Equal(t, 0, result.BulkCount)
	})

	t.Run("success - reset if yesterday", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 3)
		// 昨日のデータ
		createUsage(t, adapter, user.ID, 2025, 1, 1, 5, 3)

		result, err := repo.IncQuizOr429(ctx, user.ID, now, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.QuizCount) // リセットされて1
		assert.Equal(t, 0, result.BulkCount) // bulk_countも0にリセット
	})

	t.Run("error - quota exceeded", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		// JSTで2025/1/3 0:00にする
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 3, 0, 0, 0, 0, jst)
		user := createUser(t, adapter.EntClient(), 4)
		// 上限に達している（quiz_count = 10, dailyCap = 10）
		createUsage(t, adapter, user.ID, 2025, 1, 3, 10, 0)

		// 10回目は成功するが、11回目はエラーになる
		// 実際にインクリメントを試す
		result, err := repo.IncQuizOr429(ctx, user.ID, now, 10)

		// SQLの仕様により、quiz_count < $3 のWHERE句で更新されない
		// これは sql.ErrNoRows として返される
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("success - default daily cap when 0", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 5)
		err := repo.CreateIfNotExists(ctx, user.ID, now)
		require.NoError(t, err)

		// dailyCap=0 の場合はデフォルト値が使われる（quiz なら 20）
		result, err := repo.IncQuizOr429(ctx, user.ID, now, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestEntUserDailyUsageRepo_IncBulkOr429(t *testing.T) {
	ctx := context.Background()

	t.Run("success - first bulk today", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 11)
		err := repo.CreateIfNotExists(ctx, user.ID, now)
		require.NoError(t, err)

		result, err := repo.IncBulkOr429(ctx, user.ID, now, 5)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.QuizCount)
		assert.Equal(t, 1, result.BulkCount)
	})

	t.Run("success - second bulk today", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 12)
		// 1回目のIncBulkOr429を実行してbulk_countを1にする
		err := repo.CreateIfNotExists(ctx, user.ID, now)
		require.NoError(t, err)
		_, err = repo.IncBulkOr429(ctx, user.ID, now, 5)
		require.NoError(t, err)

		// 2回目のIncBulkOr429を実行
		result, err := repo.IncBulkOr429(ctx, user.ID, now, 5)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.QuizCount)
		assert.Equal(t, 2, result.BulkCount)
	})

	t.Run("success - reset if yesterday", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 13)
		// 昨日のデータ
		createUsage(t, adapter, user.ID, 2025, 1, 1, 3, 4)

		result, err := repo.IncBulkOr429(ctx, user.ID, now, 5)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.QuizCount) // quiz_countも0にリセット
		assert.Equal(t, 1, result.BulkCount) // リセットされて1
	})

	t.Run("error - quota exceeded", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 3, 0, 0, 0, 0, jst)
		user := createUser(t, adapter.EntClient(), 14)
		// 上限に達している
		createUsage(t, adapter, user.ID, 2025, 1, 3, 0, 5)

		result, err := repo.IncBulkOr429(ctx, user.ID, now, 5)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("success - default daily cap when 0", func(t *testing.T) {
		adapter, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)
		user := createUser(t, adapter.EntClient(), 15)
		err := repo.CreateIfNotExists(ctx, user.ID, now)
		require.NoError(t, err)

		// dailyCap=0 の場合はデフォルト値が使われる（bulk なら 5）
		result, err := repo.IncBulkOr429(ctx, user.ID, now, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("error - user not exists", func(t *testing.T) {
		_, _, repo := newRepo(t)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Date(2025, 1, 2, 15, 30, 45, 0, jst)

		result, err := repo.IncBulkOr429(ctx, 9999, now, 5)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
