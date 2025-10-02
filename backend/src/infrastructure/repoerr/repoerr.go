// infra/repoerr/repoerr.go
package repoerr

import (
	"context"
	"errors"
	"fmt"

	"word_app/backend/ent"
	"word_app/backend/ent/privacy" // 使っていれば
	"word_app/backend/src/usecase/apperror"
)

// Entの生エラーをここで終了。Usecase以上にEntが漏れない。
// メッセージはUIに出して良い短文（またはi18nキー）にしておく。
func FromEnt(err error, notFoundMsg, conflictMsg string) error {
	if err == nil {
		return nil
	}
	switch {
	case ent.IsNotFound(err):
		return apperror.NotFoundf(notFoundMsg, err)
	case ent.IsConstraintError(err):
		return apperror.Conflictf(conflictMsg, err)
	case errors.Is(err, context.DeadlineExceeded):
		return apperror.Internalf("database timeout", err)
	case errors.Is(err, privacy.Deny): // 使っていなければこのケースは削除
		return apperror.Forbiddenf("forbidden", err)
	default:
		return apperror.Internalf("database error", err)
	}
}

// 共通保存ラッパ：各Repoでの Save/Update/Exec 後に呼ぶだけ
func SaveErr(err error) error {
	return FromEnt(err, "resource not found", "conflict")
}

// 例：メッセージを具体化して使いたい場合
func NotFound(entity string) error {
	return apperror.NotFoundf(fmt.Sprintf("%s not found", entity), nil)
}
