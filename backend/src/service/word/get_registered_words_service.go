package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"
)

// word_list
func (s *WordServiceImpl) GetRegisteredWords(ctx context.Context, AllWordListRequest *models.AllWordListRequest) (*models.AllWordListResponse, error) {
	query := s.client.Word.Query()
	userID := AllWordListRequest.UserID
	search := AllWordListRequest.Search
	order := AllWordListRequest.Order
	page := AllWordListRequest.Page
	limit := AllWordListRequest.Limit

	// 検索条件の追加
	query = addSearchRegisteredWordsFilter(query, search)

	// ページネーション機能
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Wordに紐づくデータを取得 (WordInfoとRegisteredWordを含める)
	query = query.WithWordInfos(func(wiQuery *ent.WordInfoQuery) {
		wiQuery.WithJapaneseMeans().WithPartOfSpeech()
	}).WithRegisteredWords(func(rwQuery *ent.RegisteredWordQuery) {
		rwQuery.Where(registeredword.UserID(userID))
	})

	// ページネーション前のフィルタリング用クエリを作成
	filteredQuery := s.client.Word.Query().Where(
		word.HasRegisteredWordsWith(
			registeredword.UserID(userID),
			registeredword.IsActive(true),
		),
	)

	// フィルタリング後のデータ総数を取得
	totalCount, err := filteredQuery.Count(ctx)
	if err != nil {
		return nil, errors.New("failed to count filtered words")
	}

	// 現在のページのデータを取得するためのクエリ設定
	query = query.Where(
		word.HasRegisteredWordsWith(
			registeredword.UserID(userID),
			registeredword.IsActive(true),
		),
	)
	// ソート条件の設定
	if order == "asc" {
		query = query.Order(ent.Asc(word.FieldName))
	} else {
		query = query.Order(ent.Desc(word.FieldName))
	}

	// クエリ実行
	entWords, err := query.All(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch words")
	}

	words := convertEntRegisteredWordsToResponse(entWords)

	// 総ページ数を計算
	totalPages := (totalCount + limit - 1) / limit

	response := &models.AllWordListResponse{
		Words:      words,
		TotalPages: totalPages,
	}
	return response, nil
}

// 検索条件の追加
func addSearchRegisteredWordsFilter(query *ent.WordQuery, search string) *ent.WordQuery {
	if search != "" {
		query = query.Where(word.NameContains(search))
	}
	return query
}

// エンティティからレスポンス形式に変換
func convertEntRegisteredWordsToResponse(entWords []*ent.Word) []models.Word {
	words := make([]models.Word, len(entWords))

	for i, entWord := range entWords {
		// WordInfoの変換
		wordInfos := make([]models.WordInfo, len(entWord.Edges.WordInfos))
		for j, wordInfo := range entWord.Edges.WordInfos {
			partOfSpeech := models.PartOfSpeech{
				ID:   wordInfo.Edges.PartOfSpeech.ID,
				Name: wordInfo.Edges.PartOfSpeech.Name,
			}
			japaneseMeans := make([]models.JapaneseMean, len(wordInfo.Edges.JapaneseMeans))
			for k, mean := range wordInfo.Edges.JapaneseMeans {
				japaneseMeans[k] = models.JapaneseMean{
					ID:   mean.ID,
					Name: mean.Name,
				}
			}
			wordInfos[j] = models.WordInfo{
				ID:            wordInfo.ID,
				PartOfSpeech:  partOfSpeech,
				JapaneseMeans: japaneseMeans,
			}
		}

		// RegisteredWordのデータを設定
		isRegistered := false
		attentionLevel, testCount, checkCount := 0, 0, 0
		if len(entWord.Edges.RegisteredWords) == 1 {
			registeredWord := entWord.Edges.RegisteredWords[0]
			isRegistered = registeredWord.IsActive
			attentionLevel = registeredWord.AttentionLevel
			testCount = registeredWord.TestCount
			checkCount = registeredWord.CheckCount
		}

		words[i] = models.Word{
			ID:                entWord.ID,
			Name:              entWord.Name,
			RegistrationCount: entWord.RegistrationCount,
			WordInfos:         wordInfos,
			IsRegistered:      isRegistered,
			AttentionLevel:    attentionLevel,
			TestCount:         testCount,
			CheckCount:        checkCount,
		}
	}
	return words
}
