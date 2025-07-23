package word

import (
	"context"

	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

// create_word
func (s *ServiceImpl) CreateWord(ctx context.Context, CreateWordRequest *models.CreateWordRequest) (*models.CreateWordResponse, error) {
	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDatabaseFailure
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				logrus.Error(err)
			}
		}
	}()

	// 管理者チェック
	userEntity, err := tx.User.Get(ctx, CreateWordRequest.UserID)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDatabaseFailure
	}
	if !userEntity.IsAdmin {
		return nil, ErrUnauthorized
	}

	// 単語の存在確認
	exists, err := tx.Word.Query().Where(word.Name(CreateWordRequest.Name)).Exist(ctx)
	if err != nil {
		logrus.Errorf("failed to query word existence: %v", err)
		return nil, ErrDatabaseFailure
	}
	if exists {
		return nil, ErrWordExists
	}

	// 新しい単語を作成
	createdWord, err := tx.Word.Create().
		SetName(CreateWordRequest.Name).
		Save(ctx)
	if err != nil {
		return nil, ErrCreateWord
	}

	// WordInfoとJapaneseMeanを作成
	for _, wordInfo := range CreateWordRequest.WordInfos {
		createdWordInfo, err := tx.WordInfo.Create().
			SetWordID(createdWord.ID).
			SetPartOfSpeechID(wordInfo.PartOfSpeechID).
			Save(ctx)
		if err != nil {
			return nil, ErrCreateWordInfo
		}

		for _, JapaneseMean := range wordInfo.JapaneseMeans {
			_, err = tx.JapaneseMean.Create().
				SetWordInfoID(createdWordInfo.ID).
				SetName(JapaneseMean.Name).
				Save(ctx)
			if err != nil {
				return nil, ErrCreateJapaneseMean
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
