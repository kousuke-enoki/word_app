package setting_test

import (
	"context"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3" // SQLite ドライバ登録
	"github.com/stretchr/testify/assert"

	"word_app/backend/ent/enttest"

	"word_app/backend/ent"
	"word_app/backend/src/domain"
	setrepo "word_app/backend/src/infrastructure/repository/setting"
	si "word_app/backend/src/interfaces/service_interfaces"
)

/* -------------------------------------------------------------------------- */
/*                              adapter for test                              */
/* -------------------------------------------------------------------------- */

// entAdapter は service_interfaces.EntClientInterface を *ent.Client* で満たす薄いラッパー。
// すべて単純委譲なのでテスト用にしか使わない。
type entAdapter struct{ *ent.Client }

func (a entAdapter) EntClient() *ent.Client                    { return a.Client }
func (a entAdapter) ExternalAuth() *ent.ExternalAuthClient     { return a.Client.ExternalAuth }
func (a entAdapter) JapaneseMean() *ent.JapaneseMeanClient     { return a.Client.JapaneseMean }
func (a entAdapter) Quiz() *ent.QuizClient                     { return a.Client.Quiz }
func (a entAdapter) QuizQuestion() *ent.QuizQuestionClient     { return a.Client.QuizQuestion }
func (a entAdapter) RegisteredWord() *ent.RegisteredWordClient { return a.Client.RegisteredWord }
func (a entAdapter) RootConfig() *ent.RootConfigClient         { return a.Client.RootConfig }
func (a entAdapter) Tx(ctx context.Context) (*ent.Tx, error)   { return a.Client.Tx(ctx) }
func (a entAdapter) User() *ent.UserClient                     { return a.Client.User }
func (a entAdapter) UserConfig() *ent.UserConfigClient         { return a.Client.UserConfig }
func (a entAdapter) Word() *ent.WordClient                     { return a.Client.Word }
func (a entAdapter) WordInfo() *ent.WordInfoClient             { return a.Client.WordInfo }

/* -------------------------------------------------------------------------- */
/*                               helper factory                               */
/* -------------------------------------------------------------------------- */

func newRootRepo(t *testing.T) (si.EntClientInterface, *setrepo.EntRootConfigRepo) {
	cli := enttest.Open(t, "sqlite3", "file:memdb?mode=memory&_fk=1")
	return entAdapter{cli}, setrepo.NewEntRootConfigRepo(entAdapter{cli})
}

/* ========================================================================== */
/*                                   tests                                    */
/* ========================================================================== */

func TestEntRootConfigRepo_Get(t *testing.T) {
	ctx := context.Background()

	t.Run("存在すれば domain へマッピング", func(t *testing.T) {
		cli, repo := newRootRepo(t)
		_, _ = cli.RootConfig().
			Create().
			SetEditingPermission("admin").
			SetIsTestUserMode(true).
			SetIsEmailAuthenticationCheck(false).
			SetIsLineAuthentication(false).
			Save(ctx)

		got, err := repo.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "admin", got.EditingPermission)
		assert.True(t, got.IsTestUserMode)
	})

	t.Run("存在しなければ ent.NotFound を透過", func(t *testing.T) {
		_, repo := newRootRepo(t)

		_, err := repo.Get(ctx)
		assert.Error(t, err)
		assert.True(t, ent.IsNotFound(err))
	})
}

/* -------------------------------------------------------------------------- */

func TestEntRootConfigRepo_Upsert(t *testing.T) {
	ctx := context.Background()

	input := &domain.RootConfig{
		EditingPermission:          "admin",
		IsTestUserMode:             true,
		IsEmailAuthenticationCheck: false,
		IsLineAuthentication:       false,
	}

	t.Run("レコードが無い場合は Create", func(t *testing.T) {
		cli, repo := newRootRepo(t)

		got, err := repo.Upsert(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, 1, got.ID)

		// DB に 1 件だけ出来ていることを確認
		cnt, _ := cli.RootConfig().Query().Count(ctx)
		assert.Equal(t, 1, cnt)
	})

	t.Run("既存レコードがあれば Update", func(t *testing.T) {
		cli, repo := newRootRepo(t)

		// 先に既存行を作る
		_, _ = cli.RootConfig().
			Create().
			SetEditingPermission("user").
			SetIsTestUserMode(false).
			SetIsEmailAuthenticationCheck(true).
			SetIsLineAuthentication(true).
			Save(ctx)

		// 更新内容
		updateIn := &domain.RootConfig{
			EditingPermission:          "root",
			IsTestUserMode:             true,
			IsEmailAuthenticationCheck: false,
			IsLineAuthentication:       false,
		}
		got, err := repo.Upsert(ctx, updateIn)
		assert.NoError(t, err)
		assert.Equal(t, "root", got.EditingPermission)

		// DB の値も更新されているか
		dbRec, _ := cli.RootConfig().Query().Only(ctx)
		assert.Equal(t, "root", dbRec.EditingPermission)
		assert.True(t, dbRec.IsTestUserMode)
	})

	t.Run("Query エラー時はそのまま返す", func(t *testing.T) {
		// キャンセル済み context を渡して Query エラーを発生させる
		_, repo := newRootRepo(t)
		cctx, cancel := context.WithCancel(ctx)
		cancel() // immediately cancel

		_, err := repo.Upsert(cctx, input)
		assert.Error(t, err)
		// ent 内部では context.Canceled がそのまま返る
		assert.True(t, errors.Is(err, context.Canceled))
	})

	t.Run("Save エラー時も透過", func(t *testing.T) {
		cli, repo := newRootRepo(t)

		// rootconfig.ID(1) を作ってロック → Save が deadlock 的に失敗させる
		tx, _ := cli.Tx(ctx)
		_, _ = tx.RootConfig.
			Create().
			SetEditingPermission("user").
			SetIsTestUserMode(false).
			SetIsEmailAuthenticationCheck(false).
			SetIsLineAuthentication(false).
			Save(ctx)

		// Tx をコミットせずロックを保持したまま Upsert を呼ぶ
		// SQLite のメモリ DB ではロック競合を投げてくれる
		_, err := repo.Upsert(ctx, input)
		assert.Error(t, err)
		_ = tx.Rollback()
	})
}
