package setting_test

import (
	"context"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3" // SQLite ドライバ登録
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_app/backend/ent/enttest"

	"word_app/backend/ent"
	"word_app/backend/ent/userconfig"
	"word_app/backend/src/domain"
	setrepo "word_app/backend/src/infrastructure/repository/setting"
	si "word_app/backend/src/interfaces/service_interfaces"
)

/* -------------------------------------------------------------------------- */
/*                               factory helper                               */
/* -------------------------------------------------------------------------- */

func newRepo(t *testing.T) (si.EntClientInterface, *setrepo.EntUserConfigRepo) {
	cli := enttest.Open(t, "sqlite3", "file:memdb?mode=memory&_fk=1")
	return entAdapter{cli}, setrepo.NewEntUserConfigRepo(entAdapter{cli})
}

/* --------------------------- GetByUserID ---------------------------------- */

func TestEntUserConfigRepo_GetByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("存在するユーザ → domain マッピング", func(t *testing.T) {
		cli, repo := newRepo(t)

		// fk を満たす user を作成
		// u, err := cli.User().Create().Save(ctx)
		u, err := cli.User().
			Create().
			SetEmail("alpha@mail.com").
			SetName("seed").
			SetPassword("Password123$").
			Save(ctx)
		require.NoError(t, err)

		_, err = cli.UserConfig().
			Create().
			SetUserID(u.ID).
			SetIsDarkMode(true).
			Save(ctx)
		require.NoError(t, err)

		got, err := repo.GetByUserID(ctx, u.ID)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.True(t, got.IsDarkMode)
	})

	t.Run("存在しなければ ent.IsNotFound", func(t *testing.T) {
		_, repo := newRepo(t)
		_, err := repo.GetByUserID(ctx, 999)
		assert.True(t, ent.IsNotFound(err))
	})
}

/* ------------------------------ Upsert ------------------------------------ */

func TestEntUserConfigRepo_Upsert(t *testing.T) {
	ctx := context.Background()

	t.Run("Create 新規挿入", func(t *testing.T) {
		cli, repo := newRepo(t)
		// u, _ := cli.User().Create().Save(ctx)
		u, _ := cli.User().
			Create().
			SetEmail("alpha@mail.com").
			SetName("seed").
			SetPassword("Password123$").
			Save(ctx)

		in := &domain.UserConfig{UserID: u.ID, IsDarkMode: true}
		got, err := repo.Upsert(ctx, in)
		assert.NoError(t, err)
		assert.True(t, got.IsDarkMode)

		cnt, _ := cli.UserConfig().Query().Count(ctx)
		assert.Equal(t, 1, cnt)
	})

	t.Run("Update 既存更新", func(t *testing.T) {
		cli, repo := newRepo(t)
		// u, _ := cli.User().Create().Save(ctx)
		u, _ := cli.User().
			Create().
			SetEmail("alpha@mail.com").
			SetName("seed").
			SetPassword("Password123$").
			Save(ctx)

		// 既存行（false）
		_, _ = cli.UserConfig().
			Create().
			SetUserID(u.ID).
			SetIsDarkMode(false).
			Save(ctx)

		in := &domain.UserConfig{UserID: u.ID, IsDarkMode: true}
		got, err := repo.Upsert(ctx, in)
		assert.NoError(t, err)
		assert.True(t, got.IsDarkMode)

		dbRec, _ := cli.UserConfig().Query().Where(userconfig.UserID(u.ID)).Only(ctx)
		assert.True(t, dbRec.IsDarkMode)
	})

	t.Run("Query エラーを透過 (ctx cancel)", func(t *testing.T) {
		cli, repo := newRepo(t)
		// u, _ := cli.User().Create().Save(ctx)
		u, _ := cli.User().
			Create().
			SetEmail("alpha@mail.com").
			SetName("seed").
			SetPassword("Password123$").
			Save(ctx)

		cctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := repo.Upsert(cctx, &domain.UserConfig{UserID: u.ID})
		assert.True(t, errors.Is(err, context.Canceled))
	})
}
