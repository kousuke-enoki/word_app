package word_service

import (
	"context"
	"word_app/backend/src/models"
)

func (s *WordServiceImpl) RegisteredWordCount(ctx context.Context, wordID int, isRegistered bool) (*models.RegisteredWordCountResponse, error) {
	var registrationCount int
	if isRegistered {
		// Word の registration_count を +1 更新
		word, err := s.client.Word.
			UpdateOneID(wordID).
			AddRegistrationCount(1).
			Save(ctx)
		registrationCount = word.RegistrationCount
		if err != nil {
			return nil, err
		}
	} else {
		// Word の registration_count を -1 更新
		word, err := s.client.Word.
			UpdateOneID(wordID).
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
