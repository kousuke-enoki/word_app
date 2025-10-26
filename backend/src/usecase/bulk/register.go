// src/usecase/bulk/register.go
package bulk

import (
	"context"
	"strings"

	"word_app/backend/config"
	"word_app/backend/src/infrastructure/repository/registeredword"
	"word_app/backend/src/infrastructure/repository/tx"
	"word_app/backend/src/infrastructure/repository/user"
	"word_app/backend/src/infrastructure/repository/word"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/shared/ucerr"
)

// 返却は既存 models に合わせる
type RegisterUsecase interface {
	Register(ctx context.Context, userID int, words []string) (*models.BulkRegisterResponse, error)
}

type registerUsecase struct {
	wordsRead word.ReadRepository
	rwRead    registeredword.ReadRepository
	rwWrite   registeredword.WriteRepository
	txm       tx.Manager
	locker    user.Repository
	limits    *config.LimitsCfg
}

func NewRegisterUsecase(
	wordsRead word.ReadRepository,
	rwRead registeredword.ReadRepository,
	rwWrite registeredword.WriteRepository,
	txm tx.Manager,
	locker user.Repository,
	limits *config.LimitsCfg,
) RegisterUsecase {
	return &registerUsecase{
		wordsRead: wordsRead,
		rwRead:    rwRead,
		rwWrite:   rwWrite,
		txm:       txm,
		locker:    locker,
		limits:    limits,
	}
}

func (uc *registerUsecase) Register(
	ctx context.Context,
	userID int,
	payload []string,
) (
	*models.BulkRegisterResponse, error,
) {
	// 0) リクエスト件数制限（1リクエストの上限。総量200とは別）
	maxPerReq := uc.limits.BulkRegisterMaxItems
	if maxPerReq <= 0 {
		maxPerReq = 200
	}
	if len(payload) == 0 {
		return nil, ucerr.BadRequest("empty payload")
	}
	if len(payload) > maxPerReq {
		return nil, ucerr.BadRequest("too many words in request")
	}

	// 1) 正規化 + 1リクエスト上限
	norm := normalizeLowerUnique(payload, maxPerReq)

	// 2) master word → map[name]id
	nameToID, err := uc.wordsRead.FindIDsByNames(ctx, norm)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(nameToID))
	for _, id := range nameToID {
		ids = append(ids, id)
	}

	// 3) Tx
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	commit := false
	defer func() { _ = done(commit) }()

	// 4) ユーザーロック（同一ユーザーの登録競合を直列化）
	if err := uc.locker.LockByID(txCtx, userID); err != nil {
		return nil, err
	}

	// 5) 現在の active 総数と 残り枠
	maxTotal := uc.limits.RegisteredWordsPerUser
	if maxTotal <= 0 {
		maxTotal = 200
	}
	curActive, err := uc.rwRead.CountActiveByUser(txCtx, userID)
	if err != nil {
		return nil, err
	}
	remain := maxTotal - curActive
	if remain <= 0 {
		// 全部枠超過
		failed := make([]models.FailedWord, 0, len(norm))
		for _, w := range norm {
			failed = append(failed, models.FailedWord{Word: w, Reason: "limit_reached"})
		}
		return &models.BulkRegisterResponse{Success: nil, Failed: failed}, nil
	}

	// 6) 既存( active/非active )の状態取得
	existMap, err := uc.rwRead.FindActiveMapByUserAndWordIDs(txCtx, userID, ids)
	if err != nil {
		return nil, err
	}

	// 7) 登録実行（枠を超えない範囲で）
	var okWords []string
	var failed []models.FailedWord

	for _, w := range norm {
		wordID, ok := nameToID[w]
		if !ok {
			failed = append(failed, models.FailedWord{Word: w, Reason: "not_exists"})
			continue
		}
		if remain <= 0 {
			failed = append(failed, models.FailedWord{Word: w, Reason: "limit_reached"})
			continue
		}
		// 既に行はある
		if isActive, had := existMap[wordID]; had {
			if isActive {
				failed = append(failed, models.FailedWord{Word: w, Reason: "already_registered"})
				continue
			}
			// 非アクティブ → アクティブ化
			if err := uc.rwWrite.Activate(txCtx, userID, wordID); err != nil {
				failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
				continue
			}
			okWords = append(okWords, w)
			existMap[wordID] = true // 状態更新
			remain--
			continue
		}
		// 新規作成
		if err := uc.rwWrite.CreateActive(txCtx, userID, wordID); err != nil {
			// UNIQUE競合(同時実行) なら Activate 試行でリカバリしても良い
			failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
			continue
		}
		okWords = append(okWords, w)
		existMap[wordID] = true
		remain--
	}

	// 8) commit
	commit = true
	if err := done(commit); err != nil {
		return nil, err
	}

	return &models.BulkRegisterResponse{Success: okWords, Failed: failed}, nil
}

func normalizeLowerUnique(words []string, limit int) []string {
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
