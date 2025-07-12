package tx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"word_app/backend/ent/enttest"

	txrepo "word_app/backend/src/infrastructure/repository/tx"
	mockCli "word_app/backend/src/mocks/service_interfaces"

	_ "github.com/mattn/go-sqlite3"
)

func newMgr(c *mockCli.MockEntClientInterface) *txrepo.EntTxManager {
	return txrepo.NewEntTxManager(c)
}

func TestEntTxManager_WithTx(t *testing.T) {
	ctx := context.Background()
	errBegin := errors.New("begin fail")
	errFn := errors.New("fn fail")

	tests := []struct {
		name        string
		setupMocks  func(*mockCli.MockEntClientInterface)
		fn          func(context.Context) error
		expectPanic bool
		expectErr   error
	}{
		{
			name: "Tx 開始エラーを透過",
			setupMocks: func(cli *mockCli.MockEntClientInterface) {
				cli.On("Tx", ctx).Return(nil, errBegin)
			},
			fn:        func(context.Context) error { return nil },
			expectErr: errBegin,
		},
		{
			name: "fn エラー→Rollback→返却",
			setupMocks: func(cli *mockCli.MockEntClientInterface) {
				realCli := enttest.Open(t, "sqlite3", "file:mem?mode=memory&_fk=1")
				defer realCli.Close()
				txObj, _ := realCli.Tx(ctx)
				cli.On("Tx", ctx).Return(txObj, nil)
			},
			fn:        func(context.Context) error { return errFn },
			expectErr: errFn,
		},
		{
			name: "fn panic→Rollback→panic 伝播",
			setupMocks: func(cli *mockCli.MockEntClientInterface) {
				realCli := enttest.Open(t, "sqlite3", "file:mem?mode=memory&_fk=1")
				defer realCli.Close()
				txObj, _ := realCli.Tx(ctx)
				cli.On("Tx", ctx).Return(txObj, nil)
			},
			fn:          func(context.Context) error { panic("boom") },
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := mockCli.NewMockEntClientInterface(t)
			tt.setupMocks(cli)

			mgr := newMgr(cli)

			if tt.expectPanic {
				assert.Panics(t, func() { _ = mgr.WithTx(ctx, tt.fn) })
				return
			}

			err := mgr.WithTx(ctx, tt.fn)
			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
			cli.AssertExpectations(t)
		})
	}
}
