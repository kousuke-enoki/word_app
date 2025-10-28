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
) (*models.BulkRegisterResponse, error) {
	// 0) リクエスト件数制限（元挙動と同じ）
	if err := uc.validatePayloadLen(payload); err != nil {
		return nil, err
	}

	// 1) 正規化 + 1リクエスト上限（元の normalize 維持）
	maxPerReq := uc.effectiveMaxPerReq()
	norm := normalizeLowerUnique(payload, maxPerReq)

	// 2) master word → map[name]id（同一）
	nameToID, ids, err := uc.fetchWordIDs(ctx, norm)
	if err != nil {
		return nil, err
	}

	// 3) Tx 開始（既存と同じ：既定 rollback）
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = done(false) }()

	// 4) ユーザーロック（同一 tx 上で直列化）
	if err := uc.locker.LockByID(txCtx, userID); err != nil {
		return nil, err
	}

	// 5) 現在の active と残り枠（枠ゼロなら「全て limit_reached」で即返す＝元挙動）
	remain, earlyResp, err := uc.computeRemainOrEarlyFail(txCtx, userID, norm)
	if err != nil {
		return nil, err
	}
	if earlyResp != nil {
		// Success は nil、Failed は全件 limit_reached（元の挙動）
		return earlyResp, nil
	}

	// 6) 既存状態取得（active の有無を map で）
	existMap, err := uc.rwRead.FindActiveMapByUserAndWordIDs(txCtx, userID, ids)
	if err != nil {
		return nil, err
	}

	// 7) 登録実行（ロジックは per-word ヘルパーへ。判定順/文字列/remain 減算のタイミングは元と同じ）
	okWords := make([]string, 0, len(norm))
	failed := make([]models.FailedWord, 0)

	for _, w := range norm {
		ok, reason, err := uc.processOne(txCtx, userID, w, nameToID, existMap, &remain)
		if err != nil {
			// DB 例外は "db_error" に丸める（元の挙動）
			failed = append(failed, models.FailedWord{Word: w, Reason: "db_error"})
			continue
		}
		if ok {
			okWords = append(okWords, w)
			continue
		}
		if reason != "" {
			failed = append(failed, models.FailedWord{Word: w, Reason: reason})
		}
	}

	// 8) commit（同一）
	if err := done(true); err != nil {
		return nil, err
	}
	// 元実装と同様、okWords が 0 件でも nil のまま返る可能性がある
	return &models.BulkRegisterResponse{Success: okWords, Failed: failed}, nil
}

// --- helpers (挙動等価を保つための分割) ---

func (uc *registerUsecase) effectiveMaxPerReq() int {
	if uc.limits.BulkRegisterMaxItems > 0 {
		return uc.limits.BulkRegisterMaxItems
	}
	return 200
}

func (uc *registerUsecase) effectiveMaxTotal() int {
	if uc.limits.RegisteredWordsPerUser > 0 {
		return uc.limits.RegisteredWordsPerUser
	}
	return 200
}

func (uc *registerUsecase) validatePayloadLen(payload []string) error {
	maxPerReq := uc.effectiveMaxPerReq()
	if len(payload) == 0 {
		return ucerr.BadRequest("empty payload")
	}
	if len(payload) > maxPerReq {
		return ucerr.BadRequest("too many words in request")
	}
	return nil
}

func (uc *registerUsecase) fetchWordIDs(ctx context.Context, norm []string) (map[string]int, []int, error) {
	nameToID, err := uc.wordsRead.FindIDsByNames(ctx, norm)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int, 0, len(nameToID))
	for _, id := range nameToID {
		ids = append(ids, id)
	}
	return nameToID, ids, nil
}

// remain<=0 のときは元実装と同じ「Success=nil、Failed=全件 limit_reached」を返す
func (uc *registerUsecase) computeRemainOrEarlyFail(
	ctx context.Context,
	userID int,
	norm []string,
) (int, *models.BulkRegisterResponse, error) {
	maxTotal := uc.effectiveMaxTotal()
	curActive, err := uc.rwRead.CountActiveByUser(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	remain := maxTotal - curActive
	if remain <= 0 {
		failed := make([]models.FailedWord, 0, len(norm))
		for _, w := range norm {
			failed = append(failed, models.FailedWord{Word: w, Reason: "limit_reached"})
		}
		return 0, &models.BulkRegisterResponse{Success: nil, Failed: failed}, nil
	}
	return remain, nil, nil
}

// 1件分の処理（判定順を元コードと完全一致）
// 戻り値: ok=true なら成功、ok=false かつ reason != "" なら失敗理由、err は DB 例外など（元では "db_error" に丸め）
func (uc *registerUsecase) processOne(
	ctx context.Context,
	userID int,
	w string,
	nameToID map[string]int,
	existMap map[int]bool,
	remain *int,
) (ok bool, reason string, err error) {
	wordID, okName := nameToID[w]
	if !okName {
		return false, "not_exists", nil
	}
	if *remain <= 0 {
		return false, "limit_reached", nil
	}

	// 既存行あり
	if isActive, had := existMap[wordID]; had {
		// すでに active
		if isActive {
			return false, "already_registered", nil
		}
		// 非アクティブ -> アクティブ化
		if err := uc.rwWrite.Activate(ctx, userID, wordID); err != nil {
			return false, "db_error", nil // 呼び出し元で "db_error" に丸めているため、ここも合わせる
		}
		existMap[wordID] = true
		*remain--
		return true, "", nil
	}

	// 新規作成
	if err := uc.rwWrite.CreateActive(ctx, userID, wordID); err != nil {
		// UNIQUE競合などは元コードと同様 "db_error"
		return false, "db_error", nil
	}
	existMap[wordID] = true
	*remain--
	return true, "", nil
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
