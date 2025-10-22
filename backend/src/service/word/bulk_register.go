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

const MaxRegisteredPerUser = 200

func activeCountTx(ctx context.Context, tx *ent.Tx, userID int) (int, error) {
	return tx.RegisteredWord.
		Query().
		Where(registeredword.UserID(userID), registeredword.IsActive(true)).
		Count(ctx)
}

/*==================== Tx + upsert ====================*/

func (s *ServiceImpl) upsertRegsTx(
	ctx context.Context,
	userID int,
	norm []string,
	nameToID map[string]int,
	active map[int]bool, // これは “事前読み取りのスナップショット”。Tx内で再評価する。
) (okWords []string, failed []models.FailedWord, err error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return
	}
	defer finishTxWithLog(&err, tx)

	// 1) ユーザーロック（同一ユーザー操作を直列化）
	if err = s.userRepo.LockByID(ctx, tx, userID); err != nil {
		return
	}

	// 2) Tx内で “最新” の有効数を取得 → 残り枠算出
	curActive, err := activeCountTx(ctx, tx, userID)
	if err != nil {
		return
	}
	remain := MaxRegisteredPerUser - curActive
	if remain <= 0 {
		// もう枠がない：全部 failed にするか、422/409 を返すかは仕様で
		for _, w := range norm {
			failed = append(failed, models.FailedWord{Word: w, Reason: "limit_reached"})
		}
		return // err=nil で 200 部分失敗応答にするのも可
		// もしくは return nil, nil, ucerr.TooManyRequests("registered words limit exceeded")
	}

	// 3) 未登録/非アクティブだけを抽出（Tx時点の最新状態が理想）
	//    事前の active マップは stale の可能性があるので、Tx内でもう一度軽く確認するやり方がより堅いです。
	//    ただし性能と実装コストのトレードオフ。最小修正なら active をそのまま使いつつ UPSERT に寄せてOK。
	addQueue := make([][2]string, 0, len(norm)) // [wordName, reason?] 等
	for _, w := range norm {
		wordID, exists := nameToID[w]
		if !exists {
			failed = append(failed, models.FailedWord{Word: w, Reason: "not_exists"})
			continue
		}
		if active[wordID] { // 事前スナップショット
			failed = append(failed, models.FailedWord{Word: w, Reason: "already_registered"})
			continue
		}
		addQueue = append(addQueue, [2]string{w, ""})
	}

	// 4) 残り枠で切り詰め
	if remain < len(addQueue) {
		// 先頭 remain 件だけ登録、残りは limit_reached
		for i := remain; i < len(addQueue); i++ {
			failed = append(failed, models.FailedWord{Word: addQueue[i][0], Reason: "limit_reached"})
		}
		addQueue = addQueue[:remain]
	}

	// 5) 追加対象を UPSERT（重複/同時実行に強い）
	for _, item := range addQueue {
		w := item[0]
		wordID := nameToID[w]

		// (A) 既存があれば is_active=true にする UPDATE
		// (B) なければ INSERT
		exists, e := tx.RegisteredWord.
			Query().
			Where(
				registeredword.UserIDEQ(userID),
				registeredword.WordIDEQ(wordID),
			).
			Exist(ctx)
		if e != nil {
			// DBエラー
			failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
			continue
		}

		if exists {
			_, e = tx.RegisteredWord.
				Update().
				Where(
					registeredword.UserIDEQ(userID),
					registeredword.WordIDEQ(wordID),
				).
				SetIsActive(true).
				Save(ctx)
			if e != nil {
				failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
				continue
			}
		} else {
			_, e = tx.RegisteredWord.
				Create().
				SetUserID(userID).
				SetWordID(wordID).
				SetIsActive(true).
				Save(ctx)
			if e != nil {
				failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
				continue
			}
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
