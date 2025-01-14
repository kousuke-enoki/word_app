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

func TestRegisteredWordCount(t *testing.T) {
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
		SetAdmin(true).
		Save(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, adminUser)

	nonAdminUser, err := client.User.Create().
		SetName("Non-Admin User").
		SetEmail("nonadmin@example.com").
		SetPassword("password").
		SetAdmin(false).
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
		SetTestCount(10).
		SetCheckCount(5).
		SaveX(ctx)

	t.Run("Success_countAdd", func(t *testing.T) {
		// Create request
		req := &models.RegisteredWordCountRequest{
			WordID:       1,
			IsRegistered: true,
		}
		resp, err := wordService.RegisteredWordCount(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// assert.Equal(t, 2, len(resp.Words))
		assert.Equal(t, 6, resp.RegistrationCount)
	})

	t.Run("Success_countDown", func(t *testing.T) {
		// Create request
		req := &models.RegisteredWordCountRequest{
			WordID:       1,
			IsRegistered: false,
		}
		resp, err := wordService.RegisteredWordCount(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 5, resp.RegistrationCount)
	})

	t.Run("Fail_countAdd", func(t *testing.T) {
		req := &models.RegisteredWordCountRequest{
			WordID:       1000,
			IsRegistered: true,
		}

		resp, err := wordService.RegisteredWordCount(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.New("failed to fetch word"), err)
	})
	t.Run("Fail_countDown", func(t *testing.T) {
		req := &models.RegisteredWordCountRequest{
			WordID:       9999,
			IsRegistered: false,
		}

		resp, err := wordService.RegisteredWordCount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.New("failed to fetch word"), err)
	})
}
