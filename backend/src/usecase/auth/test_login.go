// src/usecase/auth/test_login.go
package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"word_app/backend/src/domain"
	"word_app/backend/src/usecase/shared/ucerr"

	"github.com/google/uuid"
)

// type TestLoginInput struct {
// 	Jump      string // list|bulk|quiz
// 	Size      int    // e.g. 10
// 	IP        string
// 	UAHash    string
// 	Now       time.Time
// 	RequestID string
// }

type TestLoginOutput struct {
	Token    string `json:"token"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Jump     string `json:"jump"`
}

func normalizeJump(j string) string {
	switch j {
	case "list", "bulk", "quiz":
		return j
	default:
		return "quiz"
	}
}

func (uc *AuthUsecase) TestLoginWithRateLimit(
	ctx context.Context,
	ip, uaHash, route, jump string,
) (*TestLoginOutput, []byte, int, error) {
	// 0) テストモード有効
	rootSetting, err := uc.rootSettingRepo.Get(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	if !rootSetting.IsTestUserMode {
		return nil, nil, 0, ucerr.Forbidden("test user mode is disabled")
	}

	// 1) レート制限チェック
	chk, err := uc.rateLimiter.CheckRateLimit(ctx, ip, uaHash, route)
	if err != nil {
		return nil, nil, 0, err
	}
	if !chk.Allowed {
		if chk.LastPayload != nil {
			// 超過でもラスト結果があれば 200 を返す設計
			// ただし、RetryAfterも返して、フロントエンドで判定できるようにする
			return nil, chk.LastPayload, chk.RetryAfter, nil
		}
		return nil, nil, chk.RetryAfter, ucerr.TooManyRequests("rate limited")
	}
	// 許可だが、近傍に成功レスポンスがあればそれを返す
	if chk.LastPayload != nil {
		return nil, chk.LastPayload, 0, nil
	}

	// 2) ユーザー作成（Tx内）
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	commit := false
	defer func() { _ = done(commit) }()

	uid := uuid.NewString()
	name := fmt.Sprintf("テストユーザー@%s", uid[:8])
	email := fmt.Sprintf("test+%s@taplex.local", uid)
	var pass *string = nil

	u, err := domain.NewUser(name, &email, pass)
	if err != nil {
		return nil, nil, 0, err
	}
	u.IsTest, u.IsAdmin, u.IsRoot = true, false, false

	createdUser, err := uc.userRepo.Create(txCtx, u)
	if err != nil {
		return nil, nil, 0, err
	}
	if err := uc.settingRepo.CreateDefault(txCtx, createdUser.ID); err != nil {
		return nil, nil, 0, err
	}
	now := uc.clock.Now()
	if err := uc.userDailyUsageRepo.CreateIfNotExists(txCtx, createdUser.ID, now); err != nil {
		return nil, nil, 0, err
	}
	commit = true
	if err := done(commit); err != nil {
		return nil, nil, 0, err
	}

	// 3) JWT発行
	token, err := uc.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", createdUser.ID))
	if err != nil {
		return nil, nil, 0, err
	}

	out := &TestLoginOutput{
		Token:    token,
		UserID:   createdUser.ID,
		UserName: createdUser.Name,
		Jump:     normalizeJump(jump),
	}

	// 4) last_result を保存（JSON化して汎用化）
	if b, err := json.Marshal(out); err == nil {
		_ = uc.rateLimiter.SaveLastResult(ctx, ip, uaHash, route, b)
	}

	return out, nil, 0, nil
}
