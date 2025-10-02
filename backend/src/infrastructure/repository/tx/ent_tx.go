// infrastructure/repository/tx/ent_tx.go
package tx

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/logger/logx"
	"word_app/backend/src/infrastructure/repoerr"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

// ===== ctx helpers (exported) =====

type ctxKeyTx struct{}

// TxFromContext は ctx に注入された *ent.Tx を取り出します（リポジトリの getDB から使用）
func TxFromContext(ctx context.Context) (*ent.Tx, bool) {
	tx, ok := ctx.Value(ctxKeyTx{}).(*ent.Tx)
	return tx, ok
}

// WithTxContext は ctx に *ent.Tx を注入します（TxManager 内部で使用）
func WithTxContext(ctx context.Context, tx *ent.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

type EntManager struct {
	client serviceinterfaces.EntClientInterface // *ent.Client 互換のinterface
}

type Manager interface {
	WithTx(ctx context.Context, f func(ctx context.Context) error) error
	Begin(ctx context.Context) (context.Context, func(commit bool) error, error) // 明示的UoWも提供（後述）
}

func NewEntManager(client serviceinterfaces.EntClientInterface) *EntManager {
	return &EntManager{client: client}
}

// WithTx は、与えられた関数をトランザクション内で実行します。
// すでにトランザクションが存在する場合は、新しいトランザクションを開始せずにそのまま関数を実行します。
// 関数がエラーを返した場合、トランザクションはロールバックされ、nilを返した場合はコミットされます。
// panicが発生した場合もロールバックされ、panicは再度発生します。
// “短い一連の更新”にはとても有効
func (m *EntManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	// すでにTx中ならjoin（新規開始しない）
	if parentTx, ok := TxFromContext(ctx); ok && parentTx != nil {
		return fn(ctx)
	}

	// 新規開始
	tx, err := m.client.Tx(ctx)
	if err != nil {
		// begin失敗は apperror に正規化して返す
		return repoerr.FromEnt(err, "failed to begin transaction", "")
	}
	txCtx := WithTxContext(ctx, tx)

	defer func() {
		// panic でも rollback を試みてから再panic
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logx.From(ctx).WithError(rbErr).Error("[tx] rollback after panic failed")
			}
			panic(p)
		}
	}()

	// body
	if err = fn(txCtx); err != nil {
		// fn 側のエラーは“意味”が付いている（apperror想定）ので、そのまま返す
		// ただし rollback に失敗したらログに残す
		if rbErr := tx.Rollback(); rbErr != nil {
			logx.From(txCtx).WithError(rbErr).WithField("cause", err).Error("[tx] rollback after error failed")
		}
		return err
	}

	// commit
	if err = tx.Commit(); err != nil {
		// commit失敗は apperror に正規化（通常 Internal）
		return repoerr.FromEnt(err, "failed to commit transaction", "")
	}
	return nil
}

// 併設: 明示的UoW（Begin/Commit）も使えるようにしておく
// 既存Txがあればjoinする（commit=false扱い固定）
// 既存Txがなければ新規開始し、cleanupでCommit/Rollbackを選択できる
// “長い一連の更新”や、複数のUsecaseを跨る更新にはこちらが有効
// 可読性は落ちるがシンプルに実装できる。
func (m *EntManager) Begin(ctx context.Context) (context.Context, func(commit bool) error, error) {
	// 既存TxにjoinするUoW: commit=false扱い固定
	if existing, ok := TxFromContext(ctx); ok && existing != nil {
		return ctx, func(commit bool) error { return nil }, nil
	}

	// begin
	tx, err := m.client.Tx(ctx)
	if err != nil {
		return ctx, nil, repoerr.FromEnt(err, "failed to begin transaction", "")
	}
	txCtx := WithTxContext(ctx, tx)

	cleanup := func(commit bool) error {
		if !commit {
			if rbErr := tx.Rollback(); rbErr != nil {
				// cleanupは呼び側の戻り値として返す（Handlerで500にマップされる）
				return repoerr.FromEnt(rbErr, "failed to rollback transaction", "")
			}
			return nil
		}
		if cmErr := tx.Commit(); cmErr != nil {
			return repoerr.FromEnt(cmErr, "failed to commit transaction", "")
		}
		return nil
	}
	return txCtx, cleanup, nil
}
