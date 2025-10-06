// backend/src/infrastructure/repository/user/find_detail_by_id_test.go
package user

import (
	"context"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	repo "word_app/backend/src/infrastructure/repository/user"
	serviceinterfaces "word_app/backend/src/mocks/service_interfaces"
)

// func ptr[T any](v T) *T { return &v }

// FindDetailByID はリポジトリ内で client.User() しか呼ばないので、UserClient を返すモックだけ用意すればOK
func newDetailRepo(t *testing.T, uc *ent.UserClient) *repo.EntUserRepo {
	t.Helper()
	m := serviceinterfaces.NewMockEntClientInterface(t)
	m.EXPECT().User().Return(uc)
	return repo.NewEntUserRepo(m)
}

func TestEntUserRepo_FindDetailByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_WithLINEAuth_lowercase", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fdd1?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("Alice").
			SetPassword("hashed").
			SetEmail("alice@example.com").
			Save(ctx)
		require.NoError(t, err)

		// provider = "line"（小文字）
		_, err = cli.ExternalAuth.Create().
			SetProvider("line").
			SetProviderUserID("alice-line-sub").
			SetUser(u).
			Save(ctx)
		require.NoError(t, err)

		r := newDetailRepo(t, cli.User)

		got, err := r.FindDetailByID(ctx, u.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, u.ID, got.ID)
		require.Equal(t, "Alice", got.Name)
		require.NotNil(t, got.Email)
		require.Equal(t, "alice@example.com", *got.Email)
		// WithExternalAuths で LINE のみを preload。ドメイン上の詳細フィールドは実装依存なのでここでは存在/成功の確認まで。
	})

	t.Run("Success_WithLINEAuth_uppercase(EqualFold)", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fdd2?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("Bob").
			SetPassword("hashed").
			SetEmail("bob@example.com").
			Save(ctx)
		require.NoError(t, err)

		// provider = "LINE"（大文字） ⇒ EqualFold("line") にヒットするはず
		_, err = cli.ExternalAuth.Create().
			SetProvider("LINE").
			SetProviderUserID("bob-line-sub").
			SetUser(u).
			Save(ctx)
		require.NoError(t, err)

		r := newDetailRepo(t, cli.User)

		got, err := r.FindDetailByID(ctx, u.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, u.ID, got.ID)
		require.Equal(t, "Bob", got.Name)
		require.NotNil(t, got.Email)
		require.Equal(t, "bob@example.com", *got.Email)
	})

	t.Run("Success_NoLINEAuth_loadedEmpty", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fdd3?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("Carol").
			SetPassword("hashed").
			SetEmail("carol@example.com").
			Save(ctx)
		require.NoError(t, err)

		// 外部認証は付けない（= WithExternalAuths は空のプリロード）
		r := newDetailRepo(t, cli.User)

		got, err := r.FindDetailByID(ctx, u.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, u.ID, got.ID)
		require.Equal(t, "Carol", got.Name)
		require.NotNil(t, got.Email)
		require.Equal(t, "carol@example.com", *got.Email)
	})

	t.Run("Success_WithOtherProviders_onlyLINEPreloaded", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fdd4?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("Dave").
			SetPassword("hashed").
			SetEmail("dave@example.com").
			Save(ctx)
		require.NoError(t, err)

		// GitHub 連携（LINE ではない → WithExternalAuths の条件にマッチしない）
		_, err = cli.ExternalAuth.Create().
			SetProvider("github").
			SetProviderUserID("dave-gh").
			SetUser(u).
			Save(ctx)
		require.NoError(t, err)

		r := newDetailRepo(t, cli.User)

		got, err := r.FindDetailByID(ctx, u.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, u.ID, got.ID)
		require.Equal(t, "Dave", got.Name)
		require.NotNil(t, got.Email)
		require.Equal(t, "dave@example.com", *got.Email)
		// WithExternalAuths は "line" のみを preload するため、ここでは LINE が無ければ空のはず（ドメインへの反映は mapper 実装依存）
	})

	t.Run("NotFound_NoSuchUser", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fdd5?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		r := newDetailRepo(t, cli.User)

		_, err := r.FindDetailByID(ctx, 99999) // 存在しない ID
		require.Error(t, err)
	})

	t.Run("Internal_DBClosed", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fdd6?mode=memory&cache=shared&_fk=1")
		require.NoError(t, cli.Close())

		r := newDetailRepo(t, cli.User)

		_, err := r.FindDetailByID(ctx, 1)
		require.Error(t, err)
	})
}
