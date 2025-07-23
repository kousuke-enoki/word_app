package word

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/japanesemean"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

// ======== public =========
func (s *ServiceImpl) DeleteWord(
	ctx context.Context,
	req *models.DeleteWordRequest,
) (resp *models.DeleteWordResponse, err error) {

	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDeleteWord
	}
	defer finishTxWithLog(&err, tx)

	// --- ① 管理者チェック ---
	if err = s.assertAdmin(ctx, tx, req.UserID); err != nil {
		return nil, err
	}

	// --- ② Word を論理削除 ---
	wordName, err := s.cascadeDeleteWord(ctx, tx, req.WordID)
	if err != nil {
		return nil, err
	}

	resp = &models.DeleteWordResponse{
		Name:    wordName, // Word を取得した際にセットしておく
		Message: "word delete complete",
	}
	return
}

/*========== helper: Tx wrap ==========*/
func finishTxWithLog(perr *error, tx *ent.Tx) {
	if p := recover(); p != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			logrus.Errorf("rollback failed after panic: %v", rbErr)
		}
		panic(p)
	} else if *perr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			logrus.Errorf("rollback failed: %v", rbErr)
		}
	} else if cErr := tx.Commit(); cErr != nil {
		logrus.Errorf("commit failed: %v", cErr)
		*perr = cErr
	}
}

/*========== domain-like helpers ==========*/

// 管理者チェックだけ担当
func (s *ServiceImpl) assertAdmin(
	ctx context.Context,
	tx *ent.Tx,
	userID int,
) error {
	u, err := tx.User.Get(ctx, userID)
	switch {
	case ent.IsNotFound(err):
		return ErrUnauthorized
	case err != nil:
		return ErrDeleteWord
	case !u.IsAdmin:
		return ErrUnauthorized
	}
	return nil
}

// Word, WordInfo, JapaneseMean, RegisteredWord を削除
func (s *ServiceImpl) cascadeDeleteWord(
	ctx context.Context,
	tx *ent.Tx,
	wordID int,
) (string, error) {

	// ① Word 存在確認
	w, err := tx.Word.Get(ctx, wordID)
	if ent.IsNotFound(err) {
		return "", ErrWordNotFound
	} else if err != nil {
		return "", ErrDeleteWord
	}

	// ② 関連行をまとめて削除（順序は外部キー制約に合わせる）
	queries := []func(context.Context) error{
		func(c context.Context) error {
			_, e := tx.JapaneseMean.Delete().
				Where(japanesemean.HasWordInfoWith(wordinfo.HasWordWith(word.IDEQ(wordID)))).
				Exec(c)
			return e
		},
		func(c context.Context) error {
			_, e := tx.WordInfo.Delete().
				Where(wordinfo.HasWordWith(word.IDEQ(wordID))).
				Exec(c)
			return e
		},
		func(c context.Context) error {
			_, e := tx.RegisteredWord.Delete().
				Where(registeredword.HasWordWith(word.IDEQ(wordID))).
				Exec(c)
			return e
		},
		func(c context.Context) error {
			return tx.Word.DeleteOne(w).Exec(c)
		},
	}
	for _, fn := range queries {
		if err = fn(ctx); err != nil {
			return "", ErrDeleteWord // finishTx が Rollback してくれる
		}
	}
	return w.Name, nil
}
