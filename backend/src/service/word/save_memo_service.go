package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

func (s *WordServiceImpl) SaveMemo(ctx context.Context, SaveMemoRequest *models.SaveMemoRequest) (*models.SaveMemoResponse, error) {
	wordID := SaveMemoRequest.WordID
	userID := SaveMemoRequest.UserID
	Memo := SaveMemoRequest.Memo

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

	// userチェック
	_, err = s.client.User().Get(ctx, userID)
	if err != nil {
		logrus.Error(err)
		return nil, ErrUserNotFound
	}

	word, err := s.client.Word().
		Query().
		Where(
			word.ID(wordID),
		).
		Only(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}

	registeredWord, err := s.client.RegisteredWord().
		Query().
		Where(
			registeredword.UserID(userID),
			registeredword.WordID(wordID),
		).
		Only(ctx)

	// 登録単語が存在しない場合、新規作成
	if ent.IsNotFound(err) {
		registeredWord, err = s.client.RegisteredWord().
			Create().
			SetUserID(userID).
			SetWordID(wordID).
			SetIsActive(false).
			SetMemo(Memo).
			Save(ctx)
		if err != nil {
			return nil, errors.New("Failed to create RegisteredWord")
		}

		response := &models.SaveMemoResponse{
			Name:    word.Name,
			Memo:    *registeredWord.Memo,
			Message: "Word memo saved",
		}

		return response, nil
	}

	if err != nil {
		// その他のエラー
		return nil, errors.New("Failed to query RegisteredWord")
	}

	// 既存の登録がある場合、is_active(登録or解除)は変えず、メモのみ更新
	registeredWord, err = registeredWord.Update().
		SetMemo(Memo).
		Save(ctx)
	if err != nil {
		return nil, errors.New("Failed to update RegisteredWord")
	}

	response := &models.SaveMemoResponse{
		Name:    word.Name,
		Memo:    *registeredWord.Memo,
		Message: "RegisteredWord memo updated",
	}

	return response, nil
}
