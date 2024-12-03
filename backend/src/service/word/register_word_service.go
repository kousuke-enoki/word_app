package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"
)

func (s *WordServiceImpl) RegisterWords(ctx context.Context, wordID int, userID int, IsRegistered bool, memo string) (*models.RegisterWordResponse, error) {
	word, err := s.client.Word.
		Query().
		Where(
			word.ID(wordID),
		).
		Only(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}

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
			SetMemo(memo).
			Save(ctx)
		if err != nil {
			return nil, errors.New("Failed to create RegisteredWord")
		}

		response := &models.RegisterWordResponse{
			Name:         word.Name,
			IsRegistered: registeredWord.IsActive,
			Memo:         memo,
			Message:      "RegisteredWord updated",
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

	response := &models.RegisterWordResponse{
		Name:         word.Name,
		IsRegistered: registeredWord.IsActive,
		Memo:         memo,
		Message:      "RegisteredWord updated",
	}

	return response, nil
}
