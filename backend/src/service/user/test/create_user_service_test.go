package user_service_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	"word_app/backend/src/infrastructure"
	user_service "word_app/backend/src/service/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEntUserClient_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() {
		if cerr := client.Close(); cerr != nil {
			logrus.Error("failed to close ent test client:", cerr)
		}
	}()

	clientWrapper := infrastructure.NewAppClient(client)

	usrClient := user_service.NewEntUserClient(clientWrapper)

	ctx := context.Background()

	// 共通の入力データ
	email := "test@example.com"
	name := "Test User"
	password := "password"

	t.Run("Success", func(t *testing.T) {
		createdUser, err := usrClient.Create(ctx, email, name, password)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, email, createdUser.Email)
		assert.Equal(t, name, createdUser.Name)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		entUsrClient := user_service.NewEntUserClient(clientWrapper)
		// ctx := context.Background()

		// Successで登録したuserと同じメールアドレスで再度作成
		_, err := entUsrClient.Create(ctx, email, "Another User", "anotherpassword")
		assert.ErrorIs(t, err, user_service.ErrDuplicateEmail)
	})

	t.Run("DatabaseFailure", func(t *testing.T) {
		badClient := enttest.Open(t, "sqlite3", "file:bad?mode=memory&cache=shared&_fk=1")
		badWrapper := infrastructure.NewAppClient(badClient)
		badSvc := user_service.NewEntUserClient(badWrapper)

		// ここで先に閉じる（重要）
		_ = badClient.Close()

		_, err := badSvc.Create(ctx, "new@example.com", "New User", "newpassword")
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrDatabaseFailure)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		// 無効なメールアドレス（空文字列）
		_, err := usrClient.Create(ctx, "", name, password)
		assert.Error(t, err)

		// 無効な名前（空文字列）
		_, err = usrClient.Create(ctx, email, "", password)
		assert.Error(t, err)

		// 無効なパスワード（空文字列）
		_, err = usrClient.Create(ctx, email, name, "")
		assert.Error(t, err)
	})
}
