package word_service

import (
	"context"
	"errors"

	"word_app/backend/ent/word"

	"github.com/sirupsen/logrus"
)

func (s *WordServiceImpl) RegisteredWordsCount(ctx context.Context, IsRegistered bool, words []string) ([]string, error) {

	if len(words) <= 0 {
		logrus.Error("empty words")
		return nil, errors.New("empty words")
	}
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

	if IsRegistered {
		_, err = s.client.Word().
			Update().
			Where(word.NameIn(words...)).
			AddRegistrationCount(1).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := s.client.Word().
			Update().
			Where(word.NameIn(words...)).
			AddRegistrationCount(-1).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return words, nil
}
