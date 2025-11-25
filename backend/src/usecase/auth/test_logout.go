// src/usecase/auth/test_logout.go
package auth

import (
	"context"

	"word_app/backend/src/usecase/shared/ucerr"
)

// type userRepo interface {
// 	FindByIDForUpdate(ctx context.Context, id int) (*UserLite, error)
// 	Delete(ctx context.Context, id int) error
// }

type UserLite struct {
	ID     int
	IsTest bool
}

// TestLogout: テストユーザー本人だけが、自分のアカウントを物理削除する。
// 子テーブルはDBの ON DELETE CASCADE が実行。
// 冪等: 既に削除済みでもエラーにしない（正常終了）。

func (uc *AuthUsecase) TestLogout(ctx context.Context, actorID int) error {
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return err
	}
	commit := false
	defer func() { _ = done(commit) }()

	// 1発で「is_test かつ id 一致」だけを削除（原子的）

	deleted, err := uc.userRepo.DeleteIfTest(txCtx, actorID)
	if err != nil {
		return err
	}

	if deleted {
		commit = true
		if err := done(commit); err != nil {
			return err
		}
		// ユーザー削除後に、レート制限キャッシュもクリア
		// トランザクション外で実行（キャッシュは外部サービスなので）
		_ = uc.rateLimiter.ClearCacheForUser(ctx, actorID)
		return nil
	}

	// 削除されなかった理由を区別したい場合のみ判定（不要ならここは省く）
	exists, err := uc.userRepo.Exists(txCtx, actorID)
	if err != nil {
		return err
	}
	if !exists {
		// 既に削除済み（冪等）→ 成功扱い
		commit = true
		return done(commit)
	}
	isTest, err := uc.userRepo.IsTest(txCtx, actorID)
	if err != nil {
		return err
	}
	if !isTest {
		return ucerr.Forbidden("only test user can be deleted via test-logout")
	}

	// ここには基本来ない（DeleteIfTest が既に判定済み）
	commit = true
	return done(commit)
}
