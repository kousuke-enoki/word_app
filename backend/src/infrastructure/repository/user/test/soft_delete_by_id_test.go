// backend/src/infrastructure/repository/user/soft_delete_by_id_test.go
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

func newSoftDeleteRepo(t *testing.T, uc *ent.UserClient) *repo.EntUserRepo {
	t.Helper()
	m := serviceinterfaces.NewMockEntClientInterface(t)
	m.EXPECT().User().Return(uc)
	return repo.NewEntUserRepo(m)
}

// func ptr[T any](v T) *T { return &v }

func TestEntUserRepo_SoftDeleteByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_UpdateDeletedAt_ForActiveUser", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:softdel1?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		u, err := cli.User.Create().
			SetName("Active").
			SetPassword("hashed").
			SetEmail("active@example.com").
			Save(ctx)
		require.NoError(t, err)

		r := newSoftDeleteRepo(t, cli.User)

		ts := time.Now().UTC().Truncate(time.Second)
		err = r.SoftDeleteByID(ctx, u.ID, ts)
		require.NoError(t, err)

		got, err := cli.User.Get(ctx, u.ID)
		require.NoError(t, err)
		require.NotNil(t, got.DeletedAt)
		require.WithinDuration(t, ts, *got.DeletedAt, time.Second)
	})

	t.Run("Success_UpdateDeletedAt_OverridesExisting", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:softdel2?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		initial := time.Now().Add(-24 * time.Hour).UTC().Truncate(time.Second)
		u, err := cli.User.Create().
			SetName("AlreadyDeleted").
			SetPassword("hashed").
			SetEmail("deleted@example.com").
			SetNillableDeletedAt(&initial).
			Save(ctx)
		require.NoError(t, err)
		require.NotNil(t, u.DeletedAt)

		r := newSoftDeleteRepo(t, cli.User)

		newTs := time.Now().UTC().Truncate(time.Second)
		err = r.SoftDeleteByID(ctx, u.ID, newTs)
		require.NoError(t, err)

		got, err := cli.User.Get(ctx, u.ID)
		require.NoError(t, err)
		require.NotNil(t, got.DeletedAt)
		require.WithinDuration(t, newTs, *got.DeletedAt, time.Second)
		require.NotEqual(t, initial, *got.DeletedAt, "DeletedAt should be overwritten")
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:softdel3?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })

		r := newSoftDeleteRepo(t, cli.User)
		err := r.SoftDeleteByID(ctx, 999999, time.Now())
		require.Error(t, err)
	})

	t.Run("Error_DBClosed_Internal", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:softdel4?mode=memory&cache=shared&_fk=1")
		require.NoError(t, cli.Close()) // クローズして内部エラーを誘発

		r := newSoftDeleteRepo(t, cli.User)
		err := r.SoftDeleteByID(ctx, 1, time.Now())
		require.Error(t, err)
	})
}
