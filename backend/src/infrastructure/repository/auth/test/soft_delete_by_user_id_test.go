package auth_test

import (
	"context"
	"testing"
	"time"

	"word_app/backend/ent/enttest"
	"word_app/backend/ent/externalauth"

	infra "word_app/backend/src/infrastructure"
	repo "word_app/backend/src/infrastructure/repository/auth"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntExtAuthRepo_SoftDeleteByUserID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:memdb?mode=memory&_fk=1")
	defer func() {
		if cerr := client.Close(); cerr != nil {
			t.Logf("failed to close ent test client: %v", cerr)
		}
	}()
	ctx := context.Background()

	r := repo.NewEntExtAuthRepo(infra.NewAppClient(client))

	t.Run("success - soft delete existing external auth", func(t *testing.T) {
		// ユーザーを作成
		user, err := client.User.
			Create().
			SetEmail("test@example.com").
			SetName("Test User").
			SetPassword("hashed").
			Save(ctx)
		require.NoError(t, err)

		// ExternalAuthを作成
		extAuth, err := client.ExternalAuth.
			Create().
			SetUserID(user.ID).
			SetProvider("google").
			SetProviderUserID("google-123").
			Save(ctx)
		require.NoError(t, err)
		assert.Nil(t, extAuth.DeletedAt)

		// SoftDeleteByUserIDを実行
		deleteTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
		err = r.SoftDeleteByUserID(ctx, user.ID, deleteTime)
		assert.NoError(t, err)

		// 削除されたか確認
		updated, err := client.ExternalAuth.
			Query().
			Where(externalauth.ID(extAuth.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.NotNil(t, updated.DeletedAt)
		assert.Equal(t, deleteTime, *updated.DeletedAt)
	})

	t.Run("success - soft delete one external auth", func(t *testing.T) {
		// ユーザーを作成
		user, err := client.User.
			Create().
			SetEmail("single@example.com").
			SetName("Single Auth User").
			SetPassword("hashed").
			Save(ctx)
		require.NoError(t, err)

		// ExternalAuthを作成
		extAuth, err := client.ExternalAuth.
			Create().
			SetUserID(user.ID).
			SetProvider("google").
			SetProviderUserID("google-single").
			Save(ctx)
		require.NoError(t, err)
		assert.Nil(t, extAuth.DeletedAt)

		// SoftDeleteByUserIDを実行
		deleteTime := time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC)
		err = r.SoftDeleteByUserID(ctx, user.ID, deleteTime)
		assert.NoError(t, err)

		// 削除されたか確認
		updated, err := client.ExternalAuth.
			Query().
			Where(externalauth.ID(extAuth.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.NotNil(t, updated.DeletedAt)
	})

	t.Run("success - no error if already deleted", func(t *testing.T) {
		// ユーザーを作成
		user, err := client.User.
			Create().
			SetEmail("deleted@example.com").
			SetName("Deleted Auth User").
			SetPassword("hashed").
			Save(ctx)
		require.NoError(t, err)

		// 既に削除されたExternalAuthを作成
		deletedTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		extAuth, err := client.ExternalAuth.
			Create().
			SetUserID(user.ID).
			SetProvider("google").
			SetProviderUserID("google-789").
			SetDeletedAt(deletedTime).
			Save(ctx)
		require.NoError(t, err)
		assert.NotNil(t, extAuth.DeletedAt)

		// SoftDeleteByUserIDを実行（既に削除済みなので何もしない）
		newDeleteTime := time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC)
		err = r.SoftDeleteByUserID(ctx, user.ID, newDeleteTime)
		assert.NoError(t, err)

		// 削除時刻は変更されない（既に削除済みのものは除外される）
		updated, err := client.ExternalAuth.
			Query().
			Where(externalauth.ID(extAuth.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.NotNil(t, updated.DeletedAt)
		assert.Equal(t, deletedTime, *updated.DeletedAt) // 元の削除時刻のまま
	})

	t.Run("success - no error if no external auth exists", func(t *testing.T) {
		// ユーザーを作成
		user, err := client.User.
			Create().
			SetEmail("noauth@example.com").
			SetName("No Auth User").
			SetPassword("hashed").
			Save(ctx)
		require.NoError(t, err)

		// ExternalAuthは作成しない

		// SoftDeleteByUserIDを実行（該当なし）
		deleteTime := time.Date(2025, 1, 3, 12, 0, 0, 0, time.UTC)
		err = r.SoftDeleteByUserID(ctx, user.ID, deleteTime)
		assert.NoError(t, err)
		// エラーが発生しないことを確認
	})

	t.Run("success - only delete auths for specific user", func(t *testing.T) {
		// 2つのユーザーを作成
		user1, err := client.User.
			Create().
			SetEmail("user1@example.com").
			SetName("User 1").
			SetPassword("hashed").
			Save(ctx)
		require.NoError(t, err)

		user2, err := client.User.
			Create().
			SetEmail("user2@example.com").
			SetName("User 2").
			SetPassword("hashed").
			Save(ctx)
		require.NoError(t, err)

		// 各ユーザーのExternalAuthを作成
		extAuth1, err := client.ExternalAuth.
			Create().
			SetUserID(user1.ID).
			SetProvider("google").
			SetProviderUserID("google-u1").
			Save(ctx)
		require.NoError(t, err)

		extAuth2, err := client.ExternalAuth.
			Create().
			SetUserID(user2.ID).
			SetProvider("google").
			SetProviderUserID("google-u2").
			Save(ctx)
		require.NoError(t, err)
		assert.Nil(t, extAuth1.DeletedAt)
		assert.Nil(t, extAuth2.DeletedAt)

		// user1のExternalAuthだけを削除
		deleteTime := time.Date(2025, 1, 4, 12, 0, 0, 0, time.UTC)
		err = r.SoftDeleteByUserID(ctx, user1.ID, deleteTime)
		assert.NoError(t, err)

		// user1のExternalAuthが削除されたか確認
		updated1, err := client.ExternalAuth.
			Query().
			Where(externalauth.ID(extAuth1.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.NotNil(t, updated1.DeletedAt)

		// user2のExternalAuthは削除されていないか確認
		updated2, err := client.ExternalAuth.
			Query().
			Where(externalauth.ID(extAuth2.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.Nil(t, updated2.DeletedAt)
	})

	t.Run("success - no error for non-existent user", func(t *testing.T) {
		// 存在しないユーザーIDでSoftDeleteByUserIDを実行
		deleteTime := time.Date(2025, 1, 5, 12, 0, 0, 0, time.UTC)
		err := r.SoftDeleteByUserID(ctx, 99999, deleteTime)
		assert.NoError(t, err)
		// 該当するレコードがない場合もエラーにならない
	})

	t.Run("error - database closed", func(t *testing.T) {
		// 新しいクライアントを作成
		closedClient := enttest.Open(t, "sqlite3", "file:closeddb?mode=memory&_fk=1")

		closedRepo := repo.NewEntExtAuthRepo(infra.NewAppClient(closedClient))

		// クライアントを閉じる
		_ = closedClient.Close()

		// 閉じたクライアントでSoftDeleteByUserIDを実行
		deleteTime := time.Date(2025, 1, 6, 12, 0, 0, 0, time.UTC)
		err := closedRepo.SoftDeleteByUserID(ctx, 1, deleteTime)
		assert.Error(t, err)
	})
}
