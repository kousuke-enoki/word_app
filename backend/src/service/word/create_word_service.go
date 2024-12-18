package word_service

import (
	"context"
	"errors"
	"log"

	"word_app/backend/ent"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

// create_word
func (s *WordServiceImpl) CreateWord(ctx context.Context, CreateWordRequest *models.CreateWordRequest) (*models.CreateWordResponse, error) {
	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDeleteWord
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 管理者チェック (将来の拡張を考慮)
	userEntity, err := tx.User.Get(ctx, CreateWordRequest.UserID)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}
	if !userEntity.Admin {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrUnauthorized
	}

	// 既存の単語があるかどうか確認
	existingWord, err := s.client.Word.Query().Where(word.Name(CreateWordRequest.Name)).Only(ctx)
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
			SetName(CreateWordRequest.Name).
			SetVoiceID("").
			Save(ctx)
		if err != nil {
			return nil, errors.New(`failed to create word: , "name"`)
		}
	}

	for _, wordInfo := range CreateWordRequest.WordInfos {
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

	response := &models.CreateWordResponse{
		ID:      createdWord.ID,
		Name:    createdWord.Name,
		Message: "create word complete",
	}

	return response, nil
}
