package settinguctest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	// mockery 生成先を合わせてください
	mockRepo "word_app/backend/src/mocks/infrastructure/repository/setting"
)

func TestAuthConfigInteractor_Execute(t *testing.T) {
	type want struct {
		isLineAuth bool
		err        bool
	}

	tests := []struct {
		name     string
		mockCfg  *domain.RootConfig
		mockErr  error
		expected want
	}{
		{
			name:     "Line 認証 ON → true",
			mockCfg:  &domain.RootConfig{IsLineAuthentication: true},
			expected: want{isLineAuth: true, err: false},
		},
		{
			name:     "Line 認証 OFF → false",
			mockCfg:  &domain.RootConfig{IsLineAuthentication: false},
			expected: want{isLineAuth: false, err: false},
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

			uc := settingUc.NewAuthConfig(repo)
			dto, err := uc.Execute(context.Background())

			if tt.expected.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.isLineAuth, dto.IsLineAuth)
			repo.AssertExpectations(t)
		})
	}
}
