// backend/src/infrastructure/repository/user/is_root_test.go
package user

import (
	"context"
	"testing"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	repo "word_app/backend/src/infrastructure/repository/user"
	serviceinterfaces "word_app/backend/src/mocks/service_interfaces"
)

func newIsRootRepo(t *testing.T, uc *ent.UserClient) *repo.EntUserRepo {
	t.Helper()
	m := serviceinterfaces.NewMockEntClientInterface(t)
	m.EXPECT().User().Return(uc)
	return repo.NewEntUserRepo(m)
}

// func ptr[T any](v T) *T { return &v }

func TestEntUserRepo_IsRoot(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_RootUser_ReturnsTrue", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:isroot1?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("Eve").
			SetPassword("hashed").
			SetIsRoot(true).
			SetEmail("eve@example.com").
			Save(ctx)
		require.NoError(t, err)

		r := newIsRootRepo(t, cli.User)
		got, err := r.IsRoot(ctx, u.ID)
		require.NoError(t, err)
		require.True(t, got)
	})

	t.Run("Success_NormalUser_ReturnsFalse", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:isroot2?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("Alice").
			SetPassword("hashed").
			SetIsRoot(false).
			SetEmail("alice@example.com").
			Save(ctx)
		require.NoError(t, err)

		r := newIsRootRepo(t, cli.User)
		got, err := r.IsRoot(ctx, u.ID)
		require.NoError(t, err)
		require.False(t, got)
	})

	t.Run("Success_DeletedUser_StillQueryable_ReturnsFlag", func(t *testing.T) {
		// IsRoot は DeletedAt で絞っていないため、削除済みでも取得可能な挙動を確認
		cli := enttest.Open(t, "sqlite3", "file:isroot3?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		delAt := time.Now()
		u, err := cli.User.Create().
			SetName("Mallory").
			SetPassword("hashed").
			SetIsRoot(false).
			SetEmail("mallory@example.com").
			SetNillableDeletedAt(&delAt).
			Save(ctx)
		require.NoError(t, err)

		r := newIsRootRepo(t, cli.User)
		got, err := r.IsRoot(ctx, u.ID)
		require.NoError(t, err)
		require.False(t, got)
	})

	t.Run("NotFound_NoSuchID_ReturnsError", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:isroot4?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		// 何も作らず存在しない ID を問い合わせ
		r := newIsRootRepo(t, cli.User)
		_, err := r.IsRoot(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("Internal_DBClosed_ReturnsError", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:isroot5?mode=memory&cache=shared&_fk=1")
		require.NoError(t, cli.Close()) // クローズして内部エラーを誘発

		r := newIsRootRepo(t, cli.User)
		_, err := r.IsRoot(ctx, 1)
		require.Error(t, err)
	})
}
