package word

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
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
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, errors.New("failed to start transaction")
	}
	defer func() {
		// retErr が nil なら commit、そうでなければ rollback
		if retErr != nil {
			_ = tx.Rollback()
		} else {
			if cerr := tx.Commit(); cerr != nil {
				retErr = cerr
			}
		}
	}()

	// --- 1) ユーザーの存在 + ロック（同一ユーザー操作を直列化）---
	// ent v0.14: Modify で ForUpdate()
	err = s.userRepo.LockByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// --- 2) 単語の存在 ---
	w, err := tx.Word.
		Query().
		Where(word.ID(req.WordID)).
		Only(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}

	// --- 3) 現在の登録状態を取得（1行ユニーク前提）---
	rw, err := tx.RegisteredWord.
		Query().
		Where(
			registeredword.UserID(req.UserID),
			registeredword.WordID(req.WordID),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		rw = nil
	} else if err != nil {
		return nil, errors.New("failed to query RegisteredWord")
	}

	// --- 4) 分岐ロジック ---
	switch {
	// 4-1) まだ行が無く、登録ONにしたい → 上限チェック→作成
	case rw == nil && req.IsRegistered:
		// active だけ数える
		activeCnt, err := tx.RegisteredWord.
			Query().
			Where(
				registeredword.UserID(req.UserID),
				registeredword.IsActive(true),
			).
			Count(ctx)
		if err != nil {
			return nil, err
		}
		if activeCnt >= s.limits.RegisteredWordsPerUser {
			return nil, ucerr.TooManyRequests("registered words limit exceeded")
		}
		if _, err := tx.RegisteredWord.
			Create().
			SetUserID(req.UserID).
			SetWordID(req.WordID).
			SetIsActive(true).
			Save(ctx); err != nil {
			return nil, errors.New("failed to create RegisteredWord")
		}
		return &models.RegisterWordResponse{
			Name:              w.Name,
			IsRegistered:      true,
			RegistrationCount: 1, // 必要なら別途集計、ここではダミー
			Message:           "RegisteredWord created",
		}, nil

	// 4-2) 行が無く、登録OFF → 何もしない（エラーにしたいなら 400）
	case rw == nil && !req.IsRegistered:
		return nil, ucerr.BadRequest("word is not registered")

	// 4-3) 行があり、状態変化なし → 409 or 200（ここでは 400）
	case rw.IsActive == req.IsRegistered:
		return nil, ucerr.BadRequest("no change in registration state")

	// 4-4) 行があり、OFF→ON へ → 上限チェック→更新
	case !rw.IsActive && req.IsRegistered:
		activeCnt, err := tx.RegisteredWord.
			Query().
			Where(
				registeredword.UserID(req.UserID),
				registeredword.IsActive(true),
			).
			Count(ctx)
		if err != nil {
			return nil, err
		}
		if activeCnt >= s.limits.RegisteredWordsPerUser {
			return nil, ucerr.TooManyRequests("registered words limit exceeded")
		}
		if _, err := rw.Update().SetIsActive(true).Save(ctx); err != nil {
			return nil, errors.New("failed to update RegisteredWord")
		}
		return &models.RegisterWordResponse{
			Name:              w.Name,
			IsRegistered:      true,
			RegistrationCount: 1,
			Message:           "RegisteredWord updated",
		}, nil

	// 4-5) 行があり、ON→OFF へ → 単にOFF（上限チェック不要）
	default: // rw.IsActive && !req.IsRegistered
		if _, err := rw.Update().SetIsActive(false).Save(ctx); err != nil {
			return nil, errors.New("failed to update RegisteredWord")
		}
		return &models.RegisterWordResponse{
			Name:              w.Name,
			IsRegistered:      false,
			RegistrationCount: 0,
			Message:           "RegisteredWord updated",
		}, nil
	}
}
