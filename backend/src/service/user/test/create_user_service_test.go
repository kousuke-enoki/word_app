package user_service_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	user_service "word_app/backend/src/service/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestEntUserClient_CreateUser(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	usrClient := user_service.NewEntUserClient(client)
	ctx := context.Background()

	// 共通の入力データ
	email := "test@example.com"
	name := "Test User"
	password := "password"

	t.Run("Success", func(t *testing.T) {
		createdUser, err := usrClient.CreateUser(ctx, email, name, password)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, email, createdUser.Email)
		assert.Equal(t, name, createdUser.Name)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		usrClient := user_service.NewEntUserClient(client)
		ctx := context.Background()

		// Successで登録したuserと同じメールアドレスで再度作成
		_, err := usrClient.CreateUser(ctx, email, "Another User", "anotherpassword")
		assert.ErrorIs(t, err, user_service.ErrDuplicateEmail)
	})

	t.Run("DatabaseFailure", func(t *testing.T) {
		// DBクライアントを強制的に無効化してエラーを発生させる
		client.Close()
		_, err := usrClient.CreateUser(ctx, "new@example.com", "New User", "newpassword")
		assert.ErrorIs(t, err, user_service.ErrDatabaseFailure)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		// 無効なメールアドレス（空文字列）
		_, err := usrClient.CreateUser(ctx, "", name, password)
		assert.Error(t, err)

		// 無効な名前（空文字列）
		_, err = usrClient.CreateUser(ctx, email, "", password)
		assert.Error(t, err)

		// 無効なパスワード（空文字列）
		_, err = usrClient.CreateUser(ctx, email, name, "")
		assert.Error(t, err)
	})
}
