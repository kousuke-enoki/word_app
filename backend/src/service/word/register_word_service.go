package word

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/shared/ucerr"
)

// ServiceImpl は起動時に DI で Limits を注入しておく
// type ServiceImpl struct {
// 	client *ent.Client
// 	limits struct {
// 		RegisteredWordsPerUser int
// 	}
// }

// RegisterWords: ユーザーの「この単語の登録ON/OFF」を切り替える。
// 上限チェックは「ONにする時だけ」、かつ per-user lock でレース対策。
func (s *ServiceImpl) RegisterWords(ctx context.Context, req *models.RegisterWordRequest) (_ *models.RegisterWordResponse, retErr error) {
	// 0) トランザクション開始
	txCtx, done, err := s.txm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = done(false) }()

	// --- 1) ユーザーの存在 + ロック（同一ユーザー操作を直列化）---
	// ent v0.14: Modify で ForUpdate()
	err = s.userRepo.LockByID(txCtx, req.UserID)
	if err != nil {
		return nil, err
	}

	// --- 2) 単語の存在 ---
	w, err := s.getWord(txCtx, req.WordID)
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}

	// --- 3) 現在の登録状態を取得（1行ユニーク前提）---
	rw, err := s.getRegisteredWord(txCtx, req.UserID, req.WordID)
	if err != nil {
		return nil, err
	}

	// 4) no-op / invalid transitions
	if rw == nil && !req.IsRegistered {
		return nil, ucerr.BadRequest("word is not registered")
	}
	if rw != nil && rw.IsActive == req.IsRegistered {
		return nil, ucerr.BadRequest("no change in registration state")
	}

	// 5) desired transitions
	if req.IsRegistered {
		// OFF→ON or 新規ON: 上限チェック→作成/更新
		if err := s.ensureWithinLimit(txCtx, req.UserID); err != nil {
			return nil, err
		}
		if rw == nil {
			if _, err := s.client.RegisteredWord().Create().
				SetUserID(req.UserID).
				SetWordID(req.WordID).
				SetIsActive(true).
				Save(txCtx); err != nil {
				return nil, errors.New("failed to create RegisteredWord")
			}
			return &models.RegisterWordResponse{
				Name:              w.Name,
				IsRegistered:      true,
				RegistrationCount: 1,
				Message:           "RegisteredWord created",
			}, nil
		}
		if _, err := rw.Update().SetIsActive(true).Save(txCtx); err != nil {
			return nil, errors.New("failed to update RegisteredWord")
		}
		return &models.RegisterWordResponse{
			Name:              w.Name,
			IsRegistered:      true,
			RegistrationCount: 1,
			Message:           "RegisteredWord updated",
		}, nil
	}

	// ON→OFF: 上限チェック不要
	if _, err := rw.Update().SetIsActive(false).Save(txCtx); err != nil {
		return nil, errors.New("failed to update RegisteredWord")
	}
	return &models.RegisterWordResponse{
		Name:              w.Name,
		IsRegistered:      false,
		RegistrationCount: 0,
		Message:           "RegisteredWord updated",
	}, nil
}

func (s *ServiceImpl) getWord(ctx context.Context, wordID int) (*ent.Word, error) {
	return s.client.Word().Query().Where(word.ID(wordID)).Only(ctx)
}

func (s *ServiceImpl) getRegisteredWord(ctx context.Context, userID, wordID int) (*ent.RegisteredWord, error) {
	rw, err := s.client.RegisteredWord().Query().Where(
		registeredword.UserID(userID),
		registeredword.WordID(wordID),
	).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("failed to query RegisteredWord")
	}
	return rw, nil
}

func (s *ServiceImpl) ensureWithinLimit(ctx context.Context, userID int) error {
	// テストユーザーのみ制限を適用
	// context から Principal を取得して isTest をチェック
	p, ok := jwt.GetPrincipalFromContext(ctx)
	if !ok || !p.IsTest {
		// Principal が取得できない、またはテストユーザーでない場合は制限なし
		return nil
	}

	cnt, err := s.client.RegisteredWord().Query().Where(
		registeredword.UserID(userID),
		registeredword.IsActive(true),
	).Count(ctx)
	if err != nil {
		return err
	}
	if cnt >= s.limits.RegisteredWordsPerUser {
		return ucerr.TooManyRequests("registered words limit exceeded")
	}
	return nil
}
