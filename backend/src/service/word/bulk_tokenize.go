package word_service

import (
	"context"
	"regexp"
	"strings"

	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
)

// 正規表現で単語抽出 → 小文字化 → 重複排除
var reWord = regexp.MustCompile(`[A-Za-z]+(?:'[A-Za-z]+)?`)

func unique(words []string) []string {
	set := map[string]struct{}{}
	out := make([]string, 0, len(words))
	for _, w := range words {
		lw := strings.ToLower(w)
		if _, ok := set[lw]; !ok {
			set[lw] = struct{}{}
			out = append(out, lw)
		}
	}
	return out
}

func (s *WordServiceImpl) BulkTokenize(ctx context.Context, userID int, text string) ([]string, []string, []string, error) {
	words := unique(reWord.FindAllString(text, -1))
	if len(words) == 0 {
		return nil, nil, nil, nil
	}

	/* ---------- DB 検索 ---------- */
	// ① master word テーブルに存在する単語
	// foundWords, err := s.db.Word.
	foundWords, err := s.client.Word().
		Query().
		Where(word.NameIn(words...)).
		All(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	nameToID := make(map[string]int)
	for _, w := range foundWords {
		nameToID[w.Name] = w.ID
	}

	/* ---------- registered_word 取得 (active=true) ---------- */
	var ids []int
	for _, id := range nameToID {
		ids = append(ids, id)
	}
	activeRegs, err := s.client.RegisteredWord().
		Query().
		Where(
			registeredword.UserIDEQ(userID),
			registeredword.WordIDIn(ids...),
			registeredword.IsActive(true),
		).
		All(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	activeSet := map[int]struct{}{}
	for _, r := range activeRegs {
		activeSet[r.WordID] = struct{}{}
	}

	/* ---------- 仕分け ---------- */
	var cand, reg, notExist []string
	for _, w := range words {
		id, ok := nameToID[w]
		if !ok {
			notExist = append(notExist, w)
			continue
		}
		if _, ok := activeSet[id]; ok {
			reg = append(reg, w)
		} else {
			cand = append(cand, w)
		}
	}
	return cand, reg, notExist, nil
}
