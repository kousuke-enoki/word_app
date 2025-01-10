package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

// all_word_list
func (s *WordServiceImpl) GetWords(ctx context.Context, AllWordListRequest *models.AllWordListRequest) (*models.AllWordListResponse, error) {
	query := s.client.Word.Query()
	userID := AllWordListRequest.UserID
	search := AllWordListRequest.Search
	sortBy := AllWordListRequest.SortBy
	order := AllWordListRequest.Order
	page := AllWordListRequest.Page
	limit := AllWordListRequest.Limit

	// user存在チェック
	_, err := s.client.User.Get(ctx, userID)
	if err != nil {
		logrus.Error(err)
		return nil, ErrUserNotFound
	}

	// 検索条件の追加
	query = addSearchFilter(query, search)

	var totalCount int = 0

	// 総レコード数を取得
	totalCount, err = query.Count(ctx)
	if err != nil {
		return nil, errors.New("failed to count words")
	}

	// ページネーション機能
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Wordに紐づくデータを取得 (WordInfoとRegisteredWordを含める)
	query = query.WithWordInfos(func(wiQuery *ent.WordInfoQuery) {
		wiQuery.WithJapaneseMeans().WithPartOfSpeech()
	}).WithRegisteredWords(func(rwQuery *ent.RegisteredWordQuery) {
		rwQuery.Where(registeredword.UserID(userID))
	})

	// ソート機能
	switch sortBy {
	case "name":
		if order == "asc" {
			query = query.Order(ent.Asc(word.FieldName))
		} else {
			query = query.Order(ent.Desc(word.FieldName))
		}
	case "registrationCount":
		if order == "asc" {
			query = query.Order(ent.Asc(word.FieldRegistrationCount))
		} else {
			query = query.Order(ent.Desc(word.FieldRegistrationCount))
		}
	}

	// クエリ実行
	entWords, err := query.All(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch words")
	}

	// エンティティからレスポンス形式への変換
	words := convertEntWordsToResponse(entWords)

	// 総ページ数を計算
	totalPages := (totalCount + limit - 1) / limit

	response := &models.AllWordListResponse{
		Words:      words,
		TotalPages: totalPages,
	}
	return response, nil
}

// 検索条件の追加
func addSearchFilter(query *ent.WordQuery, search string) *ent.WordQuery {
	if search != "" {
		query = query.Where(word.NameContains(search))
	}
	return query
}

// エンティティからレスポンス形式に変換
func convertEntWordsToResponse(entWords []*ent.Word) []models.Word {
	words := make([]models.Word, len(entWords))
	for i, entWord := range entWords {
		// WordInfoの変換
		wordInfos := make([]models.WordInfo, len(entWord.Edges.WordInfos))
		for j, wordInfo := range entWord.Edges.WordInfos {
			partOfSpeech := models.PartOfSpeech{
				ID: wordInfo.Edges.PartOfSpeech.ID,
			}
			japaneseMeans := make([]models.JapaneseMean, len(wordInfo.Edges.JapaneseMeans))
			for k, mean := range wordInfo.Edges.JapaneseMeans {
				japaneseMeans[k] = models.JapaneseMean{
					ID:   mean.ID,
					Name: mean.Name,
				}
			}
			wordInfos[j] = models.WordInfo{
				ID:             wordInfo.ID,
				PartOfSpeechID: partOfSpeech.ID,
				JapaneseMeans:  japaneseMeans,
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
