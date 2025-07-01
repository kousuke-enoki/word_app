package tx

import "context"

type TxManager interface {
	WithTx(ctx context.Context, f func(ctx context.Context) error) error
}
