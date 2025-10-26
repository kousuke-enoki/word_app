// src/usecase/bulk/tokenize.go
package bulk

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"word_app/backend/config"
	"word_app/backend/src/infrastructure/repository/registeredword"
	udurepo "word_app/backend/src/infrastructure/repository/userdailyusage"
	"word_app/backend/src/infrastructure/repository/word"
	"word_app/backend/src/usecase/clock"
)

// 正規表現で単語抽出（英字と省略形）
// 例: don't -> 1語として扱う
var reWord = regexp.MustCompile(`[A-Za-z]+(?:'[A-Za-z]+)?`)

type TokenizeUsecase interface {
	Execute(ctx context.Context, userID int, text string) (cands, regs, notExist []string, err error)
}

type tokenizeUsecase struct {
	wordRepo           word.ReadRepository
	regReadRepo        registeredword.ReadRepository
	userDailyUsageRepo udurepo.Repository
	clock              clock.Clock
	limits             *config.LimitsCfg
	maxTokens          int // 1リクエストで扱う最大ユニーク語数（デフォルト200）
}

func NewTokenizeUsecase(
	wordRepo word.ReadRepository,
	regReadRepo registeredword.ReadRepository,
	userDailyUsageRepo udurepo.Repository,
	clock clock.Clock,
	limits *config.LimitsCfg, // nil ならデフォルト
) TokenizeUsecase {
	max := 200
	if limits != nil && limits.BulkTokenizeMaxTokens > 0 {
		max = limits.BulkTokenizeMaxTokens
	}
	return &tokenizeUsecase{
		wordRepo:           wordRepo,
		regReadRepo:        regReadRepo,
		userDailyUsageRepo: userDailyUsageRepo,
		clock:              clock,
		limits:             limits,
		maxTokens:          max,
	}
}

func (uc *tokenizeUsecase) Execute(
	ctx context.Context,
	userID int,
	text string,
) (
	[]string,
	[]string,
	[]string,
	error,
) {
	// 日次クォータ上限取得
	cap := uc.limits.BulkMaxPerDay
	if cap <= 0 {
		cap = 5
	}

	// 0) 日次クォータ消費（原子的に +1、上限なら 429 相当エラーを返す）
	if _, err := uc.userDailyUsageRepo.IncBulkOr429(ctx, userID, uc.clock.Now(), cap); err != nil {
		// ucerr.TooManyRequests を返す
		return nil, nil, nil, err
	}

	// 1) トークン抽出
	raw := reWord.FindAllString(text, -1)
	if len(raw) == 0 {
		return nil, nil, nil, nil
	}
	// 荒いフィルタ：生トークンが上限の5倍超なら早期終了
	if len(raw) > uc.maxTokens*5 {
		return nil, nil, nil, fmt.Errorf("too many tokens: %d > %d", len(raw), uc.maxTokens*5)
	}

	// 2) 正規化（小文字化＋ユニーク＋maxTokens で打ち切り）
	words := uniqueLowerLimited(raw, uc.maxTokens)
	if len(words) == 0 {
		return nil, nil, nil, nil
	}

	// 3) master word から存在する単語の ID を取得
	nameToID, err := uc.wordRepo.FindIDsByNames(ctx, words)
	if err != nil {
		return nil, nil, nil, err
	}

	// 4) user の active registered_word セットを取得
	ids := make([]int, 0, len(nameToID))
	for _, id := range nameToID {
		ids = append(ids, id)
	}
	activeSet, err := uc.regReadRepo.ActiveWordIDSetByUser(ctx, userID, ids)
	if err != nil {
		return nil, nil, nil, err
	}

	// 5) 仕分け：候補/登録済み/存在なし
	var cand, regs, notExist []string
	for _, w := range words {
		id, ok := nameToID[w]
		if !ok {
			notExist = append(notExist, w)
			continue
		}
		if _, ok := activeSet[id]; ok {
			regs = append(regs, w)
		} else {
			cand = append(cand, w)
		}
	}
	return cand, regs, notExist, nil
}

func uniqueLowerLimited(words []string, limit int) []string {
	set := make(map[string]struct{}, limit)
	out := make([]string, 0, limit)
	for _, w := range words {
		if len(out) >= limit {
			break
		}
		lw := strings.ToLower(w)
		if _, ok := set[lw]; ok {
			continue
		}
		set[lw] = struct{}{}
		out = append(out, lw)
	}
	return out
}
