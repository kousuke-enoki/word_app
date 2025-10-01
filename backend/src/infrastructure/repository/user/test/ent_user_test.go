package user_test

import (
	"context"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"

	domain "word_app/backend/src/domain"
	repo "word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/test"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

/*****************************************************************/

func seedUserWithAuth(t *testing.T, ec *ent.Client, email, provider, sub string) *ent.User {
	ctx := context.Background()
	u, err := ec.User.
		Create().
		SetEmail(email).
		SetName("seed").
		SetPassword("Password123$").
		Save(ctx)
	require.NoError(t, err)

	_, err = ec.ExternalAuth.
		Create().
		SetUserID(u.ID).
		SetProvider(provider).
		SetProviderUserID(sub).
		Save(ctx)
	require.NoError(t, err)
	return u
}

func TestEntUserRepo_FindByProvider(t *testing.T) {
	ec := enttest.Open(t, "sqlite3", "file:memdb1?mode=memory&_fk=1")
	defer func() {
		if cerr := ec.Close(); cerr != nil {
			logrus.Error("close file:", cerr)
		}
	}()
	repo := repo.NewEntUserRepo(test.RealEntClient{Client: ec})
	ctx := context.Background()

	seed := seedUserWithAuth(t, ec, "alpha@mail.com", "google", "g-123")

	t.Run("found", func(t *testing.T) {
		got, err := repo.FindByProvider(ctx, "google", "g-123")
		require.NoError(t, err)
		require.Equal(t, seed.ID, got.ID)
		require.Equal(t, seed.Email, got.Email)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByProvider(ctx, "google", "nope")
		require.True(t, ent.IsNotFound(err))
	})
}

func TestEntUserRepo_Create(t *testing.T) {
	ec := enttest.Open(t, "sqlite3", "file:memdb2?mode=memory&_fk=1&_busy_timeout=10000")
	defer func() {
		if cerr := ec.Close(); cerr != nil {
			logrus.Error("close file:", cerr)
		}
	}()
	r := repo.NewEntUserRepo(test.RealEntClient{Client: ec})
	ctx := context.Background()

	Email := "beta@mail.com"
	u := &domain.User{Email: &Email, Name: "Beta", Password: "pw"}

	t.Run("success", func(t *testing.T) {
		_, err := r.Create(ctx, u)
		require.NoError(t, err)

		// 確認: 両テーブルとも 1 行ずつ存在
		cnt, _ := ec.User.Query().Count(ctx)
		require.Equal(t, 1, cnt)
	})

	t.Run("tx begin failure", func(t *testing.T) {
		defer func() {
			if cerr := ec.Close(); cerr != nil {
				logrus.Error("close file:", cerr)
			}
		}()
		_, err := r.Create(ctx, u)
		require.Error(t, err)
	})
}
