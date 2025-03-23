package word_service

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

func (s *WordServiceImpl) DeleteWord(ctx context.Context, DeleteWordRequest *models.DeleteWordRequest) (*models.DeleteWordResponse, error) {
	userID := DeleteWordRequest.UserID
	wordID := DeleteWordRequest.WordID
	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDeleteWord
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

	// 管理者チェック
	userEntity, err := tx.User.Get(ctx, userID)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}
	if !userEntity.IsAdmin {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrUnauthorized
	}

	// word を取得して存在確認
	wordEntity, err := tx.Word.Query().Where(word.IDEQ(wordID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			logrus.Info("word not found")
			return nil, ErrWordNotFound
		}
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}

	// wordInfo に紐づく japaneseMean を削除
	_, err = tx.JapaneseMean.Delete().Where(
		japanesemean.HasWordInfoWith(wordinfo.HasWordWith(word.IDEQ(wordID))),
	).Exec(ctx)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}

	// word に紐づく wordInfo を削除
	_, err = tx.WordInfo.Delete().Where(wordinfo.HasWordWith(word.IDEQ(wordID))).Exec(ctx)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}

	// word に紐づく registeredword を削除
	_, err = tx.RegisteredWord.Delete().Where(registeredword.HasWordWith(word.IDEQ(wordID))).Exec(ctx)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}

	// 最後に word を削除
	err = tx.Word.DeleteOne(wordEntity).Exec(ctx)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return nil, ErrDeleteWord
	}

	response := &models.DeleteWordResponse{
		Name:    wordEntity.Name,
		Message: "word delete complete",
	}

	return response, nil
}
