package settingUc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"word_app/backend/src/domain"
	settingUc "word_app/backend/src/usecase/setting"

	// mockery 生成先を合わせてください
	mockRoot "word_app/backend/src/mocks/infrastructure/repository/setting"
	mockUser "word_app/backend/src/mocks/infrastructure/repository/user"
)

func TestGetRootConfigInteractor_Execute(t *testing.T) {
	var (
		dbErr = errors.New("db down") // ★ 共有エラー
		// errors.New("db down")をもし 2 回 呼びだすと、wantErr とモックが返す err は 別インスタンス になる。
		// errors.Is は ポインタ一致 もしくは Is() 実装で比較するため、インスタンスが違うと一致せずエラー。
	)
	rootCfg := &domain.RootConfig{ID: 1, EditingPermission: "admin"}

	tests := []struct {
		name      string
		userFound *domain.User
		userErr   error
		rootCfg   *domain.RootConfig
		rootErr   error
		wantErr   error
	}{
		{
			name:      "正常 (root user + config取得)",
			userFound: &domain.User{ID: 99, IsRoot: true},
			rootCfg:   rootCfg,
			wantErr:   nil,
		},
		{
			name:    "ユーザ取得エラーを透過",
			userErr: dbErr,
			wantErr: dbErr,
		},
		{
			name:      "root 権限なし → ErrUnauthorized",
			userFound: &domain.User{ID: 2, IsRoot: false},
			wantErr:   settingUc.ErrUnauthorized,
		},
		{
			name:      "root config 取得失敗 → ErrRootConfigMissing",
			userFound: &domain.User{ID: 99, IsRoot: true},
			rootErr:   errors.New("not found"),
			wantErr:   settingUc.ErrRootConfigMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := mockUser.NewMockUserRepository(t)
			rc := mockRoot.NewMockRootConfigRepository(t)

			// --- UserRepo.Expect ---
			ur.
				On("FindByID", mock.Anything, 99).
				Return(tt.userFound, tt.userErr)

			// --- RootConfigRepo.Expect (only when user OK & root) ---
			if tt.userErr == nil && tt.userFound != nil && tt.userFound.IsRoot {
				rc.
					On("Get", mock.Anything).
					Return(tt.rootCfg, tt.rootErr)
			}

			uc := settingUc.NewGetRootConfig(ur, rc)
			out, err := uc.Execute(context.Background(), settingUc.InputGetRootConfig{UserID: 99})

			// 期待エラー
			if tt.wantErr != nil {
				// 特定の sentinel なら ErrorIs、それ以外はメッセージ比較
				switch tt.wantErr {
				case settingUc.ErrUnauthorized, settingUc.ErrRootConfigMissing:
					assert.ErrorIs(t, err, tt.wantErr)
				default:
					assert.EqualError(t, err, tt.wantErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, rootCfg.EditingPermission, out.Config.EditingPermission)
			}

			ur.AssertExpectations(t)
			rc.AssertExpectations(t)
		})
	}
}
