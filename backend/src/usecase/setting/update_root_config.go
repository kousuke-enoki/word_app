// usecase/setting/update_root_config.go
package settinguc

import (
	"context"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/infrastructure/repository/setting"
	userRepo "word_app/backend/src/infrastructure/repository/user"
)

type InputUpdateRootConfig struct {
	UserID            int
	EditingPermission string `json:"editing_permission"`
	IsTestUserMode    bool   `json:"is_test_user_mode"`
	IsEmailAuthCheck  bool   `json:"is_email_authentication_check"`
	IsLineAuth        bool   `json:"is_line_authentication"`
}

type UpdateRootConfigInteractor struct {
	rootRepo settingRepo.RootConfigRepository
	userRepo userRepo.Repository
}

type UpdateRootConfig interface {
	Execute(ctx context.Context, in InputUpdateRootConfig) (*domain.RootConfig, error)
}

func NewUpdateRootConfig(r settingRepo.RootConfigRepository, u userRepo.Repository) *UpdateRootConfigInteractor {
	return &UpdateRootConfigInteractor{rootRepo: r, userRepo: u}
}

func (uc *UpdateRootConfigInteractor) Execute(ctx context.Context, in InputUpdateRootConfig) (*domain.RootConfig, error) {
	ok, err := uc.userRepo.IsRoot(ctx, in.UserID)
	if err != nil {
		return nil, ErrDatabaseFailure
	}
	if !ok {
		return nil, ErrUnauthorized
	}

	cfg := &domain.RootConfig{
		EditingPermission:          in.EditingPermission,
		IsTestUserMode:             in.IsTestUserMode,
		IsEmailAuthenticationCheck: in.IsEmailAuthCheck,
		IsLineAuthentication:       in.IsLineAuth,
	}

	return uc.rootRepo.Upsert(ctx, cfg)
}
