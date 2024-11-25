package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"
)

// all_word_list
func (s *WordServiceImpl) GetWords(ctx context.Context, search string, sortBy string, order string, page int, limit int) ([]models.Word, int, int, error) {
	query := s.client.Word.Query()

	// 検索機能
	if search != "" {
		query = query.Where(word.NameContains(search))
	}

	// 総レコード数を取得
	totalCount, err := query.Count(ctx)
	if err != nil {
		return nil, 0, 0, errors.New("failed to count words")
	}

	// ページネーション機能
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Wordに紐づくデータを取得
	query = query.WithWordInfos(func(wiQuery *ent.WordInfoQuery) {
		wiQuery.WithJapaneseMeans().WithPartOfSpeech()
	})

	// ソート機能
	switch sortBy {
	case "name":
		if order == "asc" {
			query = query.Order(ent.Asc(word.FieldName))
		} else {
			query = query.Order(ent.Desc(word.FieldName))
		}
	default:
		if order == "asc" {
			query = query.Order(ent.Asc(sortBy))
		} else {
			query = query.Order(ent.Desc(sortBy))
		}
	}

	// クエリ実行
	entWords, err := query.All(ctx)
	if err != nil {
		return nil, 0, 0, errors.New("failed to fetch words")
	}

	// entの型からレスポンス用の型に変換
	words := make([]models.Word, len(entWords))
	for i, entWord := range entWords {
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
		words[i] = models.Word{
			ID:        entWord.ID,
			Name:      entWord.Name,
			WordInfos: wordInfos,
		}
	}

	// 総ページ数を計算
	totalPages := (totalCount + limit - 1) / limit

	return words, totalCount, totalPages, nil
}
