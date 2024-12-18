package word_service

import (
	"context"
	"errors"
	"fmt"

	"word_app/backend/ent"
	"word_app/backend/ent/japanesemean"
	"word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

func (s *WordServiceImpl) UpdateWord(ctx context.Context, req *models.UpdateWordRequest) (*models.UpdateWordResponse, error) {
	// // リクエストバリデーション
	// if req.ID == 0 || req.Name == "" || len(req.WordInfos) == 0 {
	// 	return nil, errors.New("invalid request: ID, Name, or WordInfos is missing")
	// }
	logrus.Info("asdf")
	logrus.Info("req", req)
	logrus.Info("qwer")
	// トランザクション開始
	tx, err := s.client.Tx(ctx)
	if err != nil {
		logrus.Error("failed to start transaction:", err)
		return nil, errors.New("failed to start transaction")
	}

	// トランザクションロールバック処理
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 管理者チェック
	userEntity, err := tx.User.Get(ctx, req.UserID)
	if err != nil {
		logrus.Error("failed to get user:", err)
		return nil, errors.New("failed to get user")
	}
	if !userEntity.Admin {
		logrus.Error("unauthorized user")
		return nil, errors.New("unauthorized user")
	}

	// 既存の単語を取得
	existingWord, err := tx.Word.Get(ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("word not found")
		}
		logrus.Error("failed to fetch word:", err)
		return nil, errors.New("failed to fetch word")
	}

	// 単語の名前を変える場合、同じ名前の単語があるかどうか確認。ある場合は失敗
	if req.Name != existingWord.Name {
		_, err := s.client.Word.Query().Where(word.Name(req.Name)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			logrus.Fatalf("failed to query word: %v", err)
		}
	}

	// 単語を更新
	updatedWord, err := tx.Word.UpdateOne(existingWord).
		SetName(req.Name).
		Save(ctx)
	if err != nil {
		logrus.Error("failed to update word:", err)
		return nil, errors.New("failed to update word")
	}

	// 関連する WordInfo を更新
	for _, wordInfo := range req.WordInfos {
		var wordInfoEntity *ent.WordInfo

		// 既存の WordInfo があるか確認
		wordInfoEntity, err = tx.WordInfo.Query().
			Where(wordinfo.IDEQ(wordInfo.ID)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				// 存在しない場合、新規作成
				wordInfoEntity, err = tx.WordInfo.Create().
					SetWordID(updatedWord.ID).
					SetPartOfSpeechID(wordInfo.PartOfSpeechID).
					Save(ctx)
				if err != nil {
					logrus.Error("failed to create word info:", err)
					return nil, errors.New("failed to create word info")
				}
			} else {
				logrus.Error("failed to fetch word info:", err)
				return nil, errors.New("failed to fetch word info")
			}
		} else {
			// 存在する場合、更新
			wordInfoEntity, err = tx.WordInfo.UpdateOne(wordInfoEntity).
				SetPartOfSpeechID(wordInfo.PartOfSpeechID).
				Save(ctx)
			if err != nil {
				logrus.Error("failed to update word info:", err)
				return nil, errors.New("failed to update word info")
			}
		}

		// JapaneseMean の更新
		for _, japaneseMean := range wordInfo.JapaneseMeans {
			var meanEntity *ent.JapaneseMean

			// 既存の JapaneseMean を検索
			meanEntity, err = tx.JapaneseMean.Query().
				Where(japanesemean.IDEQ(japaneseMean.ID)).
				Only(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					// 存在しない場合、新規作成
					_, err = tx.JapaneseMean.Create().
						SetWordInfoID(wordInfoEntity.ID).
						SetName(japaneseMean.Name).
						Save(ctx)
					if err != nil {
						logrus.Error("failed to create japanese mean:", err)
						return nil, errors.New("failed to create japanese mean")
					}
				} else {
					logrus.Error("failed to fetch japanese mean:", err)
					return nil, errors.New("failed to fetch japanese mean")
				}
			} else {
				// 存在する場合、更新
				_, err = tx.JapaneseMean.UpdateOne(meanEntity).
					SetName(japaneseMean.Name).
					Save(ctx)
				if err != nil {
					logrus.Error("failed to update japanese mean:", err)
					return nil, errors.New("failed to update japanese mean")
				}
			}
		}
	}

	// トランザクションコミット
	err = tx.Commit()
	if err != nil {
		logrus.Error("failed to commit transaction:", err)
		return nil, errors.New("failed to commit transaction")
	}

	// レスポンス生成
	response := &models.UpdateWordResponse{
		ID:      updatedWord.ID,
		Name:    updatedWord.Name,
		Message: fmt.Sprintf("word '%s' updated successfully", updatedWord.Name),
	}

	return response, nil
}
