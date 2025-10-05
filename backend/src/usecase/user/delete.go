// app/usecase/user/delete.go
package user

import (
	"context"
	"time"

	"word_app/backend/src/usecase/shared/ucerr"
)

func (uc *UserUsecase) Delete(ctx context.Context, in DeleteUserInput) error {
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

	// 2) ポリシー判定（機能要件）
	if err := authorizeDelete(editor.ID, editor.IsRoot, target.ID, target.IsRoot); err != nil {
		return err
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

// 認可・ポリシーを分離（ucerr で返す）
func authorizeDelete(editorID int, isEditorRoot bool, targetID int, isTargetRoot bool) error {
	if isEditorRoot {
		// 仕様：rootは「自分以外」かつ「相手がrootではない」
		if targetID == editorID {
			// セキュリティ上のふるまい：対象の存在や属性を隠したいなら NotFound にする、
			// 明示的に拒否したいなら Forbidden にする。どちらでもOK。ここでは Forbidden を採用。
			return ucerr.Forbidden("forbidden")
		}
		if isTargetRoot {
			return ucerr.Forbidden("forbidden")
		}
		return nil
	}
	// 一般ユーザーは自分のみ削除可
	if targetID != editorID {
		// 対象の存在有無を伏せたいなら NotFound を返す選択もある（ユーザー列挙対策）
		return ucerr.Forbidden("forbidden")
	}
	return nil
}
