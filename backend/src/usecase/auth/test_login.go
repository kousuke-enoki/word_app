// src/usecase/auth/test_login.go
package auth

import (
	"context"
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
	UserID   int
	UserName string
	Jump     string
}

func (uc *AuthUsecase) TestLogin(ctx context.Context) (*TestLoginOutput, error) {
	// 0) テストモード有効か（RootConfig等から）
	rootSetting, err := uc.rootSettingRepo.Get(ctx) // boolを返す想定
	if err != nil {
		return nil, err
	}
	if !rootSetting.IsTestUserMode {
		return nil, ucerr.Forbidden("test user mode is disabled")
	}

	// 1) ここでレート制限（DynamoDB等）を入れる
	//    uc.rateLimiter.AllowWithLastResult(ctx, in.IP, in.UAHash, "/auth/test-login")

	// 2) ユーザー作成（Tx内）
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	commit := false
	defer func() { _ = done(commit) }()

	uid := uuid.NewString()
	name := fmt.Sprintf("テストユーザー@%s", uid[:8])
	email := fmt.Sprintf("test+%s@taplex.local", uid) // 外部送信不可ドメイン
	var pass *string = nil                            // パスワード無し（サインイン経路不可）

	u, err := domain.NewUser(name, &email, pass)
	if err != nil {
		return nil, err
	}
	u.IsTest = true
	u.IsAdmin = false
	u.IsRoot = false

	createdUser, err := uc.userRepo.Create(txCtx, u)
	if err != nil {
		return nil, err
	}
	// 3) デフォルト設定（UserConfigなど）
	if err := uc.settingRepo.CreateDefault(txCtx, createdUser.ID); err != nil {
		return nil, err
	}
	now := uc.clock.Now()
	if err := uc.userDailyUsageRepo.CreateIfNotExists(txCtx, createdUser.ID, now); err != nil {
		return nil, err
	}

	commit = true
	if err := done(commit); err != nil {
		return nil, err
	}

	return &TestLoginOutput{
		UserID:   createdUser.ID,
		UserName: createdUser.Name,
	}, nil
}
