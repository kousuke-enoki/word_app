package user_service_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	"word_app/backend/src/infrastructure"
	user_service "word_app/backend/src/service/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestEntUserClient_FindUserByEmail(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	clientWrapper := infrastructure.NewAppClient(client)

	usrClient := user_service.NewEntUserClient(clientWrapper)
	ctx := context.Background()

	// 初期データ
	email := "test@example.com"
	name := "Test User"
	password := "password"

	// ユーザー作成
	_, err := usrClient.CreateUser(ctx, email, name, password)
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		// 正常系: ユーザーを見つける
		foundUser, err := usrClient.FindUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, email, foundUser.Email)
		assert.Equal(t, name, foundUser.Name)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// 異常系: 存在しないメールアドレス
		nonExistentEmail := "nonexistent@example.com"
		foundUser, err := usrClient.FindUserByEmail(ctx, nonExistentEmail)
		assert.Nil(t, foundUser)
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrUserNotFound)
	})

	t.Run("DatabaseFailure", func(t *testing.T) {
		// 異常系: データベースクライアントを強制的に閉じる
		client.Close()
		_, err := usrClient.FindUserByEmail(ctx, email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrUserNotFound)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		// 異常系: 空のメールアドレス
		emptyEmail := ""
		foundUser, err := usrClient.FindUserByEmail(ctx, emptyEmail)
		assert.Nil(t, foundUser)
		assert.Error(t, err)
	})
}
