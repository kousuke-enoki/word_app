package settinguctest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	mockRepo "word_app/backend/src/mocks/infrastructure/repository/setting"
	clock "word_app/backend/src/usecase/clock"
)

func TestGetRuntimeConfigInteractor_Execute(t *testing.T) {
	type want struct {
		isTestUserMode       bool
		isLineAuthentication bool
		version              string
		err                  bool
	}

	// 固定の時刻を使用
	fixedTime := time.Date(2025, 11, 21, 3, 27, 24, 0, time.UTC)

	tests := []struct {
		name     string
		mockCfg  *domain.RootConfig
		mockErr  error
		expected want
	}{
		{
			name:     "Line 認証 ON → true",
			mockCfg:  &domain.RootConfig{IsLineAuthentication: true, UpdatedAt: fixedTime},
			expected: want{isTestUserMode: false, isLineAuthentication: true, version: fixedTime.Format(time.DateTime), err: false},
		},
		{
			name:     "Line 認証 OFF → false",
			mockCfg:  &domain.RootConfig{IsLineAuthentication: false, UpdatedAt: fixedTime},
			expected: want{isTestUserMode: false, isLineAuthentication: false, version: fixedTime.Format(time.DateTime), err: false},
		},
		{
			name:     "repo.Get エラー → そのまま返す",
			mockCfg:  nil,
			mockErr:  errors.New("db down"),
			expected: want{err: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepo.NewMockRootConfigRepository(t)
			repo.
				On("Get", mock.Anything).
				Return(tt.mockCfg, tt.mockErr)

			uc := settingUc.NewRuntimeConfig(repo, clock.SystemClock{})
			config, err := uc.Execute(context.Background())

			if tt.expected.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.isLineAuthentication, config.IsLineAuthentication)
			assert.Equal(t, tt.expected.version, config.Version)
			assert.Equal(t, tt.expected.isTestUserMode, config.IsTestUserMode)
			repo.AssertExpectations(t)
		})
	}
}
