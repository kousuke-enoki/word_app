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

func TestEntUserClient_FindUserByID(t *testing.T) {
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
	createdUser, err := usrClient.CreateUser(ctx, email, name, password)
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		// 正常系: IDでユーザーを見つける
		foundUser, err := usrClient.FindUserByID(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, email, foundUser.Email)
		assert.Equal(t, name, foundUser.Name)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// 異常系: 存在しないID
		nonExistentID := createdUser.ID + 1
		foundUser, err := usrClient.FindUserByID(ctx, nonExistentID)
		assert.Nil(t, foundUser)
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrUserNotFound)
	})

	t.Run("DatabaseFailure", func(t *testing.T) {
		// 異常系: データベース接続エラー
		client.Close() // クライアントを閉じて強制的にエラーを発生させる
		foundUser, err := usrClient.FindUserByID(ctx, createdUser.ID)
		assert.Nil(t, foundUser)
		assert.Error(t, err)
	})

	t.Run("InvalidInput_NegativeID", func(t *testing.T) {
		// 異常系: 負のID
		invalidID := -1
		foundUser, err := usrClient.FindUserByID(ctx, invalidID)
		assert.Nil(t, foundUser)
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrUserNotFound)
	})

	t.Run("InvalidInput_ZeroID", func(t *testing.T) {
		// 異常系: IDが0の場合
		invalidID := 0
		foundUser, err := usrClient.FindUserByID(ctx, invalidID)
		assert.Nil(t, foundUser)
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_service.ErrUserNotFound)
	})
}
