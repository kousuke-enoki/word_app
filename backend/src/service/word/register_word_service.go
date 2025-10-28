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

func (s *ServiceImpl) RegisterWords(ctx context.Context, req *models.RegisterWordRequest) (*models.RegisterWordResponse, error) {
	return s.withTx(ctx, func(tx *ent.Tx) (*models.RegisterWordResponse, error) {
		// 1) per-user lock
		if err := s.userRepo.LockByID(ctx, req.UserID); err != nil {
			return nil, err
		}

		// 2) word must exist
		w, err := s.getWord(ctx, tx, req.WordID)
		if err != nil {
			return nil, errors.New("failed to fetch word")
		}

		// 3) current registration (may be nil)
		rw, err := s.getRegisteredWord(ctx, tx, req.UserID, req.WordID)
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
			if err := s.ensureWithinLimit(ctx, tx, req.UserID); err != nil {
				return nil, err
			}
			if rw == nil {
				if _, err := tx.RegisteredWord.Create().
					SetUserID(req.UserID).
					SetWordID(req.WordID).
					SetIsActive(true).
					Save(ctx); err != nil {
					return nil, errors.New("failed to create RegisteredWord")
				}
				return &models.RegisterWordResponse{
					Name:              w.Name,
					IsRegistered:      true,
					RegistrationCount: 1,
					Message:           "RegisteredWord created",
				}, nil
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
		}

		// ON→OFF: 上限チェック不要
		if _, err := rw.Update().SetIsActive(false).Save(ctx); err != nil {
			return nil, errors.New("failed to update RegisteredWord")
		}
		return &models.RegisterWordResponse{
			Name:              w.Name,
			IsRegistered:      false,
			RegistrationCount: 0,
			Message:           "RegisteredWord updated",
		}, nil
	})
}

// --- helpers ---

func (s *ServiceImpl) withTx(
	ctx context.Context,
	fn func(*ent.Tx) (*models.RegisterWordResponse, error),
) (_ *models.RegisterWordResponse, retErr error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, errors.New("failed to start transaction")
	}
	defer func() {
		if retErr != nil {
			_ = tx.Rollback()
			return
		}
		if cerr := tx.Commit(); cerr != nil {
			retErr = cerr
		}
	}()
	return fn(tx)
}

func (s *ServiceImpl) getWord(ctx context.Context, tx *ent.Tx, wordID int) (*ent.Word, error) {
	return tx.Word.Query().Where(word.ID(wordID)).Only(ctx)
}

func (s *ServiceImpl) getRegisteredWord(ctx context.Context, tx *ent.Tx, userID, wordID int) (*ent.RegisteredWord, error) {
	rw, err := tx.RegisteredWord.Query().Where(
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

func (s *ServiceImpl) ensureWithinLimit(ctx context.Context, tx *ent.Tx, userID int) error {
	cnt, err := tx.RegisteredWord.Query().Where(
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
