// backend/src/infrastructure/repository/user/test/find_for_update_test.go
package user_test

import (
	"context"
	"testing"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"
	"word_app/backend/src/domain"
	repo "word_app/backend/src/infrastructure/repository/user"
	simock "word_app/backend/src/mocks/service_interfaces"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// func ptr[T any](v T) *T { return &v }

func newRepoWithRealEnt(t *testing.T) (*repo.EntUserRepo, *simock.MockEntClientInterface, *ent.Client) {
	t.Helper()
	cli := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = cli.Close() })
	m := simock.NewMockEntClientInterface(t)
	m.EXPECT().User().Return(cli.User)
	return repo.NewEntUserRepo(m), m, cli
}

func mustSeedUser(t *testing.T, cli *ent.Client, u ent.User) int {
	t.Helper()
	created, err := cli.User.Create().
		SetName(u.Name).
		SetNillableEmail(u.Email).
		SetNillablePassword(u.Password).
		SetIsAdmin(u.IsAdmin).
		SetIsRoot(u.IsRoot).
		SetIsTest(u.IsTest).
		Save(context.Background())
	require.NoError(t, err)
	return created.ID
}

func TestEntUserRepo_FindForUpdate(t *testing.T) {
	repo, mock, cli := newRepoWithRealEnt(t)
	ctx := context.Background()

	t.Run("Success_AllFieldsPresent", func(t *testing.T) {
		email := "a@example.com"
		pass := "hash"
		id := mustSeedUser(t, cli, ent.User{
			Name:     "Alice",
			Email:    &email,
			Password: &pass,
			IsAdmin:  true,
			IsRoot:   false,
			IsTest:   false,
		})

		got, err := repo.FindForUpdate(ctx, id)
		require.NoError(t, err)
		require.Equal(t, &domain.User{
			ID:        id,
			Name:      "Alice",
			Email:     &email,
			Password:  pass,
			IsAdmin:   true,
			IsRoot:    false,
			IsTest:    false,
			CreatedAt: got.CreatedAt, // 自動採番なので値は存在だけチェック
			UpdatedAt: got.UpdatedAt,
		}, got)
		require.False(t, got.CreatedAt.IsZero())
		require.False(t, got.UpdatedAt.IsZero())
	})

	t.Run("Success_EmailAndPasswordNil", func(t *testing.T) {
		id := mustSeedUser(t, cli, ent.User{
			Name:     "Bob",
			Email:    nil,
			Password: nil,
			IsAdmin:  false,
			IsRoot:   true,
			IsTest:   true,
		})

		got, err := repo.FindForUpdate(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "Bob", got.Name)
		require.Nil(t, got.Email)
		require.Equal(t, "", got.Password) // パスワードnil→空文字に変換して返す実装
		require.True(t, got.IsRoot)
		require.True(t, got.IsTest)
	})

	t.Run("NotFound_NoSuchID", func(t *testing.T) {
		_, err := repo.FindForUpdate(ctx, 999999)
		require.Error(t, err)
	})

	t.Run("NotFound_DeletedUser", func(t *testing.T) {
		id := mustSeedUser(t, cli, ent.User{Name: "Carol"})
		// 論理削除
		_, err := cli.User.UpdateOneID(id).SetDeletedAt(time.Now()).Save(ctx)
		require.NoError(t, err)

		_, err = repo.FindForUpdate(ctx, id)
		require.Error(t, err) // DeletedAtIsNil で除外され見つからない扱い
	})

	t.Run("Internal_DBClosed", func(t *testing.T) {
		// 一旦通常の repo を作った後、閉じた client の UserClient を返すよう差し替える
		tmp := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		require.NoError(t, tmp.Close())
		mock.ExpectedCalls = nil
		mock.EXPECT().User().Return(tmp.User)

		_, err := repo.FindForUpdate(ctx, 1)
		require.Error(t, err)
	})
}
