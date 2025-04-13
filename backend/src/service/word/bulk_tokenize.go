package word_service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
)

// 正規表現で単語抽出 → 小文字化 → 重複排除
var reWord = regexp.MustCompile(`[A-Za-z]+(?:'[A-Za-z]+)?`)

const maxTokens = 200 // 1 リクエストで許可する最大単語数

func unique(words []string) []string {
	set := make(map[string]struct{}, len(words))
	out := make([]string, 0, maxTokens) // 上限ぶんだけ確保

	for _, w := range words {
		if len(out) >= maxTokens { // 上限超えたらここで打ち切り
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

func (s *WordServiceImpl) BulkTokenize(ctx context.Context, userID int, text string) ([]string, []string, []string, error) {
	wordsRaw := reWord.FindAllString(text, -1)
	if len(wordsRaw) == 0 {
		return nil, nil, nil, nil
	}
	if len(wordsRaw) > maxTokens*5 { // ★ さらに荒いフィルタ
		return nil, nil, nil, fmt.Errorf("too many tokens: %d > %d", len(wordsRaw), maxTokens*5)
	}

	words := unique(wordsRaw) // ここで重複排除＆maxTokens適用

	if len(words) == 0 {
		return nil, nil, nil, nil
	}

	/* ---------- DB 検索 ---------- */
	// ① master word テーブルに存在する単語
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
