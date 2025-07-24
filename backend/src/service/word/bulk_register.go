package word

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/src/models"
)

const maxBulkRegister = 200

/*==================== public ====================*/

// orchestration だけ残す
func (s *ServiceImpl) BulkRegister(
	ctx context.Context,
	userID int,
	words []string,
) (resp *models.BulkRegisterResponse, err error) {

	if err = validatePayload(words); err != nil {
		return nil, err
	}

	// ① 正規化＋重複排除
	norm := normalizeWords(words)

	// ② master word → map[name]id
	nameToID, err := s.fetchWordIDs(ctx, norm)
	if err != nil {
		return nil, err
	}

	// ③ 既存 registered_word の活性状態を取得
	activeMap, err := s.fetchActiveRegs(ctx, userID, nameToID)
	if err != nil {
		return nil, err
	}

	// ④ Tx で upsert
	okWords, failed, err := s.upsertRegsTx(ctx, userID, norm, nameToID, activeMap)
	if err != nil {
		return nil, err
	}

	// ⑤ 登録回数を +1
	if err = s.incrementRegCount(ctx, okWords); err != nil {
		return nil, err
	}

	return &models.BulkRegisterResponse{Success: okWords, Failed: failed}, nil
}

/*==================== validation & util ====================*/

func validatePayload(words []string) error {
	switch l := len(words); {
	case l == 0:
		return errors.New("empty payload")
	case l > maxBulkRegister:
		return fmt.Errorf("too many words: %d > %d", l, maxBulkRegister)
	}
	return nil
}

func normalizeWords(words []string) []string {
	set := make(map[string]struct{}, len(words))
	out := make([]string, 0, len(words))
	for _, w := range words {
		if len(out) >= maxBulkRegister {
			break
		}
		lw := strings.ToLower(w)
		if _, ok := set[lw]; !ok {
			set[lw] = struct{}{}
			out = append(out, lw)
		}
	}
	return out
}

func (s *ServiceImpl) fetchWordIDs(
	ctx context.Context,
	names []string,
) (map[string]int, error) {

	found, err := s.client.Word().
		Query().
		Where(word.NameIn(names...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	res := make(map[string]int, len(found))
	for _, w := range found {
		res[w.Name] = w.ID
	}
	return res, nil
}

func (s *ServiceImpl) fetchActiveRegs(
	ctx context.Context,
	userID int,
	nameToID map[string]int,
) (map[int]bool, error) {

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
	active := make(map[int]bool, len(regs))
	for _, r := range regs {
		active[r.WordID] = r.IsActive
	}
	return active, nil
}

/*==================== Tx + upsert ====================*/

func (s *ServiceImpl) upsertRegsTx(
	ctx context.Context,
	userID int,
	norm []string,
	nameToID map[string]int,
	active map[int]bool,
) (okWords []string, failed []models.FailedWord, err error) {

	tx, err := s.client.Tx(ctx)
	if err != nil {
		return
	}
	defer finishTxWithLog(&err, tx)

	for _, w := range norm {
		wordID, exists := nameToID[w]
		if !exists {
			failed = append(failed, models.FailedWord{Word: w, Reason: "not_exists"})
			continue
		}
		if active[wordID] {
			failed = append(failed, models.FailedWord{Word: w, Reason: "already_registered"})
			continue
		}

		if _, dup := active[wordID]; dup {
			err = s.activateReg(ctx, tx, userID, wordID)
		} else {
			err = s.createReg(ctx, tx, userID, wordID)
		}
		if err != nil {
			failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
			err = nil // 1 つ失敗しても続行
			continue
		}
		okWords = append(okWords, w)
	}
	return
}

func (s *ServiceImpl) activateReg(
	ctx context.Context,
	tx *ent.Tx,
	userID, wordID int,
) error {
	_, err := tx.RegisteredWord.
		Update().
		Where(
			registeredword.UserIDEQ(userID),
			registeredword.WordIDEQ(wordID),
		).
		SetIsActive(true).
		Save(ctx)
	return err
}

func (s *ServiceImpl) createReg(
	ctx context.Context,
	tx *ent.Tx,
	userID, wordID int,
) error {
	_, err := tx.RegisteredWord.
		Create().
		SetUserID(userID).
		SetWordID(wordID).
		SetIsActive(true).
		Save(ctx)
	return err
}

/*==================== post process ====================*/

func (s *ServiceImpl) incrementRegCount(
	ctx context.Context,
	okWords []string,
) error {
	_, err := s.RegisteredWordsCount(ctx, true, okWords)
	return err
}
