// usecase/setting/update_root_config.go
package settingUc

import (
	"context"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/interfaces/repository/setting"
	userRepo "word_app/backend/src/interfaces/repository/user"
)

type UpdateRootConfigInput struct {
	UserID            int
	EditingPermission string
	IsTestUserMode    bool
	IsEmailAuthCheck  bool
	IsLineAuth        bool
}

type UpdateRootConfigUsecase struct {
	rootRepo settingRepo.RootConfigRepository
	userRepo userRepo.UserRepository
}

func NewUpdateRootConfig(r settingRepo.RootConfigRepository, u userRepo.UserRepository) *UpdateRootConfigUsecase {
	return &UpdateRootConfigUsecase{rootRepo: r, userRepo: u}
}

func (uc *UpdateRootConfigUsecase) UpdateRootConfigExecute(ctx context.Context, in UpdateRootConfigInput) (*domain.RootConfig, error) {
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
