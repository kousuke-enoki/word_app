package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"
)

func (s *WordServiceImpl) RegisterWords(ctx context.Context, RegisterWordRequest *models.RegisterWordRequest) (*models.RegisterWordResponse, error) {
	wordID := RegisterWordRequest.WordID
	userID := RegisterWordRequest.UserID
	IsRegistered := RegisterWordRequest.IsRegistered
	// wordが存在するかどうか確認
	word, err := s.client.Word.
		Query().
		Where(
			word.ID(wordID),
		).
		Only(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}

	// すでに登録されているかどうか（registeredWordがあるか）確認
	registeredWord, err := s.client.RegisteredWord.
		Query().
		Where(
			registeredword.UserID(userID),
			registeredword.WordID(wordID),
		).
		Only(ctx)

	// 登録した単語が存在しない場合、新規作成
	if ent.IsNotFound(err) && IsRegistered {
		registeredWord, err = s.client.RegisteredWord.
			Create().
			SetUserID(userID).
			SetWordID(wordID).
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			return nil, errors.New("Failed to create RegisteredWord")
		}

		registrationCountResponse, err := s.RegisteredWordCount(ctx, wordID, IsRegistered)
		if err != nil {
			return nil, err
		}

		response := &models.RegisterWordResponse{
			Name:              word.Name,
			IsRegistered:      registeredWord.IsActive,
			RegistrationCount: registrationCountResponse.RegistrationCount,
			Message:           "RegisteredWord updated",
		}

		return response, nil
	}

	if err != nil {
		// その他のエラー
		return nil, errors.New("Failed to query RegisteredWord")
	}

	// 既存の登録がある場合、is_activeをIsRegistered(登録or解除)に更新
	registeredWord, err = registeredWord.Update().
		SetIsActive(IsRegistered).
		Save(ctx)
	if err != nil {
		return nil, errors.New("Failed to update RegisteredWord")
	}

	registrationCountResponse, err := s.RegisteredWordCount(ctx, wordID, IsRegistered)
	if err != nil {
		return nil, err
	}

	response := &models.RegisterWordResponse{
		Name:              word.Name,
		IsRegistered:      registeredWord.IsActive,
		RegistrationCount: registrationCountResponse.RegistrationCount,
		Message:           "RegisteredWord updated",
	}

	return response, nil
}
