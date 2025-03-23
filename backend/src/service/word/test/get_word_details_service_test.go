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

func TestGetWordDetails_Success(t *testing.T) {
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

	client.Word.Create().
		SetName("orange").
		SetRegistrationCount(0).
		SaveX(ctx)

	word2 := client.Word.Create().
		SetName("approach").
		SetRegistrationCount(0).
		SaveX(ctx)

	wordInfo1 := client.WordInfo.Create().
		SetWordID(word1.ID).
		SetPartOfSpeechID(1).
		SaveX(ctx)

	client.JapaneseMean.Create().
		SetWordInfoID(wordInfo1.ID).
		SetName("りんご").
		SaveX(ctx)

	wordInfo2 := client.WordInfo.Create().
		SetWordID(word2.ID).
		SetPartOfSpeechID(1).
		SaveX(ctx)
	wordInfo3 := client.WordInfo.Create().
		SetWordID(word2.ID).
		SetPartOfSpeechID(3).
		SaveX(ctx)

	client.JapaneseMean.Create().
		SetWordInfoID(wordInfo2.ID).
		SetName("方法").
		SaveX(ctx)
	client.JapaneseMean.Create().
		SetWordInfoID(wordInfo3.ID).
		SetName("近づく").
		SaveX(ctx)

	client.RegisteredWord.Create().
		SetWordID(word1.ID).
		SetUserID(userID).
		SetIsActive(true).
		SetAttentionLevel(2).
		SetTestCount(10).
		SetCheckCount(5).
		SetMemo("memo").
		SaveX(ctx)

	// client.RegisteredWord.Create().
	// 	SetWordID(word2.ID).
	// 	SetUserID(userID).
	// 	SetIsActive(false).
	// 	SetAttentionLevel(1).
	// 	SetTestCount(3).
	// 	SetCheckCount(1).
	// 	SaveX(ctx)

	// Create request
	req := &models.WordShowRequest{
		WordID: 1,
		UserID: userID,
	}

	t.Run("Success_registeredWord", func(t *testing.T) {
		resp, err := wordService.GetWordDetails(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.ID)
		assert.Equal(t, "apple", resp.Name)
		assert.Equal(t, 5, resp.RegistrationCount)
		assert.Equal(t, true, resp.IsRegistered)
		assert.Equal(t, 2, resp.AttentionLevel)
		assert.Equal(t, 10, resp.TestCount)
		assert.Equal(t, 5, resp.CheckCount)
		assert.Equal(t, "memo", resp.Memo)

		wordInfos := resp.WordInfos
		assert.Equal(t, 1, len(wordInfos))
		assert.Equal(t, 1, wordInfos[0].ID)
		assert.Equal(t, 1, wordInfos[0].PartOfSpeechID)

		japaneseMeans := wordInfos[0].JapaneseMeans
		assert.Equal(t, 1, len(japaneseMeans))
		assert.Equal(t, 1, japaneseMeans[0].ID)
		assert.Equal(t, "りんご", japaneseMeans[0].Name)
	})
	t.Run("Success_notRegisteredWord", func(t *testing.T) {
		// Create request
		req := &models.WordShowRequest{
			WordID: 4,
			UserID: userID,
		}
		resp, err := wordService.GetWordDetails(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 4, resp.ID)
		assert.Equal(t, "approach", resp.Name)
		assert.Equal(t, 0, resp.RegistrationCount)
		assert.Equal(t, false, resp.IsRegistered)
		assert.Equal(t, 0, resp.AttentionLevel)
		assert.Equal(t, 0, resp.TestCount)
		assert.Equal(t, 0, resp.CheckCount)
		assert.Equal(t, "", resp.Memo)

		wordInfos := resp.WordInfos
		assert.Equal(t, 2, len(wordInfos))
		assert.Equal(t, 2, wordInfos[0].ID)
		assert.Equal(t, 1, wordInfos[0].PartOfSpeechID)
		assert.Equal(t, 3, wordInfos[1].ID)
		assert.Equal(t, 3, wordInfos[1].PartOfSpeechID)

		japaneseMeans := wordInfos[0].JapaneseMeans
		assert.Equal(t, 1, len(japaneseMeans))
		assert.Equal(t, 2, japaneseMeans[0].ID)
		assert.Equal(t, "方法", japaneseMeans[0].Name)
		japaneseMeans2 := wordInfos[1].JapaneseMeans
		assert.Equal(t, 1, len(japaneseMeans2))
		assert.Equal(t, 3, japaneseMeans2[0].ID)
		assert.Equal(t, "近づく", japaneseMeans2[0].Name)
	})
	t.Run("NoResults", func(t *testing.T) {
		req := &models.WordShowRequest{
			WordID: 100,
			UserID: userID,
		}

		_, err := wordService.GetWordDetails(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, errors.New("failed to fetch word details"), err)
	})
	t.Run("InvalidUserID", func(t *testing.T) {
		req := &models.WordShowRequest{
			UserID: 9999, // 存在しないUserID
			WordID: 1,
		}

		_, err := wordService.GetWordDetails(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, word_service.ErrUserNotFound, err)
	})
	// t.Run("ExcludeWordsFromOtherUsers", func(t *testing.T) {
	// 	otherUserID := 2
	// 	client.RegisteredWord.Create().
	// 		SetWordID(word1.ID).
	// 		SetUserID(otherUserID).
	// 		SetIsActive(true).
	// 		SetAttentionLevel(3).
	// 		SetTestCount(7).
	// 		SetCheckCount(4).
	// 		SaveX(ctx)

	// 	req := &models.WordShowRequest{
	// 		UserID: userID,
	// 		Search: "",
	// 		SortBy: "name",
	// 		Order:  "asc",
	// 		Page:   1,
	// 		Limit:  10,
	// 	}

	// 	resp, err := wordService.GetWordDetails(ctx, req)

	// 	// Assertions
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, resp)
	// 	assert.Equal(t, 2, len(resp.Words))
	// })

}
