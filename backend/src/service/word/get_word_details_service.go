package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"
)

// word_show
func (s *WordServiceImpl) GetWordDetails(ctx context.Context, wordID int) (*models.WordResponse, error) {
	wordEntity, err := s.client.Word.
		Query().
		Where(word.ID(wordID)).
		WithWordInfos(func(wq *ent.WordInfoQuery) {
			wq.WithJapaneseMeans().WithPartOfSpeech()
		}).
		WithRegisteredWords().
		Only(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch word details")
	}

	// 登録済み情報を取得
	var isRegistered bool
	var testCount, checkCount int
	var memo string

	if len(wordEntity.Edges.RegisteredWords) > 0 {
		registeredWord := wordEntity.Edges.RegisteredWords[0]
		isRegistered = registeredWord.IsActive
		testCount = registeredWord.TestCount
		checkCount = registeredWord.CheckCount
		if registeredWord.Memo != nil {
			memo = *registeredWord.Memo
		}
	}

	// WordInfosを変換
	wordInfos := make([]models.WordInfo, len(wordEntity.Edges.WordInfos))
	for i, wordInfo := range wordEntity.Edges.WordInfos {
		partOfSpeech := models.PartOfSpeech{
			ID:   wordInfo.Edges.PartOfSpeech.ID,
			Name: wordInfo.Edges.PartOfSpeech.Name,
		}
		japaneseMeans := make([]models.JapaneseMean, len(wordInfo.Edges.JapaneseMeans))
		for j, mean := range wordInfo.Edges.JapaneseMeans {
			japaneseMeans[j] = models.JapaneseMean{
				ID:   mean.ID,
				Name: mean.Name,
			}
		}
		wordInfos[i] = models.WordInfo{
			ID:            wordInfo.ID,
			PartOfSpeech:  partOfSpeech,
			JapaneseMeans: japaneseMeans,
		}
	}

	response := &models.WordResponse{
		ID:           wordEntity.ID,
		Name:         wordEntity.Name,
		WordInfos:    wordInfos,
		IsRegistered: isRegistered,
		TestCount:    testCount,
		CheckCount:   checkCount,
		Memo:         memo,
	}

	return response, nil
}
