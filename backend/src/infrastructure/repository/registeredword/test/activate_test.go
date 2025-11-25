// backend/src/infrastructure/repository/registeredword/activate_test.go
package registeredword_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	entrw "word_app/backend/ent/registeredword"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/infrastructure/repository/registeredword"
)

func TestEntRegisteredWordWriteRepo_Activate(t *testing.T) {
	ctx := context.Background()

	t.Run("success - activate inactive word", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user := createUser(t, adapter.EntClient(), "activate@example.com", "Activate User")
		word := createWord(t, adapter.EntClient(), "apple")

		// まず非アクティブで作成
		createRegisteredWord(t, adapter.EntClient(), user.ID, word.ID, false)

		// Activateを実行
		err := repo.Activate(ctx, user.ID, word.ID)
		assert.NoError(t, err)

		// アクティブになったか確認
		rw, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.UserIDEQ(user.ID), entrw.WordIDEQ(word.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.True(t, rw.IsActive)
	})

	t.Run("success - activate multiple inactive words", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user := createUser(t, adapter.EntClient(), "multi@example.com", "Multi User")
		w1 := createWord(t, adapter.EntClient(), "apple")
		w2 := createWord(t, adapter.EntClient(), "banana")
		w3 := createWord(t, adapter.EntClient(), "cherry")

		createRegisteredWord(t, adapter.EntClient(), user.ID, w1.ID, false)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w2.ID, false)
		createRegisteredWord(t, adapter.EntClient(), user.ID, w3.ID, false)

		err1 := repo.Activate(ctx, user.ID, w1.ID)
		assert.NoError(t, err1)

		err2 := repo.Activate(ctx, user.ID, w2.ID)
		assert.NoError(t, err2)

		err3 := repo.Activate(ctx, user.ID, w3.ID)
		assert.NoError(t, err3)

		// 全てアクティブになったか確認
		count, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.UserIDEQ(user.ID), entrw.IsActive(true)).
			Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("success - activate already active word (idempotent)", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user := createUser(t, adapter.EntClient(), "idempotent@example.com", "Idempotent User")
		word := createWord(t, adapter.EntClient(), "apple")

		createRegisteredWord(t, adapter.EntClient(), user.ID, word.ID, true)

		// 既にアクティブなものを再度アクティブ化
		err := repo.Activate(ctx, user.ID, word.ID)
		assert.NoError(t, err)

		// まだアクティブか確認
		rw, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.UserIDEQ(user.ID), entrw.WordIDEQ(word.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.True(t, rw.IsActive)
	})

	t.Run("success - activate only specific user's word", func(t *testing.T) {
		adapter, repo := newWriteRepo(t)

		user1 := createUser(t, adapter.EntClient(), "user1@example.com", "User 1")
		user2 := createUser(t, adapter.EntClient(), "user2@example.com", "User 2")
		w1 := createWord(t, adapter.EntClient(), "word1")

		createRegisteredWord(t, adapter.EntClient(), user1.ID, w1.ID, false)
		createRegisteredWord(t, adapter.EntClient(), user2.ID, w1.ID, false)

		// user1の単語だけをアクティブ化
		err := repo.Activate(ctx, user1.ID, w1.ID)
		assert.NoError(t, err)

		// user1はアクティブ、user2は非アクティブ
		rw1, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.UserIDEQ(user1.ID), entrw.WordIDEQ(w1.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.True(t, rw1.IsActive)

		rw2, err := adapter.EntClient().RegisteredWord.Query().
			Where(entrw.UserIDEQ(user2.ID), entrw.WordIDEQ(w1.ID)).
			Only(ctx)
		require.NoError(t, err)
		assert.False(t, rw2.IsActive)
	})

	t.Run("success - no error for non-existent registration", func(t *testing.T) {
		_, repo := newWriteRepo(t)

		// 存在しない組み合わせでActivate
		err := repo.Activate(ctx, 99999, 99999)
		assert.NoError(t, err) // エラーにならない（entの仕様）
	})

	t.Run("error - database closed", func(t *testing.T) {
		closedClient := enttest.Open(t, "sqlite3", "file:closeddb?mode=memory&cache=shared&_fk=1")
		closedAdapter := entAdapter{closedClient}
		closedRepo := registeredword.NewEntRegisteredWordWriteRepo(closedAdapter)

		_ = closedClient.Close()

		err := closedRepo.Activate(ctx, 1, 1)
		assert.Error(t, err)
	})
}
