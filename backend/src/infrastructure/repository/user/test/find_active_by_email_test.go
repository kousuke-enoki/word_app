// backend/src/infrastructure/repository/user/find_active_by_email_test.go
package user

import (
	"context"
	"testing"
	"time"

	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/domain"
	repo "word_app/backend/src/infrastructure/repository/user"
	serviceinterfaces "word_app/backend/src/mocks/service_interfaces"
)

// 共有ヘルパ：本物 ent.Client（in-memory SQLite）を裏で使いつつ、
// リポジトリは公開 API だけを使う黒箱スタイルで検証する。
func newRepoForTests(t *testing.T) (*repo.EntUserRepo, *serviceinterfaces.MockEntClientInterface) {
	t.Helper()

	cli := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = cli.Close() })

	m := serviceinterfaces.NewMockEntClientInterface(t)
	// リポジトリの各メソッドは内部で client.User() を呼ぶので、それだけ返せば良い
	m.EXPECT().User().Return(cli.User)

	return repo.NewEntUserRepo(m), m
}

func TestEntUserRepo_FindActiveByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_FoundActive", func(t *testing.T) {
		r, _ := newRepoForTests(t)

		// 公開 API だけ使ってシード
		u, err := r.Create(ctx, &domain.User{
			Name:     "ActiveUser",
			Email:    ptr("active@example.com"),
			Password: "hashed",
			IsAdmin:  true,
			IsRoot:   false,
			IsTest:   false,
		})
		require.NoError(t, err)

		got, err := r.FindActiveByEmail(ctx, "active@example.com")
		require.NoError(t, err)
		require.NotNil(t, got)
		logrus.Info(got)
		require.Equal(t, u.ID, got.ID)
		require.Equal(t, "ActiveUser", got.Name)
		require.NotNil(t, got.Email)
		require.Equal(t, "active@example.com", *got.Email)
		require.Equal(t, "hashed", got.Password)
		require.Equal(t, true, got.IsAdmin)
		require.Equal(t, false, got.IsRoot)
		require.Equal(t, false, got.IsTest)
		require.WithinDuration(t, u.CreatedAt, got.CreatedAt, time.Second)
		require.WithinDuration(t, u.UpdatedAt, got.UpdatedAt, time.Second)
	})

	t.Run("NotFound_NoSuchEmail", func(t *testing.T) {
		r, _ := newRepoForTests(t)

		got, err := r.FindActiveByEmail(ctx, "nobody@example.com")
		require.Error(t, err)
		require.Nil(t, got)
		// repoerr の種類を厳密に見るならここで IsNotFound などをアサート
	})

	t.Run("NotFound_DeletedUser", func(t *testing.T) {
		r, _ := newRepoForTests(t)

		created, err := r.Create(ctx, &domain.User{
			Name:     "DeletedUser",
			Email:    ptr("deleted@example.com"),
			Password: "x",
		})
		require.NoError(t, err)

		// 公開 API で論理削除（DeletedAt をセット）
		err = r.SoftDeleteByID(ctx, created.ID, time.Now())
		require.NoError(t, err)

		got, err := r.FindActiveByEmail(ctx, "deleted@example.com")
		require.Error(t, err) // DeletedAt != nil なのでヒットしない
		require.Nil(t, got)
	})

	t.Run("Internal_DBClosed", func(t *testing.T) {
		// 通常のリポジトリを作る
		r, mock := newRepoForTests(t)

		// 「閉じた ent.Client の UserClient」を返させて内部エラーを強制
		tmp := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		require.NoError(t, tmp.Close())

		// 期待値を差し替え
		mock.ExpectedCalls = nil
		mock.EXPECT().User().Return(tmp.User)

		got, err := r.FindActiveByEmail(ctx, "any@example.com")
		require.Error(t, err)
		require.Nil(t, got)
	})
}
