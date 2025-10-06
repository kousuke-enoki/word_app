// backend/src/infrastructure/repository/user/create_test.go
package user

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/domain"
	"word_app/backend/src/infrastructure/repository/user"
	serviceinterfaces "word_app/backend/src/mocks/service_interfaces"
)

func newRepoWithRealEnt(t *testing.T) (*user.EntUserRepo, *serviceinterfaces.MockEntClientInterface) {
	t.Helper()

	// in-memory sqlite の本物 ent.Client
	cli := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = cli.Close() })

	// EntClientInterface のモックには、必要なクライアントだけ返させる
	m := serviceinterfaces.NewMockEntClientInterface(t)
	m.EXPECT().User().Return(cli.User)

	// Create は User() しか呼ばないので、他は設定不要
	return user.NewEntUserRepo(m), m
}

func ptr[T any](v T) *T { return &v }

// --- Success: Email あり ---
func TestEntUserCreateRepo(t *testing.T) {
	repo, _ := newRepoWithRealEnt(t)
	ctx := context.Background()

	t.Run("Create_Success_WithEmail", func(t *testing.T) {
		in := &domain.User{
			Name:     "Alice",
			Email:    ptr("alice@example.com"),
			Password: "hashed-password",
		}

		got, err := repo.Create(ctx, in)
		require.NoError(t, err, "create should succeed")
		require.NotZero(t, got.ID, "ID should be set by ent")
		require.Equal(t, "Alice", got.Name)
		require.NotNil(t, got.Email)
		require.Equal(t, "alice@example.com", *got.Email)
		require.NotZero(t, got.CreatedAt)
		require.NotZero(t, got.UpdatedAt)

		// 関数は引数の u を上書きして返す実装のため、同一参照であることも見ておく
		require.Same(t, in, got)
	})

	t.Run("Create_Success_WithoutEmail", func(t *testing.T) {
		in := &domain.User{
			Name:     "Bob",
			Email:    nil, // ← SetNillableEmail(nil) を通す
			Password: "hashed-password",
		}

		got, err := repo.Create(ctx, in)
		require.NoError(t, err)
		require.NotZero(t, got.ID)
		require.Equal(t, "Bob", got.Name)
		require.Nil(t, got.Email, "email should remain nil when not provided")
		require.NotZero(t, got.CreatedAt)
		require.NotZero(t, got.UpdatedAt)
	})

	// --- Failure: Email 重複（unique 制約違反想定）---
	// func TestEntUserRepo_Create_Fail_DuplicateEmail(t *testing.T) {
	t.Run("Create_Fail_DuplicateEmail", func(t *testing.T) {
		// 1件目は成功
		first := &domain.User{
			Name:     "Carol",
			Email:    ptr("dup@example.com"),
			Password: "x",
		}
		_, err := repo.Create(ctx, first)
		require.NoError(t, err)

		// 2件目：同じ Email で重複。ent 側の unique 制約が有効ならエラーになる。
		second := &domain.User{
			Name:     "Carol 2nd",
			Email:    ptr("dup@example.com"),
			Password: "y",
		}
		_, err = repo.Create(ctx, second)
		require.Error(t, err, "duplicate email should return error")

		// ここでエラー型/メッセージの粒度を合わせたい場合は、
		// repoerr.FromEnt の変換仕様に合わせてチェックを加えてください。
		// 例）require.True(t, repoerr.IsConflict(err))
	})

	// --- Failure: DB 障害（クライアント close 後に Create）---
	t.Run("Create_Fail_DBClosed", func(t *testing.T) {
		// まず通常通り repo を作る
		repo, mock := newRepoWithRealEnt(t)
		ctx := context.Background()

		// 直前に ent.Client を Close してから、モックに返させる UserClient を差し替え
		// enttest.Open は t.Cleanup で Close 済みだが、明示的に DB を閉じた UserClient を渡したいので
		// 新規に開いてすぐ Close → その UserClient を返すようにセットする
		tmp := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		require.NoError(t, tmp.Close())

		// 以降 User() 呼び出しで「閉じた DB の UserClient」を返す
		mock.ExpectedCalls = nil // 既存の期待値をクリアして差し替え
		mock.EXPECT().User().Return(tmp.User)

		in := &domain.User{
			Name:     "Dave",
			Email:    ptr("dave@example.com"),
			Password: "pw",
		}

		_, err := repo.Create(ctx, in)
		require.Error(t, err, "closed DB should cause an error")
	})
}
