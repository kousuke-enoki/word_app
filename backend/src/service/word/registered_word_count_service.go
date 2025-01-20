package word_service

import (
	"context"
	"errors"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

func (s *WordServiceImpl) RegisteredWordCount(ctx context.Context, req *models.RegisteredWordCountRequest) (*models.RegisteredWordCountResponse, error) {

	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error("Failed to start transaction: ", err)
		return nil, errors.New("failed to start transaction")
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

	_, err = s.client.Word().
		Query().
		Where(word.ID(req.WordID)).
		Only(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}

	var registrationCount int

	if req.IsRegistered {
		// Word の registration_count を +1 更新
		word, err := s.client.Word().
			UpdateOneID(req.WordID).
			AddRegistrationCount(1).
			Save(ctx)
		registrationCount = word.RegistrationCount
		if err != nil {
			return nil, err
		}
	} else {
		// Word の registration_count を -1 更新
		word, err := s.client.Word().
			UpdateOneID(req.WordID).
			AddRegistrationCount(-1).
			Save(ctx)
		registrationCount = word.RegistrationCount
		if err != nil {
			return nil, err
		}
	}

	response := &models.RegisteredWordCountResponse{
		RegistrationCount: registrationCount,
	}
	return response, nil
}
