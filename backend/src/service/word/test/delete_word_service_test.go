package word_service_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	"word_app/backend/ent/partofspeech"
	"word_app/backend/src/infrastructure"
	"word_app/backend/src/models"
	word_service "word_app/backend/src/service/word"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestDeleteWord(t *testing.T) {
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
	client.User.Create().
		SetName("Admin User").
		SetEmail("admin@example.com").
		SetPassword("password").
		SetIsAdmin(true).
		Save(ctx)
	// assert.NoError(t, err)
	// assert.NotNil(t, adminUser)

	client.User.Create().
		SetName("Non-Admin User").
		SetEmail("nonadmin@example.com").
		SetPassword("password").
		SetIsAdmin(false).
		Save(ctx)
	// assert.NoError(t, err)
	// assert.NotNil(t, nonAdminUser)

	// 正常なリクエストデータ
	japaneseMean := []models.JapaneseMean{
		{Name: "テスト"},
	}
	wordInfos := []models.WordInfo{
		{PartOfSpeechID: 1, JapaneseMeans: japaneseMean}, // 名詞
	}
	createReqData := &models.CreateWordRequest{
		Name:      "test",
		WordInfos: wordInfos,
		UserID:    1, // 存在するユーザーID
	}

	deleteReqData := &models.DeleteWordRequest{
		WordID: 1,
		UserID: 1, // 存在するユーザーID
	}

	t.Run("Success", func(t *testing.T) {
		// テスト実行
		wordService.CreateWord(ctx, createReqData)
		deletedWord, err := wordService.DeleteWord(ctx, deleteReqData)
		assert.NoError(t, err)
		assert.NotNil(t, deletedWord)
		assert.Equal(t, "test", deletedWord.Name)
	})
	t.Run("Error: Non-admin user attempts to delete", func(t *testing.T) {
		// 非管理者が削除を試みる
		deleteReqData.UserID = 2 // 非管理者のユーザーID
		_, err := wordService.DeleteWord(ctx, deleteReqData)
		assert.Error(t, err)
		assert.Equal(t, word_service.ErrUnauthorized, err)
	})

	t.Run("Error: Word not found", func(t *testing.T) {
		// 存在しない単語IDで削除を試みる
		deleteReqData.UserID = 1   // 管理者のユーザーID
		deleteReqData.WordID = 999 // 存在しない単語ID
		_, err := wordService.DeleteWord(ctx, deleteReqData)
		assert.Error(t, err)
		assert.Equal(t, word_service.ErrWordNotFound, err)
	})
}
