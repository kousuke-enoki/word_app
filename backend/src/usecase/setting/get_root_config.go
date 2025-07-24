package settinguc

import (
	"context"
	"errors"

	"word_app/backend/src/domain"
	settingRepo "word_app/backend/src/infrastructure/repository/setting"
	userRepo "word_app/backend/src/infrastructure/repository/user"
)

var (
	ErrRootConfigMissing = errors.New("root-config not found")
)

type InputGetRootConfig struct {
	UserID int
}

type OutputGetRootConfig struct {
	Config *domain.RootConfig
}

type GetRootConfig interface {
	Execute(ctx context.Context, in InputGetRootConfig) (*OutputGetRootConfig, error)
}

type GetRootConfigInteractor struct {
	userRepo       userRepo.Repository
	rootConfigRepo settingRepo.RootConfigRepository
}

func NewGetRootConfig(u userRepo.Repository, r settingRepo.RootConfigRepository) *GetRootConfigInteractor {
	return &GetRootConfigInteractor{userRepo: u, rootConfigRepo: r}
}

func (uc *GetRootConfigInteractor) Execute(ctx context.Context, in InputGetRootConfig) (*OutputGetRootConfig, error) {
	user, err := uc.userRepo.FindByID(ctx, in.UserID)
	if err != nil {
		return nil, err // ← DB エラーなどはそのまま
	}
	if !user.IsRoot {
		return nil, ErrUnauthorized
	}

	cfg, err := uc.rootConfigRepo.Get(ctx)
	if err != nil {
		return nil, ErrRootConfigMissing
	}
	return &OutputGetRootConfig{Config: cfg}, nil
}
