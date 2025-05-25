// src/repository/ent/word_repo_ent.go
package ent

import (
	"context"
	"math/rand"

	ent "word_app/backend/ent"
	"word_app/backend/ent/registeredword"
	"word_app/backend/ent/word"
	"word_app/backend/ent/wordinfo"
	"word_app/backend/src/domain"
	rep "word_app/backend/src/repository"

	"entgo.io/ent/dialect/sql"
)

type wordRepoEnt struct {
	c   *ent.Client
	rng *rand.Rand
}

func NewWordRepoEnt(c *ent.Client, r *rand.Rand) rep.WordRepo {
	return &wordRepoEnt{c: c, rng: r}
}

// ------------------------- ① ランダム候補取得 -------------------------
func (r *wordRepoEnt) RandomSelectableWords(
	ctx context.Context,
	userID int,
	f rep.WordFilter,
	limit int,
) ([]domain.Word, error) {

	q := r.c.Word.Query().
		Where(word.HasWordInfosWith(
			wordinfo.PartOfSpeechIDIn(f.PartsOfSpeech...),
			wordinfo.HasJapaneseMeans(),
		))

	// ―― 登録 / 未登録 フィルタ ―――――――――――――――――――――――――
	switch f.RegisteredMode {
	case 1: // 登録のみ
		q = q.Where(word.HasRegisteredWordsWith(
			registeredword.UserID(userID),
			registeredword.IsActiveEQ(true),
			registeredword.CorrectRateLTE(f.MaxCorrectRate),
			registeredword.AttentionLevelIn(f.AttentionLevels...),
		))
	case 2: // 未登録のみ
		q = q.Where(word.Not(word.HasRegisteredWordsWith(
			registeredword.UserID(userID),
			registeredword.IsActiveEQ(true))))
	}

	// ―― 慣用句 / 特殊文字 ―――――――
	if f.IncludeIdioms != 0 {
		q = q.Where(word.IsIdiomsEQ(f.IncludeIdioms == 1))
	}
	if f.IncludeSpecial != 0 {
		q = q.Where(word.IsSpecialCharactersEQ(f.IncludeSpecial == 1))
	}

	// ―― ランダム & Limit ―――――――
	words, err := q.
		Order(func(s *sql.Selector) { s.OrderBy("RANDOM()") }).
		Limit(limit).
		WithWordInfos(func(wi *ent.WordInfoQuery) { wi.WithJapaneseMeans() }).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// ―― Entモデル → ドメインモデル 変換 ――――――
	out := make([]domain.Word, len(words))
	for i, w := range words {
		out[i] = domain.FromEntWord(w) // ★ domain パッケージに 1 行 Mapper を置く
	}
	return out, nil
}
