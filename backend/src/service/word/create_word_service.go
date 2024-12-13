package word_service

import (
	"context"
	"errors"
	"log"

	"word_app/backend/ent"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"
)

// word_show
func (s *WordServiceImpl) CreateWord(ctx context.Context, WordCreateRequest *models.WordCreateRequest) (*models.WordCreateResponse, error) {
	existingWord, err := s.client.Word.Query().Where(word.Name(WordCreateRequest.Name)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		log.Fatalf("failed to query word: %v", err)
	}

	var createdWord *ent.Word
	if existingWord != nil {
		// 既存の単語がある場合は、エラー
		return nil, errors.New("There is already a word with the same name.")
	} else {
		// ない場合は新しい単語を作成
		createdWord, err = s.client.Word.Create().
			SetName(WordCreateRequest.Name).
			SetVoiceID("").
			Save(ctx)
		if err != nil {
			return nil, errors.New(`failed to create word: , "name"`)
		}
	}

	for _, wordInfo := range WordCreateRequest.WordInfos {
		var partOfSpeechId int
		partOfSpeechId = wordInfo.PartOfSpeechID
		createdWordInfo, err := s.client.WordInfo.Create().
			SetWordID(createdWord.ID).
			SetPartOfSpeechID(partOfSpeechId).
			Save(ctx)
		if err != nil {
			return nil, errors.New(`failed to create word info: , "err"`)
		}

		for _, JapaneseMean := range wordInfo.JapaneseMeans {
			var japaneseMeanName string
			japaneseMeanName = JapaneseMean.Name
			// japanese_mean テーブルに日本語の意味を追加
			_, err = s.client.JapaneseMean.Create().
				SetWordInfoID(createdWordInfo.ID).
				SetName(japaneseMeanName).
				Save(ctx)
			if err != nil {
				return nil, errors.New(`failed to create japanese mean: , "err"`)
			}
		}
	}

	response := &models.WordCreateResponse{
		ID:   createdWord.ID,
		Name: createdWord.Name,
	}

	return response, nil
}
