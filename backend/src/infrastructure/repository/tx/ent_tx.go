// infrastructure/tx/ent_tx.go
package tx

import (
	"context"

	"word_app/backend/src/interfaces/service_interfaces"
)

type EntTxManager struct {
	client service_interfaces.EntClientInterface
}

type TxManager interface {
	WithTx(ctx context.Context, f func(ctx context.Context) error) error
}

func NewEntTxManager(client service_interfaces.EntClientInterface) *EntTxManager {
	return &EntTxManager{client: client}
}

func (m *EntTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
