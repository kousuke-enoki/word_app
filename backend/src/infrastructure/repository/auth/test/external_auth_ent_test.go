// infrastructure/repository/external_auth_ent_test.go
package auth_test

import (
	"context"
	"testing"

	"word_app/backend/ent"
	"word_app/backend/ent/enttest"
	"word_app/backend/ent/externalauth"

	domain "word_app/backend/src/domain"
	infra "word_app/backend/src/infrastructure"
	repo "word_app/backend/src/infrastructure/repository/auth"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestEntExtAuthRepo_Create(t *testing.T) {
	// ❶ メモリ DB
	client := enttest.Open(t, "sqlite3", "file:memdb?mode=memory&_fk=1")
	defer func() {
		if cerr := client.Close(); cerr != nil {
			logrus.Error("failed to close ent test client:", cerr)
		}
	}()
	ctx := context.Background()

	// ❷ ユーザーを 1 行だけ作成
	u, err := client.User.
		Create().
		SetEmail("foo@example.com").
		SetName("Foo").
		SetPassword("hashed").
		Save(ctx)
	require.NoError(t, err)

	repo := repo.NewEntExtAuthRepo(infra.NewAppClient(client))

	valid := &domain.ExternalAuth{
		UserID:         u.ID,
		Provider:       "google",
		ProviderUserID: "google-123",
	}

	t.Run("success", func(t *testing.T) {
		require.NoError(t, repo.Create(ctx, valid))

		// DB に本当に入ったか確認
		_, err := client.ExternalAuth.
			Query().
			Where(
				externalauth.ProviderEQ(valid.Provider),
				externalauth.ProviderUserIDEQ(valid.ProviderUserID),
			).
			Only(ctx)
		require.NoError(t, err)
	})

	t.Run("duplicate (unique key)", func(t *testing.T) {
		err := repo.Create(ctx, valid) // 同じ provider+userID
		require.True(t, ent.IsConstraintError(err))
	})

	t.Run("fk violation (unknown user)", func(t *testing.T) {
		bad := &domain.ExternalAuth{
			UserID:         9999,
			Provider:       "google",
			ProviderUserID: "other-id",
		}
		err := repo.Create(ctx, bad)
		require.True(t, ent.IsConstraintError(err)) // FK 失敗
	})

	t.Run("db failure (client closed)", func(t *testing.T) {
		defer func() {
			if cerr := client.Close(); cerr != nil {
				logrus.Error("failed to close ent test client:", cerr)
			}
		}()
		err := repo.Create(ctx, valid)
		require.Error(t, err)
	})
}
