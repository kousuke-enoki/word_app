package settinguctest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	mockRoot "word_app/backend/src/mocks/infrastructure/repository/setting"
	mockUser "word_app/backend/src/mocks/infrastructure/repository/user"
)

func TestUpdateRootConfigInteractor_Execute(t *testing.T) {
	ctx := context.Background()
	upsertErr := errors.New("upsert fail")

	input := settingUc.InputUpdateRootConfig{
		UserID:            1,
		EditingPermission: "admin",
		IsTestUserMode:    true,
		IsEmailAuthCheck:  false,
		IsLineAuth:        true,
	}

	expectCfg := &domain.RootConfig{
		EditingPermission:          "admin",
		IsTestUserMode:             true,
		IsEmailAuthenticationCheck: false,
		IsLineAuthentication:       true,
	}

	tests := []struct {
		name      string
		isRootOK  bool
		isRootErr error
		upsertCfg *domain.RootConfig
		upsertErr error
		wantErr   error
	}{
		{
			name:      "正常更新",
			isRootOK:  true,
			upsertCfg: expectCfg,
		},
		{
			name:      "DB エラー (IsRoot)",
			isRootErr: errors.New("db down"),
			wantErr:   settingUc.ErrDatabaseFailure,
		},
		{
			name:     "root 権限なし",
			isRootOK: false,
			wantErr:  settingUc.ErrUnauthorized,
		},
		{
			name:      "Upsert 失敗を返す",
			isRootOK:  true,
			upsertErr: upsertErr,
			wantErr:   upsertErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := mockUser.NewMockRepository(t)
			rr := mockRoot.NewMockRootConfigRepository(t)

			// IsRoot expectation
			ur.
				On("IsRoot", ctx, 1).
				Return(tt.isRootOK, tt.isRootErr)

			// Upsert expectation (only when root check passes)
			if tt.isRootErr == nil && tt.isRootOK {
				rr.
					On("Upsert", ctx, mock.AnythingOfType("*domain.RootConfig")).
					Return(tt.upsertCfg, tt.upsertErr)
			}

			uc := settingUc.NewUpdateRootConfig(rr, ur)
			out, err := uc.Execute(ctx, input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, out)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, true, out.IsTestUserMode)
				assert.Equal(t, "admin", out.EditingPermission)
			}

			ur.AssertExpectations(t)
			rr.AssertExpectations(t)
		})
	}
}
