package word_service

import (
	"context"
	"errors"
	"word_app/backend/ent"
	"word_app/backend/ent/japanesemean"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

var (
	ErrWordNotFound = errors.New("word not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrDeleteWord   = errors.New("failed to delete word")
)

func (s *WordServiceImpl) DeleteWord(ctx context.Context, userID int, wordID int) (*models.WordDeleteResponse, error) {
	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDeleteWord
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// word を取得して存在確認
	wordEntity, err := tx.Word.Query().Where(word.IDEQ(wordID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrWordNotFound
		}
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}

	// 管理者チェック (将来の拡張を考慮)
	userEntity, err := tx.User.Get(ctx, userID)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrDeleteWord
	}
	if !userEntity.Admin {
		logrus.Error(err)
		tx.Rollback()
		return nil, ErrUnauthorized
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

	response := &models.WordDeleteResponse{
		Name:    wordEntity.Name,
		Message: "word delete complete",
	}

	return response, nil
}
