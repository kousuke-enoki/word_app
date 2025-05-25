package word_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"

	"github.com/sirupsen/logrus"
)

const maxBulkRegister = 200 // 1 リクエストで許可する単語数

func (s *WordServiceImpl) BulkRegister(ctx context.Context, userID int, words []string) (*models.BulkRegisterResponse, error) {
	if len(words) == 0 {
		return nil, errors.New("empty payload")
	}
	if len(words) > maxBulkRegister {
		return nil, fmt.Errorf("too many words: %d > %d", len(words), maxBulkRegister)
	}

	/* ---------- 正規化＋重複排除 ---------- */
	set := map[string]struct{}{}
	norm := make([]string, 0, len(words))
	for _, w := range words {
		if len(norm) >= maxBulkRegister { // max超えていたらここで打ち切り
			break
		}
		lw := strings.ToLower(w)
		if _, ok := set[lw]; !ok {
			set[lw] = struct{}{}
			norm = append(norm, lw)
		}
	}

	/* ---------- master word 一括取得 ---------- */
	found, err := s.client.Word().
		Query().
		Where(word.NameIn(norm...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	nameToID := make(map[string]int)
	for _, w := range found {
		nameToID[w.Name] = w.ID
	}

	/* ---------- 既存 registered_word 取得 ---------- */
	var ids []int
	for _, id := range nameToID {
		ids = append(ids, id)
	}
	regs, err := s.client.RegisteredWord().
		Query().
		Where(
			registeredword.UserIDEQ(userID),
			registeredword.WordIDIn(ids...),
		).All(ctx)
	if err != nil {
		return nil, err
	}
	active := map[int]bool{}
	for _, r := range regs {
		active[r.WordID] = r.IsActive
	}

	/* ---------- トランザクション ---------- */
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				logrus.Error(err)
			}
		}
	}()

	var okWords []string
	var failed []models.FailedWord

	for _, w := range norm {
		wordID, ok := nameToID[w]
		if !ok {
			failed = append(failed, models.FailedWord{Word: w, Reason: "not_exists"})
			continue
		}
		if active[wordID] { // すでにアクティブ
			failed = append(failed, models.FailedWord{Word: w, Reason: "already_registered"})
			continue
		}

		// 登録または更新
		if _, exists := active[wordID]; exists {
			// inactive -> active
			_, err = tx.RegisteredWord.
				Update().
				Where(
					registeredword.UserIDEQ(userID),
					registeredword.WordIDEQ(wordID),
				).
				SetIsActive(true).
				Save(ctx)
		} else {
			// 新規作成
			_, err = tx.RegisteredWord.
				Create().
				SetUserID(userID).
				SetWordID(wordID).
				SetIsActive(true).
				Save(ctx)
		}
		if err != nil {
			failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
			continue
		}
		okWords = append(okWords, w)
	}

	/* 4‑B. registration_count を +1 */
	_, err = s.RegisteredWordsCount(ctx,
		true,
		okWords,
	)
	if err != nil {
		return nil, err
	}

	response := &models.BulkRegisterResponse{
		Success: okWords,
		Failed:  failed,
	}

	return response, nil
}
