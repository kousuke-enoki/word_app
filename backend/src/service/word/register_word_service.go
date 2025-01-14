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

func (s *WordServiceImpl) RegisterWords(ctx context.Context, req *models.RegisterWordRequest) (*models.RegisterWordResponse, error) {

	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error("Failed to start transaction: ", err)
		return nil, errors.New("failed to start transaction")
	}

	// トランザクション終了時のロールバック処理（deferを使う）
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// トランザクション付きコンテキスト作成
	ctx = ent.NewTxContext(ctx, tx)

	// ユーザー確認
	if err := s.checkUserExists(ctx, req.UserID); err != nil {
		logrus.Error(err)
		return nil, ErrUserNotFound
	}

	// 単語取得
	wordEntity, err := s.getWord(ctx, req.WordID)
	if err != nil {
		return nil, err
	}

	// 登録状態の取得
	registeredWord, err := s.getRegisteredWord(ctx, req.UserID, req.WordID)

	// 登録状態の処理
	if registeredWord == nil && req.IsRegistered {
		return s.createRegisteredWord(ctx, req, wordEntity.Name)
	}

	if registeredWord == nil && !req.IsRegistered {
		return nil, errors.New("Failed to unregister: word is not registered")
	}

	if registeredWord.IsActive == req.IsRegistered {
		return nil, errors.New("No change in registration state")
	}

	// 更新処理
	return s.updateRegisteredWord(ctx, registeredWord, req.IsRegistered, wordEntity.Name)
}

// ユーザーが存在するか確認
func (s *WordServiceImpl) checkUserExists(ctx context.Context, userID int) error {
	_, err := s.client.User().Get(ctx, userID)
	return err
}

// 単語を取得
func (s *WordServiceImpl) getWord(ctx context.Context, wordID int) (*ent.Word, error) {
	wordEntity, err := s.client.Word().
		Query().
		Where(word.ID(wordID)).
		Only(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}
	return wordEntity, nil
}

// 登録状態の取得
func (s *WordServiceImpl) getRegisteredWord(ctx context.Context, userID, wordID int) (*ent.RegisteredWord, error) {
	registeredWord, err := s.client.RegisteredWord().
		Query().
		Where(
			registeredword.UserID(userID),
			registeredword.WordID(wordID),
		).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, nil // 登録が存在しない場合はnilを返す
	}
	if err != nil {
		return nil, errors.New("failed to query RegisteredWord")
	}

	return registeredWord, nil
}

// 登録単語の新規作成
func (s *WordServiceImpl) createRegisteredWord(ctx context.Context, req *models.RegisterWordRequest, wordName string) (*models.RegisterWordResponse, error) {
	_, err := s.client.RegisteredWord().
		Create().
		SetUserID(req.UserID).
		SetWordID(req.WordID).
		SetIsActive(true).
		Save(ctx)

	if err != nil {
		return nil, errors.New("Failed to create RegisteredWord")
	}

	return s.generateResponse(ctx, req.WordID, true, wordName, "RegisteredWord created")
}

// 登録単語の更新
func (s *WordServiceImpl) updateRegisteredWord(ctx context.Context, registeredWord *ent.RegisteredWord, isActive bool, wordName string) (*models.RegisterWordResponse, error) {
	_, err := registeredWord.Update().
		SetIsActive(isActive).
		Save(ctx)
	if err != nil {
		return nil, errors.New("Failed to update RegisteredWord")
	}

	return s.generateResponse(ctx, registeredWord.WordID, isActive, wordName, "RegisteredWord updated")
}

// レスポンスの生成
func (s *WordServiceImpl) generateResponse(ctx context.Context, wordID int, isRegistered bool, wordName, message string) (*models.RegisterWordResponse, error) {
	registeredWordCountRequest := &models.RegisteredWordCountRequest{
		WordID:       wordID,
		IsRegistered: isRegistered,
	}
	registrationCountResponse, err := s.RegisteredWordCount(ctx, registeredWordCountRequest)
	if err != nil {
		return nil, err
	}

	return &models.RegisterWordResponse{
		Name:              wordName,
		IsRegistered:      isRegistered,
		RegistrationCount: registrationCountResponse.RegistrationCount,
		Message:           message,
	}, nil
}
