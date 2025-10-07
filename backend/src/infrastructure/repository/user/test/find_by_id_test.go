// backend/src/infrastructure/repository/user/find_by_id_test.go
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

func newRepoWithEntUserClient(t *testing.T, uc *ent.UserClient) *repo.EntUserRepo {
	t.Helper()
	m := serviceinterfaces.NewMockEntClientInterface(t)
	m.EXPECT().User().Return(uc)
	return repo.NewEntUserRepo(m)
}

func TestEntUserRepo_FindByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_IsRootFalse", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:findbyid1?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		// 先に1件作っておく（IsRoot=false）
		u, err := cli.User.Create().
			SetName("User1").
			SetPassword("x").
			Save(ctx)
		require.NoError(t, err)

		r := newRepoWithEntUserClient(t, cli.User)
		got, err := r.FindByID(ctx, u.ID)
		require.NoError(t, err)

		require.Equal(t, u.ID, got.ID)
		require.False(t, got.IsRoot, "default IsRoot should be false")

		// FindByID は ID/IsRoot のみ返す実装なので他はゼロ値のまま
		require.Zero(t, got.Name)
		require.Nil(t, got.Email)
	})

	t.Run("Success_IsRootTrue", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:findbyid2?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		// IsRoot=true で1件
		u, err := cli.User.Create().
			SetName("Admin").
			SetPassword("x").
			SetIsRoot(true).
			Save(ctx)
		require.NoError(t, err)

		r := newRepoWithEntUserClient(t, cli.User)
		got, err := r.FindByID(ctx, u.ID)
		require.NoError(t, err)

		require.Equal(t, u.ID, got.ID)
		require.True(t, got.IsRoot, "should reflect IsRoot=true")
		// 他フィールドは返さない方針
		require.Zero(t, got.Name)
		require.Nil(t, got.Email)
	})

	t.Run("NotFound_NoSuchID", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:findbyid3?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		r := newRepoWithEntUserClient(t, cli.User)
		_, err := r.FindByID(ctx, 99999)
		require.Error(t, err, "unknown ID should return error")
		// repoerr の変換仕様に応じて型/メッセージを厳密に見るならここで追加検証
	})

	t.Run("Internal_DBClosed", func(t *testing.T) {
		// いったん開いてすぐ Close した ent.Client の UserClient を渡す
		cli := enttest.Open(t, "sqlite3", "file:findbyid4?mode=memory&cache=shared&_fk=1")
		require.NoError(t, cli.Close())

		r := newRepoWithEntUserClient(t, cli.User)
		_, err := r.FindByID(ctx, 1)
		require.Error(t, err, "closed DB should cause internal error")
	})
}
