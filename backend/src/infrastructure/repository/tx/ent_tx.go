// infrastructure/tx/ent_tx.go
package tx

import (
	"context"

	serviceinterfaces "word_app/backend/src/interfaces/service_interfaces"
)

type EntManager struct {
	client serviceinterfaces.EntClientInterface
}

type Manager interface {
	WithTx(ctx context.Context, f func(ctx context.Context) error) error
}

func NewEntManager(client serviceinterfaces.EntClientInterface) *EntManager {
	return &EntManager{client: client}
}

func (m *EntManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
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
