package word_service_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateWord(t *testing.T) {
	// client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	// defer client.Close()

	// clientWrapper := infrastructure.NewAppClient(client)

	// wordService := word_service.NewWordService(clientWrapper)
	// ctx := context.Background()

	// // 品詞データのシードを実行
	// partsOfSpeech := []string{"名詞", "代名詞", "動詞", "形容詞", "副詞",
	// 	"助動詞", "前置詞", "冠詞", "間投詞", "接続詞"}

	// for _, name := range partsOfSpeech {
	// 	exists, _ := client.PartOfSpeech.Query().Where(partofspeech.Name(name)).Exist(ctx)
	// 	if !exists {
	// 		client.PartOfSpeech.Create().
	// 			SetName(name).
	// 			Save(ctx)
	// 	}
	// }

	// // ユーザー作成 (管理者と非管理者)
	// adminUser, err := client.User.Create().
	// 	SetName("Admin User").
	// 	SetEmail("admin@example.com").
	// 	SetPassword("password").
	// 	SetIsAdmin(true).
	// 	Save(ctx)
	// assert.NoError(t, err)
	// assert.NotNil(t, adminUser)

	// nonAdminUser, err := client.User.Create().
	// 	SetName("Non-Admin User").
	// 	SetEmail("nonadmin@example.com").
	// 	SetPassword("password").
	// 	SetIsAdmin(false).
	// 	Save(ctx)
	// assert.NoError(t, err)
	// assert.NotNil(t, nonAdminUser)

	// // 正常なリクエストデータ
	// japaneseMean := []models.JapaneseMean{
	// 	{Name: "テスト"},
	// }
	// wordInfos := []models.WordInfo{
	// 	{PartOfSpeechID: 1, JapaneseMeans: japaneseMean}, // 名詞
	// }
	// reqData := &models.CreateWordRequest{
	// 	Name:      "test",
	// 	WordInfos: wordInfos,
	// 	UserID:    1, // 存在するユーザーID
	// }

	// t.Run("Success", func(t *testing.T) {
	// 	// テスト実行
	// 	createdWord, err := wordService.CreateWord(ctx, reqData)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, createdWord)
	// 	assert.Equal(t, "test", createdWord.Name)
	// })

	// t.Run("NonAdmin_error", func(t *testing.T) {
	// 	// 2. 管理者以外のユーザーでリクエスト
	// 	reqData = &models.CreateWordRequest{
	// 		Name:      "test",
	// 		UserID:    nonAdminUser.ID, // 非管理者ユーザー
	// 		WordInfos: []models.WordInfo{},
	// 	}

	// 	_, err = wordService.CreateWord(ctx, reqData)
	// 	assert.Error(t, err)
	// 	assert.Equal(t, word_service.ErrUnauthorized, err)

	// })
	// t.Run("DuplicateWord_Errors", func(t *testing.T) {
	// 	// 3. 既存の単語がある場合
	// 	_, err = client.Word.Create().SetName("duplicate").Save(ctx)
	// 	assert.NoError(t, err)

	// 	reqData = &models.CreateWordRequest{
	// 		Name:   "duplicate", // 同じ名前の単語
	// 		UserID: adminUser.ID,
	// 		WordInfos: []models.WordInfo{
	// 			{PartOfSpeechID: 1, JapaneseMeans: []models.JapaneseMean{{Name: "ダブリ"}}},
	// 		},
	// 	}

	// 	_, err = wordService.CreateWord(ctx, reqData)
	// 	assert.Error(t, err)
	// 	assert.Equal(t, word_service.ErrWordExists, err)

	// })
	// t.Run("CreateWordInfo_Errors", func(t *testing.T) {
	// 	// 4. WordInfo作成エラー
	// 	reqData = &models.CreateWordRequest{
	// 		Name:   "newword",
	// 		UserID: adminUser.ID,
	// 		WordInfos: []models.WordInfo{
	// 			{PartOfSpeechID: 999, JapaneseMeans: []models.JapaneseMean{{Name: "無効な品詞"}}}, // 存在しない品詞ID
	// 		},
	// 	}

	// 	_, err = wordService.CreateWord(ctx, reqData)
	// 	assert.Error(t, err)
	// 	assert.Equal(t, word_service.ErrCreateWordInfo, err) // 比較を修正

	// })
	// t.Run("CreateJapaneseMean_Errors", func(t *testing.T) {
	// 	// 5. JapaneseMean作成エラー
	// 	reqData = &models.CreateWordRequest{
	// 		Name:   "validword",
	// 		UserID: adminUser.ID,
	// 		WordInfos: []models.WordInfo{
	// 			{
	// 				PartOfSpeechID: 1,
	// 				JapaneseMeans: []models.JapaneseMean{
	// 					{Name: ""}, // 空文字の日本語意味
	// 				},
	// 			},
	// 		},
	// 	}

	// 	_, err = wordService.CreateWord(ctx, reqData)
	// 	assert.Error(t, err)
	// 	assert.Equal(t, word_service.ErrCreateJapaneseMean, err)
	// })
}
