package settingUc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	// mockery 生成先を調整してください
	mockRepo "word_app/backend/src/mocks/infrastructure/repository/setting"
)

func TestGetUserConfigInteractor_Execute(t *testing.T) {
	// 共有 sentinel 以外のエラーは変数で共通化
	repoErr := errors.New("db fail")

	tests := []struct {
		name    string
		repoCfg *domain.UserConfig
		repoErr error
		wantErr error
	}{
		{
			name:    "正常取得",
			repoCfg: &domain.UserConfig{ID: 1, UserID: 42, IsDarkMode: true},
			wantErr: nil,
		},
		{
			name:    "repo エラー → ErrUserConfigNotFound",
			repoErr: repoErr,
			wantErr: settingUc.ErrUserConfigNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepo.NewMockUserConfigRepository(t)
			repo.
				On("GetByUserID", mock.Anything, 42).
				Return(tt.repoCfg, tt.repoErr)

			uc := settingUc.NewGetUserConfig(repo)
			out, err := uc.Execute(context.Background(), settingUc.InputGetUserConfig{UserID: 42})

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, out)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, true, out.Config.IsDarkMode)
			}
			repo.AssertExpectations(t)
		})
	}
}
