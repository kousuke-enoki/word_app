package settinguctest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	mockRepo "word_app/backend/src/mocks/infrastructure/repository/setting"
	mockTx "word_app/backend/src/mocks/infrastructure/repository/tx"
)

func TestUpdateUserConfigInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	input := settingUc.InputUpdateUserConfig{UserID: 42, IsDarkMode: true}
	expectCfg := &domain.UserConfig{UserID: 42, IsDarkMode: true}

	errUpsert := errors.New("upsert fail")
	errTx := errors.New("tx begin fail")

	tests := []struct {
		name      string
		txErr     error // Tx.WithTx が返すエラー
		upsertCfg *domain.UserConfig
		upsertErr error
		expectErr error
	}{
		{
			name:      "正常更新",
			upsertCfg: expectCfg,
		},
		{
			name:      "Upsert 失敗 -> 伝播",
			upsertErr: errUpsert,
			expectErr: errUpsert,
		},
		{
			name:      "Tx 開始失敗",
			txErr:     errTx,
			expectErr: errTx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := mockTx.NewMockManager(t)
			repo := mockRepo.NewMockUserConfigRepository(t)
			txRet := tt.txErr    // Tx 自体のエラー
			if tt.txErr == nil { // Tx 開始は成功
				txRet = tt.upsertErr // -> Upsert の結果を返す
			}
			// Tx.WithTx モック
			tx.
				On("WithTx", ctx, mock.AnythingOfType("func(context.Context) error")).
				Run(func(args mock.Arguments) {
					if tt.txErr == nil { // Tx 成功時のみコールバックを実行
						cb := args.Get(1).(func(context.Context) error)
						_ = cb(context.Background()) // エラーは上で txRet に設定済み
					}
				}).
				Return(txRet)

			// Upsert は Tx が成功するシナリオでのみ期待
			if tt.txErr == nil {
				repo.
					On("Upsert", mock.Anything, mock.AnythingOfType("*domain.UserConfig")).
					Return(tt.upsertCfg, tt.upsertErr)
			}

			uc := settingUc.NewUpdateUserConfig(tx, repo)
			out, err := uc.Execute(ctx, input)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Nil(t, out)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, true, out.IsDarkMode)
			}

			tx.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
