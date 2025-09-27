// src/service/user/delete.go
package user

import (
	"context"
	"time"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
)

func (e *EntUserClient) Delete(ctx context.Context, editorID, targetID int) error {
	tx, err := e.client.Tx(ctx)
	if err != nil {
		return ErrDatabaseFailure
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	// 編集者（削除実行者）
	editor, err := tx.User.Query().
		Where(
			user.ID(editorID),
			user.DeletedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrUnauthorized // 実質操作不可
		}
		return ErrDatabaseFailure
	}

	// 対象（未削除のみ）
	target, err := tx.User.Query().
		Where(
			user.ID(targetID),
			user.DeletedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrUserNotFound
		}
		return ErrDatabaseFailure
	}

	// 権限チェック
	if editor.IsRoot {
		// root: 自分以外 かつ 相手が root でないこと
		if target.ID == editor.ID {
			return ErrUnauthorized // root の自分削除は不可（仕様）
		}
		if target.IsRoot {
			return ErrUnauthorized // root を root は削除不可（仕様）
		}
		// （任意拡張）最後の root 保護を入れるならここで root 残数を確認
	} else {
		// 非 root: 自分のみ
		if target.ID != editor.ID {
			return ErrUnauthorized
		}
	}

	// 論理削除
	_, err = tx.User.UpdateOneID(target.ID).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return ErrDatabaseFailure
	}

	if err := tx.Commit(); err != nil {
		return ErrDatabaseFailure
	}
	tx = nil
	return nil
}
