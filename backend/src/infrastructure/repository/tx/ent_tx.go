// infrastructure/tx/ent_tx.go
package tx

import (
	"context"
	"log"

	"word_app/backend/ent"
	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type ctxKeyTx struct{}

func txFromContext(ctx context.Context) (*ent.Tx, bool) {
	tx, ok := ctx.Value(ctxKeyTx{}).(*ent.Tx)
	return tx, ok
}

func withTx(ctx context.Context, tx *ent.Tx) context.Context {
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
	if parentTx, ok := txFromContext(ctx); ok && parentTx != nil {
		return fn(ctx)
	}

	// 新規開始
	tx, err := m.client.Tx(ctx)
	if err != nil {
		return err
	}
	txCtx := withTx(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[tx] rollback after panic failed: %v", rbErr)
			}
			panic(p)
		}
	}()

	if err = fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("[tx] rollback after error failed: %v (original: %v)", rbErr, err)
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// 併設: 明示的UoW（Begin/Commit）も使えるようにしておく
// 既存Txがあればjoinする（commit=false扱い固定）
// 既存Txがなければ新規開始し、cleanupでCommit/Rollbackを選択できる
// “長い一連の更新”や、複数のUsecaseを跨る更新にはこちらが有効
// 可読性は落ちるがシンプルに実装できる。
func (m *EntManager) Begin(ctx context.Context) (context.Context, func(commit bool) error, error) {
	if existing, ok := txFromContext(ctx); ok && existing != nil {
		// 既存TxにjoinするUoW: commit=false扱い固定
		return ctx, func(commit bool) error { return nil }, nil
	}
	tx, err := m.client.Tx(ctx)
	if err != nil {
		return ctx, nil, err
	}
	txCtx := withTx(ctx, tx)
	cleanup := func(commit bool) error {
		if !commit {
			if rbErr := tx.Rollback(); rbErr != nil {
				return rbErr
			}
			return nil
		}
		return tx.Commit()
	}
	return txCtx, cleanup, nil
}
