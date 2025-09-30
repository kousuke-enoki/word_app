// app/usecase/user/delete.go
package user

import (
	"context"
	"time"

	"word_app/backend/src/interfaces/http/user"
)

func (uc *UserUsecase) Delete(ctx context.Context, in user.DeleteUserInput) error {
	// Tx開始（既存Txがあればjoinされる実装が理想）
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return err
	}
	commit := false
	defer func() { _ = done(commit) }()

	// 1) Editor/Target の最小情報取得（ID, IsRoot, DeletedAt is nil の想定）
	editor, err := uc.userRepo.FindByID(txCtx, in.EditorID)
	if err != nil {
		// NotFound は権限なしと同義
		return err
	}
	target, err := uc.userRepo.FindByID(txCtx, in.TargetID)
	if err != nil {
		return err
	}

	// 2) ポリシー判定（仕様通り）
	if editor.IsRoot {
		// root: 自分以外 かつ 相手が root でないこと
		if target.ID == editor.ID {
			return err // rootの自分削除は不可
		}
		if target.IsRoot {
			return err // rootをrootは削除不可
		}
	} else {
		// 非root: 自分のみ削除可
		if target.ID != editor.ID {
			return err
		}
	}

	// 3) 関連も含め論理削除（同一Tx内・任意の時刻を統一）
	now := time.Now()

	// 3-1) 外部認証情報（例：LINE/OIDC）
	if err := uc.authRepo.SoftDeleteByUserID(txCtx, target.ID, now); err != nil {
		return err
	}

	// 3-2) user_config
	if err := uc.settingRepo.SoftDeleteByUserID(txCtx, target.ID, now); err != nil {
		return err
	}

	// 3-3) user（最後に削除）
	if err := uc.userRepo.SoftDeleteByID(txCtx, target.ID, now); err != nil {
		return err
	}

	commit = true
	if err := done(commit); err != nil {
		return err
	}
	return nil
}
