// backend/src/infrastructure/repository/user/find_by_provider_test.go
package user

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// func newRepoWithEntUserClient(t *testing.T, uc *ent.UserClient) *repo.EntUserRepo {
// 	t.Helper()
// 	m := serviceinterfaces.NewMockEntClientInterface(t)
// 	// FindByProvider 本体は User() しか呼ばないので User() だけ返せばOK
// 	m.EXPECT().User().Return(uc)
// 	return repo.NewEntUserRepo(m)
// }

// func ptr[T any](v T) *T { return &v }

func TestEntUserRepo_FindByProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_Linked", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fbp1?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		// ユーザー作成
		u, err := cli.User.Create().
			SetName("LinkedUser").
			SetPassword("hashed").
			SetEmail("linked@example.com").
			Save(ctx)
		require.NoError(t, err)

		// 外部認証連携（provider + subject）
		const provider = "line"
		const subject = "sub-123"
		_, err = cli.ExternalAuth.Create().
			SetProvider(provider).
			SetProviderUserID(subject).
			SetUser(u).
			Save(ctx)
		require.NoError(t, err)

		// Repo
		r := newRepoWithEntUserClient(t, cli.User)

		got, err := r.FindByProvider(ctx, provider, subject)
		require.NoError(t, err)

		// 実装は ent.User の値から ID/Email/Name を詰め替えて返す
		require.Equal(t, u.ID, got.ID)
		require.NotNil(t, got.Email)
		require.Equal(t, "linked@example.com", *got.Email)
		require.Equal(t, "LinkedUser", got.Name)
	})

	t.Run("NotFound_UserExistsButNoExternalAuth", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fbp2?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		// ユーザーは作るが外部認証を付けない
		_, err := cli.User.Create().
			SetName("NoLink").
			SetPassword("hashed").
			SetEmail("nolink@example.com").
			Save(ctx)
		require.NoError(t, err)

		r := newRepoWithEntUserClient(t, cli.User)

		_, err = r.FindByProvider(ctx, "google", "no-such-sub")
		require.Error(t, err, "未連携は NotFound 道に入るはず")
		// repoerr の型やメッセージを厳密にチェックしたい場合はここで追加
	})

	t.Run("NotFound_ProviderOrSubjectMismatch", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:fbp3?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("WithLink").
			SetPassword("hashed").
			SetEmail("withlink@example.com").
			Save(ctx)
		require.NoError(t, err)

		// 正しい連携は "google" / "sub-abc"
		_, err = cli.ExternalAuth.Create().
			SetProvider("google").
			SetProviderUserID("sub-abc").
			SetUser(u).
			Save(ctx)
		require.NoError(t, err)

		r := newRepoWithEntUserClient(t, cli.User)

		// provider 不一致
		_, err = r.FindByProvider(ctx, "github", "sub-abc")
		require.Error(t, err)

		// subject 不一致
		_, err = r.FindByProvider(ctx, "google", "sub-xyz")
		require.Error(t, err)
	})

	t.Run("Internal_DBClosed", func(t *testing.T) {
		// いったん開いてからクローズした ent.Client の UserClient を渡す
		cli := enttest.Open(t, "sqlite3", "file:fbp4?mode=memory&cache=shared&_fk=1")
		require.NoError(t, cli.Close())

		r := newRepoWithEntUserClient(t, cli.User)

		_, err := r.FindByProvider(ctx, "google", "anything")
		require.Error(t, err, "クローズ済み DB は内部エラー経路に入る想定")
	})
}
