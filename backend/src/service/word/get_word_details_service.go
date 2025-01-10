package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

// word_show
func (s *WordServiceImpl) GetWordDetails(ctx context.Context, WordShowRequest *models.WordShowRequest) (*models.WordShowResponse, error) {
	wordID := WordShowRequest.WordID
	userID := WordShowRequest.UserID

	// user存在チェック
	_, err := s.client.User.Get(ctx, userID)
	if err != nil {
		logrus.Error(err)
		return nil, ErrUserNotFound
	}

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
	var testCount, checkCount, attentionLevel int
	var memo string

	// userID に関連付けられた RegisteredWord を検索
	var registeredWord *ent.RegisteredWord
	for _, rw := range wordEntity.Edges.RegisteredWords {
		if rw.UserID == userID {
			registeredWord = rw
			break
		}
	}

	// registeredWord が存在する場合の処理
	if registeredWord != nil {
		isRegistered = registeredWord.IsActive
		attentionLevel = registeredWord.AttentionLevel
		testCount = registeredWord.TestCount
		checkCount = registeredWord.CheckCount
		if registeredWord.Memo != nil {
			memo = *registeredWord.Memo
		}
	} else {
		isRegistered = false
		attentionLevel = 0
		testCount = 0
		checkCount = 0
		memo = ""
	}

	// WordInfosを変換
	wordInfos := make([]models.WordInfo, len(wordEntity.Edges.WordInfos))
	for i, wordInfo := range wordEntity.Edges.WordInfos {
		partOfSpeech := models.PartOfSpeech{
			ID: wordInfo.Edges.PartOfSpeech.ID,
		}
		japaneseMeans := make([]models.JapaneseMean, len(wordInfo.Edges.JapaneseMeans))
		for j, mean := range wordInfo.Edges.JapaneseMeans {
			japaneseMeans[j] = models.JapaneseMean{
				ID:   mean.ID,
				Name: mean.Name,
			}
		}

		wordInfos[i] = models.WordInfo{
			ID:             wordInfo.ID,
			PartOfSpeechID: partOfSpeech.ID,
			JapaneseMeans:  japaneseMeans,
		}
	}

	response := &models.WordShowResponse{
		ID:                wordEntity.ID,
		Name:              wordEntity.Name,
		RegistrationCount: wordEntity.RegistrationCount,
		WordInfos:         wordInfos,
		IsRegistered:      isRegistered,
		AttentionLevel:    attentionLevel,
		TestCount:         testCount,
		CheckCount:        checkCount,
		Memo:              memo,
	}

	return response, nil
}
