package user_service_test

import (
	"context"
	"fmt"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"
	"word_app/backend/src/infrastructure"
	"word_app/backend/src/models"
	user_service "word_app/backend/src/service/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func sptr(s string) *string { return &s }

func seedUser(t *testing.T, cli *ent.Client, u ent.User) *ent.User {
	t.Helper()
	mut := cli.User.Create().
		SetEmail(u.Email).
		SetName(u.Name).
		SetIsAdmin(u.IsAdmin).
		SetIsRoot(u.IsRoot).
		SetIsTest(u.IsTest)
	if u.Password != nil {
		mut = mut.SetPassword(*u.Password)
	}
	got, err := mut.Save(context.Background())
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return got
}

func seedExternalAuth(t *testing.T, cli *ent.Client, user *ent.User, provider string) {
	t.Helper()
	_, err := cli.ExternalAuth.Create().
		SetProvider(provider).
		SetProviderUserID("sub-" + user.Email).
		SetUser(user).
		Save(context.Background())
	if err != nil {
		t.Fatalf("seed externalauth: %v", err)
	}
}

func TestEntUserClient_GetUsers(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:getusers?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() {
		if cerr := client.Close(); cerr != nil {
			logrus.Error("failed to close ent test client:", cerr)
		}
	})
	wrapper := infrastructure.NewAppClient(client)
	svc := user_service.NewEntUserClient(wrapper)

	ctx := context.Background()

	// 呼び出しユーザー（管理者）
	admin := seedUser(t, client, ent.User{
		Email:   "admin@example.com",
		Name:    "Admin",
		IsAdmin: true,
	})

	t.Run("AdminGate_Unauthorized_and_DBFailure", func(t *testing.T) {
		normal := seedUser(t, client, ent.User{
			Email: "u@example.com", Name: "User", IsAdmin: false,
		})

		_, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: normal.ID, Search: "", SortBy: "name", Order: "asc", Page: 1, Limit: 10,
		})
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrUnauthorized)

		_, err = svc.GetUsers(ctx, &models.UserListRequest{
			UserID: 999999, Search: "", SortBy: "name", Order: "asc", Page: 1, Limit: 10,
		})
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrDatabaseFailure)

		// DB クローズで失敗再現
		bad := enttest.Open(t, "sqlite3", "file:badgetusers?mode=memory&cache=shared&_fk=1")
		badWrapper := infrastructure.NewAppClient(bad)
		badSvc := user_service.NewEntUserClient(badWrapper)
		_ = bad.Close()
		_, err = badSvc.GetUsers(ctx, &models.UserListRequest{
			UserID: 1, Search: "", SortBy: "name", Order: "asc", Page: 1, Limit: 10,
		})
		assert.Error(t, err)
	})

	t.Run("Search_Filter_by_Name_and_Email", func(t *testing.T) {
		seedUser(t, client, ent.User{Email: "bob@example.com", Name: "Bobby"})
		seedUser(t, client, ent.User{Email: "alice@wonder.land", Name: "Alice W"})
		seedUser(t, client, ent.User{Email: "carol@sample.com", Name: "Carol"})

		resp, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: "Ali", SortBy: "name", Order: "asc", Page: 1, Limit: 10,
		})
		assert.NoError(t, err)
		if assert.Len(t, resp.Users, 1) {
			assert.Equal(t, "Alice W", resp.Users[0].Name)
		}

		resp, err = svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: "bob@", SortBy: "name", Order: "asc", Page: 1, Limit: 10,
		})
		assert.NoError(t, err)
		if assert.Len(t, resp.Users, 1) {
			assert.Equal(t, "Bobby", resp.Users[0].Name)
		}
	})

	t.Run("Pagination_TotalPages", func(t *testing.T) {
		// 追加データ（25件）: 3〜20文字の名前 & ユニーク email
		for i := 0; i < 25; i++ {
			seedUser(t, client, ent.User{
				Email: fmt.Sprintf("pg-user%02d@test.local", i),
				Name:  fmt.Sprintf("User%02d", i), // 5〜7文字でバリデーションOK
			})
		}
		resp, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: "pg-user", SortBy: "name", Order: "asc", Page: 2, Limit: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, resp.TotalPages)
		assert.Len(t, resp.Users, 10)
	})

	t.Run("Sort_by_Email_AscDesc", func(t *testing.T) {
		// このケース専用トークンを付与（Search で絞る）
		tok := "em-sort-1"
		seedUser(t, client, ent.User{Email: "a+" + tok + "@x.com", Name: "Ace"})
		seedUser(t, client, ent.User{Email: "m+" + tok + "@x.com", Name: "Mid"})
		seedUser(t, client, ent.User{Email: "z+" + tok + "@x.com", Name: "Zed"})

		asc, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: tok, SortBy: "email", Order: "asc", Page: 1, Limit: 10,
		})
		assert.NoError(t, err)
		if assert.Len(t, asc.Users, 3) {
			assert.Equal(t, []string{"a+" + tok + "@x.com", "m+" + tok + "@x.com", "z+" + tok + "@x.com"},
				[]string{asc.Users[0].Email, asc.Users[1].Email, asc.Users[2].Email})
		}

		desc, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: tok, SortBy: "email", Order: "desc", Page: 1, Limit: 10,
		})
		assert.NoError(t, err)
		if assert.Len(t, desc.Users, 3) {
			assert.Equal(t, []string{"z+" + tok + "@x.com", "m+" + tok + "@x.com", "a+" + tok + "@x.com"},
				[]string{desc.Users[0].Email, desc.Users[1].Email, desc.Users[2].Email})
		}
	})

	t.Run("Sort_by_Role_AscDesc", func(t *testing.T) {
		// このケース専用トークンを付与（Search で絞る）
		tok := "role-sort-1"
		seedUser(t, client, ent.User{Email: "r+" + tok + "@x.com", Name: "Root", IsRoot: true})
		seedUser(t, client, ent.User{Email: "a+" + tok + "@x.com", Name: "Admin", IsAdmin: true})
		seedUser(t, client, ent.User{Email: "u+" + tok + "@x.com", Name: "User"})
		seedUser(t, client, ent.User{Email: "t+" + tok + "@x.com", Name: "TestUser", IsTest: true})

		asc, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: tok, SortBy: "role", Order: "asc", Page: 1, Limit: 10,
		})
		assert.NoError(t, err)
		if assert.Len(t, asc.Users, 4) {
			var names []string
			for _, u := range asc.Users {
				names = append(names, u.Name)
			}
			assert.Equal(t, []string{"Root", "Admin", "User", "TestUser"}, names)
		}

		desc, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: tok, SortBy: "role", Order: "desc", Page: 1, Limit: 10,
		})
		assert.NoError(t, err)
		if assert.Len(t, desc.Users, 4) {
			var names []string
			for _, u := range desc.Users {
				names = append(names, u.Name)
			}
			assert.Equal(t, []string{"TestUser", "User", "Admin", "Root"}, names)
		}
	})

	t.Run("IsSettedPassword_and_IsLine", func(t *testing.T) {
		tok := "pw-line-1"
		n := seedUser(t, client, ent.User{Email: "n+" + tok + "@x.com", Name: "NilPass", Password: nil})
		e := seedUser(t, client, ent.User{Email: "e+" + tok + "@x.com", Name: "EmptyPass", Password: sptr("")})
		h := seedUser(t, client, ent.User{Email: "h+" + tok + "@x.com", Name: "HashedPass", Password: sptr("hash")})

		l1 := seedUser(t, client, ent.User{Email: "l1+" + tok + "@x.com", Name: "Line1"})
		l2 := seedUser(t, client, ent.User{Email: "l2+" + tok + "@x.com", Name: "Line2"})
		gh := seedUser(t, client, ent.User{Email: "g+" + tok + "@x.com", Name: "GitHubOnly"})
		na := seedUser(t, client, ent.User{Email: "no+" + tok + "@x.com", Name: "NoAuth"})
		_ = gh
		seedExternalAuth(t, client, l1, "line")
		seedExternalAuth(t, client, l2, "LiNe")
		seedExternalAuth(t, client, gh, "github")

		resp, err := svc.GetUsers(ctx, &models.UserListRequest{
			UserID: admin.ID, Search: tok, SortBy: "name", Order: "asc", Page: 1, Limit: 100,
		})
		assert.NoError(t, err)

		pw := map[string]bool{}
		ln := map[string]bool{}
		for _, u := range resp.Users {
			pw[u.Name] = u.IsSettedPassword
			ln[u.Name] = u.IsLine
		}

		assert.Equal(t, false, pw["NilPass"])
		assert.Equal(t, false, pw["EmptyPass"])
		assert.Equal(t, true, pw["HashedPass"])

		assert.Equal(t, true, ln["Line1"])
		assert.Equal(t, true, ln["Line2"])
		assert.Equal(t, false, ln["GitHubOnly"])
		assert.Equal(t, false, ln["NoAuth"])

		_, _, _, _ = n, e, h, na
	})
}
