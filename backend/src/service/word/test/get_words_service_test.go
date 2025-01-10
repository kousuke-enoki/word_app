package word_service_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	"word_app/backend/ent/partofspeech"
	"word_app/backend/src/models"
	word_service "word_app/backend/src/service/word"

	"github.com/stretchr/testify/assert"
)

func TestGetWords_Success(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	wordService := word_service.NewWordService(client)
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

	word2 := client.Word.Create().
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

	client.RegisteredWord.Create().
		SetWordID(word2.ID).
		SetUserID(userID).
		SetIsActive(false).
		SetAttentionLevel(1).
		SetTestCount(3).
		SetCheckCount(1).
		SaveX(ctx)

	// Create request
	req := &models.AllWordListRequest{
		UserID: userID,
		Search: "",
		SortBy: "name",
		Order:  "asc",
		Page:   1,
		Limit:  10,
	}

	t.Run("Success", func(t *testing.T) {
		resp, err := wordService.GetWords(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, len(resp.Words))
		assert.Equal(t, "apple", resp.Words[0].Name)
		assert.Equal(t, "banana", resp.Words[1].Name)
		assert.Equal(t, 1, resp.TotalPages)

		// Verify RegisteredWord data
		assert.True(t, resp.Words[0].IsRegistered)
		assert.Equal(t, 2, resp.Words[0].AttentionLevel)
		assert.Equal(t, 10, resp.Words[0].TestCount)
		assert.Equal(t, 5, resp.Words[0].CheckCount)

		assert.False(t, resp.Words[1].IsRegistered)
		assert.Equal(t, 1, resp.Words[1].AttentionLevel)
		assert.Equal(t, 3, resp.Words[1].TestCount)
		assert.Equal(t, 1, resp.Words[1].CheckCount)
	})
	t.Run("NoResults", func(t *testing.T) {
		req := &models.AllWordListRequest{
			UserID: userID,
			Search: "nonexistent",
			SortBy: "name",
			Order:  "asc",
			Page:   1,
			Limit:  10,
		}

		resp, err := wordService.GetWords(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 0, len(resp.Words))
		assert.Equal(t, 0, resp.TotalPages)
	})
	t.Run("InvalidUserID", func(t *testing.T) {
		req := &models.AllWordListRequest{
			UserID: 9999, // 存在しないUserID
			Search: "",
			SortBy: "name",
			Order:  "asc",
			Page:   1,
			Limit:  10,
		}

		_, err := wordService.GetWords(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, word_service.ErrUserNotFound, err)
	})
	t.Run("OutOfRangePage", func(t *testing.T) {
		req := &models.AllWordListRequest{
			UserID: userID,
			Search: "",
			SortBy: "name",
			Order:  "asc",
			Page:   100, // 存在しないページ
			Limit:  10,
		}

		resp, err := wordService.GetWords(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 0, len(resp.Words))
		assert.Equal(t, 1, resp.TotalPages)
	})
	t.Run("SortByNameAsc", func(t *testing.T) {
		req := &models.AllWordListRequest{
			UserID: userID,
			Search: "",
			SortBy: "name",
			Order:  "asc",
			Page:   1,
			Limit:  10,
		}

		resp, err := wordService.GetWords(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "apple", resp.Words[0].Name)
		assert.Equal(t, "banana", resp.Words[1].Name)
	})

	t.Run("SortByNameDesc", func(t *testing.T) {
		req := &models.AllWordListRequest{
			UserID: userID,
			Search: "",
			SortBy: "name",
			Order:  "desc",
			Page:   1,
			Limit:  10,
		}

		resp, err := wordService.GetWords(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "banana", resp.Words[0].Name)
		assert.Equal(t, "apple", resp.Words[1].Name)
	})
	t.Run("ExcludeWordsFromOtherUsers", func(t *testing.T) {
		otherUserID := 2
		client.RegisteredWord.Create().
			SetWordID(word1.ID).
			SetUserID(otherUserID).
			SetIsActive(true).
			SetAttentionLevel(3).
			SetTestCount(7).
			SetCheckCount(4).
			SaveX(ctx)

		req := &models.AllWordListRequest{
			UserID: userID,
			Search: "",
			SortBy: "name",
			Order:  "asc",
			Page:   1,
			Limit:  10,
		}

		resp, err := wordService.GetWords(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, len(resp.Words))
	})

}
