package word

import (
	"context"
	"errors"
	"fmt"

	"word_app/backend/ent"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

/*==================== public ====================*/

func (s *ServiceImpl) UpdateWord(
	ctx context.Context,
	req *models.UpdateWordRequest,
) (resp *models.UpdateWordResponse, err error) {

	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error("start tx:", err)
		return nil, errors.New("failed to start transaction")
	}
	defer finishTxWithLog(&err, tx)

	// ① 管理者チェック
	if err = s.assertAdmin(ctx, tx, req.UserID); err != nil {
		return nil, err
	}

	// ② 単語取得＆重複名チェック
	wordEnt, err := s.fetchAndValidateName(ctx, tx, req.ID, req.Name)
	if err != nil {
		return nil, err
	}

	// ③ 単語名更新
	wordEnt, err = s.updateWordName(ctx, tx, wordEnt, req.Name)
	if err != nil {
		return nil, err
	}

	// ④ WordInfo / JapaneseMean を upsert
	if err = s.upsertWordInfos(ctx, tx, req.WordInfos); err != nil {
		return nil, err
	}

	resp = &models.UpdateWordResponse{
		ID:      wordEnt.ID,
		Name:    wordEnt.Name,
		Message: fmt.Sprintf("word '%s' updated successfully", wordEnt.Name),
	}
	return
}

/*==================== tx wrapper ====================*/

// func finishTxWithLog(perr *error, tx *ent.Tx) {
// 	if p := recover(); p != nil {
// 		_ = tx.Rollback()
// 		panic(p)
// 	} else if *perr != nil {
// 		if rb := tx.Rollback(); rb != nil {
// 			logrus.Errorf("rollback failed: %v", rb)
// 		}
// 	} else if cErr := tx.Commit(); cErr != nil {
// 		logrus.Errorf("commit failed: %v", cErr)
// 		*perr = cErr
// 	}
// }

/*==================== domain / repository-like helpers ====================*/

// // 管理者かチェック
// func (s *ServiceImpl) assertAdmin(
// 	ctx context.Context,
// 	tx *ent.Tx,
// 	userID int,
// ) error {
// 	u, err := tx.User.Get(ctx, userID)
// 	if err != nil {
// 		return ErrUnauthorized
// 	}
// 	if !u.IsAdmin {
// 		return ErrUnauthorized
// 	}
// 	return nil
// }

// 既存 Word を取得し、名前重複を検証
func (s *ServiceImpl) fetchAndValidateName(
	ctx context.Context,
	tx *ent.Tx,
	wordID int,
	newName string,
) (*ent.Word, error) {

	w, err := tx.Word.Get(ctx, wordID)
	if ent.IsNotFound(err) {
		return nil, errors.New("word not found")
	}
	if err != nil {
		return nil, errors.New("failed to fetch word")
	}
	if w.Name != newName {
		exists, err := s.client.Word().
			Query().
			Where(word.Name(newName)).
			Exist(ctx)
		if err != nil {
			return nil, ErrDatabaseFailure
		}
		if exists {
			return nil, ErrWordExists
		}
	}
	return w, nil
}

// 単語名を更新
func (s *ServiceImpl) updateWordName(
	ctx context.Context,
	tx *ent.Tx,
	w *ent.Word,
	name string,
) (*ent.Word, error) {

	return tx.Word.UpdateOne(w).
		SetName(name).
		Save(ctx)
}

// WordInfo / JapaneseMean を upsert
func (s *ServiceImpl) upsertWordInfos(
	ctx context.Context,
	tx *ent.Tx,
	infos []models.WordInfo,
) error {

	for _, wi := range infos {
		wiEnt, err := tx.WordInfo.Get(ctx, wi.ID)
		if err != nil {
			return errors.New("failed to fetch word info")
		}
		_, err = tx.WordInfo.UpdateOne(wiEnt).
			SetPartOfSpeechID(wi.PartOfSpeechID).
			Save(ctx)
		if err != nil {
			return errors.New("failed to update word info")
		}
		if err = s.upsertMeans(ctx, tx, wi.JapaneseMeans); err != nil {
			return err
		}
	}
	return nil
}

// JapaneseMean upsert
func (s *ServiceImpl) upsertMeans(
	ctx context.Context,
	tx *ent.Tx,
	means []models.JapaneseMean,
) error {

	for _, jm := range means {
		mEnt, err := tx.JapaneseMean.Get(ctx, jm.ID)
		if err != nil {
			return errors.New("failed to fetch japanese mean")
		}
		if _, err = tx.JapaneseMean.UpdateOne(mEnt).
			SetName(jm.Name).
			Save(ctx); err != nil {
			return errors.New("failed to update japanese mean")
		}
	}
	return nil
}
