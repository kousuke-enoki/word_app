// backend/src/infrastructure/repository/user/list_users_test.go
package user

import (
	"context"
	"testing"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"word_app/backend/src/domain/repository"
	repo "word_app/backend/src/infrastructure/repository/user"
	serviceinterfaces "word_app/backend/src/mocks/service_interfaces"
)

// func ptr[T any](v T) *T { return &v }

// ListUsers はリポジトリ内で client.User() しか直接呼ばないので、UserClient を返すモックだけあればOK
func newListRepo(t *testing.T, uc *ent.UserClient) *repo.EntUserRepo {
	t.Helper()
	m := serviceinterfaces.NewMockEntClientInterface(t)
	m.EXPECT().User().Return(uc)
	return repo.NewEntUserRepo(m)
}

type seeded struct {
	ID      int
	Name    string
	Email   string
	IsRoot  bool
	IsAdmin bool
	IsTest  bool
	HasLINE bool // 検索には使わないが、データ形状のメモ
	Deleted bool
}

func seedUsersForList(t *testing.T, cli *ent.Client) []seeded {
	t.Helper()
	ctx := context.Background()

	// users:
	//   - Root: Eve (root)
	//   - Admin: Bob (admin)
	//   - Normal: Alice, Carol
	//   - Test: Trent (is_test)
	//   - Deleted: Mallory (除外)
	//   - 外部認証: Alice=LINE, Bob=GitHub(=LINEではない)
	now := time.Now()

	var out []seeded

	// Alice (normal + LINE)
	alice, err := cli.User.Create().
		SetName("Alice").
		SetPassword("hashed").
		SetEmail("alice@example.com").
		Save(ctx)
	require.NoError(t, err)
	_, err = cli.ExternalAuth.Create().
		SetProvider("line").
		SetProviderUserID("alice-line-sub").
		SetUser(alice).
		Save(ctx)
	require.NoError(t, err)
	out = append(out, seeded{ID: alice.ID, Name: "Alice", Email: "alice@example.com", HasLINE: true})

	// Bob (admin, GitHub only)
	bob, err := cli.User.Create().
		SetName("Bob").
		SetPassword("hashed").
		SetIsAdmin(true).
		SetEmail("bob@example.com").
		Save(ctx)
	require.NoError(t, err)
	_, err = cli.ExternalAuth.Create().
		SetProvider("github").
		SetProviderUserID("bob-gh").
		SetUser(bob).
		Save(ctx)
	require.NoError(t, err)
	out = append(out, seeded{ID: bob.ID, Name: "Bob", Email: "bob@example.com", IsAdmin: true})

	// Carol (normal, no auth)
	carol, err := cli.User.Create().
		SetName("Carol").
		SetPassword("hashed").
		SetEmail("carol@example.com").
		Save(ctx)
	require.NoError(t, err)
	out = append(out, seeded{ID: carol.ID, Name: "Carol", Email: "carol@example.com"})

	// Eve (root)
	eve, err := cli.User.Create().
		SetName("Eve").
		SetPassword("hashed").
		SetIsRoot(true).
		SetEmail("eve@example.com").
		Save(ctx)
	require.NoError(t, err)
	out = append(out, seeded{ID: eve.ID, Name: "Eve", Email: "eve@example.com", IsRoot: true})

	// Trent (test user)
	trent, err := cli.User.Create().
		SetName("Trent").
		SetPassword("hashed").
		SetIsTest(true).
		SetEmail("trent@example.com").
		Save(ctx)
	require.NoError(t, err)
	out = append(out, seeded{ID: trent.ID, Name: "Trent", Email: "trent@example.com", IsTest: true})

	// Mallory (deleted -> 除外)
	delAt := now
	mallory, err := cli.User.Create().
		SetName("Mallory").
		SetPassword("hashed").
		SetEmail("mallory@example.com").
		SetNillableDeletedAt(&delAt).
		Save(ctx)
	require.NoError(t, err)
	out = append(out, seeded{ID: mallory.ID, Name: "Mallory", Email: "mallory@example.com", Deleted: true})

	return out
}

func TestEntUserRepo_ListUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_NoFilter_ReturnsOnlyActive_AndTotalCount", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu1?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		all := seedUsersForList(t, cli)

		r := newListRepo(t, cli.User)
		got, err := r.ListUsers(ctx, repository.UserListFilter{
			Offset: 0, Limit: 100,
		})
		require.NoError(t, err)

		// 削除ユーザー(Mallory)は除外 → アクティブは 5 件
		require.Equal(t, 5, got.TotalCount)
		require.Len(t, got.Users, 5)

		// ざっくり存在していることを確認
		names := map[string]bool{}
		for _, u := range got.Users {
			names[u.Name] = true
		}
		require.True(t, names["Alice"])
		require.True(t, names["Bob"])
		require.True(t, names["Carol"])
		require.True(t, names["Eve"])
		require.True(t, names["Trent"])

		_ = all
	})

	t.Run("Success_SearchByName_Substring", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu2?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)

		r := newListRepo(t, cli.User)
		got, err := r.ListUsers(ctx, repository.UserListFilter{
			Search: "ali", // Alice をヒット
			Limit:  50,
		})
		require.NoError(t, err)
		require.Equal(t, 1, got.TotalCount)
		require.Len(t, got.Users, 1)
		require.Equal(t, "Alice", got.Users[0].Name)
	})

	t.Run("Success_SearchByEmail_Substring", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu3?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)

		r := newListRepo(t, cli.User)
		got, err := r.ListUsers(ctx, repository.UserListFilter{
			Search: "bob@", // bob@example.com をヒット
			Limit:  50,
		})
		require.NoError(t, err)
		require.Equal(t, 1, got.TotalCount)
		require.Len(t, got.Users, 1)
		require.Equal(t, "Bob", got.Users[0].Name)
	})

	t.Run("Success_Pagination_OffsetLimit", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu4?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)

		r := newListRepo(t, cli.User)

		// Limit=2 で1ページ目
		p1, err := r.ListUsers(ctx, repository.UserListFilter{
			Offset: 0, Limit: 2,
		})
		require.NoError(t, err)
		require.Equal(t, 5, p1.TotalCount)
		require.Len(t, p1.Users, 2)

		// 2ページ目
		p2, err := r.ListUsers(ctx, repository.UserListFilter{
			Offset: 2, Limit: 2,
		})
		require.NoError(t, err)
		require.Equal(t, 5, p2.TotalCount)
		require.Len(t, p2.Users, 2)

		// 3ページ目（残り1件）
		p3, err := r.ListUsers(ctx, repository.UserListFilter{
			Offset: 4, Limit: 2,
		})
		require.NoError(t, err)
		require.Equal(t, 5, p3.TotalCount)
		require.Len(t, p3.Users, 1)
	})

	t.Run("Success_SortByName_AscDesc", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu5?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)
		r := newListRepo(t, cli.User)

		asc, err := r.ListUsers(ctx, repository.UserListFilter{
			SortBy: "name", Order: "asc", Limit: 100,
		})
		require.NoError(t, err)
		require.Len(t, asc.Users, 5)
		// 期待: Alice, Bob, Carol, Eve, Trent
		gotAsc := []string{asc.Users[0].Name, asc.Users[1].Name, asc.Users[2].Name, asc.Users[3].Name, asc.Users[4].Name}
		require.Equal(t, []string{"Alice", "Bob", "Carol", "Eve", "Trent"}, gotAsc)

		des, err := r.ListUsers(ctx, repository.UserListFilter{
			SortBy: "name", Order: "desc", Limit: 100,
		})
		require.NoError(t, err)
		require.Len(t, des.Users, 5)
		gotDesc := []string{des.Users[0].Name, des.Users[1].Name, des.Users[2].Name, des.Users[3].Name, des.Users[4].Name}
		require.Equal(t, []string{"Trent", "Eve", "Carol", "Bob", "Alice"}, gotDesc)
	})

	t.Run("Success_SortByEmail_AscDesc", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu6?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)
		r := newListRepo(t, cli.User)

		asc, err := r.ListUsers(ctx, repository.UserListFilter{
			SortBy: "email", Order: "asc", Limit: 100,
		})
		require.NoError(t, err)
		require.Len(t, asc.Users, 5)
		gotAsc := []string{
			*asc.Users[0].Email, *asc.Users[1].Email, *asc.Users[2].Email, *asc.Users[3].Email, *asc.Users[4].Email,
		}
		require.Equal(t,
			[]string{"alice@example.com", "bob@example.com", "carol@example.com", "eve@example.com", "trent@example.com"},
			gotAsc,
		)

		des, err := r.ListUsers(ctx, repository.UserListFilter{
			SortBy: "email", Order: "desc", Limit: 100,
		})
		require.NoError(t, err)
		require.Len(t, des.Users, 5)
		gotDesc := []string{
			*des.Users[0].Email, *des.Users[1].Email, *des.Users[2].Email, *des.Users[3].Email, *des.Users[4].Email,
		}
		require.Equal(t,
			[]string{"trent@example.com", "eve@example.com", "carol@example.com", "bob@example.com", "alice@example.com"},
			gotDesc,
		)
	})

	t.Run("Success_SortByRole_Asc_CompositeOrder", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu7?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)
		r := newListRepo(t, cli.User)

		// role asc: ORDER BY IsRoot DESC, IsAdmin DESC, IsTest ASC
		// → Root(Eve) → Admin(Bob) → Normal(non-test)(Alice, Carol) → Test(Trent)
		got, err := r.ListUsers(ctx, repository.UserListFilter{
			SortBy: "role", Order: "asc", Limit: 100,
		})
		require.NoError(t, err)
		require.Len(t, got.Users, 5)
		names := []string{got.Users[0].Name, got.Users[1].Name, got.Users[2].Name, got.Users[3].Name, got.Users[4].Name}
		require.Equal(t, []string{"Eve", "Bob", "Alice", "Carol", "Trent"}, names)
	})

	t.Run("Success_SortByRole_Desc_CompositeOrder", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu8?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)
		r := newListRepo(t, cli.User)

		// role desc: ORDER BY IsRoot ASC, IsAdmin ASC, IsTest DESC
		// → Non-root & Non-admin & Test から先頭（Trent）→ Normal(non-test)(Alice, Carol) → Admin(Bob) → Root(Eve)
		got, err := r.ListUsers(ctx, repository.UserListFilter{
			SortBy: "role", Order: "desc", Limit: 100,
		})
		require.NoError(t, err)
		require.Len(t, got.Users, 5)
		names := []string{got.Users[0].Name, got.Users[1].Name, got.Users[2].Name, got.Users[3].Name, got.Users[4].Name}
		require.Equal(t, []string{"Trent", "Alice", "Carol", "Bob", "Eve"}, names)
	})

	t.Run("Success_WithVariousExternalAuths_NoPanic", func(t *testing.T) {
		// WithExternalAuths の条件が provider=line かつ DeletedAtIsNil。
		// LINE以外(auth=github)が混ざっていても取得自体は成功することを確認。
		cli := enttest.Open(t, "sqlite3", "file:lu9?mode=memory&cache=shared&_fk=1")
		t.Cleanup(func() { _ = cli.Close() })
		_ = seedUsersForList(t, cli)
		r := newListRepo(t, cli.User)

		got, err := r.ListUsers(ctx, repository.UserListFilter{Limit: 100})
		require.NoError(t, err)
		require.Equal(t, 5, got.TotalCount)
		require.Len(t, got.Users, 5)
	})

	t.Run("Internal_DBClosed_ReturnsError", func(t *testing.T) {
		cli := enttest.Open(t, "sqlite3", "file:lu10?mode=memory&cache=shared&_fk=1")
		require.NoError(t, cli.Close())

		r := newListRepo(t, cli.User)
		_, err := r.ListUsers(ctx, repository.UserListFilter{Limit: 10})
		require.Error(t, err)
	})
}
