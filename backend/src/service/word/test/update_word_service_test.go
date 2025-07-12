package word_service_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUpdateWord(t *testing.T) {
	// 	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	// 	defer client.Close()

	// 	clientWrapper := infrastructure.NewAppClient(client)

	// 	wordService := word_service.NewWordService(clientWrapper)
	// 	ctx := context.Background()

	// 	// 品詞データのシードを実行
	// 	partsOfSpeech := []string{"名詞", "代名詞", "動詞", "形容詞", "副詞",
	// 		"助動詞", "前置詞", "冠詞", "間投詞", "接続詞"}

	// 	for _, name := range partsOfSpeech {
	// 		exists, _ := client.PartOfSpeech.Query().Where(partofspeech.Name(name)).Exist(ctx)
	// 		if !exists {
	// 			client.PartOfSpeech.Create().
	// 				SetName(name).
	// 				Save(ctx)
	// 		}
	// 	}

	// 	// ユーザー作成 (管理者と非管理者)
	// 	adminUser, err := client.User.Create().
	// 		SetName("Admin User").
	// 		SetEmail("admin@example.com").
	// 		SetPassword("password").
	// 		SetIsAdmin(true).
	// 		Save(ctx)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, adminUser)

	// 	nonAdminUser, err := client.User.Create().
	// 		SetName("Non-Admin User").
	// 		SetEmail("nonadmin@example.com").
	// 		SetPassword("password").
	// 		SetIsAdmin(false).
	// 		Save(ctx)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, nonAdminUser)

	// 	userID := 1
	// 	word1 := client.Word.Create().
	// 		SetName("apple").
	// 		SetRegistrationCount(5).
	// 		SaveX(ctx)

	// 	client.Word.Create().
	// 		SetName("banana").
	// 		SetRegistrationCount(3).
	// 		SaveX(ctx)

	// 	client.Word.Create().
	// 		SetName("orange").
	// 		SetRegistrationCount(0).
	// 		SaveX(ctx)

	// 	word2 := client.Word.Create().
	// 		SetName("approach").
	// 		SetRegistrationCount(0).
	// 		SaveX(ctx)

	// 	wordInfo1 := client.WordInfo.Create().
	// 		SetWordID(word1.ID).
	// 		SetPartOfSpeechID(1).
	// 		SaveX(ctx)

	// 	client.JapaneseMean.Create().
	// 		SetWordInfoID(wordInfo1.ID).
	// 		SetName("りんご").
	// 		SaveX(ctx)

	// 	wordInfo2 := client.WordInfo.Create().
	// 		SetWordID(word2.ID).
	// 		SetPartOfSpeechID(1).
	// 		SaveX(ctx)
	// 	wordInfo3 := client.WordInfo.Create().
	// 		SetWordID(word2.ID).
	// 		SetPartOfSpeechID(3).
	// 		SaveX(ctx)

	// 	client.JapaneseMean.Create().
	// 		SetWordInfoID(wordInfo2.ID).
	// 		SetName("方法").
	// 		SaveX(ctx)
	// 	client.JapaneseMean.Create().
	// 		SetWordInfoID(wordInfo3.ID).
	// 		SetName("近づく").
	// 		SaveX(ctx)

	// 	client.RegisteredWord.Create().
	// 		SetWordID(word1.ID).
	// 		SetUserID(userID).
	// 		SetIsActive(true).
	// 		SetAttentionLevel(2).
	// 		SetQuizCount(10).
	// 		SetCorrectCount(5).
	// 		SetMemo("memo").
	// 		SaveX(ctx)

	// 	// 正常なリクエストデータ
	// 	japaneseMean := []models.JapaneseMean{
	// 		{ID: 1, Name: "りんご"},
	// 	}
	// 	wordInfos := []models.WordInfo{
	// 		{ID: 1, PartOfSpeechID: 1, JapaneseMeans: japaneseMean}, // 名詞
	// 	}
	// 	reqData := &models.UpdateWordRequest{
	// 		ID:        word1.ID,
	// 		Name:      "apple",
	// 		WordInfos: wordInfos,
	// 		UserID:    1, // 存在するユーザーID
	// 	}

	// 	t.Run("Success", func(t *testing.T) {
	// 		// テスト実行
	// 		createdWord, err := wordService.UpdateWord(ctx, reqData)
	// 		assert.NoError(t, err)
	// 		assert.NotNil(t, createdWord)
	// 		assert.Equal(t, 1, createdWord.ID)
	// 		assert.Equal(t, "apple", createdWord.Name)
	// 		assert.Equal(t, "word 'apple' updated successfully", createdWord.Message)
	// 	})

	// 	t.Run("NonAdmin_error", func(t *testing.T) {
	// 		// 2. 管理者以外のユーザーでリクエスト
	// 		reqData = &models.UpdateWordRequest{
	// 			ID:        word1.ID,
	// 			Name:      "test",
	// 			UserID:    nonAdminUser.ID, // 非管理者ユーザー
	// 			WordInfos: []models.WordInfo{},
	// 		}

	// 		_, err = wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Equal(t, word_service.ErrUnauthorized, err)

	// 	})
	// 	t.Run("DuplicateWord_Errors", func(t *testing.T) {
	// 		// 3. 既存の単語がある場合
	// 		// _, err = client.Word.Create().SetName("duplicate").Save(ctx)
	// 		// assert.NoError(t, err)

	// 		reqData = &models.UpdateWordRequest{
	// 			ID:     word2.ID,
	// 			Name:   "orange", // 同じ名前の単語
	// 			UserID: adminUser.ID,
	// 			WordInfos: []models.WordInfo{
	// 				{ID: 2, PartOfSpeechID: 3, JapaneseMeans: []models.JapaneseMean{{ID: 2, Name: "テスト"}}},
	// 				{ID: 3, PartOfSpeechID: 4, JapaneseMeans: []models.JapaneseMean{{ID: 3, Name: "テスト2"}}},
	// 			},
	// 		}

	// 		_, err = wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Equal(t, word_service.ErrWordExists, err)

	// 	})
	// 	t.Run("UpdateWordInfo_Errors", func(t *testing.T) {
	// 		// 4. WordInfo更新エラー
	// 		reqData = &models.UpdateWordRequest{
	// 			ID:     word2.ID,
	// 			Name:   "newword",
	// 			UserID: adminUser.ID,
	// 			WordInfos: []models.WordInfo{
	// 				{ID: 2, PartOfSpeechID: 3, JapaneseMeans: []models.JapaneseMean{{ID: 2, Name: "テスト"}}},
	// 				{ID: 3, PartOfSpeechID: 999, JapaneseMeans: []models.JapaneseMean{{ID: 3, Name: "無効な品詞"}}}, // 存在しない品詞ID
	// 			},
	// 		}

	// 		_, err = wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "failed to update word info")

	// 	})
	// 	t.Run("CreateJapaneseMean_Errors", func(t *testing.T) {
	// 		// 5. JapaneseMean更新エラー
	// 		reqData = &models.UpdateWordRequest{
	// 			ID:     word2.ID,
	// 			Name:   "validword",
	// 			UserID: adminUser.ID,
	// 			WordInfos: []models.WordInfo{
	// 				{ID: 2, PartOfSpeechID: 3, JapaneseMeans: []models.JapaneseMean{{ID: 2, Name: "テスト"}}},
	// 				{ID: 3, PartOfSpeechID: 4, JapaneseMeans: []models.JapaneseMean{{ID: 3, Name: ""}}},
	// 			},
	// 		}

	// 		_, err = wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "failed to update japanese mean")
	// 	})
	// 	t.Run("WordNotFound_error", func(t *testing.T) {
	// 		// 存在しない単語IDを指定
	// 		reqData := &models.UpdateWordRequest{
	// 			ID:        999, // 存在しないID
	// 			Name:      "nonexistent",
	// 			UserID:    adminUser.ID, // 管理者ユーザー
	// 			WordInfos: []models.WordInfo{},
	// 		}

	// 		_, err := wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Equal(t, "word not found", err.Error())
	// 	})

	// 	t.Run("DuplicateWordName_error", func(t *testing.T) {
	// 		// 既に存在する単語名を指定
	// 		reqData := &models.UpdateWordRequest{
	// 			ID:        word1.ID,
	// 			Name:      "banana",     // 重複する単語名
	// 			UserID:    adminUser.ID, // 管理者ユーザー
	// 			WordInfos: []models.WordInfo{},
	// 		}

	// 		_, err := wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Equal(t, word_service.ErrWordExists, err)
	// 	})

	// 	t.Run("InvalidPartOfSpeechID_error", func(t *testing.T) {
	// 		// 存在しない PartOfSpeechID を指定
	// 		reqData := &models.UpdateWordRequest{
	// 			ID:   word1.ID,
	// 			Name: "apple",
	// 			WordInfos: []models.WordInfo{
	// 				{ID: wordInfo1.ID, PartOfSpeechID: 999}, // 存在しないID
	// 			},
	// 			UserID: adminUser.ID, // 管理者ユーザー
	// 		}

	// 		_, err := wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "failed to update word info")
	// 	})

	// 	t.Run("JapaneseMeanUpdate_error", func(t *testing.T) {
	// 		// 不正な JapaneseMean.ID を指定
	// 		reqData := &models.UpdateWordRequest{
	// 			ID:   word1.ID,
	// 			Name: "apple",
	// 			WordInfos: []models.WordInfo{
	// 				{
	// 					ID:             wordInfo1.ID,
	// 					PartOfSpeechID: 1,
	// 					JapaneseMeans: []models.JapaneseMean{
	// 						{ID: 999, Name: "不明な意味"}, // 存在しないID
	// 					},
	// 				},
	// 			},
	// 			UserID: adminUser.ID, // 管理者ユーザー
	// 		}

	// 		result, err := wordService.UpdateWord(ctx, reqData)
	// 		logrus.Info(result, "reslut")
	// 		logrus.Info(err, "reslut")

	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "failed to fetch japanese mean")
	// 	})

	// 	t.Run("TransactionRollback_error", func(t *testing.T) {
	// 		// 内部でエラーを発生させてトランザクションがロールバックされることを確認
	// 		reqData := &models.UpdateWordRequest{
	// 			ID:   word1.ID,
	// 			Name: "apple",
	// 			WordInfos: []models.WordInfo{
	// 				{
	// 					ID:             wordInfo1.ID,
	// 					PartOfSpeechID: 1,
	// 					JapaneseMeans: []models.JapaneseMean{
	// 						{ID: japaneseMean[0].ID, Name: ""}, // 無効な名前
	// 					},
	// 				},
	// 			},
	// 			UserID: adminUser.ID, // 管理者ユーザー
	// 		}

	// 		_, err := wordService.UpdateWord(ctx, reqData)
	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "failed to update japanese mean")

	//		// トランザクションがロールバックされていることを確認
	//		updatedWord := client.Word.GetX(ctx, word1.ID)
	//		assert.Equal(t, "apple", updatedWord.Name) // 名前が変更されていない
	//	})
}
