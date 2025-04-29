package word_service_test

import (
	"context"
	"errors"
	"testing"

	"word_app/backend/ent/enttest"
	"word_app/backend/ent/partofspeech"
	"word_app/backend/src/infrastructure"
	"word_app/backend/src/models"
	word_service "word_app/backend/src/service/word"

	"github.com/stretchr/testify/assert"
)

func TestSaveMemo(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	clientWrapper := infrastructure.NewAppClient(client)

	wordService := word_service.NewWordService(clientWrapper)
	ctx := context.Background()

	// 品詞データのシードを実行
	partsOfSpeech := []string{"名詞", "代名詞", "動詞", "形容詞", "副詞",
		"助動詞", "前置詞", "冠詞", "間投詞", "接続詞"}
	for _, name := range partsOfSpeech {
		exists, _ := client.PartOfSpeech.Query().Where(partofspeech.Name(name)).Exist(ctx)
		if !exists {
			client.PartOfSpeech.Create().
				SetName(name).
				Save(ctx)
		}
	}

	// ユーザー作成 (管理者と非管理者)
	adminUser, err := client.User.Create().
		SetName("Admin User").
		SetEmail("admin@example.com").
		SetPassword("password").
		SetIsAdmin(true).
		Save(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, adminUser)

	nonAdminUser, err := client.User.Create().
		SetName("Non-Admin User").
		SetEmail("nonadmin@example.com").
		SetPassword("password").
		SetIsAdmin(false).
		Save(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, nonAdminUser)

	userID := 1
	word1 := client.Word.Create().
		SetName("apple").
		SetRegistrationCount(5).
		SaveX(ctx)

	client.Word.Create().
		SetName("banana").
		SetRegistrationCount(3).
		SaveX(ctx)

	client.RegisteredWord.Create().
		SetWordID(word1.ID).
		SetUserID(userID).
		SetIsActive(true).
		SetAttentionLevel(2).
		SetQuizCount(10).
		SetCorrectCount(5).
		SaveX(ctx)

	t.Run("Success", func(t *testing.T) {
		// Create request
		req := &models.SaveMemoRequest{
			WordID: 1,
			UserID: 1,
			Memo:   "memo sample",
		}
		resp, err := wordService.SaveMemo(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// assert.Equal(t, 2, len(resp.Words))
		assert.Equal(t, "apple", resp.Name)
		assert.Equal(t, "memo sample", resp.Memo)
		assert.Equal(t, "RegisteredWord memo updated", resp.Message)
	})
	// t.Run("Fail_create_RegisteredWord", func(t *testing.T) {
	// 	// モックの準備
	// 	mockEntClient := mocks.NewEntClientInterface(t)
	// 	mockRegisteredWordClient := new(mocks.RegisteredWord)
	// 	mockTx := new(ent.Tx)

	// 	// テスト対象のサービスの作成
	// 	service := word_service.WordServiceImpl{
	// 		Client: mockEntClient,
	// 	}

	// 	// テストデータの準備
	// 	ctx := context.Background()
	// 	userID := 1
	// 	wordID := 123
	// 	wordName := "test_word"
	// 	req := &models.RegisterWordRequest{
	// 		UserID:       userID,
	// 		WordID:       wordID,
	// 		IsRegistered: true,
	// 	}

	// 	// モックの期待値を設定
	// 	mockEntClient.On("Tx", ctx).Return(mockTx, nil)
	// 	mockEntClient.On("User").Return(nil) // 省略可能：Userクエリが必要な場合ここで定義
	// 	mockEntClient.On("Word").Return(nil) // 省略可能：Wordクエリが必要な場合ここで定義
	// 	mockEntClient.On("RegisteredWord").Return(mockRegisteredWordClient)

	// 	mockRegisteredWordClient.On("Create").
	// 		Return(mockRegisteredWordClient)
	// 	mockRegisteredWordClient.On("SetUserID", userID).
	// 		Return(mockRegisteredWordClient)
	// 	mockRegisteredWordClient.On("SetWordID", wordID).
	// 		Return(mockRegisteredWordClient)
	// 	mockRegisteredWordClient.On("SetIsActive", true).
	// 		Return(mockRegisteredWordClient)
	// 	mockRegisteredWordClient.On("Save", ctx).
	// 		Return(&ent.RegisteredWord{
	// 			ID:       1,
	// 			UserID:   userID,
	// 			WordID:   wordID,
	// 			IsActive: true,
	// 		}, nil)

	// 	// テスト実行
	// 	res, err := service.RegisterWords(ctx, req)

	// 	// 結果の検証
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, res)
	// 	assert.Equal(t, wordName, res.Name)
	// 	assert.Equal(t, true, res.IsRegistered)

	// 	// モックの呼び出しを検証
	// 	mockEntClient.AssertExpectations(t)
	// 	mockRegisteredWordClient.AssertExpectations(t)
	// })

	t.Run("no_word", func(t *testing.T) {
		// Create request
		req := &models.SaveMemoRequest{
			WordID: 1000,
			UserID: 1,
			Memo:   "memo sample",
		}
		resp, err := wordService.SaveMemo(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.New("failed to fetch word"), err)
	})
}
